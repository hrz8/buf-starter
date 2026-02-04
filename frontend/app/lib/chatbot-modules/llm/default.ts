import type { MessageInitShape } from '@bufbuild/protobuf';
import type { GenMessage } from '@bufbuild/protobuf/codegenv2';
import type { LlmConfig } from '~~/gen/chatbot/modules/v1/llm_pb';

export type LlmConfigInit = MessageInitShape<GenMessage<LlmConfig>>;

export const DEFAULT_LLM_CONFIG: LlmConfigInit = {
  enabled: false,
  sdk: 'ai-sdk',
  provider: 'bedrock',
  model: '',
  temperature: 0.7,
  maxSteps: 5,
};
