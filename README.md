# aurora-ai-agent


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
