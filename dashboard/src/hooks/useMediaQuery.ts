import { useSyncExternalStore, useCallback } from 'react';

/**
 * Tailwind CSS breakpoint values
 */
const BREAKPOINTS = {
  sm: 640,
  md: 768,
  lg: 1024,
  xl: 1280,
  '2xl': 1536,
} as const;

type Breakpoint = keyof typeof BREAKPOINTS;

/**
 * Hook to check if a media query matches
 * Uses useSyncExternalStore for proper React 18+ compatibility
 * @param query - CSS media query string (e.g., "(min-width: 768px)")
 * @returns boolean indicating if the query matches
 */
export function useMediaQuery(query: string): boolean {
  const subscribe = useCallback(
    (callback: () => void) => {
      const mq = window.matchMedia(query);
      mq.addEventListener('change', callback);
      return () => mq.removeEventListener('change', callback);
    },
    [query]
  );

  const getSnapshot = useCallback(() => {
    return window.matchMedia(query).matches;
  }, [query]);

  const getServerSnapshot = useCallback(() => {
    return false; // Default to false on server
  }, []);

  return useSyncExternalStore(subscribe, getSnapshot, getServerSnapshot);
}

/**
 * Hook to check if screen is at or above a specific breakpoint
 * @param breakpoint - Tailwind breakpoint name
 * @returns boolean indicating if screen is >= breakpoint
 */
export function useBreakpoint(breakpoint: Breakpoint): boolean {
  return useMediaQuery(`(min-width: ${BREAKPOINTS[breakpoint]}px)`);
}

/**
 * Hook to check if screen is below a specific breakpoint
 * @param breakpoint - Tailwind breakpoint name
 * @returns boolean indicating if screen is < breakpoint
 */
export function useBreakpointBelow(breakpoint: Breakpoint): boolean {
  return useMediaQuery(`(max-width: ${BREAKPOINTS[breakpoint] - 1}px)`);
}

/**
 * Hook to check if screen is mobile (below md breakpoint)
 * @returns boolean indicating if screen is mobile
 */
export function useIsMobile(): boolean {
  return useBreakpointBelow('md');
}

/**
 * Hook to check if screen is tablet (between md and lg breakpoints)
 * @returns boolean indicating if screen is tablet
 */
export function useIsTablet(): boolean {
  const isMd = useBreakpoint('md');
  const isLg = useBreakpoint('lg');
  return isMd && !isLg;
}

/**
 * Hook to check if screen is desktop (at or above lg breakpoint)
 * @returns boolean indicating if screen is desktop
 */
export function useIsDesktop(): boolean {
  return useBreakpoint('lg');
}

/**
 * Hook to get current breakpoint name
 * @returns current breakpoint or undefined for xs
 */
export function useCurrentBreakpoint(): Breakpoint | 'xs' {
  const isSm = useBreakpoint('sm');
  const isMd = useBreakpoint('md');
  const isLg = useBreakpoint('lg');
  const isXl = useBreakpoint('xl');
  const is2xl = useBreakpoint('2xl');

  if (is2xl) return '2xl';
  if (isXl) return 'xl';
  if (isLg) return 'lg';
  if (isMd) return 'md';
  if (isSm) return 'sm';
  return 'xs';
}
