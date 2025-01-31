# 使用官方的 Go 基础镜像
FROM golang:1.23-alpine as builder

# 设置工作目录
WORKDIR /app

# 复制项目文件到工作目录
COPY . .

# 下载依赖
RUN go mod tidy

# 构建可执行文件
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# 使用轻量级的基础镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /app

COPY .env .

# 从 builder 阶段复制可执行文件
COPY --from=builder /app/main .

# 暴露应用程序的端口
EXPOSE 8080

# 运行应用程序
CMD ["./main"]
