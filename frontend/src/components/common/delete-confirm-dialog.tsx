interface DeleteConfirmDialogProps {
  count: number;
  onConfirm: () => void;
  onCancel: () => void;
}

export function DeleteConfirmDialog({ count, onConfirm, onCancel }: DeleteConfirmDialogProps) {
  return (
    <div class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div class="bg-white rounded-lg shadow-xl max-w-md w-full mx-4 overflow-hidden">
        <div class="px-6 py-4">
          <div class="flex items-center gap-3 mb-4">
            <div class="flex-shrink-0 w-10 h-10 rounded-full bg-orange-100 flex items-center justify-center">
              <svg class="w-6 h-6 text-orange-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-gray-900">确认删除</h3>
          </div>
          <div class="mt-3">
            <p class="text-gray-700">
              将删除 <span class="font-medium text-gray-900">{count}</span> 首歌曲，文件和数据库记录都将被删除。
            </p>
            <p class="mt-2 text-sm text-red-600 font-medium">此操作不可撤销</p>
          </div>
        </div>
        <div class="px-6 py-3 bg-gray-50 flex justify-end gap-3">
          <button
            onClick={onCancel}
            class="px-4 py-2 bg-gray-300 text-gray-700 text-sm font-medium rounded-lg hover:bg-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500"
          >
            取消
          </button>
          <button
            onClick={onConfirm}
            class="px-4 py-2 bg-red-500 text-white text-sm font-medium rounded-lg hover:bg-red-600 focus:outline-none focus:ring-2 focus:ring-red-500"
          >
            确认删除
          </button>
        </div>
      </div>
    </div>
  );
}
