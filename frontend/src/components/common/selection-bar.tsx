import { useState } from 'preact/hooks';
import { useSelection } from '../../contexts/selection-context';
import { DeleteConfirmDialog } from './delete-confirm-dialog';
import { DeleteResultDialog } from './delete-result-dialog';

interface DeleteResult {
  total: number;
  succeeded: number;
  failed: number;
  results: Array<{
    id: number;
    file_path: string;
    status: string;
    error?: string;
  }>;
}

interface SelectionBarProps {
  totalCount: number;
  onBatchEdit?: () => void;
  onSelectAll?: () => void;
  onDeleteSuccess?: () => void;
}

export function SelectionBar({ totalCount, onBatchEdit, onSelectAll, onDeleteSuccess }: SelectionBarProps) {
  const { selected, selectAll, clear, count } = useSelection();
  const [showConfirm, setShowConfirm] = useState(false);
  const [showResult, setShowResult] = useState(false);
  const [deleteResult, setDeleteResult] = useState<DeleteResult | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);

  if (count === 0) {
    return null;
  }

  const handleSelectAll = () => {
    if (onSelectAll) {
      onSelectAll();
    } else {
      console.warn('onSelectAll not provided');
    }
  };

  const handleDelete = async () => {
    setShowConfirm(false);
    setIsDeleting(true);

    try {
      const ids = [...selected];
      const res = await fetch('/api/songs/delete', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ ids }),
      });

      if (!res.ok) {
        throw new Error('Delete failed');
      }

      const result = await res.json();
      setDeleteResult(result.data);
      setShowResult(true);

      if (result.data.failed === 0) {
        clear();
        if (onDeleteSuccess) {
          onDeleteSuccess();
        }
      }
    } catch (e) {
      console.error('Delete error:', e);
      // Show error as failed result
      setDeleteResult({
        total: selected.size,
        succeeded: 0,
        failed: selected.size,
        results: [...selected].map(id => ({
          id,
          file_path: '',
          status: 'failed',
          error: 'Network error',
        })),
      });
      setShowResult(true);
    } finally {
      setIsDeleting(false);
    }
  };

  const handleResultClose = () => {
    setShowResult(false);
    setDeleteResult(null);
  };

  return (
    <>
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
          <div class="flex items-center space-x-3">
            {onBatchEdit && (
              <button
                onClick={onBatchEdit}
                class="px-4 py-2 bg-green-500 text-white text-sm font-medium rounded-lg hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-green-500"
              >
                批量编辑
              </button>
            )}
            <button
              onClick={() => setShowConfirm(true)}
              disabled={isDeleting}
              class="px-4 py-2 bg-red-500 text-white text-sm font-medium rounded-lg hover:bg-red-600 focus:outline-none focus:ring-2 focus:ring-red-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isDeleting ? '删除中...' : '删除'}
            </button>
          </div>
        </div>
      </div>

      {showConfirm && (
        <DeleteConfirmDialog
          count={count}
          onConfirm={handleDelete}
          onCancel={() => setShowConfirm(false)}
        />
      )}

      {showResult && deleteResult && (
        <DeleteResultDialog
          succeeded={deleteResult.succeeded}
          failed={deleteResult.failed}
          results={deleteResult.results}
          onClose={handleResultClose}
        />
      )}
    </>
  );
}
