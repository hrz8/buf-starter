import { Outlet, createRootRoute } from '@tanstack/react-router';
import { useEffect } from 'react';

import { Navbar } from '@/components/navbar';
import { useAuthStore } from '@/lib/auth';

export const Route = createRootRoute({
  component: RootComponent,
});

function RootComponent() {
  const initFromCookies = useAuthStore((state) => state.initFromCookies);

  useEffect(() => {
    initFromCookies();
  }, [initFromCookies]);

  return (
    <div className="min-h-screen flex flex-col bg-background">
      <Navbar />
      <main className="flex-1 container mx-auto px-4 py-6">
        <Outlet />
      </main>
    </div>
  );
}
