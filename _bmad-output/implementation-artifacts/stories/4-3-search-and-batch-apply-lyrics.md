# Story 4.3: 搜索并批量应用歌词

**Story ID:** 4.3
**Epic:** Epic 4 - 批量编辑与撤销
**Status:** partial
**Created:** 2026-04-16

---

## Story Requirements

### User Story Statement

As a 用户，
I want 搜索并批量应用歌词，
So that 快速为一批缺少歌词的歌曲补充歌词。

### Acceptance Criteria

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

## Implementation

Story 4.3 需要在线歌词搜索服务，当前版本暂未实现此功能。

### UI Placeholder

在 SideEditPanel 的"歌词"Tab 中，显示提示信息：
- "批量歌词功能需要在线搜索歌词服务，暂未实现。"
- "您可以在播放器中单独编辑每首歌曲的歌词。"

### Future Implementation

需要实现：
1. 歌词搜索 API 服务（如 QQ 音乐、网易云音乐等）
2. 前端歌词搜索 UI
3. 歌词预览和确认流程

---

## Completion Criteria

- [x] UI 占位符显示
- [ ] 歌词搜索 API 集成
- [ ] 歌词预览功能
- [ ] 批量应用歌词

---

## Dev Agent Record

### Implementation

当前实现是 UI 占位符，提示用户此功能暂未实现。

### Status

**Status:** partial

**Note:** Story 4.3 需要第三方歌词搜索 API 服务，当前版本显示提示信息。
