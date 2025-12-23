# DdOaListDownload 部署与运维指南

作者: cjx
邮箱: xx4125517@126.com
时间: 2025-12-23 11:20:00

## 1. 部署架构
系统通常采用经典的三层架构部署：
- **前端 (Web)**: Vue 3 静态编译，Nginx 托管并在反向代理层分发 API 流量。
- **后端 (API)**: Go 二进制服务运行，与遗留或特定的 Python 服务协同。
- **持久层**: MySQL 处理持久化，Redis 处理会话与缓存。

## 2. 生产环境准备

### 2.1 依赖安装
- **MySQL**: 确保创建库并导入 `schema.sql`（推荐 MySQL 8.0+，使用 utf8mb4 字符集）。
- **Redis**: 启动 6379 端口服务。
- **Golang**: 1.20+ 环境。

### 2.2 环境变量配置
在后端运行环境配置以下关键变量：
```bash
# JWT 加密密钥（必填）
export JWT_SECRET=您的安全字符串

# 数据库 DSN
export DB_HOST=localhost
export DB_PORT=3306
# ... 其他 DB 参数 ...
```

## 3. 部署步骤

### 3.1 前端构建
```bash
cd frontend
npm install
npm run build
# 将 dist/ 文件夹下的内容拷贝至 Nginx 静态文件根目录
```

### 3.2 后端 Go 服务部署
```bash
cd backend/go
go mod tidy
go build -o dd-oa-server main.go
# 建议使用 systemd 或 pm2 守护进程运行
# pm2 start ./dd-oa-server --name ddoa-go
```

### 3.3 后端 Python 服务部署 (遗留系统支持)
```bash
cd backend/python
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
# 生产环境建议使用 Gunicorn
# gunicorn -w 4 -b 127.0.0.1:8081 app:app
```

## 4. Nginx 生产环境配置
```nginx
server {
    listen 80;
    server_name your-oa-domain.com;

    # 前端静态资源
    location / {
        root /var/www/dd-oa/frontend;
        try_files $uri $uri/ /index.html;
        index index.html;
    }

    # Go 后端 API
    location /api/v1/ {
        proxy_pass http://127.0.0.1:8080/api/v1/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # Python 辅助 API
    location /api/python/ {
        proxy_pass http://127.0.0.1:8081/;
        proxy_set_header Host $host;
    }
}
```

## 5. 生产环境 Checklist
- [ ] **密钥安全**: 确保 `JWT_SECRET` 不是默认值。
- [ ] **数据库备份**: 配置定时 `mysqldump` 任务。
- [ ] **日志滚动**: 生产环境应开启日志文件滚动，防止磁盘空间耗尽。
- [ ] **防火墙**: 仅开放 80/443 端口，关闭 3306/6379 外部访问。
- [ ] **HTTPS**: 强烈建议配置 SSL 证书（Let's Encrypt 等）。

## 6. 监控与排错
- **日志查看**: 通过 `journalctl -u ddoa-go` 或查看 log 目录文件。
- **健康检查**: 监控系统内存、磁盘空闲率及 Redis 连接情况。
- **Panic 处理**: 系统已内置 Recovery，异常会记录在 error 日志中。
