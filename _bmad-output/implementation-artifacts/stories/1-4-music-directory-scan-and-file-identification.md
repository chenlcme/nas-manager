# Story 1.4: 音乐目录扫描与文件识别

Status: review

## Story

As a 用户，
I want 指定音乐目录并扫描其中的音乐文件，
So that 将音乐文件导入到数据库。

## Acceptance Criteria

1. **Given** 用户已配置音乐目录路径 **When** 用户点击"扫描"按钮触发扫描 **Then** 遍历音乐目录及其子目录

2. **And** 识别支持的音乐格式文件（.mp3, .flac, .ape, .ogg, .m4a, .wav 等）

3. **And** 支持的格式覆盖率 ≥ 90%（按格式种类计算）

4. **Given** 扫描过程 **When** 发现音乐文件 **Then** 检查文件是否已存在于数据库（根据 file_path 判断）

5. **And** 新文件：创建数据库记录，标记需要解析 ID3

6. **And** 已存在文件：跳过或根据文件修改时间判断是否需要重新解析

## Tasks / Subtasks

- [x] Task 1: 创建音乐扫描服务 (AC: 1-6)
  - [x] 创建 `internal/service/scanner.go` - 扫描服务
  - [x] 实现 WalkDirectory 遍历目录
  - [x] 实现 FilterMusicFiles 识别音乐文件

- [x] Task 2: 定义支持的音频格式 (AC: 2-3)
  - [x] 定义支持格式列表（.mp3, .flac, .ape, .ogg, .m4a, .wav, .aiff, .wma, .aac）
  - [x] 确保格式覆盖率 ≥ 90%

- [x] Task 3: 创建扫描 Handler 和 API (AC: 1-6)
  - [x] 创建 `internal/handler/scan.go` - 扫描处理器
  - [x] 实现 POST /api/songs/scan 触发扫描

- [x] Task 4: 创建歌曲仓储 (AC: 4-6)
  - [x] 创建 `internal/repository/song.go` - 歌曲仓储
  - [x] 实现 SongExists 检查歌曲是否存在
  - [x] 实现 CreateSong 创建歌曲记录

- [x] Task 5: 集成扫描服务到主程序 (AC: 1)
  - [x] 更新 main.go 注册扫描路由

- [x] Task 6: 编写扫描模块单元测试 (AC: 1-6)
  - [x] 测试目录遍历
  - [x] 测试文件过滤
  - [x] 测试扫描逻辑

## Dev Notes

### 技术要求

**扫描实现：**
- 使用 `filepath.Walk` 遍历目录
- 支持的音频格式：.mp3, .flac, .ape, .ogg, .m4a, .wav, .aiff, .wma, .aac
- 格式识别通过文件扩展名判断

**ID3 解析（Story 1.5）：**
- 扫描阶段只负责发现文件和创建记录
- ID3 标签解析在 Story 1.5 中实现

### 项目结构延续

- `internal/service/scanner.go` - 扫描服务
- `internal/handler/scan.go` - 扫描处理器
- `internal/repository/song.go` - 歌曲仓储

### API 设计

```
POST /api/songs/scan
Response: { "found": 100, "new": 50, "errors": [] }
```

### 音乐格式支持

主流音频格式覆盖率 ≥ 90%：
- MP3 (.mp3) - 最常见
- FLAC (.flac) - 无损压缩
- APE (.ape) - 无损压缩
- OGG (.ogg) - 开源格式
- M4A (.m4a) - AAC 编码
- WAV (.wav) - 无压缩 PCM
- AIFF (.aiff) - Apple 无损
- WMA (.wma) - Windows 媒体
- AAC (.aac) - 高级音频编码

## Dev Agent Record

### Agent Model Used

MiniMax-M2.7-highspeed

### Completion Notes List

- 创建了 `internal/service/scanner.go` 实现目录遍历和文件过滤
- 支持 9 种主流音频格式，覆盖率 ≥ 90%
- 创建了 `internal/repository/song.go` 实现歌曲数据访问
- 创建了 `internal/handler/scan.go` 实现扫描 API
- 添加了扫描服务和歌曲仓储的单元测试，全部通过

## File List

1. `internal/service/scanner.go` - 扫描服务实现
2. `internal/service/scanner_test.go` - 扫描服务单元测试
3. `internal/handler/scan.go` - 扫描处理器实现
4. `internal/repository/song.go` - 歌曲仓储实现
5. `internal/repository/song_test.go` - 歌曲仓储单元测试

## Change Log

- 2026-04-15: 初始实现 Story 1.4 所有任务

### Review Findings

- [x] [Review][Patch] `filepath.Walk` 替换为 `filepath.WalkDir` [internal/service/scanner.go] — 已修复，防止符号链接循环
- [x] [Review][Decision] AC 修正：增量扫描判断依据改为 file_path + 修改时间（而非仅 file_path）— 已决策并修正
