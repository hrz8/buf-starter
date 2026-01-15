import { createFileRoute, Navigate, useSearch } from '@tanstack/react-router';
import { Shield } from 'lucide-react';

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Separator } from '@/components/ui/separator';
import { useAuthStore } from '@/lib/auth';

interface LoginSearch {
  returnTo?: string;
}

export const Route = createFileRoute('/login')({
  validateSearch: (search: Record<string, unknown>): LoginSearch => ({
    returnTo: typeof search.returnTo === 'string' ? search.returnTo : undefined,
  }),
  component: LoginPage,
});

function LoginPage() {
  const { returnTo } = useSearch({ from: '/login' });
  const { isAuthenticated, login } = useAuthStore();

  if (isAuthenticated) {
    return <Navigate to={returnTo || '/'} />;
  }

  const handleOAuthLogin = () => {
    login(returnTo || '/');
  };

  return (
    <div className="flex min-h-[70vh] items-center justify-center">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <CardTitle className="text-2xl">Welcome Back</CardTitle>
          <CardDescription>
            Sign in to access protected resources
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          {/* Dummy login form (disabled) */}
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="user@example.com"
                disabled
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                type="password"
                placeholder="••••••••"
                disabled
              />
            </div>
            <Button className="w-full" disabled>
              Sign In
            </Button>
            <p className="text-center text-xs text-muted-foreground">
              Local authentication is disabled in this demo
            </p>
          </div>

          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <Separator className="w-full" />
            </div>
            <div className="relative flex justify-center text-xs uppercase">
              <span className="bg-card px-2 text-muted-foreground">
                Or continue with
              </span>
            </div>
          </div>

          {/* OAuth Login Button */}
          <Button
            variant="outline"
            className="w-full"
            onClick={handleOAuthLogin}
          >
            <Shield className="mr-2 h-4 w-4" />
            Login with Altalune
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
