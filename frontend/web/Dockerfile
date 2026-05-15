# web 静态产物由本机脚本预先构建（部署服务器跑 pnpm build 会崩）。
# 构建镜像前请先执行：bash scripts/build-frontend.sh web
FROM nginx:1.27-alpine

COPY nginx/default.conf /etc/nginx/conf.d/default.conf
COPY dist /usr/share/nginx/html

EXPOSE 8080
