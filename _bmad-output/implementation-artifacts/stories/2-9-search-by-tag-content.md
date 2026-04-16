# Story 2.9: 按标签内容搜索

Status: done

## Story

As a 用户，
I want 按 ID3 标签内容搜索音乐，
so that 通过歌曲信息找到目标歌曲。

## Acceptance Criteria

1. **Given** 用户在搜索框输入标签关键词 **When** 用户提交搜索 **Then** 搜索匹配标题、艺术家、专辑的歌曲

2. **And** 支持多条件组合搜索

3. **And** 搜索结果高亮匹配文字

4. **And** 无结果时显示"未找到匹配的歌曲"

## Tasks / Subtasks

- [x] Task 1: 添加 Song Repository 标签搜索方法 (AC: 1-2)
  - [x] 在 `internal/repository/song.go` 添加 `SearchByTagContent` 方法
  - [x] 使用 GORM 的 LIKE 查询搜索 title、artist、album 字段
  - [x] 支持多条件组合搜索（关键词匹配这三个字段中的任意一个）
  - [x] 添加 SQL 注入防护

- [x] Task 2: 添加 Song Handler 标签搜索端点 (AC: 1-2)
  - [x] 在 `internal/handler/song.go` 添加 `SearchSongsByTag` 处理方法
  - [x] 实现 `GET /api/songs/search/by-tag?q=keyword` 端点（与 2.8 的 filename 搜索区分）
  - [x] 添加请求参数校验
  - [x] 注册路由到 `cmd/server/main.go`

- [x] Task 3: 更新 SearchBar 支持搜索类型切换 (AC: 1, 3-4)
  - [x] 修改 `frontend/src/components/common/search-bar.tsx`
  - [x] 添加搜索类型选择（文件名 / 标签内容）
  - [x] 调用不同的 API 端点

- [x] Task 4: 创建标签搜索结果视图 (AC: 1, 3-4)
  - [x] 创建或扩展 `frontend/src/views/search-results-view.tsx`
  - [x] 支持高亮 title、artist、album 中的匹配文字
  - [x] 无结果时显示"未找到匹配的歌曲"

- [x] Task 5: 集成搜索功能到 App (AC: 1, 3-4)
  - [x] 在 App 中添加搜索类型状态管理
  - [x] 根据搜索类型调用不同的 API

- [x] Task 6: 编写测试 (AC: 1-2)
  - [x] 后端: SongRepository SearchByTagContent 单元测试
  - [x] 后端: SongHandler SearchSongsByTag 单元测试

## Dev Notes

### 技术要求

**Song Repository - SearchByTagContent:**
- 使用 GORM 的 `OR` 条件查询 title、artist、album 字段
- 搜索语法: `WHERE title LIKE ? OR artist LIKE ? OR album LIKE ?`
- 多个关键词使用空格分隔，支持 AND 逻辑
- 防止 SQL 注入：keyword 仅作为 LIKE 参数，不拼接 SQL

**API 端点:**
```
GET /api/songs/search/by-tag?q={keyword}
Response: {
  "data": [
    { "id": 1, "title": "晴天", "artist": "周杰伦", "album": "叶惠美", ... },
    ...
  ]
}

Query Parameters:
- q (required): 搜索关键词，支持多条件空格分隔

Error Response:
- 400: 关键词为空或无效
- 500: 服务器内部错误
```

**搜索结果高亮:**
- 前端高亮 title、artist、album 中的关键词
- 使用 `<mark>` 标签包裹匹配文字
- 样式：背景色 #FEF08A (黄色高亮)

### Project Structure Notes

**现有结构:**
```
nas-manager/
├── cmd/server/main.go              # 需添加 /api/songs/search/by-tag 路由
├── internal/
│   ├── handler/
│   │   └── song.go                # 需添加 SearchSongsByTag 方法
│   ├── repository/
│   │   └── song.go                # 需添加 SearchByTagContent 方法
│   └── model/
│       └── song.go                # 已存在 Song 模型
├── frontend/
│   ├── src/
│   │   ├── views/
│   │   │   ├── artists-view.tsx   # Story 2.1 已创建
│   │   │   ├── albums-view.tsx    # Story 2.2 已创建
│   │   │   ├── folders-view.tsx   # Story 2.3 已创建
│   │   │   └── search-results-view.tsx  # Story 2.8 已创建
│   │   ├── components/
│   │   │   └── common/
│   │   │       └── search-bar.tsx  # Story 2.8 已创建
│   │   ├── app.tsx               # 已集成搜索功能
│   │   └── types/
│   │       └── song.ts           # Story 2.1 已创建
```

### Architecture Compliance

**遵循规范:**
1. **命名规范:** API 使用 snake_case，Go 代码使用 PascalCase/camelCase
2. **响应格式:** 使用 `pkg/response/response.go` 的统一响应格式
3. **错误处理:** 使用 `pkg/response` 的错误响应格式
4. **分层架构:** Handler → Repository (无 Service 层对于只读查询)
5. **前端状态:** 使用 Preact hooks 管理搜索状态

**UX 规范:**
- 搜索框宽度：桌面端 300px，平板 200px，手机 100%
- 搜索图标：放大镜 SVG
- 高亮颜色：#FEF08A 黄色背景
- 文字颜色：#1E293B (深色)
- 触摸目标：最小 44x44px

### 与其他 Story 的差异

| 方面 | Story 2.8 (文件名搜索) | Story 2.9 (标签搜索) |
|------|------------------------|----------------------|
| 搜索字段 | file_path | title, artist, album |
| API 端点 | /api/songs/search | /api/songs/search/by-tag |
| 高亮内容 | 文件名 | 标题、艺术家、专辑 |
| Repository 方法 | SearchByFileName | SearchByTagContent |

### Previous Story Intelligence (Story 2.8)

**经验总结:**
- Story 2.8 实现了文件名搜索功能，使用 `/api/songs/search?q=keyword` 端点
- SearchResultsView 复用 SearchBar 组件，支持高亮和无结果提示
- Repository SearchByFileName 使用 `db.Where("file_path LIKE ?", "%"+keyword+"%")` 模式
- Handler SearchSongs 处理分页参数和错误响应

**复用 Story 2.8 的实现:**
- SearchBar 组件结构（需扩展支持搜索类型切换）
- SearchResultsView 组件（需扩展高亮字段）
- 高亮工具函数 `highlightMatch` 和 `escapeRegExp`
- SongTableRow 的 `highlightedText` prop

**代码模式参考:**
```go
// SearchByFileName 在 SongRepository 中 (Story 2.8)
func (r *SongRepository) SearchByFileName(keyword string) ([]model.Song, error) {
    var songs []model.Song
    err := r.db.Where("file_path LIKE ?", "%"+keyword+"%").Find(&songs).Error
    if err != nil {
        return nil, err
    }
    return songs, nil
}
```

## Dev Agent Record

### Agent Model Used

MiniMax-M2.7-highspeed

### Debug Log References

- Repository tests: SearchByTagContent, SearchByTagContentMulti, SearchByTagContent_NoResults pass
- Handler tests: SearchSongsByTag_Success, SearchSongsByTag_MultiKeyword, SearchSongsByTag_NoResults, SearchSongsByTag_MissingQuery pass
- Build passes without errors
- All tests pass

### Completion Notes List

- 创建 `internal/repository/song.go` - 添加 `SearchByTagContent` 方法（单关键词 OR 查询）和 `SearchByTagContentMulti` 方法（多关键词 AND 查询）
- 创建 `internal/handler/song.go` - 添加 `SearchSongsByTag` 处理方法，实现 `GET /api/songs/search/by-tag?q=keyword` 端点
- 更新 `cmd/server/main.go` - 注册 `/api/songs/search/by-tag` 路由
- 更新 `frontend/src/components/common/search-bar.tsx` - 添加搜索类型选择（标签/文件名），调用不同的 API 端点
- 更新 `frontend/src/views/search-results-view.tsx` - 根据搜索类型高亮不同字段（文件名 vs 标题/艺术家/专辑）
- 更新 `frontend/src/app.tsx` - 添加搜索类型状态管理
- 创建 `internal/repository/song_test.go` - 添加 3 个 SearchByTagContent 单元测试
- 创建 `internal/handler/song_test.go` - 添加 4 个 SearchSongsByTag 单元测试
- 所有测试通过

### File List

1. `internal/repository/song.go` - 添加 SearchByTagContent 和 SearchByTagContentMulti 方法
2. `internal/handler/song.go` - 添加 SearchSongsByTag 处理方法
3. `cmd/server/main.go` - 注册 /api/songs/search/by-tag 路由
4. `frontend/src/components/common/search-bar.tsx` - 添加搜索类型选择
5. `frontend/src/views/search-results-view.tsx` - 支持标签搜索高亮
6. `frontend/src/app.tsx` - 搜索类型状态管理
7. `internal/repository/song_test.go` - 添加 SearchByTagContent 测试
8. `internal/handler/song_test.go` - 添加 SearchSongsByTag 测试

## Change Log

- 2026-04-16: 创建 Story 2.9 故事文件
- 2026-04-16: 实现后端 SearchByTagContent 和 SearchSongsByTag
- 2026-04-16: 实现前端 SearchBar 搜索类型切换和 SearchResultsView 标签高亮
- 2026-04-16: 集成搜索类型管理到 App
- 2026-04-16: 编写后端单元测试，所有测试通过
- 2026-04-16: Story 完成，标记为 review
