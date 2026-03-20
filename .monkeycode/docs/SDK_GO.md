# Go SDK 接入手册

## 环境要求

- Go 1.18 或更高版本

## 安装

```bash
go get github.com/ptm/go-sdk
```

## 快速开始

```go
package main

import (
    "fmt"
    "github.com/ptm/go-sdk"
)

func main() {
    // 初始化客户端
    client := ptm.NewClient("your_api_key_here")
    
    // 验证密码
    result, err := client.Verify("pkg-uuid-xxx", "user_input_password")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("验证结果: %v\n", result.Valid)
    if result.Valid {
        fmt.Printf("解密密钥: %s\n", result.Key)
    }
}
```

## 客户端初始化

### 使用 API Key 初始化

```go
import "github.com/ptm/go-sdk"

client := ptm.NewClient("ptm_sk_a1b2c3d4e5f6...")
```

### 自定义配置

```go
client := ptm.NewClient(
    ptm.WithAPIKey("your_api_key"),
    ptm.WithBaseURL("http://localhost:8080/api/v1"),
    ptm.WithTimeout(30),
    ptm.WithMaxRetries(3),
)
```

## 文件包管理

### 创建文件包

```go
import (
    "os"
    "github.com/ptm/go-sdk"
)

// 上传文件创建
file, err := os.Open("document.pdf")
if err != nil {
    panic(err)
}
defer file.Close()

pkg, err := client.Packages.Create(ptm.CreatePackageRequest{
    Name:   "敏感文档",
    Format: "exe",
    File:   file,
})
if err != nil {
    panic(err)
}
fmt.Printf("创建成功: %s\n", pkg.ID)

// 通过 URL 创建
pkgFromURL, err := client.Packages.CreateFromURL(ptm.CreatePackageFromURLRequest{
    Name:      "远程文件",
    SourceURL: "https://example.com/file.pdf",
    Format:    "zip",
})
```

### 列出文件包

```go
import "github.com/ptm/go-sdk"

// 分页列出
list, err := client.Packages.List(&ptm.ListPackagesRequest{
    Page:     1,
    PageSize: 20,
})

for _, pkg := range list.Items {
    fmt.Printf("%s - %s\n", pkg.Name, pkg.Status)
}

// 筛选
activeExe, err := client.Packages.List(&ptm.ListPackagesRequest{
    Status: "active",
    Format: "exe",
})
```

### 获取文件包详情

```go
pkg, err := client.Packages.Get("pkg-uuid-xxx")
if err != nil {
    panic(err)
}
fmt.Printf("名称: %s\n", pkg.Name)
fmt.Printf("格式: %s\n", pkg.Format)
fmt.Printf("密码数量: %d\n", len(pkg.Passwords))
```

### 更新文件包

```go
err = client.Packages.Update("pkg-uuid-xxx", &ptm.UpdatePackageRequest{
    Name:        "更新后的名称",
    Description: "新的描述",
})
```

### 删除文件包

```go
err = client.Packages.Delete("pkg-uuid-xxx")
```

### 下载文件包

```go
import "github.com/ptm/go-sdk"

// 下载到文件
err = client.Packages.Download("pkg-uuid-xxx", "output.exe")
if err != nil {
    panic(err)
}

// 获取下载链接
urlData, err := client.Packages.GetDownloadURL("pkg-uuid-xxx", 3600)
fmt.Printf("下载链接: %s\n", urlData.DownloadURL)
```

## 密码策略管理

### 添加密码

```go
import (
    "time"
    "github.com/ptm/go-sdk"
)

start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
end := start.AddDate(0, 0, 7)

pwd, err := client.Passwords.Create("pkg-uuid-xxx", &ptm.CreatePasswordRequest{
    Password:   "Week1Password!",
    Priority:   1,
    ValidFrom:  &start,
    ValidUntil: &end,
})
```

### 批量添加密码

```go
import "github.com/ptm/go-sdk"

start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

passwords := make([]*ptm.PasswordItem, 4)
for week := 0; week < 4; week++ {
    validFrom := start.AddDate(0, 0, week*7)
    validUntil := start.AddDate(0, 0, (week+1)*7)
    passwords[week] = &ptm.PasswordItem{
        Password:   fmt.Sprintf("Week%dPassword!", week+1),
        Priority:   week + 1,
        ValidFrom:  &validFrom,
        ValidUntil: &validUntil,
    }
}

result, err := client.Passwords.CreateBatch("pkg-uuid-xxx", passwords)
fmt.Printf("成功创建 %d 个密码\n", result.CreatedCount)
```

### 激活/停用密码

```go
// 激活
err = client.Passwords.Activate("pwd-uuid-xxx")

// 停用
err = client.Passwords.Deactivate("pwd-uuid-xxx")
```

### 获取当前有效密码

```go
current, err := client.Passwords.GetCurrent("pkg-uuid-xxx")
fmt.Printf("当前密码: %s\n", current.Password)
```

## 验证服务

### 验证密码

```go
result, err := client.Verify(&ptm.VerifyRequest{
    PackageID: "pkg-uuid-xxx",
    Password:  "user_input",
})
if err != nil {
    panic(err)
}

if result.Valid {
    fmt.Printf("解密密钥: %s\n", result.Key)
} else {
    fmt.Printf("验证失败: %s\n", result.Message)
}
```

### 批量验证

```go
result, err := client.VerifyBatch(&ptm.VerifyBatchRequest{
    PackageID: "pkg-uuid-xxx",
    Passwords: []string{"pwd1", "pwd2", "pwd3"},
})

if result.Valid {
    fmt.Printf("匹配到密码: %s\n", result.MatchedPassword)
}
```

### 检查包状态

```go
status, err := client.Verify.GetStatus("pkg-uuid-xxx")
fmt.Printf("状态: %s\n", status.Status)
fmt.Printf("有效密码数: %d\n", status.CurrentPasswordCount)
```

## 审计日志

### 查询日志

```go
import "github.com/ptm/go-sdk"

// 查询所有日志
logs, err := client.Audit.List(&ptm.ListAuditRequest{
    Page:     1,
    PageSize: 20,
})

// 按条件筛选
logsFiltered, err := client.Audit.List(&ptm.ListAuditRequest{
    Action:    "VERIFY_SUCCESS",
    PackageID: "pkg-uuid-xxx",
    StartTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
})

for _, log := range logsFiltered.Items {
    fmt.Printf("%s - %s\n", log.Action, log.CreatedAt)
}
```

### 导出日志

```go
import "github.com/ptm/go-sdk"

// 导出为 JSON
jsonLogs, err := client.Audit.Export(&ptm.ExportAuditRequest{
    Format:     "json",
    StartTime:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
    EndTime:    time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
})

// 导出为 CSV 到文件
err = client.Audit.ExportToFile(&ptm.ExportAuditRequest{
    Format:     "csv",
    SavePath:   "audit_log.csv",
    StartTime:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
    EndTime:    time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
})
```

## 错误处理

```go
import (
    "errors"
    "github.com/ptm/go-sdk"
    ptmerrors "github.com/ptm/go-sdk/errors"
)

result, err := client.Verify(&ptm.VerifyRequest{
    PackageID: "pkg-xxx",
    Password:  "password",
})

if err != nil {
    var authErr *ptmerrors.AuthenticationError
    var permErr *ptmerrors.PermissionError
    var notFoundErr *ptmerrors.NotFoundError
    var validationErr *ptmerrors.ValidationError
    
    switch {
    case errors.As(err, &authErr):
        fmt.Println("认证失败，请检查 API Key")
    case errors.As(err, &permErr):
        fmt.Println("权限不足")
    case errors.As(err, &notFoundErr):
        fmt.Printf("资源不存在: %s\n", notFoundErr.ResourceID)
    case errors.As(err, &validationErr):
        fmt.Printf("参数错误: %s\n", validationErr.Message)
    default:
        fmt.Printf("其他错误: %s\n", err)
    }
}
```

## 完整示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/ptm/go-sdk"
)

func main() {
    // 1. 初始化
    fmt.Println("1. 初始化客户端...")
    client := ptm.NewClient("your_api_key")
    
    // 2. 创建文件包
    fmt.Println("2. 创建文件包...")
    file, err := os.Open("sensitive_file.pdf")
    if err != nil {
        panic(err)
    }
    defer file.Close()
    
    pkg, err := client.Packages.Create(ptm.CreatePackageRequest{
        Name:   "Q4财务报告",
        Format: "exe",
        File:   file,
    })
    if err != nil {
        panic(err)
    }
    fmt.Printf("   文件包 ID: %s\n", pkg.ID)
    
    // 3. 批量添加密码策略
    fmt.Println("3. 添加密码策略...")
    start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    passwords := make([]*ptm.PasswordItem, 4)
    for week := 0; week < 4; week++ {
        validFrom := start.AddDate(0, 0, week*7)
        validUntil := start.AddDate(0, 0, (week+1)*7)
        passwords[week] = &ptm.PasswordItem{
            Password:   fmt.Sprintf("Week%d#Pass!", week+1),
            Priority:   week + 1,
            ValidFrom:  &validFrom,
            ValidUntil: &validUntil,
        }
    }
    
    batchResult, err := client.Passwords.CreateBatch(pkg.ID, passwords)
    if err != nil {
        panic(err)
    }
    fmt.Printf("   成功创建 %d 个密码策略\n", batchResult.CreatedCount)
    
    // 4. 下载文件包
    fmt.Println("4. 下载文件包...")
    err = client.Packages.Download(pkg.ID, "Q4_Report.exe")
    if err != nil {
        panic(err)
    }
    fmt.Println("   下载完成")
    
    // 5. 查询审计日志
    fmt.Println("5. 审计日志...")
    logs, err := client.Audit.List(&ptm.ListAuditRequest{
        PackageID: pkg.ID,
        PageSize:  10,
    })
    if err != nil {
        panic(err)
    }
    fmt.Printf("   记录数: %d\n", logs.Total)
    
    fmt.Println("\n全部完成!")
}
```

## Context 支持

所有 API 调用都支持 context：

```go
import (
    "context"
    "time"
)

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

result, err := client.VerifyWithContext(ctx, &ptm.VerifyRequest{
    PackageID: "pkg-uuid-xxx",
    Password:  "password",
})
```

## 并发安全

Go SDK 是并发安全的，可以安全地在多个 goroutine 中使用同一个客户端实例。

```go
import "sync"

var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        result, _ := client.Verify(&ptm.VerifyRequest{
            PackageID: fmt.Sprintf("pkg-%d", id),
            Password:  "password",
        })
        fmt.Println(id, result.Valid)
    }(i)
}
wg.Wait()
```
