# Multi-stage Dockerfile for Laravel Backend (Optimized for minimal size)
# Build context = backend/ directory
#
# Optimization techniques:
# 1. --no-dev: Excludes PHPUnit/PHPStan/etc from production image (~80MB saved)
# 2. --virtual .build-deps: Dev headers removed after compilation (~100MB saved)
# 3. Runtime libs only: libpng instead of libpng-dev, etc.
# 4. No git/unzip/composer in final image (build-time only)
# 5. Parallel compilation with -j$(nproc)
# 6. Aggressive cache cleanup

# ============================================================
# Stage 1: Install Composer dependencies (production only)
# ============================================================
# Use the same PHP 8.3 CLI variant as the runtime stage to ensure PHP version
# parity. Composer is copied as a static binary from the official composer image.
FROM php:8.3-cli-alpine AS composer

# Install the composer binary from the official composer image and unzip (needed
# to extract dist archives). unzip is removed before the stage ends.
RUN apk add --no-cache unzip \
    && rm -rf /var/cache/apk/*
COPY --from=composer:2.8 /usr/bin/composer /usr/bin/composer

WORKDIR /app

COPY composer.json composer.lock* ./

# Install ONLY production dependencies (no dev deps = no PHPUnit/PHPStan/Faker)
# --no-scripts: artisan not available in this stage
# --ignore-platform-reqs: extensions built in next stage, skip checks here
RUN composer install \
        --no-interaction \
        --no-progress \
        --prefer-dist \
        --optimize-autoloader \
        --no-dev \
        --no-scripts \
        --ignore-platform-reqs \
    && composer dump-autoload --optimize --no-scripts --no-dev \
    && rm -rf /root/.composer/cache

# ============================================================
# Stage 2: Production runtime image (minimal)
# ============================================================
FROM php:8.3-fpm-alpine AS app

# Install build deps as virtual package, compile PHP extensions,
# then remove build deps. Only runtime shared libraries remain.
RUN apk add --no-cache --virtual .build-deps \
        $PHPIZE_DEPS \
        postgresql-dev \
        libpng-dev \
        libjpeg-turbo-dev \
        freetype-dev \
        libzip-dev \
        oniguruma-dev \
    && apk add --no-cache \
        nginx \
        supervisor \
        curl \
        libzip \
        libpng \
        libjpeg-turbo \
        freetype \
        oniguruma \
        postgresql-libs \
    && docker-php-ext-configure gd --with-freetype --with-jpeg \
    && docker-php-ext-install -j"$(nproc)" \
        pdo_mysql \
        pdo_pgsql \
        zip \
        gd \
        mbstring \
        pcntl \
    && pecl install redis \
    && docker-php-ext-enable redis \
    && apk del --no-cache .build-deps \
    && rm -rf /var/cache/apk/* /tmp/* /var/tmp/* /usr/src/php*

# Set working directory
WORKDIR /var/www/html

# Copy production vendor (no dev dependencies)
COPY --from=composer /app/vendor ./vendor

# Copy application code (context = backend/)
COPY . .

# Copy configuration files
COPY docker/nginx/default.conf /etc/nginx/http.d/default.conf
COPY docker/supervisor/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY docker/scripts/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

# Ensure storage directories exist with proper permissions
RUN mkdir -p storage/app/public \
    storage/framework/cache/data \
    storage/framework/sessions \
    storage/framework/views \
    storage/logs \
    bootstrap/cache \
    && chown -R www-data:www-data /var/www/html \
    && chmod -R 755 /var/www/html/storage \
    && chmod -R 755 /var/www/html/bootstrap/cache

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Expose port
EXPOSE 8080

# Run entrypoint (starts nginx + php-fpm + queue worker via supervisor)
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["supervisord", "-n", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
