# DdOaListDownload 前端应用开发文档

## 项目概述

DdOaListDownload 是一套完整的钉钉 API 对接系统，前端采用 Vue 3 + Vite 开发，支持集团公司多级管理架构，提供权限管理、API 测试、下载任务管理等核心功能。

## 技术栈

- **框架**: Vue 3 (Composition API)
- **构建工具**: Vite
- **状态管理**: Pinia
- **路由**: Vue Router
- **HTTP 客户端**: Axios
- **CSS 预处理器**: 原生 CSS (支持 CSS 变量和模块化)

## 项目结构

```
frontend/
├── src/
│   ├── api/              # API 调用模块
│   ├── components/       # 通用组件
│   ├── views/            # 页面组件
│   │   ├── auth/         # 认证相关页面
│   │   ├── permission/   # 权限管理页面
│   │   ├── api/          # API 测试页面
│   │   └── download/     # 下载任务页面
│   ├── router/           # 路由配置
│   ├── store/            # 状态管理
│   ├── utils/            # 工具函数
│   ├── App.vue           # 根组件
│   └── main.js           # 入口文件
├── index.html            # HTML 模板
├── vite.config.js        # Vite 配置
└── package.json         # 项目依赖
```

## 功能模块

### 1. 权限管理

#### 用户管理
- 用户列表展示
- 用户创建、编辑、删除
- 多角色分配
- 用户状态管理

#### 角色管理
- 角色列表展示
- 角色创建、编辑、删除
- 菜单权限配置

#### 字段权限设置
- 字段权限列表
- 字段权限编辑
- 角色字段权限管理

#### 数据字典管理
- 数据字典列表
- 字典项管理
- 字典状态管理

### 2. API 测试

#### API 测试
- 支持多种 HTTP 方法（GET、POST、PUT、DELETE）
- 动态请求参数配置
- 请求头管理
- 请求体编辑
- 实时测试结果展示

#### 测试用例管理
- 测试用例列表
- 测试用例创建、编辑、删除
- 测试用例执行

#### 测试历史记录
- 测试历史列表
- 测试结果查看
- 历史记录搜索

### 3. 下载任务

#### 任务管理
- 任务列表展示
- 任务创建
- 任务状态管理

#### 任务进度
- 实时进度展示
- 进度更新机制

#### 结果下载
- 支持多种格式下载（Excel、CSV、JSON）
- 下载历史记录

## 开发流程

### 1. 环境搭建

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build

# 预览生产构建
npm run preview
```

### 2. 代码规范

- 使用 Vue 3 Composition API
- 组件设计遵循单一职责原则
- 状态管理使用 Pinia
- 路由配置使用 Vue Router
- API 调用封装在单独的模块中
- CSS 使用模块化和 CSS 变量

### 3. API 集成

所有 API 调用都封装在 `src/api` 目录下，每个模块对应一个 API 文件。API 调用使用 Axios，配置了请求拦截器和响应拦截器，统一处理认证和错误。

### 4. 权限控制

- 路由权限控制：通过路由守卫实现
- 页面权限控制：通过组件内的权限判断实现
- 字段权限控制：通过字段级权限配置实现

## API 集成文档

### 1. API 基础配置

```javascript
// src/api/index.js
import axios from 'axios'

const service = axios.create({
  baseURL: '/api',
  timeout: 10000
})

// 请求拦截器
service.interceptors.request.use(
  config => {
    // 添加认证 token
    const token = localStorage.getItem('token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  response => {
    const res = response.data
    if (res.code !== 200) {
      // 处理错误
      return Promise.reject(new Error(res.message || 'Error'))
    } else {
      return res
    }
  },
  error => {
    return Promise.reject(error)
  }
)
```

### 2. API 模块

| 模块 | 文件 | 主要功能 |
| --- | --- | --- |
| 认证 | auth.js | 登录、退出、获取用户信息 |
| 用户管理 | user.js | 用户列表、创建、编辑、删除 |
| 角色管理 | role.js | 角色列表、创建、编辑、删除、权限配置 |
| 字段权限 | fieldPermission.js | 字段权限列表、编辑、角色字段权限管理 |
| 数据字典 | dataDictionary.js | 数据字典列表、创建、编辑、删除 |
| API 测试 | apiTest.js | 测试执行、测试用例管理、测试历史 |
| 下载任务 | download.js | 任务列表、创建、执行、下载结果 |

## 功能使用说明

### 1. 权限管理

#### 用户管理
1. 进入「权限管理」->「用户管理」页面
2. 点击「添加用户」按钮创建新用户
3. 填写用户信息并分配角色
4. 点击「保存」按钮完成创建
5. 在用户列表中可以编辑或删除用户

#### 角色管理
1. 进入「权限管理」->「角色管理」页面
2. 点击「添加角色」按钮创建新角色
3. 填写角色信息
4. 点击「权限设置」按钮配置菜单权限
5. 在角色列表中可以编辑或删除角色

### 2. API 测试

#### 执行 API 测试
1. 进入「API 测试」->「API 测试」页面
2. 选择 HTTP 方法（GET、POST、PUT、DELETE）
3. 输入 API URL
4. 配置请求参数和请求头
5. 输入请求体（如果需要）
6. 点击「运行测试」按钮执行测试
7. 查看测试结果

#### 管理测试用例
1. 进入「API 测试」->「测试用例管理」页面
2. 点击「添加测试用例」按钮创建新用例
3. 填写测试用例信息
4. 在测试用例列表中可以执行、编辑或删除用例

### 3. 下载任务

#### 创建下载任务
1. 进入「下载任务」->「任务管理」页面
2. 点击「创建新任务」按钮
3. 填写任务信息
4. 选择下载格式
5. 点击「提交」按钮创建任务

#### 查看任务进度
1. 在任务列表中查看任务进度
2. 已完成的任务可以点击「下载」按钮下载结果
3. 可以删除已完成的任务

## 开发规范

### 1. 组件开发

- 组件命名：PascalCase（如 `UserManagement.vue`）
- 组件Props：使用类型检查
- 组件事件：使用明确的事件名称
- 组件样式：使用 scoped 样式

### 2. 代码风格

- 缩进：2 个空格
- 变量命名：camelCase
- 常量命名：UPPER_CASE_WITH_UNDERSCORES
- 函数命名：camelCase

### 3. 注释规范

- 组件注释：包含组件功能、Props、Events
- 函数注释：包含功能、参数、返回值
- 复杂逻辑注释：解释代码逻辑

## 部署说明

1. 构建生产版本：
   ```bash
   npm run build
   ```

2. 将构建产物（`dist` 目录）部署到 Web 服务器

3. 配置 Nginx 或其他 Web 服务器：
   ```nginx
   server {
     listen 80;
     server_name example.com;
     root /path/to/dist;
     index index.html;
     
     location / {
       try_files $uri $uri/ /index.html;
     }
     
     location /api {
       proxy_pass http://localhost:8080;
       proxy_set_header Host $host;
       proxy_set_header X-Real-IP $remote_addr;
     }
   }
   ```

## 版本历史

- v1.0.0 - 初始版本
  - 实现权限管理功能
  - 实现 API 测试功能
  - 实现下载任务管理功能
  - 支持响应式设计
  - 集成后端 API

## 作者信息

- 作者: cjx
- 邮箱: xx4125517@126.com
- 时间: 2025-12-22
