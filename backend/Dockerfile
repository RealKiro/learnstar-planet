# Multi-stage Dockerfile for Laravel Backend
# Build context = backend/ directory

# Stage 1: Install Composer dependencies
FROM composer:2.8 AS composer

WORKDIR /app

COPY composer.json composer.lock* ./

# Install PHP extensions (gd for intervention/image, redis via PECL, pdo_pgsql needs postgresql-dev, pcntl is core ext)
RUN apk add --no-cache $PHPIZE_DEPS postgresql-dev libpng-dev libjpeg-turbo-dev freetype-dev \
    && docker-php-ext-configure gd --with-freetype --with-jpeg \
    && docker-php-ext-install gd pdo_mysql pdo_pgsql pcntl \
    && pecl install redis \
    && docker-php-ext-enable redis

# Install Composer dependencies (--no-scripts: artisan not available in this stage)
RUN composer install --no-interaction --no-progress --prefer-dist --optimize-autoloader --no-scripts \
    && composer dump-autoload --optimize

# Stage 2: Production image
FROM php:8.3-fpm-alpine AS app

# Install system dependencies (postgresql-dev needed for pdo_pgsql)
RUN apk add --no-cache \
    git \
    curl \
    zip \
    unzip \
    libzip-dev \
    libpng-dev \
    libjpeg-turbo-dev \
    freetype-dev \
    oniguruma-dev \
    postgresql-dev \
    nginx \
    supervisor \
    && docker-php-ext-configure gd --with-freetype --with-jpeg \
    && docker-php-ext-install \
        pdo_mysql \
        pdo_pgsql \
        zip \
        gd \
        mbstring \
        pcntl \
    && pecl install redis \
    && docker-php-ext-enable redis \
    && rm -rf /var/cache/apk/*

# Install Composer
COPY --from=composer /usr/bin/composer /usr/bin/composer

# Set working directory
WORKDIR /var/www/html

# Copy dependencies from composer stage
COPY --from=composer /app/vendor ./vendor

# Copy application code (context = backend/)
COPY . .

# Copy Nginx config
COPY docker/nginx/default.conf /etc/nginx/http.d/default.conf

# Copy supervisor configuration
COPY docker/supervisor/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Copy entrypoint script
COPY docker/scripts/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

# Set proper permissions and ensure storage directories exist
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
