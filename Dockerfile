FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
ENV GOPROXY=https://goproxy.cn,direct \
    GO111MODULE=on \
    CGO_ENABLED=0
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 运行阶段
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .

# 安装等待依赖服务的工具
RUN apk add --no-cache bash
COPY wait-for.sh /app/wait-for.sh
RUN chmod +x /app/wait-for.sh

# 初始化脚本
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

EXPOSE 8080
CMD ["/app/entrypoint.sh"]