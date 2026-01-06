<script setup lang="ts">
import { Loader2 } from 'lucide-vue-next';
import { computed, ref, watch } from 'vue';
import { toast } from 'vue-sonner';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { usePermissionService } from '@/composables/services/usePermissionService';

const props = defineProps<{
  open: boolean;
}>();

const emit = defineEmits<{
  'update:open': [value: boolean];
  'created': [permission: any];
}>();

const permissionService = usePermissionService();

const name = ref('');
const description = ref('');
const isCreating = ref(false);

const isOpen = computed({
  get: () => props.open,
  set: value => emit('update:open', value),
});

// Reset form when dialog closes
watch(isOpen, (newValue) => {
  if (!newValue) {
    name.value = '';
    description.value = '';
  }
});

async function handleSubmit() {
  if (!name.value)
    return;

  isCreating.value = true;

  try {
    const created = await permissionService.createPermission({
      name: name.value,
      description: description.value || undefined,
      effect: 'allow',
    });

    if (created) {
      toast.success('Permission created', {
        description: `Permission "${name.value}" has been created successfully.`,
      });

      emit('created', created);
      isOpen.value = false;
    }
  }
  catch (error: any) {
    toast.error('Failed to create permission', {
      description: error.message || 'An unexpected error occurred',
    });
  }
  finally {
    isCreating.value = false;
  }
}

function handleCancel() {
  isOpen.value = false;
}
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogContent>
      <DialogHeader>
        <DialogTitle>Create New Permission</DialogTitle>
        <DialogDescription>
          Create a permission to assign immediately. Name is required
          (e.g., "project:read"), description is optional.
        </DialogDescription>
      </DialogHeader>

      <form class="space-y-4" @submit.prevent="handleSubmit">
        <div class="space-y-2">
          <Label for="permission-name">Permission Name *</Label>
          <Input
            id="permission-name"
            v-model="name"
            placeholder="project:read"
            :disabled="isCreating"
          />
          <p class="text-xs text-muted-foreground">
            Format: alphanumeric, underscores, and colons (e.g., "project:read:metadata")
          </p>
        </div>

        <div class="space-y-2">
          <Label for="permission-description">Description (Optional)</Label>
          <Textarea
            id="permission-description"
            v-model="description"
            placeholder="Access read project"
            :disabled="isCreating"
          />
        </div>

        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            :disabled="isCreating"
            @click="handleCancel"
          >
            Cancel
          </Button>
          <Button type="submit" :disabled="!name || isCreating">
            <Loader2 v-if="isCreating" class="mr-2 h-4 w-4 animate-spin" />
            Create & Assign
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
