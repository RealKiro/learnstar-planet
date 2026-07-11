# Multi-stage Dockerfile for LearnStar Planet Backend + Vue Frontend
# Build context = learnstar-planet/ (repository root)
#
# Optimization techniques:
# 1. --no-dev: Excludes PHPUnit/PHPStan/etc from production image
# 2. --virtual .build-deps: Dev headers removed after compilation
# 3. Runtime libs only: libpng instead of libpng-dev
# 4. Node stage builds Vue SPA, output copied to public/
# 5. No git/unzip/composer/node in final image (build-time only)
# 6. Aggressive cache cleanup

# ============================================================
# Stage 1: Build Vue frontend (Node)
# ============================================================
FROM node:22-alpine AS frontend-builder

WORKDIR /app

# Copy package files and install dependencies with npm
COPY frontend-vue/package.json frontend-vue/package-lock.json* ./
RUN if [ -f package-lock.json ]; then \
        npm ci; \
    else \
        npm install; \
    fi

# Copy source and build
COPY frontend-vue/ ./
RUN npm run build

# ============================================================
# Stage 2: Install Composer dependencies (production only)
# ============================================================
FROM php:8.3-cli-alpine AS composer

RUN apk add --no-cache unzip \
    && rm -rf /var/cache/apk/*
COPY --from=composer:2.8 /usr/bin/composer /usr/bin/composer

WORKDIR /app

COPY backend/composer.json backend/composer.lock* ./

# Install ONLY production dependencies
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
# Stage 3: Production runtime image (minimal)
# ============================================================
FROM php:8.3-fpm-alpine AS app

# Install build deps as virtual package, compile PHP extensions,
# then remove build deps.
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

# Copy backend application code
COPY backend/ .

# Copy Vue frontend build output to public directory
COPY --from=frontend-builder /app/dist ./public

# Regenerate autoload after code copy
RUN composer dump-autoload --optimize --no-scripts --no-dev 2>/dev/null || true

# Copy configuration files
COPY backend/docker/nginx/default.conf /etc/nginx/http.d/default.conf
COPY backend/docker/supervisor/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY backend/docker/scripts/entrypoint.sh /usr/local/bin/entrypoint.sh
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

# Allow Nginx to serve public/ (ensure read access)
RUN chmod -R 755 /var/www/html/public

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Expose port
EXPOSE 8080

# Run entrypoint (starts nginx + php-fpm + queue worker via supervisor)
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
C