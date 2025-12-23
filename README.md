# DdOaListDownload 项目

## 项目简介

**DdOaListDownload** 是一款企业级钉钉数据集成与管理平台，采用前后端分离架构（Go + Vue 3），辅以 Python 服务支持遗留系统对接。系统专为集团型企业设计，支持多级组织架构管理，实现数据的安全隔离与统一管控。核心功能涵盖钉钉 API 对接、自动化数据同步、细粒度权限控制（RBAC + 字段级）、API 调试与监控等。

## 核心特性

### 🏢 多租户与组织架构
- **集团管控**：支持“集团-分子公司”多级架构，实现数据与配置的逻辑隔离。
- **独立配置**：各分子公司可独立配置钉钉应用凭证与同步策略。

### 🔐 权限与安全
- **RBAC 权限体系**：灵活的角色与菜单权限分配。
- **字段级权限控制**：支持查看、编辑、报表可见性的细粒度控制，包含“特殊权限”覆盖机制。
- **安全性增强**：JWT 认证，动态密钥加载，敏感操作审计。

### 🔌 API 集成与调试
- **API 配置中心**：可视化管理 API 版本、参数与频率限制。
- **真实 API 调试**：内置 HTTP 客户端，支持在线发起真实 API 请求，实时查看响应。
- **免登与鉴权**：自动化 AccessToken 刷新与钉钉免登流程支持。

### 📊 数据管理
- **智能数据字典**：自定义字段属性、必填项、默认值及数据类型。
- **多格式导出**：支持 Excel、CSV、JSON 格式的任务化异步导出。

## 技术栈

- **后端 (Core)**: Go 1.20+ (Gin, GORM, JWT-v5)
- **后端 (Legacy/Aux)**: Python 3.9+ (Flask)
- **前端**: Node.js 18+ (Vue 3, Element Plus, Pinia, Vite)
- **数据存储**: MySQL 8.0+
- **缓存**: Redis 7.0+

## 快速开始

### 1. 环境准备
确保本地已安装 Go, Node.js, MySQL 和 Redis。

### 2. 数据库初始化
本项目采用 **Drop and Recreate** 策略进行初始化，运行以下命令将重置数据库并填充测试数据：
```bash
cd backend/go
go run main.go -init
```

### 3. 启动服务
- **后端**: `cd backend/go && go run main.go` (运行在 :8080)
- **前端**: `cd frontend && npm run dev` (运行在 :3002)

详细部署指南请参考 [安装与部署指南](docs/setup_guide.md)。

## 文档导航

- 📘 [开发环境搭建](docs/setup_guide.md)
- 🏗️ [数据库设计](docs/database_design.md)
- 🔌 [API 接口设计](docs/api_design.md)
- 🚀 [部署运维指南](docs/deployment_guide.md)
- 📝 [变更日志](CHANGELOG.md)

## 许可证

MIT License

---
**作者**: cjx  
**邮箱**: xx4125517@126.com  
**最后更新**: 2025-12-23 15:50:00
