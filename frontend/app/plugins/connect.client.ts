import type { Interceptor } from '@connectrpc/connect';
import { createValidator } from '@bufbuild/protovalidate';
import { Code, ConnectError, createClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { ApiKeyService } from '~~/gen/altalune/v1/api_key_pb';
import { EmployeeService } from '~~/gen/altalune/v1/employee_pb';
import { IAMMapperService } from '~~/gen/altalune/v1/iam_mapper_pb';
import { OAuthClientService } from '~~/gen/altalune/v1/oauth_client_pb';
import { OAuthProviderService } from '~~/gen/altalune/v1/oauth_provider_pb';
import { PermissionService } from '~~/gen/altalune/v1/permission_pb';
import { ProjectService } from '~~/gen/altalune/v1/project_pb';
import { RoleService } from '~~/gen/altalune/v1/role_pb';
import { UserService } from '~~/gen/altalune/v1/user_pb';
import { GreeterService } from '~~/gen/greeter/v1/greeter_pb';
import { useAuthService } from '@/composables/useAuthService';
import { useAuthStore } from '@/stores/auth';

function createAuthInterceptor(): Interceptor {
  return next => async (req) => {
    try {
      return await next(req);
    }
    catch (error) {
      if (error instanceof ConnectError && error.code === Code.Unauthenticated) {
        try {
          const authService = useAuthService();
          const refreshed = await authService.checkAndRefreshIfNeeded();

          if (refreshed) {
            return await next(req);
          }
        }
        catch (refreshError) {
          console.error('[Connect] Token refresh failed:', refreshError);
        }

        const authStore = useAuthStore();
        authStore.clearAuth();
        navigateTo('/auth/login');
      }

      throw error;
    }
  };
}

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig();

  const authInterceptor = createAuthInterceptor();

  const transport = createConnectTransport({
    baseUrl: config.public.apiUrl,
    interceptors: [authInterceptor],
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
  const oauthClientClient = createClient(OAuthClientService, transport);
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
      oauthClientClient,
      oauthProviderClient,
    },
  };
});
