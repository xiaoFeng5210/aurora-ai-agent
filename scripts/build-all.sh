#!/usr/bin/env bash
set -euo pipefail

# 本机一次性构建后端 + 两个前端，全部重活在本地完成，
# 之后服务器只需 docker compose build 做轻量镜像打包。

cd "$(dirname "$0")"

bash build-backend.sh
bash build-frontend.sh all

echo
echo ">> 后端 + 前端 全部构建完成。"
echo ">> 服务器侧执行：docker compose build && docker compose up -d"
