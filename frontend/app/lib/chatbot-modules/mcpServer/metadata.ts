import type { ModuleMetadata } from '../types';

export const mcpServerMetadata: ModuleMetadata = {
  key: 'mcpServer',
  title: 'MCP Server',
  description: 'Configure Model Context Protocol server connections for tool access.',
  icon: 'Server',
  fieldOrder: ['enabled', 'urls', 'structuredOutputs'],
  fields: {
    enabled: {},
    urls: {
      title: 'MCP Server URLs',
      titleKey: 'name',
      items: {
        fieldOrder: ['name', 'url', 'apiKey'],
        properties: {
          name: {
            placeholder: 'booking',
          },
          url: {
            title: 'URL',
            placeholder: 'https://mcp-server.example.com/mcp',
          },
          apiKey: {
            title: 'API Key',
            placeholder: 'sk-...',
          },
        },
      },
    },
    structuredOutputs: {
      titleKey: 'target',
      items: {
        fieldOrder: ['target', 'model', 'prompt', 'inputType', 'outputSchema', 'ui'],
        properties: {
          target: {
            placeholder: 'booking__search_flights',
          },
          model: {
            placeholder: 'openai.gpt-4o-mini',
          },
          prompt: {
            format: 'textarea',
            placeholder: 'Extract the data according to the schema...',
          },
          inputType: {
            enum: ['string', 'json'],
          },
          outputSchema: {
            additionalTypeInfo: 'json',
            placeholder: '{\n  "type": "object",\n  "properties": {}\n}',
          },
          ui: {
            title: 'UI Configuration',
            fieldOrder: ['enabled', 'component', 'isArray', 'dataPath', 'itemConfig'],
            properties: {
              enabled: {
                title: 'Enable UI',
              },
              component: {
                enum: ['carousel', 'card'],
              },
              isArray: {},
              dataPath: {
                placeholder: 'flights',
              },
              itemConfig: {
                additionalTypeInfo: 'json',
              },
            },
          },
        },
      },
    },
  },
};
