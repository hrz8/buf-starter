export type QueryOptions = {
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
};
