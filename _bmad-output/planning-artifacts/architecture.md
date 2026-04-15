---
stepsCompleted: [1, 2, 3, 4, 5, 6, 7, 8]
inputDocuments:
  - "/home/chenlichao/workspace/nas-manager/_bmad-output/planning-artifacts/prd.md"
  - "/home/chenlichao/workspace/nas-manager/_bmad-output/planning-artifacts/ux-design-specification.md"
workflowType: 'architecture'
project_name: 'nas-manager'
user_name: '立子'
date: '2026-04-15'
lastStep: 8
status: 'complete'
completedAt: '2026-04-16'
---

# Architecture Decision Document

_This document builds collaboratively through step-by-step discovery. Sections are appended as we work through each architectural decision together._

## Project Context Analysis

### Requirements Overview

**Functional Requirements:**
- **FR1-7 音乐扫描与导入**: 目录遍历、ID3解析、SQLite存储、增量扫描、异常处理
- **FR8-14 音乐浏览与组织**: 歌手/专辑/文件夹视图、多选、排序、详情查看、删除
- **FR15-19 播放器与现场编辑**: 点播播放、封面/歌词/时间展示、播放中编辑元数据
- **FR20-23 批量编辑**: 批量修改标签/封面/歌词、在线歌词搜索、撤销支持
- **FR24-25 搜索**: 按文件名/标签搜索
- **FR26-27 系统设置**: 加密密码、首次配置引导

**Non-Functional Requirements:**
- **性能**: 启动≤3秒、二进制≤50MB、UI响应≤200ms、1000首扫描≤10分钟
- **安全**: 凭证加密存储（用户密码派生密钥）
- **平台**: AMD64 + ARM64 架构支持
- **部署**: 单一Go二进制文件，无需Docker

**Scale & Complexity:**
- Primary domain: Web App (Go嵌入式前端)
- Complexity level: Low-Medium
- Estimated architectural components: 8-10 major components
- Cross-cutting concerns: 安全加密、嵌入式打包、跨平台编译

### Technical Constraints & Dependencies

- Go 内嵌前端方案（embed）需要精确的打包策略
- ID3 解析需要可靠的Go库（如ID3库）
- 在线歌词搜索依赖第三方API，需要降级策略
- 凭证加密使用 golang.org/x/crypto (轻量级标准库扩展)

### Cross-Cutting Concerns Identified

1. **嵌入式前端打包**: Go二进制内嵌Preact静态文件
2. **凭证安全**: 用户密码派生密钥，golang.org/x/crypto 加密存储
3. **跨平台构建**: AMD64 + ARM64 单一代码库多平台编译
4. **音乐格式支持**: MP3/FLAC/APE/OGG等主流格式解析覆盖≥90%

## Starter Template Evaluation

### Primary Technology Domain

**Go Web Application with Embedded Frontend** - 基于 Go + Preact 技术栈的 NAS 管理工具

### Technology Stack Decisions (from PRD)

| 组件 | 技术选择 | 理由 |
|------|---------|------|
| 后端语言 | Go 1.24+ | 跨平台编译、单一二进制、嵌入式前端 |
| Web 框架 | Gin | 轻量、成熟、性能好 |
| 前端框架 | Preact | 轻量级 React 替代品 |
| 页面交互 | Preact | 组件化状态驱动 |
| CSS 框架 | Tailwind CSS | 极简主义设计、快速开发 |
| 数据库 | SQLite | 本地存储、无依赖 |
| 构建工具 | Vite | 现代、快速、Preact 官方支持 |

### Selected Approach: Go + Gin + Vite(Preact)

### Project Structure

```
nas-manager/
├── cmd/
│   └── server/
│       └── main.go          # 应用入口
├── internal/
│   ├── handler/             # HTTP handlers
│   ├── service/             # 业务逻辑
│   ├── repository/          # 数据访问
│   └── model/               # 数据模型
├── frontend/
│   ├── src/
│   │   ├── components/      # Preact 组件
│   │   ├── views/           # 页面视图
│   │   ├── app.tsx          # Preact 入口
│   │   └── main.ts          # Vite 入口
│   ├── public/
│   ├── index.html
│   ├── vite.config.ts
│   └── package.json
├── embed/
│   └── static/              # 嵌入式静态文件
├── go.mod
├── go.sum
└── Makefile
```

### Build & Deployment Strategy

**原生构建（用户直接运行）：**
- `make build` → 编译单一二进制文件
- 支持 AMD64 + ARM64 平台

**Docker 构建（跨平台镜像）：**
- 使用 Docker Buildx 支持多平台镜像构建
- 分离构建阶段（builder）和运行阶段
- 多阶段 Dockerfile 优化镜像大小

**Dockerfile 设计：**
```dockerfile
# ---- Builder Stage ----
FROM golang:1.24-alpine AS builder
# 安装 Node.js 用于构建前端
RUN apk add --no-cache nodejs npm
WORKDIR /app
# 复制前端源码并构建
COPY frontend/ ./frontend/
RUN cd frontend && npm install && npm run build
# 复制 Go 源码并交叉编译
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o nas-manager ./cmd/server

# ---- Runtime Stage ----
FROM scratch
COPY --from=builder /app/nas-manager /nas-manager
EXPOSE 8080
ENTRYPOINT ["/nas-manager"]
```

**跨平台 Docker 构建命令：**
```bash
# 启用 buildx
docker buildx create --use
docker buildx inspect --bootstrap

# 构建 AMD64 + ARM64 多平台镜像
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  --tag nas-manager:latest \
  --push \
  .
```

**Makefile 自动化：**
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

### Architectural Decisions Provided

**后端架构：**
- 经典三层架构 (Handler → Service → Repository)
- Gin 中间件处理错误、日志
- SQLite + 加密存储凭证

**前端架构：**
- Preact 组件化开发
- Tailwind CSS 原子化样式
- SPA-like 用户体验

**前后端集成：**
- Go embed 打包前端资源
- API 通过 Gin REST endpoints
- 前端通过 fetch 调用后端 API

**部署架构：**
- 单一 `nas-manager` 二进制文件
- AMD64 + ARM64 多平台原生构建
- Docker 多平台镜像支持

## Core Architectural Decisions

### Decision Priority Analysis

**Critical Decisions (Block Implementation):**
- 数据库 ORM: GORM
- 加密方案: golang.org/x/crypto (标准库扩展)
- 配置管理: 无配置文件，SQLite 存储

**Important Decisions (Shape Architecture):**
- API 设计风格: 混合模式
- 前端状态管理: Context API

**Deferred Decisions (Post-MVP):**
- AI 元数据自动补全
- 歌词翻译功能

### Data Architecture

**ORM: GORM**

| 属性 | 决策 |
|------|------|
| 库 | gorm.io/gorm |
| 版本 | v2.x (最新稳定) |
| 理由 | 功能丰富、自动迁移、关联查询方便 |

**数据模型：**

```go
// Song - 音乐文件元数据
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
    CoverPath string    // 封面图片路径
    Lyrics    string    // 内嵌歌词
    FileHash  string    `gorm:"index"`
    FileSize  int64
    CreatedAt time.Time
    UpdatedAt time.Time
}

// Artist - 艺术家（冗余存储用于快速查询）
type Artist struct {
    ID   uint   `gorm:"primaryKey"`
    Name string `gorm:"uniqueIndex;not null"`
}

// Album - 专辑（冗余存储用于快速查询）
type Album struct {
    ID     uint   `gorm:"primaryKey"`
    Name   string `gorm:"index"`
    Artist string
}

// Setting - 应用配置（存储在 SQLite 中）
type Setting struct {
    Key   string `gorm:"primaryKey"`
    Value string
}

// BatchOperation - 批量操作记录（支持撤销）
type BatchOperation struct {
    ID        uint      `gorm:"primaryKey"`
    Type      string    // "update", "delete"
    TargetIDs string    // JSON array of song IDs
    OldValues string    // JSON of previous values
    NewValues string    // JSON of new values
    CreatedAt time.Time
}
```

### Authentication & Security

**凭证加密: golang.org/x/crypto**

| 属性 | 决策 |
|------|------|
| 库 | golang.org/x/crypto |
| 算法 | AES-256-GCM |
| 密钥派生 | PBKDF2 (用户密码 → 加密密钥) |
| 理由 | Go 官方维护、成熟稳定、仅对凭证加密 |

**加密实现策略：**

```go
import "golang.org/x/crypto/pbkdf2"

// 首次设置密码时
func DeriveKey(password string, salt []byte) []byte {
    return pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)
}

// 加密敏感数据 (使用 golang.org/x/crypto/chacha20poly1305 或 aes/gcm)
func Encrypt(plaintext, key []byte) ([]byte, error) {
    // ChaCha20-Poly1305 或 AES-GCM
}

// 解密敏感数据
func Decrypt(ciphertext, key []byte) ([]byte, error) {
    // ChaCha20-Poly1305 或 AES-GCM
}
```

**配置存储策略：**
- **SQLite 路径**: 通过命令行参数 `--db` 或环境变量 `NAS_MANAGER_DB` 指定
- **其他配置**: 全部存储在 SQLite `settings` 表中
- **无配置文件**: 不使用 YAML/JSON/TOML 配置文件

### API & Communication Patterns

**API 设计: 混合模式**

| 资源 | 风格 | 理由 |
|------|------|------|
| `/songs` | 扁平 | 简单列表查询 |
| `/artists` | 扁平 | 简单列表查询 |
| `/albums` | 扁平 | 简单列表查询 |
| `/songs/:id` | 扁平 | 单个资源操作 |
| `/songs/:id/lyrics` | 嵌套 | 歌词是歌曲的子资源 |
| `/batch` | 扁平 | 批量操作端点 |
| `/settings` | 扁平 | 应用配置 |

**主要 API 端点：**

```
GET    /api/songs              # 列表（支持分页、筛选）
GET    /api/songs/:id          # 单曲详情
PUT    /api/songs/:id          # 更新单曲
DELETE /api/songs/:id          # 删除单曲
POST   /api/songs/scan         # 触发扫描

GET    /api/artists            # 艺术家列表
GET    /api/artists/:id/songs  # 艺术家的歌曲

GET    /api/albums             # 专辑列表

POST   /api/batch/update       # 批量更新
POST   /api/batch/delete       # 批量删除
POST   /api/batch/undo/:id     # 撤销批量操作

GET    /api/settings            # 获取配置
PUT    /api/settings            # 更新配置

POST   /api/auth/setup          # 首次设置密码
POST   /api/auth/verify        # 验证密码
```

### Frontend Architecture

**状态管理: Preact Context API**

| 属性 | 决策 |
|------|------|
| 方案 | Preact Context |
| 理由 | 内置、无额外依赖、够用 |

**Context 结构：**

```typescript
// PlayerContext - 播放器状态
const PlayerContext = createContext<{
    currentSong: Song | null;
    isPlaying: boolean;
    play: (song: Song) => void;
    pause: () => void;
}>(null);

// SelectionContext - 多选状态
const SelectionContext = createContext<{
    selected: Set<number>;
    toggle: (id: number) => void;
    selectAll: () => void;
    clear: () => void;
}>(null);

// EditContext - 编辑状态
const EditContext = createContext<{
    editing: Song | null;
    startEdit: (song: Song) => void;
    cancelEdit: () => void;
}>(null);
```

**HTMX 局部刷新场景：**
- 列表切换（歌手/专辑/文件夹视图）
- 搜索结果更新
- 批量操作后的列表刷新

**Preact 组件：**
- 播放器组件（动态交互）
- 批量编辑弹窗（动态交互）
- Toast 通知（动态交互）

### Configuration Management

**无配置文件原则：**

| 配置项 | 来源 |
|--------|------|
| SQLite 路径 | `--db` 参数 或 `NAS_MANAGER_DB` 环境变量 |
| 服务端口 | `--port` 参数 或 `NAS_MANAGER_PORT` 环境变量（默认 8080） |
| 音乐目录 | 首次配置时写入 SQLite settings 表 |
| 加密密码 | 首次配置时设置，用于派生密钥 |
| 其他配置 | 存储在 SQLite settings 表 |

**启动流程：**

```
1. 检查 --db 参数或 NAS_MANAGER_DB 环境变量
2. 如未设置，检查默认路径 ~/.nas-manager/nas-manager.db
3. 如数据库不存在，触发首次配置向导
4. 首次配置写入 music_dir 到 SQLite settings 表
```

**命令行参数：**

```bash
nas-manager --db /path/to/database.db --port 8080
```

### Decision Impact Analysis

**实现顺序：**

1. 项目结构 + GORM 模型定义
2. 加密模块实现（AES-256-GCM）
3. SQLite 初始化 + settings 表
4. Gin 路由 + API 端点
5. 前端项目初始化（Vite + Preact）
6. Preact Context 状态管理
7. HTMX 页面结构
8. 播放器组件
9. 批量编辑组件
10. Docker 构建配置

**跨组件依赖：**

- 加密模块 → 所有敏感数据存储
- Settings 表 → 应用配置管理
- GORM → 数据库操作层

## Implementation Patterns & Consistency Rules

### Naming Conventions

**数据库命名 (snake_case):**

| 元素 | 格式 | 示例 |
|------|------|------|
| 表名 | 复数 | songs, artists, albums |
| 列名 | snake_case | file_path, cover_path |
| 外键 | resource_id | artist_id, album_id |
| 索引 | idx_resource_column | idx_songs_artist |

**API 命名:**

| 元素 | 格式 | 示例 |
|------|------|------|
| 端点 | /复数 | /api/songs, /api/artists |
| 路由参数 | :id | /songs/:id |
| JSON 字段 | snake_case | song_id, artist_name |

**Go 代码命名 (标准 Go 风格):**

| 元素 | 格式 | 示例 |
|------|------|------|
| 结构体 | PascalCase | Song, Artist |
| 函数 | PascalCase | GetSongByID |
| 变量 | camelCase | songID, filePath |
| 包 | 简短小写 | handler, service |

**TypeScript 代码命名:**

| 元素 | 格式 | 示例 |
|------|------|------|
| 组件 | PascalCase | SongCard, PlayerPanel |
| 函数/变量 | camelCase | getSongs, selectedIds |
| 文件 | kebab-case | song-card.tsx |
| 常量 | UPPER_SNAKE | MAX_BATCH_SIZE |

### Project Structure

```
nas-manager/
├── cmd/
│   └── server/
│       └── main.go           # 应用入口
├── internal/
│   ├── handler/              # HTTP handlers
│   │   ├── song.go           # 歌曲相关 handlers
│   │   ├── artist.go         # 艺术家相关 handlers
│   │   ├── batch.go          # 批量操作 handlers
│   │   └── setting.go        # 设置 handlers
│   ├── service/              # 业务逻辑
│   │   ├── song.go
│   │   ├── scanner.go        # 音乐扫描服务
│   │   ├── encrypt.go        # 加密服务
│   │   └── lyrics.go         # 歌词搜索服务
│   ├── repository/           # 数据访问层
│   │   ├── song.go
│   │   ├── artist.go
│   │   ├── album.go
│   │   └── setting.go
│   └── model/                # 数据模型
│       ├── song.go
│       ├── artist.go
│       └── setting.go
├── pkg/
│   ├── crypto/               # 加密工具包
│   │   └── encrypt.go
│   ├── id3/                  # ID3 解析工具包
│   │   └── parser.go
│   └── response/             # API 响应工具
│       └── response.go
├── frontend/
│   ├── src/
│   │   ├── components/       # Preact 组件
│   │   │   ├── player/
│   │   │   │   ├── player.tsx
│   │   │   │   └── lyrics.tsx
│   │   │   ├── edit/
│   │   │   │   └── batch-edit.tsx
│   │   │   └── common/
│   │   │       ├── toast.tsx
│   │   │       └── loading.tsx
│   │   ├── views/            # 页面视图
│   │   │   ├── songs.tsx
│   │   │   ├── artists.tsx
│   │   │   └── settings.tsx
│   │   ├── contexts/         # Preact Context
│   │   │   ├── player.tsx
│   │   │   ├── selection.tsx
│   │   │   └── edit.tsx
│   │   ├── hooks/            # 自定义 hooks
│   │   ├── utils/            # 工具函数
│   │   ├── app.tsx
│   │   └── main.ts
│   ├── public/
│   │   └── index.html
│   ├── vite.config.ts
│   └── package.json
├── embed/
│   └── static/               # 嵌入式静态文件
├── migrations/               # 数据库迁移
│   └── 001_initial.sql
├── Dockerfile
├── Makefile
├── go.mod
└── go.sum
```

### API Response Format

**成功响应:**

```json
{
  "data": {
    "id": 1,
    "title": "歌曲名",
    "artist": "艺术家"
  }
}
```

**错误响应:**

```json
{
  "error": {
    "code": "SONG_NOT_FOUND",
    "message": "歌曲不存在"
  }
}
```

**错误码规范:**

| 错误码 | 说明 |
|--------|------|
| SONG_NOT_FOUND | 歌曲不存在 |
| ARTIST_NOT_FOUND | 艺术家不存在 |
| BATCH_OPERATION_FAILED | 批量操作失败 |
| ENCRYPTION_FAILED | 加密/解密失败 |
| VALIDATION_ERROR | 参数验证失败 |
| SCAN_FAILED | 扫描失败 |
| UNAUTHORIZED | 未授权访问 |

### Error Handling Patterns

**Go 错误处理:**

```go
// 使用 errors.Wrap 添加上下文
if err != nil {
    return nil, fmt.Errorf("failed to get song: %w", err)
}

// 自定义错误类型
type NotFoundError struct {
    Resource string
    ID      uint
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s with id %d not found", e.Resource, e.ID)
}
```

**前端错误处理:**

```typescript
// API 调用错误处理
async function fetchSongs() {
    try {
        const res = await fetch('/api/songs');
        if (!res.ok) {
            const err = await res.json();
            throw new Error(err.error.message);
        }
        return await res.json();
    } catch (e) {
        showToast(e.message, 'error');
    }
}
```

### Process Patterns

**加载状态处理:**

- 列表加载: 骨架屏 (Skeleton)
- 操作中: 按钮禁用 + 加载指示器
- 全局加载: 顶部进度条

**重试策略:**

- 网络请求失败: 自动重试 3 次，间隔 1s
- 扫描操作: 失败继续，不阻塞其他文件

**验证时机:**

- 前端: 提交前即时验证
- 后端: 所有输入强制验证

### Enforcement Guidelines

**所有 AI Agent 必须:**

1. 遵循上述命名规范
2. 使用 GORM 的 snake_case 列标签
3. API 响应使用统一格式
4. 错误码使用大写下划线格式
5. Preact 组件放在 components 目录
6. Handler → Service → Repository 分层调用

## Project Structure & Boundaries

### Complete Project Directory Structure

```
nas-manager/
├── cmd/
│   └── server/
│       └── main.go              # 应用入口
├── internal/
│   ├── handler/                # HTTP handlers (Gin)
│   │   ├── song.go            # 歌曲 CRUD
│   │   ├── artist.go          # 艺术家查询
│   │   ├── album.go           # 专辑查询
│   │   ├── batch.go           # 批量操作
│   │   ├── setting.go         # 设置管理
│   │   ├── auth.go            # 认证
│   │   └── scan.go            # 扫描触发
│   ├── service/                # 业务逻辑层
│   │   ├── song.go            # 歌曲业务逻辑
│   │   ├── scanner.go         # 音乐文件扫描
│   │   ├── encrypt.go         # 加密/解密服务
│   │   ├── batch.go           # 批量操作逻辑
│   │   └── setting.go         # 设置业务逻辑
│   ├── repository/             # 数据访问层
│   │   ├── song.go            # 歌曲数据访问
│   │   ├── artist.go          # 艺术家数据访问
│   │   ├── album.go           # 专辑数据访问
│   │   ├── batch.go           # 批量操作记录
│   │   └── setting.go         # 设置数据访问
│   └── model/                  # 数据模型
│       ├── song.go             # Song 模型 + GORM tag
│       ├── artist.go          # Artist 模型
│       ├── album.go           # Album 模型
│       ├── batch.go           # BatchOperation 模型
│       └── setting.go         # Setting 模型
├── pkg/                        # 公共工具包
│   ├── crypto/                # 加密工具
│   │   └── encrypt.go        # golang.org/x/crypto
│   ├── id3/                   # ID3 解析工具
│   │   └── parser.go         # MP3/FLAC/APE/OGG 解析
│   ├── response/              # API 响应工具
│   │   └── response.go       # 统一响应格式
│   └── middleware/            # Gin 中间件
│       └── error.go          # 错误处理中间件
├── frontend/
│   ├── src/
│   │   ├── components/        # Preact 组件
│   │   │   ├── player/
│   │   │   │   ├── player.tsx       # 播放器主组件
│   │   │   │   ├── player-controls.tsx # 播放控制
│   │   │   │   ├── player-progress.tsx # 进度条
│   │   │   │   └── lyrics.tsx        # 歌词展示
│   │   │   ├── edit/
│   │   │   │   ├── batch-edit-panel.tsx # 批量编辑面板
│   │   │   │   ├── edit-field.tsx     # 编辑字段组件
│   │   │   │   └── edit-preview.tsx   # 编辑预览
│   │   │   ├── song/
│   │   │   │   ├── song-table.tsx     # 歌曲列表表格
│   │   │   │   ├── song-row.tsx       # 歌曲行
│   │   │   │   └── song-checkbox.tsx   # 选择框
│   │   │   └── common/
│   │   │       ├── toast.tsx          # Toast 通知
│   │   │       ├── loading.tsx        # 加载状态
│   │   │       ├── modal.tsx          # 模态框
│   │   │       └── spinner.tsx        # 旋转指示器
│   │   ├── views/              # 页面视图
│   │   │   ├── songs-view.tsx        # 歌曲列表页
│   │   │   ├── artists-view.tsx      # 艺术家页
│   │   │   ├── albums-view.tsx       # 专辑页
│   │   │   ├── folders-view.tsx      # 文件夹视图
│   │   │   ├── settings-view.tsx     # 设置页
│   │   │   └── setup-view.tsx       # 首次配置向导
│   │   ├── contexts/          # Preact Context
│   │   │   ├── player-context.tsx    # 播放器状态
│   │   │   ├── selection-context.tsx # 多选状态
│   │   │   └── app-context.tsx      # 全局应用状态
│   │   ├── hooks/             # 自定义 hooks
│   │   │   ├── use-player.ts       # 播放器 hook
│   │   │   ├── use-selection.ts    # 选择 hook
│   │   │   └── use-api.ts          # API 调用 hook
│   │   ├── utils/             # 工具函数
│   │   │   ├── api.ts              # API 请求封装
│   │   │   ├── formatters.ts       # 格式化工具
│   │   │   └── constants.ts        # 常量定义
│   │   ├── types/             # TypeScript 类型
│   │   │   ├── song.ts             # Song 类型定义
│   │   │   ├── api.ts              # API 响应类型
│   │   │   └── context.ts          # Context 类型
│   │   ├── app.tsx           # Preact 应用入口
│   │   ├── main.ts           # Vite 入口
│   │   └── index.css         # Tailwind 入口
│   ├── public/
│   │   └── index.html        # HTML 模板
│   ├── tailwind.config.js    # Tailwind 配置
│   ├── postcss.config.js      # PostCSS 配置
│   ├── vite.config.ts        # Vite 配置
│   ├── tsconfig.json         # TypeScript 配置
│   └── package.json          # 前端依赖
├── embed/
│   └── static/               # Go embed 静态文件
│       ├── index.html        # 嵌入式 HTML
│       ├── assets/           # 编译后的 JS/CSS
│       └── htmx/            # HTMX 片段
├── migrations/                # 数据库迁移
│   └── 001_initial.sql      # 初始数据库结构
├── Dockerfile                # Docker 多平台构建
├── docker-compose.yml        # 本地开发 Docker
├── Makefile                  # 构建自动化
├── .gitignore
├── go.mod
└── go.sum
```

### Architectural Boundaries

**API Boundaries:**

```
外部请求 (浏览器)
    ↓
Gin Router (cmd/server/main.go)
    ↓
┌─────────────────────────────────────┐
│           Handler Layer              │
│  song.go, artist.go, batch.go ...   │
└─────────────────────────────────────┘
    ↓
┌─────────────────────────────────────┐
│           Service Layer              │
│  song.go, scanner.go, encrypt.go   │
└─────────────────────────────────────┘
    ↓
┌─────────────────────────────────────┐
│          Repository Layer           │
│  song.go, artist.go, setting.go    │
└─────────────────────────────────────┘
    ↓
┌─────────────────────────────────────┐
│             SQLite                  │
│     (GORM ORM / gorm.io/gorm)       │
└─────────────────────────────────────┘
```

**Component Boundaries (Frontend):**

```
Preact Contexts (状态管理)
    ↓
┌─────────────────────────────────────┐
│         Preact Components           │
│  player/, edit/, song/, common/    │
└─────────────────────────────────────┘
    ↓
Preact 状态驱动 DOM 更新
```

**Service Boundaries:**

| Service | 职责 | 依赖 |
|---------|------|------|
| scanner.go | 遍历目录、解析 ID3 | id3/parser.go |
| encrypt.go | AES-256-GCM 加密/解密 | pkg/crypto |
| song.go | 歌曲 CRUD | repository |
| batch.go | 批量操作 + 撤销 | repository, song.go |
| setting.go | 配置读写 | repository |

**Data Boundaries:**

| Model | 存储内容 | 访问层 |
|-------|---------|--------|
| Song | 音乐元数据 | repository/song.go |
| Artist | 艺术家名 | repository/artist.go |
| Album | 专辑信息 | repository/album.go |
| Setting | 应用配置 | repository/setting.go |
| BatchOperation | 批量操作记录 | repository/batch.go |

### Requirements to Structure Mapping

**FR1-7 音乐扫描与导入 → scanner.go + id3/parser.go**
- scanner.go: 目录遍历、文件发现
- id3/parser.go: ID3 标签解析
- repository/song.go: 元数据存储

**FR8-14 音乐浏览与组织 → handler/song.go + handler/artist.go + handler/album.go**
- 视图切换 (歌手/专辑/文件夹)
- 列表排序、多选

**FR15-19 播放器与现场编辑 → player/* components**
- PlayerContext: 播放状态
- player.tsx: 播放器组件
- lyrics.tsx: 歌词展示

**FR20-23 批量编辑 → batch.go + edit/* components**
- batch.go: 批量逻辑 + 撤销
- batch-edit-panel.tsx: 编辑面板

**FR24-25 搜索 → handler/song.go (带筛选参数)**

**FR26-27 系统设置 → handler/setting.go + setup-view.tsx**

### Integration Points

**内部通信:**

- Handler → Service: 函数调用
- Service → Repository: 函数调用
- Frontend → Backend: REST API (fetch)

**外部集成:**

- ID3 解析: go.id3 库 (github.com/mikkyang/id3-go)
- 歌词搜索: 第三方 API (待定)

**数据流:**

```
用户选择音乐目录
    ↓
scanner.go 遍历文件
    ↓
id3/parser.go 解析标签
    ↓
repository/song.go 存入 SQLite
    ↓
handler/song.go 提供 API
    ↓
frontend/songs-view.tsx 展示
```

### File Organization Patterns

**配置:**

- SQLite 路径: `--db` 参数 / `NAS_MANAGER_DB` 环境变量
- 无配置文件: 所有配置存 SQLite settings 表

**源码组织:**

- Go: internal/pkg 分离 (internal: 应用代码, pkg: 可复用工具)
- TypeScript: features-based (components/, views/, contexts/)

**测试组织:**

- Go: `*_test.go` 同目录共置
- 前端: `__tests__/` 子目录

**静态资源:**

- 嵌入式: `embed/static/` (编译后)
- 前端: `frontend/public/` (开发时)

## Architecture Validation Results

### Party Mode Review Adjustments

基于 Party Mode 评审反馈，已做以下调整：

| 问题 | 调整 |
|------|------|
| HTMX 测试困难 | 移除 HTMX，改为纯 Preact 组件 |
| 手写加密实现风险 | 改用 golang.org/x/crypto (官方维护) |
| SQLite 并发 | 确认为单用户场景，SQLite 无问题 |

### Coherence Validation ✅

**决策兼容性:**
- Go 1.24+ / Gin / GORM / SQLite — 兼容
- Preact / Tailwind / Vite — 兼容
- golang.org/x/crypto — Go 官方库，稳定可靠
- 所有技术选择无冲突

**模式一致性:**
- 命名规范贯穿数据库 / API / JSON / Go / TS
- 三层架构结构清晰
- API 响应格式统一

### Requirements Coverage Validation ✅

| 功能需求 | 架构支持 |
|---------|---------|
| FR1-7 音乐扫描 | scanner.go + id3/parser.go |
| FR8-14 音乐浏览 | handler/{song,artist,album}.go |
| FR15-19 播放器 | player/* components + PlayerContext |
| FR20-23 批量编辑 | batch.go + edit/* components |
| FR24-25 搜索 | handler/song.go |
| FR26-27 系统设置 | handler/setting.go |

**非功能需求覆盖:**

| NFR | 架构支持 |
|-----|---------|
| 启动≤3秒 | Go 编译 + 嵌入式前端 |
| 二进制≤50MB | 无 CGO，静态编译 |
| 凭证加密 | golang.org/x/crypto (ChaCha20-Poly1305 或 AES-GCM) |
| AMD64+ARM64 | Go cross-compilation + Docker Buildx |

### Implementation Readiness Validation ✅

**决策完整性:**
- 技术栈已明确，无歧义
- 命名规范清晰
- 加密方案使用标准库，有成熟实现参考

**结构完整性:**
- 目录树完整
- 组件边界清晰
- 集成点已映射

### Architecture Completeness Checklist

**✅ 需求分析**
- [x] 项目上下文分析
- [x] 规模和复杂度评估
- [x] 技术约束识别
- [x] 跨领域关注点映射

**✅ 架构决策**
- [x] 技术栈完全指定
- [x] 集成模式定义
- [x] 性能考虑已处理
- [x] 安全方案确定 (golang.org/x/crypto)

**✅ 实现模式**
- [x] 命名规范确立
- [x] 结构模式定义
- [x] 通信模式指定
- [x] 错误处理模式完整

**✅ 项目结构**
- [x] 完整目录结构定义
- [x] 组件边界建立
- [x] 需求到结构映射完成

### Architecture Readiness Assessment

**总体状态:** ✅ 可开始实现

**置信度:** 高

**关键调整 (Party Mode 后):**
- 移除 HTMX，简化前端架构
- 使用 golang.org/x/crypto 替代手写加密
- 确认单用户场景，SQLite 无并发问题

**优势:**
- 技术栈成熟稳定
- 架构简洁清晰
- 标准库加密方案可靠
- 测试友好 (纯 Preact)

**待完善 (实现时):**
- 歌词搜索 API 源选择
- ID3 解析库最终选型
