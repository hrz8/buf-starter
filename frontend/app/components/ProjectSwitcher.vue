<script setup lang="ts">
import {
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuContent,
  DropdownMenuTrigger,
  DropdownMenuLabel,
  DropdownMenuItem,
  DropdownMenu,
} from '@/components/ui/dropdown-menu';
import {
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenu,
  useSidebar,
} from '@/components/ui/sidebar';
import { useProjectStore } from '@/stores/project';

const DEFAULT_PROJECT_ICON = 'lucide:folder';
const RANDOM_ICONS = [
  'lucide:aperture',
  'lucide:biohazard',
  'lucide:cat',
  'lucide:club',
  'lucide:component',
  'lucide:crown',
  'lucide:fan',
  'lucide:fingerprint',
  'lucide:paw-print',
  'lucide:gem',
  'lucide:pizza',
  'lucide:ship-wheel',
  'lucide:leaf',
  'lucide:rainbow',
  'lucide:snowflake',
  'lucide:sailboat',
];

const getProjectIcon = (publicId: string) => {
  const iconIndex = publicId.charCodeAt(0) % RANDOM_ICONS.length;
  return RANDOM_ICONS[iconIndex];
};

const { isMobile } = useSidebar();
const projectStore = useProjectStore();

const pending = computed(() => projectStore.pending);
const error = computed(() => projectStore.error);
</script>

<template>
  <SidebarMenu>
    <SidebarMenuItem>
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <SidebarMenuButton
            size="lg"
            class="
              data-[state=open]:bg-sidebar-accent
              data-[state=open]:text-sidebar-accent-foreground
            "
          >
            <div
              class="
                flex aspect-square size-8 items-center
                justify-center rounded-lg bg-sidebar-primary
                text-sidebar-primary-foreground
              "
            >
              <Icon
                :name="projectStore.currentProject
                  ? getProjectIcon(projectStore.currentProject.id) ?? DEFAULT_PROJECT_ICON
                  : DEFAULT_PROJECT_ICON"
                size="1.5em"
              />
            </div>
            <div class="grid flex-1 text-left text-sm leading-tight">
              <template v-if="projectStore.currentProject">
                <span class="truncate font-medium">
                  {{ projectStore.currentProject.name }}
                </span>
                <span class="truncate text-xs">
                  {{ projectStore.currentProject.environment }}
                </span>
              </template>
              <template v-else-if="pending">
                <span class="font-medium">
                  <Skeleton class="h-4 w-24 rounded" />
                </span>
                <span>
                  <Skeleton class="h-3 w-32 rounded mt-1" />
                </span>
              </template>
              <template v-else>
                <span class="truncate font-medium text-muted-foreground">
                  No project selected
                </span>
              </template>
            </div>
            <Icon
              name="lucide:chevrons-up-down"
              size="1.5em"
            />
          </SidebarMenuButton>
        </DropdownMenuTrigger>
        <DropdownMenuContent
          class="w-[--reka-dropdown-menu-trigger-width] min-w-56 rounded-lg"
          align="start"
          :side="isMobile ? 'bottom' : 'right'"
          :side-offset="4"
        >
          <DropdownMenuLabel class="text-xs text-muted-foreground">
            Projects
          </DropdownMenuLabel>

          <div
            v-if="pending && !projectStore.projects.length"
            class="p-2"
          >
            <div class="flex flex-col gap-2">
              <div
                v-for="i in 3"
                :key="i"
                class="flex items-center gap-2"
              >
                <Skeleton class="h-5 w-5 rounded-full" />
                <Skeleton class="h-4 w-32 rounded" />
                <Skeleton class="h-3 w-8 rounded ml-auto" />
              </div>
            </div>
          </div>

          <div
            v-else-if="error"
            class="p-2 text-sm text-destructive"
          >
            {{ error.message || 'Oops, something went wrong' }}
          </div>

          <template v-else-if="projectStore.projects.length > 0">
            <div
              v-if="pending"
              class="px-2 py-1"
            >
              <div class="h-1 bg-muted rounded overflow-hidden">
                <div class="h-full bg-primary animate-pulse w-1/3" />
              </div>
            </div>

            <DropdownMenuItem
              v-for="(project, index) in projectStore.projects"
              :key="project.id"
              class="gap-2 p-2"
              :class="{ 'opacity-60': pending }"
              @click="projectStore.setActiveProject(project.id)"
            >
              <Icon
                :name="getProjectIcon(project.id) ?? DEFAULT_PROJECT_ICON"
                size="1em"
              />
              {{ project.name }}
              <DropdownMenuShortcut>⌘{{ index + 1 }}</DropdownMenuShortcut>
            </DropdownMenuItem>
          </template>

          <div
            v-else
            class="p-2 text-sm text-muted-foreground"
          >
            No projects found
          </div>

          <DropdownMenuSeparator />
          <DropdownMenuItem class="gap-2 p-2">
            <div class="flex size-6 items-center justify-center rounded-md border bg-transparent">
              <Icon
                name="lucide:plus"
                class="size-4"
              />
            </div>
            <div class="font-medium text-muted-foreground">
              Add project
            </div>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </SidebarMenuItem>
  </SidebarMenu>
</template>
