import { useEffect, useRef, useState } from 'preact/hooks';
import { Song } from '../../types/song';

interface SongDetailPanelProps {
  song: Song | null;
  loading?: boolean;
  onClose: () => void;
  onError?: (message: string) => void;
}

// 格式化时长 (秒 -> mm:ss)
function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

// 格式化文件大小
function formatFileSize(bytes: number): string {
  if (bytes < 1024) {
    return bytes + ' B';
  } else if (bytes < 1024 * 1024) {
    return (bytes / 1024).toFixed(1) + ' KB';
  } else if (bytes < 1024 * 1024 * 1024) {
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
  } else {
    return (bytes / (1024 * 1024 * 1024)).toFixed(1) + ' GB';
  }
}

// 获取文件格式
function getFileFormat(filePath: string): string {
  const ext = filePath.split('.').pop();
  if (!ext) return '未知';
  return ext.toUpperCase();
}

export function SongDetailPanel({ song, loading, onClose, onError }: SongDetailPanelProps) {
  const panelRef = useRef<HTMLDivElement>(null);
  const closeButtonRef = useRef<HTMLButtonElement>(null);
  const [isVisible, setIsVisible] = useState(false);

  // Handle escape key
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && song) {
        onClose();
      }
    };

    if (song) {
      document.addEventListener('keydown', handleKeyDown);
      // Animate in
      requestAnimationFrame(() => setIsVisible(true));
      // Focus trap - focus the close button
      closeButtonRef.current?.focus();
    }

    return () => {
      document.removeEventListener('keydown', handleKeyDown);
    };
  }, [song, onClose]);

  // Handle click outside and focus trap
  useEffect(() => {
    if (!song) return;

    const handleFocusTrap = (e: FocusEvent) => {
      if (!panelRef.current) return;
      const focusable = panelRef.current.querySelectorAll<HTMLElement>(
        'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
      );
      const first = focusable[0];
      const last = focusable[focusable.length - 1];

      if (e.target === document.body || !panelRef.current.contains(e.target as Node)) {
        first?.focus();
      }
    };

    document.addEventListener('focus', handleFocusTrap, true);
    return () => document.removeEventListener('focus', handleFocusTrap, true);
  }, [song]);

  if (!song && !loading) return null;

  return (
    <div class={`fixed inset-0 z-50 flex justify-end transition-opacity duration-300 ${song ? (isVisible ? 'opacity-100' : 'opacity-0') : 'opacity-0 pointer-events-none'}`}>
      {/* 背景遮罩 */}
      <div
        class="absolute inset-0 bg-black bg-opacity-30"
        onClick={onClose}
      />

      {/* 滑入面板 - 移动端全屏，桌面端侧边栏 */}
      <div
        ref={panelRef}
        class={`relative w-full max-w-md bg-white shadow-xl h-full overflow-y-auto transform transition-transform duration-300 ease-out ${song ? (isVisible ? 'translate-x-0' : 'translate-x-full') : 'translate-x-full'} md:transform-none md:h-full md:translate-x-0`}
        role="dialog"
        aria-modal="true"
        aria-labelledby="song-detail-title"
      >
        {/* Loading state */}
        {loading && (
          <div class="absolute inset-0 bg-white bg-opacity-80 flex items-center justify-center z-10">
            <div class="flex flex-col items-center gap-2">
              <div class="w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full animate-spin" />
              <span class="text-gray-500 text-sm">加载中...</span>
            </div>
          </div>
        )}

        {/* 头部 */}
        <div class="sticky top-0 bg-white border-b border-gray-200 px-4 py-3 flex items-center justify-between">
          <h2 id="song-detail-title" class="text-lg font-semibold text-gray-900">歌曲详情</h2>
          <button
            ref={closeButtonRef}
            onClick={onClose}
            class="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
            title="关闭 (Esc)"
            aria-label="关闭详情面板"
          >
            ✕
          </button>
        </div>

        {/* 内容 */}
        <div class="p-4 space-y-6">
          {/* 封面 */}
          <div class="flex justify-center">
            {song.coverPath ? (
              <img
                src={song.coverPath}
                alt={song.title}
                class="w-48 h-48 rounded-lg object-cover shadow-md"
                onError={() => onError?.('封面加载失败')}
              />
            ) : (
              <div class="w-48 h-48 rounded-lg bg-gray-200 border-2 border-dashed border-orange-400 flex items-center justify-center">
                <span class="text-gray-400 text-sm">无封面</span>
              </div>
            )}
          </div>

          {/* ID3 信息 */}
          <div class="space-y-3">
            <h3 class="text-sm font-medium text-gray-500 uppercase tracking-wider">基本信息</h3>
            <div class="bg-gray-50 rounded-lg p-4 space-y-2">
              <div class="flex justify-between">
                <span class="text-gray-500">标题</span>
                <span class="text-gray-900 font-medium">{song.title || '未知歌曲'}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-500">艺术家</span>
                <span class="text-gray-900">{song.artist || '未知艺术家'}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-500">专辑</span>
                <span class="text-gray-900">{song.album || '未知专辑'}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-500">年份</span>
                <span class="text-gray-900">{song.year || '-'}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-500">流派</span>
                <span class="text-gray-900">{song.genre || '-'}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-500">曲目号</span>
                <span class="text-gray-900">{song.trackNum || '-'}</span>
              </div>
            </div>
          </div>

          {/* 文件信息 */}
          <div class="space-y-3">
            <h3 class="text-sm font-medium text-gray-500 uppercase tracking-wider">文件信息</h3>
            <div class="bg-gray-50 rounded-lg p-4 space-y-2">
              <div class="flex justify-between">
                <span class="text-gray-500">格式</span>
                <span class="text-gray-900">{getFileFormat(song.filePath)}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-500">大小</span>
                <span class="text-gray-900">{formatFileSize(song.fileSize)}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-500">时长</span>
                <span class="text-gray-900">{formatDuration(song.duration)}</span>
              </div>
              <div class="mt-2">
                <span class="text-gray-500 block mb-1">路径</span>
                <span class="text-gray-700 text-sm break-words overflow-wrap-anywhere max-w-full">{song.filePath}</span>
              </div>
            </div>
          </div>

          {/* 歌词 */}
          <div class="space-y-3">
            <h3 class="text-sm font-medium text-gray-500 uppercase tracking-wider">歌词</h3>
            <div class="bg-gray-50 rounded-lg p-4">
              {song.lyrics ? (
                <pre class="text-gray-700 text-sm whitespace-pre-wrap font-sans overflow-auto max-h-48">
                  {song.lyrics}
                </pre>
              ) : (
                <span class="text-gray-400 italic">暂无歌词</span>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
