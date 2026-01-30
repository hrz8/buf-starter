import type { Message } from '@bufbuild/protobuf';
import type { LanguageModel, ToolSet } from 'ai';
import type { LlmConfig } from '../../../gen/chatbot/modules/v1/llm_pb.js';
import type { McpServerConfig } from '../../../gen/chatbot/modules/v1/mcp_server_pb.js';
import type { PromptConfig } from '../../../gen/chatbot/modules/v1/prompt_pb.js';
import type { WidgetConfig } from '../../../gen/chatbot/modules/v1/widget_pb.js';
import type { ChatbotNode } from '../../../gen/chatbot/nodes/v1/node_pb.js';

export type { LlmConfig } from '../../../gen/chatbot/modules/v1/llm_pb.js';
export type {
  McpServerConfig,
  McpServerUrl,
  StructuredOutput,
  StructuredOutputUi,
} from '../../../gen/chatbot/modules/v1/mcp_server_pb.js';
export type { PromptConfig } from '../../../gen/chatbot/modules/v1/prompt_pb.js';
export type {
  WidgetConfig,
  WidgetCors,
} from '../../../gen/chatbot/modules/v1/widget_pb.js';
export type {
  ChatbotNode,
  ChatbotNodeMessage,
  ChatbotNodeTrigger,
} from '../../../gen/chatbot/nodes/v1/node_pb.js';

export type PlainMessage<T> = T extends Message<infer _> ? Omit<T, keyof Message<string>> : T;

export type ModuleName = 'prompt' | 'mcpServer' | 'llm' | 'widget';

export interface ModuleConfigMap {
  prompt: PlainMessage<PromptConfig>;
  mcpServer: PlainMessage<McpServerConfig>;
  llm: PlainMessage<LlmConfig>;
  widget: PlainMessage<WidgetConfig>;
}

export type ModuleConfig = ModuleConfigMap[ModuleName];

export interface BlueprintData {
  projectName: string;
  modules: Partial<ModuleConfigMap>;
  nodes: Map<string, PlainMessage<ChatbotNode>> | Record<string, PlainMessage<ChatbotNode>>;
}

export type BotStatus = 'idle' | 'loading' | 'ready' | 'error';

export interface BotState {
  status: BotStatus;
  error: string | null;
}

export interface ILlmProvider {
  readonly sdk: string;
  readonly providerType: string;
  readonly modelId: string;
  readonly temperature: number;
  readonly maxSteps: number;
  readonly systemPrompt: string;

  getModel: () => LanguageModel;
}

export interface IMcpClient {
  isEnabled: () => boolean;
  getTools: () => Promise<ToolSet>;
  hasTools: () => boolean;
  close: () => Promise<void>;
}
