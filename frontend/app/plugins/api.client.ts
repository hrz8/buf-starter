function mergeHeaders(
  baseHeaders: Record<string, string>,
  optionsHeaders: HeadersInit,
): HeadersInit {
  let headers: HeadersInit = { ...baseHeaders };

  if (optionsHeaders instanceof Headers) {
    optionsHeaders.forEach((value, key) => {
      (headers as Record<string, string>)[key] = value;
    });
  }
  else if (Array.isArray(optionsHeaders)) {
    optionsHeaders.forEach(([key, value]) => {
      (headers as Record<string, string>)[key] = value;
    });
  }
  else {
    headers = { ...headers, ...optionsHeaders };
  }

  return headers;
}

// Usage Example:
// const { $api } = useNuxtApp();
// const apiClient = $api.createClient();
//
// // Create a repository that accepts the API client
// function userRepository(f: $Fetch) {
//   return {
//     async getUsers() {
//       return await f<User[]>('/api/users', {
//         method: 'GET',
//       });
//     },
//     async getUserById(id: string) {
//       return await f<User>(`/api/users/${id}`, {
//         method: 'GET',
//       });
//     },
//     async createUser(data: CreateUserInput) {
//       return await f<User>('/api/users', {
//         method: 'POST',
//         body: data,
//       });
//     },
//   };
// }
//
// const repo = userRepository(apiClient);
// await repo.getUsers();
export default defineNuxtPlugin(() => {
  function createApiClient(
    baseUrl?: string,
    opts?: { locale?: string },
  ) {
    const config = useRuntimeConfig();

    const baseHeaders: Record<string, string> = {
      'User-Agent': import.meta.client && typeof navigator !== 'undefined' ? navigator.userAgent : '',
    };

    let baseURL: string | undefined = baseUrl;
    if (typeof baseUrl === 'undefined' && typeof config.public.apiUrl === 'string') {
      baseURL = config.public.apiUrl;
    }

    if (opts?.locale) {
      baseHeaders['Accept-Language'] = opts.locale;
    }

    async function fetcher<T>(
      url: string,
      options?: Parameters<typeof $fetch>[1],
    ) {
      try {
        return await $fetch<T>(
          url,
          {
            ...options,
            baseURL,
            headers: mergeHeaders(baseHeaders, options?.headers ?? {}),
          },
        );
      }
      catch (error) {
        console.error('API fetch error:', error);
        throw error;
      }
    }

    return fetcher as typeof $fetch;
  }

  return {
    provide: {
      api: {
        createClient: createApiClient,
      },
    },
  };
});
