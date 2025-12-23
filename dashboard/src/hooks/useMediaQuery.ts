import { useState, useEffect, useCallback } from 'react';

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
 * @param query - CSS media query string (e.g., "(min-width: 768px)")
 * @returns boolean indicating if the query matches
 */
export function useMediaQuery(query: string): boolean {
  const getMatches = useCallback((): boolean => {
    if (typeof window === 'undefined') {
      return false;
    }
    return window.matchMedia(query).matches;
  }, [query]);

  const [matches, setMatches] = useState<boolean>(getMatches);

  useEffect(() => {
    const mq = window.matchMedia(query);

    const handleChange = (event: MediaQueryListEvent) => {
      setMatches(event.matches);
    };

    // Set initial value
    setMatches(mq.matches);

    // Listen for changes
    mq.addEventListener('change', handleChange);

    return () => {
      mq.removeEventListener('change', handleChange);
    };
  }, [query]);

  return matches;
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
