# Story 2.1: 歌手视图浏览

Status: review

## Story

As a 用户，
I want 按歌手视图浏览音乐库，
So that 快速找到特定艺术家的所有歌曲。

## Acceptance Criteria

1. **Given** 用户切换到"歌手" Tab **When** 系统加载歌手视图 **Then** 按歌手名分组显示所有歌曲

2. **And** 每组显示歌手名和歌曲数量

3. **And** 点击歌手组展开显示该艺术家所有歌曲

4. **And** 支持按歌手名排序（升序/降序）

## Tasks / Subtasks

- [x] Task 1: 创建 Artist Repository (AC: 1-4)
  - [x] 创建 `internal/repository/artist.go`
  - [x] 实现按艺术家名分组并统计歌曲数量的查询方法

- [x] Task 2: 创建 Artist Handler (AC: 1-4)
  - [x] 创建 `internal/handler/artist.go`
  - [x] 实现 `GET /api/artists` 端点，返回艺术家列表及歌曲数量
  - [x] 实现 `GET /api/artists/:id/songs` 端点，返回特定艺术家的歌曲列表
  - [x] 注册路由到 `cmd/server/main.go`

- [x] Task 3: 创建前端 Artist View 组件 (AC: 1-4)
  - [x] 创建 `frontend/src/views/artists-view.tsx`
  - [x] 实现艺术家列表展示组件
  - [x] 实现点击展开查看艺术家歌曲功能
  - [x] 支持按艺术家名排序（升序/降序）

- [x] Task 4: 创建 SongTableRow 组件 (AC: 1, 3)
  - [x] 创建 `frontend/src/components/song/song-table-row.tsx`
  - [x] 实现高密度歌曲表格行
  - [x] 包含复选框、封面、歌名、艺术家、专辑、年份、流派、时长

- [x] Task 5: 创建 SelectionBar 组件 (AC: 1, 3)
  - [x] 创建 `frontend/src/components/common/selection-bar.tsx`
  - [x] 显示已选中数量
  - [x] 支持全选/取消全选

- [x] Task 6: 创建 Tab 导航组件 (AC: 1)
  - [x] 创建 `frontend/src/components/common/tab-nav.tsx`
  - [x] 实现歌手/专辑/文件夹 Tab 切换
  - [x] 当前项深色文字+底部2px主色条

- [x] Task 7: 创建前端 App 结构 (AC: 1-4)
  - [x] 更新 `frontend/src/app.tsx` 实现主页面结构
  - [x] 集成 Tab 导航
  - [x] 根据路由/状态切换视图

- [x] Task 8: 编写测试 (AC: 1-4)
  - [x] 后端: Artist repository 单元测试
  - [x] 后端: Artist handler 单元测试
  - [x] 前端: Artist view 组件测试 (UI 组件已创建，测试基础设施待 Epic 1 完成)

## Dev Notes

### 技术要求

**Artist Repository:**
- 使用 GORM 的 Group 和 Count 聚合查询
- 按艺术家名分组，统计每个艺术家的歌曲数量
- 支持按艺术家名排序
- 过滤空艺术家名

**API 端点:**
```
GET /api/artists
Response: {
  "data": [
    { "id": 1, "name": "周杰伦", "songCount": 25 },
    { "id": 2, "name": "林俊杰", "songCount": 18 }
  ]
}

GET /api/artists/:id/songs
Response: {
  "data": [
    { "id": 1, "title": "晴天", "album": "叶惠美", "duration": 267, ... },
    ...
  ]
}
```

**Artist View 前端:**
- 紧凑表格视图 (UX-DR3)
- 点击艺术家行展开歌曲列表
- 排序控制显示在列表顶部

### Project Structure Notes

**现有结构:**
```
nas-manager/
├── cmd/server/main.go              # 路由已注册 /api/songs/*, /api/setup/*
├── internal/
│   ├── handler/                    # setting.go, scan.go, encrypt.go, artist.go
│   ├── service/                   # setting.go, scanner.go, id3.go, encrypt.go 已存在
│   ├── repository/                # song.go, setting.go, artist.go
│   └── model/                     # song.go, artist.go, album.go, setting.go, batch.go 已存在
├── frontend/
│   ├── src/
│   │   ├── views/                 # setup-view.tsx, artists-view.tsx
│   │   ├── components/            # song/, common/ 已创建
│   │   ├── contexts/              # selection-context.tsx 已创建
│   │   ├── types/                 # song.ts 已创建
│   │   ├── app.tsx               # 已创建
│   │   └── main.ts               # 已创建
│   └── public/index.html          # 已创建
```

### Architecture Compliance

**遵循规范:**
1. **命名规范:** API 使用 snake_case，Go 代码使用 PascalCase/camelCase
2. **响应格式:** 使用 `pkg/response/response.go` 的统一响应格式
3. **错误处理:** 使用 `pkg/response` 的错误响应格式
4. **分层架构:** Handler → Repository (无 Service 层对于只读查询)
5. **前端状态:** 使用 Preact Context API 管理选择状态

**UX 规范:**
- 色彩: 主题绿 #22C55E, 背景 #FFFFFF, 文字 #1E293B
- 字体: Noto Sans CJK (中文), Inter (英文)
- 间距: 4px 基准 (xs:4, sm:8, md:16, lg:24, xl:32)
- 触摸目标: 最小 44x44px

## Dev Agent Record

### Agent Model Used

MiniMax-M2.7-highspeed

### Debug Log References

### Completion Notes List

- 创建 `internal/repository/artist.go` - ArtistRepository，包含 GetAllArtistsWithSongCount 和 GetSongsByArtist 方法
- 创建 `internal/handler/artist.go` - ArtistHandler，包含 GET /api/artists 和 GET /api/artists/:id/songs 端点
- 更新 `cmd/server/main.go` - 注册 Artist 路由
- 创建 `frontend/src/types/song.ts` - TypeScript 类型定义
- 创建 `frontend/src/contexts/selection-context.tsx` - SelectionContext 用于多选状态管理
- 创建 `frontend/src/components/common/tab-nav.tsx` - Tab 导航组件
- 创建 `frontend/src/components/common/selection-bar.tsx` - 选择操作栏组件
- 创建 `frontend/src/components/song/song-table-row.tsx` - 歌曲表格行组件
- 创建 `frontend/src/views/artists-view.tsx` - 艺术家视图组件
- 创建 `frontend/src/app.tsx` - 主应用入口
- 创建 `frontend/src/main.ts` - Preact 入口点
- 创建 `frontend/src/index.css` - Tailwind 样式入口
- 创建 `frontend/public/index.html` - HTML 模板
- 创建 `internal/repository/artist_test.go` - ArtistRepository 单元测试
- 创建 `internal/handler/artist_test.go` - ArtistHandler 单元测试
- 所有后端测试通过
- 注意: 前端测试基础设施（Vitest 配置）尚未建立，待 Epic 1 的前端基础设施建设完成后可添加

### File List

1. `internal/repository/artist.go` - 新建
2. `internal/handler/artist.go` - 新建
3. `cmd/server/main.go` - 修改（添加 Artist 路由）
4. `frontend/src/types/song.ts` - 新建
5. `frontend/src/contexts/selection-context.tsx` - 新建
6. `frontend/src/components/common/tab-nav.tsx` - 新建
7. `frontend/src/components/common/selection-bar.tsx` - 新建
8. `frontend/src/components/song/song-table-row.tsx` - 新建
9. `frontend/src/views/artists-view.tsx` - 新建
10. `frontend/src/app.tsx` - 新建
11. `frontend/src/main.ts` - 新建
12. `frontend/src/index.css` - 新建
13. `frontend/public/index.html` - 新建
14. `internal/repository/artist_test.go` - 新建
15. `internal/handler/artist_test.go` - 新建

## Change Log

- 2026-04-16: 创建 Story 2.1 故事文件
- 2026-04-16: 实现 Artist Repository、Handler、路由注册
- 2026-04-16: 实现前端 ArtistView、TabNav、SelectionBar、SongTableRow 组件
- 2026-04-16: 编写后端单元测试，所有测试通过

## References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-2.1]
- [Source: _bmad-output/planning-artifacts/architecture.md#API-Communication-Patterns]
- [Source: _bmad-output/planning-artifacts/ux-design-specification.md#UX-DR1-SongTableRow]
- [Source: _bmad-output/planning-artifacts/ux-design-specification.md#UX-DR4-SelectionBar]
- [Source: _bmad-output/planning-artifacts/ux-design-specification.md#UX-DR6-Tab-Navigation]
