# 认证授权

## 概述

系统支持两种认证方式：
1. **JWT Token** - 适用于 Web 前端管理界面
2. **API Key** - 适用于 SDK 集成和后端服务调用

## 用户管理

### 注册用户（仅管理员可操作）

```http
POST /api/v1/admin/users
Content-Type: application/json
X-API-Key: <admin_api_key>

{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "SecurePassword123!",
  "role": "operator"
}
```

**角色可选值：** `super_admin`, `admin`, `operator`, `viewer`

### 获取用户列表

```http
GET /api/v1/admin/users?page=1&page_size=20
X-API-Key: <admin_api_key>
```

### 获取当前用户信息

```http
GET /api/v1/auth/me
Authorization: Bearer <jwt_token>
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "uuid-xxx",
    "username": "admin",
    "email": "admin@example.com",
    "role": "super_admin",
    "created_at": "2024-01-01T00:00:00Z",
    "last_login": "2024-01-05T10:30:00Z"
  }
}
```

### 修改密码

```http
PUT /api/v1/auth/password
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "old_password": "old_password",
  "new_password": "new_password"
}
```

## API Key 管理

### 生成 API Key

```http
POST /api/v1/auth/api-keys
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "name": "Production Key",
  "expires_in": 365
}
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "key-xxx",
    "name": "Production Key",
    "api_key": "ptm_sk_a1b2c3d4e5f6...",
    "created_at": "2024-01-01T00:00:00Z",
    "expires_at": "2025-01-01T00:00:00Z"
  }
}
```

**注意：** `api_key` 只在创建时返回一次，请妥善保管。

### 列出 API Keys

```http
GET /api/v1/auth/api-keys
Authorization: Bearer <jwt_token>
```

### 删除 API Key

```http
DELETE /api/v1/auth/api-keys/{key_id}
Authorization: Bearer <jwt_token>
```

## 权限说明

### 权限列表

| 权限代码 | 说明 | super_admin | admin | operator | viewer |
|----------|------|:-----------:|:-----:|:---------:|:------:|
| package:create | 创建文件包 | ✅ | ✅ | ✅ | ❌ |
| package:read | 查看文件包 | ✅ | ✅ | ✅ | ✅ |
| package:update | 更新文件包 | ✅ | ✅ | ✅ | ❌ |
| package:delete | 删除文件包 | ✅ | ✅ | ❌ | ❌ |
| password:manage | 管理密码策略 | ✅ | ✅ | ✅ | ❌ |
| password:activate | 激活/停用密码 | ✅ | ✅ | ✅ | ❌ |
| audit:read | 查看审计日志 | ✅ | ✅ | ❌ | ❌ |
| user:manage | 管理用户 | ✅ | ❌ | ❌ | ❌ |
| encryption:manage | 管理加密配置 | ✅ | ✅ | ❌ | ❌ |

## 错误码

| 错误码 | 说明 |
|--------|------|
| 10002 | 认证失败（token 无效或过期） |
| 10003 | 权限不足 |
| 10004 | 用户不存在 |
| 10005 | 用户已存在 |
| 10006 | 密码错误 |
