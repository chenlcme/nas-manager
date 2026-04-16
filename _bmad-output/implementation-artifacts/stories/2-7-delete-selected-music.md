# Story 2.7: 删除选中音乐

**Story ID:** 2.7
**Epic:** Epic 2 - 音乐库浏览与搜索
**Status:** ready-for-dev
**Created:** 2026-04-16

---

## Story Requirements

### User Story Statement

As a 用户，
I want 删除选中的音乐文件，
So that 清理不需要的歌曲（同时删除文件和数据库记录）。

### Acceptance Criteria

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

## Technical Requirements

### Backend API Design

**Batch Delete Endpoint:**

```
POST /api/songs/delete
Content-Type: application/json

Request Body:
{
  "ids": [1, 2, 3]  // 要删除的歌曲 ID 数组
}

Response (成功):
{
  "data": {
    "total": 3,
    "succeeded": 2,
    "failed": 1,
    "results": [
      { "id": 1, "file_path": "/music/song1.mp3", "status": "deleted" },
      { "id": 2, "file_path": "/music/song2.mp3", "status": "deleted" },
      { "id": 3, "file_path": "/music/song3.mp3", "status": "failed", "error": "permission denied" }
    ]
  }
}
```

**Error Codes:**

| 错误码 | 说明 |
|--------|------|
| VALIDATION_ERROR | 参数验证失败（空数组） |
| DELETE_FILE_FAILED | 文件删除失败 |
| DELETE_DB_FAILED | 数据库记录删除失败 |

### Backend Components

**Handler Layer (`internal/handler/song.go`):**
- 添加 `DeleteSongs` 方法处理 `POST /api/songs/delete`

**Service Layer (`internal/service/song.go`):**
- 添加 `DeleteSongs` 方法处理批量删除逻辑
- 返回详细的删除结果报告

**Repository Layer (`internal/repository/song.go`):**
- 已有 `Delete(id uint)` 方法
- 考虑添加 `DeleteBatch(ids []uint)` 批量删除方法

### Frontend Components

**SelectionBar 增强 (`frontend/src/components/common/selection-bar.tsx`):**
- 添加 `onDelete` prop
- 添加"删除"按钮（红色风格）
- 点击后调用确认对话框

**Delete Confirmation Dialog (`frontend/src/components/common/delete-confirm-dialog.tsx`):**
- 新组件
- 显示删除数量 N
- 警告文本："文件和数据库记录都将被删除"
- 确认/取消按钮

**API Call (`frontend/src/utils/api.ts`):**
- 添加 `deleteSongs(ids: number[])` 函数
- 调用 `POST /api/songs/delete`

**Views 更新:**
- ArtistsView, AlbumsView, FoldersView 传递 `onDelete` 到 SelectionBar

### File Structure

```
frontend/src/
├── components/
│   ├── common/
│   │   ├── selection-bar.tsx    # 修改 - 添加删除按钮
│   │   └── delete-confirm-dialog.tsx  # 新增
│   └── ...
├── utils/
│   └── api.ts                  # 修改 - 添加 deleteSongs
├── views/
│   ├── artists-view.tsx        # 修改 - 传递 onDelete
│   ├── albums-view.tsx         # 修改 - 传递 onDelete
│   └── folders-view.tsx        # 修改 - 传递 onDelete
```

### Styling Requirements

**Delete Button (SelectionBar):**
- 背景: bg-red-500
- Hover: bg-red-600
- 文字: 白色

**Delete Confirm Dialog:**
- 警告图标: 橙色
- 确认按钮: bg-red-500
- 取消按钮: bg-gray-300

### State Management

- 使用现有的 `SelectionContext` 获取选中歌曲 ID
- 删除成功后调用 `clear()` 清除选中状态
- 删除成功后刷新当前视图列表

### UX Requirements

1. **确认对话框:**
   - 模态对话框，居中显示
   - 标题："确认删除"
   - 内容："将删除 N 首歌曲，文件和数据库记录都将被删除"
   - 操作不可逆提示（红色警告文字）

2. **删除进度:**
   - 显示整体进度（成功/失败数量）
   - 完成后显示 Toast 提示

3. **失败处理:**
   - 显示失败文件列表及原因
   - 提供"关闭"按钮

---

## Implementation Notes

### Previous Story Learnings

Story 2.6 已实现：
- SelectionContext 管理选中状态
- SelectionBar 组件固定在列表上方
- 视图组件正确集成 SelectionBar

Story 2.7 需要在 SelectionBar 添加删除按钮，并实现后端批量删除 API。

### Key Patterns

1. **API 调用模式:**
   ```typescript
   // frontend/src/utils/api.ts
   export async function deleteSongs(ids: number[]) {
     const res = await fetch('/api/songs/delete', {
       method: 'POST',
       headers: { 'Content-Type': 'application/json' },
       body: JSON.stringify({ ids }),
     });
     if (!res.ok) throw new Error('Delete failed');
     return res.json();
   }
   ```

2. **删除结果处理:**
   ```typescript
   const handleDelete = async () => {
     try {
       const result = await deleteSongs([...selected]);
       if (result.data.failed > 0) {
         // 显示失败列表
       } else {
         // 显示成功 Toast
         clear();
       }
     } catch (e) {
       // 显示错误 Toast
     }
   };
   ```

### Backend Delete Logic

```go
// Service Layer - 删除逻辑
func (s *SongService) DeleteSongs(ids []uint) (*DeleteResult, error) {
    result := &DeleteResult{
        Results: make([]SongDeleteResult, 0, len(ids)),
    }

    for _, id := range ids {
        song, err := s.songRepo.GetByID(id)
        if err != nil {
            result.Results = append(result.Results, SongDeleteResult{
                ID: id, Status: "failed", Error: "song not found",
            })
            result.Failed++
            continue
        }

        // 删除文件
        if err := os.Remove(song.FilePath); err != nil {
            result.Results = append(result.Results, SongDeleteResult{
                ID: id, FilePath: song.FilePath, Status: "failed", Error: err.Error(),
            })
            result.Failed++
            continue
        }

        // 删除数据库记录
        if err := s.songRepo.Delete(id); err != nil {
            // 文件已删除，记录警告但继续
            result.Results = append(result.Results, SongDeleteResult{
                ID: id, FilePath: song.FilePath, Status: "failed", Error: "db delete failed: " + err.Error(),
            })
            result.Failed++
            continue
        }

        result.Results = append(result.Results, SongDeleteResult{
            ID: id, FilePath: song.FilePath, Status: "deleted",
        })
        result.Succeeded++
    }

    result.Total = len(ids)
    return result, nil
}
```

---

## Dependencies

- Story 2.6 (多选歌曲) 已完成
- Story 2.5 (查看单曲详情) 已完成
- 后端三层架构已建立

---

## Testing Requirements

1. 选中 N 首歌曲，点击删除按钮显示确认对话框
2. 确认删除后显示删除进度
3. 删除成功：显示成功 Toast，清除选中，刷新列表
4. 删除失败：显示失败文件列表及原因
5. 取消删除：关闭对话框，选中状态保持
6. 权限不足场景：文件无法删除但数据库记录删除成功

---

## Completion Criteria

- [x] 后端 `POST /api/songs/delete` 端点实现
- [x] 后端批量删除逻辑（文件 + 数据库）
- [x] 前端删除确认对话框组件
- [x] SelectionBar 添加删除按钮
- [x] API 调用和结果处理
- [x] 失败文件列表展示
- [x] 视图正确传递 onDelete prop
- [x] 单元测试（Handler 层）

---

## Dev Agent Record

### Implementation Plan

Story 2.7 删除选中音乐实现说明：

**后端实现:**
1. 在 `internal/handler/song.go` 添加 `DeleteSongs` 方法
2. 添加 `DeleteResult` 和 `SongDeleteResult` 结构体
3. 在 `cmd/server/main.go` 注册 `POST /api/songs/delete` 路由

**前端实现:**
1. 创建 `DeleteConfirmDialog` 组件 - 确认删除对话框
2. 创建 `DeleteResultDialog` 组件 - 删除结果展示
3. 更新 `SelectionBar` 组件 - 添加删除按钮和处理逻辑
4. 更新 `ArtistsView`, `AlbumsView`, `FoldersView` - 添加 onDeleteSuccess 回调

**测试:**
1. `TestDeleteSongs_Success` - 验证成功删除
2. `TestDeleteSongs_EmptyIDs` - 验证空数组校验
3. `TestDeleteSongs_PartialFailure` - 验证部分失败场景

### Files Modified/Created

- `internal/handler/song.go` - 添加 DeleteSongs 方法
- `internal/handler/song_test.go` - 添加删除功能测试
- `cmd/server/main.go` - 注册删除路由
- `frontend/src/components/common/selection-bar.tsx` - 添加删除功能
- `frontend/src/components/common/delete-confirm-dialog.tsx` - 新增
- `frontend/src/components/common/delete-result-dialog.tsx` - 新增
- `frontend/src/views/artists-view.tsx` - 添加 onDeleteSuccess
- `frontend/src/views/albums-view.tsx` - 添加 onDeleteSuccess
- `frontend/src/views/folders-view.tsx` - 添加 onDeleteSuccess

### Status

**Status:** review
