import type { ChatbotNode, PlainMessage } from '../bot/types.js';

export class NodeMatcher {
  private nodes: Map<string, PlainMessage<ChatbotNode>>;

  public constructor(nodes: Map<string, PlainMessage<ChatbotNode>>) {
    this.nodes = nodes;
  }

  public findMatch(userInput: string, lang?: string): PlainMessage<ChatbotNode> | undefined {
    const normalizedInput = userInput.toLowerCase().trim();

    for (const node of this.nodes.values()) {
      if (!node.enabled) {
        continue;
      }

      if (lang && node.lang && node.lang !== lang) {
        continue;
      }

      const matched = this.matchTriggers(node, normalizedInput, userInput);
      if (matched) {
        return node;
      }
    }

    return undefined;
  }

  private matchTriggers(
    node: PlainMessage<ChatbotNode>,
    normalizedInput: string,
    originalInput: string,
  ): boolean {
    for (const trigger of node.triggers ?? []) {
      const triggerValue = trigger.value.toLowerCase();

      switch (trigger.type) {
        case 'equals':
          if (normalizedInput === triggerValue) {
            return true;
          }
          break;

        case 'contains':
          if (normalizedInput.includes(triggerValue)) {
            return true;
          }
          break;

        case 'keyword': {
          const words = normalizedInput.split(/\s+/);
          if (words.includes(triggerValue)) {
            return true;
          }
          break;
        }

        case 'regex':
          try {
            const regex = new RegExp(trigger.value, 'i');
            if (regex.test(originalInput)) {
              return true;
            }
          } catch {
            console.warn(`[node] invalid regex pattern: ${trigger.value}`);
          }
          break;

        default:
          console.warn(`[node] unknown trigger type: ${trigger.type}`);
      }
    }

    return false;
  }

  public static getResponse(node: PlainMessage<ChatbotNode>): string {
    const messages = node.messages ?? [];
    if (messages.length === 0) {
      return '';
    }

    return messages
      .filter(m => m.role === 'assistant')
      .map(m => m.content)
      .join('\n\n');
  }
}
