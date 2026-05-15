#!/usr/bin/env bash
set -euo pipefail

# 本地交叉编译 aurora-agent 后端二进制，供 Dockerfile.backend 直接拷贝使用。
# 原因：部署服务器 CPU 较弱，在 Docker 内编译 Go 非常吃力，
# 因此改为在本机编译好，镜像只做轻量打包。

cd "$(dirname "$0")/.."

OUT_DIR="bin"
OUT_BIN="${OUT_DIR}/aurora-agent"

# 目标为部署服务器架构（linux/amd64），可用环境变量覆盖 GOARCH。
TARGET_GOOS="linux"
TARGET_GOARCH="${GOARCH:-amd64}"

if ! command -v go >/dev/null 2>&1; then
  echo "错误：未找到 go 命令，请先安装 Go。" >&2
  exit 1
fi

mkdir -p "${OUT_DIR}"

echo ">> 编译 ${OUT_BIN} (${TARGET_GOOS}/${TARGET_GOARCH}) ..."
CGO_ENABLED=0 GOOS="${TARGET_GOOS}" GOARCH="${TARGET_GOARCH}" \
  go build -trimpath -ldflags="-s -w" -o "${OUT_BIN}" .

echo ">> 完成：$(ls -lh "${OUT_BIN}" | awk '{print $5, $9}')"
echo ">> 下一步：docker compose build backend  (或 docker build -f Dockerfile.backend .)"
