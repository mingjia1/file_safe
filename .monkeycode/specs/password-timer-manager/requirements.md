# 需求文档 - Password Timer Manager

## 1. 项目概述

**项目名称**：Password Timer Manager（密码时效管理器）

**项目类型**：企业级文件安全分享平台

**核心功能摘要**：一个支持多密码时效控制的文件安全分享平台，可生成带密码策略的自解压程序或加密 ZIP 包，通过服务器或本地验证控制文件访问。

**目标用户**：
- 企业内部需要安全分享敏感文件的员工
- 需要向外部合作伙伴提供有时效性的文件访问的企业
- 需要 SDK/API 集成到内部系统的企业

---

## 2. 功能需求

### 2.1 文件包管理

| 需求ID | 需求描述 | 优先级 |
|--------|----------|--------|
| F001 | 支持上传文件/文件夹创建文件包 | P0 |
| F002 | 支持通过远程 URL 创建文件包 | P1 |
| F003 | 支持生成 EXE 自解压程序格式 | P0 |
| F004 | 支持生成 ZIP + 内嵌验证程序格式 | P0 |
| F005 | 支持查看文件包列表（分页、筛选） | P0 |
| F006 | 支持查看文件包详情（含密码列表） | P0 |
| F007 | 支持更新文件包名称、描述、状态 | P1 |
| F008 | 支持删除文件包（级联删除密码策略） | P0 |
| F009 | 支持下载生成的文件包 | P0 |
| F010 | 支持生成临时下载链接 | P2 |

### 2.2 密码策略管理

| 需求ID | 需求描述 | 优先级 |
|--------|----------|--------|
| P001 | 支持为文件包添加多个密码 | P0 |
| P002 | 支持设置密码优先级（1=最高） | P0 |
| P003 | 支持设置密码有效期（valid_from, valid_until） | P0 |
| P004 | 支持批量添加密码策略 | P1 |
| P005 | 支持手动激活密码（跳过有效期检查） | P0 |
| P006 | 支持手动停用密码 | P0 |
| P007 | 支持删除密码策略 | P1 |
| P008 | 支持查看当前有效密码 | P1 |
| P009 | 密码自动状态转换（pending→active→expired） | P0 |

### 2.3 验证服务

| 需求ID | 需求描述 | 优先级 |
|--------|----------|--------|
| V001 | 支持在线验证密码（优先） | P0 |
| V002 | 支持离线本地验证（网络异常时兜底） | P0 |
| V003 | 支持批量尝试验证多个密码 | P2 |
| V004 | 支持检查文件包状态 | P1 |
| V005 | 验证频率限制（防暴力破解） | P1 |

### 2.4 客户端验证程序

| 需求ID | 需求描述 | 优先级 |
|--------|----------|--------|
| C001 | 客户端防篡改（哈希校验） | P0 |
| C002 | 配置文件加密存储（AES） | P0 |
| C003 | 配置签名验证（RSA） | P0 |
| C004 | 隐藏文件属性保护 | P1 |
| C005 | 验证失败友好提示 | P1 |

### 2.5 用户与权限

| 需求ID | 需求描述 | 优先级 |
|--------|----------|--------|
| U001 | 支持用户账号密码登录 | P0 |
| U002 | 支持 JWT Token 认证 | P0 |
| U003 | 支持 API Key 认证（SDK 接入） | P0 |
| U004 | 支持预定义角色（super_admin, admin, operator, viewer） | P0 |
| U005 | 支持权限验证中间件 | P0 |
| U006 | 支持 API Key 管理（生成、列表、删除） | P1 |

### 2.6 审计日志

| 需求ID | 需求描述 | 优先级 |
|--------|----------|--------|
| A001 | 记录文件包下载事件 | P1 |
| A002 | 记录密码验证成功/失败 | P0 |
| A003 | 记录管理员操作（密码策略变更等） | P1 |
| A004 | 支持查询审计日志（分页、筛选） | P1 |
| A005 | 支持导出审计日志（CSV/JSON） | P2 |

### 2.7 系统管理

| 需求ID | 需求描述 | 优先级 |
|--------|----------|--------|
| S001 | 全局加密配置（AES/RSA 密钥长度） | P0 |
| S002 | 全局密码策略配置（最小长度、复杂度要求） | P0 |
| S003 | 加密配置验证接口（测试配置有效性） | P0 |
| S004 | 管理员可重新生成加密密钥 | P1 |
| S005 | 系统健康检查接口 | P0 |
| S006 | 仪表盘统计数据 | P2 |

---

## 3. 非功能需求

### 3.1 性能需求

| 项目 | 指标 |
|------|------|
| API 响应时间 | P95 < 200ms |
| 并发用户数 | 支持 100+ 并发 |
| 文件上传大小 | 最大 500MB |

### 3.2 安全需求

| 项目 | 要求 |
|------|------|
| 密码存储 | bcrypt 哈希 |
| JWT 过期 | 可配置，默认 24h |
| API 限流 | 10 次/分钟/IP |
| 敏感数据 | 加密存储 |

### 3.3 兼容性需求

| 项目 | 要求 |
|------|------|
| 后端 | Go 1.21+ |
| 数据库 | PostgreSQL 15+ |
| 缓存 | Redis 7+ |
| 客户端 | Windows 7+ |

---

## 4. 数据模型

### 4.1 User

| 字段 | 类型 | 约束 |
|------|------|------|
| id | UUID | PK |
| username | string | unique, required |
| email | string | unique, required |
| password_hash | string | required |
| role | enum | required |
| status | enum | default ACTIVE |
| created_at | datetime | required |
| last_login | datetime | nullable |

### 4.2 FilePackage

| 字段 | 类型 | 约束 |
|------|------|------|
| id | UUID | PK |
| name | string | required |
| format | enum | EXE/ZIP |
| description | string | nullable |
| status | enum | default ACTIVE |
| file_path | string | required |
| file_hash | string | required |
| file_size | int64 | required |
| created_by | UUID | FK → User |
| created_at | datetime | required |
| updated_at | datetime | nullable |

### 4.3 PasswordPolicy

| 字段 | 类型 | 约束 |
|------|------|------|
| id | UUID | PK |
| package_id | UUID | FK → FilePackage |
| password_hash | string | required |
| priority | int | default 1 |
| valid_from | datetime | nullable |
| valid_until | datetime | nullable |
| status | enum | default PENDING |
| created_at | datetime | required |

### 4.4 AuditLog

| 字段 | 类型 | 约束 |
|------|------|------|
| id | UUID | PK |
| action | enum | required |
| package_id | UUID | FK → FilePackage, nullable |
| user_id | UUID | FK → User, nullable |
| ip_address | string | nullable |
| user_agent | string | nullable |
| detail | json | nullable |
| created_at | datetime | required |

### 4.5 EncryptionConfig

| 字段 | 类型 | 约束 |
|------|------|------|
| id | UUID | PK |
| aes_key_length | int | default 256 |
| rsa_key_length | int | default 2048 |
| password_min_length | int | default 8 |
| password_require_special | bool | default true |
| config_encrypt | bool | default true |
| enable_signature | bool | default true |
| updated_at | datetime | required |
| updated_by | UUID | FK → User |

---

## 5. API 接口

详见 `API.md` 文档。

---

## 6. 验收标准

### 6.1 文件包管理

- [ ] 可以上传文件创建 EXE 格式文件包
- [ ] 可以上传文件创建 ZIP 格式文件包
- [ ] 可以查看文件包列表（分页）
- [ ] 可以查看文件包详情（含密码）
- [ ] 可以删除文件包

### 6.2 密码策略

- [ ] 可以添加单个密码
- [ ] 可以批量添加多个密码
- [ ] 可以手动激活/停用密码
- [ ] 密码状态自动转换正确

### 6.3 验证

- [ ] 在线验证返回正确结果
- [ ] 离线验证在网络异常时正常工作
- [ ] 防篡改机制有效

### 6.4 客户端

- [ ] EXE 格式可正常解压
- [ ] ZIP 格式可正常解压
- [ ] 隐藏文件属性生效

### 6.5 认证授权

- [ ] 用户可以登录
- [ ] API Key 可以正常认证
- [ ] 权限验证正确

---

## 7. 术语表

| 术语 | 定义 |
|------|------|
| 文件包 | 用户上传并设置密码策略的文件 |
| 密码策略 | 定义密码内容、优先级和有效期 |
| 自解压程序 | 带验证程序的可执行压缩包 |
| 在线验证 | 通过 API 服务器验证密码 |
| 离线验证 | 客户端本地验证密码 |
