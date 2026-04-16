# Story 3.4: 展示播放时间

**Story ID:** 3.4
**Epic:** Epic 3 - 播放器与现场编辑
**Status:** done
**Created:** 2026-04-16

---

## Story Requirements

### User Story Statement

As a 用户，
I want 在播放器中看到播放时间，
So that 了解播放进度。

### Acceptance Criteria

**Given** 歌曲正在播放
**When** 播放器加载歌曲
**Then** 显示当前播放时间（格式：mm:ss）

**And** 显示总时长（格式：mm:ss）

**And** 显示播放进度条

**And** 进度条可拖动跳转播放位置

---

## Implementation

播放时间展示已集成到 SidePlayer 组件中（Story 3.1）。

### Technical Implementation

**HTML5 Audio Events:**
- `onTimeUpdate`: 更新当前播放时间
- `onLoadedMetadata`: 获取总时长

**进度条:**
- 使用 `<input type="range">` 实现
- 支持 onChange 事件跳转播放位置

**时间格式化:**
```typescript
function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}
```

---

## Completion Criteria

- [x] 当前播放时间显示
- [x] 总时长显示
- [x] 进度条显示
- [x] 进度条可拖动跳转

---

## Dev Agent Record

### Implementation

播放时间展示已集成到 `frontend/src/components/player/side-player.tsx` 中。

### Status

**Status:** done
