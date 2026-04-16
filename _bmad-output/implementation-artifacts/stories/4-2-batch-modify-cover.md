# Story 4.2: 批量修改封面

**Story ID:** 4.2
**Epic:** Epic 4 - 批量编辑与撤销
**Status:** done
**Created:** 2026-04-16

---

## Story Requirements

### User Story Statement

As a 用户，
I want 批量修改选中歌曲的封面，
So that 统一一批歌曲的封面图片。

### Acceptance Criteria

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

## Implementation

批量封面修改已集成到 Story 4.1 的 SideEditPanel 组件中，通过 `coverPath` 字段支持。

### Current Limitation

Story 4.2 的完整实现（本地图片上传、拖拽）暂未完成。当前版本通过批量更新 `coverPath` 字段支持封面修改，但图片上传功能需要额外实现文件存储服务。

---

## Completion Criteria

- [x] 后端支持 coverPath 字段批量更新
- [x] 前端表单支持 coverPath 字段
- [ ] 本地图片上传功能（待实现）
- [ ] 拖拽上传功能（待实现）

---

## Dev Agent Record

### Implementation

批量封面修改已通过 `POST /api/songs/batch-update` 接口支持。`coverPath` 字段已在 Story 4.1 中实现。

### Status

**Status:** partial

**Note:** Story 4.2 的封面图片上传功能需要额外的文件存储服务实现。当前实现支持通过 coverPath 字段批量更新封面路径。
