# Password Timer Manager 文档中心

## 概述

Password Timer Manager 是一个企业级文件安全分享平台，支持生成带有多密码时效控制的文件包。管理员可以通过 Web 平台或 API 管理密码策略，实现文件的安全分发与访问控制。

## 核心特性

- 多密码时效管理：支持设置多个密码，每个密码有独立的有效期
- 混合验证模式：优先在线验证，网络异常时自动切换本地验证
- 多种输出格式：支持 EXE 自解压程序和 ZIP + 内嵌验证程序
- 防篡改机制：数字签名 + 加密存储，防止配置被修改
- 完整审计日志：记录所有下载、验证和管理操作

---

## 文档目录

### 快速开始
- [快速开始指南](./QUICKSTART.md) - 5 分钟快速上手
- [系统架构](./ARCHITECTURE.md) - 系统架构与组件
- [部署指南](./DEPLOY.md) - Linux/Windows 一键部署

### API 文档
- [API 概述](./API.md) - API 基础信息
- [认证授权](./API_AUTH.md) - 认证方式与权限控制
- [文件包管理](./API_PACKAGE.md) - 文件包 CRUD 操作
- [密码策略管理](./API_PASSWORD.md) - 密码策略管理接口
- [验证服务](./API_VERIFY.md) - 密码验证接口
- [审计日志](./API_AUDIT.md) - 审计日志接口
- [系统管理](./API_ADMIN.md) - 系统配置接口

### SDK 接入
- [Python SDK](./SDK_PYTHON.md) - Python 语言接入
- [Java SDK](./SDK_JAVA.md) - Java 语言接入
- [Go SDK](./SDK_GO.md) - Go 语言接入

### 客户端使用
- [自解压程序用户指南](./CLIENT_EXE.md) - EXE 格式使用说明
- [ZIP 验证程序用户指南](./CLIENT_ZIP.md) - ZIP 格式使用说明

---

## 版本信息

当前版本：1.0.0

最后更新：2026-03-20
