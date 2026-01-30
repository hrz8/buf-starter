import type {
  IMcpClient,
  McpServerConfig,
  McpServerUrl,
  PlainMessage,
} from '../../bot/types.js';
import { experimental_createMCPClient as createMCPClient } from '@ai-sdk/mcp';

type McpClientInstance = Awaited<ReturnType<typeof createMCPClient>>;
type McpToolsRecord = Awaited<ReturnType<McpClientInstance['tools']>>;

export class McpClient implements IMcpClient {
  private config: PlainMessage<McpServerConfig>;
  private clients: McpClientInstance[] = [];
  private toolsCache: McpToolsRecord | null = null;
  private isConnected = false;

  public constructor(config: PlainMessage<McpServerConfig>) {
    this.config = config;
  }

  public getServerUrls(): PlainMessage<McpServerUrl>[] {
    return this.config.urls ?? [];
  }

  public isEnabled(): boolean {
    return this.config.enabled && (this.config.urls?.length ?? 0) > 0;
  }

  public async connect(): Promise<void> {
    if (this.isConnected) {
      return;
    }

    const urls = this.config.urls ?? [];
    if (urls.length === 0) {
      return;
    }

    const wrappedTools: McpToolsRecord = {};

    for (const serverConfig of urls) {
      try {
        const client = await this.connectToServer(serverConfig);
        this.clients.push(client);

        const tools = await client.tools();
        for (const [toolName, tool] of Object.entries(tools)) {
          const prefixedName = `${serverConfig.name}__${toolName}`;
          wrappedTools[prefixedName] = tool;
        }

        console.info(`[mcp] loaded ${Object.keys(tools).length} tools from: ${serverConfig.name}`);
      } catch (error) {
        console.error(`[mcp] failed to connect to ${serverConfig.name}:`, error);
      }
    }

    this.toolsCache = wrappedTools;
    this.isConnected = true;
    console.info(`[mcp] total ${Object.keys(wrappedTools).length} tools loaded`);
  }

  public async getTools(): Promise<McpToolsRecord> {
    if (!this.isConnected) {
      await this.connect();
    }
    return this.toolsCache ?? {};
  }

  public hasTools(): boolean {
    return this.toolsCache !== null && Object.keys(this.toolsCache).length > 0;
  }

  public async close(): Promise<void> {
    for (const client of this.clients) {
      try {
        await client.close();
      } catch (error) {
        console.error('[mcp] error closing client:', error);
      }
    }

    if (this.clients.length > 0) {
      console.info(`[mcp] closed ${this.clients.length} client(s)`);
    }

    this.clients = [];
    this.toolsCache = null;
    this.isConnected = false;
  }

  private async connectToServer(
    serverConfig: PlainMessage<McpServerUrl>,
  ): Promise<McpClientInstance> {
    console.info(`[mcp] connecting to: ${serverConfig.name} at ${serverConfig.url}`);

    const headers: Record<string, string> = {};
    if (serverConfig.apiKey) {
      headers.Authorization = `Bearer ${serverConfig.apiKey}`;
    }

    return createMCPClient({
      transport: {
        type: 'http',
        url: serverConfig.url,
        ...(Object.keys(headers).length > 0 && { headers }),
      },
    });
  }
}
