import { SORT_FIELDS, SortField, SortOrder } from '../../constants/sort';

interface SortSelectorProps {
  sortBy: SortField;
  order: SortOrder;
  onSortChange: (sortBy: SortField, order: SortOrder) => void;
}

const sortLabels: Record<SortField, string> = {
  title: '名称',
  duration: '时长',
  created_at: '添加时间',
};

export function SortSelector({ sortBy, order, onSortChange }: SortSelectorProps) {
  return (
    <div class="flex items-center gap-3">
      {/* 排序字段选择 */}
      <select
        value={sortBy}
        onChange={(e) => onSortChange((e.target as HTMLSelectElement).value as SortField, order)}
        class="text-sm border border-gray-300 rounded px-2 py-1 text-gray-700 focus:outline-none focus:ring-2 focus:ring-green-500"
      >
        {SORT_FIELDS.map((field) => (
          <option key={field} value={field}>
            按{sortLabels[field]}
          </option>
        ))}
      </select>

      {/* 升序/降序切换 */}
      <button
        onClick={() => onSortChange(sortBy, order === 'asc' ? 'desc' : 'asc')}
        class="flex items-center gap-1 text-sm text-gray-600 hover:text-gray-900 px-2 py-1 rounded hover:bg-gray-100"
        title={order === 'asc' ? '升序' : '降序'}
      >
        <span>{sortLabels[sortBy]}</span>
        <span class={order === 'asc' ? '' : 'rotate-180'}>
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12" />
          </svg>
        </span>
      </button>
    </div>
  );
}