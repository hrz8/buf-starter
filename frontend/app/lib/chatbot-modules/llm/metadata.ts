import type { ModuleMetadata } from '../types';

export const llmMetadata: ModuleMetadata = {
  key: 'llm',
  title: 'LLM',
  description: 'Configure the Large Language Model settings for your chatbot.',
  icon: 'Brain',
  fieldOrder: ['enabled', 'model', 'temperature', 'maxToolCalls'],
  fields: {
    enabled: {},
    model: {
      placeholder: 'e.g., gpt-4-turbo',
    },
    temperature: {
      step: 0.1,
    },
    maxToolCalls: {
      step: 1,
    },
  },
};
