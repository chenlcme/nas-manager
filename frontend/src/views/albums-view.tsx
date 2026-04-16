import { useState, useEffect, Fragment } from 'preact/hooks';
import { AlbumWithCount, Song } from '../types/song';
import { SongTableRow } from '../components/song/song-table-row';

interface AlbumsViewProps {
  onPlaySong: (song: Song) => void;
}

export function AlbumsView({ onPlaySong }: AlbumsViewProps) {
  const [albums, setAlbums] = useState<AlbumWithCount[]>([]);
  const [expandedAlbum, setExpandedAlbum] = useState<number | null>(null);
  const [albumSongs, setAlbumSongs] = useState<Song[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [orderAsc, setOrderAsc] = useState(false);

  // 获取专辑列表
  useEffect(() => {
    setLoading(true);
    setError('');
    setExpandedAlbum(null); // 排序切换时清除展开状态
    setAlbumSongs([]);
    fetchAlbums();
  }, [orderAsc]);

  async function fetchAlbums() {
    try {
      const res = await fetch(`/api/albums?order=${orderAsc ? 'asc' : 'desc'}`);
      if (!res.ok) {
        throw new Error('获取专辑列表失败');
      }
      const data = await res.json();
      setAlbums(data.data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : '获取专辑列表失败');
    } finally {
      setLoading(false);
    }
  }

  // 展开专辑获取歌曲
  async function toggleAlbum(albumId: number) {
    if (loading) return; // 防止在加载时点击

    if (expandedAlbum === albumId) {
      setExpandedAlbum(null);
      setAlbumSongs([]);
      return;
    }

    setExpandedAlbum(albumId);
    setAlbumSongs([]); // 清除之前的歌曲列表
    try {
      const res = await fetch(`/api/albums/${albumId}/songs`);
      if (!res.ok) {
        throw new Error('获取歌曲列表失败');
      }
      const data = await res.json();
      setAlbumSongs(data.data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : '获取歌曲列表失败');
      setExpandedAlbum(null);
      setAlbumSongs([]);
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
          共 {albums.length} 张专辑
        </span>
        <button
          onClick={() => setOrderAsc(!orderAsc)}
          class="text-sm text-green-600 hover:text-green-700"
        >
          按专辑名{orderAsc ? '升序' : '降序'} ↕
        </button>
      </div>

      {/* 专辑列表 */}
      <div class="flex-1 overflow-auto">
        {albums.length === 0 ? (
          <div class="flex flex-col items-center justify-center py-12 text-gray-500">
            <div class="text-4xl mb-4">💿</div>
            <div class="text-lg">暂无专辑</div>
            <div class="text-sm">请先扫描音乐目录</div>
          </div>
        ) : (
          <table class="w-full">
            <thead class="bg-gray-50 sticky top-0">
              <tr>
                <th class="w-8 px-2 py-2"></th>
                <th class="px-2 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  专辑
                </th>
                <th class="px-2 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  艺术家
                </th>
                <th class="px-2 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  歌曲数量
                </th>
              </tr>
            </thead>
            <tbody>
              {albums.map((album) => (
                <Fragment key={album.id}>
                  {/* 专辑行 */}
                  <tr
                    onClick={() => toggleAlbum(album.id)}
                    class={{
                      'cursor-pointer hover:bg-gray-50 border-b border-gray-100': true,
                      'bg-gray-50': expandedAlbum === album.id,
                    }}
                  >
                    <td class="w-8 px-2 py-3">
                      <span class={{ 'text-gray-400': true, 'transform rotate-90': expandedAlbum === album.id }}>
                        ›
                      </span>
                    </td>
                    <td class="px-2 py-3 font-medium text-gray-900">
                      {album.name}
                    </td>
                    <td class="px-2 py-3 text-sm text-gray-500">
                      {album.artist}
                    </td>
                    <td class="px-2 py-3 text-sm text-gray-500">
                      {album.songCount} 首歌曲
                    </td>
                  </tr>

                  {/* 展开的歌曲列表 */}
                  {expandedAlbum === album.id && (
                    <tr key={`${album.id}-songs`}>
                      <td colSpan={4} class="bg-gray-50 px-4 py-2">
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
                                时长
                              </th>
                              <th class="w-12 px-2"></th>
                            </tr>
                          </thead>
                          <tbody>
                            {albumSongs.map((song) => (
                              <SongTableRow
                                key={song.id}
                                song={song}
                                onPlay={onPlaySong}
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
