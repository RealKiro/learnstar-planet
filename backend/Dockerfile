# Multi-stage Dockerfile for Laravel Backend
# Stage 1: Dependencies
FROM composer:2.6 AS composer

WORKDIR /app

COPY backend/composer.json backend/composer.lock* ./

# Install PHP extensions required by dependencies
RUN apk add --no-cache $PHPIZE_DEPS \
    && docker-php-ext-install pdo_mysql pdo_pgsql redis

# Install Composer dependencies
RUN composer install --no-interaction --no-progress --prefer-dist --optimize-autoloader \
    && composer dump-autoload --optimize

# Stage 2: Build frontend assets (if using Mix/Vite)
FROM node:20-alpine AS frontend

WORKDIR /app

COPY frontend/package*.json ./
RUN npm ci --only=production

COPY frontend/ .
RUN npm run build

# Stage 3: Production image
FROM php:8.2-fpm-alpine AS app

# Install system dependencies
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
    && docker-php-ext-install \
        pdo_mysql \
        pdo_pgsql \
        redis \
        zip \
        gd \
        mbstring \
    && rm -rf /var/cache/apk/*

# Install Composer
COPY --from=composer /usr/bin/composer /usr/bin/composer

# Set working directory
WORKDIR /var/www/html

# Copy dependencies from composer stage
COPY --from=composer /app/vendor ./vendor

# Copy application code
COPY backend/ .

# Copy frontend assets (if built)
COPY --from=frontend /app/public/build ./public/build

# Set proper permissions
RUN chown -R www-data:www-data /var/www/html \
    && chmod -R 755 /var/www/html/storage \
    && chmod -R 755 /var/www/html/bootstrap/cache

# Install supervisord for running queue workers
RUN apk add --no-cache supervisor

# Copy supervisor configuration
COPY docker/supervisor/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Copy entrypoint script
COPY docker/scripts/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:9000/health || exit 1

# Expose port
EXPOSE 9000

# Run entrypoint
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

# Start PHP-FPM
CMD ["php-fpm"]
