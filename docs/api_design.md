# DdOaListDownload API 设计规格

作者: cjx
邮箱: xx4125517@126.com
时间: 2025-12-23 11:05:00

## 1. 基础信息
- **Base URL**: `/api/v1`
- **认证方式**: Bearer Token (JWT)
- **数据格式**: JSON

## 2. 状态码与错误处理
所有响应遵循以下结构：
```json
{
  "code": 200,
  "data": {},
  "message": "success"
}
```
### 核心错误码：
| Code | 说明 | 处理建议 |
| :--- | :--- | :--- |
| 200 | 成功 | - |
| 400 | 请求参数错误 | 检查参数格式或逻辑 |
| 401 | 令牌无效/过期 | 清除本地缓存并重新登录 |
| 403 | 权限不足 | 联系管理员分配权限 |
| 500 | 服务器内部错误 | 联系开发人员，查看日志挂载的 TraceID |

## 3. 核心 API 概览

### 3.1 身份认证 (Auth)
- `POST /user/login`: 用户登录，返回 Token 及基础信息。
- `GET /user/info`: 获取当前登录用户的详尽配置。

### 3.2 权限管理 (Permission)
- `GET /field-permission`: 获取字段权限列表（支持按公司、角色、模块筛选）。
- `PUT /field-permission/:id`: 更新细粒度权限，包括 `special_edit` 开关。

### 3.3 API 测试与执行 (API Test)
- `POST /api-test/run/:id`: 执行指定测试用例。后端将实时抓取 APIConfig 并通过 HTTP 客户端发送真实请求。
- `GET /api-test/history`: 检索执行历史及耗时统计。

### 3.4 下载任务 (Download)
- `POST /download/task`: 创建新导出任务。
- `GET /download/task/list`: 根据用户角色获取任务列表（Admin 可看全量，普通用户仅看自身）。
