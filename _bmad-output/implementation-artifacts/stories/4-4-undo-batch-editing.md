# Story 4.4: 撤销批量编辑

**Story ID:** 4.4
**Epic:** Epic 4 - 批量编辑与撤销
**Status:** done
**Created:** 2026-04-16

---

## Story Requirements

### User Story Statement

As a 用户，
I want 撤销批量编辑操作，
So that 修正错误的批量修改。

### Acceptance Criteria

**Given** 用户执行过批量编辑
**When** 用户点击"撤销"
**Then** 恢复最近一次批量编辑之前的状态

**And** 包括标签修改、封面修改、歌词修改

**And** 重复撤销可逐步恢复更早的操作

---

**Given** 撤销执行
**When** 系统恢复历史状态
**Then** 从 BatchOperation 表获取上一次的 OldValues

**And** 将歌曲记录恢复为 OldValues

**And** 显示 Toast 提示"已撤销"

**And** 更新列表显示恢复后的状态

---

**Given** 无可撤销的操作
**When** 用户点击"撤销"
**Then** 按钮禁用或显示"无可撤销"

**And** 提示用户没有可撤销的操作

---

## Technical Requirements

### Backend API Design

**Undo Batch Endpoint:**

```
POST /api/songs/undo/:batchId
```

Response:
```json
{
  "data": {
    "succeeded": 3,
    "failed": 0
  }
}
```

**Get Latest Batch Endpoint:**

```
GET /api/batches/latest
```

Response:
```json
{
  "data": {
    "id": 1,
    "type": "update",
    "target_ids": "[1,2,3]",
    "old_values": "{...}",
    "new_values": "{...}",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### Backend Components

**Handler Layer (`internal/handler/batch.go`):**
- `UndoBatch` 方法处理 `POST /api/songs/undo/:batchId`
- `GetLatestBatch` 方法处理 `GET /api/batches/latest`

**Repository Layer:**
- `BatchRepository.GetLatest()` - 获取最新批量操作
- `BatchRepository.Delete()` - 删除已撤销的批量操作

### Frontend Components

**SideEditPanel 撤销按钮:**
- 集成在 SideEditPanel 中
- 点击"撤销上次的批量编辑"按钮
- 显示撤销结果 Toast

### Data Model

**BatchOperation 表:**
```go
type BatchOperation struct {
    ID        uint      `gorm:"primaryKey"`
    Type      string    // "update", "delete"
    TargetIDs string    // JSON array of song IDs
    OldValues string    // JSON of previous values
    NewValues string    // JSON of new values
    CreatedAt time.Time
}
```

---

## Completion Criteria

- [x] 后端 `POST /api/songs/undo/:batchId` 端点实现
- [x] 后端 `GET /api/batches/latest` 端点实现
- [x] 撤销功能正确恢复旧值
- [x] 撤销后删除 BatchOperation 记录
- [x] 前端撤销按钮集成

---

## Dev Agent Record

### Implementation

**后端实现:**
1. 在 `internal/handler/batch.go` 添加 `UndoBatch` 和 `GetLatestBatch` 方法
2. 撤销时解析 OldValues JSON，恢复歌曲记录
3. 撤销后删除 BatchOperation 记录

**前端实现:**
1. 在 `SideEditPanel` 中添加"撤销"按钮
2. 调用 `GET /api/batches/latest` 检查是否有可撤销操作
3. 调用 `POST /api/songs/undo/:batchId` 执行撤销

### Files Modified

- `internal/handler/batch.go` - 添加 UndoBatch, GetLatestBatch
- `frontend/src/components/edit/side-edit-panel.tsx` - 添加撤销按钮

### Status

**Status:** done
