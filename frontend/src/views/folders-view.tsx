import { useState, useEffect, useRef, Fragment } from 'preact/hooks';
import { FolderWithCount, Song } from '../types/song';
import { SongTableRow } from '../components/song/song-table-row';
import { SortSelector } from '../components/common/sort-selector';
import { DEFAULT_SORT_FIELD, DEFAULT_SORT_ORDER, SORT_BY_PARAM, ORDER_PARAM, SortField, SortOrder, REQUEST_TIMEOUT_MS } from '../constants/sort';

interface FoldersViewProps {
  onPlaySong: (song: Song) => void;
}

export function FoldersView({ onPlaySong }: FoldersViewProps) {
  const [folders, setFolders] = useState<FolderWithCount[]>([]);
  const [expandedFolder, setExpandedFolder] = useState<number | null>(null);
  const [folderSongs, setFolderSongs] = useState<Song[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [orderAsc, setOrderAsc] = useState(false);
  // 歌曲排序状态
  const [songSortBy, setSongSortBy] = useState<SortField>(DEFAULT_SORT_FIELD);
  const [songOrder, setSongOrder] = useState<SortOrder>(DEFAULT_SORT_ORDER);
  // 展开时显示 loading
  const [loadingSongs, setLoadingSongs] = useState(false);
  // 排序切换时显示 loading
  const [sorting, setSorting] = useState(false);
  // 用于取消旧请求
  const abortControllerRef = useRef<AbortController | null>(null);
  // 用于防止竞态的请求 ID
  const requestIdRef = useRef<number>(0);
  // 防抖定时器
  const debounceTimerRef = useRef<number | null>(null);

  // 获取文件夹列表
  useEffect(() => {
    setLoading(true);
    setError('');
    setExpandedFolder(null); // 排序切换时清除展开状态
    setFolderSongs([]);
    setSorting(false);
    fetchFolders();

    return () => {
      // 组件卸载时取消所有待处理的请求和定时器
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
      if (debounceTimerRef.current !== null) {
        clearTimeout(debounceTimerRef.current);
      }
    };
  }, [orderAsc]);

  async function fetchFolders() {
    try {
      const res = await fetch(`/api/folders?order=${orderAsc ? 'asc' : 'desc'}`);
      if (!res.ok) {
        throw new Error('获取文件夹列表失败');
      }
      const data = await res.json();
      setFolders(data.data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : '获取文件夹列表失败');
    } finally {
      setLoading(false);
    }
  }

  // 展开文件夹获取歌曲
  async function toggleFolder(folderId: number) {
    if (loading) return; // 防止在加载时点击

    if (expandedFolder === folderId) {
      setExpandedFolder(null);
      setFolderSongs([]);
      setLoadingSongs(false);
      setSorting(false);
      setError('');
      // 重置排序状态
      setSongSortBy(DEFAULT_SORT_FIELD);
      setSongOrder(DEFAULT_SORT_ORDER);
      return;
    }

    // 取消之前的请求
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
    // 清除防抖定时器
    if (debounceTimerRef.current !== null) {
      clearTimeout(debounceTimerRef.current);
      debounceTimerRef.current = null;
    }

    setExpandedFolder(folderId);
    setFolderSongs([]); // 清除之前的歌曲列表
    setLoadingSongs(true);
    setError('');

    const controller = new AbortController();
    abortControllerRef.current = controller;
    const currentRequestId = ++requestIdRef.current;

    try {
      const url = `/api/folders/${folderId}/songs?${SORT_BY_PARAM}=${encodeURIComponent(songSortBy)}&${ORDER_PARAM}=${encodeURIComponent(songOrder)}`;
      const res = await fetchWithTimeout(url, { signal: controller.signal }, REQUEST_TIMEOUT_MS);
      if (currentRequestId !== requestIdRef.current) return; // 请求已过时
      if (!res.ok) {
        throw new Error('获取歌曲列表失败');
      }
      const data = await res.json();
      setFolderSongs(data.data || []);
    } catch (err) {
      if (err instanceof Error && err.name === 'AbortError') {
        // 请求被取消，忽略
        return;
      }
      if (currentRequestId !== requestIdRef.current) return; // 请求已过时
      setError(err instanceof Error ? err.message : '获取歌曲列表失败');
      setExpandedFolder(null);
      setFolderSongs([]);
    } finally {
      if (currentRequestId === requestIdRef.current) {
        setLoadingSongs(false);
      }
    }
  }

  // 歌曲排序变化时重新获取歌曲（带防抖）
  function handleSongSortChange(newSortBy: SortField, newOrder: SortOrder) {
    setSongSortBy(newSortBy);
    setSongOrder(newOrder);
    if (expandedFolder !== null) {
      // 清除之前的防抖定时器
      if (debounceTimerRef.current !== null) {
        clearTimeout(debounceTimerRef.current);
      }

      debounceTimerRef.current = setTimeout(() => {
        // 取消之前的请求
        if (abortControllerRef.current) {
          abortControllerRef.current.abort();
        }

        setSorting(true);
        setError(''); // 清除之前的错误

        const controller = new AbortController();
        abortControllerRef.current = controller;
        const currentRequestId = ++requestIdRef.current;

        const url = `/api/folders/${expandedFolder}/songs?${SORT_BY_PARAM}=${encodeURIComponent(newSortBy)}&${ORDER_PARAM}=${encodeURIComponent(newOrder)}`;
        fetchWithTimeout(url, { signal: controller.signal }, REQUEST_TIMEOUT_MS)
          .then((res) => {
            if (currentRequestId !== requestIdRef.current) return null;
            if (!res.ok) {
              throw new Error('获取歌曲列表失败');
            }
            return res.json();
          })
          .then((data) => {
            if (data && currentRequestId === requestIdRef.current) {
              setFolderSongs(data.data || []);
              setError(''); // 成功获取后清除错误状态
            }
          })
          .catch((err) => {
            if (err instanceof Error && err.name === 'AbortError') {
              // 请求被取消，忽略
              return;
            }
            if (currentRequestId !== requestIdRef.current) return;
            setError(err instanceof Error ? err.message : '获取歌曲列表失败');
          })
          .finally(() => {
            if (currentRequestId === requestIdRef.current) {
              setSorting(false);
            }
          });
      }, 300); // 300ms 防抖
    }
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
      {/* 排序控制 */}
      <div class="px-4 py-3 flex items-center justify-between border-b border-gray-200">
        <span class="text-sm text-gray-600">
          共 {folders.length} 个文件夹
        </span>
        <button
          onClick={() => setOrderAsc(!orderAsc)}
          class="text-sm text-green-600 hover:text-green-700"
        >
          按路径{orderAsc ? '升序' : '降序'} ↕
        </button>
      </div>

      {/* 文件夹列表 */}
      <div class="flex-1 overflow-auto">
        {folders.length === 0 ? (
          <div class="flex flex-col items-center justify-center py-12 text-gray-500">
            <div class="text-4xl mb-4">📁</div>
            <div class="text-lg">暂无文件夹</div>
            <div class="text-sm">请先扫描音乐目录</div>
          </div>
        ) : (
          <table class="w-full">
            <thead class="bg-gray-50 sticky top-0">
              <tr>
                <th class="w-8 px-2 py-2"></th>
                <th class="px-2 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  文件夹路径
                </th>
                <th class="px-2 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  歌曲数量
                </th>
              </tr>
            </thead>
            <tbody>
              {folders.map((folder) => (
                <Fragment key={folder.id}>
                  {/* 文件夹行 */}
                  <tr
                    onClick={() => toggleFolder(folder.id)}
                    class={{
                      'cursor-pointer hover:bg-gray-50 border-b border-gray-100': true,
                      'bg-gray-50': expandedFolder === folder.id,
                    }}
                  >
                    <td class="w-8 px-2 py-3">
                      <span class={{ 'text-gray-400': true, 'transform rotate-90': expandedFolder === folder.id }}>
                        ›
                      </span>
                    </td>
                    <td class="px-2 py-3 font-medium text-gray-900">
                      {folder.path}
                    </td>
                    <td class="px-2 py-3 text-sm text-gray-500">
                      {folder.songCount} 首歌曲
                    </td>
                  </tr>

                  {/* 展开的歌曲列表 */}
                  {expandedFolder === folder.id && (
                    <tr key={`${folder.id}-songs`}>
                      <td colSpan={3} class="bg-gray-50 px-4 py-2">
                        {/* 排序控制 */}
                        <div class="flex items-center justify-end mb-2">
                          <SortSelector
                            sortBy={songSortBy}
                            order={songOrder}
                            onSortChange={handleSongSortChange}
                          />
                        </div>
                        {/* 展开/排序时的 loading 状态 */}
                        {(loadingSongs || sorting) && (
                          <div class="flex items-center justify-center py-4 text-gray-500">
                            {loadingSongs ? '加载中...' : '排序中...'}
                          </div>
                        )}
                        {!loadingSongs && !sorting && (
                          <table class="w-full">
                            <thead>
                              <tr>
                                <th class="w-8 px-2"></th>
                                <th class="w-12 px-2"></th>
                                <th class="px-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                  歌名
                                </th>
                                <th class="px-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                  艺术家
                                </th>
                                <th class="px-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                  专辑
                                </th>
                                <th class="px-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                  年份
                                </th>
                                <th class="px-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                  流派
                                </th>
                                <th class="px-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                  文件路径
                                </th>
                                <th class="px-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                  时长
                                </th>
                                <th class="w-12 px-2"></th>
                              </tr>
                            </thead>
                            <tbody>
                              {folderSongs.map((song) => (
                                <SongTableRow
                                  key={song.id}
                                  song={song}
                                  onPlay={onPlaySong}
                                  showPath={true}
                                />
                              ))}
                            </tbody>
                          </table>
                        )}
                      </td>
                    </tr>
                  )}
                </Fragment>
              ))}
            </tbody>
          </table>
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
