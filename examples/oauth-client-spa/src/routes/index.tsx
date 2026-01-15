import { Link, useNavigate } from '@tanstack/react-router';
import { createFileRoute } from '@tanstack/react-router';
import { Key, LayoutDashboard, Lock, Settings, Shield, User } from 'lucide-react';

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useAuthStore } from '@/lib/auth';

export const Route = createFileRoute('/')({
  component: HomePage,
});

const pages = [
  {
    to: '/dashboard',
    title: 'Dashboard',
    description: 'View statistics and analytics',
    icon: LayoutDashboard,
    protected: true,
  },
  {
    to: '/profile',
    title: 'Profile',
    description: 'Your account information',
    icon: User,
    protected: true,
  },
  {
    to: '/settings',
    title: 'Settings',
    description: 'Token and session info',
    icon: Settings,
    protected: true,
  },
];

function HomePage() {
  const navigate = useNavigate();
  const { isAuthenticated, user, accessToken } = useAuthStore();

  return (
    <div className="space-y-8">
      <div className="text-center space-y-4">
        <div className="flex items-center justify-center gap-3">
          <Shield className="h-10 w-10 text-primary" />
          <h1 className="text-4xl font-bold">OAuth SPA Demo</h1>
        </div>
        <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
          A React SPA demonstrating OAuth 2.0 authentication with PKCE flow
          using Altalune's auth server.
        </p>
      </div>

      {isAuthenticated && user && (
        <Card className="max-w-2xl mx-auto">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Key className="h-5 w-5" />
              Authenticated Session
            </CardTitle>
            <CardDescription>
              You are signed in as {user.name || user.email}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              <div>
                <p className="text-sm text-muted-foreground mb-1">Access Token (truncated)</p>
                <code className="block rounded bg-muted p-2 text-xs break-all">
                  {accessToken?.slice(0, 100)}...
                </code>
              </div>
              <div>
                <p className="text-sm text-muted-foreground mb-1">JWT Claims</p>
                <pre className="rounded bg-muted p-2 text-xs overflow-auto max-h-40">
                  {JSON.stringify(user, null, 2)}
                </pre>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      <div className="grid gap-4 md:grid-cols-3 max-w-4xl mx-auto">
        {pages.map((page) => {
          const isLocked = page.protected && !isAuthenticated;

          return (
            <Link key={page.to} to={page.to}>
              <Card className={`h-full transition-colors hover:border-primary ${isLocked ? 'opacity-60' : ''}`}>
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <page.icon className="h-6 w-6 text-primary" />
                    {isLocked && <Lock className="h-4 w-4 text-muted-foreground" />}
                  </div>
                  <CardTitle className="text-lg">{page.title}</CardTitle>
                  <CardDescription>{page.description}</CardDescription>
                </CardHeader>
              </Card>
            </Link>
          );
        })}
      </div>

      {!isAuthenticated && (
        <div className="text-center space-y-4">
          <p className="text-muted-foreground">
            Sign in to access protected pages
          </p>
          <Button size="lg" onClick={() => navigate({ to: '/login' })}>
            <Shield className="mr-2 h-5 w-5" />
            Login
          </Button>
        </div>
      )}

      <Card className="max-w-2xl mx-auto">
        <CardHeader>
          <CardTitle>Features</CardTitle>
          <CardDescription>
            This demo showcases OAuth 2.0 PKCE flow for SPAs
          </CardDescription>
        </CardHeader>
        <CardContent>
          <ul className="space-y-2 text-sm text-muted-foreground">
            <li>• PKCE (Proof Key for Code Exchange) for public clients</li>
            <li>• No client secret required (public client)</li>
            <li>• Cookie-based token storage for session persistence</li>
            <li>• Protected routes with automatic redirect to login</li>
            <li>• Multi-provider ready architecture (easily extensible)</li>
            <li>• Built with React, TanStack Router, and Zustand</li>
          </ul>
        </CardContent>
      </Card>
    </div>
  );
}
