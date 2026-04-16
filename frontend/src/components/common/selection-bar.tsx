import { useSelection } from '../../contexts/selection-context';

interface SelectionBarProps {
  totalCount: number;
  onBatchEdit: () => void;
}

export function SelectionBar({ totalCount, onBatchEdit }: SelectionBarProps) {
  const { selected, selectAll, clear, count } = useSelection();

  if (count === 0) {
    return null;
  }

  return (
    <div class="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 shadow-lg z-50">
      <div class="max-w-7xl mx-auto px-4 py-3 flex items-center justify-between">
        <div class="flex items-center space-x-4">
          <span class="text-sm text-gray-700">
            已选中 <span class="font-medium text-gray-900">{count}</span> 首歌曲
          </span>
          <button
            onClick={() => selectAll(Array.from({ length: totalCount }, (_, i) => i + 1))}
            class="text-sm text-green-600 hover:text-green-700"
          >
            全选
          </button>
          <button
            onClick={clear}
            class="text-sm text-gray-500 hover:text-gray-700"
          >
            取消全选
          </button>
        </div>
        <div class="flex items-center space-x-3">
          <button
            onClick={onBatchEdit}
            class="px-4 py-2 bg-green-500 text-white text-sm font-medium rounded-lg hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-green-500"
          >
            批量编辑
          </button>
        </div>
      </div>
    </div>
  );
}
