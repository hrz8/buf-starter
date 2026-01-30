import type { Bot } from '../../bot/bot.js';
import type { LlmConfig, PlainMessage } from '../../bot/types.js';
import { LlmProvider } from './provider.js';

export { LlmProvider } from './provider.js';

export function setup(bot: Bot, config: PlainMessage<LlmConfig>): void {
  const systemPrompt = bot.getSystemPrompt();

  const provider = new LlmProvider(config, systemPrompt);
  bot.llmProvider = provider;
}
