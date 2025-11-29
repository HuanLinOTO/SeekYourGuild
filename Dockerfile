# 多阶段构建

# 阶段1: 构建前端
FROM node:20-alpine AS frontend-builder

WORKDIR /app

# 复制依赖文件
COPY frontend/package.json frontend/yarn.lock ./
RUN yarn install --frozen-lockfile

# 复制源码并构建
COPY frontend/ ./
RUN yarn build

# 阶段2: 构建后端（包含嵌入的前端）
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app

# 复制依赖文件
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# 复制源码
COPY backend/ ./

# 复制前端构建产物到静态目录
COPY --from=frontend-builder /app/dist ./internal/static/dist

# 构建后端（纯 Go，无需 CGO）
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /server ./cmd/server/

# 阶段3: 最终镜像
FROM alpine:3.19

WORKDIR /app

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 从构建阶段复制文件（只需要一个二进制文件）
COPY --from=backend-builder /server ./server

# 设置时区
ENV TZ=Asia/Shanghai

# 环境变量
ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=8080
ENV DB_DRIVER=sqlite
ENV DB_DSN=file:/data/seekyourguild.db?cache=shared&mode=rwc

# 创建数据目录
RUN mkdir -p /data

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/groups || exit 1

# 启动命令
CMD ["./server"]
