# Java SDK 接入手册

## 环境要求

- JDK 8 或更高版本
- Maven 3.6+ 或 Gradle 5+

## 安装

### Maven

```xml
<dependency>
    <groupId>com.ptm</groupId>
    <artifactId>ptm-sdk</artifactId>
    <version>1.0.0</version>
</dependency>
```

### Gradle

```groovy
implementation 'com.ptm:ptm-sdk:1.0.0'
```

## 快速开始

```java
import com.ptm.client.PTMClient;
import com.ptm.model.VerifyResult;

public class QuickStart {
    public static void main(String[] args) {
        // 初始化客户端
        PTMClient client = new PTMClient("your_api_key_here");
        
        // 验证密码
        VerifyResult result = client.verify()
            .packageId("pkg-uuid-xxx")
            .password("user_input_password")
            .execute();
        
        System.out.println("验证结果: " + result.isValid());
        if (result.isValid()) {
            System.out.println("解密密钥: " + result.getKey());
        }
    }
}
```

## 客户端初始化

### 使用 API Key 初始化

```java
import com.ptm.client.PTMClient;

PTMClient client = new PTMClient("ptm_sk_a1b2c3d4e5f6...");
```

### 使用 JWT Token 初始化

```java
PTMClient client = new PTMClient()
    .withToken("jwt_token_here");
```

### 自定义配置

```java
PTMClient client = new PTMClient()
    .withApiKey("your_api_key")
    .withBaseUrl("http://localhost:8080/api/v1")
    .withTimeout(30)
    .withMaxRetries(3);
```

## 文件包管理

### 创建文件包

```java
import java.io.File;
import com.ptm.model.Package;

// 上传文件创建
File file = new File("document.pdf");
Package pkg = client.packages()
    .create()
    .name("敏感文档")
    .file(file)
    .format(Package.Format.EXE)
    .execute();

System.out.println("创建成功: " + pkg.getId());

// 通过 URL 创建
Package pkgFromUrl = client.packages()
    .create()
    .name("远程文件")
    .sourceUrl("https://example.com/file.pdf")
    .format(Package.Format.ZIP)
    .execute();
```

### 列出文件包

```java
import com.ptm.model.PackageList;

// 分页列出
PackageList list = client.packages()
    .list()
    .page(1)
    .pageSize(20)
    .execute();

for (Package pkg : list.getItems()) {
    System.out.println(pkg.getName() + " - " + pkg.getStatus());
}

// 筛选
PackageList activeExe = client.packages()
    .list()
    .status(Package.Status.ACTIVE)
    .format(Package.Format.EXE)
    .execute();
```

### 获取文件包详情

```java
Package pkg = client.packages()
    .get("pkg-uuid-xxx")
    .execute();

System.out.println("名称: " + pkg.getName());
System.out.println("格式: " + pkg.getFormat());
System.out.println("密码数量: " + pkg.getPasswords().size());
```

### 更新文件包

```java
client.packages()
    .update("pkg-uuid-xxx")
    .name("更新后的名称")
    .description("新的描述")
    .execute();
```

### 删除文件包

```java
client.packages()
    .delete("pkg-uuid-xxx")
    .execute();
```

### 下载文件包

```java
// 下载到文件
client.packages()
    .download("pkg-uuid-xxx")
    .savePath("output.exe")
    .execute();

// 获取下载链接
DownloadUrl url = client.packages()
    .getDownloadUrl("pkg-uuid-xxx")
    .expiresIn(3600)
    .execute();

System.out.println("下载链接: " + url.getDownloadUrl());
```

## 密码策略管理

### 添加密码

```java
import java.time.LocalDateTime;
import com.ptm.model.Password;

Password pwd = client.passwords()
    .create("pkg-uuid-xxx")
    .password("Week1Password!")
    .priority(1)
    .validFrom(LocalDateTime.of(2024, 1, 1, 0, 0))
    .validUntil(LocalDateTime.of(2024, 1, 8, 0, 0))
    .execute();
```

### 批量添加密码

```java
import java.time.LocalDateTime;
import java.util.Arrays;
import com.ptm.model.PasswordBatchRequest;

LocalDateTime start = LocalDateTime.of(2024, 1, 1, 0, 0);

PasswordBatchRequest.PasswordItem[] passwords = {
    PasswordBatchRequest.PasswordItem.builder()
        .password("Week1Password!")
        .priority(1)
        .validFrom(start)
        .validUntil(start.plusWeeks(1))
        .build(),
    PasswordBatchRequest.PasswordItem.builder()
        .password("Week2Password!")
        .priority(2)
        .validFrom(start.plusWeeks(1))
        .validUntil(start.plusWeeks(2))
        .build()
};

PasswordBatchResult result = client.passwords()
    .createBatch("pkg-uuid-xxx")
    .passwords(Arrays.asList(passwords))
    .execute();

System.out.println("成功创建 " + result.getCreatedCount() + " 个密码");
```

### 激活/停用密码

```java
// 激活
client.passwords()
    .activate("pwd-uuid-xxx")
    .execute();

// 停用
client.passwords()
    .deactivate("pwd-uuid-xxx")
    .execute();
```

### 获取当前有效密码

```java
Password current = client.passwords()
    .getCurrent("pkg-uuid-xxx")
    .execute();

System.out.println("当前密码: " + current.getPassword());
```

## 验证服务

### 验证密码

```java
VerifyResult result = client.verify()
    .packageId("pkg-uuid-xxx")
    .password("user_input")
    .execute();

if (result.isValid()) {
    System.out.println("解密密钥: " + result.getKey());
} else {
    System.out.println("验证失败: " + result.getMessage());
}
```

### 批量验证

```java
VerifyResult result = client.verify()
    .packageId("pkg-uuid-xxx")
    .passwords(Arrays.asList("pwd1", "pwd2", "pwd3"))
    .execute();

if (result.isValid()) {
    System.out.println("匹配到密码: " + result.getMatchedPassword());
}
```

### 检查包状态

```java
PackageStatus status = client.verify()
    .getStatus("pkg-uuid-xxx")
    .execute();

System.out.println("状态: " + status.getStatus());
System.out.println("有效密码数: " + status.getCurrentPasswordCount());
```

## 审计日志

### 查询日志

```java
import com.ptm.model.AuditLogList;

// 查询所有日志
AuditLogList logs = client.audit()
    .list()
    .page(1)
    .pageSize(20)
    .execute();

// 按条件筛选
AuditLogList logsFiltered = client.audit()
    .list()
    .action("VERIFY_SUCCESS")
    .packageId("pkg-uuid-xxx")
    .startTime(LocalDateTime.of(2024, 1, 1, 0, 0))
    .execute();

for (AuditLog log : logsFiltered.getItems()) {
    System.out.println(log.getAction() + " - " + log.getCreatedAt());
}
```

### 导出日志

```java
// 导出为 JSON
String jsonLogs = client.audit()
    .export()
    .format("json")
    .startTime(LocalDateTime.of(2024, 1, 1, 0, 0))
    .endTime(LocalDateTime.of(2024, 1, 31, 23, 59))
    .execute();

// 导出为 CSV 到文件
client.audit()
    .export()
    .format("csv")
    .savePath("audit_log.csv")
    .startTime(LocalDateTime.of(2024, 1, 1, 0, 0))
    .endTime(LocalDateTime.of(2024, 1, 31, 23, 59))
    .execute();
```

## 错误处理

```java
import com.ptm.client.PTMClient;
import com.ptm.exceptions.*;

try {
    PTMClient client = new PTMClient("invalid_key");
    VerifyResult result = client.verify()
        .packageId("pkg-xxx")
        .password("password")
        .execute();
} catch (AuthenticationException e) {
    System.out.println("认证失败，请检查 API Key");
} catch (PermissionException e) {
    System.out.println("权限不足");
} catch (NotFoundException e) {
    System.out.println("资源不存在: " + e.getResourceId());
} catch (ValidationException e) {
    System.out.println("参数错误: " + e.getMessage());
} catch (PTMException e) {
    System.out.println("其他错误: " + e.getMessage());
}
```

## 完整示例

```java
import com.ptm.client.PTMClient;
import com.ptm.model.*;
import java.time.LocalDateTime;
import java.util.Arrays;

public class FullExample {
    public static void main(String[] args) {
        // 1. 初始化
        System.out.println("1. 初始化客户端...");
        PTMClient client = new PTMClient("your_api_key");
        
        // 2. 创建文件包
        System.out.println("2. 创建文件包...");
        Package pkg = client.packages()
            .create()
            .name("Q4财务报告")
            .file(new java.io.File("sensitive_file.pdf"))
            .format(Package.Format.EXE)
            .execute();
        System.out.println("   文件包 ID: " + pkg.getId());
        
        // 3. 批量添加密码策略
        System.out.println("3. 添加密码策略...");
        LocalDateTime start = LocalDateTime.of(2024, 1, 1, 0, 0);
        PasswordItem[] passwords = new PasswordItem[4];
        for (int week = 0; week < 4; week++) {
            passwords[week] = PasswordItem.builder()
                .password("Week" + (week+1) + "#Pass!")
                .priority(week + 1)
                .validFrom(start.plusWeeks(week))
                .validUntil(start.plusWeeks(week + 1))
                .build();
        }
        
        PasswordBatchResult batchResult = client.passwords()
            .createBatch(pkg.getId())
            .passwords(Arrays.asList(passwords))
            .execute();
        System.out.println("   成功创建 " + batchResult.getCreatedCount() + " 个密码策略");
        
        // 4. 下载文件包
        System.out.println("4. 下载文件包...");
        client.packages()
            .download(pkg.getId())
            .savePath("Q4_Report.exe")
            .execute();
        System.out.println("   下载完成");
        
        // 5. 查询审计日志
        System.out.println("5. 审计日志...");
        AuditLogList logs = client.audit()
            .list()
            .packageId(pkg.getId())
            .pageSize(10)
            .execute();
        System.out.println("   记录数: " + logs.getTotal());
        
        System.out.println("\n全部完成!");
    }
}
```

## Spring Boot 集成

### 配置

```yaml
# application.yml
ptm:
  api-key: ${PTM_API_KEY}
  base-url: http://localhost:8080/api/v1
  timeout: 30
```

### 配置类

```java
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class PTMConfig {
    
    @Bean
    @ConfigurationProperties(prefix = "ptm")
    public PTMProperties ptmProperties() {
        return new PTMProperties();
    }
    
    @Bean
    public PTMClient ptmClient(PTMProperties properties) {
        return new PTMClient()
            .withApiKey(properties.getApiKey())
            .withBaseUrl(properties.getBaseUrl())
            .withTimeout(properties.getTimeout());
    }
}
```

### 使用

```java
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@RestController
public class FileController {
    
    @Autowired
    private PTMClient ptmClient;
    
    @PostMapping("/verify")
    public VerifyResult verify(@RequestParam String packageId, 
                               @RequestParam String password) {
        return ptmClient.verify()
            .packageId(packageId)
            .password(password)
            .execute();
    }
}
```
