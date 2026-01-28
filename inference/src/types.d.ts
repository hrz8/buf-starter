export interface Server<T = unknown> {
  name: string;
  instance: T | null;
  start: () => Promise<T>;
  stop: () => Promise<void>;
}
