<script setup lang="ts">
import type { ChatbotNode } from '~~/gen/chatbot/nodes/v1/node_pb';
import { toTypedSchema } from '@vee-validate/zod';
import { useForm } from 'vee-validate';
import { toast } from 'vue-sonner';
import { Button } from '@/components/ui/button';
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { useChatbotNodeService } from '@/composables/services/useChatbotNodeService';
import { useChatbotNodeStore } from '@/stores/chatbot-node';
import { useProjectStore } from '@/stores/project';
import { getConnectRPCError, hasConnectRPCError, LANGUAGES, nodeCreateSchema } from '.';

const emit = defineEmits<{
  success: [node: ChatbotNode];
  cancel: [];
}>();

const { t } = useI18n();
const projectStore = useProjectStore();
const nodeStore = useChatbotNodeStore();
const nodeService = useChatbotNodeService();

const form = useForm({
  validationSchema: toTypedSchema(nodeCreateSchema),
  initialValues: {
    name: '',
    lang: 'en-US',
    tags: [],
  },
});

const isSubmitting = computed(() => nodeService.createLoading.value);

const onSubmit = form.handleSubmit(async (values) => {
  const projectId = projectStore.activeProjectId;
  if (!projectId) {
    toast.error(t('common.error'), {
      description: t('features.chatbotNode.messages.noProject'),
    });
    return;
  }

  try {
    const node = await nodeService.createNode(
      projectId,
      values.name,
      values.lang,
      values.tags || [],
    );

    if (node) {
      nodeStore.addNode(node);
      toast.success(t('common.success'), {
        description: t('features.chatbotNode.messages.createSuccess'),
      });
      emit('success', node);
    }
  }
  catch {
    toast.error(t('common.error'), {
      description: nodeService.createError.value || t('features.chatbotNode.messages.createError'),
    });
  }
});

function handleCancel() {
  form.resetForm();
  nodeService.resetCreateState();
  emit('cancel');
}
</script>

<template>
  <form class="space-y-6" @submit="onSubmit">
    <FormField v-slot="{ componentField }" name="name">
      <FormItem>
        <FormLabel>{{ t('features.chatbotNode.form.name') }}</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            :placeholder="t('features.chatbotNode.form.namePlaceholder')"
            :disabled="isSubmitting"
          />
        </FormControl>
        <FormDescription>
          {{ t('features.chatbotNode.form.nameHelp') }}
        </FormDescription>
        <FormMessage />
        <p
          v-if="hasConnectRPCError(nodeService.createValidationErrors, 'name')"
          class="text-sm font-medium text-destructive"
        >
          {{ getConnectRPCError(nodeService.createValidationErrors, 'name') }}
        </p>
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="lang">
      <FormItem>
        <FormLabel>{{ t('features.chatbotNode.form.lang') }}</FormLabel>
        <Select v-bind="componentField" :disabled="isSubmitting">
          <FormControl>
            <SelectTrigger>
              <SelectValue :placeholder="t('features.chatbotNode.form.langPlaceholder')" />
            </SelectTrigger>
          </FormControl>
          <SelectContent>
            <SelectItem v-for="lang in LANGUAGES" :key="lang.value" :value="lang.value">
              {{ lang.label }}
            </SelectItem>
          </SelectContent>
        </Select>
        <FormMessage />
        <p
          v-if="hasConnectRPCError(nodeService.createValidationErrors, 'lang')"
          class="text-sm font-medium text-destructive"
        >
          {{ getConnectRPCError(nodeService.createValidationErrors, 'lang') }}
        </p>
      </FormItem>
    </FormField>

    <div class="flex justify-end gap-3 pt-4">
      <Button type="button" variant="outline" :disabled="isSubmitting" @click="handleCancel">
        {{ t('common.cancel') }}
      </Button>
      <Button type="submit" :disabled="isSubmitting">
        <span v-if="isSubmitting">{{ t('common.creating') }}</span>
        <span v-else>{{ t('features.chatbotNode.form.create') }}</span>
      </Button>
    </div>
  </form>
</template>
