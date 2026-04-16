import { useSelection } from '../../contexts/selection-context';

interface SelectionBarProps {
  totalCount: number;
  onBatchEdit?: () => void;
  onSelectAll?: () => void;
}

export function SelectionBar({ totalCount, onBatchEdit, onSelectAll }: SelectionBarProps) {
  const { selected, selectAll, clear, count } = useSelection();

  if (count === 0) {
    return null;
  }

  const handleSelectAll = () => {
    if (onSelectAll) {
      onSelectAll();
    } else {
      // Fallback: this shouldn't happen if views implement correctly
      console.warn('onSelectAll not provided');
    }
  };

  return (
    <div class="bg-green-50 border-b border-gray-200 sticky top-0 z-40">
      <div class="max-w-7xl mx-auto px-4 py-3 flex items-center justify-between">
        <div class="flex items-center space-x-4">
          <span class="text-sm text-gray-700">
            已选中 <span class="font-medium text-gray-900">{count}</span> / {totalCount} 首歌曲
          </span>
          <button
            onClick={handleSelectAll}
            class="text-sm text-green-600 hover:text-green-700 font-medium"
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
        {onBatchEdit && (
          <div class="flex items-center space-x-3">
            <button
              onClick={onBatchEdit}
              class="px-4 py-2 bg-green-500 text-white text-sm font-medium rounded-lg hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-green-500"
            >
              批量编辑
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
