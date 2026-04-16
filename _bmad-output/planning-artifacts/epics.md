---
stepsCompleted:
  - "step-01-validate-prerequisites"
  - "step-02-design-epics"
  - "step-03-create-stories"
inputDocuments:
  - "/home/chenlichao/workspace/nas-manager/_bmad-output/planning-artifacts/prd.md"
  - "/home/chenlichao/workspace/nas-manager/_bmad-output/planning-artifacts/architecture.md"
  - "/home/chenlichao/workspace/nas-manager/_bmad-output/planning-artifacts/ux-design-specification.md"
---

# nas-manager - Epic Breakdown

## Overview

This document provides the complete epic and story breakdown for nas-manager, decomposing the requirements from the PRD, UX Design if it exists, and Architecture requirements into implementable stories.

## Requirements Inventory

### Functional Requirements

FR1: 用户可以指定音乐目录路径
FR2: 系统可以遍历指定目录，识别音乐文件（支持 MP3/FLAC/APE/OGG 等主流格式）
FR3: 系统可以解析音乐文件的 ID3 标签（标题/艺术家/专辑/年份/曲目号等）
FR4: 系统可以将解析后的元数据存储到 SQLite 数据库
FR5: 用户可以触发重新扫描（支持增量扫描和全量扫描，增量扫描依据文件创建时间和最后修改时间判定）
FR6: 系统可以检测并清理孤岛记录（文件已删除但数据库仍有记录）
FR7: 系统可以处理扫描异常（文件损坏、权限不足、网络超时），并向用户报告
FR8: 用户可以按歌手视图浏览音乐库
FR9: 用户可以按专辑视图浏览音乐库
FR10: 用户可以按文件夹结构浏览音乐库
FR11: 用户可以在浏览视图中按名称/时长/添加时间排序
FR12: 用户可以查看单曲的详细信息（ID3标签/封面/歌词）
FR13: 用户可以在列表中多选多个音乐文件
FR14: 用户可以删除选中的音乐文件（同时删除文件和数据库记录）
FR15: 用户可以播放选中的音乐文件（点一首播一首，无播放队列）
FR16: 播放器可以展示专辑封面图片
FR17: 播放器可以展示歌词（同步或静态）
FR18: 播放器可以展示播放时间信息
FR19: 用户在播放过程中可以直接编辑当前歌曲的标签/封面/歌词
FR20: 用户可以批量修改选中音乐的标签（艺术家/专辑/标题等）
FR21: 用户可以批量修改选中音乐的封面图片
FR22: 用户可以搜索并批量应用歌词（系统按文件名分隔的词片段在线搜索歌词，搜索失败时用户可手动重试，API来源在实现时确定）
FR23: 系统支持依次撤销批量编辑操作（每次撤销一个批量记录）
FR24: 用户可以按文件名搜索音乐
FR25: 用户可以按 ID3 标签内容搜索音乐
FR26: 用户可以设置加密密码，系统基于密码生成加密密钥
FR27: 系统可以在首次访问时引导用户完成基础配置（设置音乐目录和SQLite存储路径两步完成）

### NonFunctional Requirements

NFR1: 启动时间 ≤ 3秒
NFR2: 二进制文件大小 ≤ 50MB
NFR3: UI 操作响应 ≤ 200ms
NFR4: 首屏加载（扫描结果展示） ≤ 2秒
NFR5: 1000首音乐扫描入库 ≤ 10分钟
NFR6: 凭证加密 - 云存储凭证使用用户密码生成的密钥加密存储
NFR7: 密钥管理 - 密钥由用户密码派生，丢失后凭证无法解密（用户需重新配置凭证）
NFR8: 单用户本地部署，无并发用户需求
NFR9: SQLite 存储，本地 NAS 性能足够
NFR10: 架构支持 AMD64 + ARM64
NFR11: 部署方式 - 单一 Go 二进制文件，直接运行
NFR12: 浏览器支持 - Chrome（桌面端 + 移动端）

### Additional Requirements

- 起始模板：Go + Gin + Vite(Preact) 技术栈
- 项目结构：cmd/server/main.go + internal/{handler,service,repository,model} + frontend/
- ORM：使用 GORM (gorm.io/gorm)
- 加密方案：golang.org/x/crypto (AES-256-GCM 或 ChaCha20-Poly1305)
- 密钥派生：PBKDF2 (用户密码 → 加密密钥)
- 配置管理：无配置文件，SQLite 存储，--db 参数或 NAS_MANAGER_DB 环境变量指定路径
- 前端架构：Preact + Tailwind CSS + Vite，SPA-like 体验
- 状态管理：Preact Context API
- API 风格：扁平 REST API，/api/songs, /api/artists, /api/albums 等
- 音乐格式支持覆盖率 ≥ 90%（按格式种类计算）
- 嵌入式前端：Go embed 打包前端资源
- 跨平台构建：支持 AMD64 + ARM64 原生构建 + Docker 多平台镜像

### UX Design Requirements

UX-DR1: 实现 SongTableRow 组件 - 高密度歌曲表格行，包含选择复选框、封面缩略图(40x40px)、歌名、艺术家、专辑、年份、流派、时长、操作菜单
UX-DR2: 实现 SideEditPanel 组件 - 右侧滑入批量编辑面板，宽度320px，包含已选中歌曲预览、编辑字段表单（专辑/艺术家/年份/风格）、预览区域、操作按钮；字段留空=保持不变
UX-DR3: 实现 SidePlayer 组件 - 右侧固定播放器区域，展示封面大图(200x200px)、歌名/艺术家/专辑、播放进度条、播放控制（上一首/播放暂停/下一首）、歌词显示、编辑按钮
UX-DR4: 实现 SelectionBar 组件 - 固定在列表上方，显示已选中数量、全选/取消全选、批量编辑按钮、取消选择
UX-DR5: 实现 Toast 提示组件 - 右上角显示，成功(绿色)、错误(红色)、警告(橙色)，3秒自动消失
UX-DR6: 实现顶部 Tab 导航 - 歌手/专辑/文件夹切换，当前项深色文字+底部2px主色条
UX-DR7: 实现首次配置向导 - 两步完成：设置音乐目录和SQLite存储路径
UX-DR8: 色彩系统实现 - 主题绿色#22C55E、蓝色#3B82F6、橙色#F97316、白色#FFFFFF、浅绿底色#F0FDF4、深色文字#1E293B、次要文字#64748B
UX-DR9: 字体系统实现 - 中文字体：思源黑体/Noto Sans CJK、英文字体：Inter、等宽字体：JetBrains Mono/SF Mono
UX-DR10: 间距系统实现 - 基于4px基准，xs:4px, sm:8px, md:16px, lg:24px, xl:32px
UX-DR11: 响应式布局实现 - 桌面(≥1024px)紧凑表格视图、平板(768-1023px)表格视图、手机(<768px)卡片视图
UX-DR12: 无障碍支持 - WCAG 2.1 AA级，对比度≥4.5:1，触摸目标最小44x44px，键盘导航支持，ARIA标签

### FR Coverage Map

FR1: Epic 1 - 指定音乐目录路径
FR2: Epic 1 - 遍历识别音乐文件
FR3: Epic 1 - 解析 ID3 标签
FR4: Epic 1 - 存储到 SQLite
FR5: Epic 1 - 增量/全量扫描
FR6: Epic 1 - 清理孤岛记录
FR7: Epic 1 - 扫描异常处理
FR8: Epic 2 - 按歌手视图浏览
FR9: Epic 2 - 按专辑视图浏览
FR10: Epic 2 - 按文件夹浏览
FR11: Epic 2 - 排序
FR12: Epic 2 - 查看单曲详情
FR13: Epic 2 - 多选音乐文件
FR14: Epic 2 - 删除选中音乐
FR15: Epic 3 - 播放选中音乐
FR16: Epic 3 - 展示专辑封面
FR17: Epic 3 - 展示歌词
FR18: Epic 3 - 展示播放时间
FR19: Epic 3 - 播放中编辑元数据
FR20: Epic 4 - 批量修改标签
FR21: Epic 4 - 批量修改封面
FR22: Epic 4 - 搜索并批量应用歌词
FR23: Epic 4 - 撤销批量编辑
FR24: Epic 2 - 按文件名搜索
FR25: Epic 2 - 按标签搜索
FR26: Epic 1 - 设置加密密码
FR27: Epic 1 - 首次配置引导

## Epic List

### Epic 1: 项目初始化与音乐扫描

用户可以完成基础配置，导入音乐目录并扫描所有音乐文件到数据库。

**FRs covered:** FR1, FR2, FR3, FR4, FR5, FR6, FR7, FR26, FR27

### Epic 2: 音乐库浏览与搜索

用户可以高效浏览、搜索和组织音乐库中的音乐。

**FRs covered:** FR8, FR9, FR10, FR11, FR12, FR13, FR14, FR24, FR25

### Epic 3: 播放器与现场编辑

用户可以播放音乐并在播放过程中查看和编辑元数据。

**FRs covered:** FR15, FR16, FR17, FR18, FR19

### Epic 4: 批量编辑与撤销

用户可以高效批量修改元数据并支持撤销操作。

**FRs covered:** FR20, FR21, FR22, FR23

---

## Epic 1: 项目初始化与音乐扫描

用户可以完成基础配置，导入音乐目录并扫描所有音乐文件到数据库。

**FRs covered:** FR1, FR2, FR3, FR4, FR5, FR6, FR7, FR26, FR27

---

### Story 1.1: 项目基础结构与数据库模型

As a 系统，
I want 建立项目基础结构和数据库模型，
So that 为后续功能开发提供坚实基础。

**Acceptance Criteria:**

**Given** 空项目目录
**When** 开发者初始化 Go 项目并创建基础目录结构
**Then** 创建 cmd/server/main.go、internal/{handler,service,repository,model}、frontend/、pkg/ 等目录结构

**And** 创建 go.mod 文件，依赖包括 gin、gorm、golang.org/x/crypto 等

**And** 创建 Makefile 支持 build、build-docker、build-all 等构建目标

---

**Given** 项目结构已建立
**When** 开发者创建 GORM 数据模型
**Then** 创建 Song 模型（ID、FilePath、Title、Artist、Album、Year、Genre、TrackNum、Duration、CoverPath、Lyrics、FileHash、FileSize、CreatedAt、UpdatedAt）

**And** 创建 Artist 模型（ID、Name）

**And** 创建 Album 模型（ID、Name、Artist）

**And** 创建 Setting 模型（Key、Value）

**And** 创建 BatchOperation 模型（ID、Type、TargetIDs、OldValues、NewValues、CreatedAt）

**And** GORM auto-migration 正确创建所有表

---

### Story 1.2: 首次配置向导

As a 用户，
I want 在首次访问时通过引导完成基础配置，
So that 快速开始使用音乐管理功能。

**Acceptance Criteria:**

**Given** 用户首次访问应用（数据库为空或无 music_dir 配置）
**When** 系统检测到未配置状态
**Then** 显示首次配置向导界面，引导用户完成两步配置

**And** 引导步骤1：设置音乐目录路径（支持输入路径或选择文件夹）

**And** 引导步骤2：确认 SQLite 存储路径（默认 ~/.nas-manager/nas-manager.db，可自定义）

**And** 配置完成后，music_dir 和 db_path 保存到 SQLite settings 表

---

**Given** 配置向导完成
**When** 用户提交配置
**Then** 验证路径有效性（音乐目录存在且可读、SQLite 路径可写）

**And** 保存配置到 settings 表

**And** 跳转到音乐库浏览页面

**And** 显示 Toast 提示"配置已保存"

---

### Story 1.3: 加密密码设置

As a 用户，
I want 设置加密密码来保护敏感数据，
So that 云存储凭证等敏感信息可以被安全存储。

**Acceptance Criteria:**

**Given** 用户已完成首次配置
**When** 用户在设置页面设置加密密码
**Then** 密码长度 ≥ 8 字符

**And** 使用 PBKDF2 派生加密密钥（100000 次迭代，32 字节输出）

**And** 密钥和盐值存储在 settings 表中（盐值随机生成）

**And** 后续凭证加密使用 AES-256-GCM 或 ChaCha20-Poly1305

---

**Given** 用户已设置加密密码
**When** 用户修改加密密码
**Then** 验证原密码正确性

**And** 重新派生密钥

**And** 使用新密钥重新加密已有凭证

**And** 原密钥加密的数据无法被新密钥解密时，提示用户重新配置凭证

---

### Story 1.4: 音乐目录扫描与文件识别

As a 用户，
I want 指定音乐目录并扫描其中的音乐文件，
So that 将音乐文件导入到数据库。

**Acceptance Criteria:**

**Given** 用户已配置音乐目录路径
**When** 用户点击"扫描"按钮触发扫描
**Then** 遍历音乐目录及其子目录

**And** 识别支持的音乐格式文件（.mp3, .flac, .ape, .ogg, .m4a, .wav 等）

**And** 支持的格式覆盖率 ≥ 90%（按格式种类计算）

---

**Given** 扫描过程
**When** 发现音乐文件
**Then** 检查文件是否已存在于数据库（根据 file_path 判断）

**And** 新文件：创建数据库记录，标记需要解析 ID3

**And** 已存在文件：跳过或根据文件修改时间判断是否需要重新解析

---

### Story 1.5: ID3 标签解析

As a 系统，
I want 解析音乐文件的 ID3 标签，
So that 提取歌曲元数据用于展示和管理。

**Acceptance Criteria:**

**Given** 扫描到新的音乐文件
**When** 系统解析 ID3 标签
**Then** 提取标题（Title）

**And** 提取艺术家（Artist）

**And** 提取专辑（Album）

**And** 提取年份（Year）

**And** 提取曲目号（TrackNum）

**And** 提取流派（Genre）

**And** 提取时长（Duration，秒）

**And** 提取封面图片（Cover）

**And** 提取内嵌歌词（Lyrics）

**And** 计算文件哈希（FileHash）用于去重

**And** 记录文件大小（FileSize）

---

**Given** ID3 解析过程
**When** 遇到不支持的编码格式或损坏的标签
**Then** 使用默认值（空字符串或 0）

**And** 记录解析警告日志

**And** 继续解析其他文件

---

### Story 1.6: 增量扫描与全量扫描

As a 用户，
I want 选择增量扫描或全量扫描，
So that 高效更新音乐库而不重复处理未变化的歌曲。

**Acceptance Criteria:**

**Given** 用户触发重新扫描
**When** 用户选择"全量扫描"
**Then** 重新解析所有音乐文件的 ID3 标签

**And** 更新数据库中所有相关记录

**And** 保留已有的用户编辑数据（如有）

---

**Given** 用户触发重新扫描
**When** 用户选择"增量扫描"
**Then** 仅处理修改时间晚于上次扫描时间的文件

**And** 新文件：创建新记录

**And** 修改过的文件：更新元数据

**And** 删除过的文件：不处理（由孤岛清理处理）

**And** 增量扫描耗时 < 全量扫描的 10%

---

### Story 1.7: 孤岛记录清理与扫描异常处理

As a 系统，
I want 清理孤岛记录并处理扫描异常，
So that 保证数据库与实际文件系统一致，并给用户清晰的错误反馈。

**Acceptance Criteria:**

**Given** 扫描完成
**When** 系统检测到孤岛记录（数据库有记录但文件不存在）
**Then** 标记这些记录为"待删除"

**And** 提供界面让用户确认删除

**And** 用户确认后，同时删除数据库记录（文件本已不存在）

---

**Given** 扫描过程中遇到异常文件
**When** 遇到文件损坏、权限不足或格式不支持
**Then** 记录错误日志（文件名、错误类型、发生时间）

**And** 继续扫描其余文件

**And** 扫描完成后，返回扫描报告

---

**Given** 扫描报告生成
**When** 扫描包含失败文件
**Then** 显示成功数量和失败数量

**And** 列出失败文件及其错误原因

**And** 提供"重新扫描失败文件"选项

**And** 提供"忽略并继续"选项

---

## Epic 2: 音乐库浏览与搜索

用户可以高效浏览、搜索和组织音乐库中的音乐。

**FRs covered:** FR8, FR9, FR10, FR11, FR12, FR13, FR14, FR24, FR25

**相关 UX-DRs：**
- UX-DR1: SongTableRow 组件
- UX-DR4: SelectionBar 组件
- UX-DR6: 顶部 Tab 导航
- UX-DR11: 响应式布局

---

### Story 2.1: 歌手视图浏览

As a 用户，
I want 按歌手视图浏览音乐库，
So that 快速找到特定艺术家的所有歌曲。

**Acceptance Criteria:**

**Given** 用户切换到"歌手" Tab
**When** 系统加载歌手视图
**Then** 按歌手名分组显示所有歌曲

**And** 每组显示歌手名和歌曲数量

**And** 点击歌手组展开显示该艺术家所有歌曲

**And** 支持按歌手名排序（升序/降序）

---

### Story 2.2: 专辑视图浏览

As a 用户，
I want 按专辑视图浏览音乐库，
So that 快速找到特定专辑的所有歌曲。

**Acceptance Criteria:**

**Given** 用户切换到"专辑" Tab
**When** 系统加载专辑视图
**Then** 按专辑名分组显示所有歌曲

**And** 每组显示专辑名、艺术家名和歌曲数量

**And** 点击专辑组展开显示该专辑所有歌曲

**And** 支持按专辑名排序（升序/降序）

---

### Story 2.3: 文件夹视图浏览

As a 用户，
I want 按文件夹结构浏览音乐库，
So that 按照文件系统组织查看音乐。

**Acceptance Criteria:**

**Given** 用户切换到"文件夹" Tab
**When** 系统加载文件夹视图
**Then** 按音乐目录的文件夹结构显示

**And** 显示文件夹路径和包含的歌曲数量

**And** 点击文件夹展开显示其中的歌曲

**And** 显示文件在磁盘上的实际路径

---

### Story 2.4: 歌曲列表排序

As a 用户，
I want 在浏览视图中对歌曲列表排序，
So that 按照我关心的方式组织音乐。

**Acceptance Criteria:**

**Given** 用户在歌曲列表页面
**When** 用户选择排序方式
**Then** 支持按歌曲名称排序

**And** 支持按时长排序

**And** 支持按添加时间排序

**And** 支持升序/降序切换

**And** 排序选项在列表顶部显示

---

### Story 2.5: 查看单曲详情

As a 用户，
I want 查看单曲的详细信息，
So that 了解歌曲的完整元数据。

**Acceptance Criteria:**

**Given** 用户在歌曲列表中
**When** 用户点击某首歌曲
**Then** 显示歌曲详情面板或弹窗

**And** 展示完整 ID3 信息（标题、艺术家、专辑、年份、流派、曲目号）

**And** 展示封面图片

**And** 展示内嵌歌词

**And** 展示文件信息（路径、大小、时长、格式）

---

### Story 2.6: 多选歌曲文件

As a 用户，
I want 在列表中多选多个音乐文件，
So that 对批量歌曲进行操作。

**Acceptance Criteria:**

**Given** 用户在歌曲列表中
**When** 用户点击歌曲的复选框
**Then** 添加/移除该歌曲的选中状态

**And** 选中状态高亮显示（主题色背景）

**And** SelectionBar 显示已选中的歌曲数量

---

**Given** 用户已选中部分歌曲
**When** 用户点击"全选"
**Then** 选中当前视图中的所有歌曲

**And** SelectionBar 显示已选中数量

---

**Given** 用户已选中歌曲
**When** 用户点击"取消全选"
**Then** 清除所有选中状态

**And** SelectionBar 隐藏或显示数量为 0

---

### Story 2.7: 删除选中音乐

As a 用户，
I want 删除选中的音乐文件，
So that 清理不需要的歌曲（同时删除文件和数据库记录）。

**Acceptance Criteria:**

**Given** 用户选中了 N 首歌曲
**When** 用户点击"删除"按钮
**Then** 显示确认对话框

**And** 提示"将删除 N 首歌曲，文件和数据库记录都将被删除"

**And** 要求用户确认删除

---

**Given** 用户确认删除
**When** 系统执行删除
**Then** 同时删除磁盘上的音频文件

**And** 删除数据库中的歌曲记录

**And** 返回删除结果报告

**And** 展示成功删除数量和失败数量

---

**Given** 删除失败
**When** 部分文件无法删除（如权限不足）
**Then** 记录错误

**And** 继续删除其他文件

**And** 显示失败文件列表及原因

---

### Story 2.8: 按文件名搜索

As a 用户，
I want 按文件名搜索音乐，
So that 快速找到特定文件。

**Acceptance Criteria:**

**Given** 用户在搜索框输入文件名关键词
**When** 用户提交搜索（回车或点击搜索按钮）
**Then** 搜索匹配文件名的歌曲

**And** 支持中文和英文文件名

**And** 搜索结果高亮匹配文字

**And** 无结果时显示"未找到匹配的歌曲"

---

### Story 2.9: 按标签内容搜索

As a 用户，
I want 按 ID3 标签内容搜索音乐，
So that 通过歌曲信息找到目标歌曲。

**Acceptance Criteria:**

**Given** 用户在搜索框输入标签关键词
**When** 用户提交搜索
**Then** 搜索匹配标题、艺术家、专辑的歌曲

**And** 支持多条件组合搜索

**And** 搜索结果高亮匹配文字

**And** 无结果时显示"未找到匹配的歌曲"

---

## Epic 3: 播放器与现场编辑

用户可以播放音乐并在播放过程中查看和编辑元数据。

**FRs covered:** FR15, FR16, FR17, FR18, FR19

**相关 UX-DRs：**
- UX-DR3: SidePlayer 组件
- UX-DR8: 色彩系统实现

---

### Story 3.0: SQLite 连接池修复（技术准备）

As a 开发者，
I want 修复 SQLite 并发访问的连接池问题，
So that 消除测试中的竞态条件并提升数据库稳定性。

**Background:**

Epic 1&2 的代码审查发现 `TestDeleteSongs_*` 测试存在竞态条件。SQLite :memory: 数据库在并发 GORM 连接访问时不稳定。此外，DeleteSongs handler 使用无缓冲 channel 实现超时控制存在 goroutine 泄漏风险。

**Acceptance Criteria:**

**Given** 当前代码库 **When** 运行 `TestDeleteSongs_*` 测试 **Then** 测试稳定通过，无竞态条件

**And** 并发删除操作不会导致数据库连接错误

---

**Given** DeleteSongs handler **When** 执行文件删除超时 **Then** goroutine 不会阻塞泄漏

**And** 使用带缓冲的 channel 正确处理超时

---

**Technical Implementation:**

1. 评估 GORM SQLite 驱动的连接池配置
2. 如需要，配置合理的最大连接数（`SetMaxOpenConns`/`SetMaxIdleConns`）
3. 或考虑使用文件数据库（`:memory:` → 临时文件）进行测试
4. 确保 DeleteSongs 的 channel 实现正确（已修复）

**Dependencies:** 无

**Estimated Effort:** 1 day

---

### Story 3.1: 播放选中音乐

As a 用户，
I want 播放选中的音乐文件，
So that 试听歌曲并验证元数据。

**Acceptance Criteria:**

**Given** 用户在歌曲列表中
**When** 用户点击歌曲的播放按钮
**Then** 启动音频播放器播放该歌曲

**And** 播放器显示在界面固定位置（右侧 SidePlayer）

**And** 点一首播一首，无播放队列

**And** 播放时列表中该歌曲显示播放中状态（左侧主题色条）

---

**Given** 歌曲正在播放
**When** 用户点击"暂停"按钮
**Then** 暂停播放

**And** 播放按钮变为"播放"图标

**And** 进度条停止

---

**Given** 歌曲正在播放或暂停
**When** 用户点击"播放"按钮
**Then** 恢复播放

**And** 进度条继续更新

---

### Story 3.2: 展示专辑封面

As a 用户，
I want 在播放器中看到专辑封面，
So that 直观识别当前播放的歌曲。

**Acceptance Criteria:**

**Given** 歌曲正在播放
**When** 播放器加载歌曲信息
**Then** 展示专辑封面图片（200x200px）

**And** 封面居中显示在播放器顶部

**And** 有封面：显示实际封面图

**And** 无封面：显示灰色占位 + 音符图标 + 橙色虚线边框

---

### Story 3.3: 展示歌词

As a 用户，
I want 在播放器中看到歌词，
So that 跟随歌曲播放查看歌词内容。

**Acceptance Criteria:**

**Given** 歌曲正在播放
**When** 播放器加载歌曲信息
**Then** 如果有内嵌歌词，显示歌词文本

**And** 支持静态歌词显示（逐行滚动）

**And** 如果无歌词，显示"暂无歌词"提示

**And** 歌词区域可滚动

---

### Story 3.4: 展示播放时间

As a 用户，
I want 在播放器中看到播放时间，
So that 了解播放进度。

**Acceptance Criteria:**

**Given** 歌曲正在播放
**When** 播放器加载歌曲
**Then** 显示当前播放时间（格式：mm:ss）

**And** 显示总时长（格式：mm:ss）

**And** 显示播放进度条

**And** 进度条可拖动跳转播放位置

---

### Story 3.5: 播放中编辑元数据

As a 用户，
I want 在播放过程中直接编辑歌曲元数据，
So that 边听边验证边修改，高效完成元数据整理。

**Acceptance Criteria:**

**Given** 歌曲正在播放
**When** 用户点击"编辑"按钮
**Then** 播放器展开编辑面板

**And** 显示当前歌曲的可编辑字段（标题、艺术家、专辑、年份、流派）

**And** 封面和歌词也可编辑

**And** 已修改的字段高亮显示

---

**Given** 用户修改了字段
**When** 用户点击"保存"
**Then** 验证输入有效性

**And** 更新数据库中的歌曲记录

**And** 刷新播放器显示新内容

**And** 继续播放（不中断）

**And** 显示 Toast 提示"已保存"

---

**Given** 用户修改了字段
**When** 用户点击"取消"
**Then** 放弃修改

**And** 恢复原始显示

---

## Epic 4: 批量编辑与撤销

用户可以高效批量修改元数据并支持撤销操作。

**FRs covered:** FR20, FR21, FR22, FR23

**相关 UX-DRs：**
- UX-DR2: SideEditPanel 组件
- UX-DR4: SelectionBar 组件
- UX-DR5: Toast 提示组件

---

### Story 4.1: 批量修改标签

As a 用户，
I want 批量修改选中歌曲的标签，
So that 高效整理一批歌曲的元数据。

**Acceptance Criteria:**

**Given** 用户选中了 N 首歌曲
**When** 用户点击"批量编辑"
**Then** 右侧滑入 SideEditPanel（宽度 320px）

**And** 显示"已选中 N 首"

**And** 显示选中歌曲列表预览

---

**Given** 批量编辑面板展开
**When** 用户填写标签字段（艺术家/专辑/标题等）
**Then** 留空的字段保持不变

**And** 填写的字段显示预览值

**And** 预览区域显示将要做的更改

---

**Given** 用户确认预览
**When** 用户点击"应用"
**Then** 批量更新数据库中的歌曲记录

**And** 创建 BatchOperation 记录用于撤销

**And** 关闭编辑面板

**And** 显示 Toast 提示"已更新 N 首"

**And** 保持选中状态（可继续操作）

---

### Story 4.2: 批量修改封面

As a 用户，
I want 批量修改选中歌曲的封面，
So that 统一一批歌曲的封面图片。

**Acceptance Criteria:**

**Given** 用户选中了 N 首歌曲
**When** 用户在批量编辑面板选择"更换封面"
**Then** 支持选择本地图片文件

**And** 支持拖拽上传封面

**And** 预览区域显示新封面效果

---

**Given** 用户确认封面
**When** 用户点击"应用"
**Then** 将新封面应用到所有选中歌曲

**And** 如果歌曲已有封面，覆盖原封面

**And** 创建 BatchOperation 记录用于撤销

**And** 显示 Toast 提示"已更新 N 首封面"

---

### Story 4.3: 搜索并批量应用歌词

As a 用户，
I want 搜索并批量应用歌词，
So that 快速为一批缺少歌词的歌曲补充歌词。

**Acceptance Criteria:**

**Given** 用户选中了 N 首无歌词的歌曲
**When** 用户点击"搜索歌词"
**Then** 使用文件名分隔的词片段作为搜索关键词

**And** 调用在线歌词搜索 API

**And** 显示搜索结果列表

---

**Given** 搜索结果返回
**When** 用户选择某个结果
**Then** 预览歌词内容

**And** 显示匹配度说明

**And** 用户确认后批量应用歌词到选中歌曲

---

**Given** 歌词搜索失败
**When** API 返回错误或无结果
**Then** 显示"搜索失败，请重试"

**And** 提供"重试"按钮

**And** 用户可手动输入歌词

---

### Story 4.4: 撤销批量编辑

As a 用户，
I want 撤销批量编辑操作，
So that 修正错误的批量修改。

**Acceptance Criteria:**

**Given** 用户执行过批量编辑
**When** 用户点击"撤销"
**Then** 恢复最近一次批量编辑之前的状态

**And** 包括标签修改、封面修改、歌词修改

**And** 重复撤销可逐步恢复更早的操作

---

**Given** 撤销执行
**When** 系统恢复历史状态
**Then** 从 BatchOperation 表获取上一次的 OldValues

**And** 将歌曲记录恢复为 OldValues

**And** 显示 Toast 提示"已撤销"

**And** 更新列表显示恢复后的状态

---

**Given** 无可撤销的操作
**When** 用户点击"撤销"
**Then** 按钮禁用或显示"无可撤销"

**And** 提示用户没有可撤销的操作

---
