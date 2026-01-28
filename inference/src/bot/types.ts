import type { Message } from '@bufbuild/protobuf';
import type { LlmConfig } from '../../gen/chatbot/modules/v1/llm_pb.js';
import type { McpServerConfig } from '../../gen/chatbot/modules/v1/mcp_server_pb.js';
import type { PromptConfig } from '../../gen/chatbot/modules/v1/prompt_pb.js';
import type { WidgetConfig } from '../../gen/chatbot/modules/v1/widget_pb.js';
import type { ChatbotNode } from '../../gen/chatbot/nodes/v1/node_pb.js';

export type { LlmConfig } from '../../gen/chatbot/modules/v1/llm_pb.js';
export type {
  McpServerConfig,
  McpServerUrl,
  StructuredOutput,
  StructuredOutputUi,
} from '../../gen/chatbot/modules/v1/mcp_server_pb.js';
export type { PromptConfig } from '../../gen/chatbot/modules/v1/prompt_pb.js';
export type {
  WidgetConfig,
  WidgetCors,
} from '../../gen/chatbot/modules/v1/widget_pb.js';
export type {
  ChatbotNode,
  ChatbotNodeMessage,
  ChatbotNodeTrigger,
} from '../../gen/chatbot/nodes/v1/node_pb.js';

type PlainMessage<T> = T extends Message<infer _> ? Omit<T, keyof Message<string>> : T;

export type ModuleName = 'llm' | 'mcpServer' | 'prompt' | 'widget';

export interface ModuleConfigMap {
  llm: PlainMessage<LlmConfig>;
  mcpServer: PlainMessage<McpServerConfig>;
  prompt: PlainMessage<PromptConfig>;
  widget: PlainMessage<WidgetConfig>;
}

export type ModuleConfig = ModuleConfigMap[ModuleName];

export interface BlueprintData {
  projectName: string;
  modules: Partial<ModuleConfigMap>;
  nodes: Map<string, PlainMessage<ChatbotNode>> | Record<string, PlainMessage<ChatbotNode>>;
}

export type { PlainMessage };
