import { createContext, useContext, useCallback, useState } from 'react';
import type { ReactNode } from 'react';
import { api } from '@/api/client';
import { toast } from '@/lib/toast';

export interface ImpersonationTarget {
  userId: string;
  name: string;
  email: string;
  schools: string[];
}

interface ImpersonationContextType {
  isImpersonating: boolean;
  targetUser: ImpersonationTarget | null;
  reason: string;
  startImpersonation: (userId: string, reason?: string) => Promise<boolean>;
  stopImpersonation: () => void;
  setReason: (reason: string) => void;
}

const ImpersonationContext = createContext<ImpersonationContextType | undefined>(undefined);

// Storage key for persisting impersonation state
const IMPERSONATION_STORAGE_KEY = 'impersonation_state';

interface StoredImpersonationState {
  targetUser: ImpersonationTarget;
  reason: string;
}

interface ImpersonationProviderProps {
  children: ReactNode;
}

export function ImpersonationProvider({ children }: ImpersonationProviderProps) {
  // Initialize from localStorage if available
  const [targetUser, setTargetUser] = useState<ImpersonationTarget | null>(() => {
    try {
      const stored = localStorage.getItem(IMPERSONATION_STORAGE_KEY);
      if (stored) {
        const state: StoredImpersonationState = JSON.parse(stored);
        return state.targetUser;
      }
    } catch {
      // Ignore parse errors
    }
    return null;
  });

  const [reason, setReasonState] = useState<string>(() => {
    try {
      const stored = localStorage.getItem(IMPERSONATION_STORAGE_KEY);
      if (stored) {
        const state: StoredImpersonationState = JSON.parse(stored);
        return state.reason;
      }
    } catch {
      // Ignore parse errors
    }
    return '';
  });

  const startImpersonation = useCallback(async (userId: string, impersonationReason?: string): Promise<boolean> => {
    try {
      // Validate the impersonation target
      const response = await api.post<{
        valid: boolean;
        userId: string;
        name: string;
        email: string;
        schools: string[];
        error?: string;
      }>('/impersonate/validate', { targetUserId: userId });

      if (!response.valid) {
        toast.error('Cannot impersonate user', response.error || 'Unknown error');
        return false;
      }

      const target: ImpersonationTarget = {
        userId: response.userId,
        name: response.name,
        email: response.email,
        schools: response.schools,
      };

      const reasonToUse = impersonationReason || '';

      // Store in localStorage for persistence across page refreshes
      const state: StoredImpersonationState = {
        targetUser: target,
        reason: reasonToUse,
      };
      localStorage.setItem(IMPERSONATION_STORAGE_KEY, JSON.stringify(state));

      setTargetUser(target);
      setReasonState(reasonToUse);

      toast.success(`Now acting as ${target.name}`);
      return true;
    } catch (error) {
      console.error('Failed to start impersonation:', error);
      toast.error('Failed to start impersonation');
      return false;
    }
  }, []);

  const stopImpersonation = useCallback(() => {
    localStorage.removeItem(IMPERSONATION_STORAGE_KEY);
    setTargetUser(null);
    setReasonState('');
    toast.success('Stopped impersonation');
  }, []);

  const setReason = useCallback((newReason: string) => {
    setReasonState(newReason);
    // Update stored state
    if (targetUser) {
      const state: StoredImpersonationState = {
        targetUser,
        reason: newReason,
      };
      localStorage.setItem(IMPERSONATION_STORAGE_KEY, JSON.stringify(state));
    }
  }, [targetUser]);

  const value: ImpersonationContextType = {
    isImpersonating: !!targetUser,
    targetUser,
    reason,
    startImpersonation,
    stopImpersonation,
    setReason,
  };

  return (
    <ImpersonationContext.Provider value={value}>
      {children}
    </ImpersonationContext.Provider>
  );
}

export function useImpersonation(): ImpersonationContextType {
  const context = useContext(ImpersonationContext);
  if (context === undefined) {
    throw new Error('useImpersonation must be used within an ImpersonationProvider');
  }
  return context;
}

// Hook to get impersonation headers for API calls
export function useImpersonationHeaders(): Record<string, string> {
  const { isImpersonating, targetUser, reason } = useImpersonation();

  if (!isImpersonating || !targetUser) {
    return {};
  }

  const headers: Record<string, string> = {
    'X-Impersonate-User': targetUser.userId,
  };

  if (reason) {
    headers['X-Impersonate-Reason'] = reason;
  }

  return headers;
}
