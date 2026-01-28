import type { Server as HttpServer } from 'node:http';
import type { Server } from './types.js';
import { app } from './server/express.js';

const PORT = process.env.PORT || 3400;

export const httpServer: Server<HttpServer> = {
  name: 'HTTP Server',
  instance: null,
  async start() {
    return new Promise((resolve) => {
      const server = app.listen(PORT, () => {
        console.info(`🚀 starting HTTP server at port: ${PORT}`);
        this.instance = server;
        resolve(server);
      });
    });
  },
  async stop() {
    if (this.instance) {
      return new Promise((resolve, reject) => {
        this.instance!.close((err) => {
          if (err) {
            console.error(`⚠️ failed to stop ${this.name}:`, err);
            reject(err);
          } else {
            resolve();
          }
        });
      });
    }
  },
};
