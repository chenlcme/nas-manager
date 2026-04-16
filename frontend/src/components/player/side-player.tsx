import { useState, useRef, useEffect } from 'preact/hooks';
import { Song } from '../../types/song';

interface SidePlayerProps {
  song: Song | null;
  onClose: () => void;
  onEdit: () => void;
  onError?: (message: string) => void;
}

// 格式化时长 (秒 -> mm:ss)
function formatDuration(seconds: number): string {
  if (!seconds || seconds <= 0) return '0:00';
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

export function SidePlayer({ song, onClose, onEdit, onError }: SidePlayerProps) {
  const audioRef = useRef<HTMLAudioElement>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [currentTime, setCurrentTime] = useState(0);
  const [duration, setDuration] = useState(0);
  const [volume, setVolume] = useState(1);
  const [showEditPanel, setShowEditPanel] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [editForm, setEditForm] = useState({
    title: '',
    artist: '',
    album: '',
    year: '',
    genre: '',
    lyrics: '',
  });
  const [isSaving, setIsSaving] = useState(false);

  // Reset state when song changes
  useEffect(() => {
    if (song) {
      setEditForm({
        title: song.title || '',
        artist: song.artist || '',
        album: song.album || '',
        year: song.year?.toString() || '',
        genre: song.genre || '',
        lyrics: song.lyrics || '',
      });
      setShowEditPanel(false);
      setIsEditing(false);
    }
    setIsPlaying(false);
    setCurrentTime(0);
  }, [song?.id]);

  // Load audio when song changes
  useEffect(() => {
    if (song && audioRef.current) {
      audioRef.current.src = `/api/songs/${song.id}/stream`;
      audioRef.current.load();
    }
  }, [song?.id]);

  // Auto-play when song is set
  useEffect(() => {
    if (song && audioRef.current) {
      audioRef.current.play().catch((err) => {
        console.error('Failed to play:', err);
        onError?.('播放失败');
      });
      setIsPlaying(true);
    }
  }, [song?.id]);

  const handlePlayPause = () => {
    if (!audioRef.current) return;
    if (isPlaying) {
      audioRef.current.pause();
    } else {
      audioRef.current.play().catch((err) => {
        console.error('Failed to play:', err);
        onError?.('播放失败');
      });
    }
    setIsPlaying(!isPlaying);
  };

  const handleTimeUpdate = () => {
    if (audioRef.current) {
      setCurrentTime(audioRef.current.currentTime);
    }
  };

  const handleLoadedMetadata = () => {
    if (audioRef.current) {
      setDuration(audioRef.current.duration);
    }
  };

  const handleSeek = (e: Event) => {
    const target = e.target as HTMLInputElement;
    const time = parseFloat(target.value);
    if (audioRef.current) {
      audioRef.current.currentTime = time;
      setCurrentTime(time);
    }
  };

  const handleVolumeChange = (e: Event) => {
    const target = e.target as HTMLInputElement;
    const vol = parseFloat(target.value);
    if (audioRef.current) {
      audioRef.current.volume = vol;
      setVolume(vol);
    }
  };

  const handleEnded = () => {
    setIsPlaying(false);
    setCurrentTime(0);
  };

  const handleSaveEdit = async () => {
    if (!song) return;
    setIsSaving(true);

    try {
      // Validate year if provided
      if (editForm.year) {
        const yearNum = parseInt(editForm.year);
        if (isNaN(yearNum)) {
          onError?.('年份必须是有效数字');
          setIsSaving(false);
          return;
        }
      }

      const res = await fetch(`/api/songs/${song.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          title: editForm.title || null,
          artist: editForm.artist || null,
          album: editForm.album || null,
          year: editForm.year ? parseInt(editForm.year) : null,
          genre: editForm.genre || null,
          lyrics: editForm.lyrics || null,
        }),
      });

      if (!res.ok) {
        throw new Error('保存失败');
      }

      setIsEditing(false);
      setShowEditPanel(false);
      onError?.('已保存');
    } catch (err) {
      onError?.(err instanceof Error ? err.message : '保存失败');
    } finally {
      setIsSaving(false);
    }
  };

  const handleEditButton = () => {
    setShowEditPanel(true);
    setIsEditing(true);
  };

  if (!song) return null;

  return (
    <div class="fixed right-0 top-0 bottom-0 w-80 bg-white shadow-xl z-50 flex flex-col">
      {/* Header */}
      <div class="flex items-center justify-between px-4 py-3 border-b border-gray-200">
        <h2 class="text-lg font-semibold text-gray-900">正在播放</h2>
        <button
          onClick={onClose}
          class="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg"
          title="关闭"
        >
          ✕
        </button>
      </div>

      {/* Audio element */}
      <audio
        ref={audioRef}
        onTimeUpdate={handleTimeUpdate}
        onLoadedMetadata={handleLoadedMetadata}
        onEnded={handleEnded}
        onError={() => onError?.('音频加载失败')}
      />

      {/* Content */}
      <div class="flex-1 overflow-y-auto">
        {/* Cover */}
        <div class="p-4 flex justify-center">
          {song.coverPath ? (
            <img
              src={song.coverPath}
              alt={song.title}
              class="w-48 h-48 rounded-lg object-cover shadow-md"
              onError={() => onError?.('封面加载失败')}
            />
          ) : (
            <div class="w-48 h-48 rounded-lg bg-gray-200 border-2 border-dashed border-orange-400 flex items-center justify-center">
              <span class="text-6xl text-gray-400">♪</span>
            </div>
          )}
        </div>

        {/* Song Info */}
        <div class="px-4 text-center">
          <h3 class="text-xl font-semibold text-gray-900 truncate">{song.title || '未知歌曲'}</h3>
          <p class="text-gray-600 truncate">{song.artist || '未知艺术家'}</p>
          <p class="text-sm text-gray-500 truncate">{song.album || '未知专辑'}</p>
        </div>

        {/* Progress Bar */}
        <div class="px-4 py-4 space-y-2">
          <input
            type="range"
            min="0"
            max={duration || 0}
            value={currentTime}
            onChange={handleSeek}
            class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-green-500"
          />
          <div class="flex justify-between text-sm text-gray-500">
            <span>{formatDuration(currentTime)}</span>
            <span>{formatDuration(duration)}</span>
          </div>
        </div>

        {/* Controls */}
        <div class="px-4 flex items-center justify-center gap-4">
          <button
            onClick={() => audioRef.current && (audioRef.current.currentTime -= 10)}
            class="p-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg"
            title="后退10秒"
          >
            ⏪
          </button>
          <button
            onClick={handlePlayPause}
            class="p-4 bg-green-500 text-white rounded-full hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-green-500"
          >
            {isPlaying ? '⏸' : '▶'}
          </button>
          <button
            onClick={() => audioRef.current && (audioRef.current.currentTime += 10)}
            class="p-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg"
            title="前进10秒"
          >
            ⏩
          </button>
        </div>

        {/* Volume */}
        <div class="px-4 py-2 flex items-center gap-2">
          <span class="text-gray-400">🔈</span>
          <input
            type="range"
            min="0"
            max="1"
            step="0.1"
            value={volume}
            onChange={handleVolumeChange}
            class="flex-1 h-1 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-green-500"
          />
          <span class="text-gray-400">🔊</span>
        </div>

        {/* Edit Button */}
        <div class="px-4 py-2">
          <button
            onClick={handleEditButton}
            class="w-full py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg text-sm font-medium"
          >
            编辑信息
          </button>
        </div>

        {/* Lyrics */}
        <div class="px-4 py-4">
          <h4 class="text-sm font-medium text-gray-500 uppercase tracking-wider mb-2">歌词</h4>
          <div class="bg-gray-50 rounded-lg p-3 max-h-48 overflow-y-auto">
            {song.lyrics ? (
              <pre class="text-gray-700 text-sm whitespace-pre-wrap font-sans">{song.lyrics}</pre>
            ) : (
              <span class="text-gray-400 italic">暂无歌词</span>
            )}
          </div>
        </div>

        {/* Edit Panel */}
        {showEditPanel && (
          <div class="px-4 py-4 border-t border-gray-200">
            <h4 class="text-sm font-medium text-gray-500 uppercase tracking-wider mb-3">编辑信息</h4>
            <div class="space-y-3">
              <div>
                <label class="block text-xs text-gray-500 mb-1">标题</label>
                <input
                  type="text"
                  value={editForm.title}
                  onChange={(e) => setEditForm({ ...editForm, title: (e.target as HTMLInputElement).value })}
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500 text-sm"
                />
              </div>
              <div>
                <label class="block text-xs text-gray-500 mb-1">艺术家</label>
                <input
                  type="text"
                  value={editForm.artist}
                  onChange={(e) => setEditForm({ ...editForm, artist: (e.target as HTMLInputElement).value })}
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500 text-sm"
                />
              </div>
              <div>
                <label class="block text-xs text-gray-500 mb-1">专辑</label>
                <input
                  type="text"
                  value={editForm.album}
                  onChange={(e) => setEditForm({ ...editForm, album: (e.target as HTMLInputElement).value })}
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500 text-sm"
                />
              </div>
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <label class="block text-xs text-gray-500 mb-1">年份</label>
                  <input
                    type="text"
                    value={editForm.year}
                    onChange={(e) => setEditForm({ ...editForm, year: (e.target as HTMLInputElement).value })}
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500 text-sm"
                  />
                </div>
                <div>
                  <label class="block text-xs text-gray-500 mb-1">流派</label>
                  <input
                    type="text"
                    value={editForm.genre}
                    onChange={(e) => setEditForm({ ...editForm, genre: (e.target as HTMLInputElement).value })}
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500 text-sm"
                  />
                </div>
              </div>
              <div>
                <label class="block text-xs text-gray-500 mb-1">歌词</label>
                <textarea
                  value={editForm.lyrics}
                  onChange={(e) => setEditForm({ ...editForm, lyrics: (e.target as HTMLTextAreaElement).value })}
                  rows={4}
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500 text-sm resize-none"
                />
              </div>
              <div class="flex gap-2 pt-2">
                <button
                  onClick={() => {
                    setShowEditPanel(false);
                    setIsEditing(false);
                  }}
                  class="flex-1 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 text-sm font-medium"
                >
                  取消
                </button>
                <button
                  onClick={handleSaveEdit}
                  disabled={isSaving}
                  class="flex-1 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 text-sm font-medium disabled:opacity-50"
                >
                  {isSaving ? '保存中...' : '保存'}
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}