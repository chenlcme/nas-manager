# Story 2.5: 查看单曲详情

Status: done

## Story

As a 用户，
I want 查看单曲的详细信息，
So that 了解歌曲的完整元数据。

## Acceptance Criteria

1. **Given** 用户在歌曲列表中 **When** 用户点击某首歌曲 **Then** 显示歌曲详情面板或弹窗

2. **And** 展示完整 ID3 信息（标题、艺术家、专辑、年份、流派、曲目号）

3. **And** 展示封面图片

4. **And** 展示内嵌歌词

5. **And** 展示文件信息（路径、大小、时长、格式）

## Tasks / Subtasks

- [x] Task 1: 实现后端单曲详情 API (AC: 1-5)
  - [x] 创建 `internal/handler/song.go` - SongHandler
  - [x] 实现 `GET /api/songs/:id` 端点获取单曲详情
  - [x] 在 router 中注册路由
  - [x] 添加错误处理（歌曲不存在等）

- [x] Task 2: 创建前端单曲详情组件 (AC: 1-5)
  - [x] 创建 `frontend/src/components/song/song-detail-panel.tsx` - 详情面板组件
  - [x] 实现 ID3 信息展示（标题、艺术家、专辑、年份、流派、曲目号）
  - [x] 实现封面图片展示
  - [x] 实现内嵌歌词展示
  - [x] 实现文件信息展示（路径、大小、时长、格式）

- [x] Task 3: 集成详情面板到视图 (AC: 1-5)
  - [x] 在 App 组件中管理详情面板显示状态
  - [x] 在 SongTableRow 添加查看详情按钮
  - [x] 将详情面板集成到 ArtistsView, AlbumsView, FoldersView

- [x] Task 4: 编写测试 (AC: 1-5)
  - [x] 后端: SongHandler 单曲详情 API 单元测试

## Dev Notes

### 技术要求

**API 设计:**
```
GET /api/songs/:id
Response:
{
  "data": {
    "id": 1,
    "file_path": "/path/to/song.mp3",
    "title": "歌曲名",
    "artist": "艺术家",
    "album": "专辑",
    "year": 2024,
    "genre": "流行",
    "track_num": 1,
    "duration": 240,
    "cover_path": "/path/to/cover.jpg",
    "lyrics": "歌词内容",
    "file_hash": "abc123",
    "file_size": 8563214,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**详情面板 UI:**
- 使用右侧滑入面板或模态框形式
- 面板宽度建议 400px
- 封面图片居中显示，建议 200x200px
- ID3 信息使用标签形式展示
- 文件信息在底部显示
- 支持关闭面板

**文件格式判断:**
- 根据文件扩展名判断格式 (.mp3, .flac, .ape, .ogg 等)
- 文件大小转换为人类可读格式 (KB/MB/GB)

**时长格式化:**
- 秒数转换为 mm:ss 格式

### Project Structure Notes

**需要创建的文件:**
```
nas-manager/
├── cmd/server/main.go              # 需要注册 SongHandler 路由
├── internal/
│   ├── handler/
│   │   ├── song.go               # 新建 - SongHandler
│   │   └── song_test.go          # 新建 - 测试文件
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   └── song/
│   │   │       └── song-detail-panel.tsx  # 新建 - 详情面板组件
│   │   └── app.tsx               # 修改 - 添加详情面板状态
```

**需要修改的文件:**
- `cmd/server/main.go` - 添加 SongHandler 路由注册
- `frontend/src/app.tsx` - 添加详情面板状态和渲染
- `frontend/src/views/artists-view.tsx` - 传递 onShowSongDetail 回调
- `frontend/src/views/albums-view.tsx` - 传递 onShowSongDetail 回调
- `frontend/src/views/folders-view.tsx` - 传递 onShowSongDetail 回调

### Architecture Compliance

**遵循规范:**
1. **命名规范:** API 使用 snake_case，Go 代码使用 PascalCase/camelCase
2. **响应格式:** 使用 `pkg/response/response.go` 的统一响应格式
3. **错误处理:** 使用 `pkg/response` 的错误响应格式
4. **分层架构:** Handler → Repository
5. **前端状态:** 使用 Preact Context API + 组件内局部状态

**UX 规范:**
- 色彩: 主题绿 #22C55E, 背景 #FFFFFF, 文字 #1E293B
- 字体: Noto Sans CJK (中文), Inter (英文)
- 间距: 4px 基准 (xs:4, sm:8, md:16, lg:24, xl:32)
- 触摸目标: 最小 44x44px

### Previous Story Intelligence (Stories 2.1, 2.2, 2.3, 2.4)

**经验总结:**
- ArtistRepository/AlbumRepository/FolderRepository 使用 `db.Model(&Song{}).Select(...).Group(...)` 模式
- Handler 直接调用 Repository，无需 Service 层（只读查询）
- 前端使用 Preact functional 组件 + hooks
- 统一响应格式: `{"data": [...]}`
- 错误响应格式: `{"error": {"code": "...", "message": "..."}}`
- 使用 `fetchWithTimeout` 辅助函数处理请求超时
- 使用 `AbortController` 取消旧请求

**复用的组件:**
- `song-table-row.tsx` - 高密度歌曲表格行
- `selection-bar.tsx` - 选择操作栏
- `tab-nav.tsx` - Tab 导航
- `selection-context.tsx` - 选择状态管理
- `sort-selector.tsx` - 排序选择组件

**代码模式参考:**
```go
// Handler 模式 (参考 artist.go)
func (h *ArtistHandler) GetArtistSongs(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        response.Error(c, http.StatusBadRequest, "INVALID_ID", "无效的ID")
        return
    }
    // ... 业务逻辑
    response.Success(c, songs)
}
```

### 特殊注意事项

1. **性能考虑:**
   - 歌词可能很长，需要限制显示区域高度并支持滚动
   - 封面图片加载时显示占位符

2. **错误处理:**
   - 歌曲不存在返回 404
   - 网络错误显示友好提示

3. **与其他功能的交互:**
   - 查看详情不应影响播放状态
   - 查看详情不应影响选中状态

## Dev Agent Record

### Agent Model Used

MiniMax-M2.7-highspeed

### Debug Log References

### Completion Notes List

- 创建 `internal/handler/song.go` - SongHandler，包含 GetSong 方法获取单曲详情
- 创建 `internal/handler/song_test.go` - 包含 3 个测试用例（成功、歌曲不存在、无效 ID）
- 在 `cmd/server/main.go` 注册 SongHandler 和 GET /api/songs/:id 路由
- 创建 `frontend/src/components/song/song-detail-panel.tsx` - 右侧滑入详情面板组件
- 更新 `SongTableRow` 组件添加 onShowDetail 回调和详情按钮
- 更新 `frontend/src/app.tsx` 添加 detailSong 状态和 SongDetailPanel 渲染
- 更新 `ArtistsView`, `AlbumsView`, `FoldersView` 传递 onShowSongDetail 回调
- 所有后端测试通过

### File List

- `internal/handler/song.go` - 新建
- `internal/handler/song_test.go` - 新建
- `frontend/src/components/song/song-detail-panel.tsx` - 新建
- `cmd/server/main.go` - 修改
- `frontend/src/app.tsx` - 修改
- `frontend/src/views/artists-view.tsx` - 修改
- `frontend/src/views/albums-view.tsx` - 修改
- `frontend/src/views/folders-view.tsx` - 修改
- `frontend/src/components/song/song-table-row.tsx` - 修改

## Change Log

- 2026-04-16: 创建 Story 2.5 故事文件
- 2026-04-16: 实现后端 SongHandler 和 GET /api/songs/:id API
- 2026-04-16: 创建 SongDetailPanel 前端组件
- 2026-04-16: 集成详情面板到 App 和各视图
- 2026-04-16: 添加 SongHandler 单元测试
- 2026-04-16: 所有后端测试通过
- 2026-04-16: 代码审查发现 4 个 patch 级别问题
- 2026-04-16: 修复所有 4 个 patch 问题（race condition 和 cleanup）

### Review Findings

- [x] [Review][Patch] Race: fetch success after panel close reopens panel with stale data [app.tsx:75]
- [x] [Review][Patch] Race: sequential song detail requests overwrite each other [app.tsx:75]
- [x] [Review][Patch] No AbortController cleanup on component unmount [app.tsx:75]
- [x] [Review][Patch] Toast setTimeout captures stale toasts closure [app.tsx:103]
- [x] [Review][Defer] Panel shows basic data instantly, no loading distinction [app.tsx:69] — deferred, design preference

## References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-2.5]
- [Source: _bmad-output/planning-artifacts/architecture.md#API-Communication-Patterns]
- [Source: _bmad-output/implementation-artifacts/stories/2-4-song-list-sorting.md]
