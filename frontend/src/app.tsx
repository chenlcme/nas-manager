import { useState, useEffect, useRef } from 'preact/hooks';
import { SetupView } from './views/setup-view';
import { ArtistsView } from './views/artists-view';
import { AlbumsView } from './views/albums-view';
import { FoldersView } from './views/folders-view';
import { SearchResultsView } from './views/search-results-view';
import { TabNav } from './components/common/tab-nav';
import { SearchBar, SearchType } from './components/common/search-bar';
import { SongDetailPanel } from './components/song/song-detail-panel';
import { SidePlayer } from './components/player/side-player';
import { SideEditPanel } from './components/edit/side-edit-panel';
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
  const [detailSong, setDetailSong] = useState<Song | null>(null);
  const [detailLoading, setDetailLoading] = useState(false);
  const [toasts, setToasts] = useState<Toast[]>([]);
  const [searchKeyword, setSearchKeyword] = useState('');
  const [searchType, setSearchType] = useState<SearchType>('tag');
  const [searchLoading, setSearchLoading] = useState(false);
  const [scanning, setScanning] = useState(false);
  const toastIdRef = useRef(0);
  const abortControllerRef = useRef<AbortController | null>(null);
  const fetchSongIdRef = useRef<number | null>(null);
  const toastTimeoutRefs = useRef<Map<number, ReturnType<typeof setTimeout>>>(new Map());
  const [showPlayer, setShowPlayer] = useState(false);
  const [showEditPanel, setShowEditPanel] = useState(false);

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
        setView(data.data.needs_setup ? 'setup' : 'main');
      } else {
        setView('setup');
      }
    } catch {
      setView('setup');
    }
  }

  function handlePlaySong(song: Song) {
    setCurrentSong(song);
    setShowPlayer(true);
  }

  function handleClosePlayer() {
    setShowPlayer(false);
  }

  function handleBatchEdit() {
    setShowEditPanel(true);
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

  function handleSearch(keyword: string, searchType: SearchType) {
    setSearchKeyword(keyword);
    setSearchType(searchType);
    setSearchLoading(true);
  }

  function handleBackFromSearch() {
    setSearchKeyword('');
    setSearchLoading(false);
  }

  async function handleScan() {
    setScanning(true);
    try {
      const res = await fetch('/api/songs/scan', { method: 'POST' });
      const data = await res.json();
      if (!res.ok) {
        showToast(data.error?.message || '扫描失败', 'error');
      } else {
        showToast(`扫描完成：${data.data?.added || 0} 首新歌曲`, 'success');
        // 刷新当前视图
        checkSetupStatus();
      }
    } catch (err) {
      showToast('扫描请求失败', 'error');
    } finally {
      setScanning(false);
    }
  }

  if (view === 'setup' || setupStatus?.needs_setup) {
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
        <div class="flex items-center justify-between px-4 py-2 border-b border-gray-200 bg-white">
          <TabNav activeTab={activeTab} onTabChange={setActiveTab} />
          <div class="flex items-center gap-2">
            <button
              onClick={handleScan}
              disabled={scanning}
              class="px-3 py-1.5 text-sm bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg disabled:opacity-50 flex items-center gap-1"
            >
              {scanning ? (
                <>
                  <span class="animate-spin">⟳</span>
                  扫描中...
                </>
              ) : (
                <>
                  <span>🔄</span>
                  扫描
                </>
              )}
            </button>
            <SearchBar onSearch={handleSearch} loading={searchLoading} />
          </div>
        </div>

        {/* 主内容区 */}
        <main class="flex-1 overflow-hidden">
          {searchKeyword ? (
            <SearchResultsView
              keyword={searchKeyword}
              searchType={searchType}
              onPlaySong={handlePlaySong}
              onShowSongDetail={handleShowSongDetail}
              onBatchEdit={handleBatchEdit}
              onBack={handleBackFromSearch}
              playingSongId={currentSong?.id ?? null}
            />
          ) : activeTab === 'artists' ? (
            <ArtistsView onPlaySong={handlePlaySong} onShowSongDetail={handleShowSongDetail} onBatchEdit={handleBatchEdit} playingSongId={currentSong?.id ?? null} />
          ) : activeTab === 'albums' ? (
            <AlbumsView onPlaySong={handlePlaySong} onShowSongDetail={handleShowSongDetail} onBatchEdit={handleBatchEdit} playingSongId={currentSong?.id ?? null} />
          ) : (
            <FoldersView onPlaySong={handlePlaySong} onShowSongDetail={handleShowSongDetail} onBatchEdit={handleBatchEdit} playingSongId={currentSong?.id ?? null} />
          )}
        </main>

        {/* 批量编辑面板 */}
        {showEditPanel && (
          <SideEditPanel
            onClose={() => setShowEditPanel(false)}
            onSuccess={() => {
              showToast('批量更新成功', 'success');
            }}
            onError={(msg) => showToast(msg, 'error')}
          />
        )}

        {/* 歌曲详情面板 */}
        <SongDetailPanel
          song={detailSong}
          loading={detailLoading}
          onClose={handleCloseSongDetail}
          onError={showToast}
        />

        {/* 播放器面板 */}
        {showPlayer && currentSong && (
          <SidePlayer
            song={currentSong}
            onClose={handleClosePlayer}
            onEdit={handleClosePlayer}
            onError={(msg) => showToast(msg, 'error')}
          />
        )}
      </div>
    </SelectionProvider>
  );
}
