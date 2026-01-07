import { createValidator } from '@bufbuild/protovalidate';
import { createClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { ApiKeyService } from '~~/gen/altalune/v1/api_key_pb';
import { EmployeeService } from '~~/gen/altalune/v1/employee_pb';
import { IAMMapperService } from '~~/gen/altalune/v1/iam_mapper_pb';
import { OAuthProviderService } from '~~/gen/altalune/v1/oauth_provider_pb';
import { PermissionService } from '~~/gen/altalune/v1/permission_pb';
import { ProjectService } from '~~/gen/altalune/v1/project_pb';
import { RoleService } from '~~/gen/altalune/v1/role_pb';
import { UserService } from '~~/gen/altalune/v1/user_pb';
import { GreeterService } from '~~/gen/greeter/v1/greeter_pb';

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig();

  const transport = createConnectTransport({
    baseUrl: config.public.apiUrl,
  });

  const validator = createValidator();
  const apiKeyClient = createClient(ApiKeyService, transport);
  const greeterClient = createClient(GreeterService, transport);
  const employeeClient = createClient(EmployeeService, transport);
  const projectClient = createClient(ProjectService, transport);
  const userClient = createClient(UserService, transport);
  const roleClient = createClient(RoleService, transport);
  const permissionClient = createClient(PermissionService, transport);
  const iamMapperClient = createClient(IAMMapperService, transport);
  const oauthProviderClient = createClient(OAuthProviderService, transport);

  return {
    provide: {
      validator,
      apiKeyClient,
      greeterClient,
      employeeClient,
      projectClient,
      userClient,
      roleClient,
      permissionClient,
      iamMapperClient,
      oauthProviderClient,
    },
  };
});
