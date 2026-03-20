# 审计日志

## 概述

系统记录所有关键操作的审计日志，包括文件包下载、密码验证、管理员操作等。

## 动作类型

| 类型 | 说明 |
|------|------|
| DOWNLOAD | 文件包下载 |
| VERIFY_SUCCESS | 密码验证成功 |
| VERIFY_FAIL | 密码验证失败 |
| POLICY_CREATE | 创建密码策略 |
| POLICY_UPDATE | 更新密码策略 |
| POLICY_DELETE | 删除密码策略 |
| POLICY_ACTIVATE | 激活密码 |
| POLICY_DEACTIVATE | 停用密码 |
| PACKAGE_CREATE | 创建文件包 |
| PACKAGE_UPDATE | 更新文件包 |
| PACKAGE_DELETE | 删除文件包 |
| USER_LOGIN | 用户登录 |
| USER_LOGOUT | 用户登出 |
| USER_CREATE | 创建用户 |
| USER_UPDATE | 更新用户 |
| USER_DELETE | 删除用户 |
| CONFIG_UPDATE | 修改系统配置 |

## 查询审计日志

```http
GET /api/v1/audit?page=1&page_size=20
Authorization: Bearer <jwt_token>
```

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| page | int | 页码，默认 1 |
| page_size | int | 每页数量，默认 20 |
| action | string | 筛选动作类型 |
| package_id | string | 筛选文件包 |
| user_id | string | 筛选操作用户 |
| start_time | datetime | 筛选开始时间 |
| end_time | datetime | 筛选结束时间 |
| ip_address | string | 筛选 IP 地址 |

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "audit-uuid-xxx",
        "action": "VERIFY_SUCCESS",
        "package_id": "pkg-uuid-xxx",
        "package_name": "敏感文档",
        "user_id": "user-xxx",
        "username": "external_user",
        "ip_address": "192.168.1.100",
        "user_agent": "PTM-Client/1.0",
        "detail": {
          "password_id": "pwd-uuid-xxx",
          "password_priority": 1
        },
        "created_at": "2024-01-05T10:30:00Z"
      },
      {
        "id": "audit-uuid-yyy",
        "action": "DOWNLOAD",
        "package_id": "pkg-uuid-xxx",
        "package_name": "敏感文档",
        "user_id": "admin",
        "username": "admin",
        "ip_address": "192.168.1.10",
        "user_agent": "Mozilla/5.0...",
        "detail": {
          "format": "exe",
          "file_size": 1048576
        },
        "created_at": "2024-01-05T09:15:00Z"
      }
    ],
    "total": 150,
    "page": 1,
    "page_size": 20,
    "total_pages": 8
  }
}
```

## 获取文件包的审计记录

```http
GET /api/v1/audit/package/{package_id}?page=1&page_size=20
Authorization: Bearer <jwt_token>
```

## 获取用户的审计记录

```http
GET /api/v1/audit/user/{user_id}?page=1&page_size=20
Authorization: Bearer <jwt_token>
```

## 获取验证失败记录

```http
GET /api/v1/audit/verify-fails?page=1&page_size=20
Authorization: Bearer <jwt_token>
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "audit-uuid-zzz",
        "action": "VERIFY_FAIL",
        "package_id": "pkg-uuid-xxx",
        "ip_address": "192.168.1.100",
        "user_agent": "PTM-Client/1.0",
        "detail": {
          "attempted_password": "wrong_password",
          "remaining_attempts": 2
        },
        "created_at": "2024-01-05T10:35:00Z"
      }
    ],
    "total": 10,
    "page": 1,
    "page_size": 20,
    "total_pages": 1
  }
}
```

## 导出审计日志

将审计日志导出为 CSV 或 JSON 格式。

```http
GET /api/v1/audit/export?format=csv&start_time=2024-01-01T00:00:00Z&end_time=2024-01-31T23:59:59Z
Authorization: Bearer <jwt_token>
```

**查询参数：**

| 参数 | 类型 | 说明 |
|------|------|------|
| format | string | 导出格式：`csv` 或 `json` |
| start_time | datetime | 开始时间 |
| end_time | datetime | 结束时间 |
| package_id | string | 筛选文件包（可选） |
| action | string | 筛选动作类型（可选） |

**响应：**

- `format=csv`：返回 CSV 文件下载
- `format=json`：返回 JSON 格式数据

## 审计日志保留策略

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| retention_days | 365 | 日志保留天数 |
| archive_enabled | true | 是否归档旧日志 |
| archive_after_days | 90 | 超过此天数后归档 |

## detail 字段说明

根据不同的 `action` 类型，`detail` 字段包含不同的详细信息：

### VERIFY_SUCCESS / VERIFY_FAIL

```json
{
  "password_id": "pwd-uuid-xxx",
  "password_priority": 1,
  "verify_mode": "online"
}
```

### DOWNLOAD

```json
{
  "format": "exe",
  "file_size": 1048576,
  "download_url_type": "direct"
}
```

### POLICY_CREATE / POLICY_UPDATE

```json
{
  "policy_id": "pwd-uuid-xxx",
  "changes": {
    "password": "***",
    "valid_until": "2024-01-14T00:00:00Z"
  }
}
```

### USER_LOGIN

```json
{
  "login_type": "password",
  "mfa_used": false
}
```
