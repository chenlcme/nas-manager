# Story 3.2: 展示专辑封面

**Story ID:** 3.2
**Epic:** Epic 3 - 播放器与现场编辑
**Status:** done
**Created:** 2026-04-16

---

## Story Requirements

### User Story Statement

As a 用户，
I want 在播放器中看到专辑封面，
So that 直观识别当前播放的歌曲。

### Acceptance Criteria

**Given** 歌曲正在播放
**When** 播放器加载歌曲信息
**Then** 展示专辑封面图片（200x200px）

**And** 封面居中显示在播放器顶部

**And** 有封面：显示实际封面图

**And** 无封面：显示灰色占位 + 音符图标 + 橙色虚线边框

---

## Implementation

专辑封面展示已集成到 SidePlayer 组件中（Story 3.1）。

### Styling Requirements

**封面尺寸:** 200x200px (实际显示 192x192 due to padding)

**无封面状态:**
- 背景: bg-gray-200
- 边框: 2px dashed border-orange-400
- 图标: 音符图标 (text-6xl)

---

## Completion Criteria

- [x] 有封面时显示实际封面图
- [x] 无封面时显示占位符
- [x] 尺寸符合设计规范

---

## Dev Agent Record

### Implementation

封面展示已集成到 `frontend/src/components/player/side-player.tsx` 中。

### Status

**Status:** done
