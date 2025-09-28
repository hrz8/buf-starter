import type { MessageInitShape } from '@bufbuild/protobuf';
import type { StringListSchema } from '~~/gen/altalune/v1/common_pb';
import type { QueryOptions } from '../types/query';

/**
 * @deprecated since using structure from proto
 */
export function serializeFilters(filters: NonNullable<QueryOptions['filters']>): string {
  return Object.keys(filters)
    .sort()
    .map(key => `${key}:${filters[key]}`)
    .join('|');
}

export function serializeProtoFilters(
  filters: Record<string, MessageInitShape<typeof StringListSchema>> | undefined,
): string | null {
  if (!filters)
    return null;

  return Object.keys(filters)
    .sort()
    .map((key) => {
      const values = filters?.[key]?.values || [];
      return `${key}:${values.join(',')}`;
    })
    .join('|');
}
