import { createFileRoute } from '@tanstack/react-router';
import { Clock, Key, RefreshCw, Shield } from 'lucide-react';

import { ProtectedRoute } from '@/components/protected-route';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useAuthStore } from '@/lib/auth';

export const Route = createFileRoute('/settings')({
  component: SettingsPage,
});

function SettingsPage() {
  const { accessToken, refreshToken, expiresAt, provider } = useAuthStore();

  const formatExpiry = (timestamp: number | null) => {
    if (!timestamp) return 'Unknown';
    const date = new Date(timestamp);
    return date.toLocaleString();
  };

  const truncateToken = (token: string | null, length = 20) => {
    if (!token) return 'N/A';
    if (token.length <= length) return token;
    return `${token.slice(0, length)}...`;
  };

  return (
    <ProtectedRoute>
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold">Settings</h1>
          <p className="text-muted-foreground">
            Authentication and session information
          </p>
        </div>

        <div className="grid gap-6 md:grid-cols-2">
          <Card>
            <CardHeader>
              <CardTitle>Session Information</CardTitle>
              <CardDescription>
                Details about your current session
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center gap-3">
                <div className="rounded-full bg-primary/10 p-2">
                  <Shield className="h-4 w-4 text-primary" />
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Provider</p>
                  <p className="font-medium capitalize">{provider || 'Unknown'}</p>
                </div>
              </div>

              <div className="flex items-center gap-3">
                <div className="rounded-full bg-primary/10 p-2">
                  <Clock className="h-4 w-4 text-primary" />
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Expires At</p>
                  <p className="font-medium">{formatExpiry(expiresAt)}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Token Information</CardTitle>
              <CardDescription>
                Your OAuth tokens (truncated for security)
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center gap-3">
                <div className="rounded-full bg-primary/10 p-2">
                  <Key className="h-4 w-4 text-primary" />
                </div>
                <div className="flex-1 overflow-hidden">
                  <p className="text-sm text-muted-foreground">Access Token</p>
                  <p className="truncate font-mono text-xs">
                    {truncateToken(accessToken, 40)}
                  </p>
                </div>
              </div>

              <div className="flex items-center gap-3">
                <div className="rounded-full bg-primary/10 p-2">
                  <RefreshCw className="h-4 w-4 text-primary" />
                </div>
                <div className="flex-1 overflow-hidden">
                  <p className="text-sm text-muted-foreground">Refresh Token</p>
                  <p className="truncate font-mono text-xs">
                    {truncateToken(refreshToken, 40)}
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Protected Content</CardTitle>
            <CardDescription>
              This content is only visible to authenticated users
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="rounded-lg border-2 border-dashed border-yellow-500/50 bg-yellow-50 p-6 text-center dark:bg-yellow-950/20">
              <p className="text-sm text-yellow-700 dark:text-yellow-300">
                You are viewing protected settings.
                This page requires authentication.
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </ProtectedRoute>
  );
}
