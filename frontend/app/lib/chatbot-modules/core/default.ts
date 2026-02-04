import type { MessageInitShape } from '@bufbuild/protobuf';
import type { GenMessage } from '@bufbuild/protobuf/codegenv2';
import type { CoreConfig } from '~~/gen/chatbot/modules/v1/core_pb';

export type CoreConfigInit = MessageInitShape<GenMessage<CoreConfig>>;

export const DEFAULT_CORE_CONFIG: CoreConfigInit = {
  enabled: true,
  session: {
    ttlSeconds: 1800,
    cleanupIntervalSeconds: 300,
    defaultMode: 'assistant',
  },
  streaming: {
    textDelayMs: 15,
    chunkSize: 5,
    typingEffectEnabled: true,
  },
  flow: {
    maxGotoDepth: 10,
    useFallbackWhenNoMatch: true,
  },
};
