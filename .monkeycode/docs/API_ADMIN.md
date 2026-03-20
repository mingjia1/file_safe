# 系统管理

## 全局加密配置

### 获取当前加密配置

```http
GET /api/v1/admin/encryption/config
Authorization: Bearer <jwt_token>
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "aes_key_length": 256,
    "rsa_key_length": 2048,
    "password_min_length": 8,
    "password_require_special": true,
    "password_require_uppercase": true,
    "password_require_lowercase": true,
    "password_require_digit": true,
    "config_encrypt": true,
    "enable_signature": true,
    "updated_at": "2024-01-01T00:00:00Z",
    "updated_by": "admin"
  }
}
```

### 更新加密配置

```http
PUT /api/v1/admin/encryption/config
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "aes_key_length": 256,
  "rsa_key_length": 4096,
  "password_min_length": 12,
  "password_require_special": true,
  "password_require_uppercase": true,
  "password_require_lowercase": true,
  "password_require_digit": true,
  "config_encrypt": true,
  "enable_signature": true
}
```

**参数说明：**

| 参数 | 类型 | 说明 | 可选值 |
|------|------|------|--------|
| aes_key_length | int | AES 密钥长度 | 128, 192, 256 |
| rsa_key_length | int | RSA 密钥长度 | 1024, 2048, 4096 |
| password_min_length | int | 密码最小长度 | 4-64 |
| password_require_special | bool | 需要特殊字符 | true/false |
| password_require_uppercase | bool | 需要大写字母 | true/false |
| password_require_lowercase | bool | 需要小写字母 | true/false |
| password_require_digit | bool | 需要数字 | true/false |
| config_encrypt | bool | 加密配置文件 | true/false |
| enable_signature | bool | 启用签名验证 | true/false |

**注意：** 修改加密配置后，之后生成的文件包会使用新的配置，已生成的文件包不受影响。

### 验证加密配置有效性

在保存配置之前，测试配置是否有效。

```http
POST /api/v1/admin/encryption/validate
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "aes_key_length": 256,
  "rsa_key_length": 4096,
  "password_min_length": 12,
  "password_require_special": true,
  "config_encrypt": true,
  "enable_signature": true
}
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "valid": true,
    "tests": [
      {
        "name": "AES-256 加密/解密",
        "status": "PASS",
        "duration_ms": 5
      },
      {
        "name": "RSA-4096 签名/验签",
        "status": "PASS",
        "duration_ms": 120
      },
      {
        "name": "密码强度验证",
        "status": "PASS",
        "test_cases": [
          {"password": "Weak", "result": "FAIL"},
          {"password": "StrongPass123!", "result": "PASS"}
        ]
      },
      {
        "name": "配置加密/解密",
        "status": "PASS",
        "duration_ms": 3
      }
    ],
    "message": "所有测试通过"
  }
}
```

**测试失败响应：**
```json
{
  "code": 40001,
  "message": "encryption configuration is invalid",
  "data": {
    "valid": false,
    "tests": [
      {
        "name": "RSA-4096 签名/验签",
        "status": "FAIL",
        "error": "key generation failed: insufficient entropy"
      }
    ]
  }
}
```

### 重新生成加密密钥

重新生成 RSA 密钥对（用于签名验证）。

```http
POST /api/v1/admin/encryption/regenerate-keys
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "reason": "安全轮转"
}
```

**警告：** 此操作会导致使用旧密钥生成的文件包无法进行离线验证。请确保在执行此操作前已了解影响范围。

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "old_key_expires_at": "2024-01-15T00:00:00Z",
    "new_key_active_at": "2024-01-01T12:00:00Z",
    "message": "旧密钥将在 14 天后失效，请在失效前重新生成相关文件包"
  }
}
```

## 系统设置

### 获取系统信息

```http
GET /api/v1/admin/system
Authorization: Bearer <jwt_token>
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "version": "1.0.0",
    "environment": "production",
    "database": {
      "type": "postgresql",
      "version": "15.0",
      "host": "localhost"
    },
    "redis": {
      "version": "7.0",
      "connected": true
    },
    "storage": {
      "type": "local",
      "used_bytes": 1073741824,
      "total_bytes": 10737418240
    }
  }
}
```

### 获取系统配置

```http
GET /api/v1/admin/system/config
Authorization: Bearer <jwt_token>
```

### 更新系统配置

```http
PUT /api/v1/admin/system/config
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "verify_rate_limit": 10,
  "verify_daily_limit": 1000,
  "file_max_size": 524288000,
  "log_retention_days": 365,
  "offline_config_validity_days": 30
}
```

## 角色管理

### 获取角色列表

```http
GET /api/v1/admin/roles
Authorization: Bearer <jwt_token>
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": "role-super-admin",
      "name": "super_admin",
      "description": "超级管理员",
      "permissions": ["*"],
      "user_count": 1
    },
    {
      "id": "role-admin",
      "name": "admin",
      "description": "管理员",
      "permissions": [
        "package:create",
        "package:read",
        "package:update",
        "password:manage",
        "password:activate",
        "audit:read",
        "encryption:manage"
      ],
      "user_count": 3
    }
  ]
}
```

## 统计数据

### 获取仪表盘数据

```http
GET /api/v1/admin/stats/dashboard
Authorization: Bearer <jwt_token>
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_packages": 150,
    "active_packages": 120,
    "total_passwords": 450,
    "active_passwords": 180,
    "total_downloads": 3200,
    "total_verifies": 15000,
    "verify_success_rate": 0.95,
    "recent_activity": [
      {
        "action": "DOWNLOAD",
        "package_name": "财务报告 Q4",
        "user": "zhangsan",
        "created_at": "2024-01-05T10:30:00Z"
      }
    ]
  }
}
```

### 获取使用趋势

```http
GET /api/v1/admin/stats/trends?period=30d
Authorization: Bearer <jwt_token>
```

**period 可选值：** `7d`, `30d`, `90d`, `1y`

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period": "30d",
    "downloads": [
      {"date": "2024-01-01", "count": 45},
      {"date": "2024-01-02", "count": 62}
    ],
    "verifies": [
      {"date": "2024-01-01", "success": 100, "fail": 5},
      {"date": "2024-01-02", "success": 120, "fail": 8}
    ]
  }
}
```
