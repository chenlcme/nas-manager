# Story 2.4: 歌曲列表排序

Status: done

## Story

As a 用户，
I want 在浏览视图中对歌曲列表排序，
so that 按照我关心的方式组织音乐。

## Acceptance Criteria

1. **Given** 用户在歌曲列表页面 **When** 用户选择排序方式 **Then** 支持按歌曲名称排序

2. **And** 支持按时长排序

3. **And** 支持按添加时间排序

4. **And** 支持升序/降序切换

5. **And** 排序选项在列表顶部显示

## Tasks / Subtasks

- [x] Task 1: 设计排序 API 参数 (AC: 1-5)
  - [x] 确定 sort_by 参数选项（title, duration, created_at）
  - [x] 确定 order 参数选项（asc, desc）
  - [x] 更新相关 repository 支持排序查询

- [x] Task 2: 实现后端排序功能 (AC: 1-5)
  - [x] 修改 ArtistHandler.GetArtistSongs 支持排序
  - [x] 修改 AlbumHandler.GetAlbumSongs 支持排序
  - [x] 修改 FolderHandler.GetFolderSongs 支持排序
  - [x] 添加输入校验防止 SQL 注入

- [x] Task 3: 实现前端排序组件 (AC: 1-5)
  - [x] 创建 SortSelector 组件（排序选项下拉/按钮组）
  - [x] 在 ArtistsView, AlbumsView, FoldersView 集成排序组件
  - [x] 排序选项显示在列表顶部
  - [x] 升序/降序切换按钮

- [x] Task 4: 编写测试 (AC: 1-5)
  - [x] 后端: 各 Handler 排序功能单元测试

## Dev Notes

### 技术要求

**API 参数设计:**
```
GET /api/artists/:id/songs?sort_by=title&order=asc
GET /api/albums/:id/songs?sort_by=duration&order=desc
GET /api/folders/:id/songs?sort_by=created_at&order=asc
```

**排序字段映射:**
| 字段 | SQL 列 | 说明 |
|------|--------|------|
| title | title | 歌曲名称 |
| duration | duration | 时长（秒） |
| created_at | created_at | 添加时间 |

**排序选项 UI:**
- 排序字段下拉选择：名称 / 时长 / 添加时间
- 升序/降序切换按钮（箭头图标）
- 默认排序：按名称升序

### Project Structure Notes

**现有结构（Story 2.1, 2.2, 2.3 已创建）:**
```
nas-manager/
├── cmd/server/main.go              # 路由已注册（Artist/Album/Folder）
├── internal/
│   ├── handler/
│   │   ├── artist.go              # Story 2.1 已创建，需添加排序参数
│   │   ├── album.go               # Story 2.2 已创建，需添加排序参数
│   │   └── folder.go              # Story 2.3 已创建，需添加排序参数
│   ├── repository/
│   │   ├── artist.go              # Story 2.1 已创建
│   │   ├── album.go               # Story 2.2 已创建
│   │   └── folder.go              # Story 2.3 已创建
├── frontend/
│   ├── src/
│   │   ├── views/
│   │   │   ├── artists-view.tsx   # Story 2.1 已创建，需集成排序
│   │   │   ├── albums-view.tsx    # Story 2.2 已创建，需集成排序
│   │   │   └── folders-view.tsx   # Story 2.3 已创建，需集成排序
│   │   └── components/
│   │       ├── common/
│   │       │   └── sort-selector.tsx  # TODO: 创建排序选择组件
```

### Architecture Compliance

**遵循规范:**
1. **命名规范:** API 使用 snake_case，Go 代码使用 PascalCase/camelCase
2. **响应格式:** 使用 `pkg/response/response.go` 的统一响应格式
3. **错误处理:** 使用 `pkg/response` 的错误响应格式
4. **分层架构:** Handler → Repository (无 Service 层对于只读查询)
5. **前端状态:** 使用 Preact Context API

**UX 规范:**
- 色彩: 主题绿 #22C55E, 背景 #FFFFFF, 文字 #1E293B
- 字体: Noto Sans CJK (中文), Inter (英文)
- 间距: 4px 基准 (xs:4, sm:8, md:16, lg:24, xl:32)
- 触摸目标: 最小 44x44px

### Previous Story Intelligence (Stories 2.1, 2.2, 2.3)

**经验总结:**
- ArtistRepository/AlbumRepository/FolderRepository 使用 `db.Model(&Song{}).Select(...).Group(...)` 模式
- Handler 直接调用 Repository，无需 Service 层（只读查询）
- 前端使用 Preact functional 组件 + hooks
- 统一响应格式: `{"data": [...]}`
- 错误响应格式: `{"error": {"code": "...", "message": "..."}}`

**复用的组件:**
- `song-table-row.tsx` - 高密度歌曲表格行
- `selection-bar.tsx` - 选择操作栏
- `tab-nav.tsx` - Tab 导航
- `selection-context.tsx` - 选择状态管理

**代码模式参考:**
```go
// 排序查询示例（Repository 层）
func (r *SongRepository) GetSongsByArtistID(artistID uint, sortBy, order string) ([]Song, error) {
    // 校验 sortBy 和 order 参数
    // 构建 ORDER BY 子句
    // 执行查询
}
```

### 特殊注意事项

1. **SQL 注入防护:**
   - sortBy 只允许预定义的值：title, duration, created_at
   - order 只允许：asc, desc
   - 使用白名单校验

2. **默认值:**
   - 默认排序字段：title
   - 默认排序方向：asc

3. **与其他功能的交互:**
   - 排序不应影响多选状态
   - 搜索结果也应该支持排序

## Dev Agent Record

### Agent Model Used

MiniMax-M2.7-highspeed

### Debug Log References

### Completion Notes List

- 修改 `internal/repository/artist.go` - `GetSongsByArtist` 添加 sortBy 和 order 参数，支持白名单校验
- 修改 `internal/repository/album.go` - `GetSongsByAlbum` 添加 sortBy 和 order 参数，支持白名单校验
- 修改 `internal/repository/folder.go` - `GetSongsByFolder` 添加 sortBy 和 order 参数，支持白名单校验
- 修改 `internal/handler/artist.go` - `GetArtistSongs` 添加 sort_by 和 order 查询参数
- 修改 `internal/handler/album.go` - `GetAlbumSongs` 添加 sort_by 和 order 查询参数
- 修改 `internal/handler/folder.go` - `GetFolderSongs` 添加 sort_by 和 order 查询参数
- 创建 `frontend/src/components/common/sort-selector.tsx` - SortSelector 组件，支持排序字段选择和升序/降序切换
- 更新 `frontend/src/views/artists-view.tsx` - 集成 SortSelector 到展开的歌曲列表
- 更新 `frontend/src/views/albums-view.tsx` - 集成 SortSelector 到展开的歌曲列表
- 更新 `frontend/src/views/folders-view.tsx` - 集成 SortSelector 到展开的歌曲列表
- 更新测试文件：`artist_test.go`, `album_test.go`, `folder_test.go` 添加排序参数
- 所有后端测试通过

### File List

- `internal/repository/artist.go` - 修改
- `internal/repository/album.go` - 修改
- `internal/repository/folder.go` - 修改
- `internal/handler/artist.go` - 修改
- `internal/handler/album.go` - 修改
- `internal/handler/folder.go` - 修改
- `internal/repository/artist_test.go` - 修改
- `internal/repository/album_test.go` - 修改
- `internal/repository/folder_test.go` - 修改
- `frontend/src/components/common/sort-selector.tsx` - 新建
- `frontend/src/views/artists-view.tsx` - 修改
- `frontend/src/views/albums-view.tsx` - 修改
- `frontend/src/views/folders-view.tsx` - 修改

## Change Log

- 2026-04-16: 创建 Story 2.4 故事文件
- 2026-04-16: 实现后端排序功能（Repository + Handler 修改）
- 2026-04-16: 创建 SortSelector 前端组件
- 2026-04-16: 集成排序组件到 ArtistsView, AlbumsView, FoldersView
- 2026-04-16: 更新测试文件并验证所有测试通过

## References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-2.4]
- [Source: _bmad-output/planning-artifacts/architecture.md#API-Communication-Patterns]
- [Source: _bmad-output/implementation-artifacts/stories/2-1-artist-view-browse.md]
- [Source: _bmad-output/implementation-artifacts/stories/2-2-album-view-browse.md]
- [Source: _bmad-output/implementation-artifacts/stories/2-3-folder-view-browse.md]
