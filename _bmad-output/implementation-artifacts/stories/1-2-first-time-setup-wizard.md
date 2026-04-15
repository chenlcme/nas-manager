# Story 1.2: 首次配置向导

Status: review

## Story

As a 用户，
I want 在首次访问时通过引导完成基础配置，
So that 快速开始使用音乐管理功能。

## Acceptance Criteria

1. **Given** 用户首次访问应用（数据库为空或无 music_dir 配置） **When** 系统检测到未配置状态 **Then** 显示首次配置向导界面，引导用户完成两步配置

2. **And** 引导步骤1：设置音乐目录路径（支持输入路径或选择文件夹）

3. **And** 引导步骤2：确认 SQLite 存储路径（默认 ~/.nas-manager/nas-manager.db，可自定义）

4. **And** 配置完成后，music_dir 和 db_path 保存到 SQLite settings 表

5. **Given** 配置向导完成 **When** 用户提交配置 **Then** 验证路径有效性（音乐目录存在且可读、SQLite 路径可写）

6. **And** 保存配置到 settings 表

7. **And** 跳转到音乐库浏览页面

8. **And** 显示 Toast 提示"配置已保存"

## Tasks / Subtasks

- [x] Task 1: 创建设置数据访问层 (AC: 1-8)
  - [x] 创建 `internal/repository/setting.go` - Setting 仓储
  - [x] 实现 GetSetting / SetSetting 方法
  - [x] 实现检查配置是否完成的逻辑

- [x] Task 2: 创建设置服务层 (AC: 1-8)
  - [x] 创建 `internal/service/setting.go` - 设置服务
  - [x] 实现 CheckSetupRequired 检查是否需要首次配置
  - [x] 实现 SaveSetupConfig 保存配置

- [x] Task 3: 创建设置 Handler 和 API 路由 (AC: 1-8)
  - [x] 创建 `internal/handler/setting.go` - 设置处理器
  - [x] 实现 GET /api/setup/status 端点（检查是否需要配置）
  - [x] 实现 POST /api/setup 完成配置

- [x] Task 4: 实现前端首次配置向导组件 (AC: 1-8)
  - [x] 创建 `frontend/src/views/setup-view.tsx` - 配置向导视图
  - [x] 实现步骤1：音乐目录设置
  - [x] 实现步骤2：SQLite 路径设置
  - [x] 实现表单验证和提交

- [x] Task 5: 集成配置检查到主入口 (AC: 1)
  - [x] 修改 `cmd/server/main.go` - 添加 Gin 路由和处理器

## Dev Notes

### 技术要求

**后端实现：**
- 使用 Story 1.1 创建的 GORM 模型和数据库连接
- Setting 仓储负责 settings 表的读写
- API 端点返回 JSON 响应

**前端实现：**
- Preact 组件（遵循 Story 1.1 建立的模式）
- 两步表单：步骤1音乐目录，步骤2数据库路径
- 表单验证：目录存在性检查、可写性检查

### API 设计

```
GET /api/setup/status
Response: { "needs_setup": true/false, "music_dir": "...", "db_path": "..." }

POST /api/setup
Body: { "music_dir": "...", "db_path": "..." }
Response: { "success": true }
```

### 配置存储

配置保存在 SQLite settings 表：
- `music_dir`: 音乐目录路径
- `db_path`: SQLite 数据库路径（可选）

### 项目结构延续

- Handler → Service → Repository 分层结构（Story 1.1 建立）
- Gin 路由模式
- GORM 数据库操作

## Dev Agent Record

### Agent Model Used

MiniMax-M2.7-highspeed

### Completion Notes List

- 创建了 SettingRepository 实现 settings 表的 CRUD 操作
- 创建了 SettingService 实现配置检查和保存逻辑
- 创建了 SettingHandler 实现 GET /api/setup/status 和 POST /api/setup 端点
- 创建了前端 SetupView 组件实现两步配置向导
- 更新了 main.go 集成 Gin 路由和处理器
- 添加了仓储层和服务层单元测试，全部通过

## File List

1. `internal/repository/setting.go` - Setting 仓储实现
2. `internal/repository/setting_test.go` - 仓储单元测试
3. `internal/service/setting.go` - 设置服务实现
4. `internal/service/setting_test.go` - 服务单元测试
5. `internal/handler/setting.go` - 设置 Handler 实现
6. `pkg/response/response.go` - API 响应工具
7. `frontend/src/views/setup-view.tsx` - 配置向导前端组件
8. `cmd/server/main.go` - 更新：添加 Gin 路由

## Change Log

- 2026-04-15: 初始实现 Story 1.2 所有任务

### Review Findings

- [x] [Review][Defer] `SaveSetupConfig` db_path 可写性验证缺失 [internal/service/setting.go] — deferred，代码已通过 `isDirWritable` 检查，非阻塞
- [x] [Review][Defer] `isDirWritable` 原子性问题 [internal/service/setting.go] — deferred，`defer os.Remove` 可处理
