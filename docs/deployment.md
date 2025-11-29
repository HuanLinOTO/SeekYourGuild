# SeekYourGuild 部署指南

本文档详细介绍如何在生产环境中部署 SeekYourGuild。

## 目录

- [系统要求](#系统要求)
- [Docker 部署](#docker-部署)
- [手动部署](#手动部署)
- [反向代理配置](#反向代理配置)
- [数据库配置](#数据库配置)
- [机器人接入](#机器人接入)
- [常见问题](#常见问题)

---

## 系统要求

### 最低配置

- CPU: 1 核
- 内存: 512MB
- 存储: 1GB
- 系统: Linux / Windows / macOS

### 推荐配置

- CPU: 2 核
- 内存: 1GB
- 存储: 10GB (取决于数据量)
- 系统: Linux (Ubuntu 22.04+)

---

## Docker 部署

### 方式一：使用 docker-compose（推荐）

1. **下载配置文件**

```bash
# 创建目录
mkdir -p /opt/seekyourguild && cd /opt/seekyourguild

# 下载 docker-compose.yml
wget https://raw.githubusercontent.com/yourname/seekyourguild/main/docker-compose.yml
```

2. **创建环境变量文件**

```bash
cat > .env << EOF
# 机器人上报密钥（请修改为复杂的随机字符串）
API_KEY=sk-your-very-secret-key-here

# 允许的群号列表（逗号分隔，为空则允许所有）
ALLOWED_GROUPS=123456789,987654321

# 可选：使用 PostgreSQL
# DB_DRIVER=postgres
# DB_DSN=postgres://user:pass@localhost:5432/seekyourguild?sslmode=disable
EOF
```

3. **启动服务**

```bash
docker-compose up -d
```

4. **查看日志**

```bash
docker-compose logs -f
```

### 方式二：使用预构建镜像

```bash
docker run -d \
  --name seekyourguild \
  -p 8080:8080 \
  -v /opt/seekyourguild/data:/data \
  -e API_KEY=sk-your-secret-key \
  -e ALLOWED_GROUPS=123456789,987654321 \
  ghcr.io/yourname/seekyourguild:latest
```

### 更新

```bash
# 拉取最新镜像
docker-compose pull

# 重启服务
docker-compose up -d
```

---

## 手动部署

### 1. 后端部署

**下载预编译二进制**

从 [Releases](https://github.com/yourname/seekyourguild/releases) 下载对应平台的二进制文件：

- `seekyourguild-linux-amd64` - Linux x64
- `seekyourguild-linux-arm64` - Linux ARM64
- `seekyourguild-windows-amd64.exe` - Windows
- `seekyourguild-darwin-amd64` - macOS Intel
- `seekyourguild-darwin-arm64` - macOS Apple Silicon

**Linux 部署示例**

```bash
# 下载
wget https://github.com/yourname/seekyourguild/releases/latest/download/seekyourguild-linux-amd64
chmod +x seekyourguild-linux-amd64

# 创建配置
export API_KEY="your-secret-key"
export ALLOWED_GROUPS="123456789,987654321"

# 运行
./seekyourguild-linux-amd64
```

**使用 Systemd 管理（推荐）**

```bash
# 创建服务文件
sudo cat > /etc/systemd/system/seekyourguild.service << EOF
[Unit]
Description=SeekYourGuild Server
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/seekyourguild
ExecStart=/opt/seekyourguild/seekyourguild-linux-amd64
Restart=always
RestartSec=5

Environment=API_KEY=your-secret-key
Environment=ALLOWED_GROUPS=123456789,987654321
Environment=DB_DSN=file:/opt/seekyourguild/data/seekyourguild.db?cache=shared&mode=rwc

[Install]
WantedBy=multi-user.target
EOF

# 启动服务
sudo systemctl daemon-reload
sudo systemctl enable seekyourguild
sudo systemctl start seekyourguild

# 查看状态
sudo systemctl status seekyourguild
```

### 2. 前端部署

前端已经内嵌在后端二进制中，无需单独部署。

如需自定义前端：

```bash
# 下载前端构建产物
wget https://github.com/yourname/seekyourguild/releases/latest/download/frontend-dist.zip
unzip frontend-dist.zip -d /var/www/seekyourguild

# 使用 Nginx 托管
```

---

## 反向代理配置

### Nginx

```nginx
server {
    listen 80;
    server_name seekyourguild.example.com;

    # 重定向到 HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name seekyourguild.example.com;

    # SSL 证书
    ssl_certificate /etc/letsencrypt/live/seekyourguild.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/seekyourguild.example.com/privkey.pem;

    # 安全头
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket 支持（未来版本）
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

### Caddy

```caddyfile
seekyourguild.example.com {
    reverse_proxy localhost:8080
}
```

---

## 数据库配置

### SQLite（默认）

适合小规模部署，无需额外配置。

```bash
export DB_DRIVER=sqlite
export DB_DSN="file:/data/seekyourguild.db?cache=shared&mode=rwc"
```

### PostgreSQL（推荐生产环境）

```bash
# 创建数据库
createdb seekyourguild

# 配置连接
export DB_DRIVER=postgres
export DB_DSN="postgres://user:password@localhost:5432/seekyourguild?sslmode=disable"
```

使用 Docker Compose 部署 PostgreSQL：

```yaml
services:
  seekyourguild:
    # ...
    environment:
      - DB_DRIVER=postgres
      - DB_DSN=postgres://seekyourguild:password@postgres:5432/seekyourguild?sslmode=disable
    depends_on:
      - postgres

  postgres:
    image: postgres:16-alpine
    environment:
      - POSTGRES_USER=seekyourguild
      - POSTGRES_PASSWORD=password
      - POSTGRES_DB=seekyourguild
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

---

## 机器人接入

### 1. 获取 API Key

在部署时设置的 `API_KEY` 环境变量即为机器人上报使用的密钥。

### 2. 配置机器人

在机器人端配置以下信息：

```
API 地址: https://your-domain.com/api/report
API Key: your-api-key
```

### 3. 上报示例

参考 [机器人开发指南](../机器人端/开发指南.md)

---

## 常见问题

### Q: 端口被占用怎么办？

修改 `SERVER_PORT` 环境变量：

```bash
export SERVER_PORT=3000
```

### Q: 如何备份数据？

**SQLite:**

```bash
cp /data/seekyourguild.db /backup/seekyourguild-$(date +%Y%m%d).db
```

**PostgreSQL:**

```bash
pg_dump seekyourguild > /backup/seekyourguild-$(date +%Y%m%d).sql
```

### Q: 如何查看日志？

**Docker:**

```bash
docker-compose logs -f seekyourguild
```

**Systemd:**

```bash
journalctl -u seekyourguild -f
```

### Q: 如何更改允许的群？

修改 `ALLOWED_GROUPS` 环境变量后重启服务：

```bash
export ALLOWED_GROUPS="111,222,333"
docker-compose restart
```

### Q: 内存占用过高怎么办？

SQLite 默认会缓存数据，可以限制缓存大小：

```bash
export DB_DSN="file:/data/seekyourguild.db?cache=shared&mode=rwc&_cache_size=-20000"
```

---

## 监控与告警

### 健康检查

```bash
curl http://localhost:8080/api/groups
```

### Prometheus 指标（计划中）

未来版本将支持 `/metrics` 端点。

---

如有问题，请提交 [Issue](https://github.com/yourname/seekyourguild/issues)。
