import { test, expect } from '@playwright/test';

test.describe('Authentication', () => {
  test.beforeEach(async ({ page }) => {
    // Clear any existing auth state
    await page.context().clearCookies();
  });

  test('should display login page when not authenticated', async ({ page }) => {
    await page.goto('/');

    // Should redirect to login or show login form
    await expect(page).toHaveURL(/login/);
    await expect(page.getByRole('heading', { name: /login|sign in/i })).toBeVisible();
  });

  test('should show login form with username and password fields', async ({ page }) => {
    await page.goto('/login');

    await expect(page.getByLabel(/username/i)).toBeVisible();
    await expect(page.getByLabel(/password/i)).toBeVisible();
    await expect(page.getByRole('button', { name: /login|sign in/i })).toBeVisible();
  });

  test('should show error on invalid credentials', async ({ page }) => {
    await page.goto('/login');

    await page.getByLabel(/username/i).fill('invalid_user');
    await page.getByLabel(/password/i).fill('wrong_password');
    await page.getByRole('button', { name: /login|sign in/i }).click();

    // Wait for error message
    await expect(page.getByText(/invalid|error|failed/i)).toBeVisible();
  });

  test('should login successfully with valid credentials', async ({ page }) => {
    await page.goto('/login');

    // Use demo admin credentials
    await page.getByLabel(/username/i).fill('admin');
    await page.getByLabel(/password/i).fill('admin123');
    await page.getByRole('button', { name: /login|sign in/i }).click();

    // Should redirect to dashboard after successful login
    await expect(page).toHaveURL(/dashboard|home|\//);

    // User should be logged in - check for user menu or logout button
    await expect(page.getByText(/admin|logout|dashboard/i).first()).toBeVisible();
  });

  test('should logout successfully', async ({ page }) => {
    // First login
    await page.goto('/login');
    await page.getByLabel(/username/i).fill('admin');
    await page.getByLabel(/password/i).fill('admin123');
    await page.getByRole('button', { name: /login|sign in/i }).click();

    // Wait for dashboard to load
    await page.waitForURL(/dashboard|home|\//);

    // Find and click logout
    const userMenu = page.getByRole('button', { name: /user|profile|menu/i });
    if (await userMenu.isVisible()) {
      await userMenu.click();
    }

    await page.getByRole('button', { name: /logout|sign out/i }).click();

    // Should redirect to login
    await expect(page).toHaveURL(/login/);
  });
});

test.describe('Role-based Access', () => {
  test('admin user should see admin menu items', async ({ page }) => {
    await page.goto('/login');
    await page.getByLabel(/username/i).fill('admin');
    await page.getByLabel(/password/i).fill('admin123');
    await page.getByRole('button', { name: /login|sign in/i }).click();

    await page.waitForURL(/dashboard|home|\//);

    // Admin should see settings/admin options
    await expect(page.getByText(/settings|admin|configuration/i).first()).toBeVisible();
  });

  test('support agent should see appropriate menu items', async ({ page }) => {
    await page.goto('/login');
    await page.getByLabel(/username/i).fill('support_agent');
    await page.getByLabel(/password/i).fill('support123');
    await page.getByRole('button', { name: /login|sign in/i }).click();

    await page.waitForURL(/dashboard|home|\//);

    // Support agent should see tickets/incidents
    await expect(page.getByText(/incident|ticket|work order/i).first()).toBeVisible();
  });
});
