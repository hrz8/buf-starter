import type { MessageInitShape } from '@bufbuild/protobuf';
import type { QueryMetaResponseSchema } from '~~/gen/altalune/v1/common_pb';

import type { Project } from '~~/gen/altalune/v1/project_pb';
import { projectRepository } from '#shared/repository/project';

import { create } from '@bufbuild/protobuf';
import {
  CreateProjectRequestSchema,

  QueryProjectsRequestSchema,
} from '~~/gen/altalune/v1/project_pb';
import { useConnectValidator } from '../useConnectValidator';
import { useErrorMessage } from '../useErrorMessage';

export function useProjectService() {
  const { $projectClient } = useNuxtApp();
  const project = projectRepository($projectClient);
  const { parseError } = useErrorMessage();

  const queryValidator = useConnectValidator(QueryProjectsRequestSchema);
  const createValidator = useConnectValidator(CreateProjectRequestSchema);

  // Create state for form submission
  const createState = reactive({
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
  };
}
