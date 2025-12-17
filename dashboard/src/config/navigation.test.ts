import { describe, it, expect } from 'vitest';
import { getDefaultRouteForRoles, roleDefaultRoutes, navItems, navGroups } from './navigation';

describe('getDefaultRouteForRoles', () => {
  it('should return /incidents for school_contact', () => {
    expect(getDefaultRouteForRoles(['ssp_school_contact'])).toBe('/incidents');
  });

  it('should return /incidents for support_agent', () => {
    expect(getDefaultRouteForRoles(['ssp_support_agent'])).toBe('/incidents');
  });

  it('should return /work-orders for lead_tech', () => {
    expect(getDefaultRouteForRoles(['ssp_lead_tech'])).toBe('/work-orders');
  });

  it('should return /work-orders for field_tech', () => {
    expect(getDefaultRouteForRoles(['ssp_field_tech'])).toBe('/work-orders');
  });

  it('should return /parts-catalog for warehouse_manager', () => {
    expect(getDefaultRouteForRoles(['ssp_warehouse_manager'])).toBe('/parts-catalog');
  });

  it('should return /projects for demo_team', () => {
    expect(getDefaultRouteForRoles(['ssp_demo_team'])).toBe('/projects');
  });

  it('should return /sales for sales_marketing', () => {
    expect(getDefaultRouteForRoles(['ssp_sales_marketing'])).toBe('/sales');
  });

  it('should return /work-orders for supplier', () => {
    expect(getDefaultRouteForRoles(['ssp_supplier'])).toBe('/work-orders');
  });

  it('should return /work-orders for contractor', () => {
    expect(getDefaultRouteForRoles(['ssp_contractor'])).toBe('/work-orders');
  });

  it('should return /overview for empty roles', () => {
    expect(getDefaultRouteForRoles([])).toBe('/overview');
  });

  it('should return /overview for unknown roles', () => {
    expect(getDefaultRouteForRoles(['unknown_role'])).toBe('/overview');
  });

  it('should use first matching role in priority order', () => {
    // school_contact comes before lead_tech in roleDefaultRoutes
    expect(getDefaultRouteForRoles(['ssp_lead_tech', 'ssp_school_contact'])).toBe('/incidents');
  });

  it('should handle multiple roles with priority', () => {
    // The order in roleDefaultRoutes determines priority
    const roles = ['ssp_field_tech', 'ssp_support_agent'];
    const result = getDefaultRouteForRoles(roles);
    // support_agent comes before field_tech in the roleDefaultRoutes object
    expect(result).toBe('/incidents');
  });
});

describe('roleDefaultRoutes', () => {
  it('should have all expected roles defined', () => {
    const expectedRoles = [
      'ssp_school_contact',
      'ssp_support_agent',
      'ssp_lead_tech',
      'ssp_field_tech',
      'ssp_warehouse_manager',
      'ssp_demo_team',
      'ssp_sales_marketing',
      'ssp_supplier',
      'ssp_contractor',
    ];

    expectedRoles.forEach((role) => {
      expect(roleDefaultRoutes).toHaveProperty(role);
    });
  });

  it('should only contain valid routes', () => {
    const validRoutes = navItems.map((item) => item.href);

    Object.values(roleDefaultRoutes).forEach((route) => {
      expect(validRoutes).toContain(route);
    });
  });
});

describe('navItems', () => {
  it('should have required properties for each item', () => {
    navItems.forEach((item) => {
      expect(item).toHaveProperty('title');
      expect(item).toHaveProperty('href');
      expect(item).toHaveProperty('icon');
      expect(item).toHaveProperty('color');
      expect(item).toHaveProperty('bgColor');
    });
  });

  it('should have unique hrefs', () => {
    const hrefs = navItems.map((item) => item.href);
    const uniqueHrefs = new Set(hrefs);
    expect(hrefs.length).toBe(uniqueHrefs.size);
  });

  it('should have Overview as first item', () => {
    expect(navItems[0].title).toBe('Overview');
    expect(navItems[0].href).toBe('/overview');
  });

  it('should have Settings as last item', () => {
    const lastItem = navItems[navItems.length - 1];
    expect(lastItem.title).toBe('Settings');
    expect(lastItem.href).toBe('/settings');
  });
});

describe('navGroups', () => {
  it('should have required properties for each group', () => {
    navGroups.forEach((group) => {
      expect(group).toHaveProperty('id');
      expect(group).toHaveProperty('title');
      expect(group).toHaveProperty('icon');
      expect(group).toHaveProperty('color');
      expect(group).toHaveProperty('items');
      expect(Array.isArray(group.items)).toBe(true);
    });
  });

  it('should have unique group ids', () => {
    const ids = navGroups.map((group) => group.id);
    const uniqueIds = new Set(ids);
    expect(ids.length).toBe(uniqueIds.size);
  });

  it('should have main group first', () => {
    expect(navGroups[0].id).toBe('main');
  });

  it('should have admin group with proper roles', () => {
    const adminGroup = navGroups.find((g) => g.id === 'admin');
    expect(adminGroup).toBeDefined();
    expect(adminGroup?.roles).toContain('ssp_admin');
  });

  it('should have all items with valid properties', () => {
    navGroups.forEach((group) => {
      group.items.forEach((item) => {
        expect(item).toHaveProperty('title');
        expect(item).toHaveProperty('href');
        expect(item).toHaveProperty('icon');
        expect(typeof item.title).toBe('string');
        expect(item.href).toMatch(/^\//);
      });
    });
  });
});
