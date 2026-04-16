import { useState, useEffect } from 'preact/hooks';
import { SetupView } from './views/setup-view';
import { ArtistsView } from './views/artists-view';
import { TabNav } from './components/common/tab-nav';
import { SelectionBar } from './components/common/selection-bar';
import { SelectionProvider } from './contexts/selection-context';
import { Song } from './types/song';

type View = 'setup' | 'main';
type Tab = 'artists' | 'albums' | 'folders';

interface SetupStatus {
  configured: boolean;
  hasPassword: boolean;
}

export function App() {
  const [view, setView] = useState<View>('setup');
  const [activeTab, setActiveTab] = useState<Tab>('artists');
  const [setupStatus, setSetupStatus] = useState<SetupStatus | null>(null);
  const [currentSong, setCurrentSong] = useState<Song | null>(null);
  const [showBatchEdit, setShowBatchEdit] = useState(false);

  // 检查设置状态
  useEffect(() => {
    checkSetupStatus();
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

  if (view === 'setup' || !setupStatus?.configured) {
    return <SetupView />;
  }

  return (
    <SelectionProvider>
      <div class="min-h-screen bg-white flex flex-col">
        {/* 顶部 Tab 导航 */}
        <TabNav activeTab={activeTab} onTabChange={setActiveTab} />

        {/* 主内容区 */}
        <main class="flex-1 overflow-hidden">
          {activeTab === 'artists' && (
            <ArtistsView onPlaySong={handlePlaySong} />
          )}
          {activeTab === 'albums' && (
            <div class="flex items-center justify-center h-full text-gray-500">
              专辑视图 (Epic 2.2)
            </div>
          )}
          {activeTab === 'folders' && (
            <div class="flex items-center justify-center h-full text-gray-500">
              文件夹视图 (Epic 2.3)
            </div>
          )}
        </main>

        {/* 选择操作栏 */}
        <SelectionBar totalCount={0} onBatchEdit={handleBatchEdit} />

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
      </div>
    </SelectionProvider>
  );
}
