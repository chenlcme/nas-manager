import { JSX } from 'preact/hooks';

interface TabNavProps {
  activeTab: 'artists' | 'albums' | 'folders';
  onTabChange: (tab: 'artists' | 'albums' | 'folders') => void;
}

export function TabNav({ activeTab, onTabChange }: TabNavProps) {
  const tabs: { id: 'artists' | 'albums' | 'folders'; label: string }[] = [
    { id: 'artists', label: '歌手' },
    { id: 'albums', label: '专辑' },
    { id: 'folders', label: '文件夹' },
  ];

  return (
    <div class="flex border-b border-gray-200">
      {tabs.map((tab) => (
        <button
          key={tab.id}
          onClick={() => onTabChange(tab.id)}
          class={{
            'px-6 py-3 text-sm font-medium border-b-2 transition-colors': true,
            'border-green-500 text-gray-900': activeTab === tab.id,
            'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300':
              activeTab !== tab.id,
          }}
        >
          {tab.label}
        </button>
      ))}
    </div>
  );
}
