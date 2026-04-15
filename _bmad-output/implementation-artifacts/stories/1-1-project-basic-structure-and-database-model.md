# Story 1.1: 项目基础结构与数据库模型

Status: review

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a 系统，
I want 建立项目基础结构和数据库模型，
So that 为后续功能开发提供坚实基础。

## Acceptance Criteria

1. **Given** 空项目目录 **When** 开发者初始化 Go 项目并创建基础目录结构 **Then** 创建 `cmd/server/main.go`、`internal/{handler,service,repository,model}`、`frontend/`、`pkg/` 等目录结构

2. **And** 创建 `go.mod` 文件，依赖包括 gin、gorm、golang.org/x/crypto 等

3. **And** 创建 `Makefile` 支持 `build`、`build-docker`、`build-all` 等构建目标

4. **Given** 项目结构已建立 **When** 开发者创建 GORM 数据模型 **Then** 创建 Song 模型（ID、FilePath、Title、Artist、Album、Year、Genre、TrackNum、Duration、CoverPath、Lyrics、FileHash、FileSize、CreatedAt、UpdatedAt）

5. **And** 创建 Artist 模型（ID、Name）

6. **And** 创建 Album 模型（ID、Name、Artist）

7. **And** 创建 Setting 模型（Key、Value）

8. **And** 创建 BatchOperation 模型（ID、Type、TargetIDs、OldValues、NewValues、CreatedAt）

9. **And** GORM auto-migration 正确创建所有表

## Tasks / Subtasks

- [x] Task 1: 创建项目目录结构 (AC: 1)
  - [x] 创建 `cmd/server/` 目录和 `main.go` 入口文件
  - [x] 创建 `internal/{handler,service,repository,model}` 目录结构
  - [x] 创建 `pkg/` 目录结构
  - [x] 创建 `frontend/` 目录结构
  - [x] 创建 `embed/` 目录结构

- [x] Task 2: 初始化 Go 模块和依赖 (AC: 2)
  - [x] 运行 `go mod init`
  - [x] 添加依赖：gin、gorm、sqlite3、golang.org/x/crypto

- [x] Task 3: 创建 Makefile 构建脚本 (AC: 3)
  - [x] 实现 `build` 目标
  - [x] 实现 `build-docker` 目标（多平台）
  - [x] 实现 `build-all` 目标（amd64 + arm64）
  - [x] 实现 `clean` 目标

- [x] Task 4: 创建 GORM 数据模型 (AC: 4-8)
  - [x] 创建 `internal/model/song.go` - Song 模型
  - [x] 创建 `internal/model/artist.go` - Artist 模型
  - [x] 创建 `internal/model/album.go` - Album 模型
  - [x] 创建 `internal/model/setting.go` - Setting 模型
  - [x] 创建 `internal/model/batch.go` - BatchOperation 模型

- [x] Task 5: 实现 GORM Auto-Migration (AC: 9)
  - [x] 在 `cmd/server/main.go` 中初始化数据库连接
  - [x] 配置 GORM 自动迁移所有模型
  - [x] 处理数据库连接错误

## Dev Notes

### 技术栈与依赖版本

| 组件 | 技术选择 | 版本 |
|------|---------|------|
| Go | Go 1.24+ | 最新稳定 |
| Web 框架 | Gin | github.com/gin-gonic/gin v1.10+ |
| ORM | GORM | gorm.io/gorm v2.x |
| SQLite 驱动 | modernc.org/sqlite | 纯 Go 实现，无 CGO |
| 加密 | golang.org/x/crypto | Go 官方库 |

### 项目结构规范（Architecture §Project Structure）

```
nas-manager/
├── cmd/
│   └── server/
│       └── main.go              # 应用入口
├── internal/
│   ├── handler/                # HTTP handlers (Gin)
│   ├── service/                # 业务逻辑层
│   ├── repository/             # 数据访问层
│   └── model/                  # 数据模型
├── pkg/                        # 公共工具包
│   ├── crypto/                # 加密工具
│   ├── id3/                   # ID3 解析工具
│   ├── response/              # API 响应工具
│   └── middleware/            # Gin 中间件
├── frontend/                   # 前端源码
│   ├── src/
│   ├── public/
│   ├── vite.config.ts
│   └── package.json
├── embed/                     # 嵌入式静态文件
│   └── static/
├── go.mod
├── go.sum
├── Makefile
└── Dockerfile
```

### 命名规范（Architecture §Naming Conventions）

**数据库命名 (snake_case):**
- 表名：复数形式（songs, artists, albums）
- 列名：snake_case（file_path, cover_path）
- 外键：resource_id（artist_id, album_id）

**Go 代码命名：**
- 结构体：PascalCase（Song, Artist）
- 函数：PascalCase（GetSongByID）
- 变量：camelCase（songID, filePath）
- 包：简短小写（handler, service）

### GORM 模型规范

所有模型必须使用 GORM tag 标注：
```go
type Song struct {
    ID        uint      `gorm:"primaryKey"`
    FilePath  string    `gorm:"uniqueIndex;not null"`
    Title     string
    Artist    string
    Album     string
    Year      int
    Genre     string
    TrackNum  int
    Duration  int       // 秒
    CoverPath string
    Lyrics    string
    FileHash  string    `gorm:"index"`
    FileSize  int64
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### 数据库连接配置

- SQLite 路径通过 `--db` 参数或 `NAS_MANAGER_DB` 环境变量指定
- 默认路径：`~/.nas-manager/nas-manager.db`
- 启动时如果数据库不存在，触发首次配置向导（Story 1.2）

### Makefile 构建目标

```makefile
.PHONY: build build-docker build-all clean

build:
    go build -o nas-manager ./cmd/server

build-docker:
    docker buildx build \
        --platform linux/amd64,linux/arm64 \
        --tag nas-manager:latest \
        --push \
        .

build-all:
    GOOS=linux GOARCH=amd64 go build -o nas-manager-linux-amd64 ./cmd/server
    GOOS=linux GOARCH=arm64 go build -o nas-manager-linux-arm64 ./cmd/server
    GOOS=darwin GOARCH=amd64 go build -o nas-manager-darwin-amd64 ./cmd/server
    GOOS=darwin GOARCH=arm64 go build -o nas-manager-darwin-arm64 ./cmd/server

clean:
    rm -f nas-manager*
```

### 错误处理模式

- 使用 `fmt.Errorf` 和 `errors.Wrap` 添加上下文
- 数据库连接失败应输出友好错误信息并退出
- 关键错误应记录日志

### 项目上下文参考

- **PRD**: `_bmad-output/planning-artifacts/prd.md`
- **Architecture**: `_bmad-output/planning-artifacts/architecture.md`
- **Epic 1 Stories**: `_bmad-output/planning-artifacts/epics.md` (§Epic 1: 项目初始化与音乐扫描)

### 测试标准

- GORM 模型应有基本的单元测试验证模型定义正确
- 数据库连接应有初始化测试
- 遵循 Go 测试规范：`*_test.go` 同目录共置

## Dev Agent Record

### Agent Model Used

MiniMax-M2.7-highspeed

### Debug Log References

### Completion Notes List

- 创建了完整的项目目录结构（cmd/server, internal/{handler,service,repository,model}, pkg/{crypto,id3,response,middleware}, frontend, embed）
- 初始化 Go 模块并添加依赖：gin、gorm、sqlite driver
- 创建 Makefile 支持 build、build-docker、build-all、clean 目标
- 创建了所有 GORM 模型：Song、Artist、Album、Setting、BatchOperation
- 实现了 main.go 包含数据库初始化和 AutoMigrate
- 添加了完整的模型测试（6个测试用例，全部通过）
- 项目成功编译为二进制文件 nas-manager

### File List

1. `cmd/server/main.go` - 应用入口，初始化数据库和 Gin 路由
2. `internal/model/song.go` - Song 模型定义
3. `internal/model/artist.go` - Artist 模型定义
4. `internal/model/album.go` - Album 模型定义
5. `internal/model/setting.go` - Setting 模型定义
6. `internal/model/batch.go` - BatchOperation 模型定义
7. `internal/model/model_test.go` - 模型单元测试
8. `go.mod` - Go 模块定义
9. `go.sum` - Go 依赖锁定
10. `Makefile` - 构建脚本

## Change Log

- 2026-04-15: 初始实现 Story 1.1 所有任务

### Review Findings

- [x] [Review][Patch] `/api/auth/*` 路由未在 main.go 中注册 [cmd/server/main.go] — 已修复
- [x] [Review][Patch] `WalkDir` 替代 `Walk` 防止符号链接循环 [internal/service/scanner.go] — 已修复
- [x] [Review][Defer] `lastScanTime==0` 语义不清晰 [internal/service/scanner.go] — deferred，需要添加 `hasScannedBefore` 标志
- [x] [Review][Decision] Story 1.4 AC 描述修正：增量扫描判断依据改为 file_path + 修改时间 — 已决策
