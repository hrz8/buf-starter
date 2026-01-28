import type { Server } from './types.js';
import { BotFactory } from './bot/factory.js';
import { AppConfig } from './config/config.js';
import { httpServer } from './http.js';

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
        console.error(`cleaning up error ${s.name}:`, err),
      ),
    ),
  );
  console.info('✨ cleanup done');
}

async function main() {
  const config = new AppConfig();

  await BotFactory.initialize(config);
  await startServers();
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
