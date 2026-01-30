import type { ModelMessage, ToolSet, UIMessage } from 'ai';
import type { Response } from 'express';
import type {
  BlueprintData,
  BotState,
  BotStatus,
  ChatbotNode,
  ILlmProvider,
  IMcpClient,
  ModuleConfigMap,
  ModuleName,
  PlainMessage,
  PromptConfig,
} from './types.js';
import { convertToModelMessages, stepCountIs, streamText } from 'ai';
import { getRandomFallbackResponse } from '../../helpers/fallback.js';
import { streamStaticMessage } from '../../helpers/stream.js';
import { getModuleSetup } from '../modules/index.js';
import { NodeMatcher } from '../node/index.js';
import { Blueprint } from './blueprint.js';
import { MODULE_SETUP_ORDER } from './constants.js';

export class Bot {
  public readonly projectId: string;
  public projectName: string | null = null;

  public modules = new Map<ModuleName, ModuleConfigMap[ModuleName]>();
  public nodes = new Map<string, PlainMessage<ChatbotNode>>();

  public llmProvider: ILlmProvider | null = null;
  public mcpClient: IMcpClient | null = null;

  private nodeMatcher: NodeMatcher | null = null;
  private state: BotState = {
    status: 'idle',
    error: null,
  };

  public constructor(projectId: string) {
    this.projectId = projectId;
  }

  public async load(): Promise<void> {
    this.setStatus('loading');

    const data = await Blueprint.load(this.projectId);
    if (!data) {
      this.setError(`failed to load blueprint for project: ${this.projectId}`);
      return;
    }

    this.projectName = data.projectName;
    console.info(`loading bot for '${data.projectName}'`);

    await this.loadBluePrint(data);
    this.setStatus('ready');
  }

  public async loadBluePrint(data: BlueprintData): Promise<void> {
    const { modules, nodes } = data;
    this.modules.clear();
    this.nodes.clear();

    for (const moduleId of MODULE_SETUP_ORDER) {
      const moduleConfig = modules[moduleId];
      if (moduleConfig?.enabled) {
        this.modules.set(moduleId, moduleConfig);
        await this.setupModule(moduleId, moduleConfig);
        console.info(`module '${moduleId}' initialized`);
      }
    }

    const nodeEntries = nodes instanceof Map
      ? nodes.entries()
      : Object.entries(nodes);
    for (const [nodeId, node] of nodeEntries) {
      this.nodes.set(nodeId, node);
      console.info(`node '${node.name}' loaded`);
    }

    this.nodeMatcher = new NodeMatcher(this.nodes);
  }

  private async setupModule<K extends ModuleName>(
    moduleId: K,
    config: ModuleConfigMap[K],
  ): Promise<void> {
    const setupFn = getModuleSetup(moduleId);
    if (setupFn) {
      await setupFn(this, config);
    }
  }

  public async chat(
    messages: UIMessage[],
    res: Response,
    options?: { lang?: string },
  ): Promise<void> {
    const modelMessages = await convertToModelMessages(messages);
    const lastUserContent = this.getLastUserMessage(modelMessages);

    if (!this.llmProvider) {
      await this.respondWithoutLlm(lastUserContent, res, options?.lang);
      return;
    }

    await this.streamLlmResponse(modelMessages, res);
  }

  private async respondWithoutLlm(
    userInput: string,
    res: Response,
    lang?: string,
  ): Promise<void> {
    if (this.nodeMatcher) {
      const matchedNode = this.nodeMatcher.findMatch(userInput, lang);

      if (matchedNode) {
        const nodeResponse = NodeMatcher.getResponse(matchedNode);
        streamStaticMessage(res, nodeResponse);
        return;
      }
    }

    streamStaticMessage(res, getRandomFallbackResponse());
  }

  private async streamLlmResponse(
    messages: ModelMessage[],
    res: Response,
  ): Promise<void> {
    if (!this.llmProvider) {
      throw new Error('LLM provider not initialized');
    }

    let mcpTools: ToolSet | undefined;
    if (this.mcpClient?.isEnabled()) {
      try {
        mcpTools = await this.mcpClient.getTools();
        console.info(`[bot] loaded ${Object.keys(mcpTools).length} MCP tools`);
      } catch (error) {
        console.warn('[bot] failed to load MCP tools:', error);
      }
    }

    const model = this.llmProvider.getModel();
    const { systemPrompt, temperature, maxSteps } = this.llmProvider;
    const tools = mcpTools && Object.keys(mcpTools).length > 0 ? mcpTools : undefined;

    const result = streamText({
      model,
      system: systemPrompt,
      messages,
      tools,
      temperature,
      stopWhen: stepCountIs(maxSteps),
      onFinish: async () => {
        if (this.mcpClient) {
          await this.mcpClient.close();
        }
      },
    });

    result.pipeUIMessageStreamToResponse(res);
  }

  public isModuleEnabled(moduleId: ModuleName): boolean {
    const module = this.modules.get(moduleId);
    return module?.enabled ?? false;
  }

  public getSystemPrompt(): string {
    const promptConfig = this.modules.get('prompt') as PlainMessage<PromptConfig> | undefined;
    return promptConfig?.systemPrompt ?? 'You are a helpful assistant.';
  }

  public getState(): BotState {
    return { ...this.state };
  }

  public isReady(): boolean {
    return this.state.status === 'ready';
  }

  private setStatus(status: BotStatus): void {
    this.state.status = status;
    if (status !== 'error') {
      this.state.error = null;
    }
  }

  private setError(message: string): void {
    this.state.status = 'error';
    this.state.error = message;
  }

  private getLastUserMessage(messages: ModelMessage[]): string {
    for (let i = messages.length - 1; i >= 0; i--) {
      const msg = messages[i];
      if (msg && msg.role === 'user') {
        if (typeof msg.content === 'string') {
          return msg.content;
        }
        if (Array.isArray(msg.content)) {
          const textPart = msg.content.find(p => p.type === 'text');
          if (textPart && textPart.type === 'text') {
            return textPart.text;
          }
        }
      }
    }
    return '';
  }
}
