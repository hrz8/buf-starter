import type { MessageInitShape } from '@bufbuild/protobuf';
import type { ProjectMemberWithUser, UserProjectMembership } from '~~/gen/altalune/v1/iam_mapper_pb';
import type { Permission } from '~~/gen/altalune/v1/permission_pb';
import type { Role } from '~~/gen/altalune/v1/role_pb';
import { iamMapperRepository } from '#shared/repository/iam_mapper';
import { create } from '@bufbuild/protobuf';
import {
  AssignProjectMembersRequestSchema,
  AssignRolePermissionsRequestSchema,
  AssignUserPermissionsRequestSchema,
  AssignUserRolesRequestSchema,
  GetProjectMembersRequestSchema,
  GetRolePermissionsRequestSchema,
  GetUserPermissionsRequestSchema,
  GetUserProjectsRequestSchema,
  GetUserRolesRequestSchema,
  RemoveProjectMembersRequestSchema,
  RemoveRolePermissionsRequestSchema,
  RemoveUserPermissionsRequestSchema,
  RemoveUserRolesRequestSchema,
} from '~~/gen/altalune/v1/iam_mapper_pb';
import { useConnectValidator } from '../useConnectValidator';
import { useErrorMessage } from '../useErrorMessage';

export function useIAMMapperService() {
  const { $iamMapperClient } = useNuxtApp();
  const iamMapper = iamMapperRepository($iamMapperClient);
  const { parseError } = useErrorMessage();

  const assignUserRolesValidator = useConnectValidator(AssignUserRolesRequestSchema);
  const removeUserRolesValidator = useConnectValidator(RemoveUserRolesRequestSchema);
  const getUserRolesValidator = useConnectValidator(GetUserRolesRequestSchema);
  const assignRolePermissionsValidator = useConnectValidator(AssignRolePermissionsRequestSchema);
  const removeRolePermissionsValidator = useConnectValidator(RemoveRolePermissionsRequestSchema);
  const getRolePermissionsValidator = useConnectValidator(GetRolePermissionsRequestSchema);
  const assignUserPermissionsValidator = useConnectValidator(AssignUserPermissionsRequestSchema);
  const removeUserPermissionsValidator = useConnectValidator(RemoveUserPermissionsRequestSchema);
  const getUserPermissionsValidator = useConnectValidator(GetUserPermissionsRequestSchema);
  const assignProjectMembersValidator = useConnectValidator(AssignProjectMembersRequestSchema);
  const removeProjectMembersValidator = useConnectValidator(RemoveProjectMembersRequestSchema);
  const getProjectMembersValidator = useConnectValidator(GetProjectMembersRequestSchema);
  const getUserProjectsValidator = useConnectValidator(GetUserProjectsRequestSchema);

  const mappingState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  // User-Role Mappings
  async function assignUserRoles(
    req: MessageInitShape<typeof AssignUserRolesRequestSchema>,
  ): Promise<boolean> {
    mappingState.loading = true;
    mappingState.error = '';
    mappingState.success = false;

    assignUserRolesValidator.reset();

    if (!assignUserRolesValidator.validate(req)) {
      mappingState.loading = false;
      return false;
    }

    try {
      const message = create(AssignUserRolesRequestSchema, req);
      await iamMapper.assignUserRoles(message);
      mappingState.success = true;
      return true;
    }
    catch (err) {
      mappingState.error = parseError(err);
      throw new Error(mappingState.error);
    }
    finally {
      mappingState.loading = false;
    }
  }

  async function removeUserRoles(
    req: MessageInitShape<typeof RemoveUserRolesRequestSchema>,
  ): Promise<boolean> {
    mappingState.loading = true;
    mappingState.error = '';
    mappingState.success = false;

    removeUserRolesValidator.reset();

    if (!removeUserRolesValidator.validate(req)) {
      mappingState.loading = false;
      return false;
    }

    try {
      const message = create(RemoveUserRolesRequestSchema, req);
      await iamMapper.removeUserRoles(message);
      mappingState.success = true;
      return true;
    }
    catch (err) {
      mappingState.error = parseError(err);
      throw new Error(mappingState.error);
    }
    finally {
      mappingState.loading = false;
    }
  }

  async function getUserRoles(
    req: MessageInitShape<typeof GetUserRolesRequestSchema>,
  ): Promise<Role[]> {
    getUserRolesValidator.reset();

    if (!getUserRolesValidator.validate(req)) {
      return [];
    }

    try {
      const message = create(GetUserRolesRequestSchema, req);
      const result = await iamMapper.getUserRoles(message);
      return result.roles;
    }
    catch (err) {
      const errorMessage = parseError(err);
      throw new Error(errorMessage);
    }
  }

  // Role-Permission Mappings
  async function assignRolePermissions(
    req: MessageInitShape<typeof AssignRolePermissionsRequestSchema>,
  ): Promise<boolean> {
    mappingState.loading = true;
    mappingState.error = '';
    mappingState.success = false;

    assignRolePermissionsValidator.reset();

    if (!assignRolePermissionsValidator.validate(req)) {
      mappingState.loading = false;
      return false;
    }

    try {
      const message = create(AssignRolePermissionsRequestSchema, req);
      await iamMapper.assignRolePermissions(message);
      mappingState.success = true;
      return true;
    }
    catch (err) {
      mappingState.error = parseError(err);
      throw new Error(mappingState.error);
    }
    finally {
      mappingState.loading = false;
    }
  }

  async function removeRolePermissions(
    req: MessageInitShape<typeof RemoveRolePermissionsRequestSchema>,
  ): Promise<boolean> {
    mappingState.loading = true;
    mappingState.error = '';
    mappingState.success = false;

    removeRolePermissionsValidator.reset();

    if (!removeRolePermissionsValidator.validate(req)) {
      mappingState.loading = false;
      return false;
    }

    try {
      const message = create(RemoveRolePermissionsRequestSchema, req);
      await iamMapper.removeRolePermissions(message);
      mappingState.success = true;
      return true;
    }
    catch (err) {
      mappingState.error = parseError(err);
      throw new Error(mappingState.error);
    }
    finally {
      mappingState.loading = false;
    }
  }

  async function getRolePermissions(
    req: MessageInitShape<typeof GetRolePermissionsRequestSchema>,
  ): Promise<Permission[]> {
    getRolePermissionsValidator.reset();

    if (!getRolePermissionsValidator.validate(req)) {
      return [];
    }

    try {
      const message = create(GetRolePermissionsRequestSchema, req);
      const result = await iamMapper.getRolePermissions(message);
      return result.permissions;
    }
    catch (err) {
      const errorMessage = parseError(err);
      throw new Error(errorMessage);
    }
  }

  // User-Permission Mappings
  async function assignUserPermissions(
    req: MessageInitShape<typeof AssignUserPermissionsRequestSchema>,
  ): Promise<boolean> {
    mappingState.loading = true;
    mappingState.error = '';
    mappingState.success = false;

    assignUserPermissionsValidator.reset();

    if (!assignUserPermissionsValidator.validate(req)) {
      mappingState.loading = false;
      return false;
    }

    try {
      const message = create(AssignUserPermissionsRequestSchema, req);
      await iamMapper.assignUserPermissions(message);
      mappingState.success = true;
      return true;
    }
    catch (err) {
      mappingState.error = parseError(err);
      throw new Error(mappingState.error);
    }
    finally {
      mappingState.loading = false;
    }
  }

  async function removeUserPermissions(
    req: MessageInitShape<typeof RemoveUserPermissionsRequestSchema>,
  ): Promise<boolean> {
    mappingState.loading = true;
    mappingState.error = '';
    mappingState.success = false;

    removeUserPermissionsValidator.reset();

    if (!removeUserPermissionsValidator.validate(req)) {
      mappingState.loading = false;
      return false;
    }

    try {
      const message = create(RemoveUserPermissionsRequestSchema, req);
      await iamMapper.removeUserPermissions(message);
      mappingState.success = true;
      return true;
    }
    catch (err) {
      mappingState.error = parseError(err);
      throw new Error(mappingState.error);
    }
    finally {
      mappingState.loading = false;
    }
  }

  async function getUserPermissions(
    req: MessageInitShape<typeof GetUserPermissionsRequestSchema>,
  ): Promise<Permission[]> {
    getUserPermissionsValidator.reset();

    if (!getUserPermissionsValidator.validate(req)) {
      return [];
    }

    try {
      const message = create(GetUserPermissionsRequestSchema, req);
      const result = await iamMapper.getUserPermissions(message);
      return result.permissions;
    }
    catch (err) {
      const errorMessage = parseError(err);
      throw new Error(errorMessage);
    }
  }

  // Project Members
  async function assignProjectMembers(
    req: MessageInitShape<typeof AssignProjectMembersRequestSchema>,
  ): Promise<boolean> {
    mappingState.loading = true;
    mappingState.error = '';
    mappingState.success = false;

    assignProjectMembersValidator.reset();

    if (!assignProjectMembersValidator.validate(req)) {
      mappingState.loading = false;
      return false;
    }

    try {
      const message = create(AssignProjectMembersRequestSchema, req);
      await iamMapper.assignProjectMembers(message);
      mappingState.success = true;
      return true;
    }
    catch (err) {
      mappingState.error = parseError(err);
      throw new Error(mappingState.error);
    }
    finally {
      mappingState.loading = false;
    }
  }

  async function removeProjectMembers(
    req: MessageInitShape<typeof RemoveProjectMembersRequestSchema>,
  ): Promise<boolean> {
    mappingState.loading = true;
    mappingState.error = '';
    mappingState.success = false;

    removeProjectMembersValidator.reset();

    if (!removeProjectMembersValidator.validate(req)) {
      mappingState.loading = false;
      return false;
    }

    try {
      const message = create(RemoveProjectMembersRequestSchema, req);
      await iamMapper.removeProjectMembers(message);
      mappingState.success = true;
      return true;
    }
    catch (err) {
      mappingState.error = parseError(err);
      throw new Error(mappingState.error);
    }
    finally {
      mappingState.loading = false;
    }
  }

  async function getProjectMembers(
    req: MessageInitShape<typeof GetProjectMembersRequestSchema>,
  ): Promise<ProjectMemberWithUser[]> {
    getProjectMembersValidator.reset();

    if (!getProjectMembersValidator.validate(req)) {
      return [];
    }

    try {
      const message = create(GetProjectMembersRequestSchema, req);
      const result = await iamMapper.getProjectMembers(message);
      return result.members;
    }
    catch (err) {
      const errorMessage = parseError(err);
      throw new Error(errorMessage);
    }
  }

  // User Projects (reverse lookup - projects a user belongs to)
  async function getUserProjects(
    req: MessageInitShape<typeof GetUserProjectsRequestSchema>,
  ): Promise<UserProjectMembership[]> {
    getUserProjectsValidator.reset();

    if (!getUserProjectsValidator.validate(req)) {
      return [];
    }

    try {
      const message = create(GetUserProjectsRequestSchema, req);
      const result = await iamMapper.getUserProjects(message);
      return result.projects;
    }
    catch (err) {
      const errorMessage = parseError(err);
      throw new Error(errorMessage);
    }
  }

  return {
    // User-Role Mappings
    assignUserRoles,
    removeUserRoles,
    getUserRoles,

    // Role-Permission Mappings
    assignRolePermissions,
    removeRolePermissions,
    getRolePermissions,

    // User-Permission Mappings
    assignUserPermissions,
    removeUserPermissions,
    getUserPermissions,

    // Project Members
    assignProjectMembers,
    removeProjectMembers,
    getProjectMembers,

    // User Projects
    getUserProjects,

    // State
    mappingLoading: computed(() => mappingState.loading),
    mappingError: computed(() => mappingState.error),
    mappingSuccess: computed(() => mappingState.success),
  };
}
