#!/usr/bin/env bash
set -euo pipefail

# 本地构建两个前端工程（web / admin），产物供各自 Dockerfile 直接拷贝。
# 原因：部署服务器 CPU/内存较弱，在 Docker 内跑 pnpm build 会崩，
# 因此改为本机构建好，镜像只做 nginx 静态托管打包。
#
# 用法：
#   bash scripts/build-frontend.sh          # 构建 web + admin
#   bash scripts/build-frontend.sh web      # 只构建 web
#   bash scripts/build-frontend.sh admin    # 只构建 admin

cd "$(dirname "$0")/.."
ROOT="$(pwd)"

if ! command -v corepack >/dev/null 2>&1; then
  echo "错误：未找到 corepack（随 Node.js 提供），请先安装 Node.js 18+。" >&2
  exit 1
fi

build_one() {
  local name="$1" dir="$2" pnpm_ver="$3" out="$4"
  echo
  echo ">> [${name}] 构建中 (pnpm@${pnpm_ver}) ..."
  cd "${ROOT}/${dir}"
  corepack prepare "pnpm@${pnpm_ver}" --activate
  corepack pnpm install --frozen-lockfile
  corepack pnpm build
  if [[ ! -d "${out}" ]]; then
    echo "错误：[${name}] 构建未生成产物目录 ${dir}/${out}" >&2
    exit 1
  fi
  echo ">> [${name}] 完成：${dir}/${out} ($(du -sh "${out}" | awk '{print $1}'))"
}

target="${1:-all}"

case "${target}" in
  web)   build_one web   frontend/web   8.15.9  dist  ;;
  admin) build_one admin frontend/admin 10.28.0 build ;;
  all)
    build_one web   frontend/web   8.15.9  dist
    build_one admin frontend/admin 10.28.0 build
    ;;
  *)
    echo "未知参数：${target}（可选 web | admin | all）" >&2
    exit 1
    ;;
esac

echo
echo ">> 全部完成。下一步：docker compose build web admin  (或 docker compose build)"
