import { test as base, expect, Page } from '@playwright/test';

// Demo user credentials
export const users = {
  admin: { username: 'admin', password: 'admin123' },
  support_agent: { username: 'support_agent', password: 'support123' },
  lead_tech: { username: 'lead_tech', password: 'lead123' },
  field_tech: { username: 'field_tech', password: 'tech123' },
  warehouse: { username: 'warehouse', password: 'warehouse123' },
  school_contact: { username: 'school_contact', password: 'school123' },
  sales_marketing: { username: 'sales_marketing', password: 'sales123' },
} as const;

export type UserRole = keyof typeof users;

/**
 * Login as a specific user role
 */
export async function loginAs(page: Page, role: UserRole) {
  const user = users[role];
  await page.goto('/login');
  await page.getByLabel(/username/i).fill(user.username);
  await page.getByLabel(/password/i).fill(user.password);
  await page.getByRole('button', { name: /login|sign in/i }).click();
  await page.waitForURL(/dashboard|home|\//);
}

/**
 * Logout the current user
 */
export async function logout(page: Page) {
  // Try to find and click user menu
  const userMenu = page.getByRole('button', { name: /user|profile|menu/i });
  if (await userMenu.isVisible()) {
    await userMenu.click();
  }

  await page.getByRole('button', { name: /logout|sign out/i }).click();
  await page.waitForURL(/login/);
}

/**
 * Wait for API response
 */
export async function waitForApi(page: Page, urlPattern: string | RegExp) {
  return page.waitForResponse((resp) => {
    if (typeof urlPattern === 'string') {
      return resp.url().includes(urlPattern);
    }
    return urlPattern.test(resp.url());
  });
}

/**
 * Check if element exists and is visible
 */
export async function isVisible(page: Page, selector: string): Promise<boolean> {
  try {
    await page.waitForSelector(selector, { state: 'visible', timeout: 3000 });
    return true;
  } catch {
    return false;
  }
}

/**
 * Fill form fields from an object
 */
export async function fillForm(page: Page, fields: Record<string, string>) {
  for (const [label, value] of Object.entries(fields)) {
    const input = page.getByLabel(new RegExp(label, 'i'));
    if (await input.isVisible()) {
      await input.fill(value);
    }
  }
}

/**
 * Extended test fixture with authentication
 */
export const test = base.extend<{
  authenticatedPage: Page;
  loginAsUser: (role: UserRole) => Promise<void>;
}>({
  authenticatedPage: async ({ page }, use) => {
    // Login as admin by default
    await loginAs(page, 'admin');
    await use(page);
  },

  loginAsUser: async ({ page }, use) => {
    const login = async (role: UserRole) => {
      await loginAs(page, role);
    };
    await use(login);
  },
});

export { expect };
