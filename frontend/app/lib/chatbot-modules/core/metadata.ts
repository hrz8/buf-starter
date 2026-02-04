import type { ModuleMetadata } from '../types';

export const coreMetadata: ModuleMetadata = {
  key: 'core',
  title: 'Core Settings',
  description: 'Configure core runtime settings including session management, streaming behavior, and flow handling.',
  icon: 'Settings2',
  fieldOrder: ['enabled', 'session', 'streaming', 'flow'],
  fields: {
    enabled: { title: 'Enabled' },
    session: {
      title: 'Session Management',
      fieldOrder: ['ttlSeconds', 'cleanupIntervalSeconds', 'defaultMode'],
      properties: {
        ttlSeconds: {
          title: 'Session TTL (seconds)',
          placeholder: '1800 (30 minutes)',
          step: 60,
        },
        cleanupIntervalSeconds: {
          title: 'Cleanup Interval (seconds)',
          placeholder: '300 (5 minutes)',
          step: 30,
        },
        defaultMode: {
          title: 'Default Mode',
          enum: ['assistant', 'flow', 'liveChat'],
          enumLabels: {
            assistant: 'AI Assistant',
            flow: 'Static Flow',
            liveChat: 'Live Chat',
          },
        },
      },
    },
    streaming: {
      title: 'Streaming',
      fieldOrder: ['typingEffectEnabled', 'textDelayMs', 'chunkSize'],
      properties: {
        typingEffectEnabled: { title: 'Enable Typing Effect' },
        textDelayMs: { title: 'Delay Between Chunks (ms)', step: 5 },
        chunkSize: { title: 'Chunk Size (characters)', step: 1 },
      },
    },
    flow: {
      title: 'Flow Handler',
      fieldOrder: ['maxGotoDepth', 'useFallbackWhenNoMatch'],
      properties: {
        maxGotoDepth: { title: 'Max Goto Depth', step: 1 },
        useFallbackWhenNoMatch: { title: 'Use Fallback When No Match' },
      },
    },
  },
};
