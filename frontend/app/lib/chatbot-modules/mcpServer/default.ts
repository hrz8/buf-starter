import type { MessageInitShape } from '@bufbuild/protobuf';
import type { GenMessage } from '@bufbuild/protobuf/codegenv2';
import type { McpServerConfig } from '~~/gen/chatbot/modules/v1/mcp_server_pb';

export type McpServerConfigInit = MessageInitShape<GenMessage<McpServerConfig>>;

export const mcpServerDefaults: McpServerConfigInit = {
  enabled: false,
  urls: [],
  structuredOutputs: [],
};
