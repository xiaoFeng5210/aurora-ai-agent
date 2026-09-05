# Ariadne

> 新的产品名称：Ariadne


## 开发
### 启动
```bash
air
```




## 部署
### 后端
docker build -f Dockerfile.backend -t ariadne-backend:v0.0.1 .

docker run -d --name ariadne-backend \
  -p 1119:1119 \
  -v "$(pwd)/log:/log" \
  --env-file .env \
  ariadne-backend:v0.0.1

### 前端
镜像采用多阶段构建：`node:22-alpine` 里 `pnpm install` + `pnpm build`，再把 `dist/` 拷到 `nginx:1.27-alpine` 由 nginx 渲染。

构建上下文必须是仓库根目录（Dockerfile 里需要访问 `frontend/web/`）。nginx 配置统一放在 `frontend/web/nginx/`，整个目录会被 `COPY` 到容器的 `/etc/nginx/conf.d/`，以后新增 `xxx.conf` 直接丢进去即可。

```bash
# 构建（根目录执行）
docker build -f Dockerfile.frontend -t ariadne-frontend:v0.0.1 .

# 本地运行验证
docker run -d --name ariadne-frontend -p 8080:8080 ariadne-frontend:v0.0.1
# 访问 http://localhost:8080
```

打 tag 并推送到阿里云 ACR（把变量替换成你自己的命名空间/区域）：

```bash
export ACR_REGISTRY=registry.cn-hangzhou.aliyuncs.com
export ACR_NAMESPACE=<your-namespace>
export IMAGE_TAG=v0.0.1

docker login --username=<your-acr-username> ${ACR_REGISTRY}

docker tag ariadne-frontend:${IMAGE_TAG} \
  ${ACR_REGISTRY}/${ACR_NAMESPACE}/ariadne-frontend:${IMAGE_TAG}

docker push ${ACR_REGISTRY}/${ACR_NAMESPACE}/ariadne-frontend:${IMAGE_TAG}
```

如需多架构镜像（ACK 节点常见 amd64，本地 Mac 是 arm64），用 buildx：

```bash
docker buildx create --use --name ariadne-builder 2>/dev/null || true

docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -f Dockerfile.frontend \
  -t ${ACR_REGISTRY}/${ACR_NAMESPACE}/ariadne-frontend:${IMAGE_TAG} \
  --push .
```

> nginx 配置在 `frontend/web/nginx/default.conf`：容器内监听 **8080**（非特权端口，便于 K8s 以非 root 用户跑），已开启 gzip、SPA 路由 fallback（`try_files ... /index.html`），`/healthz` 可用作 ACK 探针。
>
> ACK 部署时 `containerPort` 填 `8080`，`Service.targetPort` 也填 `8080`，对外的 `Service.port` / SLB 监听端口（80、443 等）随意。



## 回答模型：DeepSeek V4 Flash

回答与分析统一使用 `deepseek-v4-flash`，通过官方 `https://api.deepseek.com/chat/completions` 调用。
在本地 `.env` 或部署环境配置 `DEEPSEEK_API_KEY`，不要将真实 Key 写入源码。
Docker Compose 已通过 `env_file: .env` 加载；更新后需重建并重启 backend：

```bash
docker compose up -d --build backend
```

`GLM_API_KEY` 暂时仅供原有 `embedding-3` 向量化使用，仍须保留。向量数据库、维度和存量向量均未迁移；不能直接把 DeepSeek Key 填到向量化配置中。

- `ai/llm/model.go` 定义供应商无关的 `Model`、`ChatOptions` 和 `ChatResult`；agent 只依赖该接口。
- `ai/llm/deepseek.go` 封装请求映射、HTTP 超时/取消、SSE 解码、工具分片合并、思考上下文回传及错误处理。工具定义由 agent 注入。
- 默认 `thinking.type=disabled`，沿用原有行为；显式 `enabled` 可开启思考。思考内容仅在当前工具循环内部传递，不展示或保存到聊天历史。历史普通回答以空思考字段兼容回传。
- SSE 保留 `start / delta / tool_call / tool_result / done / error`；截断、内容过滤、资源不足、连接提前关闭不会伪装成正常完成，也不会执行未完整生成的工具参数。
- 新接口：`POST /api/v1/chat/stream/:document_id`。旧 `/api/v1/chat/glm/stream/:document_id` 作为兼容别名，也实际使用 DeepSeek。

验证：

```bash
# 无外部模型请求的适配层与 agent 回归测试
go test ./ai/llm ./ai/agent
# 显式真实 API 测试：普通/思考模式各执行一次工具调用和结果回答；会产生少量费用
DEEPSEEK_LIVE_TEST=1 go test ./ai/llm -run '^TestDeepSeekLive$' -v -count=1
```

官方参考：[模型与价格](https://api-docs.deepseek.com/quick_start/pricing/)、[Chat Completions](https://api-docs.deepseek.com/api/create-chat-completion/)、[思考与工具调用上下文](https://api-docs.deepseek.com/guides/thinking_mode/)。
