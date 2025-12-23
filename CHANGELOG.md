# DdOaListDownload 变更日志 (CHANGELOG)

作者: cjx
邮箱: xx4125517@126.com
时间: 2025-12-23 11:15:00

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
