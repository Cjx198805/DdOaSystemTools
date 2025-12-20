# DdOaListDownload 项目实施方案

## 1. 项目概述

DdOaListDownload 是一套完整的钉钉 API 对接系统，采用 Go + Vue 3 开发，支持集团公司多级管理架构，集团总部可以管理多个不同分子公司的钉钉对接、人员、权限等，分子公司可以管理自己公司的钉钉对接、人员、权限等。系统包含权限管理、API 配置、API 测试、身份验证、accessToken 获取等核心功能，为企业提供完整的钉钉 API 对接解决方案。

## 2. 核心功能

### 2.1 基础功能
- 支持从钉钉 OA 系统获取列表数据
- 提供数据下载和导出功能
- 支持多种文件格式导出（Excel、CSV、JSON）

### 2.2 核心功能
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
  - 菜单权限控制
  - API 访问权限
  - 操作日志记录

- **API 配置**：
  - 钉钉应用配置
    - 旧版 API：AppKey、AppSecret
    - 新版 API：根据官方文档配置相应参数
  - API 版本管理（区分旧版和新版 API）
  - 接口参数配置（根据 API 版本动态调整）
  - 调用频率限制（根据官方文档要求）
  - 严格参考钉钉官方开发文档

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

## 3. 技术栈

- **后端语言**：
  - **主要语言**：Go 1.20+（处理新版 API）
  - **支撑语言**：Python 3.9+（处理旧版 API，旧版 SDK 不支持 Go）
- **前端框架**：Vue 3
- **钉钉 SDK 版本**：
  - **新版服务端 SDK**：dingtalk-sdk-go v1.5.0+（官方推荐最新稳定版）
  - **旧版服务端 SDK**：dingtalk-sdk-python v1.5.0+（用于处理旧版 API）
  - **客户端 SDK**：dingtalk-jsapi v3.0.0+（官方推荐最新稳定版）
- **数据库**：
  - **主数据库**：MySQL 8.0+
  - **缓存数据库**：Redis 7.0+
  - **ORM 框架**：
    - Go：gorm.io/gorm v1.25.0+（官方推荐最新稳定版）
    - Python：SQLAlchemy v2.0+（用于 Python 服务）
- **后端依赖**：
  - **Go 依赖**：
    - github.com/gin-gonic/gin：Web 框架
    - github.com/sirupsen/logrus：日志库
    - github.com/go-playground/validator/v10：参数验证
    - github.com/dingtalk/dingtalk-sdk-go：钉钉服务端 SDK（新版 API）
    - gorm.io/gorm：ORM 框架
    - gorm.io/driver/mysql：MySQL 驱动
    - github.com/go-redis/redis/v8：Redis 客户端
  - **Python 依赖**：
    - dingtalk-sdk-python：钉钉服务端 SDK（旧版 API）
    - flask：Web 框架（用于 Python 服务）
    - sqlalchemy：ORM 框架
    - pymysql：MySQL 驱动
    - redis：Redis 客户端
- **前端依赖**：
  - Vue 3
  - Element Plus：UI 组件库
  - Axios：HTTP 客户端
  - Vue Router：路由管理
  - Pinia：状态管理
  - dingtalk-jsapi：钉钉客户端 SDK

### 3.1 官方文档参考

钉钉官方开发文档：[https://open.dingtalk.com/document/development/](https://open.dingtalk.com/document/development/)

所有 API 开发必须严格参照钉钉官方文档进行，确保与官方 API 保持一致，避免使用非官方或已弃用的 API。

## 4. 项目结构

```
DdOaListDownload/
├── README.md              # 项目说明文档
├── backend/               # 后端目录
│   ├── go/                # Go 服务（处理新版 API）
│   │   ├── main.go           # Go 后端入口文件
│   │   ├── go.mod            # Go 模块文件
│   │   ├── go.sum            # Go 依赖校验文件
│   │   ├── config/           # 配置目录
│   │   │   └── config.go     # 配置管理
│   │   ├── controller/       # 控制器目录
│   │   │   ├── downloader.go       # 下载控制器
│   │   │   ├── auth.go             # 认证控制器
│   │   │   ├── company.go           # 集团公司控制器
│   │   │   ├── sso.go               # 身份验证（免登）控制器
│   │   │   ├── access_token.go      # accessToken 控制器
│   │   │   ├── permission.go        # 权限控制器
│   │   │   ├── api_config.go        # API 配置控制器
│   │   │   ├── api_test.go          # API 测试控制器
│   │   │   └── system.go            # 系统管理控制器
│   │   ├── service/          # 服务层目录
│   │   │   ├── downloader.go       # 下载服务
│   │   │   ├── auth.go             # 认证服务
│   │   │   ├── company.go           # 集团公司服务
│   │   │   ├── sso.go               # 身份验证（免登）服务
│   │   │   ├── access_token.go      # accessToken 服务
│   │   │   ├── permission.go        # 权限服务
│   │   │   ├── api_config.go        # API 配置服务
│   │   │   ├── api_test.go          # API 测试服务
│   │   │   └── system.go            # 系统管理服务
│   │   ├── model/            # 数据模型目录
│   │   │   ├── data.go             # 数据模型
│   │   │   ├── company.go          # 集团公司模型
│   │   │   ├── user.go             # 用户模型
│   │   │   ├── role.go             # 角色模型
│   │   │   ├── menu.go             # 菜单模型
│   │   │   ├── sso_config.go       # 身份验证（免登）配置模型
│   │   │   ├── access_token.go     # accessToken 配置模型
│   │   │   ├── api_config.go       # API 配置模型
│   │   │   └── log.go              # 日志模型
│   │   ├── middleware/       # 中间件目录
│   │   │   ├── auth.go       # 认证中间件
│   │   │   ├── cors.go       # CORS 中间件
│   │   │   └── logger.go     # 日志中间件
│   │   ├── utils/            # 工具函数目录
│   │   │   ├── http_client.go # HTTP 客户端
│   │   │   ├── jwt.go        # JWT 工具
│   │   │   └── encrypt.go    # 加密工具
│   │   └── database/         # 数据库目录
│   │       ├── db.go         # 数据库连接
│   │       └── redis.go      # Redis 连接
│   └── python/             # Python 服务（处理旧版 API）
│       ├── app.py            # Python 后端入口文件
│       ├── requirements.txt  # Python 依赖
│       ├── config/           # 配置目录
│       │   └── config.py     # 配置管理
│       ├── controller/       # 控制器目录
│       │   └── legacy_api.py # 旧版 API 控制器
│       ├── service/          # 服务层目录
│       │   └── legacy_api.py # 旧版 API 服务
│       ├── model/            # 数据模型目录
│       │   └── legacy.py     # 旧版 API 数据模型
│       ├── utils/            # 工具函数目录
│       │   └── http_client.py # HTTP 客户端
│       └── database/         # 数据库目录
│           ├── db.py         # 数据库连接
│           └── redis.py      # Redis 连接
├── frontend/              # 前端目录
│   ├── package.json      # 前端依赖
│   ├── vite.config.js    # Vite 配置
│   ├── index.html        # HTML 入口
│   ├── src/              # 前端源代码
│   │   ├── main.js       # 前端入口文件
│   │   ├── App.vue       # 根组件
│   │   ├── router/       # 路由配置
│   │   ├── store/        # 状态管理
│   │   ├── components/   # 组件目录
│   │   │   ├── common/   # 通用组件
│   │   │   └── business/ # 业务组件
│   │   ├── views/        # 页面目录
│   │   │   ├── auth/           # 认证页面
│   │   │   │   └── login.vue   # 登录页面
│   │   │   ├── company/        # 集团公司管理
│   │   │   │   ├── index.vue   # 集团公司列表
│   │   │   │   └── detail.vue  # 集团公司详情
│   │   │   ├── sso/            # 身份验证（免登）
│   │   │   │   ├── config.vue  # 免登配置
│   │   │   │   └── test.vue    # 免登测试
│   │   │   ├── access_token/   # accessToken 管理
│   │   │   │   ├── config.vue  # accessToken 配置
│   │   │   │   ├── test.vue    # accessToken 测试
│   │   │   │   └── monitor.vue # accessToken 监控
│   │   │   ├── download/       # 下载页面
│   │   │   │   └── index.vue   # 下载首页
│   │   │   ├── permission/     # 权限管理
│   │   │   │   ├── user.vue    # 用户管理
│   │   │   │   ├── role.vue    # 角色管理
│   │   │   │   └── menu.vue    # 菜单管理
│   │   │   ├── api/            # API 管理
│   │   │   │   ├── config.vue  # API 配置
│   │   │   │   ├── test.vue    # API 测试
│   │   │   │   └── legacy.vue  # 旧版 API 管理
│   │   │   └── system/         # 系统管理
│   │   │       ├── parameter.vue # 参数配置
│   │   │       ├── log.vue       # 日志管理
│   │   │       └── backup.vue     # 备份恢复
│   │   ├── api/          # API 调用
│   │   │   ├── go.js        # Go API 调用
│   │   │   └── python.js    # Python API 调用
│   │   └── utils/        # 工具函数
│   └── public/           # 静态资源目录
└── docs/                 # 文档目录
    ├── development_guide.md  # 开发规范指南
    └── implementation_plan.md  # 实施方案
```

## 5. 开发规范

### 5.1 命名规范

#### Go 后端
- **文件命名**：使用小写字母和下划线组合，如 `downloader.go`
- **包命名**：使用小写字母，如 `config`、`service`
- **函数命名**：使用驼峰命名法，如 `DownloadList`
- **变量命名**：使用驼峰命名法，如 `appKey`
- **常量命名**：使用全大写字母和下划线组合，如 `LOG_FORMAT`

#### Vue 前端
- **文件命名**：组件使用 PascalCase，如 `Downloader.vue`；工具函数使用小写字母和下划线组合，如 `http_client.js`
- **组件命名**：使用 PascalCase，如 `Downloader`
- **函数命名**：使用驼峰命名法，如 `downloadList`
- **变量命名**：使用驼峰命名法，如 `appKey`
- **常量命名**：使用全大写字母和下划线组合，如 `API_BASE_URL`

### 5.2 注释规范

#### Go 后端
- **文件头部**：每个文件必须包含作者、邮箱、时间和功能说明
- **函数注释**：每个导出函数必须包含功能说明、参数、返回值和异常说明
- **代码注释**：复杂代码段必须添加注释，说明代码逻辑

#### Vue 前端
- **文件头部**：每个组件文件必须包含功能说明
- **组件注释**：每个组件必须包含 props、emits、data 等说明
- **函数注释**：复杂函数必须包含功能说明、参数、返回值
- **代码注释**：复杂逻辑必须添加注释

### 5.3 代码风格

#### Go 后端
- 缩进：使用 4 个空格
- 行长度：每行不超过 120 个字符
- 空行：函数之间、逻辑块之间使用空行分隔
- 导入顺序：标准库 → 第三方库 → 本地库
- 使用 `go fmt` 格式化代码

#### Vue 前端
- 缩进：使用 2 个空格
- 行长度：每行不超过 120 个字符
- 空行：组件之间、逻辑块之间使用空行分隔
- 导入顺序：Vue 核心 → 第三方库 → 本地组件 → 样式文件
- 使用 ESLint 和 Prettier 格式化代码

## 6. 开发计划与任务列表

### 6.1 第一阶段：项目初始化（已完成）

- [x] 创建项目目录结构
- [x] 编写 README.md 文档
- [x] 编写开发规范指南
- [x] 编写项目实施方案

### 6.2 第二阶段：基础架构搭建

#### 6.2.1 后端开发（Go）
- [ ] 初始化 Go 模块
  - 创建 go.mod 文件
  - 添加依赖
- [ ] 实现数据库连接
  - 数据库初始化
  - ORM 框架配置
- [ ] 实现工具函数
  - HTTP 客户端
  - 日志配置
  - JWT 工具
  - 加密工具
- [ ] 实现中间件
  - 认证中间件
  - CORS 中间件
  - 日志中间件

#### 6.2.2 前端开发（Vue 3）
- [ ] 初始化 Vue 项目
  - 创建 Vue 项目结构
  - 添加依赖
- [ ] 实现基础组件
  - 布局组件
  - 表格组件
  - 表单组件
- [ ] 实现路由配置
  - 路由定义
  - 导航守卫
- [ ] 实现状态管理
  - Pinia 配置
  - 状态定义

### 6.3 第三阶段：核心功能开发

#### 6.3.1 集团公司管理模块
- [ ] 后端开发
  - 集团公司管理 API
  - 分子公司管理 API
  - 数据隔离实现
  - 多级权限控制
- [ ] 前端开发
  - 集团公司列表页面
  - 集团公司详情页面
  - 分子公司管理页面

#### 6.3.2 身份验证（免登）模块
- [ ] 后端开发
  - 免登配置 API
  - 免登测试 API
  - 免登结果记录 API
- [ ] 前端开发
  - 免登配置页面
  - 免登测试页面

#### 6.3.3 accessToken 获取模块
- [ ] 后端开发
  - accessToken 配置 API
  - accessToken 自动刷新机制
  - accessToken 测试 API
  - accessToken 状态监控 API
- [ ] 前端开发
  - accessToken 配置页面
  - accessToken 测试页面
  - accessToken 状态监控页面

#### 6.3.4 权限管理模块
- [ ] 后端开发
  - 用户管理 API
  - 角色管理 API
  - 菜单管理 API
  - 权限分配 API
  - 操作日志 API
- [ ] 前端开发
  - 用户管理页面
  - 角色管理页面
  - 菜单管理页面
  - 权限分配页面
  - 操作日志页面

#### 6.3.5 API 配置模块
- [ ] 后端开发
  - 钉钉应用配置 API
  - API 版本管理 API
  - 接口参数配置 API
  - 调用频率限制 API
- [ ] 前端开发
  - 钉钉应用配置页面
  - API 版本管理页面
  - 接口参数配置页面
  - 调用频率限制页面

#### 6.3.6 API 测试模块
- [ ] 后端开发
  - 在线 API 调试 API
  - 测试用例管理 API
  - 响应结果查看 API
  - 调用历史记录 API
- [ ] 前端开发
  - 在线 API 调试页面
  - 测试用例管理页面
  - 响应结果查看页面
  - 调用历史记录页面

#### 6.3.7 系统管理模块
- [ ] 后端开发
  - 系统参数配置 API
  - 数据库管理 API
  - 日志管理 API
  - 备份恢复 API
- [ ] 前端开发
  - 系统参数配置页面
  - 数据库管理页面
  - 日志管理页面
  - 备份恢复页面

### 6.4 第四阶段：业务功能开发

#### 6.4.1 下载功能模块
- [ ] 后端开发
  - 列表数据获取 API
  - 数据导出 API
  - 文件下载 API
- [ ] 前端开发
  - 列表数据查询页面
  - 数据导出配置页面
  - 文件下载页面

### 6.5 第五阶段：测试与优化

- [ ] 后端单元测试
  - 服务层测试
  - 控制器测试
  - 中间件测试
- [ ] 前端单元测试
  - 组件测试
  - 工具函数测试
  - 页面测试
- [ ] 集成测试
  - API 测试
  - 端到端测试
  - 权限测试
- [ ] 性能优化
  - 响应时间优化
  - 内存使用优化
  - 数据库查询优化
- [ ] 安全优化
  - 权限控制优化
  - 数据加密优化
  - 防止 SQL 注入

### 6.6 第六阶段：部署与发布

- [ ] 编写部署文档
  - 后端部署
  - 前端部署
  - 环境配置
  - 数据库配置
- [ ] 构建与部署
  - 后端编译
  - 前端打包
  - 部署到服务器
  - 配置 Nginx
- [ ] 发布第一版本
  - 版本号定义
  - 发布日志
  - 用户文档
  - 培训文档

## 7. 预期成果

- 完整的钉钉 API 对接系统
- 包含权限管理、API 配置、API 测试和系统管理等核心功能
- 支持多种文件格式导出（Excel、CSV、JSON）
- 详细的文档和开发规范
- 可直接部署使用的代码
- 友好的用户界面
- 完善的权限控制
- 安全可靠的 API 调用
- 方便的 API 测试功能

## 8. 项目风险

- 钉钉 API 变更可能导致功能失效
- 网络不稳定可能影响数据下载
- 不同版本的依赖库可能存在兼容性问题

## 9. 应对措施

- 定期检查并更新钉钉 SDK
- 实现网络错误重试机制
- 锁定依赖库版本
- 编写详细的异常处理代码

# 任务列表

| 任务 ID | 任务名称 | 优先级 | 状态 | 预计完成时间 |
|---------|----------|--------|------|--------------|
| T001 | 创建项目目录结构 | 高 | 已完成 | 2025-12-20 |
| T002 | 编写 README.md 文档 | 高 | 已完成 | 2025-12-20 |
| T003 | 编写 requirements.txt 文件 | 高 | 已完成 | 2025-12-20 |
| T004 | 编写开发规范指南 | 高 | 已完成 | 2025-12-20 |
| T005 | 实现工具函数模块 | 高 | 待开发 | 2025-12-21 |
| T006 | 实现下载器模块 | 高 | 待开发 | 2025-12-22 |
| T007 | 实现主入口文件 | 高 | 待开发 | 2025-12-23 |
| T008 | 编写单元测试 | 中 | 待开发 | 2025-12-24 |
| T009 | 测试核心功能 | 中 | 待开发 | 2025-12-25 |
| T010 | 优化代码结构 | 中 | 待开发 | 2025-12-26 |
| T011 | 完善文档 | 中 | 待开发 | 2025-12-27 |
| T012 | 编写部署文档 | 低 | 待开发 | 2025-12-28 |
| T013 | 发布第一版本 | 低 | 待开发 | 2025-12-29 |

作者: cjx
邮箱: xx4125517@126.com
时间: 2025-12-20 11:45:00