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

- **API 配置**：
  - 钉钉应用配置（AppKey、AppSecret）
  - API 版本管理
  - 接口参数配置
  - 调用频率限制

- **API 测试**：
  - 在线 API 调试
  - 测试用例管理
  - 响应结果查看
  - 调用历史记录

- **系统管理**：
  - 系统参数配置
  - 数据库管理
  - 日志管理
  - 备份恢复

## 快速开始

### 环境要求

- **后端**：
  - Go 1.20+
  - Python 3.9+
- **前端**：
  - Node.js 18+
- **数据库**：
  - MySQL 8.0+
  - Redis 7.0+

### 开发环境搭建

请参考 [开发环境搭建指南](docs/setup_guide.md)。

## 项目结构

```
DdOaListDownload/
├── README.md              # 项目说明文档
├── backend/               # 后端目录
│   ├── go/                # Go 服务（处理新版 API）
│   └── python/             # Python 服务（处理旧版 API）
├── frontend/              # 前端目录
└── docs/                 # 文档目录
    ├── setup_guide.md      # 开发环境搭建指南
    ├── development_guide.md  # 开发规范指南
    └── implementation_plan.md  # 实施方案
```

## 开发指南

请参考 [开发规范](docs/development_guide.md)。

## 许可证

MIT License

作者: cjx
邮箱: xx4125517@126.com
时间: 2025-12-20 14:45:00