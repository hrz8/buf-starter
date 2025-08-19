import { QueryProjectsRequestSchema, type Project } from '~~/gen/altalune/v1/project_pb';
import { type MessageInitShape, create } from '@bufbuild/protobuf';

import type { QueryMetaResponseSchema } from '~~/gen/altalune/v1/common_pb';

import { projectRepository } from '#shared/repository/project';
import { useConnectValidator } from '../useConnectValidator';
import { useErrorMessage } from '../useErrorMessage';

export function useProjectService() {
  const { $projectClient } = useNuxtApp();
  const project = projectRepository($projectClient);

  const { parseError } = useErrorMessage();

  const queryValidator = useConnectValidator(QueryProjectsRequestSchema);

  async function query(
    req: MessageInitShape<typeof QueryProjectsRequestSchema>,
  ): Promise<{
    data: Project[];
    meta: MessageInitShape<typeof QueryMetaResponseSchema> | undefined;
  }> {
    queryValidator.reset();

    if (!queryValidator.validate(req)) {
      console.warn('Validation failed for QueryEmployeesRequest:', queryValidator.errors.value);
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
    } catch (err) {
      const errorMessage = parseError(err);
      throw new Error(errorMessage);
    }
  }

  return {
    // Query
    query,

    // Validation
    queryValidationErrors: queryValidator.errors,
  };
}
