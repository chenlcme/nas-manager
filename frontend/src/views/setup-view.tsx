import { useState } from 'preact/hooks';

// SetupView - 首次配置向导组件
export function SetupView() {
  const [step, setStep] = useState(1);
  const [musicDir, setMusicDir] = useState('');
  const [dbPath, setDBPath] = useState('~/.nas-manager/nas-manager.db');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleNext = () => {
    if (!musicDir.trim()) {
      setError('请输入音乐目录路径');
      return;
    }
    setError('');
    setStep(2);
  };

  const handleBack = () => {
    setStep(1);
    setError('');
  };

  const handleSubmit = async () => {
    setLoading(true);
    setError('');

    try {
      const response = await fetch('/api/setup', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          music_dir: musicDir,
          db_path: dbPath || undefined,
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error?.message || '配置保存失败');
      }

      // 配置成功，跳转到主页面
      window.location.href = '/';
    } catch (err) {
      setError(err instanceof Error ? err.message : '配置保存失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div class="min-h-screen bg-gray-50 flex items-center justify-center p-4">
      <div class="bg-white rounded-lg shadow-lg p-8 max-w-md w-full">
        <h1 class="text-2xl font-bold text-gray-900 mb-2">首次配置向导</h1>
        <p class="text-gray-600 mb-6">让我们开始配置您的音乐管理器</p>

        {/* Progress indicator */}
        <div class="flex mb-6">
          <div class={`flex-1 h-2 rounded-l ${step >= 1 ? 'bg-green-500' : 'bg-gray-200'}`} />
          <div class={`flex-1 h-2 rounded-r ${step >= 2 ? 'bg-green-500' : 'bg-gray-200'}`} />
        </div>

        {step === 1 && (
          <div>
            <h2 class="text-lg font-semibold text-gray-900 mb-4">步骤 1：设置音乐目录</h2>
            <p class="text-gray-600 mb-4">请选择包含您音乐文件的目录</p>

            <div class="mb-4">
              <label class="block text-gray-700 text-sm font-bold mb-2" for="music-dir">
                音乐目录路径
              </label>
              <input
                id="music-dir"
                type="text"
                value={musicDir}
                onInput={(e) => setMusicDir((e.target as HTMLInputElement).value)}
                placeholder="/mnt/music"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
              />
            </div>

            {error && (
              <p class="text-red-500 text-sm mb-4">{error}</p>
            )}

            <div class="flex justify-end">
              <button
                onClick={handleNext}
                class="px-6 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-green-500"
              >
                下一步
              </button>
            </div>
          </div>
        )}

        {step === 2 && (
          <div>
            <h2 class="text-lg font-semibold text-gray-900 mb-4">步骤 2：确认数据库路径</h2>
            <p class="text-gray-600 mb-4">SQLite 数据库将存储在这里（可选）</p>

            <div class="mb-4">
              <label class="block text-gray-700 text-sm font-bold mb-2" for="db-path">
                数据库路径
              </label>
              <input
                id="db-path"
                type="text"
                value={dbPath}
                onInput={(e) => setDBPath((e.target as HTMLInputElement).value)}
                placeholder="~/.nas-manager/nas-manager.db"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
              />
              <p class="text-gray-500 text-xs mt-1">留空使用默认路径</p>
            </div>

            {error && (
              <p class="text-red-500 text-sm mb-4">{error}</p>
            )}

            <div class="flex justify-between">
              <button
                onClick={handleBack}
                class="px-6 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-gray-500"
                disabled={loading}
              >
                返回
              </button>
              <button
                onClick={handleSubmit}
                disabled={loading}
                class="px-6 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-green-500 disabled:opacity-50"
              >
                {loading ? '保存中...' : '完成配置'}
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
