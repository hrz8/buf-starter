import type { MessageInitShape } from '@bufbuild/protobuf';
import type { GenMessage } from '@bufbuild/protobuf/codegenv2';
import type { WidgetConfig } from '~~/gen/chatbot/modules/v1/widget_pb';

export type WidgetConfigInit = MessageInitShape<GenMessage<WidgetConfig>>;

export const widgetDefaults: WidgetConfigInit = {
  enabled: false,
  cors: {
    allowedOrigins: [],
    allowedHeaders: ['Content-Type'],
    credentials: false,
  },
};
