import { Navigate, useLocation } from '@tanstack/react-router';
import { Loader2 } from 'lucide-react';
import { type ReactNode, useEffect, useState } from 'react';

import { useAuthStore } from '@/lib/auth';

interface ProtectedRouteProps {
  children: ReactNode;
}

export function ProtectedRoute({ children }: ProtectedRouteProps) {
  const location = useLocation();
  const { isAuthenticated, isLoading, isRefreshing, checkAndRefreshIfNeeded } = useAuthStore();
  const [isChecking, setIsChecking] = useState(true);

  useEffect(() => {
    let mounted = true;

    const checkAuth = async () => {
      setIsChecking(true);
      await checkAndRefreshIfNeeded();
      if (mounted) {
        setIsChecking(false);
      }
    };

    checkAuth();

    return () => {
      mounted = false;
    };
  }, [location.pathname, checkAndRefreshIfNeeded]);

  if (isLoading || isChecking || isRefreshing) {
    return (
      <div className="flex min-h-[50vh] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (!isAuthenticated) {
    return (
      <Navigate
        to="/login"
        search={{ returnTo: location.pathname }}
        replace
      />
    );
  }

  return <>{children}</>;
}
