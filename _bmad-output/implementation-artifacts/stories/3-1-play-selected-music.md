# Story 3.1: 播放选中音乐

**Story ID:** 3.1
**Epic:** Epic 3 - 播放器与现场编辑
**Status:** done
**Created:** 2026-04-16

---

## Story Requirements

### User Story Statement

As a 用户，
I want 播放选中的音乐文件，
So that 试听歌曲并验证元数据。

### Acceptance Criteria

**Given** 用户在歌曲列表中
**When** 用户点击歌曲的播放按钮
**Then** 启动音频播放器播放该歌曲

**And** 播放器显示在界面固定位置（右侧 SidePlayer）

**And** 点一首播一首，无播放队列

**And** 播放时列表中该歌曲显示播放中状态

---

**Given** 歌曲正在播放
**When** 用户点击"暂停"按钮
**Then** 暂停播放

**And** 播放按钮变为"播放"图标

**And** 进度条停止

---

**Given** 歌曲正在播放或暂停
**When** 用户点击"播放"按钮
**Then** 恢复播放

**And** 进度条继续更新

---

## Technical Requirements

### Backend API Design

**Streaming Endpoint:**

```
GET /api/songs/:id/stream
```

Response:
- Content-Type: audio/mpeg (or based on file extension)
- Streams the audio file content

### Backend Components

**Handler Layer (`internal/handler/song.go`):**
- 添加 `StreamSong` 方法处理 `GET /api/songs/:id/stream`
- 根据文件扩展名返回正确的 Content-Type

**Routes (`cmd/server/main.go`):**
- 注册 `GET /api/songs/:id/stream` 路由

### Frontend Components

**SidePlayer 组件** (`frontend/src/components/player/side-player.tsx`):
- 新组件
- 展示封面、歌名、艺术家、专辑
- 播放/暂停控制
- 进度条和时间显示
- 音量控制
- 编辑按钮

### File Structure

```
frontend/src/
├── components/
│   ├── player/
│   │   └── side-player.tsx    # 新增
├── views/
│   └── ...                    # 修改 - 集成播放器
```

---

## Completion Criteria

- [x] 后端 `GET /api/songs/:id/stream` 流式端点实现
- [x] SidePlayer 组件实现
- [x] 播放/暂停控制功能
- [x] 进度条和时间显示

---

## Dev Agent Record

### Implementation

**后端实现:**
1. 在 `internal/handler/song.go` 添加 `StreamSong` 方法
2. 支持多种音频格式（MP3, FLAC, OGG, M4A, WAV, APE）
3. 在 `cmd/server/main.go` 注册路由

**前端实现:**
1. 创建 `frontend/src/components/player/side-player.tsx` 组件
2. 使用 HTML5 Audio API 播放音乐
3. 实现播放/暂停、进度拖动、音量控制

### Files Modified/Created

- `internal/handler/song.go` - 添加 StreamSong 方法
- `cmd/server/main.go` - 注册流式路由
- `frontend/src/components/player/side-player.tsx` - 新增

### Status

**Status:** done
