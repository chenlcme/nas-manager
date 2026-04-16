# Story 1.7: 孤立记录清理与扫描异常处理

Status: done

## Story

As a 用户，
I want 系统自动清理孤立的数据库记录并优雅处理扫描异常，
So that 音乐库始终保持一致且扫描过程不会因错误中断。

## Acceptance Criteria

1. **Given** 数据库中存在歌曲记录 **When** 对应的音乐文件已被删除 **Then** 系统自动标记并清理该记录

2. **And** 孤立记录清理不影響其他正常记录

3. **Given** 扫描过程中遇到权限不足的文件 **When** 系统扫描该文件 **Then** 记录错误并继续扫描其他文件

4. **And** 不中断整个扫描过程

5. **Given** 扫描过程中遇到损坏的音乐文件 **When** 系统尝试解析 ID3 标签 **Then** 跳过该文件并记录错误日志

6. **And** 继续处理其他文件

7. **Given** 扫描过程中遇到无法识别的文件格式 **When** 系统判断文件类型 **Then** 跳过该文件并继续

8. **And** 提供友好的错误摘要信息

## Tasks / Subtasks

- [x] Task 1: 实现孤立记录清理功能 (AC: 1-2)
  - [x] 扫描前获取所有数据库记录的文件路径
  - [x] 对比文件系统，发现不存在的文件
  - [x] 删除孤立记录

- [x] Task 2: 更新扫描异常处理 (AC: 3-8)
  - [x] 权限错误处理
  - [x] 损坏文件处理
  - [x] 未知格式处理
  - [x] 错误聚合与报告

- [x] Task 3: 添加孤立清理 API (AC: 1-2)
  - [x] POST /api/songs/cleanup - 手动触发清理
  - [x] 返回清理结果

- [x] Task 4: 编写测试 (AC: 1-8)
  - [x] 测试孤立记录清理
  - [x] 测试扫描异常处理

## Dev Notes

### 技术要求

**孤立记录清理：**
- 在扫描前或扫描后执行
- 比对数据库记录与实际文件系统
- 标记已删除文件的记录
- 提供手动清理接口

**扫描异常处理：**
- 权限错误：记录并继续
- 损坏文件：记录跳过并继续
- 未知格式：记录并继续
- 所有错误聚合到最终结果

### API 设计

```
POST /api/songs/cleanup
Response: { "cleaned": 5, "errors": [] }
```

## Dev Agent Record

### Agent Model Used

MiniMax-M2.7-highspeed

### Debug Log References

### Completion Notes List

- 实现 CleanupOrphanRecords 方法清理孤立记录
- 添加 POST /api/songs/cleanup 端点
- 扫描异常处理已在 scanner.go 中实现（GetMusicFilesWithErrors）

### File List

1. `internal/service/scanner.go` - 添加 CleanupOrphanRecords
2. `internal/handler/scan.go` - 添加 Cleanup 端点
3. `cmd/server/main.go` - 注册 cleanup 路由

## Change Log

- 2026-04-15: 初始实现 Story 1.7 所有任务

### Review Findings

- [x] [Review][Defer] `WalkDir` 符号链接安全 [internal/service/scanner.go] — deferred，作为 Story 1.4 的通用修复，已修复但标记为 deferred
