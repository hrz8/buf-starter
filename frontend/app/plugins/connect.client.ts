import { createConnectTransport } from '@connectrpc/connect-web';
import { EmployeeService } from '~~/gen/altalune/v1/employee_pb';
import { ProjectService } from '~~/gen/altalune/v1/project_pb';
import { GreeterService } from '~~/gen/greeter/v1/greeter_pb';
import { createValidator } from '@bufbuild/protovalidate';
import { createClient } from '@connectrpc/connect';

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig();

  const transport = createConnectTransport({
    baseUrl: config.public.apiUrl,
  });

  const validator = createValidator();
  const greeterClient = createClient(GreeterService, transport);
  const employeeClient = createClient(EmployeeService, transport);
  const projectClient = createClient(ProjectService, transport);

  return {
    provide: {
      validator,
      greeterClient,
      employeeClient,
      projectClient,
    },
  };
});
