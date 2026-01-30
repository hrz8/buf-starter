import type { Bot } from '../../bot/bot.js';
import type { McpServerConfig, PlainMessage } from '../../bot/types.js';
import { McpClient } from './client.js';

export { McpClient } from './client.js';

export function setup(bot: Bot, config: PlainMessage<McpServerConfig>): void {
  const client = new McpClient(config);
  bot.mcpClient = client;
}
