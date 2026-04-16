interface DeleteResultDialogProps {
  succeeded: number;
  failed: number;
  results: Array<{
    id: number;
    file_path: string;
    status: string;
    error?: string;
  }>;
  onClose: () => void;
}

export function DeleteResultDialog({ succeeded, failed, results, onClose }: DeleteResultDialogProps) {
  const failedResults = results.filter(r => r.status === 'failed');

  return (
    <div class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div class="bg-white rounded-lg shadow-xl max-w-lg w-full mx-4 overflow-hidden">
        <div class="px-6 py-4">
          <div class="flex items-center gap-3 mb-4">
            {failed === 0 ? (
              <>
                <div class="flex-shrink-0 w-10 h-10 rounded-full bg-green-100 flex items-center justify-center">
                  <svg class="w-6 h-6 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                </div>
                <h3 class="text-lg font-semibold text-gray-900">删除完成</h3>
              </>
            ) : (
              <>
                <div class="flex-shrink-0 w-10 h-10 rounded-full bg-red-100 flex items-center justify-center">
                  <svg class="w-6 h-6 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </div>
                <h3 class="text-lg font-semibold text-gray-900">删除完成（部分失败）</h3>
              </>
            )}
          </div>

          <div class="mt-3 space-y-3">
            <div class="flex gap-4">
              <div class="flex items-center gap-2">
                <span class="text-sm text-gray-600">成功:</span>
                <span class="text-sm font-medium text-green-600">{succeeded}</span>
              </div>
              <div class="flex items-center gap-2">
                <span class="text-sm text-gray-600">失败:</span>
                <span class="text-sm font-medium text-red-600">{failed}</span>
              </div>
            </div>

            {failedResults.length > 0 && (
              <div class="mt-4">
                <h4 class="text-sm font-medium text-gray-700 mb-2">失败文件:</h4>
                <div class="max-h-48 overflow-y-auto border border-gray-200 rounded-md">
                  <table class="min-w-full divide-y divide-gray-200">
                    <thead class="bg-gray-50">
                      <tr>
                        <th class="px-3 py-2 text-left text-xs font-medium text-gray-500">文件</th>
                        <th class="px-3 py-2 text-left text-xs font-medium text-gray-500">原因</th>
                      </tr>
                    </thead>
                    <tbody class="bg-white divide-y divide-gray-200">
                      {failedResults.map((r) => (
                        <tr key={r.id}>
                          <td class="px-3 py-2 text-xs text-gray-900 truncate max-w-xs" title={r.file_path}>
                            {r.file_path}
                          </td>
                          <td class="px-3 py-2 text-xs text-red-600">{r.error}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            )}
          </div>
        </div>
        <div class="px-6 py-3 bg-gray-50 flex justify-end">
          <button
            onClick={onClose}
            class="px-4 py-2 bg-gray-300 text-gray-700 text-sm font-medium rounded-lg hover:bg-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500"
          >
            关闭
          </button>
        </div>
      </div>
    </div>
  );
}
