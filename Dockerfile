# Dockerfile
# 运行镜像
FROM alpine:latest

WORKDIR /app

# 安装基本工具
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 复制编译好的应用
COPY --from=builder /app/new-api-proxy .
COPY --from=builder /app/config/config.yaml ./config/

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./new-api-proxy"]