import type { Bot } from '../../bot/bot.js';
import type { PlainMessage, WidgetConfig } from '../../bot/types.js';

/**
 * Widget Module
 *
 * Manages embeddable chat widget configuration and CORS settings.
 * The widget config is stored and used by HTTP handlers for CORS validation.
 */
export function setup(bot: Bot, config: PlainMessage<WidgetConfig>): void {
  const allowedOrigins = config.cors?.allowedOrigins ?? [];
  console.info(`widget module initialized: ${allowedOrigins.length} allowed origin(s)`);
  // The widget config is already stored in bot.modules via loadModule
  // HTTP handlers should read it from there for CORS validation
}
