import type { UIMessage } from 'ai';
import type { Request, Response } from 'express';
import { BotFactory } from '../../app/bot/factory.js';

interface SimpleMessage {
  role: 'user' | 'assistant' | 'system';
  content: string;
}

/**
 * Convert simple message format to UIMessage format
 */
function toUIMessages(messages: (SimpleMessage | UIMessage)[]): UIMessage[] {
  return messages.map((msg, index) => {
    // Already in UIMessage format
    if ('parts' in msg || 'id' in msg) {
      return {
        id: (msg as UIMessage).id || String(index),
        role: msg.role,
        parts: (msg as UIMessage).parts || [{ type: 'text', text: (msg as SimpleMessage).content || '' }],
      } as UIMessage;
    }

    // Simple format: { role, content }
    return {
      id: String(index),
      role: msg.role,
      parts: [{ type: 'text', text: msg.content }],
    } as UIMessage;
  });
}

/**
 * Chat Handler
 *
 * POST /chat endpoint handler.
 *
 * Flow:
 * 1. Get bot instance (default or by projectId)
 * 2. Validate request (messages required)
 * 3. Extract language from Accept-Language header
 * 4. Delegate to bot.chat() which handles:
 *    - If LLM disabled: node matching → fallback
 *    - If LLM enabled: MCP tools → LLM streaming
 *
 * Request body:
 *   { messages: UIMessage[] }
 *
 * Query params:
 *   - projectId: Project ID (uses default if not provided)
 *
 * Headers:
 *   - Accept-Language: Used for node language matching (e.g., "en-US", "id-ID")
 */
export async function chatHandler(req: Request, res: Response): Promise<void> {
  try {
    const projectId = req.query.projectId as string | undefined;
    const bot = projectId
      ? await BotFactory.getOrCreateBot(projectId)
      : BotFactory.getDefaultBot();

    if (!bot.isReady()) {
      res.status(503).json({ error: 'Bot is not ready', projectId });
      return;
    }

    const { messages } = req.body || {};

    if (!messages || !Array.isArray(messages) || messages.length === 0) {
      res.status(400).json({ error: 'No messages provided' });
      return;
    }

    const acceptLanguage = req.headers['accept-language'] ?? '';
    const lang = acceptLanguage.split(',')[0]?.trim() || undefined;

    const uiMessages = toUIMessages(messages);
    await bot.chat(uiMessages, res, { lang });
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unknown error';
    console.error('[chat] handler error:', message, error);

    if (!res.headersSent) {
      res.status(500).json({ error: 'Internal server error', message });
    }
  }
}
