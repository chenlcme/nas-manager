# Story 1.6: 增量扫描与全量扫描

Status: done

## Story

As a 用户，
I want 选择增量扫描或全量扫描，
So that 高效更新音乐库而不重复处理未变化的歌曲。

## Acceptance Criteria

1. **Given** 用户触发重新扫描 **When** 用户选择"全量扫描" **Then** 重新解析所有音乐文件的 ID3 标签

2. **And** 更新数据库中所有相关记录

3. **And** 保留已有的用户编辑数据（如有）

4. **Given** 用户触发重新扫描 **When** 用户选择"增量扫描" **Then** 仅处理修改时间晚于上次扫描时间的文件

5. **And** 新文件：创建新记录

6. **And** 修改过的文件：更新元数据

7. **And** 删除过的文件：不处理（由孤岛清理处理）

8. **And** 增量扫描耗时 < 全量扫描的 10%

## Tasks / Subtasks

- [x] Task 1: 更新设置表存储扫描时间 (AC: 1-8)
  - [x] 添加 last_scan_time 设置项
  - [x] 更新扫描后保存时间戳

- [x] Task 2: 更新扫描服务支持全量/增量模式 (AC: 1-3)
  - [x] 修改 scanner.go 支持全量扫描
  - [x] 实现 ForceUpdateSong 强制更新

- [x] Task 3: 实现增量扫描逻辑 (AC: 4-8)
  - [x] 获取上次扫描时间
  - [x] 比较文件修改时间
  - [x] 只处理新增和修改的文件

- [x] Task 4: 更新扫描 API 支持扫描模式参数 (AC: 1-8)
  - [x] 修改 POST /api/songs/scan 支持 mode 参数
  - [x] GET 参数或 POST body

- [x] Task 5: 编写增量扫描单元测试 (AC: 4-8)
  - [x] 测试增量扫描逻辑

## Dev Notes

### 技术要求

**全量扫描：**
- 遍历所有音乐文件
- 重新解析每个文件的 ID3 标签
- 更新数据库记录

**增量扫描：**
- 记录上次扫描时间戳
- 只处理修改时间 > 上次扫描时间的文件
- 新文件直接创建
- 修改过的文件更新元数据

### API 设计

```
POST /api/songs/scan
Body: { "mode": "full" | "incremental" }
Response: { "found": 100, "new": 10, "updated": 5, "errors": [] }
```

### 存储结构

Settings 表：
- `last_scan_time`: 上次扫描时间戳（Unix）

## Dev Agent Record

### Agent Model Used

MiniMax-M2.7-highspeed

### Debug Log References

### Completion Notes List

- 实现全量/增量扫描模式
- 更新 scanner.go 添加 ScanMode 和 ScanFiles 方法
- 更新 scan.go handler 支持 mode 参数
- 添加 GetFileModTime 和 SetLastScanTime 方法
- 添加增量扫描单元测试

### File List

1. `internal/service/scanner.go` - 添加增量/全量扫描支持
2. `internal/service/scanner_test.go` - 添加增量扫描测试
3. `internal/handler/scan.go` - 更新支持 mode 参数
4. `internal/repository/setting.go` - 添加 GetLastScanTime/SetLastScanTime

## Change Log

- 2026-04-15: 初始实现 Story 1.6 所有任务

### Review Findings

- [x] [Review][Defer] `lastScanTime==0` 语义歧义（"从未扫描" vs "Unix时间0"）[internal/service/scanner.go] — deferred，需添加 `hasScannedBefore` 布尔标志
