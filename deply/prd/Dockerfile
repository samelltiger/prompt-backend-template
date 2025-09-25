# 使用Alpine作为基础镜像
FROM alpine:latest

# 安装必要的依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 创建应用用户
RUN addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -s /bin/sh -D appuser

WORKDIR /app

# 复制编译好的应用程序
COPY nano-banana-prompt-api .

# 复制配置文件夹
COPY config/ ./config/

# 复制web文件夹
COPY web/ ./web/

# 设置文件权限
RUN chown -R appuser:appgroup /app && \
    chmod +x /app/nano-banana-prompt-api

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 9701

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:9701/health || exit 1

# 运行应用
CMD ["./nano-banana-prompt-api"]