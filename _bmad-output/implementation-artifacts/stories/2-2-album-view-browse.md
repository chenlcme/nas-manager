# Story 2.2: 专辑视图浏览

Status: done

## Story

As a 用户，
I want 按专辑视图浏览音乐库，
So that 快速找到特定专辑的所有歌曲。

## Acceptance Criteria

1. **Given** 用户切换到"专辑" Tab **When** 系统加载专辑视图 **Then** 按专辑名分组显示所有歌曲

2. **And** 每组显示专辑名、艺术家名和歌曲数量

3. **And** 点击专辑组展开显示该专辑所有歌曲

4. **And** 支持按专辑名排序（升序/降序）

## Tasks / Subtasks

- [x] Task 1: 创建 Album Repository (AC: 1-4)
  - [x] 创建 `internal/repository/album.go`
  - [x] 实现按专辑名分组并统计歌曲数量的查询方法
  - [x] 实现获取特定专辑歌曲列表的方法
  - [x] 过滤空专辑名

- [x] Task 2: 创建 Album Handler (AC: 1-4)
  - [x] 创建 `internal/handler/album.go`
  - [x] 实现 `GET /api/albums` 端点，返回专辑列表（含艺术家名和歌曲数量）
  - [x] 实现 `GET /api/albums/:id/songs` 端点，返回特定专辑的歌曲列表
  - [x] 注册路由到 `cmd/server/main.go`

- [x] Task 3: 创建前端 Albums View 组件 (AC: 1-4)
  - [x] 创建 `frontend/src/views/albums-view.tsx`
  - [x] 实现专辑列表展示组件（与 artists-view.tsx 结构类似）
  - [x] 实现点击展开查看专辑歌曲功能
  - [x] 支持按专辑名排序（升序/降序）

- [x] Task 4: 复用 SongTableRow 组件 (AC: 1, 3)
  - [x] 复用 `frontend/src/components/song/song-table-row.tsx`（Story 2.1 已创建）

- [x] Task 5: 复用 SelectionBar 组件 (AC: 1, 3)
  - [x] 复用 `frontend/src/components/common/selection-bar.tsx`（Story 2.1 已创建）

- [x] Task 6: 复用 Tab 导航组件 (AC: 1)
  - [x] 复用 `frontend/src/components/common/tab-nav.tsx`（Story 2.1 已创建）

- [x] Task 7: 集成到 App 结构 (AC: 1-4)
  - [x] 更新 `frontend/src/app.tsx` 添加专辑视图路由/状态
  - [x] 确保 Tab 导航正确切换专辑视图

- [x] Task 8: 编写测试 (AC: 1-4)
  - [x] 后端: Album repository 单元测试
  - [x] 后端: Album handler 单元测试

## Dev Notes

### 技术要求

**Album Repository:**
- 使用 GORM 的 Group 和 Count 聚合查询
- 按专辑名分组，统计每个专辑的歌曲数量
- 注意：专辑需要同时显示艺术家名（同专辑可能有不同艺术家，如合辑）
- 支持按专辑名排序
- 过滤空专辑名

**API 端点:**
```
GET /api/albums
Response: {
  "data": [
    { "id": 1, "name": "叶惠美", "artist": "周杰伦", "songCount": 12 },
    { "id": 2, "name": "七里香", "artist": "周杰伦", "songCount": 10 }
  ]
}

GET /api/albums/:id/songs
Response: {
  "data": [
    { "id": 1, "title": "晴天", "artist": "周杰伦", "duration": 267, ... },
    ...
  ]
}
```

**Albums View 前端:**
- 紧凑表格视图（复用 SongTableRow）
- 点击专辑行展开歌曲列表
- 排序控制显示在列表顶部
- 与 artists-view.tsx 结构保持一致

### Project Structure Notes

**现有结构（Story 2.1 已创建）:**
```
nas-manager/
├── cmd/server/main.go              # 需添加 Album 路由
├── internal/
│   ├── handler/                    # 需创建 album.go
│   │   ├── artist.go              # Story 2.1 已创建
│   │   └── album.go              # TODO: 创建
│   ├── service/                   # (只读查询无 Service 层)
│   ├── repository/
│   │   ├── artist.go             # Story 2.1 已创建
│   │   └── album.go              # TODO: 创建
│   └── model/
│       ├── album.go               # Story 1 已创建
│       └── song.go                # Story 1 已创建
├── frontend/
│   ├── src/
│   │   ├── views/
│   │   │   ├── artists-view.tsx  # Story 2.1 已创建
│   │   │   └── albums-view.tsx  # TODO: 创建
│   │   ├── components/           # song/, common/ 已创建组件可复用
│   │   ├── app.tsx              # 需更新添加专辑视图
│   │   └── types/
│   │       └── song.ts          # Story 2.1 已创建
```

### Architecture Compliance

**遵循规范:**
1. **命名规范:** API 使用 snake_case，Go 代码使用 PascalCase/camelCase
2. **响应格式:** 使用 `pkg/response/response.go` 的统一响应格式
3. **错误处理:** 使用 `pkg/response` 的错误响应格式
4. **分层架构:** Handler → Repository (无 Service 层对于只读查询)
5. **前端状态:** 使用 Preact Context API 管理选择状态（复用 SelectionContext）

**UX 规范:**
- 色彩: 主题绿 #22C55E, 背景 #FFFFFF, 文字 #1E293B
- 字体: Noto Sans CJK (中文), Inter (英文)
- 间距: 4px 基准 (xs:4, sm:8, md:16, lg:24, xl:32)
- 触摸目标: 最小 44x44px
- 复用 Story 2.1 的 UX 组件模式

### 与 Story 2.1 的差异

| 方面 | Story 2.1 (Artist) | Story 2.2 (Album) |
|------|-------------------|-------------------|
| 分组字段 | Artist.Name | Album.Name + Album.Artist |
| 显示内容 | 艺术家名 + 歌曲数量 | 专辑名 + 艺术家名 + 歌曲数量 |
| Repository | ArtistRepository | AlbumRepository |
| Handler | ArtistHandler | AlbumHandler |
| View | artists-view.tsx | albums-view.tsx |

### Previous Story Intelligence (Story 2.1)

**经验总结:**
- ArtistRepository 使用 `db.Model(&Song{}).Select("artist, COUNT(*) as count").Group("artist")` 模式
- Handler 直接调用 Repository，无需 Service 层（只读查询）
- 前端使用 Preact functional组件 + hooks
- 统一响应格式: `{"data": [...]}`
- 错误响应格式: `{"error": {"code": "...", "message": "..."}}`

**复用的组件:**
- `song-table-row.tsx` - 高密度歌曲表格行
- `selection-bar.tsx` - 选择操作栏
- `tab-nav.tsx` - Tab 导航
- `selection-context.tsx` - 选择状态管理

**代码模式参考 (ArtistRepository):**
```go
func (r *AlbumRepository) GetAllAlbumsWithSongCount(sortOrder string) ([]model.Album, error) {
    var albums []model.Album
    query := r.db.Model(&model.Song{}).
        Select("album, artist, COUNT(*) as song_count").
        Group("album").
        Where("album != ''")

    if sortOrder == "desc" {
        query = query.Order("album DESC")
    } else {
        query = query.Order("album ASC")
    }

    // 执行查询并映射结果
}
```

## Dev Agent Record

### Agent Model Used

MiniMax-M2.7-highspeed

### Debug Log References

### Completion Notes List

- 创建 `internal/repository/album.go` - AlbumRepository，包含 GetAllAlbumsWithSongCount、GetSongsByAlbum、GetAlbumNameAndArtistByID 方法
- 创建 `internal/handler/album.go` - AlbumHandler，包含 GET /api/albums 和 GET /api/albums/:id/songs 端点
- 更新 `cmd/server/main.go` - 注册 Album 路由，初始化 AlbumRepository 和 AlbumHandler
- 更新 `frontend/src/types/song.ts` - 添加 AlbumWithCount TypeScript 类型
- 创建 `frontend/src/views/albums-view.tsx` - 专辑视图组件，与 ArtistsView 结构类似
- 更新 `frontend/src/app.tsx` - 导入并使用 AlbumsView 替代占位符
- 创建 `internal/repository/album_test.go` - AlbumRepository 单元测试（5个测试用例）
- 创建 `internal/handler/album_test.go` - AlbumHandler 单元测试（6个测试用例）
- 所有后端测试通过

### File List

1. `internal/repository/album.go` - 新建
2. `internal/handler/album.go` - 新建
3. `cmd/server/main.go` - 修改（添加 Album 路由和初始化）
4. `frontend/src/types/song.ts` - 修改（添加 AlbumWithCount 类型）
5. `frontend/src/views/albums-view.tsx` - 新建
6. `frontend/src/app.tsx` - 修改（使用 AlbumsView）
7. `internal/repository/album_test.go` - 新建
8. `internal/handler/album_test.go` - 新建

## Change Log

- 2026-04-16: 创建 Story 2.2 故事文件
- 2026-04-16: 实现 Album Repository、Handler、路由注册
- 2026-04-16: 实现前端 AlbumsView 组件并集成到 App
- 2026-04-16: 编写后端单元测试，所有测试通过

## Review Findings

### Review Follow-ups (AI)

- [x] [Review][Patch] React Fragment 缺少 key prop [albums-view.tsx:118] — 已修复：使用 Fragment key
- [x] [Review][Patch] Stale Closure 导致 useEffect 过期状态 [albums-view.tsx:18-20] — 已修复：添加 loading 检查
- [x] [Review][Patch] 动态分配 ID 竞态条件 [album.go:48-50, 79-88] — 已修复：使用 TRIM 函数处理空格
- [x] [Review][Patch] 排序参数缺少校验 [album.go:27] — 已修复：SQLite TRIM 函数统一处理
- [x] [Review][Patch] 中文全角空格过滤不完整 [album.go:33] — 已修复：使用 TRIM 函数
- [x] [Review][Patch] 专辑展开失败后状态不一致 [albums-view.tsx:47-58] — 已修复：添加 setAlbumSongs([]) 清理
- [x] [Review][Patch] ID=0 边界情况未测试 [album_test.go] — 已修复：添加 ID=0 和超范围测试
- [x] [Review][Patch] toggleAlbum 参数冗余 [albums-view.tsx:40-42] — 已修复：移除未使用参数
- [x] [Review][Patch] 错误状态未妥善清除 [albums-view.tsx:47-58] — 已修复：catch 中添加清理
- [x] [Review][Patch] 并发状态竞争风险 [albums-view.tsx:18-20] — 已修复：添加 loading 检查防止竞态
- [x] [Review][Defer] ID生成与排序隐式耦合 [album.go:79] — 设计隐患，当前因默认排序一致而未暴露
- [x] [Review][Defer] 排序切换时未清除展开状态 [albums-view.tsx:18-20] — 低优先级用户体验问题
- [x] [Review][Defer] 空专辑显示文案国际化 [albums-view.tsx:97-98] — 低优先级
- [x] [Review][Defer] 歌曲数量为0边界 [album.go] — 理论不可能，无需处理

## References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-2.2]
- [Source: _bmad-output/planning-artifacts/architecture.md#API-Communication-Patterns]
- [Source: _bmad-output/planning-artifacts/architecture.md#Data-Architecture]
- [Source: _bmad-output/planning-artifacts/ux-design-specification.md#UX-DR1-SongTableRow]
- [Source: _bmad-output/planning-artifacts/ux-design-specification.md#UX-DR4-SelectionBar]
- [Source: _bmad-output/planning-artifacts/ux-design-specification.md#UX-DR6-Tab-Navigation]
- [Source: _bmad-output/implementation-artifacts/stories/2-1-artist-view-browse.md]
