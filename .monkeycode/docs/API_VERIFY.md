# 验证服务

## 概述

验证服务是客户端验证程序调用的核心接口，用于验证密码是否有效并获取解密密钥。

## 验证密码

```http
POST /api/v1/verify
Content-Type: application/json

{
  "package_id": "pkg-uuid-xxx",
  "password": "UserEnteredPassword"
}
```

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| package_id | string | 是 | 文件包 ID |
| password | string | 是 | 用户输入的密码 |

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "valid": true,
    "key": "base64_encoded_decryption_key",
    "expires_at": "2024-01-01T12:00:00Z"
  }
}
```

**失败响应（密码错误）：**
```json
{
  "code": 30001,
  "message": "password is incorrect",
  "data": {
    "valid": false,
    "remaining_attempts": 2
  }
}
```

**失败响应（密码已过期）：**
```json
{
  "code": 30002,
  "message": "password has expired",
  "data": {
    "valid": false
  }
}
```

**失败响应（密码未生效）：**
```json
{
  "code": 30003,
  "message": "password is not yet valid",
  "data": {
    "valid": false,
    "valid_from": "2024-01-07T00:00:00Z"
  }
}
```

## 批量验证（尝试验证多个密码）

```http
POST /api/v1/verify/batch
Content-Type: application/json

{
  "package_id": "pkg-uuid-xxx",
  "passwords": ["pwd1", "pwd2", "pwd3"]
}
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "valid": true,
    "matched_password": "pwd2",
    "key": "base64_encoded_decryption_key"
  }
}
```

## 检查文件包状态

获取文件包的当前状态，用于验证程序启动时检查。

```http
GET /api/v1/verify/status/{package_id}
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "package_id": "pkg-uuid-xxx",
    "status": "active",
    "current_password_count": 2,
    "next_password_change": "2024-01-07T00:00:00Z",
    "offline_mode_available": true
  }
}
```

**状态说明：**

| 状态 | 说明 |
|------|------|
| active | 正常，可验证 |
| archived | 已归档，不可验证 |
| no_active_password | 没有激活的密码 |

## 获取离线验证配置

获取用于离线验证的加密配置和内置策略（用于网络异常时的兜底验证）。

```http
GET /api/v1/verify/offline-config/{package_id}
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "version": 1,
    "algorithm": "AES-256-GCM",
    "encrypted_config": "base64_encrypted_config_data",
    "signature": "base64_rsa_signature",
    "policies": [
      {
        "password_hash": "sha256_of_password",
        "valid_from": "2024-01-01T00:00:00Z",
        "valid_until": "2024-01-07T00:00:00Z",
        "key_encrypted": "encrypted_key_for_this_password"
      }
    ]
  }
}
```

**说明：**
- `password` 字段存储的是密码的哈希值（SHA-256），而不是明文
- `key_encrypted` 是使用该密码对应的密钥加密后的数据
- 验证程序通过验证输入密码的哈希值来确认密码正确性

## 客户端验证流程

```
客户端验证程序工作流程：

1. 启动时（可选）
   GET /api/v1/verify/status/{package_id}
   └── 检查包状态和网络可用性

2. 在线验证（优先）
   POST /api/v1/verify
   ├── body: {"package_id": xxx, "password": "用户输入"}
   └── 成功：获得解密密钥

3. 离线验证（兜底）
   ├── 读取内置的 encrypted_config
   ├── 使用 RSA 公钥验证签名
   ├── 解密 config 获取策略
   ├── 计算输入密码的 SHA-256 与内置哈希比对
   └── 验证通过后获取解密密钥

4. 解密文件
   └── 使用获取的密钥解密并解压文件
```

## 验证限制

| 项目 | 限制 |
|------|------|
| 单 IP 验证频率 | 10 次/分钟 |
| 单包每日验证次数 | 1000 次 |
| 离线配置有效期 | 与密码有效期一致 |

## 错误码

| 错误码 | 说明 |
|--------|------|
| 30001 | 密码错误 |
| 30002 | 密码已过期 |
| 30003 | 密码未生效 |
| 30004 | 没有有效的密码（全部过期或被禁用） |
| 20001 | 文件包不存在 |
| 20003 | 文件包已归档 |

## SDK 使用示例

### Python SDK

```python
from ptm_sdk import PTMClient

client = PTMClient(api_key="your_api_key")

# 验证密码
result = client.verify(
    package_id="pkg-uuid-xxx",
    password="user_input_password"
)

if result.valid:
    print(f"验证成功，解密密钥: {result.key}")
else:
    print(f"验证失败: {result.error_message}")
```

### Java SDK

```java
import com.ptm.client.PTMClient;

PTMClient client = new PTMClient("your_api_key");

VerifyResult result = client.verify("pkg-uuid-xxx", "user_input_password");

if (result.isValid()) {
    System.out.println("验证成功，解密密钥: " + result.getKey());
} else {
    System.out.println("验证失败: " + result.getErrorMessage());
}
```

### Go SDK

```go
import "github.com/ptm/go-sdk"

client := ptm.NewClient("your_api_key")

result, err := client.Verify("pkg-uuid-xxx", "user_input_password")
if err != nil {
    log.Fatal(err)
}

if result.Valid {
    fmt.Println("验证成功，解密密钥:", result.Key)
} else {
    fmt.Println("验证失败:", result.Message)
}
```
