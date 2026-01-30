import cors from 'cors';
import express from 'express';
import { ERROR_CODE } from '../errors/errors.js';
import { chatHandler } from '../handlers/chat/index.js';

const app = express();

app.use(express.json());
app.use(cors());

// Health check endpoint
app.get('/healthz', (_req, res) => {
  res.status(200).json({
    result: {
      status: 'ok',
    },
  });
});

// Chat endpoint - stateless, accepts message history from client
app.post('/chat', chatHandler);

// 404 handler
app.use((_req, res) => {
  res.status(404).json({
    code: ERROR_CODE.NOT_FOUND,
    error: 'Resource not found',
  });
});

// Error handler
app.use((
  err: unknown,
  _req: express.Request,
  res: express.Response,
  _next: express.NextFunction,
) => {
  console.error('Unhandled error:', err);
  return res.status(500).json({
    code: ERROR_CODE.INTERNAL_SERVER_ERROR,
    message: 'Internal server error',
  });
});

export { app };
