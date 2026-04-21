import { Song } from '../../types/song';
import { useSelection } from '../../contexts/selection-context';

interface SongTableRowProps {
  song: Song;
  onPlay: (song: Song) => void;
  onShowDetail?: (song: Song) => void;
  showPath?: boolean; // 显示完整文件路径（用于文件夹视图）
  showDir?: boolean; // 显示相对目录（用于歌曲列表视图）
  highlightedText?: string | preact.JSX.Element; // 高亮文本（用于搜索结果）
  playingSongId?: number | null; // 当前播放的歌曲ID
}

// 格式化时长 (秒 -> mm:ss)
function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

export function SongTableRow({ song, onPlay, onShowDetail, showPath, showDir, highlightedText, playingSongId }: SongTableRowProps) {
  const { isSelected, toggle } = useSelection();
  const selected = isSelected(song.id);
  const isPlaying = playingSongId === song.id;

  return (
    <tr
      class={{
        'border-b border-gray-100 hover:bg-gray-50': true,
        'bg-green-50': selected,
        'bg-green-100': isPlaying,
      }}
    >
      {/* 复选框 - 缩小尺寸 */}
      <td class="w-6 px-1 py-2">
        <input
          type="checkbox"
          checked={selected}
          onChange={() => toggle(song.id)}
          class="w-3 h-3 text-green-500 rounded border-gray-300 focus:ring-green-500"
        />
      </td>

      {/* 封面 */}
      <td class="w-10 px-1 py-2">
        {song.coverPath ? (
          <img
            src={song.coverPath}
            alt=""
            class="w-8 h-8 rounded object-cover"
          />
        ) : (
          <div class="w-8 h-8 rounded bg-gray-200 border-2 border-dashed border-orange-400 flex items-center justify-center">
            <span class="text-gray-400 text-[8px]">无</span>
          </div>
        )}
      </td>

      {/* 歌名 */}
      <td class="px-2 py-2">
        <div class="font-medium text-gray-900 truncate max-w-[200px]">
          {highlightedText || song.title || '未知歌曲'}
        </div>
      </td>

      {/* 艺术家 */}
      <td class="px-2 py-2 text-sm text-gray-600 truncate max-w-[120px]">
        {song.artist || '未知艺术家'}
      </td>

      {/* 专辑 */}
      <td class="px-2 py-2 text-sm text-gray-600 truncate max-w-[150px]">
        {song.album || '未知专辑'}
      </td>

      {/* 年份 */}
      <td class="px-2 py-2 text-sm text-gray-500 w-14">
        {song.year || '-'}
      </td>

      {/* 流派 */}
      <td class="px-2 py-2 text-sm text-gray-500 w-16 truncate">
        {song.genre || '-'}
      </td>

      {/* 文件路径 (仅文件夹视图显示) */}
      {showPath && (
        <td class="px-2 py-2 text-sm text-gray-400 truncate max-w-[200px]" title={song.filePath}>
          {song.filePath}
        </td>
      )}

      {/* 所属目录 (用于歌曲列表视图，显示相对路径) */}
      {showDir && (
        <td class="px-2 py-2 text-sm text-gray-400 truncate max-w-[150px]" title={song.filePath}>
          {(() => {
            const parts = song.filePath.split('/');
            parts.pop(); // 移除文件名
            return parts.join('/') || '/';
          })()}
        </td>
      )}

      {/* 时长 */}
      <td class="px-2 py-2 text-sm text-gray-500 w-12">
        {formatDuration(song.duration)}
      </td>

      {/* 操作 */}
      <td class="px-2 py-2 w-16">
        <div class="flex items-center gap-1">
          <button
            onClick={() => onPlay(song)}
            class={`p-1.5 rounded transition-colors ${isPlaying ? 'text-green-500' : 'text-gray-400 hover:text-white hover:bg-green-500'}`}
            title={isPlaying ? "播放中" : "播放"}
            aria-label={isPlaying ? "播放中" : "播放歌曲"}
          >
            {isPlaying ? (
              <svg class="w-4 h-4 animate-pulse" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 3v9.28a4.39 4.39 0 0 0-1.5-.28C8.01 12 6 14.01 6 16.5S8.01 21 10.5 21c2.31 0 4.2-1.72 4.45-3.94.57-.3 1.05-.75 1.05-1.78V6h4V3h-8z"/>
              </svg>
            ) : (
              <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
                <path d="M8 5v14l11-7z"/>
              </svg>
            )}
          </button>
          {onShowDetail && (
            <button
              onClick={() => onShowDetail(song)}
              class="p-1.5 rounded text-gray-400 hover:text-white hover:bg-blue-500 transition-colors"
              title="详情"
              aria-label="查看歌曲详情"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth="2">
                <path strokeLinecap="round" strokeLinejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
            </button>
          )}
        </div>
      </td>
    </tr>
  );
}
