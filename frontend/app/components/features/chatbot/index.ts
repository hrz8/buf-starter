// Re-export field components
export * from './fields';
// Barrel exports for chatbot components
export { default as ModuleConfigForm } from './ModuleConfigForm.vue';
export { default as ModuleToggle } from './ModuleToggle.vue';
export { default as SchemaField } from './SchemaField.vue';

export { default as SchemaForm } from './SchemaForm.vue';

// Re-export schemas from lib (convenience re-export)
export * from '@/lib/chatbot-modules';
