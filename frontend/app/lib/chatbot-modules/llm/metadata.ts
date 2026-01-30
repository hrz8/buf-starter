import type { ModuleMetadata } from '../types';

export const llmMetadata: ModuleMetadata = {
  key: 'llm',
  title: 'LLM',
  description: 'Configure the Large Language Model settings for your chatbot.',
  icon: 'Brain',
  fieldOrder: ['enabled', 'sdk', 'provider', 'model', 'temperature', 'maxSteps'],
  fields: {
    enabled: {},
    sdk: {
      title: 'SDK',
      enum: ['ai-sdk'],
      enumLabels: {
        'ai-sdk': 'Vercel AI SDK',
      },
    },
    provider: {
      title: 'Provider',
      enum: ['bedrock'],
      enumLabels: {
        bedrock: 'AWS Bedrock',
      },
    },
    model: {
      placeholder: 'e.g., us.anthropic.claude-sonnet-4-20250514-v1:0',
    },
    temperature: {
      step: 0.1,
    },
    maxSteps: {
      title: 'Max Steps',
      step: 1,
    },
  },
};
