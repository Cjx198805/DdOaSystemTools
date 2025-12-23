# DdOaListDownload 调试说明文档

作者: cjx
邮箱: xx4125517@126.com
时间: 2025-12-23 11:50:00

## 1. 简介
本指南旨在帮助开发人员有效地调试 DdOaListDownload 系统，涵盖前后端及数据库的常见调试方法。

## 2. 环境准备：启动全线服务 (Prerequisites)

在进行任何形式的调试（尤其是人工实机调试）之前，必须确保所有关联服务已在开发环境下正常运行。

### 2.1 启动基础设施 (Infrastructure)
- **MySQL**: 确保服务已启动并能通过配置的账号连接。
- **Redis**: 确保 6379 端口正常监听。

### 2.2 启动 Go 后端 (主业务)
1. **路径**: `cd backend/go`
2. **设置环境变量 (PowerShell)**:
   ```powershell
   $env:JWT_SECRET="ddoalistdownload-debug-secret-key"
   ```
3. **启动**:
   ```bash
   go run main.go
   ```
   *控制台应输出 “MySQL连接初始化成功” 及 Gin 路由注册信息。*

### 2.3 启动 Python 后端 (辅助/旧版 API)
1. **路径**: `cd backend/python`
2. **激活环境并启动**:
   ```bash
   .\venv\Scripts\activate
   python app.py
   ```

### 2.4 启动前端应用
1. **路径**: `cd frontend`
2. **启动开发服务器**:
   ```bash
   npm run dev
   ```
   *默认地址: http://localhost:5173*

## 3. 后端调试 (Go)

### 2.1 VS Code 调试配置
在项目根目录创建 `.vscode/launch.json`:
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Go Server",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceRoot}/backend/go/main.go",
            "cwd": "${workspaceRoot}/backend/go",
            "env": {
                "JWT_SECRET": "your-debug-secret-key"
            }
        }
    ]
}
```

### 2.2 日志调试
系统使用 `logrus` 进行日志记录。
- **控制台日志**: 开发环境下日志会直接输出到控制台。
- **错误堆栈**: 全局错误处理中间件会打印详细的 Panic 堆栈信息，定位代码行号。
- **自定义日志**:
  ```go
  logrus.WithFields(logrus.Fields{
      "user_id": userID,
      "module": "field_permission",
  }).Infof("正在执行调试信息: %s", "some data")
  ```

### 2.3 单元测试调试
运行特定函数的测试并查看详细输出：
```bash
cd backend/go
go test -v ./service -run TestCheckFieldEditable
```

## 4. 前端调试 (Vue 3)

### 3.4 浏览器开发者工具
- **Network 面板**: 监控 API 请求，检查 Request Headers (Authorization Token) 和 Response Body。
- **Console 面板**: 系统会自动捕获并打印接口异常及 LocalStorage 解析错误。

### 3.5 Vue Devtools
安装 Chrome 插件 **Vue.js devtools** (支持 Vue 3)：
- 检查 **Pinia Store** 中的 `user` 状态（token, userInfo）。
- 查看组件的 `props`、`data` 和 `computed` 属性。

### 3.6 环境变量调试
检查 `.env` 文件中的 `VITE_API_BASE_URL` 是否正确指向后端地址。

## 5. 开发环境支持操作 (Dev Ops)

### 4.1 开启 SQL 全日志调试 (GORM Debug)
在调试数据库逻辑时，建议临时开启 GORM 的调试模式，这将在控制台打印每一条执行的 SQL 语句：
- **操作**: 修改 `backend/go/database/db.go` 中的 `gorm.Config`。
- **配置**:
  ```go
  // 将 LogLevel 设置为 logger.Info 且在 GetDB 时调用 .Debug()
  db.Debug().Where(...).Find(...)
  ```

### 4.2 Redis 数据实时监控
系统中的 JWT 状态、异步任务进度存储在 Redis 中。
- **监控命令**:
  ```bash
  redis-cli
  monitor  # 实时查看所有读写操作
  keys *   # 查看当前存储的所有键 (如 ddoa:task:*)
  ```

### 4.3 前端热更新与反向代理
- **Vite 热更新**: 修改 `.vue` 文件后，浏览器会自动局部刷新。如果未生效，请检查 `package.json` 中的 `vite` 版本。
- **接口代理**: 如果遇到跨域问题，检查 `frontend/vite.config.js` 中的 `server.proxy`，确保请求被转发到 `http://localhost:8080`。

### 4.4 快速数据重置 (Clean Test)
在进行复杂权限测试前，如果需要清空并重新初始化数据库：
```bash
# 后端执行
cd backend/go
go run main.go -migrate=true -force=true  # (如果实现了该参数)
# 或者直接在 MySQL 中执行
TRUNCATE TABLE user_role;
TRUNCATE TABLE field_permission;
```

## 6. 接口调试 (API)

### 4.1 内置 API 测试工具
系统在“API 测试”模块提供了真实请求执行器：
- 可以实时修改参数并观察返回结果。
- 自动记录响应时间，便于性能调优。

### 4.2 Postman / cURL
调试需要带上 JWT Token:
```bash
curl -H "Authorization: Bearer <your_token>" http://localhost:8080/api/v1/user/info
```

## 7. 常见调试场景
- **401 Unauthorized**: 检查前端是否成功在请求拦截器中注入了 Token，或者 Token 是否已过期。
- **数据库零值问题**: 检查模型字段是否使用了指针类型（如 `*int`），防止 GORM 忽略 0 值更新。
- **跨域问题**: 检查 Vite 开发服务的 `proxy` 配置。

## 8. 人工操作实机调试场景

为了确保业务链路完整，建议除单元测试外，进行以下人工实机调试：

### 6.1 场景一：权限优先级“特权覆盖”验证
**目标**: 验证 `SpecialEdit` 权限是否能在字典限制的情况下强制开启编辑功能。
1. **人工操作**:
   - 在“数据字典”中将 `user` 模块的 `status` 字段设置为 `Editable = 0` (不可编辑)。
   - 使用普通角色账号登录，进入用户管理，观察 `status` 字段是否已被禁用（禁掉 Input 或 Switch）。
   - 切换到管理员，在“字段权限设置”中，为该普通角色勾选 `status` 字段的“特殊编辑”。
   - 再次以普通角色登录。
2. **预期观察**: 此时该字段应变为“可编辑”状态，证明后端 `CheckFieldEditable` 逻辑生效且前端 UI 响应准确。

### 6.2 场景二：Token 失效与自动拦截验证
**目标**: 验证系统对 401 状态码的全局处理。
1. **人工操作**:
   - 登录系统，保持在任意业务页面。
   - 打开浏览器开发者工具 (F12) -> Application -> LocalStorage。
   - 手动修改 `token` 的值为任意乱码，或者在后端重启服务使 JWT Secret 位移。
   - 在前端点击“筛选”或“刷新”触发 API 请求。
2. **预期观察**: 
   - 页面弹出 `ElMessage.error` 提示登录过期或无效。
   - 页面自动清除本地用户信息并重定向到 `/login` 登录页。

### 6.3 场景三：异步下载任务全链路追踪
**目标**: 验证前端轮询与后端异步执行器的配合。
1. **人工操作**:
   - 进入“任务管理”，发起一个数据量较大的导出请求。
   - 观察“任务进度”页面。
2. **预期观察**:
   - 查看浏览器 Network 面板，确认前端每隔 N 秒发送一次 `GET /download/task/progress` 状态查询。
   - 查看后端控制台日志，确认 `ExecuteTask` 正在并发执行，且 Redis 中的状态从 `pending` 变为 `processing` 最后变为 `completed`。
   - 最终“下载”按钮变为可用状态，点击可获取真实文件。

### 6.4 场景四：API 测试参数合并调试
**目标**: 验证全局配置与实时参数的合并逻辑。
1. **人工操作**:
   - 在“API 配置”中设置一个通用的 `Authorization: AppKey_XXX` Header。
   - 在“测试用例管理”中，针对该 API 建立用例，并覆盖该 Header 为 `Authorization: Case_YYY`。
   - 点击“执行测试”。
2. **预期观察**: 在“测试历史记录”中查看 `Actual Request` 详情，确认发送的 Header 值为用例中的 `Case_YYY`（遵循用例优先级高于全局配置的原则）。
