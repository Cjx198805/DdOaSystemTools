# DdOaListDownload 数据库设计文档

作者: cjx
邮箱: xx4125517@126.com
时间: 2025-12-23 11:00:00

## 1. 设计概述

本系统采用 MySQL 作为主关系型数据库，Redis 作为缓存和会话存储。设计核心围绕“集团-公司”多租户隔离架构及“角色-权限-字段”的三维权限体系展开。

## 2. 核心实体模型

### 2.1 集团与公司结构 (Group & Company)
- **Company 表**: 存储组织信息。
  - `id`: 主键
  - `parent_id`: 支持多级组织树。
  - `name`, `code`: 公司名称与唯一标识。

### 2.2 身份与权限 (IAM)
- **User 表**: 用户核心信息。
- **Role 表**: 角色定义。
- **UserRole 表**: 用户-角色多对多映射。
- **Menu 表**: 资源/功能权限菜单树。
- **FieldPermission 表**: 细粒度字段权限控制。
  - `special_edit`: 是否拥有特殊编辑权（优先级高于字典配置）。

### 2.3 业务配置 (Business Config)
- **APIConfig 表**: 钉钉接口对接规格。
- **DataDictionary 表**: 业务字段元数据定义（可编辑性、必填性等）。
- **APITestCase & APITestHistory**: 接口测试资产。

### 2.4 下载引擎 (Download Engine)
- **DownloadTask 表**: 异步任务状态追踪。
- **DownloadResult 表**: 任务结果存储/索引。

## 3. 关键特性实现方案

### 3.1 数据隔离
所有业务相关表（如 `User`, `APIConfig`, `DownloadTask`）均包含 `company_id` 字段。后端通过中间件或 Service 层提取 Context 中的 `CompanyID` 自动注入查询条件，确保不同分子公司间的数据物理/逻辑隔离。

### 3.2 字段权限优先级联动逻辑
系统在判定一个字段是否允许用户编辑时，遵循以下伪代码逻辑：
```sql
Result = (EXISTS FieldPermission WHERE userID=? AND field=? AND SpecialEdit=1) 
         OR 
         (NOT EXISTS DataDictionary WHERE field=? AND Editable=0)
```
即：**用户特权 > 数据字典限制 > 默认放行**。

## 4. 索引优化策略
- 为所有 `company_id` 添加普通索引。
- 为 `username`, `role_code`, `api_code` 添加唯一索引。
- 为 `created_at` 增加索引以支持高效的任务历史查询。
