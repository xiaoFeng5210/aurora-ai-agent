# Docker Compose 单机部署说明

这套部署用于单机生产或准生产环境，包含：

- `frontend`: Nginx 托管 `frontend/web` 构建产物，并代理 `/api/` 到后端
- `admin`: Nginx 托管 `frontend/admin` 构建产物，并代理 `/api/` 到后端
- `backend`: Go 后端服务
- `postgres`: 业务数据库
- `redis`: 缓存和百度网盘 token 存储
- `rabbitmq`: RAG 向量化任务队列

Qdrant 使用外部服务，不在本 compose 中启动。

## 1. 准备配置

复制外部平台密钥配置：

```bash
cp .env.example .env
```

`.env` 只放当前代码仍通过环境变量读取的外部平台配置：

```bash
GLM_API_KEY=...
BAIDU_NETWORKDISK_CLIENT_ID=...
BAIDU_NETWORKDISK_DEVICE_ID=...
BAIDU_NETWORKDISK_CLIENT_SECRET=...
QDRANT_HOST=...
QDRANT_API_KEY=...
```

后端业务配置从 `/app/conf/*.yaml` 读取，也就是宿主机的 `deploy/conf/*.yaml`。这些文件已经放在 `deploy/conf/` 下，结构和根目录 `conf/` 保持一致。

## 2. 配置对应关系

`docker-compose.yml` 里会初始化 Postgres、Redis、RabbitMQ。它们的账号密码必须和 `deploy/conf/*.yaml` 保持一致。

默认 Postgres：

```yaml
postgres:
  host: postgres
  port: 5432
  user: aurora
  password: aurora_postgres_change_me
  dbname: aurora_ai_agent
```

对应 compose：

```yaml
POSTGRES_DB: aurora_ai_agent
POSTGRES_USER: aurora
POSTGRES_PASSWORD: aurora_postgres_change_me
```

默认 Redis：

```yaml
redis:
  host: redis
  port: 6379
  username: ""
  password: aurora_redis_change_me
  db: 0
```

对应 compose 的 Redis 启动参数：

```yaml
--requirepass aurora_redis_change_me
```

默认 RabbitMQ：

```yaml
rabbitmq:
  url: amqp://aurora:aurora_rabbitmq_change_me@rabbitmq:5672/
```

对应 compose：

```yaml
RABBITMQ_DEFAULT_USER: aurora
RABBITMQ_DEFAULT_PASS: aurora_rabbitmq_change_me
```

如果你改了任意一边，另一边也必须同步修改。

## 3. 持久化目录

compose 使用宿主机目录保存数据和日志：

```text
deploy/data/postgres       Postgres 数据
deploy/data/redis          Redis AOF/RDB 数据
deploy/data/rabbitmq       RabbitMQ 数据
deploy/data/backend/log    后端 zap/db 日志
deploy/data/nginx/log      Nginx access/error 日志
deploy/data/admin-nginx/log Admin Nginx access/error 日志
```

这些目录不要删除。删除后对应服务数据会丢失。

Postgres 的 `assets/sql/table.sql` 只会在 `deploy/data/postgres` 为空、数据库首次初始化时自动执行。后续升级脚本不会自动执行，需要手动处理。

## 4. 启动

先检查 compose 配置：

```bash
docker compose config
```

构建并启动：

```bash
docker compose up -d --build
```

查看状态：

```bash
docker compose ps
```

查看后端日志：

```bash
docker compose logs -f backend
```

## 5. 端口说明

Compose 默认只把前端入口暴露到宿主机：

```text
frontend  宿主机 5019  -> 容器 8080
admin     宿主机 50190 -> 容器 8080
rabbitmq  宿主机 127.0.0.1:15672 -> 容器 15672
```

`backend` 在 Docker 内部网络监听 `1119`，没有直接暴露到宿主机；`frontend` 和 `admin` 都通过服务名 `backend:1119` 访问它。

如果宿主机端口冲突，只改冒号左边的宿主机端口。例如 admin 改成 `50191`：

```yaml
admin:
  ports:
    - "50191:8080"
```

访问地址也相应变成 `http://服务器IP:50191`。

## 6. 验证

检查 Web 前端 Nginx：

```bash
curl http://localhost:5019/healthz
```

检查 Admin 前端 Nginx：

```bash
curl http://localhost:50190/healthz
```

检查后端是否通过 Nginx 可访问：

```bash
curl http://localhost:5019/ping
curl http://localhost:50190/ping
```

浏览器访问：

```text
Web 前端:   http://localhost:5019/
Admin 前端: http://localhost:50190/
```

RabbitMQ 管理后台默认只绑定本机：

```text
http://127.0.0.1:15672
```

账号密码见 `docker-compose.yml`，生产环境建议替换。

## 7. 常用运维命令

重启全部服务：

```bash
docker compose restart
```

只重启后端：

```bash
docker compose restart backend
```

停止服务但保留数据：

```bash
docker compose down
```

升级代码后重新构建：

```bash
git pull
docker compose up -d --build
```

查看容器资源：

```bash
docker stats
```

## 8. 备份和恢复

备份 Postgres：

```bash
docker compose exec postgres pg_dump -U aurora -d aurora_ai_agent > aurora_ai_agent.sql
```

恢复 Postgres：

```bash
cat aurora_ai_agent.sql | docker compose exec -T postgres psql -U aurora -d aurora_ai_agent
```

Redis 使用 AOF 持久化，备份时可以先触发保存：

```bash
docker compose exec redis redis-cli -a aurora_redis_change_me BGSAVE
```

然后备份：

```text
deploy/data/redis
```

RabbitMQ 备份：

```text
deploy/data/rabbitmq
```

日志备份：

```text
deploy/data/backend/log
deploy/data/nginx/log
```

## 9. 安全注意

- `.env` 不要提交。
- `deploy/conf/*.yaml` 已按 Docker Compose 部署提供默认值；如果你在服务器上改成真实生产密钥，注意不要把敏感值提交到公开仓库。
- 真实 API key 不要写入 README、issue、commit message 或文档。
- 如果密钥已经暴露，建议在对应平台立刻轮换。
- 生产环境必须替换示例密码和 JWT secret。
- 不建议把 Postgres、Redis、RabbitMQ 的业务端口直接暴露到公网。

## 10. 常见问题

后端启动失败并提示数据库连接失败：

- 检查 `deploy/conf/db.yaml` 的账号密码。
- 检查它是否和 `docker-compose.yml` 里的 Postgres 初始化配置一致。
- 如果 Postgres 已经初始化过，改 compose 里的 `POSTGRES_PASSWORD` 不会自动修改旧数据库密码。

Redis 认证失败：

- 检查 `deploy/conf/redis.yaml` 的 `password`。
- 检查它是否和 Redis `--requirepass` 参数一致。

RabbitMQ 连接失败：

- 检查 `deploy/conf/rabbitmq.yaml` 的 URL。
- 检查账号密码是否和 `RABBITMQ_DEFAULT_USER` / `RABBITMQ_DEFAULT_PASS` 一致。

修改 `assets/sql/table.sql` 后不生效：

- 这是预期行为。Postgres 初始化脚本只在数据目录为空时执行。
- 已有数据库需要手动执行升级 SQL。
