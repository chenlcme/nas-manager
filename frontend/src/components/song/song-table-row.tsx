import { Song } from '../../types/song';
import { useSelection } from '../../contexts/selection-context';

interface SongTableRowProps {
  song: Song;
  onPlay: (song: Song) => void;
  showPath?: boolean; // 显示文件路径（用于文件夹视图）
}

// 格式化时长 (秒 -> mm:ss)
function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

export function SongTableRow({ song, onPlay, showPath }: SongTableRowProps) {
  const { isSelected, toggle } = useSelection();
  const selected = isSelected(song.id);

  return (
    <tr
      class={{
        'border-b border-gray-100 hover:bg-gray-50': true,
        'bg-green-50': selected,
      }}
    >
      {/* 复选框 */}
      <td class="w-8 px-2">
        <input
          type="checkbox"
          checked={selected}
          onChange={() => toggle(song.id)}
          class="w-4 h-4 text-green-500 rounded border-gray-300 focus:ring-green-500"
        />
      </td>

      {/* 封面 */}
      <td class="w-12 px-2">
        {song.coverPath ? (
          <img
            src={song.coverPath}
            alt=""
            class="w-10 h-10 rounded object-cover"
          />
        ) : (
          <div class="w-10 h-10 rounded bg-gray-200 border-2 border-dashed border-orange-400 flex items-center justify-center">
            <span class="text-gray-400 text-xs">无封面</span>
          </div>
        )}
      </td>

      {/* 歌名 */}
      <td class="px-2 py-3">
        <div class="font-medium text-gray-900 truncate max-w-[200px]">
          {song.title || '未知歌曲'}
        </div>
      </td>

      {/* 艺术家 */}
      <td class="px-2 py-3 text-sm text-gray-600 truncate max-w-[120px]">
        {song.artist || '未知艺术家'}
      </td>

      {/* 专辑 */}
      <td class="px-2 py-3 text-sm text-gray-600 truncate max-w-[150px]">
        {song.album || '未知专辑'}
      </td>

      {/* 年份 */}
      <td class="px-2 py-3 text-sm text-gray-500 w-16">
        {song.year || '-'}
      </td>

      {/* 流派 */}
      <td class="px-2 py-3 text-sm text-gray-500 w-20 truncate">
        {song.genre || '-'}
      </td>

      {/* 文件路径 (仅文件夹视图显示) */}
      {showPath && (
        <td class="px-2 py-3 text-sm text-gray-400 truncate max-w-[200px]" title={song.filePath}>
          {song.filePath}
        </td>
      )}

      {/* 时长 */}
      <td class="px-2 py-3 text-sm text-gray-500 w-14">
        {formatDuration(song.duration)}
      </td>

      {/* 操作 */}
      <td class="px-2 py-3 w-12">
        <button
          onClick={() => onPlay(song)}
          class="p-2 text-gray-400 hover:text-green-500 transition-colors"
          title="播放"
        >
          ▶
        </button>
      </td>
    </tr>
  );
}
