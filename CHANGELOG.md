# DdOaListDownload 变更日志 (CHANGELOG)

作者: cjx
邮箱: xx4125517@126.com
时间: 2025-12-23 17:25:00

## [Unreleased]

## [1.2.0] - 2025-12-23
### 增加
- **前端**: 全新重构，由传统 Vue 3 迁移至 **Vue Vben Admin 5 (Ant Design Vue 版)**。
- **前端**: 引入 Pnpm Monorepo 架构，提升包管理效率与构建速度。
- **前端**: 实现动态菜单对接（Backend Access Mode），支持从后端加载 RBAC 菜单。
- **前端**: 业务功能完整迁移（用户管理、角色管理、菜单管理、公司管理、数据字典、字段权限、下载任务、API 测试）。
- **前端**: 引入 `IconifyIcon` + Lucide 方案，实现全系统图标标准化。

### 修复
- **前端**: 解决 ESM 模式下由于 Namespace 引用导致的 `SyntaxError` (统一 import type)。
- **前端**: 修复页面刷新 404 及旧路径（如 `/analytics`）自动重定向逻辑。
- **前端**: 解决 `@vben/icons` 导出缺失问题。

### 优化
- **清理**: 移除 `frontend_legacy` 冗余目录，精简工作空间。
- **文档**: 更新 `README.md` 与 `setup_guide.md` 以适配 Vben Admin 架构。

### Added
- **Backend**: Implemented core business logic for Phase 1.
    - `FieldPermission` service with priority logic (SpecialEdit > Dictionary > Default).
    - `APITest` service with real HTTP request capability (URL construction, parameter merging).
    - `DownloadTask` service with permission isolation and robust path handling.
    - Dynamic loading of JWT Secret from environment variables.
- **Backend/CLI**: Added `-init` flag to `main.go` for "Drop and Recreate" database initialization strategy.
- **Frontend**: Added missing API modules (`menu.js`, `download.js`).

### Fixed
- **Backend**: Resolved GORM Error 1091 (duplicate key name) by adopting "Drop and Recreate" strategy.
- **Backend**: Fixed `user_roles` table name mismatch in `FieldPermission` query.
- **Frontend**: Fixed `Layout.vue` empty page issue by replacing `<slot>` with `<router-view>`.
- **Frontend**: Fixed `user.js` store crash due to invalid JSON in localStorage.
- **Frontend**: Fixed router path mismatches in `router/index.js` and `Layout.vue` menu links.
- **Frontend**: Fixed 500 errors in components by creating missing API files and correcting verification imports.

## [1.1.0] - 2025-12-23
### 增加
- **后端**: 实现 `RecoverMiddleware` 全局异常捕获逻辑。
- **后端**: `API测试` 模块由 Mock 升级为真实 HTTP 调用引擎。
- **前端**: 重构登录页面与字段权限管理页面（Element Plus 风格）。
- **前端**: Axios 拦截器增加 401 自动过期处理及 `ElMessage` 全局反馈。

### 优化
- **安全性**: JWT Secret 迁移至环境变量配置，支持默认值回退。
- **权限引擎**: `FieldPermission` 整合 `SpecialEdit` 优先级覆盖逻辑。
- **UI/UX**: 布局组件支持基于路由元数据的动态标题，移除硬编码映射。
- **架构**: 修正大量 API 路径不一致问题（Singular vs Plural）。

### 文档
- 补充 `database_design.md`。
- 完善 `api_design.md`。
- 更新 `README.md` 与 `task.md`。

---

## [1.0.0] - 2025-12-20
### 初始发布
- 基础 RBAC 权限框架。
- 钉钉应用配置与 Token 管理基础功能。
- 基础 Vue 3 + Gin 骨架搭建。
