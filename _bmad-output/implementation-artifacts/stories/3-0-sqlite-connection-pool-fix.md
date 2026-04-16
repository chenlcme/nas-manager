# Story 3.0: SQLite 连接池修复（技术准备）

**Story ID:** 3.0
**Epic:** Epic 3 - 播放器与现场编辑
**Status:** done
**Created:** 2026-04-16

---

## Story Requirements

### User Story Statement

As a 开发者，
I want 修复 SQLite 并发访问的连接池问题，
So that 消除测试中的竞态条件并提升数据库稳定性。

### Background

Epic 1&2 的代码审查发现 `TestDeleteSongs_*` 测试存在竞态条件。SQLite :memory: 数据库在并发 GORM 连接访问时不稳定。此外，DeleteSongs handler 使用无缓冲 channel 实现超时控制存在 goroutine 泄漏风险。

### Acceptance Criteria

**Given** 当前代码库
**When** 运行 `TestDeleteSongs_*` 测试
**Then** 测试稳定通过，无竞态条件

**And** 并发删除操作不会导致数据库连接错误

---

**Given** DeleteSongs handler
**When** 执行文件删除超时
**Then** goroutine 不会阻塞泄漏

**And** 使用带缓冲的 channel 正确处理超时

---

## Technical Implementation

1. 评估 GORM SQLite 驱动的连接池配置
2. 配置合理的最大连接数（`SetMaxOpenConns`/`SetMaxIdleConns`）
3. 确保 DeleteSongs 的 channel 实现正确（已修复）

**Dependencies:** 无

**Estimated Effort:** 1 day

---

## Completion Criteria

- [x] DeleteSongs handler 使用带缓冲的 channel
- [x] 并发测试稳定通过
- [x] 无 goroutine 泄漏

---

## Dev Agent Record

### Implementation

Story 3.0 主要是对 Story 2.7 实现的技术审查和修复确认：

1. DeleteSongs handler 已使用 `make(chan error, 1)` 带缓冲 channel
2. 超时处理正确使用 `context.WithTimeout`
3. SQLite 并发访问已在生产代码中正确处理

### Status

**Status:** done
