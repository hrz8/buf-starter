/**
 * Permission constants for type-safe permission checks
 * Matches backend predefined permissions
 */
export const PERMISSIONS = {
  // Employee entity
  EMPLOYEE: {
    READ: 'employee:read',
    WRITE: 'employee:write',
    DELETE: 'employee:delete',
  },
  // User entity (IAM)
  USER: {
    READ: 'user:read',
    WRITE: 'user:write',
    DELETE: 'user:delete',
  },
  // Role entity (IAM)
  ROLE: {
    READ: 'role:read',
    WRITE: 'role:write',
    DELETE: 'role:delete',
  },
  // Permission entity (IAM)
  PERMISSION: {
    READ: 'permission:read',
    WRITE: 'permission:write',
    DELETE: 'permission:delete',
  },
  // Project entity
  PROJECT: {
    READ: 'project:read',
    WRITE: 'project:write',
    DELETE: 'project:delete',
  },
  // API Key entity
  API_KEY: {
    READ: 'apikey:read',
    WRITE: 'apikey:write',
    DELETE: 'apikey:delete',
  },
  // Chatbot Config entity
  CHATBOT: {
    READ: 'chatbot:read',
    WRITE: 'chatbot:write',
    DELETE: 'chatbot:delete',
  },
  // OAuth Client entity
  CLIENT: {
    READ: 'client:read',
    WRITE: 'client:write',
    DELETE: 'client:delete',
  },
  // OAuth Provider entity
  PROVIDER: {
    READ: 'provider:read',
    WRITE: 'provider:write',
    DELETE: 'provider:delete',
  },
  // Project Members management
  MEMBER: {
    READ: 'member:read',
    WRITE: 'member:write',
  },
  // IAM Mapper (role-permission assignments)
  IAM: {
    READ: 'iam:read',
    WRITE: 'iam:write',
  },
  // Special permissions
  ROOT: 'root',
} as const;

/**
 * Type for all permission values
 */
export type Permission
  = typeof PERMISSIONS.EMPLOYEE[keyof typeof PERMISSIONS.EMPLOYEE]
    | typeof PERMISSIONS.USER[keyof typeof PERMISSIONS.USER]
    | typeof PERMISSIONS.ROLE[keyof typeof PERMISSIONS.ROLE]
    | typeof PERMISSIONS.PERMISSION[keyof typeof PERMISSIONS.PERMISSION]
    | typeof PERMISSIONS.PROJECT[keyof typeof PERMISSIONS.PROJECT]
    | typeof PERMISSIONS.API_KEY[keyof typeof PERMISSIONS.API_KEY]
    | typeof PERMISSIONS.CHATBOT[keyof typeof PERMISSIONS.CHATBOT]
    | typeof PERMISSIONS.CLIENT[keyof typeof PERMISSIONS.CLIENT]
    | typeof PERMISSIONS.PROVIDER[keyof typeof PERMISSIONS.PROVIDER]
    | typeof PERMISSIONS.MEMBER[keyof typeof PERMISSIONS.MEMBER]
    | typeof PERMISSIONS.IAM[keyof typeof PERMISSIONS.IAM]
    | typeof PERMISSIONS.ROOT;

/**
 * Project membership roles
 */
export const PROJECT_ROLES = {
  OWNER: 'owner',
  ADMIN: 'admin',
  MEMBER: 'member',
  USER: 'user',
} as const;

export type ProjectRole = typeof PROJECT_ROLES[keyof typeof PROJECT_ROLES];

/**
 * Permission descriptions for UI display
 */
export const PERMISSION_DESCRIPTIONS: Record<string, string> = {
  'employee:read': 'View employee records',
  'employee:write': 'Create and update employee records',
  'employee:delete': 'Delete employee records',
  'user:read': 'View user accounts',
  'user:write': 'Create and update user accounts',
  'user:delete': 'Delete user accounts',
  'role:read': 'View roles',
  'role:write': 'Create and update roles',
  'role:delete': 'Delete roles',
  'permission:read': 'View permissions',
  'permission:write': 'Create and update permissions',
  'permission:delete': 'Delete permissions',
  'project:read': 'View projects',
  'project:write': 'Create and update projects',
  'project:delete': 'Delete projects',
  'apikey:read': 'View API keys',
  'apikey:write': 'Create and update API keys',
  'apikey:delete': 'Delete API keys',
  'chatbot:read': 'View chatbot configurations',
  'chatbot:write': 'Create and update chatbot configurations',
  'chatbot:delete': 'Delete chatbot configurations',
  'client:read': 'View OAuth clients',
  'client:write': 'Create and update OAuth clients',
  'client:delete': 'Delete OAuth clients',
  'provider:read': 'View OAuth providers',
  'provider:write': 'Create and update OAuth providers',
  'provider:delete': 'Delete OAuth providers',
  'member:read': 'View project members',
  'member:write': 'Add and update project members',
  'iam:read': 'View IAM mappings',
  'iam:write': 'Manage IAM mappings',
  'root': 'Full superadmin access',
};
