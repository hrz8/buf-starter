import { Bot } from './bot.js';

class BotProxy {
  get(target: Bot, prop: keyof Bot) {
    return target[prop];
  }
}

interface Config {
  getDefaultProjectId: () => string;
}

function proxifyBot(bot: Bot): Bot {
  return new Proxy(bot, new BotProxy());
}

async function createAndLoadBot(projectId: string): Promise<Bot> {
  const bot = new Bot(projectId);
  await bot.load();
  return proxifyBot(bot);
}

export class BotFactory {
  private static bots: Bot[] = [];

  public static async initialize(config: Config): Promise<Bot> {
    const bots = [await createAndLoadBot(config.getDefaultProjectId())];
    this.bots.push(...bots);
    return proxifyBot(this.bots[0]!);
  }

  public static async initializeForProject(projectId: string): Promise<Bot> {
    const bot = await createAndLoadBot(projectId);
    this.bots.push(bot);
    return proxifyBot(bot);
  }

  public static getDefaultBot(): Bot {
    if (this.bots.length === 0) {
      throw new Error('BotFactory not initialized.');
    }

    return proxifyBot(this.bots[0]!);
  }

  public static getBotByProjectId(projectId: string): Bot | null {
    const bot = this.bots.find(bot => bot.projectId === projectId);
    if (!bot) {
      return null;
    }
    return proxifyBot(bot);
  }
}
