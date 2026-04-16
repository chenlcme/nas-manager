# Story 2.3: 文件夹视图浏览

Status: done

## Story

As a 用户，
I want 按文件夹结构浏览音乐库，
So that 按照文件系统组织查看音乐。

## Acceptance Criteria

1. **Given** 用户切换到"文件夹" Tab **When** 系统加载文件夹视图 **Then** 按音乐目录的文件夹结构显示

2. **And** 每层显示文件夹路径和包含的歌曲数量

3. **And** 点击文件夹展开显示其中的歌曲

4. **And** 显示文件在磁盘上的实际路径

## Tasks / Subtasks

- [x] Task 1: 创建 Folder Repository (AC: 1-4)
  - [x] 创建 `internal/repository/folder.go`
  - [x] 实现按文件夹路径分组并统计歌曲数量的查询方法
  - [x] 实现获取特定文件夹内歌曲列表的方法
  - [x] 获取文件夹的相对路径（相对于音乐目录）

- [x] Task 2: 创建 Folder Handler (AC: 1-4)
  - [x] 创建 `internal/handler/folder.go`
  - [x] 实现 `GET /api/folders` 端点，返回文件夹列表（含路径和歌曲数量）
  - [x] 实现 `GET /api/folders/:id/songs` 端点，返回特定文件夹的歌曲列表
  - [x] 注册路由到 `cmd/server/main.go`

- [x] Task 3: 创建前端 Folders View 组件 (AC: 1-4)
  - [x] 创建 `frontend/src/views/folders-view.tsx`
  - [x] 实现文件夹列表展示组件（与 artists-view.tsx/albums-view.tsx 结构类似）
  - [x] 实现点击展开查看文件夹内歌曲功能
  - [x] 显示文件的实际磁盘路径

- [x] Task 4: 复用 SongTableRow 组件 (AC: 1, 3)
  - [x] 复用 `frontend/src/components/song/song-table-row.tsx`（Story 2.1 已创建）

- [x] Task 5: 复用 SelectionBar 组件 (AC: 1, 3)
  - [x] 复用 `frontend/src/components/common/selection-bar.tsx`（Story 2.1 已创建）

- [x] Task 6: 复用 Tab 导航组件 (AC: 1)
  - [x] 复用 `frontend/src/components/common/tab-nav.tsx`（Story 2.1 已创建）

- [x] Task 7: 集成到 App 结构 (AC: 1-4)
  - [x] 更新 `frontend/src/app.tsx` 添加文件夹视图路由/状态
  - [x] 确保 Tab 导航正确切换文件夹视图

- [x] Task 8: 编写测试 (AC: 1-4)
  - [x] 后端: Folder repository 单元测试
  - [x] 后端: Folder handler 单元测试

## Dev Notes

### 技术要求

**Folder Repository:**
- 使用 GORM 的 SUBSTR 和 INSTR 函数提取父目录路径
- 按父目录分组，统计每个文件夹的歌曲数量
- 需要从 `file_path` 字段提取目录路径
- 支持按文件夹路径排序
- 过滤空路径或根目录

**API 端点:**
```
GET /api/folders
Response: {
  "data": [
    { "id": 1, "path": "/music/rock", "songCount": 15 },
    { "id": 2, "path": "/music/pop", "songCount": 23 }
  ]
}

GET /api/folders/:id/songs
Response: {
  "data": [
    { "id": 1, "title": "歌曲1", "file_path": "/music/rock/song1.mp3", "duration": 267, ... },
    ...
  ]
}
```

**Folders View 前端:**
- 紧凑表格视图（复用 SongTableRow）
- 点击文件夹行展开歌曲列表
- 文件夹路径作为主要显示内容
- 歌曲列表中显示文件的完整磁盘路径

### Project Structure Notes

**现有结构（Story 2.1, 2.2 已创建）:**
```
nas-manager/
├── cmd/server/main.go              # 需添加 Folder 路由
├── internal/
│   ├── handler/
│   │   ├── artist.go              # Story 2.1 已创建
│   │   ├── album.go               # Story 2.2 已创建
│   │   └── folder.go              # TODO: 创建
│   ├── service/                   # (只读查询无 Service 层)
│   ├── repository/
│   │   ├── artist.go              # Story 2.1 已创建
│   │   ├── album.go               # Story 2.2 已创建
│   │   └── folder.go              # TODO: 创建
│   └── model/
│       ├── album.go               # Story 1 已创建
│       └── song.go                # Story 1 已创建
├── frontend/
│   ├── src/
│   │   ├── views/
│   │   │   ├── artists-view.tsx  # Story 2.1 已创建
│   │   │   ├── albums-view.tsx   # Story 2.2 已创建
│   │   │   └── folders-view.tsx  # TODO: 创建
│   │   ├── components/           # song/, common/ 已创建组件可复用
│   │   ├── app.tsx               # 需更新添加文件夹视图
│   │   └── types/
│   │       └── song.ts           # Story 2.1 已创建
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
- 复用 Story 2.1/2.2 的 UX 组件模式

### 与 Story 2.1/2.2 的差异

| 方面 | Story 2.1 (Artist) | Story 2.2 (Album) | Story 2.3 (Folder) |
|------|-------------------|-------------------|-------------------|
| 分组字段 | Artist.Name | Album.Name + Album.Artist | SUBSTR(FilePath, 0, LENGTH(FilePath) - LENGTH(FileName)) |
| 显示内容 | 艺术家名 + 歌曲数量 | 专辑名 + 艺术家名 + 歌曲数量 | 文件夹路径 + 歌曲数量 |
| Repository | ArtistRepository | AlbumRepository | FolderRepository |
| Handler | ArtistHandler | AlbumHandler | FolderHandler |
| View | artists-view.tsx | albums-view.tsx | folders-view.tsx |
| 路径获取 | N/A | N/A | 需要从 file_path 提取父目录 |

### Previous Story Intelligence (Stories 2.1, 2.2)

**经验总结:**
- ArtistRepository/AlbumRepository 使用 `db.Model(&Song{}).Select(...).Group(...)` 模式
- Handler 直接调用 Repository，无需 Service 层（只读查询）
- 前端使用 Preact functional 组件 + hooks
- 统一响应格式: `{"data": [...]}`
- 错误响应格式: `{"error": {"code": "...", "message": "..."}}`

**复用的组件:**
- `song-table-row.tsx` - 高密度歌曲表格行
- `selection-bar.tsx` - 选择操作栏
- `tab-nav.tsx` - Tab 导航
- `selection-context.tsx` - 选择状态管理

**代码模式参考 (AlbumRepository -> FolderRepository 适配):**
```go
// AlbumRepository 使用 SUBSTR(file_path, 1, LENGTH(file_path) - LENGTH(REPLACE(file_path, '/', '')) 来获取目录
// SQLite 中提取文件路径的父目录：
// SUBSTR(file_path, 0, LENGTH(file_path) - LENGTH(filename))
// 或者使用 INSTR 找到最后一个斜杠的位置

func (r *FolderRepository) GetAllFoldersWithSongCount(sortOrder string) ([]FolderInfo, error) {
    // 提取文件路径的目录部分
    // 使用 SUBSTR 和 INSTR 函数
    // 示例: SUBSTR(file_path, 1, LENGTH(file_path) - LENGTH(REPLACE(file_path, '/', '')))
}
```

**Review 问题修复记录（来自 Story 2.2）:**
- [x] React Fragment 缺少 key prop — 使用 Fragment key
- [x] Stale Closure 导致 useEffect 过期状态 — 添加 loading 检查
- [x] 动态分配 ID 竞态条件 — 使用稳定的方式获取文件夹 ID
- [x] 排序参数缺少校验 — 添加 sortOrder 参数校验
- [x] 空值状态未妥善清除 — catch 中添加清理
- [x] 并发状态竞争风险 — 添加 loading 检查防止竞态

### 特殊注意事项

1. **路径处理:**
   - Windows 路径分隔符为 `\`，Linux/Mac 为 `/`
   - 需要兼容处理音乐目录跨平台场景
   - 相对路径 vs 绝对路径的处理

2. **文件夹层级:**
   - 当前实现仅支持一级目录分组
   - 不需要支持嵌套展开（仅展开当前层级）

3. **音乐目录配置:**
   - 需要从 settings 表获取 music_dir 配置
   - 文件夹路径应显示为相对于 music_dir 的路径

## Review Findings

### Senior Developer Review (AI)

Date: 2026-04-16

**Review Outcome:** Changes Requested

**Action Items:**
- [x] [Review][Patch] AC4 Violation: SongTableRow 不显示 filePath [folders-view.tsx] — 已修复：添加 showPath prop
- [x] [Review][Patch] LIKE 模式未转义特殊字符 [_%.] [folder.go] — 已修复：添加 ESCAPE 子句
- [x] [Review][Patch] 根文件路径返回 "." 而非 "/" [folder.go] — 已修复：path.Dir 返回 "." 时返回 "/"
- [x] [Review][Patch] 尾部斜杠导致父目录错误 [folder.go] — 已修复：trim 尾部斜杠
- [x] [Review][Patch] Windows 路径分隔符不兼容 [folder.go] — 已修复：统一使用正斜杠处理
- [x] [Review][Patch] 双斜杠 // 未规范化 [folder.go] — 已修复：替换双斜杠为单斜杠
- [x] [Review][Defer] GetFolderPathByID 重复全表扫描 [folder.go:98] — deferred, consistent with Artist/Album pattern
- [x] [Review][Defer] 全表扫描加载所有歌曲到内存 [folder.go:48] — deferred, known trade-off
- [x] [Review][Defer] 无分页支持 [folder.go:26,47] — deferred, local NAS scope
- [x] [Review][Defer] 动态 ID 排序变化时不稳定 [folder.go:86] — deferred, consistent with existing pattern
- [x] [Review][Defer] . 文件夹名处理 [folder.go:33] — deferred, edge case
- [x] [Review][Defer] 非 ASCII Unicode 大小写排序 [folder.go:69,73] — deferred, Go known limitation

## Dev Agent Record

### Agent Model Used

MiniMax-M2.7-highspeed

### Debug Log References

- SQLite REVERSE function not available - switched to Go-level path extraction
- Empty slice vs nil issue when serializing to JSON - initialized results as empty slice

### Completion Notes List

- 创建 `internal/repository/folder.go` - FolderRepository，包含 GetAllFoldersWithSongCount、GetSongsByFolder、GetFolderPathByID 方法
- 创建 `internal/repository/folder_test.go` - FolderRepository 单元测试（7个测试用例，全部通过）
- 创建 `internal/handler/folder.go` - FolderHandler，包含 GET /api/folders 和 GET /api/folders/:id/songs 端点
- 创建 `internal/handler/folder_test.go` - FolderHandler 单元测试（6个测试用例，全部通过）
- 更新 `cmd/server/main.go` - 注册 Folder 路由，初始化 FolderRepository 和 FolderHandler
- 更新 `frontend/src/types/song.ts` - 添加 FolderWithCount TypeScript 类型
- 创建 `frontend/src/views/folders-view.tsx` - 文件夹视图组件，与 AlbumsView 结构类似
- 更新 `frontend/src/app.tsx` - 导入并使用 FoldersView 替代占位符
- 所有后端测试通过
- 技术说明：由于 SQLite 不支持 REVERSE 函数，文件夹路径提取使用 Go 的 path.Dir 函数实现

## Change Log

- 2026-04-16: 创建 Story 2.3 故事文件
- 2026-04-16: 实现 Folder Repository、Handler、路由注册
- 2026-04-16: 实现前端 FoldersView 组件并集成到 App
- 2026-04-16: 编写后端单元测试，所有测试通过

## References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-2.3]
- [Source: _bmad-output/planning-artifacts/architecture.md#API-Communication-Patterns]
- [Source: _bmad-output/planning-artifacts/architecture.md#Data-Architecture]
- [Source: _bmad-output/planning-artifacts/ux-design-specification.md#UX-DR1-SongTableRow]
- [Source: _bmad-output/planning-artifacts/ux-design-specification.md#UX-DR4-SelectionBar]
- [Source: _bmad-output/planning-artifacts/ux-design-specification.md#UX-DR6-Tab-Navigation]
- [Source: _bmad-output/implementation-artifacts/stories/2-1-artist-view-browse.md]
- [Source: _bmad-output/implementation-artifacts/stories/2-2-album-view-browse.md]
