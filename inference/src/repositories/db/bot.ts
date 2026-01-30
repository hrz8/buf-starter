import type { Knex } from 'knex';

import type {
  ChatbotNode,
  ChatbotNodeMessage,
  ChatbotNodeTrigger,
  ModuleConfigMap,
  PlainMessage,
} from '../../app/bot/types.js';

// Database row types
interface ChatbotConfigRow {
  id: string;
  public_id: string;
  project_id: string;
  modules_config: ModuleConfigMap;
  created_at: Date;
  updated_at: Date;
}

interface ChatbotNodeRow {
  id: string;
  public_id: string;
  project_id: string;
  name: string;
  lang: string;
  tags: string[];
  enabled: boolean;
  triggers: ChatbotNodeTrigger[];
  messages: ChatbotNodeMessage[];
  created_at: Date;
  updated_at: Date;
}

export class BotRepository {
  private readonly db: Knex;

  public constructor(db: Knex) {
    this.db = db;
  }

  public async getModulesConfig(
    projectInternalId: string,
  ): Promise<Partial<ModuleConfigMap> | null> {
    const row = await this.db<ChatbotConfigRow>('altalune_chatbot_configs')
      .select('id', 'public_id', 'project_id', 'modules_config', 'created_at', 'updated_at')
      .where('project_id', projectInternalId)
      .first();

    if (!row) {
      return null;
    }

    return row.modules_config;
  }

  public async getNodes(
    projectInternalId: string,
  ): Promise<Map<string, PlainMessage<ChatbotNode>>> {
    const rows = await this.db<ChatbotNodeRow>('altalune_chatbot_nodes')
      .select(
        'id',
        'public_id',
        'project_id',
        'name',
        'lang',
        'tags',
        'enabled',
        'triggers',
        'messages',
        'created_at',
        'updated_at',
      )
      .where('project_id', projectInternalId)
      .andWhere('enabled', true);

    const nodes = new Map<string, PlainMessage<ChatbotNode>>();

    for (const row of rows) {
      const node: PlainMessage<ChatbotNode> = {
        id: row.public_id,
        name: row.name,
        lang: row.lang,
        tags: row.tags,
        enabled: row.enabled,
        triggers: row.triggers,
        messages: row.messages,
      };
      nodes.set(row.public_id, node);
    }

    return nodes;
  }
}
