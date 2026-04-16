# Story 2.8: 按文件名搜索

Status: review

## Story

As a 用户，
I want 按文件名搜索音乐，
So that 快速找到特定文件。

## Acceptance Criteria

1. **Given** 用户在搜索框输入文件名关键词 **When** 用户提交搜索（回车或点击搜索按钮） **Then** 搜索匹配文件名的歌曲

2. **And** 支持中文和英文文件名

3. **And** 搜索结果高亮匹配文字

4. **And** 无结果时显示"未找到匹配的歌曲"

## Tasks / Subtasks

- [x] Task 1: 添加 Song Repository 搜索方法 (AC: 1-2)
  - [x] 在 `internal/repository/song.go` 添加 `SearchByFileName` 方法
  - [x] 使用 GORM 的 LIKE 查询支持中文和英文文件名
  - [x] 添加 SQL 注入防护

- [x] Task 2: 添加 Song Handler 搜索端点 (AC: 1-2)
  - [x] 在 `internal/handler/song.go` 添加 `SearchSongs` 处理方法
  - [x] 实现 `GET /api/songs/search?q=keyword` 端点
  - [x] 添加请求参数校验
  - [x] 注册路由到 `cmd/server/main.go`

- [x] Task 3: 创建 SearchBar 组件 (AC: 1, 3-4)
  - [x] 创建 `frontend/src/components/common/search-bar.tsx`
  - [x] 实现搜索输入框
  - [x] 实现回车和点击按钮提交搜索
  - [x] 显示加载状态

- [x] Task 4: 创建搜索结果视图 (AC: 1, 3-4)
  - [x] 创建 `frontend/src/views/search-results-view.tsx`
  - [x] 展示搜索结果列表
  - [x] 高亮匹配文字
  - [x] 无结果时显示"未找到匹配的歌曲"

- [x] Task 5: 集成搜索功能到 App (AC: 1, 3-4)
  - [x] 在 TabNav 旁边添加搜索输入框
  - [x] 实现搜索状态管理
  - [x] 点击搜索结果显示搜索结果视图

- [x] Task 6: 编写测试 (AC: 1-2)
  - [x] 后端: SongRepository SearchByFileName 单元测试
  - [x] 后端: SongHandler SearchSongs 单元测试

## Dev Notes

### 技术要求

**Song Repository - SearchByFileName:**
- 使用 GORM 的 `Where("file_path LIKE ?", "%"+keyword+"%")` 查询
- file_path 字段存储完整文件路径，搜索时匹配文件名部分
- 支持中文和英文：SQLite LIKE 本身支持 Unicode，无需特殊处理
- 防止 SQL 注入：keyword 仅作为 LIKE 参数，不拼接 SQL

**API 端点:**
```
GET /api/songs/search?q={keyword}
Response: {
  "data": [
    { "id": 1, "file_path": "/music/rock/晴天.mp3", "title": "晴天", ... },
    ...
  ]
}

Query Parameters:
- q (required): 搜索关键词，最小 1 字符

Error Response:
- 400: 关键词为空或无效
- 500: 服务器内部错误
```

**搜索结果高亮:**
- 前端使用正则匹配文件名中的关键词
- 使用 `<mark>` 标签包裹匹配文字
- 样式：背景色 #FEF08A (黄色高亮)

**搜索输入框位置:**
- 位于 TabNav 右侧或下方
- 搜索图标 + 输入框 + 搜索按钮

### Project Structure Notes

**现有结构:**
```
nas-manager/
├── cmd/server/main.go              # 需添加 /api/songs/search 路由
├── internal/
│   ├── handler/
│   │   └── song.go                # 需添加 SearchSongs 方法
│   ├── repository/
│   │   └── song.go                # 需添加 SearchByFileName 方法
│   └── model/
│       └── song.go                # 已存在 Song 模型
├── frontend/
│   ├── src/
│   │   ├── views/
│   │   │   ├── artists-view.tsx   # Story 2.1 已创建
│   │   │   ├── albums-view.tsx   # Story 2.2 已创建
│   │   │   ├── folders-view.tsx  # Story 2.3 已创建
│   │   │   └── search-results-view.tsx  # TODO: 创建
│   │   ├── components/
│   │   │   └── common/
│   │   │       └── search-bar.tsx  # TODO: 创建
│   │   ├── app.tsx               # 需集成搜索功能
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

| 方面 | Story 2.1-2.3 (浏览) | Story 2.8 (搜索) |
|------|----------------------|------------------|
| 数据获取 | 按分组条件查询 | 模糊匹配文件名 |
| 展示方式 | 分组列表 | 扁平列表 |
| API 参数 | 路径参数 | 查询参数 q |
| Repository | Group + Count | LIKE 查询 |

### Previous Story Intelligence (Stories 2.1-2.7)

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
// SearchByFileName 在 SongRepository 中
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

- Repository SearchByFileName tests: all 4 tests pass
- Handler SearchSongs tests: all 6 tests pass
- Build passes without errors

### Completion Notes List

- 创建 `internal/repository/song.go` - 添加 `SearchByFileName` 方法，使用 GORM LIKE 查询支持模糊匹配
- 创建 `internal/handler/song.go` - 添加 `SearchSongs` 处理方法，实现 `GET /api/songs/search?q=keyword` 端点
- 更新 `cmd/server/main.go` - 注册 `/api/songs/search` 路由
- 创建 `frontend/src/components/common/search-bar.tsx` - 搜索输入框组件，支持回车和按钮提交
- 创建 `frontend/src/views/search-results-view.tsx` - 搜索结果视图，支持高亮显示和无结果提示
- 更新 `frontend/src/app.tsx` - 集成搜索功能，在 TabNav 旁边添加搜索框
- 更新 `frontend/src/components/song/song-table-row.tsx` - 添加 `highlightedText` prop 支持高亮
- 创建 `internal/repository/song_test.go` - 添加 4 个 SearchByFileName 单元测试
- 创建 `internal/handler/song_test.go` - 添加 6 个 SearchSongs 单元测试
- 所有搜索相关测试通过

### File List

1. `internal/repository/song.go` - 添加 SearchByFileName 方法
2. `internal/handler/song.go` - 添加 SearchSongs 处理方法
3. `cmd/server/main.go` - 注册搜索路由
4. `frontend/src/components/common/search-bar.tsx` - 新建搜索栏组件
5. `frontend/src/views/search-results-view.tsx` - 新建搜索结果视图
6. `frontend/src/app.tsx` - 集成搜索功能
7. `frontend/src/components/song/song-table-row.tsx` - 添加高亮文本支持
8. `internal/repository/song_test.go` - 添加 SearchByFileName 测试
9. `internal/handler/song_test.go` - 添加 SearchSongs 测试

## Change Log

- 2026-04-16: 创建 Story 2.8 故事文件
- 2026-04-16: 实现后端 SearchByFileName 和 SearchSongs
- 2026-04-16: 实现前端 SearchBar 和 SearchResultsView 组件
- 2026-04-16: 集成搜索功能到 App
- 2026-04-16: 编写后端单元测试，所有测试通过
- 2026-04-16: Story 完成，标记为 review
