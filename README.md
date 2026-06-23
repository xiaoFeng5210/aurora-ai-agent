# aurora-ai-agent

> 新的产品名称：Ariadne


## 开发
### 启动
```bash
air
```




## 部署
### 后端
docker build -f Dockerfile.backend -t aurora-backend:v0.0.1 .

docker run -d --name aurora-backend \
  -p 1119:1119 \
  -v "$(pwd)/log:/log" \
  --env-file .env \
  aurora-backend:v0.0.1

### 前端
镜像采用多阶段构建：`node:22-alpine` 里 `pnpm install` + `pnpm build`，再把 `dist/` 拷到 `nginx:1.27-alpine` 由 nginx 渲染。

构建上下文必须是仓库根目录（Dockerfile 里需要访问 `frontend/web/`）。nginx 配置统一放在 `frontend/web/nginx/`，整个目录会被 `COPY` 到容器的 `/etc/nginx/conf.d/`，以后新增 `xxx.conf` 直接丢进去即可。

```bash
# 构建（根目录执行）
docker build -f Dockerfile.frontend -t aurora-frontend:v0.0.1 .

# 本地运行验证
docker run -d --name aurora-frontend -p 8080:8080 aurora-frontend:v0.0.1
# 访问 http://localhost:8080
```

打 tag 并推送到阿里云 ACR（把变量替换成你自己的命名空间/区域）：

```bash
export ACR_REGISTRY=registry.cn-hangzhou.aliyuncs.com
export ACR_NAMESPACE=<your-namespace>
export IMAGE_TAG=v0.0.1

docker login --username=<your-acr-username> ${ACR_REGISTRY}

docker tag aurora-frontend:${IMAGE_TAG} \
  ${ACR_REGISTRY}/${ACR_NAMESPACE}/aurora-frontend:${IMAGE_TAG}

docker push ${ACR_REGISTRY}/${ACR_NAMESPACE}/aurora-frontend:${IMAGE_TAG}
```

如需多架构镜像（ACK 节点常见 amd64，本地 Mac 是 arm64），用 buildx：

```bash
docker buildx create --use --name aurora-builder 2>/dev/null || true

docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -f Dockerfile.frontend \
  -t ${ACR_REGISTRY}/${ACR_NAMESPACE}/aurora-frontend:${IMAGE_TAG} \
  --push .
```

> nginx 配置在 `frontend/web/nginx/default.conf`：容器内监听 **8080**（非特权端口，便于 K8s 以非 root 用户跑），已开启 gzip、SPA 路由 fallback（`try_files ... /index.html`），`/healthz` 可用作 ACK 探针。
>
> ACK 部署时 `containerPort` 填 `8080`，`Service.targetPort` 也填 `8080`，对外的 `Service.port` / SLB 监听端口（80、443 等）随意。



