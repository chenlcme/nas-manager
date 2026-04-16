import { useState, useEffect, Fragment } from 'preact/hooks';
import { FolderWithCount, Song } from '../types/song';
import { SongTableRow } from '../components/song/song-table-row';

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

  // 获取文件夹列表
  useEffect(() => {
    setLoading(true);
    setError('');
    setExpandedFolder(null); // 排序切换时清除展开状态
    setFolderSongs([]);
    fetchFolders();
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
      return;
    }

    setExpandedFolder(folderId);
    setFolderSongs([]); // 清除之前的歌曲列表
    try {
      const res = await fetch(`/api/folders/${folderId}/songs`);
      if (!res.ok) {
        throw new Error('获取歌曲列表失败');
      }
      const data = await res.json();
      setFolderSongs(data.data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : '获取歌曲列表失败');
      setExpandedFolder(null);
      setFolderSongs([]);
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
