import type { ModuleMetadata } from '../types';

export const widgetMetadata: ModuleMetadata = {
  key: 'widget',
  title: 'Widget',
  description: 'Configure the embeddable chat widget settings.',
  icon: 'MessageSquare',
  fieldOrder: ['enabled', 'cors'],
  fields: {
    enabled: {},
    cors: {
      title: 'CORS Settings',
      fieldOrder: ['allowedOrigins', 'allowedHeaders', 'credentials'],
      properties: {
        allowedOrigins: {
          items: {
            placeholder: 'https://example.com',
          },
        },
        allowedHeaders: {
          items: {},
        },
        credentials: {},
      },
    },
  },
};
