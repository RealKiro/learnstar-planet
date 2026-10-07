# Multi-stage Dockerfile for LearnStar Planet（Go 后端）
#
# Stage 1: Node 构建前端（Vue 3 + Vite）
# Stage 2: Go 编译（CGO_ENABLED=0，纯 Go SQLite 驱动，无 cgo 依赖）
# Stage 3: alpine 运行时（ca-certificates 供第三方 HTTPS 调用；tzdata 兜底，
#          二进制内已 embed time/tzdata 双保险）
#
# Location: docker/Dockerfile
# Build context = learnstar-planet/（仓库根）

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
# Stage 2: Build Go binary
# ============================================================
FROM golang:1.27-alpine AS go-builder

WORKDIR /src

# 依赖独立成层：go.mod/go.sum 不变时命中缓存
COPY backend-go/go.mod backend-go/go.sum ./
RUN go mod download

COPY backend-go/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/learnstar-go .

# ============================================================
# Stage 3: Runtime
# ============================================================
FROM alpine:3.21 AS app

RUN apk add --no-cache ca-certificates tzdata && rm -rf /var/cache/apk/*

WORKDIR /app

COPY --from=go-builder /out/learnstar-go ./learnstar-go

# 前端构建产物：Go 端 PUBLIC_DIR=public 由它服务（SPA + index.html 兜底）
COPY --from=frontend-builder /app/dist ./public

# 运行时目录：data（SQLite + .jwt_secret）/ uploads（学校 LOGO）/ logs
RUN mkdir -p data storage/app/uploads logs

COPY docker/scripts/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

# alpine busybox 自带 wget，无需额外安装 curl
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -q -O /dev/null http://localhost:8080/health || exit 1

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
