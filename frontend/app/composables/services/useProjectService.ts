import type { MessageInitShape } from '@bufbuild/protobuf';
import type { QueryMetaResponseSchema } from '~~/gen/altalune/v1/common_pb';

import type { Project } from '~~/gen/altalune/v1/project_pb';
import { projectRepository } from '#shared/repository/project';

import { create } from '@bufbuild/protobuf';
import {
  CreateProjectRequestSchema,
  DeleteProjectRequestSchema,
  GetProjectRequestSchema,
  QueryProjectsRequestSchema,
  UpdateProjectRequestSchema,
} from '~~/gen/altalune/v1/project_pb';
import { useConnectValidator } from '../useConnectValidator';
import { useErrorMessage } from '../useErrorMessage';

export function useProjectService() {
  const { $projectClient } = useNuxtApp();
  const project = projectRepository($projectClient);
  const { parseError } = useErrorMessage();

  const queryValidator = useConnectValidator(QueryProjectsRequestSchema);
  const createValidator = useConnectValidator(CreateProjectRequestSchema);
  const getValidator = useConnectValidator(GetProjectRequestSchema);
  const updateValidator = useConnectValidator(UpdateProjectRequestSchema);
  const deleteValidator = useConnectValidator(DeleteProjectRequestSchema);

  // State management
  const createState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  const getState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  const updateState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  const deleteState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  async function query(
    req: MessageInitShape<typeof QueryProjectsRequestSchema>,
  ): Promise<{
    data: Project[];
    meta: MessageInitShape<typeof QueryMetaResponseSchema> | undefined;
  }> {
    queryValidator.reset();

    if (!queryValidator.validate(req)) {
      console.warn('Validation failed for QueryProjectsRequest:', queryValidator.errors.value);
      return {
        data: [],
        meta: {
          rowCount: 0,
          pageCount: 0,
          filters: {},
        },
      };
    }

    try {
      const message = create(QueryProjectsRequestSchema, req);
      const result = await project.queryProjects(message);
      return {
        data: result.data,
        meta: result.meta,
      };
    }
    catch (err) {
      const errorMessage = parseError(err);
      throw new Error(errorMessage);
    }
  }

  async function createProject(
    req: MessageInitShape<typeof CreateProjectRequestSchema>,
  ): Promise<Project | null> {
    createState.loading = true;
    createState.error = '';
    createState.success = false;

    createValidator.reset();

    if (!createValidator.validate(req)) {
      createState.loading = false;
      return null;
    }

    try {
      const message = create(CreateProjectRequestSchema, req);
      const result = await project.createProject(message);
      createState.success = true;
      return result.project || null;
    }
    catch (err) {
      createState.error = parseError(err);
      throw new Error(createState.error);
    }
    finally {
      createState.loading = false;
    }
  }

  function resetCreateState() {
    createState.loading = false;
    createState.error = '';
    createState.success = false;
    createValidator.reset();
  }

  async function getProject(
    req: MessageInitShape<typeof GetProjectRequestSchema>,
  ): Promise<Project | null> {
    getState.loading = true;
    getState.error = '';
    getState.success = false;

    getValidator.reset();

    if (!getValidator.validate(req)) {
      getState.loading = false;
      return null;
    }

    try {
      const message = create(GetProjectRequestSchema, req);
      const result = await project.getProject(message);
      getState.success = true;
      return result.project || null;
    }
    catch (err) {
      getState.error = parseError(err);
      throw new Error(getState.error);
    }
    finally {
      getState.loading = false;
    }
  }

  function resetGetState() {
    getState.loading = false;
    getState.error = '';
    getState.success = false;
    getValidator.reset();
  }

  async function updateProject(
    req: MessageInitShape<typeof UpdateProjectRequestSchema>,
  ): Promise<Project | null> {
    updateState.loading = true;
    updateState.error = '';
    updateState.success = false;

    updateValidator.reset();

    if (!updateValidator.validate(req)) {
      updateState.loading = false;
      return null;
    }

    try {
      const message = create(UpdateProjectRequestSchema, req);
      const result = await project.updateProject(message);
      updateState.success = true;
      return result.project || null;
    }
    catch (err) {
      updateState.error = parseError(err);
      throw new Error(updateState.error);
    }
    finally {
      updateState.loading = false;
    }
  }

  function resetUpdateState() {
    updateState.loading = false;
    updateState.error = '';
    updateState.success = false;
    updateValidator.reset();
  }

  async function deleteProject(
    req: MessageInitShape<typeof DeleteProjectRequestSchema>,
  ): Promise<boolean> {
    deleteState.loading = true;
    deleteState.error = '';
    deleteState.success = false;

    deleteValidator.reset();

    if (!deleteValidator.validate(req)) {
      deleteState.loading = false;
      return false;
    }

    try {
      const message = create(DeleteProjectRequestSchema, req);
      await project.deleteProject(message);
      deleteState.success = true;
      return true;
    }
    catch (err) {
      deleteState.error = parseError(err);
      throw new Error(deleteState.error);
    }
    finally {
      deleteState.loading = false;
    }
  }

  function resetDeleteState() {
    deleteState.loading = false;
    deleteState.error = '';
    deleteState.success = false;
    deleteValidator.reset();
  }

  return {
    // Query
    query,
    queryValidationErrors: queryValidator.errors,

    // Create
    createProject,
    createLoading: computed(() => createState.loading),
    createError: computed(() => createState.error),
    createSuccess: computed(() => createState.success),
    createValidationErrors: createValidator.errors,
    resetCreateState,

    // Get
    getProject,
    getLoading: computed(() => getState.loading),
    getError: computed(() => getState.error),
    getSuccess: computed(() => getState.success),
    getValidationErrors: getValidator.errors,
    resetGetState,

    // Update
    updateProject,
    updateLoading: computed(() => updateState.loading),
    updateError: computed(() => updateState.error),
    updateSuccess: computed(() => updateState.success),
    updateValidationErrors: updateValidator.errors,
    resetUpdateState,

    // Delete
    deleteProject,
    deleteLoading: computed(() => deleteState.loading),
    deleteError: computed(() => deleteState.error),
    deleteSuccess: computed(() => deleteState.success),
    deleteValidationErrors: deleteValidator.errors,
    resetDeleteState,
  };
}
