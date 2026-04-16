import { useState, useRef } from 'preact/hooks';

interface SearchBarProps {
  onSearch: (keyword: string) => void;
  loading?: boolean;
}

function doSearch(keyword: string, onSearch: (keyword: string) => void) {
  const trimmed = keyword.trim();
  if (trimmed) {
    onSearch(trimmed);
  }
}

export function SearchBar({ onSearch, loading = false }: SearchBarProps) {
  const [keyword, setKeyword] = useState('');
  const inputRef = useRef<HTMLInputElement>(null);

  function handleSubmit(e: Event) {
    e.preventDefault();
    doSearch(keyword, onSearch);
  }

  function handleKeyDown(e: KeyboardEvent) {
    if (e.key === 'Enter') {
      e.preventDefault();
      doSearch(keyword, onSearch);
    }
  }

  return (
    <form onSubmit={handleSubmit} class="flex items-center">
      <div class="relative">
        <input
          ref={inputRef}
          type="text"
          value={keyword}
          onInput={(e) => setKeyword((e.target as HTMLInputElement).value)}
          onKeyDown={handleKeyDown}
          placeholder="搜索文件名..."
          class="w-full sm:w-48 md:w-64 pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent text-sm"
          disabled={loading}
        />
        {/* Search icon */}
        <svg
          class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
          />
        </svg>
        {loading && (
          <div class="absolute right-3 top-1/2 -translate-y-1/2">
            <div class="w-4 h-4 border-2 border-green-500 border-t-transparent rounded-full animate-spin" />
          </div>
        )}
      </div>
      <button
        type="submit"
        disabled={loading || !keyword.trim()}
        class="ml-2 px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors text-sm min-h-[44px] min-w-[44px]"
      >
        搜索
      </button>
    </form>
  );
}
