# 部署指南

## 概述

Password Timer Manager 支持以下部署方式：
- Docker Compose（推荐）
- Kubernetes
- 手动源码部署

---

## Docker Compose 部署

### 环境要求

- Docker 20.10+
- Docker Compose 2.0+

### 快速部署

```bash
# 1. 创建部署目录
sudo mkdir -p /opt/ptm
cd /opt/ptm

# 2. 下载 docker-compose.yml 和环境配置
sudo wget https://your-domain.com/docker/docker-compose.yml
sudo wget https://your-domain.com/docker/.env.example -O .env

# 3. 编辑环境配置
sudo nano .env

# 4. 启动服务
sudo docker-compose up -d

# 5. 检查服务状态
sudo docker-compose ps
```

### .env 配置示例

```env
# 数据库配置
DB_USER=ptm_user
DB_PASSWORD=your_secure_password
DB_NAME=ptm_db

# Redis 配置
REDIS_PASSWORD=your_redis_password

# JWT 配置
JWT_SECRET=your-32-character-secret-key-here

# API 配置
API_PORT=8080
API_BASE_URL=http://your-domain.com

# 前端配置
FRONTEND_PORT=80

# 加密配置
ENCRYPTION_AES_KEY_LENGTH=256
ENCRYPTION_RSA_KEY_LENGTH=2048
ENCRYPTION_PASSWORD_MIN_LENGTH=8
ENCRYPTION_CONFIG_ENCRYPT=true
ENCRYPTION_ENABLE_SIGNATURE=true

# SSL 配置（可选）
SSL_ENABLED=true
SSL_CERT_PATH=/ssl/cert.pem
SSL_KEY_PATH=/ssl/key.pem
```

### 暴露服务到外网

生产环境建议使用 Nginx 反向代理并配置 SSL：

```nginx
# /etc/nginx/conf.d/ptm.conf

upstream api {
    server 127.0.0.1:8080;
}

upstream frontend {
    server 127.0.0.1:80;
}

server {
    listen 80;
    server_name your-domain.com;
    
    # 重定向到 HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;
    
    ssl_certificate /etc/ssl/certs/your-cert.pem;
    ssl_certificate_key /etc/ssl/private/your-key.pem;
    
    # 前端
    location / {
        proxy_pass http://frontend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
    
    # API
    location /api/ {
        proxy_pass http://api;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_read_timeout 300;
        proxy_connect_timeout 300;
    }
}
```

---

## Kubernetes 部署

### 前置条件

- Kubernetes 1.20+
- Helm 3.0+
- PV (PersistentVolume) 用于存储

### 使用 Helm 部署

```bash
# 添加 Helm 仓库
helm repo add ptm https://charts.ptm.com
helm repo update

# 创建 namespace
kubectl create namespace ptm

# 安装
helm install ptm ptm/ptm \
  --namespace ptm \
  --set api.image.tag=v1.0.0 \
  --set api.service.type=LoadBalancer \
  --set database.storage.size=10Gi \
  --set redis.storage.size=1Gi
```

### values.yaml 配置

```yaml
# values.yaml
api:
  image:
    repository: ptm/api
    tag: v1.0.0
  replicaCount: 2
  service:
    type: ClusterIP
    port: 8080
  env:
    JWT_SECRET: "your-secret-key"
    ENCRYPTION_AES_KEY_LENGTH: 256
    ENCRYPTION_RSA_KEY_LENGTH: 2048

frontend:
  image:
    repository: ptm/frontend
    tag: v1.0.0
  replicaCount: 2
  service:
    type: LoadBalancer
    port: 80

database:
  enabled: true
  image: postgres:15-alpine
  storage:
    size: 10Gi
    storageClass: standard
  env:
    POSTGRES_USER: ptm_user
    POSTGRES_PASSWORD: your-password
    POSTGRES_DB: ptm_db

redis:
  enabled: true
  image: redis:7-alpine
  storage:
    size: 1Gi
    storageClass: standard
  env:
    REDIS_PASSWORD: your-redis-password
```

---

## Linux 手动部署

### 环境要求

- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Nginx

### 安装依赖

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y curl wget git nginx postgresql-15 redis-server

# CentOS/RHEL
sudo yum install -y curl wget git nginx postgresql-server redis
```

### 初始化数据库

```bash
# 启动 PostgreSQL
sudo systemctl start postgresql
sudo systemctl enable postgresql

# 创建数据库和用户
sudo -u postgres psql << EOF
CREATE USER ptm_user WITH PASSWORD 'your_secure_password';
CREATE DATABASE ptm_db OWNER ptm_user;
EOF
```

### 初始化 Redis

```bash
sudo systemctl start redis
sudo systemctl enable redis
```

### 部署后端

```bash
# 下载源码
cd /opt
sudo git clone https://your-repo/password-timer-manager.git ptm
cd ptm/backend

# 编译
sudo go build -o api ./cmd/server

# 创建运行用户
sudo useradd -r -s /sbin/nologin ptm
sudo chown -R ptm:ptm /opt/ptm

# 创建 systemd 服务
sudo cat > /etc/systemd/system/ptm-api.service << EOF
[Unit]
Description=PTM API Server
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=ptm
WorkingDirectory=/opt/ptm/backend
ExecStart=/opt/ptm/backend/api
Restart=always

[Install]
WantedBy=multi-user.target
EOF

# 启动服务
sudo systemctl daemon-reload
sudo systemctl start ptm-api
sudo systemctl enable ptm-api
```

### 部署前端

```bash
cd /opt/ptm/frontend

# 安装依赖并构建
npm install
npm run build

# 配置 Nginx
sudo cat > /etc/nginx/sites-available/ptm << EOF
server {
    listen 80;
    server_name _;
    
    root /opt/ptm/frontend/dist;
    index index.html;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    location /api/ {
        proxy_pass http://127.0.0.1:8080/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
EOF

sudo ln -s /etc/nginx/sites-available/ptm /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

---

## Windows 手动部署

### 环境要求

- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- IIS 或 Nginx

### 使用 PowerShell 脚本部署

```powershell
# 以管理员身份运行 PowerShell

# 1. 下载源码
cd C:\Program Files
git clone https://your-repo/password-timer-manager.git PTM

# 2. 安装 PostgreSQL
# 下载并安装: https://www.postgresql.org/download/windows/

# 3. 创建数据库
psql -U postgres -c "CREATE USER ptm_user WITH PASSWORD 'your_password';"
psql -U postgres -c "CREATE DATABASE ptm_db OWNER ptm_user;"

# 4. 安装 Redis（使用 Chocolatey）
choco install redis-64 -y

# 5. 编译后端
cd C:\Program Files\PTM\backend
go build -o api.exe .\cmd\server

# 6. 创建 Windows 服务
# 使用 NSSM: choco install nssm
nssm install PTM-API "C:\Program Files\PTM\backend\api.exe"
nssm set PTM-API AppDirectory "C:\Program Files\PTM\backend"
nssm set PTM-API Start SERVICE_AUTO_START
nssm start PTM-API

# 7. 构建前端
cd C:\Program Files\PTM\frontend
npm install
npm run build

# 8. 使用 IIS 部署前端
# - 安装 IIS
# - 将前端 dist 目录发布为网站
# - 配置反向代理到 localhost:8080
```

---

## 验证部署

### 检查服务状态

```bash
# API 健康检查
curl http://localhost:8080/api/v1/health

# 预期响应
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "healthy",
    "version": "1.0.0"
  }
}
```

### 登录测试

```bash
# 登录获取 Token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

---

## 数据备份

### 数据库备份

```bash
# PostgreSQL
pg_dump -U ptm_user -d ptm_db > backup_$(date +%Y%m%d).sql

# 恢复
psql -U ptm_user -d ptm_db < backup_20240101.sql
```

### 文件存储备份

```bash
# 备份上传的文件
tar -czf files_backup_$(date +%Y%m%d).tar.gz /opt/ptm/storage
```

---

## 升级

### Docker Compose 升级

```bash
cd /opt/ptm

# 拉取新版本
git pull

# 重新构建并启动
docker-compose down
docker-compose build
docker-compose up -d
```

### 数据库迁移

```bash
# 运行迁移脚本
docker-compose exec api ./api migrate
```

---

## 卸载

### Docker Compose 卸载

```bash
cd /opt/ptm

# 停止并删除容器
docker-compose down

# 删除数据卷（谨慎操作）
docker-compose down -v

# 删除源码目录
sudo rm -rf /opt/ptm
```

### 系统服务卸载

```bash
# 停止服务
sudo systemctl stop ptm-api
sudo systemctl disable ptm-api
sudo rm /etc/systemd/system/ptm-api.service

# 删除用户
sudo userdel ptm
```
