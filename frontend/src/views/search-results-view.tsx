import { useState, useEffect, useCallback } from 'preact/hooks';
import { SongTableRow } from '../components/song/song-table-row';
import { Song } from '../types/song';
import { SearchType } from '../components/common/search-bar';

// Extract filename from full path (module-level utility)
function getFileName(filePath: string): string {
  if (!filePath) return '';
  const parts = filePath.split(/[/\\]/);
  return parts[parts.length - 1] || filePath;
}

// Escape special regex characters
function escapeRegExp(str: string): string {
  return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

interface SearchResultsViewProps {
  keyword: string;
  searchType: SearchType;
  onPlaySong: (song: Song) => void;
  onShowSongDetail: (song: Song) => void;
  onBatchEdit: () => void;
  onBack: () => void;
}

export function SearchResultsView({
  keyword,
  searchType,
  onPlaySong,
  onShowSongDetail,
  onBatchEdit,
  onBack,
}: SearchResultsViewProps) {
  const [results, setResults] = useState<Song[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!keyword) {
      setResults([]);
      setLoading(false);
      return;
    }

    setLoading(true);
    setError(null);

    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 10000);

    // Determine API endpoint based on search type
    const apiUrl = searchType === 'tag'
      ? `/api/songs/search/by-tag?q=${encodeURIComponent(keyword)}`
      : `/api/songs/search?q=${encodeURIComponent(keyword)}`;

    fetch(apiUrl, {
      signal: controller.signal,
    })
      .then((res) => {
        clearTimeout(timeoutId);
        if (!res.ok) {
          throw new Error('搜索请求失败');
        }
        return res.json();
      })
      .then((data) => {
        setResults(data.data || []);
        setLoading(false);
      })
      .catch((err) => {
        clearTimeout(timeoutId);
        if (err.name === 'AbortError') {
          setError('搜索超时，请重试');
        } else {
          setError(err.message || '搜索失败');
        }
        setLoading(false);
      });

    return () => {
      clearTimeout(timeoutId);
      controller.abort();
    };
  }, [keyword, searchType]);

  // Highlight matching text in string
  const highlightMatch = useCallback((text: string, keyword: string): preact.JSX.Element => {
    if (!keyword) return <span>{text}</span>;

    const parts = text.split(new RegExp(`(${escapeRegExp(keyword)})`, 'gi'));
    return (
      <span>
        {parts.map((part, i) =>
          part.toLowerCase() === keyword.toLowerCase() ? (
            <mark key={i} class="bg-yellow-200 text-gray-900 px-0.5 rounded">
              {part}
            </mark>
          ) : (
            <span key={i}>{part}</span>
          )
        )}
      </span>
    );
  }, []);

  // Get highlighted title text
  const getHighlightedTitle = useCallback((song: Song): preact.JSX.Element => {
    if (searchType === 'tag') {
      return highlightMatch(song.title || getFileName(song.file_path), keyword);
    }
    return <span>{song.title || getFileName(song.file_path)}</span>;
  }, [searchType, keyword, highlightMatch]);

  // Get highlighted artist text
  const getHighlightedArtist = useCallback((song: Song): preact.JSX.Element => {
    if (searchType === 'tag') {
      return highlightMatch(song.artist || '-', keyword);
    }
    return <span>{song.artist || '-'}</span>;
  }, [searchType, keyword, highlightMatch]);

  // Get highlighted album text
  const getHighlightedAlbum = useCallback((song: Song): preact.JSX.Element => {
    if (searchType === 'tag') {
      return highlightMatch(song.album || '-', keyword);
    }
    return <span>{song.album || '-'}</span>;
  }, [searchType, keyword, highlightMatch]);

  if (loading) {
    return (
      <div class="flex flex-col h-full">
        {/* Header */}
        <div class="flex items-center justify-between px-4 py-3 border-b border-gray-200 bg-gray-50">
          <button
            onClick={onBack}
            class="flex items-center text-gray-600 hover:text-gray-900 min-h-[44px]"
          >
            <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
            返回
          </button>
          <span class="text-sm text-gray-500">搜索: "{keyword}"</span>
        </div>

        {/* Loading state */}
        <div class="flex-1 flex items-center justify-center">
          <div class="text-center">
            <div class="w-8 h-8 border-3 border-green-500 border-t-transparent rounded-full animate-spin mx-auto mb-4" />
            <p class="text-gray-500">搜索中...</p>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div class="flex flex-col h-full">
        {/* Header */}
        <div class="flex items-center justify-between px-4 py-3 border-b border-gray-200 bg-gray-50">
          <button
            onClick={onBack}
            class="flex items-center text-gray-600 hover:text-gray-900 min-h-[44px]"
          >
            <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
            返回
          </button>
          <span class="text-sm text-gray-500">搜索: "{keyword}"</span>
        </div>

        {/* Error state */}
        <div class="flex-1 flex items-center justify-center">
          <div class="text-center">
            <p class="text-red-500 mb-2">{error}</p>
            <button
              onClick={onBack}
              class="px-4 py-2 text-sm text-green-600 hover:text-green-700"
            >
              返回重试
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div class="flex flex-col h-full">
      {/* Header */}
      <div class="flex items-center justify-between px-4 py-3 border-b border-gray-200 bg-gray-50">
        <button
          onClick={onBack}
          class="flex items-center text-gray-600 hover:text-gray-900 min-h-[44px]"
        >
          <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
          返回
        </button>
        <span class="text-sm text-gray-500">
          搜索 "{keyword}" - 找到 {results.length} 首歌曲
        </span>
      </div>

      {/* Results */}
      <div class="flex-1 overflow-auto">
        {results.length === 0 ? (
          <div class="flex flex-col items-center justify-center h-full text-gray-500">
            <svg class="w-12 h-12 sm:w-16 sm:h-16 mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <p class="text-lg font-medium">未找到匹配的歌曲</p>
            <p class="text-sm mt-1">请尝试其他关键词</p>
          </div>
        ) : (
          <table class="w-full">
            <thead class="bg-gray-50 sticky top-0">
              <tr class="text-left text-xs text-gray-500 uppercase">
                <th class="px-4 py-2 w-8"></th>
                <th class="px-4 py-2">歌曲</th>
                <th class="px-4 py-2">歌手</th>
                <th class="px-4 py-2">专辑</th>
                <th class="px-4 py-2">时长</th>
              </tr>
            </thead>
            <tbody>
              {results.map((song) => (
                <tr
                  key={song.id}
                  class="border-b border-gray-100 hover:bg-gray-50 cursor-pointer"
                  onClick={() => onShowSongDetail(song)}
                >
                  <td class="px-4 py-2">
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        onPlaySong(song);
                      }}
                      class="p-1 hover:bg-green-100 rounded-full"
                    >
                      <svg class="w-5 h-5 text-green-600" fill="currentColor" viewBox="0 0 24 24">
                        <path d="M8 5v14l11-7z" />
                      </svg>
                    </button>
                  </td>
                  <td class="px-4 py-2">
                    {searchType === 'filename' ? (
                      highlightMatch(getFileName(song.file_path), keyword)
                    ) : (
                      getHighlightedTitle(song)
                    )}
                  </td>
                  <td class="px-4 py-2 text-gray-600">
                    {getHighlightedArtist(song)}
                  </td>
                  <td class="px-4 py-2 text-gray-600">
                    {getHighlightedAlbum(song)}
                  </td>
                  <td class="px-4 py-2 text-gray-500 text-sm">
                    {song.duration ? formatDuration(song.duration) : '-'}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}

// Format duration in seconds to mm:ss
function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}