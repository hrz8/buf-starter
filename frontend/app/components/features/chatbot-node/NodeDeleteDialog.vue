<script setup lang="ts">
import type { ChatbotNode } from '~~/gen/chatbot/nodes/v1/node_pb';
import { Trash2 } from 'lucide-vue-next';
import { toast } from 'vue-sonner';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import { Button } from '@/components/ui/button';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { useChatbotNodeService } from '@/composables/services/useChatbotNodeService';
import { useChatbotNodeStore } from '@/stores/chatbot-node';
import { useProjectStore } from '@/stores/project';

interface Props {
  node: ChatbotNode;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  deleted: [];
}>();

const { t } = useI18n();
const projectStore = useProjectStore();
const nodeStore = useChatbotNodeStore();
const nodeService = useChatbotNodeService();

const isOpen = ref(false);
const isDeleting = computed(() => nodeService.deleteLoading.value);

const nodeName = computed(() => `${props.node.name}_${props.node.lang}`);

// Predefined nodes cannot be deleted
const isPredefined = computed(() => props.node.isPredefined);

async function handleDelete() {
  const projectId = projectStore.activeProjectId;
  if (!projectId || isPredefined.value) {
    return;
  }

  try {
    const success = await nodeService.deleteNode(projectId, props.node.id);
    if (success) {
      nodeStore.removeNode(props.node.id);
      toast.success(t('common.success'), {
        description: t('features.chatbotNode.messages.deleteSuccess'),
      });
      isOpen.value = false;
      emit('deleted');
    }
  }
  catch {
    toast.error(t('common.error'), {
      description: nodeService.deleteError.value || t('features.chatbotNode.messages.deleteError'),
    });
  }
}
</script>

<template>
  <TooltipProvider v-if="isPredefined">
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="destructive" size="sm" disabled>
          <Trash2 class="h-4 w-4 mr-2" />
          {{ t('features.chatbotNode.delete.button') }}
        </Button>
      </TooltipTrigger>
      <TooltipContent>
        <p>{{ t('features.chatbotNode.delete.predefinedTooltip') }}</p>
      </TooltipContent>
    </Tooltip>
  </TooltipProvider>

  <AlertDialog v-else v-model:open="isOpen">
    <AlertDialogTrigger as-child>
      <Button variant="destructive" size="sm">
        <Trash2 class="h-4 w-4 mr-2" />
        {{ t('features.chatbotNode.delete.button') }}
      </Button>
    </AlertDialogTrigger>
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>{{ t('features.chatbotNode.delete.title') }}</AlertDialogTitle>
        <AlertDialogDescription>
          {{ t('features.chatbotNode.delete.description', { name: nodeName }) }}
        </AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel :disabled="isDeleting">
          {{ t('common.cancel') }}
        </AlertDialogCancel>
        <AlertDialogAction
          class="bg-destructive text-destructive-foreground hover:bg-destructive/90"
          :disabled="isDeleting"
          @click="handleDelete"
        >
          <span v-if="isDeleting">{{ t('common.deleting') }}</span>
          <span v-else>{{ t('features.chatbotNode.delete.confirm') }}</span>
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
