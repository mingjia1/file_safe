# Python SDK 接入手册

## 安装

```bash
pip install ptm-sdk
```

## 快速开始

```python
from ptm_sdk import PTMClient

# 初始化客户端
client = PTMClient(api_key="your_api_key_here")

# 验证密码
result = client.verify(
    package_id="pkg-uuid-xxx",
    password="user_input_password"
)

print(f"验证结果: {result.valid}")
if result.valid:
    print(f"解密密钥: {result.key}")
```

## 客户端初始化

### 使用 API Key 初始化

```python
from ptm_sdk import PTMClient

client = PTMClient(api_key="ptm_sk_a1b2c3d4e5f6...")
```

### 使用 JWT Token 初始化

```python
from ptm_sdk import PTMClient

client = PTMClient(token="jwt_token_here")
```

### 自定义配置

```python
from ptm_sdk import PTMClient

client = PTMClient(
    api_key="your_api_key",
    base_url="http://localhost:8080/api/v1",
    timeout=30,
    max_retries=3
)
```

## 文件包管理

### 创建文件包

```python
# 上传文件创建
with open("document.pdf", "rb") as f:
    package = client.packages.create(
        file=f,
        name="敏感文档",
        format="exe"
    )
print(f"创建成功: {package.id}")

# 通过 URL 创建
package = client.packages.create_from_url(
    name="远程文件",
    source_url="https://example.com/file.pdf",
    format="zip"
)
```

### 列出文件包

```python
# 分页列出
result = client.packages.list(page=1, page_size=20)

for package in result.items:
    print(f"{package.name} - {package.status}")

# 筛选
active_packages = client.packages.list(status="active", format="exe")
```

### 获取文件包详情

```python
package = client.packages.get("pkg-uuid-xxx")
print(f"名称: {package.name}")
print(f"格式: {package.format}")
print(f"密码数量: {len(package.passwords)}")
```

### 更新文件包

```python
client.packages.update(
    "pkg-uuid-xxx",
    name="更新后的名称",
    description="新的描述"
)
```

### 删除文件包

```python
client.packages.delete("pkg-uuid-xxx")
```

### 下载文件包

```python
# 直接下载到文件
client.packages.download(
    "pkg-uuid-xxx",
    save_path="output.exe"
)

# 获取下载链接
url_data = client.packages.get_download_url("pkg-uuid-xxx", expires_in=3600)
print(f"下载链接: {url_data.download_url}")
```

## 密码策略管理

### 添加密码

```python
from datetime import datetime, timedelta

# 创建密码
password = client.passwords.create(
    package_id="pkg-uuid-xxx",
    password="Week1Password!",
    priority=1,
    valid_from=datetime(2024, 1, 1),
    valid_until=datetime(2024, 1, 8)
)
```

### 批量添加密码

```python
from datetime import datetime, timedelta

start = datetime(2024, 1, 1)
passwords = []
for week in range(4):
    passwords.append({
        "password": f"Week{week+1}Password!",
        "priority": week + 1,
        "valid_from": start + timedelta(weeks=week),
        "valid_until": start + timedelta(weeks=week+1)
    })

result = client.passwords.create_batch("pkg-uuid-xxx", passwords)
print(f"成功创建 {result.created_count} 个密码")
```

### 激活/停用密码

```python
# 激活
client.passwords.activate("pwd-uuid-xxx")

# 停用
client.passwords.deactivate("pwd-uuid-xxx")
```

### 获取当前有效密码

```python
current = client.passwords.get_current("pkg-uuid-xxx")
print(f"当前密码: {current.password}")
```

## 验证服务

### 验证密码

```python
result = client.verify(
    package_id="pkg-uuid-xxx",
    password="user_input"
)

if result.valid:
    # 使用密钥解密
    print(f"解密密钥: {result.key}")
else:
    print(f"验证失败: {result.message}")
```

### 批量验证

```python
result = client.verify_batch(
    package_id="pkg-uuid-xxx",
    passwords=["pwd1", "pwd2", "pwd3"]
)

if result.valid:
    print(f"匹配到密码: {result.matched_password}")
```

### 检查包状态

```python
status = client.verify.get_status("pkg-uuid-xxx")
print(f"状态: {status.status}")
print(f"有效密码数: {status.current_password_count}")
```

## 审计日志

### 查询日志

```python
# 查询所有日志
logs = client.audit.list(page=1, page_size=20)

# 按条件筛选
logs = client.audit.list(
    action="VERIFY_SUCCESS",
    package_id="pkg-uuid-xxx",
    start_time=datetime(2024, 1, 1)
)

for log in logs.items:
    print(f"{log.action} - {log.created_at}")
```

### 导出日志

```python
# 导出为 JSON
logs_json = client.audit.export(
    format="json",
    start_time=datetime(2024, 1, 1),
    end_time=datetime(2024, 1, 31)
)

# 导出为 CSV
client.audit.export(
    format="csv",
    save_path="audit_log.csv",
    start_time=datetime(2024, 1, 1),
    end_time=datetime(2024, 1, 31)
)
```

## 错误处理

```python
from ptm_sdk import PTMClient
from ptm_sdk.exceptions import (
    PTMError,
    AuthenticationError,
    PermissionError,
    NotFoundError,
    ValidationError
)

try:
    client = PTMClient(api_key="invalid_key")
    result = client.verify("pkg-xxx", "password")
except AuthenticationError:
    print("认证失败，请检查 API Key")
except PermissionError:
    print("权限不足")
except NotFoundError as e:
    print(f"资源不存在: {e.resource_id}")
except ValidationError as e:
    print(f"参数错误: {e.message}")
except PTMError as e:
    print(f"其他错误: {e.message}")
```

## 完整示例

```python
from ptm_sdk import PTMClient
from datetime import datetime, timedelta
import time

def main():
    # 初始化
    client = PTMClient(api_key="your_api_key")
    
    # 1. 创建文件包
    print("1. 创建文件包...")
    with open("sensitive_file.pdf", "rb") as f:
        package = client.packages.create(
            file=f,
            name="Q4财务报告",
            format="exe"
        )
    print(f"   文件包 ID: {package.id}")
    
    # 2. 批量添加密码策略（4 周）
    print("2. 添加密码策略...")
    start = datetime(2024, 1, 1)
    passwords = []
    for week in range(4):
        passwords.append({
            "password": f"Week{week+1}#Pass!",
            "priority": week + 1,
            "valid_from": start + timedelta(weeks=week),
            "valid_until": start + timedelta(weeks=week+1)
        })
    
    result = client.passwords.create_batch(package.id, passwords)
    print(f"   成功创建 {result.created_count} 个密码策略")
    
    # 3. 下载生成的文件包
    print("3. 下载文件包...")
    client.packages.download(
        package.id,
        save_path=f"Q4_Report_Week{int(time.time())}.exe"
    )
    print("   下载完成")
    
    # 4. 查询审计日志
    print("4. 审计日志...")
    logs = client.audit.list(package_id=package.id, page_size=10)
    print(f"   记录数: {logs.total}")
    
    print("\n全部完成!")

if __name__ == "__main__":
    main()
```

## SDK 配置项

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| base_url | `http://localhost:8080/api/v1` | API 基础 URL |
| timeout | 30 | 请求超时（秒） |
| max_retries | 3 | 最大重试次数 |
| retry_delay | 1 | 重试间隔（秒） |

## 异步客户端

```python
import asyncio
from ptm_sdk.async_client import AsyncPTMClient

async def main():
    async with AsyncPTMClient(api_key="your_api_key") as client:
        # 并发请求
        results = await asyncio.gather(
            client.packages.list(),
            client.passwords.list("pkg-xxx"),
            client.audit.list(page_size=5)
        )
        
        packages, passwords, logs = results
        print(f"包数量: {packages.total}")
        print(f"密码数量: {len(passwords)}")
        print(f"日志数量: {logs.total}")

asyncio.run(main())
```

## 类型定义

```python
from ptm_sdk.types import (
    Package,
    Password,
    AuditLog,
    VerifyResult,
    EncryptionConfig
)

# 使用类型提示
def process_package(pkg: Package) -> str:
    return f"{pkg.name} ({pkg.format})"
```
