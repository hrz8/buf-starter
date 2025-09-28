<script setup lang="ts">
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import { Separator } from '@/components/ui/separator';

interface Option {
  label: string;
  value: string;
  icon?: Component;
}

interface Props {
  title: string;
  options: Option[];
  modelValue: string[];
  multiple?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  multiple: true,
});

const emit = defineEmits<{
  'update:modelValue': [value: string[]];
  'update': [value: string[]];
  'clear': [];
}>();

const selectedValues = ref<string[]>(props.modelValue ?? []);
const open = ref(false);

watch(() => props.modelValue, (newValue) => {
  selectedValues.value = newValue ?? [];
}, { deep: true });

function toggleValue(value: string) {
  if (props.multiple) {
    if (selectedValues.value.includes(value)) {
      selectedValues.value = selectedValues.value.filter(v => v !== value);
    }
    else {
      selectedValues.value = [...selectedValues.value, value];
    }
  }
  else {
    selectedValues.value = selectedValues.value[0] === value ? [] : [value];
  }

  emit('update:modelValue', [...selectedValues.value]);
  emit('update', [...selectedValues.value]);
}

function clearAll() {
  selectedValues.value = [];
  emit('update:modelValue', []);
  emit('update', []);
  emit('clear');
  open.value = false;
}

const selectedLabels = computed(() => {
  return selectedValues.value
    .map(value => props.options.find(option => option.value === value)?.label)
    .filter(Boolean);
});
</script>

<template>
  <Popover v-model:open="open">
    <PopoverTrigger as-child>
      <Button
        variant="outline"
        size="sm"
        class="h-8 border-dashed"
      >
        <Icon
          name="radix-icons:plus-circled"
          class="mr-2 h-4 w-4"
        />
        {{ title }}
        <template v-if="selectedValues.length > 0">
          <Separator
            orientation="vertical"
            class="mx-2 h-4"
          />
          <Badge
            variant="secondary"
            class="rounded-sm px-1 font-normal lg:hidden"
          >
            {{ selectedValues.length }}
          </Badge>
          <div class="hidden space-x-1 lg:flex">
            <Badge
              v-for="label in (selectedValues.length <= 2 ? selectedLabels : [])"
              :key="label"
              variant="secondary"
              class="rounded-sm px-1 font-normal"
            >
              {{ label }}
            </Badge>
            <Badge
              v-if="selectedValues.length > 2"
              variant="secondary"
              class="rounded-sm px-1 font-normal"
            >
              {{ selectedValues.length }} selected
            </Badge>
          </div>
        </template>
      </Button>
    </PopoverTrigger>
    <PopoverContent
      class="w-[200px] p-0"
      align="start"
    >
      <Command>
        <CommandInput :placeholder="`Search ${title?.toLowerCase()}...`" />
        <CommandList>
          <CommandEmpty>No results found.</CommandEmpty>
          <CommandGroup>
            <CommandItem
              v-for="option in options"
              :key="option.value"
              :value="option.label"
              @select="toggleValue(option.value)"
            >
              <div
                class="
                  mr-2 flex h-4 w-4 items-center
                  justify-center rounded-sm border border-primary
                "
                :class="selectedValues.includes(option.value)
                  ? 'bg-primary text-primary-foreground'
                  : 'opacity-50 [&_svg]:invisible'"
              >
                <Icon
                  name="radix-icons:check"
                  class="h-4 w-4"
                />
              </div>
              <component
                :is="option.icon"
                v-if="option.icon"
                class="mr-2 h-4 w-4 text-muted-foreground"
              />
              <span>{{ option.label }}</span>
            </CommandItem>
          </CommandGroup>
          <template v-if="selectedValues.length > 0">
            <CommandSeparator />
            <CommandGroup>
              <CommandItem
                value="clear"
                class="justify-center text-center"
                @select="clearAll"
              >
                Clear filters
              </CommandItem>
            </CommandGroup>
          </template>
        </CommandList>
      </Command>
    </PopoverContent>
  </Popover>
</template>
