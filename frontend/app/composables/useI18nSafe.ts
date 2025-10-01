import type { VNode } from 'vue';

/**
 * Safe i18n composable with inline formatting support
 *
 * Supports markdown-style formatting without v-html:
 * - **bold text** -> <strong>
 * - *italic text* -> <em>
 * - __underline text__ -> <u>
 *
 * @example
 * ```vue
 * <script setup>
 * const { t, tFormatted } = useI18nSafe();
 * </script>
 *
 * <template>
 *   <!-- Plain text -->
 *   <p>{{ t('key') }}</p>
 *
 *   <!-- With formatting (renders VNodes, no v-html) -->
 *   <p><component :is="tFormatted('key')" /></p>
 * </template>
 * ```
 */
export function useI18nSafe() {
  const { t } = useI18n();

  /**
   * Render translation with SAFE inline formatting
   * Converts markdown-style syntax to Vue VNodes
   *
   * @param key - Translation key
   * @param args - Interpolation arguments
   * @returns VNode with formatted content
   */
  function tFormatted(key: string, args?: Record<string, unknown>): VNode {
    const text = args ? t(key, args) : t(key);

    // Parse markdown-style formatting
    // Regex captures: **bold**, *italic*, __underline__
    const parts = text.split(/(\*\*[^*]+\*\*|\*[^*]+\*|__[^_]+__)/g);

    const children = parts.map((part) => {
      // Bold: **text**
      if (part.startsWith('**') && part.endsWith('**')) {
        return h('strong', part.slice(2, -2));
      }
      // Italic: *text*
      if (part.startsWith('*') && part.endsWith('*') && !part.includes('**')) {
        return h('em', part.slice(1, -1));
      }
      // Underline: __text__
      if (part.startsWith('__') && part.endsWith('__')) {
        return h('u', part.slice(2, -2));
      }
      // Plain text
      return part;
    });

    return h('span', children);
  }

  return { t, tFormatted };
}
