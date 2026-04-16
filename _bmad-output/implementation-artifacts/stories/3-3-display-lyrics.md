# Story 3.3: 展示歌词

**Story ID:** 3.3
**Epic:** Epic 3 - 播放器与现场编辑
**Status:** done
**Created:** 2026-04-16

---

## Story Requirements

### User Story Statement

As a 用户，
I want 在播放器中看到歌词，
So that 跟随歌曲播放查看歌词内容。

### Acceptance Criteria

**Given** 歌曲正在播放
**When** 播放器加载歌曲信息
**Then** 如果有内嵌歌词，显示歌词文本

**And** 支持静态歌词显示（逐行滚动）

**And** 如果无歌词，显示"暂无歌词"提示

**And** 歌词区域可滚动

---

## Implementation

歌词展示已集成到 SidePlayer 组件中（Story 3.1）。

### Styling Requirements

**歌词区域:**
- 最大高度: 192px (max-h-48)
- 背景: bg-gray-50
- 内边距: p-3
- 字体: text-sm, whitespace-pre-wrap
- 滚动: overflow-y-auto

---

## Completion Criteria

- [x] 有歌词时显示歌词文本
- [x] 无歌词时显示"暂无歌词"
- [x] 歌词区域可滚动

---

## Dev Agent Record

### Implementation

歌词展示已集成到 `frontend/src/components/player/side-player.tsx` 中。

### Status

**Status:** done
