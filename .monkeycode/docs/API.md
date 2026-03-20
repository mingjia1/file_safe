# API 概述

## 基本信息

| 项目 | 说明 |
|------|------|
| 基础 URL | `http://{host}:{port}/api/v1` |
| 数据格式 | JSON |
| 字符编码 | UTF-8 |

## 认证方式

### 1. JWT Token 认证（Web 管理平台）

用于 Web 前端登录后调用 API。

**请求头：**
```
Authorization: Bearer <jwt_token>
```

**登录接口：**
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "your_password"
}
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2024-01-08T00:00:00Z"
  }
}
```

### 2. API Key 认证（SDK/集成接入）

用于服务端或客户端程序调用 API。

**请求头：**
```
X-API-Key: <your_api_key>
```

**获取 API Key：**
通过 Web 管理平台的用户设置页面生成和管理 API Key。

## 通用响应格式

### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

### 错误响应

```json
{
  "code": 10001,
  "message": "error description",
  "data": null
}
```

## 错误码定义

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 10001 | 参数错误 |
| 10002 | 认证失败 |
| 10003 | 权限不足 |
| 10004 | 资源不存在 |
| 10005 | 资源已存在 |
| 10006 | 操作失败 |
| 20001 | 文件不存在 |
| 20002 | 文件上传失败 |
| 20003 | 文件包生成失败 |
| 30001 | 密码错误 |
| 30002 | 密码已过期 |
| 30003 | 密码未生效 |
| 30004 | 无有效密码 |
| 40001 | 加密配置无效 |
| 40002 | 加密/解密失败 |

## 分页查询

列表接口支持分页查询：

**请求参数：**
```
page=1
page_size=20
```

**响应格式：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [...],
    "total": 100,
    "page": 1,
    "page_size": 20,
    "total_pages": 5
  }
}
```

## 公共请求头

所有 API 请求都需要以下请求头：

| 请求头 | 说明 | 必填 |
|--------|------|------|
| Content-Type | application/json | 是 |
| Authorization | Bearer {token} 或 X-API-Key | 是 |

## 健康检查

```http
GET /api/v1/health
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "healthy",
    "version": "1.0.0"
  }
}
```
