import { useState, useEffect } from 'preact/hooks';
import { useSelection } from '../../contexts/selection-context';
import { Song } from '../../types/song';

interface SideEditPanelProps {
  onClose: () => void;
  onSuccess: () => void;
  onError?: (message: string) => void;
  onUndoAvailable?: (hasUndo: boolean) => void;
}

interface EditForm {
  artist: string;
  album: string;
  year: string;
  genre: string;
}

export function SideEditPanel({ onClose, onSuccess, onError, onUndoAvailable }: SideEditPanelProps) {
  const { selected, clear, count } = useSelection();
  const [editForm, setEditForm] = useState<EditForm>({
    artist: '',
    album: '',
    year: '',
    genre: '',
  });
  const [isSaving, setIsSaving] = useState(false);
  const [activeTab, setActiveTab] = useState<'tags' | 'lyrics'>('tags');
  const [hasUndo, setHasUndo] = useState(false);
  const [selectedSongs, setSelectedSongs] = useState<Song[]>([]);
  const [loadingSongs, setLoadingSongs] = useState(false);

  // Fetch selected songs for preview
  useEffect(() => {
    if (count > 0) {
      fetchSelectedSongs();
    }
  }, [count]);

  const fetchSelectedSongs = async () => {
    setLoadingSongs(true);
    try {
      const ids = [...selected].slice(0, 10); // Limit to first 10 for preview
      const res = await fetch('/api/songs/batch-get', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ ids }),
      });
      if (res.ok) {
        const data = await res.json();
        setSelectedSongs(data.data || []);
      }
    } catch {
      // Ignore fetch errors for preview
    } finally {
      setLoadingSongs(false);
    }
  };

  // Check for undo availability
  useEffect(() => {
    checkUndoAvailability();
  }, []);

  const checkUndoAvailability = async () => {
    try {
      const res = await fetch('/api/batches/latest');
      const hasUndo = res.ok;
      setHasUndo(hasUndo);
      onUndoAvailable?.(hasUndo);
    } catch {
      setHasUndo(false);
      onUndoAvailable?.(false);
    }
  };

  const handleSave = async () => {
    if (count === 0) {
      onError?.('请先选择要编辑的歌曲');
      return;
    }
    setIsSaving(true);

    try {
      const ids = [...selected];
      const updates: Record<string, unknown> = {};

      if (editForm.artist) updates.artist = editForm.artist;
      if (editForm.album) updates.album = editForm.album;
      if (editForm.year) {
        const yearNum = parseInt(editForm.year);
        if (isNaN(yearNum)) {
          onError?.('年份必须是有效数字');
          setIsSaving(false);
          return;
        }
        updates.year = yearNum;
      }
      if (editForm.genre) updates.genre = editForm.genre;

      if (Object.keys(updates).length === 0) {
        onError?.('请至少填写一个要修改的字段');
        setIsSaving(false);
        return;
      }

      const res = await fetch('/api/songs/batch-update', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          ids,
          ...updates,
        }),
      });

      if (!res.ok) {
        throw new Error('批量更新失败');
      }

      const result = await res.json();
      if (result.data.failed > 0) {
        onError?.(`部分更新失败: 成功 ${result.data.succeeded}, 失败 ${result.data.failed}`);
      } else {
        onSuccess();
        clear();
        onClose();
      }
    } catch (err) {
      onError?.(err instanceof Error ? err.message : '批量更新失败');
    } finally {
      setIsSaving(false);
    }
  };

  const handleUndo = async () => {
    try {
      const res = await fetch('/api/batches/latest');
      if (!res.ok) {
        onError?.('没有可撤销的操作');
        return;
      }

      const batch = await res.json();
      if (!batch.data?.id) {
        onError?.('没有可撤销的操作');
        return;
      }

      const undoRes = await fetch(`/api/songs/undo/${batch.data.id}`, {
        method: 'POST',
      });

      if (!undoRes.ok) {
        throw new Error('撤销失败');
      }

      const result = await undoRes.json();
      onSuccess();
      onError?.(`已撤销: 成功 ${result.data.succeeded}, 失败 ${result.data.failed}`);
      setHasUndo(false);
      checkUndoAvailability();
    } catch (err) {
      onError?.(err instanceof Error ? err.message : '撤销失败');
    }
  };

  return (
    <div class="fixed right-0 top-0 bottom-0 w-80 bg-white shadow-xl z-50 flex flex-col">
      {/* Header */}
      <div class="flex items-center justify-between px-4 py-3 border-b border-gray-200">
        <h2 class="text-lg font-semibold text-gray-900">批量编辑</h2>
        <button
          onClick={onClose}
          class="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg"
          title="关闭"
        >
          ✕
        </button>
      </div>

      {/* Content */}
      <div class="flex-1 overflow-y-auto p-4">
        {/* Selected count and preview */}
        <div class="mb-4 p-3 bg-green-50 rounded-lg">
          <span class="text-green-700 font-medium">已选中 {count} 首歌曲</span>
          {count > 0 && (
            <div class="mt-2 text-sm text-green-600">
              {loadingSongs ? (
                <span>加载预览...</span>
              ) : (
                <div class="space-y-1">
                  {selectedSongs.slice(0, 5).map((song) => (
                    <div key={song.id} class="truncate">
                      {song.title || '未知歌曲'} - {song.artist || '未知艺术家'}
                    </div>
                  ))}
                  {count > 5 && (
                    <div class="text-green-500">等 {count - 5} 首...</div>
                  )}
                </div>
              )}
            </div>
          )}
        </div>

        {/* Tabs */}
        <div class="flex border-b border-gray-200 mb-4">
          <button
            onClick={() => setActiveTab('tags')}
            class={`px-4 py-2 text-sm font-medium border-b-2 ${
              activeTab === 'tags'
                ? 'border-green-500 text-green-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            标签信息
          </button>
          <button
            onClick={() => setActiveTab('lyrics')}
            class={`px-4 py-2 text-sm font-medium border-b-2 ${
              activeTab === 'lyrics'
                ? 'border-green-500 text-green-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            歌词
          </button>
        </div>

        {activeTab === 'tags' && (
          <div class="space-y-4">
            <p class="text-sm text-gray-500">留空的字段将保持不变</p>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">艺术家</label>
              <input
                type="text"
                value={editForm.artist}
                onChange={(e) => setEditForm({ ...editForm, artist: (e.target as HTMLInputElement).value })}
                placeholder="保持不变"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">专辑</label>
              <input
                type="text"
                value={editForm.album}
                onChange={(e) => setEditForm({ ...editForm, album: (e.target as HTMLInputElement).value })}
                placeholder="保持不变"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
              />
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">年份</label>
                <input
                  type="text"
                  value={editForm.year}
                  onChange={(e) => setEditForm({ ...editForm, year: (e.target as HTMLInputElement).value })}
                  placeholder="保持不变"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">流派</label>
                <input
                  type="text"
                  value={editForm.genre}
                  onChange={(e) => setEditForm({ ...editForm, genre: (e.target as HTMLInputElement).value })}
                  placeholder="保持不变"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
                />
              </div>
            </div>
          </div>
        )}

        {activeTab === 'lyrics' && (
          <div class="space-y-4">
            <p class="text-sm text-gray-500">
              批量歌词功能需要在线搜索歌词服务，暂未实现。
            </p>
            <p class="text-sm text-gray-500">
              您可以在播放器中单独编辑每首歌曲的歌词。
            </p>
          </div>
        )}

        {/* Preview */}
        {(editForm.artist || editForm.album || editForm.year || editForm.genre) && (
          <div class="mt-6 p-3 bg-gray-50 rounded-lg">
            <h4 class="text-sm font-medium text-gray-700 mb-2">预览更改</h4>
            <ul class="text-sm text-gray-600 space-y-1">
              {editForm.artist && <li>艺术家 → {editForm.artist}</li>}
              {editForm.album && <li>专辑 → {editForm.album}</li>}
              {editForm.year && <li>年份 → {editForm.year}</li>}
              {editForm.genre && <li>流派 → {editForm.genre}</li>}
            </ul>
          </div>
        )}
      </div>

      {/* Actions */}
      <div class="p-4 border-t border-gray-200 space-y-2">
        <button
          onClick={handleSave}
          disabled={isSaving || count === 0}
          class="w-full py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 font-medium disabled:opacity-50"
        >
          {isSaving ? '保存中...' : '应用更改'}
        </button>
        <button
          onClick={handleUndo}
          disabled={!hasUndo}
          class="w-full py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 font-medium disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {!hasUndo ? '无可撤销' : '撤销上次的批量编辑'}
        </button>
      </div>
    </div>
  );
}