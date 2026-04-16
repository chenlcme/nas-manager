import { useState, useEffect } from 'preact/hooks';
import { ArtistWithCount, Song } from '../types/song';
import { useSelection } from '../contexts/selection-context';
import { SongTableRow } from '../components/song/song-table-row';

interface ArtistsViewProps {
  onPlaySong: (song: Song) => void;
}

export function ArtistsView({ onPlaySong }: ArtistsViewProps) {
  const [artists, setArtists] = useState<ArtistWithCount[]>([]);
  const [expandedArtist, setExpandedArtist] = useState<number | null>(null);
  const [artistSongs, setArtistSongs] = useState<Song[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [orderAsc, setOrderAsc] = useState(false);

  // 获取艺术家列表
  useEffect(() => {
    fetchArtists();
  }, [orderAsc]);

  async function fetchArtists() {
    setLoading(true);
    setError('');
    try {
      const res = await fetch(`/api/artists?order=${orderAsc ? 'asc' : 'desc'}`);
      if (!res.ok) {
        throw new Error('获取艺术家列表失败');
      }
      const data = await res.json();
      setArtists(data.data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : '获取艺术家列表失败');
    } finally {
      setLoading(false);
    }
  }

  // 展开艺术家获取歌曲
  async function toggleArtist(artistId: number, artistName: string) {
    if (expandedArtist === artistId) {
      setExpandedArtist(null);
      setArtistSongs([]);
      return;
    }

    setExpandedArtist(artistId);
    try {
      const res = await fetch(`/api/artists/${artistId}/songs`);
      if (!res.ok) {
        throw new Error('获取歌曲列表失败');
      }
      const data = await res.json();
      setArtistSongs(data.data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : '获取歌曲列表失败');
      setExpandedArtist(null);
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
          共 {artists.length} 位艺术家
        </span>
        <button
          onClick={() => setOrderAsc(!orderAsc)}
          class="text-sm text-green-600 hover:text-green-700"
        >
          按名称{orderAsc ? '升序' : '降序'} ↕
        </button>
      </div>

      {/* 艺术家列表 */}
      <div class="flex-1 overflow-auto">
        {artists.length === 0 ? (
          <div class="flex flex-col items-center justify-center py-12 text-gray-500">
            <div class="text-4xl mb-4">🎵</div>
            <div class="text-lg">暂无艺术家</div>
            <div class="text-sm">请先扫描音乐目录</div>
          </div>
        ) : (
          <table class="w-full">
            <thead class="bg-gray-50 sticky top-0">
              <tr>
                <th class="w-8 px-2 py-2"></th>
                <th class="px-2 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  艺术家
                </th>
                <th class="px-2 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  歌曲数量
                </th>
              </tr>
            </thead>
            <tbody>
              {artists.map((artist) => (
                <>
                  {/* 艺术家行 */}
                  <tr
                    key={artist.id}
                    onClick={() => toggleArtist(artist.id, artist.name)}
                    class={{
                      'cursor-pointer hover:bg-gray-50 border-b border-gray-100': true,
                      'bg-gray-50': expandedArtist === artist.id,
                    }}
                  >
                    <td class="w-8 px-2 py-3">
                      <span class={{ 'text-gray-400': true, 'transform rotate-90': expandedArtist === artist.id }}>
                        ›
                      </span>
                    </td>
                    <td class="px-2 py-3 font-medium text-gray-900">
                      {artist.name}
                    </td>
                    <td class="px-2 py-3 text-sm text-gray-500">
                      {artist.songCount} 首歌曲
                    </td>
                  </tr>

                  {/* 展开的歌曲列表 */}
                  {expandedArtist === artist.id && (
                    <tr key={`${artist.id}-songs`}>
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
                                时长
                              </th>
                              <th class="w-12 px-2"></th>
                            </tr>
                          </thead>
                          <tbody>
                            {artistSongs.map((song) => (
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
                </>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
