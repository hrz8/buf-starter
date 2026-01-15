import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { AlertCircle, CheckCircle, Loader2 } from 'lucide-react';
import { useEffect, useState } from 'react';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { altaluneClient, useAuthStore } from '@/lib/auth';
import { parseCallbackParams } from '@/lib/oauth';

export const Route = createFileRoute('/auth/altalune/callback')({
  component: CallbackPage,
});

function CallbackPage() {
  const navigate = useNavigate();
  const setAuth = useAuthStore((state) => state.setAuth);
  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading');
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const handleCallback = async () => {
      try {
        const params = parseCallbackParams(window.location.search);

        if (params.error) {
          throw new Error(params.errorDescription || params.error);
        }

        const { tokens, user, returnTo } = await altaluneClient.handleCallback(params);

        setAuth(tokens, user, 'altalune');
        setStatus('success');

        // Redirect after short delay
        setTimeout(() => {
          navigate({ to: returnTo || '/' });
        }, 1000);
      } catch (err) {
        console.error('OAuth callback error:', err);
        setError(err instanceof Error ? err.message : 'Authentication failed');
        setStatus('error');
      }
    };

    handleCallback();
  }, [navigate, setAuth]);

  return (
    <div className="flex min-h-[70vh] items-center justify-center">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          {status === 'loading' && (
            <>
              <Loader2 className="mx-auto h-12 w-12 animate-spin text-primary" />
              <CardTitle className="mt-4">Authenticating...</CardTitle>
              <CardDescription>
                Please wait while we complete your sign in
              </CardDescription>
            </>
          )}

          {status === 'success' && (
            <>
              <CheckCircle className="mx-auto h-12 w-12 text-green-500" />
              <CardTitle className="mt-4">Success!</CardTitle>
              <CardDescription>
                You have been signed in. Redirecting...
              </CardDescription>
            </>
          )}

          {status === 'error' && (
            <>
              <AlertCircle className="mx-auto h-12 w-12 text-destructive" />
              <CardTitle className="mt-4">Authentication Failed</CardTitle>
              <CardDescription className="text-destructive">
                {error}
              </CardDescription>
            </>
          )}
        </CardHeader>

        {status === 'error' && (
          <CardContent>
            <button
              onClick={() => navigate({ to: '/login' })}
              className="w-full rounded-md bg-primary px-4 py-2 text-primary-foreground hover:bg-primary/90"
            >
              Back to Login
            </button>
          </CardContent>
        )}
      </Card>
    </div>
  );
}
