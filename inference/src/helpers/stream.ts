import type { Response } from 'express';

export function streamStaticMessage(res: Response, message: string): void {
  res.writeHead(200, {
    'Content-Type': 'text/plain; charset=utf-8',
    'Cache-Control': 'no-cache',
    'Connection': 'keep-alive',
  });

  res.write('data: {"type":"start"}\n\n');
  res.write('data: {"type":"start-step"}\n\n');
  res.write('data: {"type":"text-start","id":"0"}\n\n');

  const chunkSize = 5;
  for (let i = 0; i < message.length; i += chunkSize) {
    const chunk = message.slice(i, i + chunkSize);
    const escapedChunk = JSON.stringify(chunk).slice(1, -1);
    res.write(`data: {"type":"text-delta","id":"0","delta":"${escapedChunk}"}\n\n`);
  }

  res.write('data: {"type":"text-end","id":"0"}\n\n');
  res.write('data: {"type":"finish-step"}\n\n');
  res.write('data: {"type":"finish","finishReason":"stop"}\n\n');
  res.write('data: [DONE]\n\n');
  res.end();
}
