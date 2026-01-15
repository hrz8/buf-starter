import { Link, useLocation, useNavigate } from '@tanstack/react-router';
import { LayoutDashboard, Lock, LogOut, Settings, Shield, User } from 'lucide-react';

import { useAuthStore } from '@/lib/auth';

import { Button } from './ui/button';

const navItems = [
  { to: '/', label: 'Home', icon: null, protected: false },
  { to: '/dashboard', label: 'Dashboard', icon: LayoutDashboard, protected: true },
  { to: '/profile', label: 'Profile', icon: User, protected: true },
  { to: '/settings', label: 'Settings', icon: Settings, protected: true },
];

export function Navbar() {
  const location = useLocation();
  const navigate = useNavigate();
  const { isAuthenticated, user, logout } = useAuthStore();

  const handleLogin = () => {
    navigate({ to: '/login' });
  };

  return (
    <nav className="border-b bg-background">
      <div className="container mx-auto flex h-14 items-center justify-between px-4">
        <div className="flex items-center gap-6">
          <Link to="/" className="flex items-center gap-2 font-semibold">
            <Shield className="h-5 w-5 text-primary" />
            <span>OAuth SPA Demo</span>
          </Link>

          <div className="flex items-center gap-1">
            {navItems.map((item) => {
              const isActive = location.pathname === item.to;
              const isLocked = item.protected && !isAuthenticated;

              return (
                <Link
                  key={item.to}
                  to={item.to}
                  className={`
                    flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm
                    transition-colors
                    ${isActive
                      ? 'bg-primary/10 text-primary'
                      : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                    }
                    ${isLocked ? 'opacity-50' : ''}
                  `}
                >
                  {item.icon && <item.icon className="h-4 w-4" />}
                  <span>{item.label}</span>
                  {isLocked && <Lock className="h-3 w-3" />}
                </Link>
              );
            })}
          </div>
        </div>

        <div className="flex items-center gap-4">
          {isAuthenticated ? (
            <>
              <span className="text-sm text-muted-foreground">
                {user?.name || user?.email || 'User'}
              </span>
              <Button variant="outline" size="sm" onClick={logout}>
                <LogOut className="mr-1.5 h-4 w-4" />
                Logout
              </Button>
            </>
          ) : (
            <Button size="sm" onClick={handleLogin}>
              Login
            </Button>
          )}
        </div>
      </div>
    </nav>
  );
}
