# InfoWeaver Backend

InfoWeaver 后端服务 - 一个基于 Go 语言构建的多模态 RAG（Retrieval-Augmented Generation）系统。

## 📖 功能特性

- 📁 文件上传与管理（存储于 MinIO）
- 📊 数据集管理（基于所有权的访问控制）
- 🔐 用户认证与授权（JWT）
- 🔍 向量检索（Milvus）
- 📦 消息队列（RabbitMQ）
- 🗄️ 多数据库支持（PostgreSQL、Redis、MinIO、Milvus）

## 🏗️ 项目架构

```txt
server/
├── api/           # API 层，REST 端点定义
│   └── v1/        # v1 版本 API
├── cmd/           # 应用入口
├── config/        # 配置管理（Viper）
├── db/            # 数据库初始化（PostgreSQL、Redis、MinIO、RabbitMQ、Milvus）
├── docs/          # Swagger 文档
├── middleware/    # 中间件（JWT、CORS、日志、错误处理）
├── models/        # 数据模型（GORM 实体）
├── service/       # 业务逻辑层
├── tests/         # 测试文件
└── utils/         # 工具函数
```

## 🛠️ 技术栈

- **Web 框架**: [Echo v5](https://github.com/labstack/echo)
- **ORM**: [GORM](https://gorm.io/)
- **配置管理**: [Viper](https://github.com/spf13/viper)
- **日志**: [Zap](https://go.uber.org/zap)
- **数据库**:
  - PostgreSQL - 元数据存储
  - Redis - 缓存与会话管理
  - MinIO - 对象存储
  - Milvus - 向量数据库
  - RabbitMQ - 消息队列
- **API 文档**: [Swagger](https://github.com/swaggo/swag)

## 📋 环境要求

- Go 1.21+
- PostgreSQL 12+
- Redis 6+
- MinIO
- Milvus 2.x
- RabbitMQ 3.x

## 🚀 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/AAAkater/InfoWeaver.git
cd InfoWeaver/server
```

### 2. 配置环境变量

复制环境变量示例文件并根据需要修改：

```bash
cp .env.example .env.local
```

### 3. 启动服务

```bash
go run cmd/main.go
```

服务将在配置的端口启动（默认 `8080`）。

### 4. 访问 API 文档

启动服务后，访问 Swagger 文档：

```
http://localhost:8080/swagger/index.html
```

## 📝 配置说明

主要配置项（`.env.local`）：

| 配置项               | 说明            | 默认值       |
| -------------------- | --------------- | ------------ |
| `SYSTEM_SERVER_PORT` | 服务端口        | `8080`       |
| `POSTGRES_HOST`      | PostgreSQL 主机 | `localhost`  |
| `POSTGRES_PORT`      | PostgreSQL 端口 | `5432`       |
| `POSTGRES_DB`        | 数据库名称      | `InfoWeaver` |
| `REDIS_HOST`         | Redis 主机      | `localhost`  |
| `REDIS_PORT`         | Redis 端口      | `6379`       |
| `MINIO_HOST`         | MinIO 主机      | `localhost`  |
| `MINIO_PORT`         | MinIO 端口      | `9000`       |
| `MILVUS_HOST`        | Milvus 主机     | `localhost`  |
| `MILVUS_PORT`        | Milvus 端口     | `19530`      |
| `JWT_SIGNING_KEY`    | JWT 签名密钥    | `KFCvME50`   |
| `JWT_EXPIRES_TIME`   | JWT 过期时间    | `7d`         |

## 🧪 测试

```bash
go test -v ./tests
```

## 📚 Swagger 文档生成

当 API 发生变更时，重新生成 Swagger 文档：

```bash
swag init -g ./cmd/main.go -o ./docs
```

## 🔒 安全特性

- JWT 身份认证
- 基于所有权的访问控制
- 所有文件/数据集操作均验证用户所有权

## 📄 License

MIT License
