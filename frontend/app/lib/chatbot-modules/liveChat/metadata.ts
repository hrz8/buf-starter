import type { ModuleMetadata } from '../types';

export const liveChatMetadata: ModuleMetadata = {
  key: 'liveChat',
  title: 'Live Chat',
  description: 'Configure live chat queue and agent assignment settings. (Coming Soon)',
  icon: 'Headphones',
  fieldOrder: ['enabled', 'queue', 'agent'],
  fields: {
    enabled: { title: 'Enable Live Chat' },
    queue: {
      title: 'Queue Settings',
      fieldOrder: ['maxSize', 'estimatedWaitPerPositionSeconds'],
      properties: {
        maxSize: { title: 'Max Queue Size', step: 10 },
        estimatedWaitPerPositionSeconds: {
          title: 'Est. Wait Per Position (seconds)',
          step: 10,
        },
      },
    },
    agent: {
      title: 'Agent Settings',
      fieldOrder: ['assignmentStrategy', 'maxConcurrentChats'],
      properties: {
        assignmentStrategy: {
          title: 'Assignment Strategy',
          enum: ['round-robin', 'least-busy', 'manual'],
          enumLabels: {
            'round-robin': 'Round Robin',
            'least-busy': 'Least Busy',
            'manual': 'Manual',
          },
        },
        maxConcurrentChats: { title: 'Max Concurrent Chats', step: 1 },
      },
    },
  },
};
