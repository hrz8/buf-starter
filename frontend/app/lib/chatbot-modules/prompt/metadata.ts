import type { ModuleMetadata } from '../types';

export const promptMetadata: ModuleMetadata = {
  key: 'prompt',
  title: 'Prompt',
  description: 'Configure the system prompt and conversation behavior.',
  icon: 'FileText',
  fieldOrder: ['enabled', 'systemPrompt'],
  fields: {
    enabled: {},
    systemPrompt: {
      format: 'textarea',
      placeholder: 'You are a helpful assistant...',
    },
  },
};
