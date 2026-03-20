# 文件包管理

## 创建文件包

### 上传文件创建

```http
POST /api/v1/packages
Authorization: Bearer <jwt_token>
Content-Type: multipart/form-data

file: <binary file>
name: "敏感文档"
format: "exe"
```

**参数说明：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | binary | 是 | 上传的文件（最大 500MB） |
| name | string | 是 | 文件包名称 |
| format | string | 是 | 输出格式：`exe` 或 `zip` |
| description | string | 否 | 文件包描述 |

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "pkg-uuid-xxx",
    "name": "敏感文档",
    "format": "exe",
    "status": "active",
    "file_size": 1048576,
    "file_hash": "sha256:abc123...",
    "created_at": "2024-01-01T00:00:00Z",
    "created_by": "admin"
  }
}
```

### 通过 URL 创建

```http
POST /api/v1/packages
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "name": "远程文件",
  "format": "zip",
  "source_url": "https://example.com/file.pdf",
  "description": "从远程 URL 导入"
}
```

## 获取文件包列表

```http
GET /api/v1/packages?page=1&page_size=20&status=active
Authorization: Bearer <jwt_token>
```

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| page | int | 页码，默认 1 |
| page_size | int | 每页数量，默认 20 |
| status | string | 筛选状态：`active`, `archived` |
| format | string | 筛选格式：`exe`, `zip` |
| created_by | string | 筛选创建人 |

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "pkg-uuid-xxx",
        "name": "敏感文档",
        "format": "exe",
        "status": "active",
        "file_size": 1048576,
        "created_at": "2024-01-01T00:00:00Z",
        "password_count": 3,
        "current_password": "pwd1"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20,
    "total_pages": 3
  }
}
```

## 获取文件包详情

```http
GET /api/v1/packages/{package_id}
Authorization: Bearer <jwt_token>
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "pkg-uuid-xxx",
    "name": "敏感文档",
    "format": "exe",
    "status": "active",
    "description": "这是一份敏感文档",
    "file_size": 1048576,
    "file_hash": "sha256:abc123...",
    "file_path": "/storage/packages/pkg-uuid-xxx.exe",
    "created_at": "2024-01-01T00:00:00Z",
    "created_by": "admin",
    "updated_at": "2024-01-02T00:00:00Z",
    "passwords": [
      {
        "id": "pwd-uuid-1",
        "password": "pwd1",
        "priority": 1,
        "valid_from": "2024-01-01T00:00:00Z",
        "valid_until": "2024-01-07T00:00:00Z",
        "status": "active"
      },
      {
        "id": "pwd-uuid-2",
        "password": "pwd2",
        "priority": 2,
        "valid_from": "2024-01-07T00:00:00Z",
        "valid_until": "2024-01-14T00:00:00Z",
        "status": "pending"
      }
    ]
  }
}
```

## 更新文件包

```http
PUT /api/v1/packages/{package_id}
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "name": "更新的名称",
  "description": "更新的描述",
  "status": "archived"
}
```

**注意：** 不允许通过此接口更新文件内容，如需更换文件请删除后重新创建。

## 删除文件包

```http
DELETE /api/v1/packages/{package_id}
Authorization: Bearer <jwt_token>
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

**注意：** 删除操作会同时删除关联的密码策略和生成的下载文件。

## 下载文件包

```http
GET /api/v1/packages/{package_id}/download
Authorization: Bearer <jwt_token>
```

**响应：**

- Headers 中包含 `Content-Disposition` 指定文件名
- 文件内容作为二进制流返回

**Python SDK 示例：**
```python
import requests

response = requests.get(
    "http://localhost:8080/api/v1/packages/pkg-uuid-xxx/download",
    headers={"Authorization": "Bearer <token>"},
    stream=True
)

with open("output.exe", "wb") as f:
    for chunk in response.iter_content(chunk_size=8192):
        f.write(chunk)
```

## 获取下载链接（生成临时 URL）

```http
POST /api/v1/packages/{package_id}/download-url
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "expires_in": 3600
}
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "download_url": "http://localhost:8080/api/v1/packages/pkg-uuid-xxx/download?token=xxx",
    "expires_at": "2024-01-01T12:00:00Z"
  }
}
```
