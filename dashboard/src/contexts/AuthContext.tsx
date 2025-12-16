import { createContext, useContext, useEffect, useState, useCallback, useMemo } from 'react';
import type { ReactNode } from 'react';
import { authApi } from '@/lib/api';
import type { User, LoginRequest, SSOUserProfile } from '@/lib/api';
import { getPermissionsForRoles } from '@/lib/permissions';

interface AuthContextType {
  user: User | null;
  profile: SSOUserProfile | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  isLoadingProfile: boolean;
  error: string | null;
  login: (credentials: LoginRequest) => Promise<boolean>;
  logout: () => Promise<void>;
  checkAuth: () => Promise<void>;
  fetchProfile: () => Promise<SSOUserProfile | null>;
  hasPermission: (permission: string) => boolean;
  hasRole: (role: string) => boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null);
  const [profile, setProfile] = useState<SSOUserProfile | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isLoadingProfile, setIsLoadingProfile] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const checkAuth = useCallback(async () => {
    try {
      setIsLoading(true);
      const response = await authApi.me();
      if (response.authenticated && response.user) {
        setUser(response.user);
        setError(null);
      } else {
        setUser(null);
        setProfile(null);
      }
    } catch (err) {
      setUser(null);
      setProfile(null);
    } finally {
      setIsLoading(false);
    }
  }, []);

  const fetchProfile = useCallback(async (): Promise<SSOUserProfile | null> => {
    if (!user) return null;

    try {
      setIsLoadingProfile(true);
      const response = await authApi.profile();
      if (response.profile) {
        setProfile(response.profile);
        return response.profile;
      }
      return null;
    } catch (err) {
      console.error('Failed to fetch profile:', err);
      return null;
    } finally {
      setIsLoadingProfile(false);
    }
  }, [user]);

  const login = useCallback(async (credentials: LoginRequest): Promise<boolean> => {
    try {
      setIsLoading(true);
      setError(null);
      const response = await authApi.login(credentials);
      if (response.success && response.user) {
        setUser(response.user);
        return true;
      }
      setError(response.message || 'Login failed');
      return false;
    } catch (err: any) {
      const message = err.response?.data?.message || err.message || 'Login failed';
      setError(message);
      return false;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const logout = useCallback(async () => {
    try {
      await authApi.logout();
    } catch (err) {
      // Ignore logout errors
    } finally {
      setUser(null);
      setProfile(null);
      setError(null);
    }
  }, []);

  // Get all roles from profile or user
  const userRoles = useMemo(() => {
    return profile?.roles || user?.roles || [];
  }, [profile, user]);

  // Get all permissions based on roles (or from profile if provided)
  const userPermissions = useMemo(() => {
    // If profile has explicit permissions, use those
    if (profile?.permissions?.length) {
      return profile.permissions;
    }
    // Otherwise, derive permissions from roles
    return getPermissionsForRoles(userRoles);
  }, [profile, userRoles]);

  // Check if user has a specific permission
  const hasPermission = useCallback((permission: string): boolean => {
    return userPermissions.includes(permission);
  }, [userPermissions]);

  // Check if user has a specific role
  const hasRole = useCallback((role: string): boolean => {
    return userRoles.includes(role);
  }, [userRoles]);

  // Check authentication on mount
  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  // Periodically refresh the token
  useEffect(() => {
    if (!user) return;

    const refreshInterval = setInterval(async () => {
      try {
        await authApi.refresh();
      } catch (err) {
        // If refresh fails, check auth status
        await checkAuth();
      }
    }, 30 * 60 * 1000); // Refresh every 30 minutes

    return () => clearInterval(refreshInterval);
  }, [user, checkAuth]);

  const value: AuthContextType = {
    user,
    profile,
    isAuthenticated: !!user,
    isLoading,
    isLoadingProfile,
    error,
    login,
    logout,
    checkAuth,
    fetchProfile,
    hasPermission,
    hasRole,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextType {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
