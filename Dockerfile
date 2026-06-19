# ------------------ 阶段 1: 编译前端 React SPA ------------------
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build

# ------------------ 阶段 2: 编译后端 Go ------------------
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app/backend

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o bitlog-server .

# ------------------ 阶段 3: 运行容器 ------------------
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 创建持久化数据与附件文件夹
RUN mkdir -p /app/data/uploads

# 拷贝后端可执行程序
COPY --from=backend-builder /app/backend/bitlog-server .

# 拷贝前端静态编译包到 Go 托管路径
COPY --from=frontend-builder /app/frontend/dist ./dist

# 默认环境变量
ENV PORT=3000
ENV GIN_MODE=release

EXPOSE 3000

# 启动 Go Web 服务 (其内部会自动进行 SQLite 建表迁移并初始化默认管理员)
CMD ["./bitlog-server"]
