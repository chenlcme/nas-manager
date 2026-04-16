# Story 4.1: 批量修改标签

**Story ID:** 4.1
**Epic:** Epic 4 - 批量编辑与撤销
**Status:** done
**Created:** 2026-04-16

---

## Story Requirements

### User Story Statement

As a 用户，
I want 批量修改选中歌曲的标签，
So that 高效整理一批歌曲的元数据。

### Acceptance Criteria

**Given** 用户选中了 N 首歌曲
**When** 用户点击"批量编辑"
**Then** 右侧滑入 SideEditPanel（宽度 320px）

**And** 显示"已选中 N 首"

**And** 显示选中歌曲列表预览

---

**Given** 批量编辑面板展开
**When** 用户填写标签字段（艺术家/专辑/年份/流派）
**Then** 留空的字段保持不变

**And** 填写的字段显示预览值

**And** 预览区域显示将要做的更改

---

**Given** 用户确认预览
**When** 用户点击"应用"
**Then** 批量更新数据库中的歌曲记录

**And** 创建 BatchOperation 记录用于撤销

**And** 关闭编辑面板

**And** 显示 Toast 提示"已更新 N 首"

**And** 保持选中状态（可继续操作）

---

## Technical Requirements

### Backend API Design

**Batch Update Endpoint:**

```
POST /api/songs/batch-update
Content-Type: application/json

Request Body:
{
  "ids": [1, 2, 3],
  "title": "新标题",        // optional
  "artist": "新艺术家",     // optional
  "album": "新专辑",       // optional
  "year": 2024,            // optional
  "genre": "新流派",       // optional
  "trackNum": 1,           // optional
  "coverPath": "...",      // optional
  "lyrics": "..."          // optional
}

Response:
{
  "data": {
    "total": 3,
    "succeeded": 3,
    "failed": 0
  }
}
```

### Backend Components

**Handler Layer (`internal/handler/batch.go`):**
- 新文件
- `BatchUpdate` 方法处理 `POST /api/songs/batch-update`
- 保存 BatchOperation 记录用于撤销

**Repository Layer:**
- `internal/repository/batch.go` - 新文件
- BatchOperation CRUD 操作

### Frontend Components

**SideEditPanel 组件** (`frontend/src/components/edit/side-edit-panel.tsx`):
- 新组件
- 右侧滑入面板
- Tab 切换：标签信息 / 歌词
- 表单字段：艺术家、专辑、年份、流派
- 预览区域
- 应用和撤销按钮

### File Structure

```
frontend/src/
├── components/
│   ├── edit/
│   │   └── side-edit-panel.tsx    # 新增
```

---

## Completion Criteria

- [x] 后端 `POST /api/songs/batch-update` 端点实现
- [x] SideEditPanel 组件实现
- [x] 批量更新逻辑
- [x] BatchOperation 记录保存
- [x] Toast 提示

---

## Dev Agent Record

### Implementation

**后端实现:**
1. 创建 `internal/repository/batch.go` - BatchOperation 数据访问
2. 创建 `internal/handler/batch.go` - BatchUpdate 和 UndoBatch 处理器
3. 在 `cmd/server/main.go` 注册路由

**前端实现:**
1. 创建 `frontend/src/components/edit/side-edit-panel.tsx` 组件
2. 使用 SelectionContext 获取选中歌曲
3. 调用 batch-update API

### Files Created/Modified

- `internal/repository/batch.go` - 新增
- `internal/handler/batch.go` - 新增
- `cmd/server/main.go` - 注册新路由
- `frontend/src/components/edit/side-edit-panel.tsx` - 新增

### Status

**Status:** done
