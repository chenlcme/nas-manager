# Story 1.3: 加密密码设置

Status: review

## Story

As a 用户，
I want 设置加密密码来保护敏感数据，
So that 云存储凭证等敏感信息可以被安全存储。

## Acceptance Criteria

1. **Given** 用户已完成首次配置 **When** 用户在设置页面设置加密密码 **Then** 密码长度 ≥ 8 字符

2. **And** 使用 PBKDF2 派生加密密钥（100000 次迭代，32 字节输出）

3. **And** 密钥和盐值存储在 settings 表中（盐值随机生成）

4. **And** 后续凭证加密使用 AES-256-GCM 或 ChaCha20-Poly1305

5. **Given** 用户已设置加密密码 **When** 用户修改加密密码 **Then** 验证原密码正确性

6. **And** 重新派生密钥

7. **And** 使用新密钥重新加密已有凭证

8. **And** 原密钥加密的数据无法被新密钥解密时，提示用户重新配置凭证

## Tasks / Subtasks

- [x] Task 1: 创建加密工具包 (AC: 1-8)
  - [x] 创建 `pkg/crypto/encrypt.go` - 加密解密工具
  - [x] 实现 PBKDF2 密钥派生
  - [x] 实现 ChaCha20-Poly1305 加密
  - [x] 实现盐值生成

- [x] Task 2: 创建加密服务层 (AC: 1-8)
  - [x] 创建 `internal/service/encrypt.go` - 加密服务
  - [x] 实现 SetupPassword 设置密码
  - [x] 实现 VerifyPassword 验证密码
  - [x] 实现 ChangePassword 修改密码

- [x] Task 3: 创建加密 Handler 和 API (AC: 1-8)
  - [x] 创建 `internal/handler/encrypt.go` - 加密处理器
  - [x] 实现 POST /api/auth/setup 设置密码
  - [x] 实现 POST /api/auth/verify 验证密码
  - [x] 实现 POST /api/auth/change 修改密码

- [x] Task 4: 添加加密相关设置常量 (AC: 3)
  - [x] 在 settings 表中存储加密盐值和验证值

- [x] Task 5: 编写加密模块单元测试 (AC: 1-8)
  - [x] 测试 PBKDF2 密钥派生
  - [x] 测试加密解密功能
  - [x] 测试密码验证

## Dev Notes

### 技术要求

**加密实现：**
- 使用 `golang.org/x/crypto` 标准库
- PBKDF2 派生密钥：100000 次迭代，SHA256，32 字节输出
- 对称加密：ChaCha20-Poly1305
- 盐值：32 字节随机数

### PBKDF2 规范

```go
func DeriveKey(password string, salt []byte) []byte {
    return pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)
}
```

### 存储结构

Settings 表存储：
- `crypto_salt`: 加密盐值（base64 编码）
- `crypto_verify`: 验证值（使用派生密钥加密固定字符串的结果）

### API 设计

```
POST /api/auth/setup
Body: { "password": "..." }
Response: { "success": true }

POST /api/auth/verify
Body: { "password": "..." }
Response: { "valid": true/false }

POST /api/auth/change
Body: { "old_password": "...", "new_password": "..." }
Response: { "success": true }
```

## Dev Agent Record

### Agent Model Used

MiniMax-M2.7-highspeed

### Completion Notes List

- 创建了 `pkg/crypto/encrypt.go` 实现 ChaCha20-Poly1305 加密解密
- 使用 PBKDF2 派生密钥（100000 次迭代，32 字节）
- 创建了 `internal/service/encrypt.go` 实现密码设置、验证、修改
- 创建了 `internal/handler/encrypt.go` 实现 API 端点
- 添加了加密模块和服务层单元测试，全部通过
- 密码验证使用派生密钥加密验证字符串的方式

## File List

1. `pkg/crypto/encrypt.go` - 加密工具实现
2. `pkg/crypto/encrypt_test.go` - 加密工具单元测试
3. `internal/service/encrypt.go` - 加密服务实现
4. `internal/service/encrypt_test.go` - 加密服务单元测试
5. `internal/handler/encrypt.go` - 加密处理器实现

## Change Log

- 2026-04-15: 初始实现 Story 1.3 所有任务

### Review Findings

- [x] [Review][Patch] `/api/auth/*` 路由未注册 [cmd/server/main.go] — 已修复
- [x] [Review][Defer] ChangePassword 未重新加密已有凭证 [internal/service/encrypt.go:833] — deferred，需等云存储功能实现（TODO 注释）
