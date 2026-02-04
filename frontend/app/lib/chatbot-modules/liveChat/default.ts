import type { MessageInitShape } from '@bufbuild/protobuf';
import type { GenMessage } from '@bufbuild/protobuf/codegenv2';
import type { LiveChatConfig } from '~~/gen/chatbot/modules/v1/live_chat_pb';

export type LiveChatConfigInit = MessageInitShape<GenMessage<LiveChatConfig>>;

export const DEFAULT_LIVECHAT_CONFIG: LiveChatConfigInit = {
  enabled: false,
  queue: {
    maxSize: 100,
    estimatedWaitPerPositionSeconds: 120, // 2 minutes
  },
  agent: {
    assignmentStrategy: 'round-robin',
    maxConcurrentChats: 5,
  },
};
