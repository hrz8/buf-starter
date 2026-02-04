<script setup lang="ts">
import type { ChatbotNode } from '~~/gen/chatbot/nodes/v1/node_pb';
import { toTypedSchema } from '@vee-validate/zod';
import { useForm } from 'vee-validate';
import { toast } from 'vue-sonner';
import { z } from 'zod';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { useChatbotNodeService } from '@/composables/services/useChatbotNodeService';
import { useChatbotNodeStore } from '@/stores/chatbot-node';
import { useProjectStore } from '@/stores/project';
import { getConnectRPCError, hasConnectRPCError } from '.';

interface Props {
  open?: boolean;
  nodeName: string;
  nodeLang: string;
}

const props = withDefaults(defineProps<Props>(), {
  open: false,
});

const emit = defineEmits<{
  'success': [node: ChatbotNode];
  'cancel': [];
  'update:open': [open: boolean];
}>();

const { t } = useI18n();
const projectStore = useProjectStore();
const nodeStore = useChatbotNodeStore();
const nodeService = useChatbotNodeService();

// Version-only schema (name and lang are fixed)
const versionSchema = z.object({
  version: z.string().min(1, 'Version is required').regex(
    /^[a-z][a-z0-9_]*$/,
    'Must be lowercase with letters, numbers, and underscores',
  ),
});

const form = useForm({
  validationSchema: toTypedSchema(versionSchema),
  initialValues: {
    version: '',
  },
});

const isOpen = computed({
  get: () => props.open,
  set: (value: boolean) => emit('update:open', value),
});

const isSubmitting = computed(() => nodeService.createLoading.value);

// Display name combining name and lang
const nodeDisplayName = computed(() => `${props.nodeName}_${props.nodeLang}`);

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
      props.nodeName,
      props.nodeLang,
      [], // Empty tags for new version
      values.version,
    );

    if (node) {
      nodeStore.addNode(node);
      toast.success(t('common.success'), {
        description: t('features.chatbotNode.messages.createSuccess'),
      });
      form.resetForm();
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
  isOpen.value = false;
}

// Reset form when dialog opens
watch(isOpen, (open) => {
  if (open) {
    form.resetForm();
    nodeService.resetCreateState();
  }
});
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogContent class="sm:max-w-[425px]">
      <DialogHeader>
        <DialogTitle>{{ t('features.chatbotNode.addVersion') }}</DialogTitle>
        <DialogDescription>
          {{ t('features.chatbotNode.addVersionDesc', { name: nodeDisplayName }) }}
        </DialogDescription>
      </DialogHeader>

      <form class="space-y-4 py-4" @submit="onSubmit">
        <!-- Node Name (readonly) -->
        <div class="space-y-2">
          <label class="text-sm font-medium leading-none">
            {{ t('features.chatbotNode.form.name') }}
          </label>
          <Input :model-value="props.nodeName" disabled class="bg-muted" />
        </div>

        <!-- Language (readonly) -->
        <div class="space-y-2">
          <label class="text-sm font-medium leading-none">
            {{ t('features.chatbotNode.form.lang') }}
          </label>
          <Input :model-value="props.nodeLang" disabled class="bg-muted" />
        </div>

        <!-- Version (required) -->
        <FormField v-slot="{ componentField }" name="version">
          <FormItem>
            <FormLabel>{{ t('features.chatbotNode.form.version') }} *</FormLabel>
            <FormControl>
              <Input
                v-bind="componentField"
                :placeholder="t('features.chatbotNode.form.versionPlaceholder')"
                :disabled="isSubmitting"
              />
            </FormControl>
            <FormDescription>
              {{ t('features.chatbotNode.form.versionHelp') }}
            </FormDescription>
            <FormMessage />
            <p
              v-if="hasConnectRPCError(nodeService.createValidationErrors, 'version')"
              class="text-sm font-medium text-destructive"
            >
              {{ getConnectRPCError(nodeService.createValidationErrors, 'version') }}
            </p>
          </FormItem>
        </FormField>

        <DialogFooter class="pt-4">
          <Button type="button" variant="outline" :disabled="isSubmitting" @click="handleCancel">
            {{ t('common.cancel') }}
          </Button>
          <Button type="submit" :disabled="isSubmitting">
            <span v-if="isSubmitting">{{ t('common.creating') }}</span>
            <span v-else>{{ t('features.chatbotNode.addVersion') }}</span>
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
