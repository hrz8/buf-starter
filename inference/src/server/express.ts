import cors from 'cors';
import express from 'express';
import { ERROR_CODE } from '../errors/errors.js';

const app = express();

app.use(express.json());
app.use(cors());

app.get('/healthz', (req, res) => {
  res.status(200).json({
    result: {
      status: 'ok',
    },
  });
});

app.use((req, res) => {
  res.status(404).json({
    code: ERROR_CODE.NOT_FOUND,
    error: 'Resource not found',
  });
});

app.use((
  err: unknown,
  req: express.Request,
  res: express.Response,
  _: express.NextFunction,
) => {
  console.error('Unhandled error:', err);
  return res.status(500).json({
    code: ERROR_CODE.INTERNAL_SERVER_ERROR,
    message: 'Internal server error',
  });
});

export { app };
