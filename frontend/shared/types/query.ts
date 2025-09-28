export interface QueryOptions {
  pagination: {
    page: number;
    pageSize: number;
  };
  keyword?: string;
  filters?: {
    [key: string]: string | string[] | number | boolean | null | undefined;
  };
  sorting?: {
    field: string;
    order: 'asc' | 'desc';
  };
}

export interface PaginatedResponse<T> {
  data: T[];
  meta: {
    rowCount: number;
    pageCount: number;
    filters?: {
      [key: string]: string[];
    };
  };
}
