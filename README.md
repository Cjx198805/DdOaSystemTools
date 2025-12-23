# DdOaListDownload 项目

## 项目简介

DdOaListDownload 是一套完整的钉钉 API 对接系统，采用 Go + Vue 3 + Python 开发，支持集团公司多级管理架构，集团总部可以管理多个不同分子公司的钉钉对接、人员、权限等，分子公司可以管理自己公司的钉钉对接、人员、权限等。系统包含权限管理、API 配置、API 测试、身份验证、accessToken 获取等核心功能，为企业提供完整的钉钉 API 对接解决方案。

## 功能特性

### 基础功能
- 支持从钉钉 OA 系统获取列表数据
- 提供数据下载和导出功能
- 支持多种文件格式导出（Excel、CSV、JSON）

### 核心功能
- **集团公司管理**：
  - 集团总部与分子公司多级架构
  - 分子公司独立管理
  - 数据隔离与权限控制

- **身份验证（免登）**：
  - 钉钉免登配置
  - 免登测试
  - 免登结果记录

- **accessToken 获取**：
  - accessToken 配置
  - accessToken 自动刷新
  - accessToken 测试
  - accessToken 状态监控

- **权限管理**：
  - 用户角色管理
  - 单用户多角色分配
  - 菜单权限控制
  - 字段级权限控制（查看、编辑、报表显示）
  - 特殊权限支持（优先级高于数据字典）
  - API 访问权限
  - 操作日志记录

- **数据字典管理**：
  - 字段属性管理
  - 字段是否必填配置
  - 创建后是否可编辑配置
  - 数据类型管理
  - 默认值设置

  - 默认值设置

- **API 配置**：
  - 钉钉应用配置（AppKey、AppSecret）
  - API 版本管理
  - 接口参数配置
  - 调用频率限制

- **API 测试**：
  - 在线 API 真实调试（支持 HTTP 客户端调用）
  - 测试用例管理与参数合并逻辑
  - 实时响应时间统计与结果持久化
  - 调用历史记录管理

- **系统管理**：
  - 系统参数配置
  - 数据库管理
  - 日志管理
  - 备份恢复

- **安全性与稳定性**：
  - JWT Secret 环境变量化配置
  - 全局 Panic 自动恢复与异常捕获
  - 前端 401 自动重定向与存储防护

## 快速开始

### 环境要求

- **后端**：
  - Go 1.20+ (Gin, GORM, JWT-v5)
  - Python 3.9+ (Flask, SQLAlchemy)
- **前端**：
  - Node.js 18+ (Vue 3, Element Plus, Pinia)
- **基础设施**：
  - MySQL 8.0+
  - Redis 7.0+

### 开发环境搭建

请参考 [开发环境搭建指南](docs/setup_guide.md)。

## 项目结构

```
DdOaListDownload/
├── README.md              # 项目主要说明文档
├── backend/               # 后端服务
│   ├── go/                # 主业务逻辑及新版 API
│   └── python/             # 遗留系统对接与辅助服务
├── frontend/              # 前端 Web 应用
└── docs/                 # 系统文档体系
    ├── setup_guide.md      # 开发环境搭建指南
    ├── database_design.md  # 数据库结构与模型设计
    ├── api_design.md       # API 接口规格说明
    ├── development_guide.md # 核心开发规范
    └── deployment_guide.md  # 生产环境部署方案
```

## 开发指南

请参考 [开发规范](docs/development_guide.md)。

## 变更记录

请参考 [CHANGELOG.md](CHANGELOG.md)。

## 许可证

MIT License

作者: cjx
邮箱: xx4125517@126.com
时间: 2025-12-23 10:55:00