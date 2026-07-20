# Multi-stage Dockerfile for LearnStar Planet
# Stage 1: Build Vue frontend (Node)
# Stage 2: PHP runtime with RoadRunner (no Nginx, no FPM, no Supervisor)
#
# Build context = learnstar-planet/ (repository root)

# ============================================================
# Stage 1: Build Vue frontend
# ============================================================
FROM node:22-alpine AS frontend-builder

WORKDIR /app

COPY frontend-vue/package.json frontend-vue/package-lock.json* ./
RUN if [ -f package-lock.json ]; then npm ci; else npm install; fi

COPY frontend-vue/ ./
RUN npm run build:ci

# ============================================================
# Stage 2: PHP + RoadRunner production runtime
# ============================================================
FROM php:8.3-cli-alpine AS app

# Install system dependencies for PHP extensions
RUN apk add --no-cache \
    libpng libjpeg-turbo freetype libzip oniguruma \
    postgresql-libs libwebp curl unzip \
    && rm -rf /var/cache/apk/*

# Compile PHP extensions
RUN apk add --no-cache --virtual .build-deps \
        $PHPIZE_DEPS \
        postgresql-dev libpng-dev libjpeg-turbo-dev \
        freetype-dev libzip-dev oniguruma-dev libwebp-dev \
    && docker-php-ext-configure gd --with-freetype --with-jpeg --with-webp \
    && docker-php-ext-install -j"$(nproc)" \
        pdo_mysql pdo_pgsql zip gd mbstring pcntl \
    && pecl install redis \
    && docker-php-ext-enable redis \
    && apk del --no-cache .build-deps \
    && rm -rf /var/cache/apk/* /tmp/* /var/tmp/* /usr/src/php*

# Install Composer
COPY --from=composer:2.8 /usr/bin/composer /usr/bin/composer

WORKDIR /app

# Copy composer files and install production dependencies
COPY backend/composer.json backend/composer.lock* ./
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

# Copy application code (需要在 octane:install 之前)
COPY backend/ .

# 创建必要目录
RUN mkdir -p bootstrap/cache storage/framework/cache/data \
    storage/framework/sessions storage/framework/views storage/logs

# Install RoadRunner via Octane
RUN php artisan package:discover --ansi 2>/dev/null || true \
    && php artisan octane:install --server=roadrunner --no-interaction

# Copy Vue frontend build to public directory
COPY --from=frontend-builder /app/dist ./public

# Copy RoadRunner config
COPY backend/docker/roadrunner/.rr.yaml ./

# Ensure storage directories exist with proper permissions
RUN mkdir -p storage/app/public \
    storage/framework/cache/data \
    storage/framework/sessions \
    storage/framework/views \
    storage/logs \
    bootstrap/cache \
    && chown -R www-data:www-data /var/www/html 2>/dev/null || true \
    && chmod -R 755 storage bootstrap/cache public

# Copy entrypoint script
COPY backend/docker/scripts/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
