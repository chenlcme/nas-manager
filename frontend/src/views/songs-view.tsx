import { useState, useEffect, useRef } from 'preact/hooks';
import { Song } from '../types/song';
import { useSelection } from '../contexts/selection-context';
import { SongTableRow } from '../components/song/song-table-row';
import { SortSelector } from '../components/common/sort-selector';
import { SelectionBar } from '../components/common/selection-bar';
import { DEFAULT_SORT_FIELD, DEFAULT_SORT_ORDER, SORT_BY_PARAM, ORDER_PARAM, FOLDER_PARAM, SortField, SortOrder, REQUEST_TIMEOUT_MS } from '../constants/sort';

interface SongsViewProps {
  onPlaySong: (song: Song) => void;
  onShowSongDetail: (song: Song) => void;
  onBatchEdit?: () => void;
  playingSongId?: number | null;
}

export function SongsView({ onPlaySong, onShowSongDetail, onBatchEdit, playingSongId }: SongsViewProps) {
  const { selectAll } = useSelection();
  const [songs, setSongs] = useState<Song[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [sortBy, setSortBy] = useState<SortField>(DEFAULT_SORT_FIELD);
  const [order, setOrder] = useState<SortOrder>(DEFAULT_SORT_ORDER);
  const [folderFilter, setFolderFilter] = useState<string | null>(null);
  const [sorting, setSorting] = useState(false);
  const abortControllerRef = useRef<AbortController | null>(null);
  const requestIdRef = useRef<number>(0);

  // 删除成功后刷新
  const handleDeleteSuccess = () => {
    fetchSongs();
  };

  // 获取所有歌曲
  useEffect(() => {
    fetchSongs();

    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
    };
  }, [sortBy, order, folderFilter]);

  async function fetchSongs() {
    // 取消之前的请求
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }

    setLoading(true);
    setError('');
    setSorting(false);

    const controller = new AbortController();
    abortControllerRef.current = controller;
    const currentRequestId = ++requestIdRef.current;

    try {
      let url = `/api/songs?${SORT_BY_PARAM}=${encodeURIComponent(sortBy)}&${ORDER_PARAM}=${encodeURIComponent(order)}`;
      if (folderFilter) {
        url += `&${FOLDER_PARAM}=${encodeURIComponent(folderFilter)}`;
      }
      const res = await fetchWithTimeout(url, { signal: controller.signal }, REQUEST_TIMEOUT_MS);
      if (currentRequestId !== requestIdRef.current) return;
      if (!res.ok) {
        throw new Error('获取歌曲列表失败');
      }
      const data = await res.json();
      if (currentRequestId === requestIdRef.current) {
        setSongs(data.data || []);
      }
    } catch (err) {
      if (err instanceof Error && err.name === 'AbortError') {
        return;
      }
      if (currentRequestId !== requestIdRef.current) return;
      setError(err instanceof Error ? err.message : '获取歌曲列表失败');
    } finally {
      if (currentRequestId === requestIdRef.current) {
        setLoading(false);
      }
    }
  }

  // 排序变化时重新获取
  function handleSortChange(newSortBy: SortField, newOrder: SortOrder) {
    setSortBy(newSortBy);
    setOrder(newOrder);
  }

  // 文件夹筛选
  function handleFolderFilter(folder: string) {
    setFolderFilter(folder);
  }

  // 清除筛选
  function handleClearFilter() {
    setFolderFilter(null);
  }

  if (loading) {
    return (
      <div class="flex items-center justify-center py-12">
        <div class="text-gray-500">加载中...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div class="flex items-center justify-center py-12">
        <div class="text-red-500">{error}</div>
      </div>
    );
  }

  return (
    <div class="flex flex-col h-full">
      {/* 筛选状态提示 */}
      {folderFilter && (
        <div class="px-4 py-2 bg-green-50 border-b border-green-200 flex items-center justify-between">
          <span class="text-sm text-green-700">
            当前筛选：<span class="font-mono">{folderFilter}</span>
          </span>
          <button
            onClick={handleClearFilter}
            class="text-sm text-green-600 hover:text-green-800 underline"
          >
            清除筛选
          </button>
        </div>
      )}

      {/* 排序控制 */}
      <div class="px-4 py-3 flex items-center justify-between border-b border-gray-200">
        <span class="text-sm text-gray-600">
          共 {songs.length} 首歌曲
        </span>
        <SortSelector
          sortBy={sortBy}
          order={order}
          onSortChange={handleSortChange}
        />
      </div>

      {/* 歌曲列表 */}
      <div class="flex-1 overflow-auto">
        {songs.length === 0 ? (
          <div class="flex flex-col items-center justify-center py-12 text-gray-500">
            <div class="text-4xl mb-4">🎵</div>
            <div class="text-lg">暂无歌曲</div>
            <div class="text-sm">请先扫描音乐目录</div>
          </div>
        ) : (
          <>
            <SelectionBar
              totalCount={songs.length}
              onBatchEdit={onBatchEdit}
              onSelectAll={() => selectAll(songs.map(s => s.id))}
              onDeleteSuccess={handleDeleteSuccess}
            />
            <table class="w-full">
              <thead class="bg-gray-50 sticky top-0">
                <tr>
                  <th class="w-6 px-1 py-2"></th>
                  <th class="w-10 px-1 py-2"></th>
                  <th class="px-2 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    歌名
                  </th>
                  <th class="px-2 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    艺术家
                  </th>
                  <th class="px-2 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    专辑
                  </th>
                  <th class="px-2 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-14">
                    年份
                  </th>
                  <th class="px-2 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-16">
                    流派
                  </th>
                  <th class="px-2 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    所属目录
                  </th>
                  <th class="px-2 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-12">
                    时长
                  </th>
                  <th class="w-16 px-2 py-2"></th>
                </tr>
              </thead>
              <tbody>
                {songs.map((song) => (
                  <SongTableRow
                    key={song.id}
                    song={song}
                    onPlay={onPlaySong}
                    onShowDetail={onShowSongDetail}
                    onFolderClick={handleFolderFilter}
                    showDir={true}
                    playingSongId={playingSongId}
                  />
                ))}
              </tbody>
            </table>
          </>
        )}
      </div>
    </div>
  );
}

// 带超时的 fetch 辅助函数
async function fetchWithTimeout(url: string, options: RequestInit, timeoutMs: number): Promise<Response> {
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), timeoutMs);

  try {
    const response = await fetch(url, {
      ...options,
      signal: controller.signal,
    });
    return response;
  } finally {
    clearTimeout(timeoutId);
  }
}