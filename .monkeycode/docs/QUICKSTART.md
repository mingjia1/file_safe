# 快速开始指南

## 5 分钟快速上手

本指南将帮助您在 5 分钟内完成 Password Timer Manager 的基本配置和使用。

---

## 前置条件

- Docker 和 Docker Compose 已安装
- 或者手动安装 PostgreSQL 15+ 和 Redis 7+

---

## 第一步：启动服务

### 使用 Docker Compose 启动（推荐）

```bash
# 克隆项目
git clone https://your-repo/password-timer-manager.git
cd password-timer-manager

# 启动所有服务
docker-compose up -d

# 检查服务状态
docker-compose ps
```

### 手动启动

```bash
# 启动 PostgreSQL
pg_ctl -D /var/lib/postgresql/data start

# 启动 Redis
redis-server --daemonize yes

# 启动 API 服务
cd backend && go run ./cmd/server

# 启动前端（另一个终端）
cd frontend && npm run dev
```

---

## 第二步：访问管理平台

打开浏览器访问：`http://localhost:3000`

首次登录使用默认管理员账号：
- 用户名：`admin`
- 密码：`admin123`（请立即修改）

---

## 第三步：配置全局加密设置

1. 进入 **系统设置** → **加密配置**
2. 根据您的安全要求配置：
   - AES 密钥长度：256（推荐）
   - RSA 密钥长度：2048（推荐）
   - 密码最小长度：8
3. 点击 **验证配置** 确保配置有效
4. 保存配置

---

## 第四步：创建文件包

### 通过 Web 界面

1. 进入 **文件包管理** → **创建文件包**
2. 上传要保护的文件
3. 选择输出格式：
   - **EXE**：自解压程序（单一文件）
   - **ZIP**：ZIP + 内嵌验证程序
4. 填写名称和描述
5. 点击 **创建**

### 通过 API

```bash
curl -X POST http://localhost:8080/api/v1/packages \
  -H "Authorization: Bearer <your_token>" \
  -F "file=@document.pdf" \
  -F "name=敏感文档" \
  -F "format=exe"
```

### 通过 Python SDK

```python
from ptm_sdk import PTMClient

client = PTMClient(api_key="your_api_key")

with open("document.pdf", "rb") as f:
    pkg = client.packages.create(
        file=f,
        name="敏感文档",
        format="exe"
    )

print(f"文件包创建成功: {pkg.id}")
```

---

## 第五步：设置密码策略

### 方式一：固定有效期

设置多个密码，每个有固定的有效期：

```python
from datetime import datetime, timedelta

start = datetime(2024, 1, 1)
passwords = [
    {
        "password": "Week1Pass!",
        "priority": 1,
        "valid_from": start,
        "valid_until": start + timedelta(weeks=1)
    },
    {
        "password": "Week2Pass!",
        "priority": 2,
        "valid_from": start + timedelta(weeks=1),
        "valid_until": start + timedelta(weeks=2)
    },
    {
        "password": "Week3Pass!",
        "priority": 3,
        "valid_from": start + timedelta(weeks=2),
        "valid_until": start + timedelta(weeks=3)
    }
]

result = client.passwords.create_batch(pkg.id, passwords)
```

### 方式二：立即切换

管理员随时手动切换当前密码：

```python
# 停用当前密码
client.passwords.deactivate("pwd-uuid-xxx")

# 激活下一个密码
client.passwords.activate("pwd-uuid-yyy")
```

### 通过 API

```bash
# 添加密码
curl -X POST http://localhost:8080/api/v1/packages/{package_id}/passwords \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "password": "Week1Pass!",
    "priority": 1,
    "valid_from": "2024-01-01T00:00:00Z",
    "valid_until": "2024-01-08T00:00:00Z"
  }'
```

---

## 第六步：下载分发

### 生成下载包

1. 在文件包详情页点击 **下载**
2. 系统自动生成带密码策略的 EXE 或 ZIP 文件
3. 将生成的文件分发给接收者

### 通过 API 下载

```bash
# 获取下载链接
curl -X POST http://localhost:8080/api/v1/packages/{package_id}/download-url \
  -H "Authorization: Bearer <token>" \
  -d '{"expires_in": 3600}'

# 下载文件
curl -O http://localhost:8080/api/v1/packages/{package_id}/download \
  -H "Authorization: Bearer <token>"
```

---

## 第七步：客户端使用

### EXE 格式

1. 接收者双击 `.exe` 文件
2. 输入当前有效密码
3. 验证通过后自动解压文件

### ZIP 格式

1. 接收者解压 ZIP 文件
2. 双击运行验证程序（显示为系统文件）
3. 输入当前有效密码
4. 验证通过后释放文件

---

## 常见问题

### Q: 忘记了密码怎么办？

A: 管理员可以在管理平台查看当前有效的密码。

### Q: 所有密码都过期了怎么办？

A: 管理员可以为该文件包添加新密码，然后接收者重新下载。

### Q: 如何实现实时密码切换？

A: 使用 API 调用 `POST /passwords/{id}/activate` 可以立即切换密码。

### Q: 离线环境下如何使用？

A: 客户端内置了离线验证机制，在网络不可用时会自动使用本地策略验证。

---

## 下一步

- [完整 API 文档](./API.md)
- [SDK 接入指南](./SDK_PYTHON.md)
- [部署指南](./DEPLOY.md)
