# 密码策略管理

## 添加密码

```http
POST /api/v1/packages/{package_id}/passwords
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "password": "Password123!",
  "priority": 1,
  "valid_from": "2024-01-01T00:00:00Z",
  "valid_until": "2024-01-07T00:00:00Z"
}
```

**参数说明：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| password | string | 是 | 密码内容 |
| priority | int | 否 | 优先级，1 为最高，默认 1 |
| valid_from | datetime | 否 | 生效时间，不设置则立即生效 |
| valid_until | datetime | 否 | 过期时间，不设置则永不过期 |

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "pwd-uuid-xxx",
    "password": "Password123!",
    "priority": 1,
    "valid_from": "2024-01-01T00:00:00Z",
    "valid_until": "2024-01-07T00:00:00Z",
    "status": "pending",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

**密码规则：**
- 最小长度：`8` 字符（可通过全局配置修改）
- 可使用：字母、数字、特殊字符

## 批量添加密码

```http
POST /api/v1/packages/{package_id}/passwords/batch
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "passwords": [
    {
      "password": "Week1Password!",
      "priority": 1,
      "valid_from": "2024-01-01T00:00:00Z",
      "valid_until": "2024-01-08T00:00:00Z"
    },
    {
      "password": "Week2Password!",
      "priority": 2,
      "valid_from": "2024-01-08T00:00:00Z",
      "valid_until": "2024-01-15T00:00:00Z"
    },
    {
      "password": "Week3Password!",
      "priority": 3,
      "valid_from": "2024-01-15T00:00:00Z",
      "valid_until": "2024-01-22T00:00:00Z"
    }
  ]
}
```

## 获取文件包的密码列表

```http
GET /api/v1/packages/{package_id}/passwords
Authorization: Bearer <jwt_token>
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": "pwd-uuid-1",
      "password": "Password123!",
      "priority": 1,
      "valid_from": "2024-01-01T00:00:00Z",
      "valid_until": "2024-01-07T00:00:00Z",
      "status": "active",
      "created_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "pwd-uuid-2",
      "password": "Password456!",
      "priority": 2,
      "valid_from": "2024-01-07T00:00:00Z",
      "valid_until": "2024-01-14T00:00:00Z",
      "status": "pending",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

## 获取单个密码详情

```http
GET /api/v1/passwords/{password_id}
Authorization: Bearer <jwt_token>
```

## 更新密码策略

```http
PUT /api/v1/passwords/{password_id}
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "password": "NewPassword789!",
  "priority": 1,
  "valid_from": "2024-01-01T00:00:00Z",
  "valid_until": "2024-01-10T00:00:00Z"
}
```

**注意：** 更新密码后，需要重新生成文件包才能生效。

## 删除密码

```http
DELETE /api/v1/passwords/{password_id}
Authorization: Bearer <jwt_token>
```

## 立即激活密码

手动将密码状态设置为激活（跳过有效期检查）。

```http
POST /api/v1/passwords/{password_id}/activate
Authorization: Bearer <jwt_token>
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "pwd-uuid-xxx",
    "status": "active"
  }
}
```

## 立即停用密码

手动将密码状态设置为停用（即使在有效期内）。

```http
POST /api/v1/passwords/{password_id}/deactivate
Authorization: Bearer <jwt_token>
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "pwd-uuid-xxx",
    "status": "disabled"
  }
}
```

## 获取当前有效密码

```http
GET /api/v1/packages/{package_id}/passwords/current
Authorization: Bearer <jwt_token>
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "pwd-uuid-xxx",
    "password": "CurrentPassword!",
    "valid_from": "2024-01-01T00:00:00Z",
    "valid_until": "2024-01-07T00:00:00Z",
    "status": "active"
  }
}
```

## 密码状态说明

| 状态 | 说明 |
|------|------|
| pending | 待生效（未到 valid_from 时间） |
| active | 激活（当前有效） |
| expired | 已过期（超过 valid_until 时间） |
| disabled | 已停用（被管理员手动停用） |

## 状态自动转换规则

```
创建时：
  - 如果未设置 valid_from，默认立即生效 → status: active
  - 如果设置了 valid_from 且在未来 → status: pending

运行时检查：
  - active + 当前时间 > valid_until → status: expired
  - pending + 当前时间 >= valid_from → status: active
```

## 完整使用示例

### 场景：设置每周自动切换密码

```python
import requests
from datetime import datetime, timedelta

API_BASE = "http://localhost:8080/api/v1"
token = "Bearer <your_token>"
headers = {"Authorization": token, "Content-Type": "application/json"}

package_id = "pkg-uuid-xxx"

# 创建 4 周的密码策略
start_date = datetime(2024, 1, 1)
passwords = []

for week in range(4):
    pwd = {
        "password": f"Week{week+1}Password!",
        "priority": week + 1,
        "valid_from": (start_date + timedelta(weeks=week)).isoformat() + "Z",
        "valid_until": (start_date + timedelta(weeks=week+1)).isoformat() + "Z"
    }
    passwords.append(pwd)

# 批量创建
response = requests.post(
    f"{API_BASE}/packages/{package_id}/passwords/batch",
    headers=headers,
    json={"passwords": passwords}
)

print(response.json())
```
