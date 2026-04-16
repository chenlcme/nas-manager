import { useState, useEffect, useRef } from 'preact/hooks';
import { SetupView } from './views/setup-view';
import { ArtistsView } from './views/artists-view';
import { AlbumsView } from './views/albums-view';
import { FoldersView } from './views/folders-view';
import { TabNav } from './components/common/tab-nav';
import { SongDetailPanel } from './components/song/song-detail-panel';
import { SelectionProvider } from './contexts/selection-context';
import { Song } from './types/song';

type View = 'setup' | 'main';
type Tab = 'artists' | 'albums' | 'folders';

interface SetupStatus {
  configured: boolean;
  hasPassword: boolean;
}

interface Toast {
  id: number;
  message: string;
  type: 'error' | 'success' | 'info';
}

const REQUEST_TIMEOUT_MS = 10000;

export function App() {
  const [view, setView] = useState<View>('setup');
  const [activeTab, setActiveTab] = useState<Tab>('artists');
  const [setupStatus, setSetupStatus] = useState<SetupStatus | null>(null);
  const [currentSong, setCurrentSong] = useState<Song | null>(null);
  const [showBatchEdit, setShowBatchEdit] = useState(false);
  const [detailSong, setDetailSong] = useState<Song | null>(null);
  const [detailLoading, setDetailLoading] = useState(false);
  const [toasts, setToasts] = useState<Toast[]>([]);
  const toastIdRef = useRef(0);
  const abortControllerRef = useRef<AbortController | null>(null);
  const fetchSongIdRef = useRef<number | null>(null);
  const toastTimeoutRefs = useRef<Map<number, ReturnType<typeof setTimeout>>>(new Map());

  // 检查设置状态
  useEffect(() => {
    checkSetupStatus();
  }, []);

  // Cleanup: abort in-flight requests and clear toast timeouts on unmount
  useEffect(() => {
    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
        abortControllerRef.current = null;
      }
      toastTimeoutRefs.current.forEach((timeoutId) => clearTimeout(timeoutId));
      toastTimeoutRefs.current.clear();
    };
  }, []);

  async function checkSetupStatus() {
    try {
      const res = await fetch('/api/setup/status');
      if (res.ok) {
        const data = await res.json();
        setSetupStatus(data.data);
        setView(data.data.configured ? 'main' : 'setup');
      } else {
        setView('setup');
      }
    } catch {
      setView('setup');
    }
  }

  function handlePlaySong(song: Song) {
    setCurrentSong(song);
    // 播放器功能将在 Epic 3 中实现
  }

  function handleBatchEdit() {
    setShowBatchEdit(true);
    // 批量编辑功能将在 Epic 4 中实现
  }

  function handleShowSongDetail(song: Song) {
    // Abort any in-flight request
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
      abortControllerRef.current = null;
    }
    setDetailSong(song);
    setDetailLoading(true);
    fetchSongDetail(song.id);
  }

  async function fetchSongDetail(songId: number) {
    const controller = new AbortController();
    abortControllerRef.current = controller;
    fetchSongIdRef.current = songId;
    const timeoutId = setTimeout(() => controller.abort(), REQUEST_TIMEOUT_MS);

    try {
      const res = await fetch(`/api/songs/${songId}`, {
        signal: controller.signal,
      });
      clearTimeout(timeoutId);

      if (!res.ok) {
        throw new Error('获取歌曲详情失败');
      }

      const data = await res.json();
      // Only update if we're still fetching the same song (panel not closed)
      if (fetchSongIdRef.current === songId) {
        setDetailSong(data.data);
        setDetailLoading(false);
      }
    } catch (err) {
      clearTimeout(timeoutId);
      // Only update state if we're still fetching the same song
      if (fetchSongIdRef.current === songId) {
        setDetailLoading(false);
        setDetailSong(null);

        const message = err instanceof Error ? err.message : '获取歌曲详情失败';
        showToast(message, 'error');
      }
    } finally {
      if (abortControllerRef.current === controller) {
        abortControllerRef.current = null;
      }
    }
  }

  function showToast(message: string, type: 'error' | 'success' | 'info' = 'error') {
    const id = ++toastIdRef.current;
    setToasts((prev) => [...prev, { id, message, type }]);
    // Auto-dismiss after 4 seconds
    const timeoutId = setTimeout(() => {
      setToasts((prev) => prev.filter((t) => t.id !== id));
      toastTimeoutRefs.current.delete(id);
    }, 4000);
    toastTimeoutRefs.current.set(id, timeoutId);
  }

  function handleCloseSongDetail() {
    // Clear fetch tracking to prevent stale fetch from updating state
    fetchSongIdRef.current = null;
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
      abortControllerRef.current = null;
    }
    setDetailSong(null);
    setDetailLoading(false);
  }

  if (view === 'setup' || !setupStatus?.configured) {
    return <SetupView />;
  }

  return (
    <SelectionProvider>
      <div class="min-h-screen bg-white flex flex-col">
        {/* Toast 通知 */}
        <div class="fixed top-4 right-4 z-[60] space-y-2">
          {toasts.map((toast) => (
            <div
              key={toast.id}
              class={`px-4 py-3 rounded-lg shadow-lg text-white text-sm animate-slide-in ${
                toast.type === 'error' ? 'bg-red-500' : toast.type === 'success' ? 'bg-green-500' : 'bg-blue-500'
              }`}
            >
              {toast.message}
            </div>
          ))}
        </div>

        {/* 顶部 Tab 导航 */}
        <TabNav activeTab={activeTab} onTabChange={setActiveTab} />

        {/* 主内容区 */}
        <main class="flex-1 overflow-hidden">
          {activeTab === 'artists' && (
            <ArtistsView onPlaySong={handlePlaySong} onShowSongDetail={handleShowSongDetail} onBatchEdit={handleBatchEdit} />
          )}
          {activeTab === 'albums' && (
            <AlbumsView onPlaySong={handlePlaySong} onShowSongDetail={handleShowSongDetail} onBatchEdit={handleBatchEdit} />
          )}
          {activeTab === 'folders' && (
            <FoldersView onPlaySong={handlePlaySong} onShowSongDetail={handleShowSongDetail} onBatchEdit={handleBatchEdit} />
          )}
        </main>

        {/* 批量编辑面板 (未来) */}
        {showBatchEdit && (
          <div class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center">
            <div class="bg-white rounded-lg p-6 max-w-md w-full mx-4">
              <h2 class="text-lg font-semibold mb-4">批量编辑</h2>
              <p class="text-gray-500 mb-4">批量编辑功能将在 Epic 4 中实现</p>
              <button
                onClick={() => setShowBatchEdit(false)}
                class="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300"
              >
                关闭
              </button>
            </div>
          </div>
        )}

        {/* 歌曲详情面板 */}
        <SongDetailPanel
          song={detailSong}
          loading={detailLoading}
          onClose={handleCloseSongDetail}
          onError={showToast}
        />
      </div>
    </SelectionProvider>
  );
}
