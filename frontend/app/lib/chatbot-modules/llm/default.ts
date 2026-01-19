import type { MessageInitShape } from '@bufbuild/protobuf';
import type { GenMessage } from '@bufbuild/protobuf/codegenv2';
import type { LlmConfig } from '~~/gen/chatbot/modules/v1/llm_pb';

export type LlmConfigInit = MessageInitShape<GenMessage<LlmConfig>>;

export const llmDefaults: LlmConfigInit = {
  enabled: false,
  model: '',
  temperature: 0.7,
  maxToolCalls: 5,
};
