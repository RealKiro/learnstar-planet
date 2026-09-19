# Multi-stage Dockerfile for LearnStar Planet
# Stage 1: Build Vue frontend (Node)
# Stage 2: Composer vendor (独立层：只在 composer.lock 变化时重建，且 composer 不进最终镜像)
# Stage 3: Runtime = Octane + FrankenPHP（常驻 worker，替代 artisan serve 开发服务器）
#
# Location: docker/Dockerfile (all Docker build files live in docker/)
# Build context = learnstar-planet/ (repository root)
#
# 为什么是 FrankenPHP：
#   - 单个静态二进制 = Web 服务器(Caddy) + PHP 运行时，无 nginx/FPM/Supervisor
#   - Octane worker 常驻内存，框架启动成本归零（PHP 唯一的架构劣势）
#   - alpine 变体 + install-php-extensions，镜像比旧的 php-cli + artisan serve 更小
#   - 官方镜像含 PHP 8.5（tag: 1-php8.5-alpine），与 config/platform.php 对齐

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
# Stage 2: Composer 依赖（vendor 独立成段，最终镜像不含 composer）
# ============================================================
FROM dunglas/frankenphp:1-php8.5-alpine AS vendor-builder

COPY --from=composer:2.8 /usr/bin/composer /usr/bin/composer

WORKDIR /app

COPY backend/composer.json backend/composer.lock* ./
RUN composer install \
        --no-interaction \
        --no-progress \
        --prefer-dist \
        --optimize-autoloader \
        --no-dev \
        --no-scripts \
        --ignore-platform-reqs

# ============================================================
# Stage 3: Runtime（Octane + FrankenPHP）
# ============================================================
FROM dunglas/frankenphp:1-php8.5-alpine AS app

# curl 仅供 HEALTHCHECK；其余 PHP 扩展用镜像自带的 install-php-extensions
# （脚本会自动装各扩展的 alpine 运行库，无需手工 apk 拉一串 -dev 包）
RUN apk add --no-cache curl && rm -rf /var/cache/apk/*

# 扩展清单与旧镜像对齐（gd 供 intervention/image 头像处理；pcntl 为 Octane 必需）：
RUN install-php-extensions \
        opcache \
        pdo_mysql \
        pdo_pgsql \
        pdo_sqlite \
        zip \
        gd \
        mbstring \
        pcntl \
        redis

WORKDIR /app

# vendor 来自独立段：业务代码改动不再打掉 composer install 层
COPY --from=vendor-builder /app/vendor ./vendor

COPY backend/ .

# Copy Vue frontend build to public directory
COPY --from=frontend-builder /app/dist ./public

# 创建必要目录、完成包发现、发布 Octane 配置
# （octane:install 若因构建期无 .env 失败不影响运行：entrypoint 启动时显式传 --server）
RUN mkdir -p bootstrap/cache storage/framework/cache/data \
    storage/framework/sessions storage/framework/views storage/logs \
    storage/app/public \
    && php artisan package:discover --ansi 2>/dev/null || true \
    && php artisan octane:install --server=frankenphp --no-interaction 2>/dev/null || true \
    && touch storage/database.sqlite && chmod 666 storage/database.sqlite \
    && chown -R www-data:www-data /app \
    && chmod -R 755 storage bootstrap/cache public

# Copy entrypoint script
COPY docker/scripts/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
