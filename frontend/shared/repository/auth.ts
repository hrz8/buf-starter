import type { $Fetch } from 'nitropack';

export interface AuthExchangeRequest {
  code: string;
  code_verifier: string;
  redirect_uri: string;
}

export interface AuthUserInfo {
  sub: string;
  email?: string;
  name?: string;
}

export interface AuthExchangeResponse {
  user: AuthUserInfo;
  expires_in: number;
}

export interface AuthErrorResponse {
  error: string;
  error_description?: string;
}

export function authRepository(f: $Fetch) {
  return {
    async exchange(req: AuthExchangeRequest): Promise<AuthExchangeResponse> {
      try {
        return await f<AuthExchangeResponse>('/exchange', {
          method: 'POST',
          body: req,
          credentials: 'include',
        });
      }
      catch (error) {
        console.error('Auth exchange error:', error);
        throw error;
      }
    },

    async refresh(): Promise<AuthExchangeResponse> {
      try {
        return await f<AuthExchangeResponse>('/refresh', {
          method: 'POST',
          credentials: 'include',
        });
      }
      catch (error) {
        console.error('Auth refresh error:', error);
        throw error;
      }
    },

    async me(): Promise<AuthExchangeResponse> {
      try {
        return await f<AuthExchangeResponse>('/me', {
          method: 'GET',
          credentials: 'include',
        });
      }
      catch (error) {
        console.error('Auth me error:', error);
        throw error;
      }
    },

    async logout(): Promise<void> {
      try {
        await f('/logout', {
          method: 'POST',
          credentials: 'include',
        });
      }
      catch (error) {
        console.error('Auth logout error:', error);
        throw error;
      }
    },
  };
}
