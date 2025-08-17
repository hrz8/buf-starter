import type { QueryOptions } from '../types/query';

export function serializeFilters(filters: NonNullable<QueryOptions['filters']>): string {
  return Object.keys(filters)
    .sort()
    .map((key) => `${key}:${filters[key]}`)
    .join('|');
}
