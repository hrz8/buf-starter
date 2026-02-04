import type { MessageInitShape } from '@bufbuild/protobuf';
import type { GenMessage } from '@bufbuild/protobuf/codegenv2';
import type { PromptConfig } from '~~/gen/chatbot/modules/v1/prompt_pb';

export type PromptConfigInit = MessageInitShape<GenMessage<PromptConfig>>;

export const DEFAULT_PROMPT_CONFIG: PromptConfigInit = {
  enabled: true,
  systemPrompt: 'You are a helpful assistant.',
};
