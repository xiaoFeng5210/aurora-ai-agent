# 镜像清单与上传

这份目录只回答两件事：Ariadne 有哪些镜像，以及怎么把它们推到阿里云 ACR，供 ACK 拉取。

不写 Deployment / Service / Ingress。单机运行仍然用仓库根目录的 `docker-compose.yml`。

镜像地址由同目录的 `docker-compose.yml` 决定。改仓库名、版本号，只改 `.env` 和这份 compose，不要在文档里手写第二套地址。

## 1. 镜像清单

和根目录 compose 对齐，一共 6 个。前 3 个自己构建，后 3 个从官方镜像同步到 ACR。

| Compose 服务 | 类型 | 构建上下文 | Dockerfile | 推到 ACR 后的仓库名 |
| --- | --- | --- | --- | --- |
| `backend` | 自建 | 仓库根目录 | `Dockerfile.backend` | `ariadne-backend` |
| `frontend` | 自建 | `frontend/web` | `frontend/web/Dockerfile` | `ariadne-frontend` |
| `admin` | 自建 | `frontend/admin` | `frontend/admin/Dockerfile` | `ariadne-admin` |
| `postgres` | 同步 | — | 上游 `postgres:16-alpine` | `postgres` |
| `redis` | 同步 | — | 上游 `redis:7-alpine` | `redis` |
| `rabbitmq` | 同步 | — | 上游 `rabbitmq:3.13-management-alpine` | `rabbitmq` |

Qdrant 继续用外部服务，这里不打包。

完整地址形态：

```text
<ACR_REGISTRY>/<ACR_NAMESPACE>/<仓库名>:<IMAGE_TAG>
```

个人版示例：

```text
registry.cn-hangzhou.aliyuncs.com/ariadne/ariadne-backend:v0.0.1
```

企业版示例：

```text
<实例名>-registry.cn-hangzhou.cr.aliyuncs.com/ariadne/ariadne-backend:v0.0.1
```

本地单机 compose 仍叫 `aurora-*:compose`。推给 ACK 的名字改成产品名 `ariadne-*`，避免和旧的 compose tag 混在一起。

## 2. 先在 ACR 准备什么

1. 开通 [容器镜像服务 ACR](https://cr.console.aliyun.com/)，个人版或企业版都可以。
2. 地域尽量和 ACK 集群相同。跨地域拉镜像更慢，还可能产生流量费。
3. 建一个命名空间，例如 `ariadne`。
4. 在「访问凭证」里设置 Registry 登录密码。用户名以控制台显示的为准，**不要用纯数字的主账号 ID**。
5. 仓库可以先建好 `ariadne-backend` / `ariadne-frontend` / `ariadne-admin` / `postgres` / `redis` / `rabbitmq`，也可以打开命名空间的自动创建仓库。

官方说明：

- [推送和拉取镜像（企业版）](https://help.aliyun.com/zh/acr/getting-started/use-a-container-registry-enterprise-edition-instance-to-push-and-pull-images)
- [配置访问凭证](https://help.aliyun.com/zh/acr/user-guide/configure-access-credentials)
- [ACK 使用私有镜像](https://help.aliyun.com/zh/ack/create-an-application-by-using-a-private-image-repository)

## 3. 本机配置

在仓库根目录：

```bash
cp .k8s/images/.env.example .k8s/images/.env
```

至少改这三项：

```bash
ACR_REGISTRY=registry.cn-hangzhou.aliyuncs.com
ACR_NAMESPACE=ariadne
IMAGE_TAG=v0.0.1
```

企业版把 `ACR_REGISTRY` 换成控制台给的实例域名。`.env` 已被 git 忽略。

先看 compose 解析出来的最终镜像名，确认没写错：

```bash
docker compose -f .k8s/images/docker-compose.yml --env-file .k8s/images/.env config
```

## 4. 登录 ACR

个人版：

```bash
docker login --username=<控制台显示的用户名> registry.cn-<地域>.aliyuncs.com
```

企业版：

```bash
docker login --username=<控制台显示的用户名> <实例名>-registry.cn-<地域>.cr.aliyuncs.com
```

密码是 ACR「访问凭证」里设置的 Registry 密码，不是阿里云控制台登录密码。看到 `Login Succeeded` 再继续。

密码里如果有 `@` 这类字符，终端可能被 Shell 吃掉，控制台会报 `unauthorized`。这时去访问凭证页重设一个不含特殊字符的密码。

## 5. 推自建镜像

镜像只打包已经编好的产物，不会在 Docker 里跑 `go build` 或 `pnpm build`。

```bash
bash scripts/build-all.sh
```

只更新其中一个时：

```bash
bash scripts/build-backend.sh
bash scripts/build-frontend.sh web
bash scripts/build-frontend.sh admin
```

后端脚本默认交叉编译 `linux/amd64`，和 ACK 常见节点一致。三个业务镜像也都写了 `platform: linux/amd64`。本机如果是 Apple Silicon，需要 Docker Desktop 打开容器的虚拟化 / qemu，才能打出 amd64 镜像。

然后构建并推送：

```bash
docker compose -f .k8s/images/docker-compose.yml --env-file .k8s/images/.env build
docker compose -f .k8s/images/docker-compose.yml --env-file .k8s/images/.env push
```

只推其中一个：

```bash
docker compose -f .k8s/images/docker-compose.yml --env-file .k8s/images/.env build backend
docker compose -f .k8s/images/docker-compose.yml --env-file .k8s/images/.env push backend
```

`build` 只会打有 `build:` 的 3 个业务镜像。`postgres` / `redis` / `rabbitmq` 挂了 `profiles: [base]`，不会被这次 build/push 带上。

推成功后，到 ACR 控制台对应仓库的「镜像版本」里应能看到 `IMAGE_TAG`。

## 6. 把基础镜像同步到 ACR

ACK 节点直拉 Docker Hub 经常慢或失败，所以官方基础镜像也推进自己的 ACR。

先按 amd64 拉上游，再打上 compose 里的 ACR 名字：

```bash
set -a
source .k8s/images/.env
set +a

docker pull --platform linux/amd64 postgres:${POSTGRES_TAG}
docker pull --platform linux/amd64 redis:${REDIS_TAG}
docker pull --platform linux/amd64 rabbitmq:${RABBITMQ_TAG}

docker tag postgres:${POSTGRES_TAG} \
  ${ACR_REGISTRY}/${ACR_NAMESPACE}/postgres:${POSTGRES_TAG}
docker tag redis:${REDIS_TAG} \
  ${ACR_REGISTRY}/${ACR_NAMESPACE}/redis:${REDIS_TAG}
docker tag rabbitmq:${RABBITMQ_TAG} \
  ${ACR_REGISTRY}/${ACR_NAMESPACE}/rabbitmq:${RABBITMQ_TAG}

docker compose -f .k8s/images/docker-compose.yml --env-file .k8s/images/.env --profile base push
```

基础镜像的 tag 跟上游走（`16-alpine` 这种），不要和业务的 `v0.0.1` 混用。

## 7. ACK 之后怎么用这些镜像

编排 YAML 还没写。到时候 Pod 里的 `image` 直接填 compose 解析出来的完整地址。

现在只需要知道这三件事：

1. **本机推送用公网域名**，也就是 `.env` 里的 `ACR_REGISTRY`。
2. **集群内拉取优先用同地域 VPC 域名**，少走公网。个人版一般是 `registry-vpc.cn-<地域>.aliyuncs.com`，企业版是 `<实例名>-registry-vpc.cn-<地域>.cr.aliyuncs.com`。仓库路径和 tag 不变。
3. **私有仓库 ACK 拉不下来时**，给工作负载配 `imagePullSecrets`，或对企业版安装免密组件。2024-09-09 之后新建的 ACR 个人版不支持免密组件，用 Secret。Secret 必须和工作负载在同一个命名空间。

## 8. 常见问题

**`COPY bin/aurora-agent` 失败**  
先跑 `bash scripts/build-backend.sh`。这个文件被 `.gitignore` 忽略，仓库里不会有。

**`COPY dist` / `COPY build` 失败**  
先跑 `bash scripts/build-frontend.sh`。web 产物在 `frontend/web/dist`，admin 产物在 `frontend/admin/build`。

**ACK 里 `exec format error`**  
镜像或后端二进制是 arm64，节点是 amd64。确认 `scripts/build-backend.sh` 没有把 `GOARCH` 改成 `arm64`，并且 compose 构建时带了 `platform: linux/amd64`。

**`denied: requested access to the resource is denied`**  
没登录、登错实例域名、命名空间不对，或仓库不存在且未打开自动创建。

**`unauthorized: authentication required`**  
用户名用了主账号 ID，或密码被 Shell 转义。回 ACR 访问凭证页核对 `docker login` 整行命令。
