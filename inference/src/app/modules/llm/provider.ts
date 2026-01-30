import type { LanguageModel } from 'ai';
import type { ILlmProvider, LlmConfig, PlainMessage } from '../../bot/types.js';
import { createAmazonBedrock } from '@ai-sdk/amazon-bedrock';
import { fromNodeProviderChain } from '@aws-sdk/credential-providers';

export const LLM_SDKS = [
  { value: 'ai-sdk', label: 'Vercel AI SDK', isDefault: true },
] as const;

export type LlmSdk = typeof LLM_SDKS[number]['value'];

export const DEFAULT_SDK = LLM_SDKS.find(s => s.isDefault)!.value;

export const LLM_PROVIDERS = [
  { value: 'bedrock', label: 'AWS Bedrock', isDefault: true },
] as const;

export type LlmProviderType = typeof LLM_PROVIDERS[number]['value'];

export const DEFAULT_PROVIDER = LLM_PROVIDERS.find(p => p.isDefault)!.value;

export const DEFAULT_MODEL = 'us.anthropic.claude-sonnet-4-20250514-v1:0';
export const DEFAULT_SYSTEM_PROMPT = 'You are a helpful assistant.';

export class LlmProvider implements ILlmProvider {
  public readonly sdk: LlmSdk;
  public readonly providerType: LlmProviderType;
  public readonly modelId: string;
  public readonly temperature: number;
  public readonly maxSteps: number;
  public readonly systemPrompt: string;

  private readonly model: LanguageModel;

  public constructor(config: PlainMessage<LlmConfig>, baseSystemPrompt?: string) {
    const sdk = (config.sdk || DEFAULT_SDK) as string;
    const providerType = (config.provider || DEFAULT_PROVIDER) as string;

    this.modelId = config.model || DEFAULT_MODEL;
    this.temperature = config.temperature ?? 0.7;
    this.maxSteps = config.maxSteps ?? 10;

    this.sdk = this.validateSdk(sdk);
    this.providerType = this.validateProvider(providerType);

    this.systemPrompt = this.buildSystemPrompt(baseSystemPrompt || DEFAULT_SYSTEM_PROMPT);

    this.model = this.createModel();
  }

  public getModel(): LanguageModel {
    return this.model;
  }

  private createModel(): LanguageModel {
    switch (this.sdk) {
      case 'ai-sdk':
        return this.createAiSdkModel();

      default:
        return this.createAiSdkModel();
    }
  }

  private createAiSdkModel(): LanguageModel {
    switch (this.providerType) {
      case 'bedrock':
        return this.createBedrockModel();

      default:
        return this.createBedrockModel();
    }
  }

  private createBedrockModel(): LanguageModel {
    const bedrock = createAmazonBedrock({
      region: process.env.AWS_REGION ?? 'us-east-1',
      credentialProvider: fromNodeProviderChain(),
    });
    return bedrock(this.modelId);
  }

  private buildSystemPrompt(basePrompt: string): string {
    const currentDate = new Date();

    const dateStr = currentDate.toLocaleDateString('en-US', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });

    const timeStr = currentDate.toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      timeZoneName: 'short',
    });

    const isoDateStr = currentDate.toISOString().split('T')[0];

    return `${basePrompt}

CURRENT DATE & TIME:
- Today's date: ${dateStr} (${isoDateStr})
- Current time: ${timeStr}`;
  }

  private validateSdk(sdk: string): LlmSdk {
    const valid = LLM_SDKS.find(s => s.value === sdk);
    return valid ? valid.value : DEFAULT_SDK;
  }

  private validateProvider(provider: string): LlmProviderType {
    const valid = LLM_PROVIDERS.find(p => p.value === provider);
    return valid ? valid.value : DEFAULT_PROVIDER;
  }
}
