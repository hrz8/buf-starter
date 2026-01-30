import type { Server } from './types.js';
import { Blueprint } from './app/bot/blueprint.js';
import { BotFactory } from './app/bot/factory.js';
import { AppConfig } from './config/config.js';
import { httpServer } from './http.js';
import { BotRepository } from './repositories/db/bot.js';
import { Database } from './repositories/db/connection.js';
import { ProjectRepository } from './repositories/db/project.js';

const servers: Server[] = [];
servers.push(httpServer);

async function startServers() {
  await Promise.all(servers.map(s => s.start()));
}

async function stopServers() {
  console.info('🍀 performing graceful shutdown...');
  await Promise.all(
    servers.map(s =>
      s.stop().catch(err =>
        console.error(`server cleaning up error ${s.name}:`, err),
      ),
    ),
  );
  await Database.close();
  console.info('✨ cleanup done');
}

async function main() {
  const config = new AppConfig();

  const db = Database.initialize(config);

  const projectRepo = new ProjectRepository(db);
  const botRepo = new BotRepository(db);

  Blueprint.setup(projectRepo, botRepo);

  await BotFactory.initialize(config);

  await startServers();

  const bot = BotFactory.getDefaultBot();
  console.info(`\n--- Bot '${bot.projectName}' ready ---`);
  console.info(`status: ${bot.getState().status}`);
  console.info(`modules: ${Array.from(bot.modules.keys()).join(', ') || '(none)'}`);
  console.info(`nodes: ${bot.nodes.size}`);
}

async function cleanup() {
  await stopServers();
  process.exit(0);
}

process.on('SIGTERM', cleanup);
process.on('SIGINT', cleanup);

process.on('uncaughtException', (error) => {
  console.error('Uncaught Exception:', error);
  cleanup();
});

process.on('unhandledRejection', (reason, promise) => {
  console.error('Unhandled Rejection at:', promise, 'reason:', reason);
  cleanup();
});

main().catch((error) => {
  console.error('failed on startup', error);
  process.exit(1);
});
