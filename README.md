# Prompt Backend Template

一个基于 Go 语言开发的 AI 大语言模型 API 后端服务模板，提供完整的 API 管理、用户认证、日志记录、计费统计等功能。

[访问网站主页](https://nano-banana.yueshu365.cn/)

## 项目简介

本项目是一个功能完整的 AI API 后端服务，主要特性包括：

- **OpenAI 兼容接口**：支持标准的 `/v1/chat/completions` 接口
- **多租户管理**：支持多个 API Key 的独立计费和日志管理
- **实时监控**：提供调用日志、费用统计、服务状态监控
- **管理员功能**：兑换码生成、用户配额管理、数据同步等
- **前端界面**：内置 Web 管理界面，支持用户自助查询

## 技术栈

- **后端框架**：Gin (Go)
- **数据库**：MySQL + GORM
- **缓存**：Redis
- **对象存储**：阿里云 OSS
- **监控**：Prometheus 指标
- **定时任务**：Cron 任务管理
- **日志**：Zap + Lumberjack

## 项目结构

```
prompt-backend-template/
├── cmd/                    # 主程序入口
│   └── main.go
├── internal/               # 内部业务逻辑
│   ├── api/               # API 路由和处理器
│   │   ├── admin/        # 管理员 API
│   │   ├── chat/         # 聊天接口
│   │   └── dashboard/    # 仪表板接口
│   ├── middleware/        # 中间件
│   ├── models/           # 数据模型
│   └── service/          # 业务服务层
├── pkg/                   # 公共库
│   ├── cache/            # 缓存组件
│   ├── cron/             # 定时任务
│   ├── database/         # 数据库连接
│   ├── logger/           # 日志组件
│   ├── oss/              # OSS 存储
│   ├── queue/            # 消息队列
│   └── util/             # 工具函数
├── config/               # 配置文件
│   ├── config.yaml      # 开发配置
│   └── config.prd.yaml  # 生产配置
├── web/                 # 静态文件（前端界面）
├── scripts/             # 脚本文件
├── docs/               # 文档
└── deploy/             # 部署文件
```

## 快速开始

### 环境要求

- Go 1.23.7+
- MySQL 5.7+
- Redis 6.0+

### 安装步骤

1. **克隆项目**
   ```bash
   git clone https://github.com/samelltiger/prompt-backend-template.git
   cd prompt-backend-template
   ```

2. **配置数据库**
   ```bash
   # 导入数据库脚本
   mysql -u root -p < scripts/prompts_tables.sql
   ```

3. **修改配置文件**
   ```bash
   cp config/config.yaml.example config/config.yaml
   # 编辑配置文件，设置数据库连接、Redis、OSS等参数
   ```

4. **安装依赖**
   ```bash
   go mod tidy
   ```

5. **启动服务**
   ```bash
   go run cmd/main.go
   ```

6. **访问界面**
   - 前端界面：http://localhost:9701
   - API 文档：http://localhost:9701/api/about
   - 监控指标：http://localhost:9701/metrics

## 其他教程/文章
[
3个月变永久！claude code给我写了一套脚本，帮我把网站HTTPS成本降为0](https://mp.weixin.qq.com/s/Yvp7UTLM4w75ILE-yAwIgA)

[Claude Code+Jenkins急速打造网站部署的上线系统](https://mp.weixin.qq.com/s/vjzw9K4TMwSnt-wel-prFg)

[“别再让AI网站丑下去了！”涛哥用Claude code实战MinIO主题移植，效果炸裂！完全去掉AI味！](https://mp.weixin.qq.com/s/3OGKZqAT7-i--Qhl19eu9w)

## API 接口

### 认证方式

所有受保护的接口都需要在 Header 中提供 API Key：
```
Authorization: Bearer <your-api-key>
```

### 主要接口

#### 1. 聊天接口
```http
POST /v1/chat/completions
Content-Type: application/json
Authorization: Bearer <api-key>

{
  "model": "deepseek-chat",
  "messages": [
    {"role": "user", "content": "Hello"}
  ]
}
```

#### 2. 账单查询
```http
GET /v1/dashboard/billing/subscription
GET /v1/dashboard/billing/usage
```

#### 3. 价格查询
```http
GET /api/pricing
```

#### 4. 兑换码管理
```http
GET /api/redeem          # 查询兑换码信息
POST /api/redeem         # 使用兑换码
```


## 配置说明

### 主要配置项

```yaml
server:
  port: 9701
  host: "0.0.0.0"

new_api:
  domain: "http://your-api-domain"
  admin_key: "your-admin-key"

database:
  gateway_dsn: "数据库连接字符串"
  new_api_dsn: "New API 数据库连接"

redis:
  addr: "redis-host:6379"
  password: "redis-password"
  db: 2

rate_limit:
  billing_query_limit: 10
  log_query_limit: 20

log:
  retention_days: 30

model_mapping:
  "deepseek-chat":
    - "deepseek-chat"
    - "hs-deepseek-v3-250324"
```

## 监控和日志

### 监控指标

服务内置 Prometheus 指标，可通过 `/metrics` 端点访问：

- HTTP 请求统计
- 数据库连接池状态
- Redis 连接状态
- 业务指标（调用次数、成功率等）

### 日志管理

- 日志级别：支持 debug、info、warn、error
- 日志轮转：自动按大小和时间轮转
- 日志保留：可配置保留天数

## 开发指南

### 添加新的 API 接口

1. 在 `internal/api/` 创建新的处理器
2. 在 `internal/service/` 实现业务逻辑
3. 在 `cmd/main.go` 中注册路由
4. 添加相应的中间件（如需要认证）

## 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查数据库服务是否启动
   - 验证连接字符串配置

2. **Redis 连接失败**
   - 检查 Redis 服务状态
   - 验证密码和数据库编号

3. **API Key 认证失败**
   - 检查 API Key 格式是否正确
   - 验证用户是否存在于数据库

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 联系方式

- 项目主页：[https://github.com/samelltiger/prompt-backend-template](https://github.com/samelltiger/prompt-backend-template)
- 微信：ithulianwang
- 公众号： 涛哥AI编程
- 网站：[https://nano-banana.yueshu365.cn/](https://nano-banana.yueshu365.cn/)

---

**提示**: 这是一个后端服务模板，实际使用时需要根据具体业务需求进行定制开发。