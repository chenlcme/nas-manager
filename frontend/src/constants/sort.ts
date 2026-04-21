// Sort field constants - must match backend validation
export const SORT_FIELDS = ['title', 'duration', 'created_at'] as const;
export type SortField = typeof SORT_FIELDS[number];

// Sort order constants
export const SORT_ORDERS = ['asc', 'desc'] as const;
export type SortOrder = typeof SORT_ORDERS[number];

// Default sort values
export const DEFAULT_SORT_FIELD: SortField = 'title';
export const DEFAULT_SORT_ORDER: SortOrder = 'asc';

// API parameter names
export const SORT_BY_PARAM = 'sort_by';
export const ORDER_PARAM = 'order';
export const FOLDER_PARAM = 'folder';

// Validate sort field
export function isValidSortField(field: string): field is SortField {
  return (SORT_FIELDS as readonly string[]).includes(field);
}

// Validate sort order
export function isValidSortOrder(order: string): order is SortOrder {
  return (SORT_ORDERS as readonly string[]).includes(order);
}

// Request timeout in milliseconds
export const REQUEST_TIMEOUT_MS = 10000;