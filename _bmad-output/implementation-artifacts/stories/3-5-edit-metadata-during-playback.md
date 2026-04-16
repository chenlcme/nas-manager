# Story 3.5: 播放中编辑元数据

**Story ID:** 3.5
**Epic:** Epic 3 - 播放器与现场编辑
**Status:** done
**Created:** 2026-04-16

---

## Story Requirements

### User Story Statement

As a 用户，
I want 在播放过程中直接编辑歌曲元数据，
So that 边听边验证边修改，高效完成元数据整理。

### Acceptance Criteria

**Given** 歌曲正在播放
**When** 用户点击"编辑"按钮
**Then** 播放器展开编辑面板

**And** 显示当前歌曲的可编辑字段（标题、艺术家、专辑、年份、流派）

**And** 封面和歌词也可编辑

**And** 已修改的字段高亮显示

---

**Given** 用户修改了字段
**When** 用户点击"保存"
**Then** 验证输入有效性

**And** 更新数据库中的歌曲记录

**And** 刷新播放器显示新内容

**And** 继续播放（不中断）

**And** 显示 Toast 提示"已保存"

---

**Given** 用户修改了字段
**When** 用户点击"取消"
**Then** 放弃修改

**And** 恢复原始显示

---

## Technical Requirements

### Backend API Design

**Update Song Endpoint:**

```
PUT /api/songs/:id
Content-Type: application/json

Request Body:
{
  "title": "新标题",
  "artist": "新艺术家",
  "album": "新专辑",
  "year": 2024,
  "genre": "新流派",
  "trackNum": 1,
  "coverPath": "封面路径",
  "lyrics": "歌词内容"
}

Response:
{
  "data": { ... updated song ... }
}
```

### Backend Components

**Handler Layer (`internal/handler/song.go`):**
- 添加 `UpdateSong` 方法处理 `PUT /api/songs/:id`
- 支持部分更新（nil 字段保持不变）

**Routes (`cmd/server/main.go`):**
- 注册 `PUT /api/songs/:id` 路由

### Frontend Components

**SidePlayer 编辑面板:**
- 在 SidePlayer 组件中集成编辑表单
- Tab 切换：标签信息 / 歌词
- 保存和取消按钮

### File Structure

```
frontend/src/
├── components/
│   └── player/
│       └── side-player.tsx    # 修改 - 添加编辑功能
```

---

## Completion Criteria

- [x] 后端 `PUT /api/songs/:id` 端点实现
- [x] SidePlayer 集成编辑面板
- [x] 保存和取消功能
- [x] Toast 提示

---

## Dev Agent Record

### Implementation

**后端实现:**
1. 在 `internal/handler/song.go` 添加 `UpdateSong` 方法
2. 支持部分更新（只更新非 nil 字段）
3. 在 `cmd/server/main.go` 注册路由

**前端实现:**
1. 在 `side-player.tsx` 中添加编辑面板
2. 使用 useState 管理编辑表单状态
3. 调用 PUT API 保存更改

### Files Modified

- `internal/handler/song.go` - 添加 UpdateSong 方法
- `cmd/server/main.go` - 注册 PUT 路由
- `frontend/src/components/player/side-player.tsx` - 添加编辑功能

### Status

**Status:** done
