# DdOaListDownload API 设计规格

作者: cjx
邮箱: xx4125517@126.com
最后更新: 2025-12-23 16:00:00

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

## 3. 详细 API 定义

### 3.1 权限管理 (Field Permission)

#### GET /field-permission
获取字段权限列表。

**Query 参数**:
| 参数名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页数量，默认 10 |
| role_id | int | 否 | 按角色筛选 |
| module | string | 否 | 按模块名筛选 (e.g., 'user', 'company') |
| field | string | 否 | 按字段名筛选 |

#### PUT /field-permission/:id
更新字段权限配置。

**Request Body**:
| 字段名 | 类型 | 说明 |
| :--- | :--- | :--- |
| viewable | int | 是否可见 (0/1) |
| editable | int | 是否可编辑 (0/1) |
| special_edit | int | **特殊编辑权限** (0/1)，优先级高于数据字典 |

### 3.2 API 测试 (API Test)

#### POST /api-test/run/:id
执行指定的 API 测试用例。

**逻辑说明**:
后端会读取用例关联的 `APIConfig`，结合用例中的 `Params` 和 `Headers`，发起真实的 HTTP 请求，并将结果写入 `APITestHistory`。

**Response Data**:
```json
{
  "id": 123,
  "status": "success",
  "status_code": 200,
  "response_time": 150, // ms
  "actual_result": "{...}" // 实际响应体
}
```

#### GET /api-test/history
获取测试执行历史。

**Query 参数**:
| 参数名 | 类型 | 说明 |
| :--- | :--- | :--- |
| test_case_id | int | 按用例 ID 筛选 |
| status | string | 按状态筛选 ('success', 'failed') |

### 3.3 下载任务 (Download Task)

#### POST /download-task
创建新的下载/导出任务。

**Request Body**:
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| name | string | 是 | 任务名称 |
| type | string | 是 | 导出类型 ('excel', 'csv', 'json') |
| params | json | 否 | 导出筛选参数 |

#### GET /download-task
获取下载任务列表。

**权限逻辑**:
- **Admin**: 可查看所有用户的任务。
- **普通用户**: 仅可查看自己创建的任务。

**Response Model**:
```json
{
  "tasks": [
    {
      "id": 1,
      "name": "用户导出_20251223",
      "status": "completed", // pending, running, completed, failed
      "progress": 100,
      "file_url": "/api/v1/download/result/1",
      "user_name": "admin"
    }
  ],
  "total": 50
}
```
