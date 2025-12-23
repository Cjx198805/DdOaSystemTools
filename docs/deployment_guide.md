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

### 3.2 后端部署
```bash
cd backend/go
go mod tidy
go build -o dd-oa-server main.go
./dd-oa-server
```

## 4. Nginx 配置示例
```nginx
server {
    listen 80;
    server_name your-domain.com;

    root /path/to/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api/v1/ {
        proxy_pass http://localhost:8080/api/v1/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 5. 监控与排错
- **日志查看**: 后端标准输出已被重定向或通过 logrus 记录。
- **Panic 处理**: 系统已内置容错，请检查 `/error` 日志中的 StackTrace。
- **健康检查**: 监控 8080 端口存活状态。
