<script setup lang="ts">
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb';
import { Separator } from '@/components/ui/separator';
import { SidebarTrigger } from '@/components/ui/sidebar';
import { useBreadcrumbs } from '~/composables/navigation/useBreadcrumbs';

const { locales, setLocale, locale } = useI18n();
const colorMode = useColorMode();

const flags = {
  'en-US': 'emojione:flag-for-united-states',
  'id-ID': 'emojione:flag-for-indonesia',
};

const currentLocale = computed(() => {
  const found = locales.value.find(l => l.code === locale.value);
  return found ? found.name : locale.value;
});

// Breadcrumb navigation
const { breadcrumbs, hasBreadcrumbs } = useBreadcrumbs();
</script>

<template>
  <header
    class="
      flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear
      group-has-data-[collapsible=icon]/sidebar-wrapper:h-12
    "
  >
    <div class="flex items-center gap-2 px-4 flex-1">
      <SidebarTrigger class="-ml-1" />
      <Separator
        orientation="vertical"
        class="mr-2 data-[orientation=vertical]:h-4"
      />
      <Breadcrumb v-if="hasBreadcrumbs">
        <BreadcrumbList>
          <template
            v-for="(crumb, index) in breadcrumbs"
            :key="crumb.path"
          >
            <BreadcrumbItem
              :class="{ 'hidden md:block': index === 0 }"
            >
              <BreadcrumbLink
                v-if="!crumb.isCurrent"
                :href="crumb.path"
              >
                {{ crumb.label }}
              </BreadcrumbLink>
              <BreadcrumbPage v-else>
                {{ crumb.label }}
              </BreadcrumbPage>
            </BreadcrumbItem>
            <BreadcrumbSeparator
              v-if="index < breadcrumbs.length - 1"
              :class="{ 'hidden md:block': index === 0 }"
            />
          </template>
        </BreadcrumbList>
      </Breadcrumb>
    </div>

    <div class="pr-4">
      <DropdownMenu>
        <DropdownMenuTrigger class="flex items-center gap-1 px-3 py-2 rounded-md hover:bg-muted">
          <Icon
            name="lucide:languages"
            class="w-4 h-4"
          />
          <span class="text-sm">{{ currentLocale }}</span>
        </DropdownMenuTrigger>

        <DropdownMenuContent class="w-32">
          <DropdownMenuItem
            v-for="l in locales"
            :key="l.code"
            @click="setLocale(l.code)"
          >
            <Icon
              :name="flags[l.code]"
              class="w-4 h-4"
            />
            {{ l.name }}
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
    <div class="pr-4">
      <DropdownMenu>
        <DropdownMenuTrigger class="flex items-center gap-1 px-3 py-2 rounded-md hover:bg-muted">
          <Icon
            name="radix-icons:moon"
            class="
              h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all
              dark:-rotate-90 dark:scale-0
            "
          />
          <Icon
            name="radix-icons:sun"
            class="
              absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0
              transition-all dark:rotate-0 dark:scale-100
            "
          />
          <span class="sr-only">Toggle theme</span>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem
            @click="colorMode = 'light'"
          >
            Light
          </DropdownMenuItem>
          <DropdownMenuItem
            @click="colorMode = 'dark'"
          >
            Dark
          </DropdownMenuItem>
          <DropdownMenuItem
            @click="colorMode = 'auto'"
          >
            System
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  </header>
</template>
