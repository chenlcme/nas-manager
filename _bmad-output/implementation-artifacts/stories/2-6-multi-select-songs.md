# Story 2.6: 多选歌曲文件

**Story ID:** 2.6
**Epic:** Epic 2 - 音乐库浏览与搜索
**Status:** ready-for-dev
**Created:** 2026-04-16

---

## Story Requirements

### User Story Statement

As a 用户，
I want 在列表中多选多个音乐文件，
So that 对批量歌曲进行操作。

### Acceptance Criteria

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

## Technical Requirements

### Frontend Architecture

**State Management:** Preact Context API
- SelectionContext 管理选中状态
- Context 结构：
  ```typescript
  interface SelectionContextType {
    selected: Set<number>;  // 选中歌曲的 ID 集合
    toggle: (id: number) => void;      // 切换单个歌曲选中状态
    selectAll: (ids: number[]) => void; // 全选当前视图所有歌曲
    clear: () => void;                 // 清除所有选中状态
    isSelected: (id: number) => boolean; // 检查歌曲是否被选中
    count: number;                      // 选中歌曲数量
  }
  ```

**Components to Implement/Modify:**

1. **SelectionBar 组件** (新组件)
   - 路径: `frontend/src/components/selection/selection-bar.tsx`
   - 位置: 固定在列表上方
   - 功能:
     - 显示已选中歌曲数量
     - "全选" 按钮
     - "取消全选" 按钮
     - 批量编辑按钮（后续 Story 2.7 使用）
   - 样式: 浅绿色背景 (#F0FDF4)，主题绿色边框

2. **SongTableRow 组件** (修改现有)
   - 路径: `frontend/src/components/song/song-table.tsx`
   - 添加复选框列（32px 宽度）
   - 复选框交互: 点击切换选中状态
   - 选中状态样式: 主题色背景 (#F0FDF4)
   - 无封面状态: 灰色占位 + 音符图标 + 橙色虚线边框

3. **SongCheckbox 组件** (可能需要)
   - 路径: `frontend/src/components/song/song-checkbox.tsx`
   - 纯复选框组件，可复用

### File Structure

```
frontend/src/
├── components/
│   ├── selection/
│   │   └── selection-bar.tsx    # 新增
│   ├── song/
│   │   ├── song-table.tsx       # 修改 - 添加复选框
│   │   └── song-checkbox.tsx    # 可选复用
├── contexts/
│   └── selection-context.tsx     # 新增 - 选中状态管理
├── views/
│   ├── artists-view.tsx         # 修改 - 集成 SelectionBar
│   ├── albums-view.tsx          # 修改 - 集成 SelectionBar
│   └── folders-view.tsx         # 修改 - 集成 SelectionBar
```

### Styling Requirements

**色彩系统:**
- 主题绿: #22C55E
- 浅绿背景: #F0FDF4
- 深色文字: #1E293B
- 次要文字: #64748B

**选中状态:**
- 背景: #F0FDF4 (浅绿色)
- 复选框: 主题绿色勾选

**SelectionBar 样式:**
- 固定顶部定位
- 背景: #F0FDF4
- 边框-bottom: 1px solid #E2E8F0
- 高度: 48px
- 内边距: 0 16px

### UX Requirements

1. **交互反馈:**
   - 点击复选框立即响应（≤200ms）
   - 选中/取消选中时有视觉反馈

2. **SelectionBar 行为:**
   - 有选中时显示，数量 > 0
   - 无选中时隐藏或显示数量为 0

3. **全选逻辑:**
   - 只选中当前视图可见的歌曲
   - 不同视图（歌手/专辑/文件夹）有各自的全选范围

### Dependencies

- Story 2.1, 2.2, 2.3 已完成（歌手/专辑/文件夹视图）
- Story 2.5 已完成（查看单曲详情）
- Story 2.7 依赖本 Story（删除选中音乐）

### Testing Requirements

1. 复选框点击切换选中状态
2. 选中状态视觉高亮正确显示
3. SelectionBar 数量显示正确
4. 全选按钮选中当前视图所有歌曲
5. 取消全选清除所有选中状态
6. 不同视图切换后选中状态独立

---

## Implementation Notes

### Previous Story Learnings

从最近的 commit 历史分析：
- 0a8fa01: feat(frontend): 集成 FoldersView 到 App 组件
- 486b1fa: feat(handler,repo,frontend): 实现文件夹视图浏览功能
- f635a85: feat(repo,handler): 实现专辑视图浏览功能
- 7bf12f9: Implement Story 2.1: Artist View Browse

前端使用 Preact + Tailwind CSS，组件化开发模式已建立。视图组件（ArtistsView, AlbumsView, FoldersView）已实现列表展示，需要在现有基础上添加多选功能。

### Key Patterns

1. **Context 模式:** 使用 Preact createContext 管理全局选中状态
2. **组件组合:** SongTableRow 内含 Checkbox，SelectionBar 在视图顶层
3. **状态驱动:** UI 状态完全由 SelectionContext 驱动

---

## Completion Criteria

- [x] SelectionContext 实现完成
- [x] SongTableRow 添加复选框和选中样式
- [x] SelectionBar 组件实现，显示选中数量
- [x] 全选/取消全选功能正常
- [x] 不同视图正确集成多选功能
- [x] 视觉样式符合设计规范

---

## Dev Agent Record

### Implementation Plan

Story 2.6 多选歌曲文件实现说明：

1. **SelectionContext** (`frontend/src/contexts/selection-context.tsx`)
   - 已是完整实现，包含 toggle, selectAll, clear, isSelected, count

2. **SongTableRow** (`frontend/src/components/song/song-table-row.tsx`)
   - 已有复选框实现，使用 useSelection hook
   - 选中状态样式: bg-green-50

3. **SelectionBar** (`frontend/src/components/common/selection-bar.tsx`)
   - 修复定位: 从 `fixed bottom-0` 改为 `sticky top-0`
   - 修复背景色: 从 `bg-white` 改为 `bg-green-50`
   - 修复全选逻辑: 使用实际的歌曲 ID 数组
   - 添加 onSelectAll 回调 prop

4. **视图集成**
   - ArtistsView: 添加 SelectionBar 到展开的艺术家歌曲列表
   - AlbumsView: 添加 SelectionBar 到展开的专辑歌曲列表
   - FoldersView: 添加 SelectionBar 到展开的文件夹歌曲列表
   - App.tsx: 移除多余的 App 级 SelectionBar，传递 onBatchEdit prop

### Files Modified

- `frontend/src/components/common/selection-bar.tsx` - 修复定位和样式
- `frontend/src/views/artists-view.tsx` - 集成 SelectionBar
- `frontend/src/views/albums-view.tsx` - 集成 SelectionBar
- `frontend/src/views/folders-view.tsx` - 集成 SelectionBar
- `frontend/src/app.tsx` - 移除多余 SelectionBar，添加 onBatchEdit prop

### Status

**Status:** review
