# 技术设计文档 - Password Timer Manager

## 1. 技术栈

| 层级 | 技术选型 |
|------|----------|
| 后端 | Go 1.21+ / Gin 框架 |
| 前端 | React 18+ / TypeScript / Vite |
| 数据库 | PostgreSQL 15+ / GORM |
| 缓存 | Redis 7+ |
| 文件存储 | 本地文件系统 / S3 兼容 |
| 容器化 | Docker / Docker Compose |

---

## 2. 项目结构

```
password-timer-manager/
├── backend/                     # 后端服务
│   ├── cmd/
│   │   └── server/              # 主入口
│   │       └── main.go
│   ├── internal/
│   │   ├── api/                # API 层
│   │   │   ├── handler/        # 请求处理
│   │   │   │   ├── auth.go
│   │   │   │   ├── package.go
│   │   │   │   ├── password.go
│   │   │   │   ├── verify.go
│   │   │   │   ├── audit.go
│   │   │   │   └── admin.go
│   │   │   ├── middleware/     # 中间件
│   │   │   │   ├── auth.go
│   │   │   │   ├── cors.go
│   │   │   │   └── ratelimit.go
│   │   │   └── router.go       # 路由定义
│   │   ├── model/              # 数据模型
│   │   │   ├── user.go
│   │   │   ├── package.go
│   │   │   ├── password.go
│   │   │   ├── audit.go
│   │   │   └── config.go
│   │   ├── repository/          # 数据访问层
│   │   │   ├── user.go
│   │   │   ├── package.go
│   │   │   ├── password.go
│   │   │   └── audit.go
│   │   ├── service/            # 业务逻辑层
│   │   │   ├── auth.go
│   │   │   ├── package.go
│   │   │   ├── password.go
│   │   │   ├── verify.go
│   │   │   ├── audit.go
│   │   │   └── crypto.go
│   │   └── pkg/
│   │       ├── response/        # 统一响应
│   │       └── errors/         # 错误定义
│   ├── migrations/             # 数据库迁移
│   ├── storage/               # 文件存储目录
│   └── go.mod
│
├── frontend/                    # Web 管理平台
│   ├── src/
│   │   ├── api/               # API 调用
│   │   ├── components/        # 通用组件
│   │   ├── pages/             # 页面
│   │   ├── stores/            # 状态管理
│   │   ├── types/             # TypeScript 类型
│   │   └── App.tsx
│   ├── package.json
│   └── vite.config.ts
│
├── generator/                   # 自解压包生成器
│   ├── cmd/
│   │   └── generator/
│   │       └── main.go
│   ├── internal/
│   │   ├── exe/               # EXE 生成
│   │   ├── zip/               # ZIP 生成
│   │   ├── crypto/            # 加密/签名
│   │   └── verifier/          # 验证程序
│   └── go.mod
│
├── sdk/                         # 多语言 SDK
│   ├── python/
│   ├── java/
│   └── go/
│
└── docs/                        # 文档
```

---

## 3. API 路由设计

### 3.1 认证相关

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/auth/login | 用户登录 |
| POST | /api/v1/auth/logout | 用户登出 |
| GET | /api/v1/auth/me | 获取当前用户 |
| PUT | /api/v1/auth/password | 修改密码 |
| POST | /api/v1/auth/api-keys | 生成 API Key |
| GET | /api/v1/auth/api-keys | 列表 API Keys |
| DELETE | /api/v1/auth/api-keys/{id} | 删除 API Key |

### 3.2 文件包管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/packages | 创建文件包 |
| GET | /api/v1/packages | 列表文件包 |
| GET | /api/v1/packages/{id} | 获取详情 |
| PUT | /api/v1/packages/{id} | 更新文件包 |
| DELETE | /api/v1/packages/{id} | 删除文件包 |
| GET | /api/v1/packages/{id}/download | 下载文件包 |
| POST | /api/v1/packages/{id}/download-url | 生成下载链接 |

### 3.3 密码策略

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/packages/{id}/passwords | 添加密码 |
| GET | /api/v1/packages/{id}/passwords | 列表密码 |
| POST | /api/v1/packages/{id}/passwords/batch | 批量添加 |
| GET | /api/v1/passwords/{pid} | 获取密码详情 |
| PUT | /api/v1/passwords/{pid} | 更新密码 |
| DELETE | /api/v1/passwords/{pid} | 删除密码 |
| POST | /api/v1/passwords/{pid}/activate | 激活密码 |
| POST | /api/v1/passwords/{pid}/deactivate | 停用密码 |
| GET | /api/v1/packages/{id}/passwords/current | 当前有效密码 |

### 3.4 验证服务

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/verify | 验证密码 |
| POST | /api/v1/verify/batch | 批量验证 |
| GET | /api/v1/verify/status/{id} | 检查包状态 |
| GET | /api/v1/verify/offline-config/{id} | 获取离线配置 |

### 3.5 审计日志

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/audit | 查询日志 |
| GET | /api/v1/audit/package/{id} | 包日志 |
| GET | /api/v1/audit/user/{id} | 用户日志 |
| GET | /api/v1/audit/verify-fails | 验证失败记录 |
| GET | /api/v1/audit/export | 导出日志 |

### 3.6 系统管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/admin/encryption/config | 获取加密配置 |
| PUT | /api/v1/admin/encryption/config | 更新加密配置 |
| POST | /api/v1/admin/encryption/validate | 验证加密配置 |
| POST | /api/v1/admin/encryption/regenerate-keys | 重新生成密钥 |
| GET | /api/v1/admin/system | 系统信息 |
| GET | /api/v1/admin/roles | 角色列表 |
| GET | /api/v1/admin/stats/dashboard | 仪表盘数据 |

---

## 4. 数据库设计

### 4.1 ER 图

```
┌─────────────┐       ┌─────────────────┐       ┌─────────────────┐
│    User     │       │  FilePackage    │       │ PasswordPolicy  │
├─────────────┤       ├─────────────────┤       ├─────────────────┤
│ id (PK)     │       │ id (PK)         │       │ id (PK)         │
│ username    │◄──────│ created_by (FK) │       │ package_id (FK) │
│ email       │       │ name            │◄──────│ password_hash   │
│ password    │       │ format          │       │ priority        │
│ role        │       │ status          │       │ valid_from      │
│ status      │       │ file_path       │       │ valid_until     │
│ created_at  │       │ file_hash       │       │ status          │
│ last_login  │       │ created_at      │       │ created_at      │
└─────────────┘       └─────────────────┘       └─────────────────┘

┌─────────────────┐       ┌─────────────────┐
│    AuditLog     │       │ EncryptionConfig │
├─────────────────┤       ├─────────────────┤
│ id (PK)         │       │ id (PK)          │
│ action           │       │ aes_key_length   │
│ package_id (FK)  │       │ rsa_key_length   │
│ user_id (FK)     │       │ password_policy  │
│ ip_address       │       │ config_encrypt   │
│ user_agent       │       │ enable_signature │
│ detail (json)    │       │ updated_at       │
│ created_at       │       │ updated_by (FK)  │
└─────────────────┘       └─────────────────┘

┌─────────────────┐
│    ApiKey       │
├─────────────────┤
│ id (PK)         │
│ user_id (FK)    │
│ name            │
│ key_hash        │
│ expires_at      │
│ created_at      │
└─────────────────┘
```

### 4.2 表结构

#### users

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(64) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(32) NOT NULL DEFAULT 'operator',
    status VARCHAR(32) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_login TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
```

#### file_packages

```sql
CREATE TABLE file_packages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    format VARCHAR(16) NOT NULL,
    description TEXT,
    status VARCHAR(32) NOT NULL DEFAULT 'active',
    file_path VARCHAR(512) NOT NULL,
    file_hash VARCHAR(128) NOT NULL,
    file_size BIGINT NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE INDEX idx_file_packages_created_by ON file_packages(created_by);
CREATE INDEX idx_file_packages_status ON file_packages(status);
```

#### password_policies

```sql
CREATE TABLE password_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    package_id UUID NOT NULL REFERENCES file_packages(id) ON DELETE CASCADE,
    password_hash VARCHAR(255) NOT NULL,
    priority INT NOT NULL DEFAULT 1,
    valid_from TIMESTAMP,
    valid_until TIMESTAMP,
    status VARCHAR(32) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_password_policies_package_id ON password_policies(package_id);
CREATE INDEX idx_password_policies_status ON password_policies(status);
```

#### audit_logs

```sql
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    action VARCHAR(64) NOT NULL,
    package_id UUID REFERENCES file_packages(id),
    user_id UUID REFERENCES users(id),
    ip_address VARCHAR(64),
    user_agent TEXT,
    detail JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_package_id ON audit_logs(package_id);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
```

#### encryption_configs

```sql
CREATE TABLE encryption_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aes_key_length INT NOT NULL DEFAULT 256,
    rsa_key_length INT NOT NULL DEFAULT 2048,
    password_min_length INT NOT NULL DEFAULT 8,
    password_require_special BOOLEAN NOT NULL DEFAULT TRUE,
    password_require_uppercase BOOLEAN NOT NULL DEFAULT TRUE,
    password_require_lowercase BOOLEAN NOT NULL DEFAULT TRUE,
    password_require_digit BOOLEAN NOT NULL DEFAULT TRUE,
    config_encrypt BOOLEAN NOT NULL DEFAULT TRUE,
    enable_signature BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES users(id)
);
```

#### api_keys

```sql
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(128) NOT NULL,
    key_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX idx_api_keys_key_hash ON api_keys(key_hash);
```

---

## 5. 安全设计

### 5.1 加密体系

```
┌─────────────────────────────────────────────────────┐
│                    密钥层次结构                       │
├─────────────────────────────────────────────────────┤
│  主密钥 (Master Key)                                 │
│  └── 存储在环境变量或密钥管理系统                     │
│      ├── 用于加密数据库敏感字段                       │
│      └── 用于加密配置文件                            │
│                                                      │
│  应用密钥                                            │
│  ├── AES-256-GCM: 加密文件内容                        │
│  ├── RSA-2048: 配置签名验证                          │
│  └── bcrypt: 密码哈希                                │
└─────────────────────────────────────────────────────┘
```

### 5.2 客户端防篡改

```
┌─────────────────────────────────────────────────────┐
│                   EXE/ZIP 内部结构                    │
├─────────────────────────────────────────────────────┤
│  verify.exe                                          │
│  ├── 内嵌 RSA 公钥                                   │
│  ├── 内嵌 AES 密钥                                   │
│  ├── 内嵌配置哈希                                    │
│  └── 解压程序                                        │
│                                                      │
│  config.bin (隐藏文件)                               │
│  └── AES 加密的配置 + RSA 签名                       │
│                                                      │
│  files/                                             │
│  └── encrypted_files.zip                            │
└─────────────────────────────────────────────────────┘
```

### 5.3 验证流程

```
1. 检查 config.bin 是否存在
2. 验证 config.bin 签名（RSA）
3. 解密 config.bin 获取策略（AES）
4. 验证配置哈希完整性
5. 尝试验证密码（在线优先）
6. 在线失败则使用本地策略
```

---

## 6. 缓存设计

### 6.1 Redis 缓存策略

| 缓存 Key | 内容 | TTL |
|----------|------|-----|
| `package:{id}` | 文件包信息 | 5 分钟 |
| `passwords:{package_id}` | 密码策略列表 | 1 分钟 |
| `current_password:{package_id}` | 当前有效密码 | 30 秒 |
| `encryption:config` | 加密配置 | 1 小时 |
| `ratelimit:{ip}` | 限流计数 | 1 分钟 |

---

## 7. 文件存储

### 7.1 目录结构

```
storage/
├── packages/              # 用户上传的原始文件
│   └── {year}/{month}/{day}/{uuid}/
├── generated/            # 生成的文件包
│   └── {year}/{month}/{day}/{uuid}.exe
│   └── {year}/{month}/{day}/{uuid}.zip
└── temp/                 # 临时文件
```

### 7.2 存储策略

- 本地存储（默认）：`./storage` 目录
- S3 兼容存储（可选）：配置 S3 相关环境变量

---

## 8. 错误码设计

| 范围 | 说明 |
|------|------|
| 0 | 成功 |
| 10001-10099 | 认证/授权相关错误 |
| 20001-20099 | 文件相关错误 |
| 30001-30099 | 密码验证相关错误 |
| 40001-40099 | 加密相关错误 |
| 90001-99999 | 系统内部错误 |
