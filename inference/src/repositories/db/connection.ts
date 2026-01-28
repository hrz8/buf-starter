import type { Knex as KnexType } from 'knex';
import Knex from 'knex';

interface Config {
  getConnectionString: () => string;
}

export class Database {
  private static instance: KnexType | null = null;

  private constructor() {}

  public static initialize(config: Config): KnexType {
    if (this.instance) {
      return this.instance;
    }

    this.instance = Knex({
      client: 'pg',
      connection: config.getConnectionString(),
      pool: {
        min: 2,
        max: 10,
      },
    });

    console.info('database connected');

    return this.instance;
  }

  public static get(): KnexType {
    if (!this.instance) {
      throw new Error('Database not initialized. Call Database.init() first.');
    }
    return this.instance;
  }

  public static async close(): Promise<void> {
    if (this.instance) {
      await this.instance.destroy();
      this.instance = null;
    }
  }
}
