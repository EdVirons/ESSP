/**
 * Mobile Responsiveness Visual Testing Configuration
 *
 * Defines viewport sizes and role-page mappings for automated screenshot testing
 */

// Mobile viewport widths to test (covering wide range of devices)
export const MOBILE_VIEWPORTS = [
  { width: 320, name: 'xs', label: 'iPhone SE (320px)' },
  { width: 360, name: 'sm', label: 'Samsung Galaxy (360px)' },
  { width: 390, name: 'md', label: 'iPhone 14 (390px)' },
  { width: 428, name: 'lg', label: 'iPhone 14 Pro Max (428px)' },
] as const;

// Standard mobile height
export const MOBILE_HEIGHT = 844;

// Priority roles and their accessible pages (field-facing roles first)
export const ROLE_PAGES = {
  field_tech: {
    label: 'Field Technician',
    pages: [
      { path: '/overview', name: 'overview' },
      { path: '/work-orders', name: 'work-orders' },
      { path: '/knowledge-base', name: 'knowledge-base' },
      { path: '/messages', name: 'messages' },
      { path: '/profile', name: 'profile' },
    ],
  },
  lead_tech: {
    label: 'Lead Technician',
    pages: [
      { path: '/overview', name: 'overview' },
      { path: '/work-orders', name: 'work-orders' },
      { path: '/service-shops', name: 'service-shops' },
      { path: '/schools', name: 'schools' },
      { path: '/devices', name: 'devices' },
      { path: '/parts-catalog', name: 'parts-catalog' },
      { path: '/knowledge-base', name: 'knowledge-base' },
      { path: '/reports', name: 'reports' },
      { path: '/messages', name: 'messages' },
      { path: '/profile', name: 'profile' },
    ],
  },
  school_contact: {
    label: 'School Contact',
    pages: [
      { path: '/overview', name: 'overview' },
      { path: '/incidents', name: 'incidents' },
      { path: '/school-inventory', name: 'school-inventory' },
      { path: '/profile', name: 'profile' },
    ],
  },
} as const;

export type RoleKey = keyof typeof ROLE_PAGES;
export type ViewportConfig = typeof MOBILE_VIEWPORTS[number];

// Screenshot output configuration
export const OUTPUT_CONFIG = {
  screenshotDir: 'mobile-screenshots',
  reportName: 'report.html',
  waitAfterNavigation: 2000, // ms to wait for page render
  waitForViewportChange: 500, // ms to wait after viewport change
};
