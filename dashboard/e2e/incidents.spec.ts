import { test, expect } from '@playwright/test';

// Helper to login as a user
async function loginAs(page: import('@playwright/test').Page, username: string, password: string) {
  await page.goto('/login');
  await page.getByLabel(/username/i).fill(username);
  await page.getByLabel(/password/i).fill(password);
  await page.getByRole('button', { name: /login|sign in/i }).click();
  await page.waitForURL(/dashboard|home|\//);
}

test.describe('Incidents', () => {
  test.beforeEach(async ({ page }) => {
    await page.context().clearCookies();
    await loginAs(page, 'support_agent', 'support123');
  });

  test('should display incidents list', async ({ page }) => {
    await page.goto('/incidents');

    // Should show incidents page
    await expect(page.getByRole('heading', { name: /incident/i })).toBeVisible();

    // Should show table or list of incidents
    await expect(page.locator('table, [data-testid="incidents-list"]')).toBeVisible();
  });

  test('should be able to create a new incident', async ({ page }) => {
    await page.goto('/incidents');

    // Click create button
    await page.getByRole('button', { name: /create|new|add/i }).click();

    // Should show create form
    await expect(page.getByRole('heading', { name: /create|new/i })).toBeVisible();

    // Fill in the form
    await page.getByLabel(/type/i).click();
    await page.getByRole('option', { name: /damage|hardware/i }).first().click();

    await page.getByLabel(/description/i).fill('Test incident created by E2E test');

    // Submit the form
    await page.getByRole('button', { name: /submit|create|save/i }).click();

    // Should show success message or redirect to incident details
    await expect(
      page.getByText(/success|created|saved/i).or(page.locator('[data-testid="incident-details"]'))
    ).toBeVisible();
  });

  test('should be able to view incident details', async ({ page }) => {
    await page.goto('/incidents');

    // Click on first incident in the list
    await page.locator('table tbody tr, [data-testid="incident-item"]').first().click();

    // Should show incident details
    await expect(page.getByText(/incident|details/i).first()).toBeVisible();
  });

  test('should be able to filter incidents by status', async ({ page }) => {
    await page.goto('/incidents');

    // Find and click status filter
    const statusFilter = page.getByRole('combobox', { name: /status/i }).or(
      page.getByLabel(/status/i)
    );

    if (await statusFilter.isVisible()) {
      await statusFilter.click();
      await page.getByRole('option', { name: /open|pending/i }).first().click();

      // Table should update (might need to wait for API response)
      await page.waitForResponse(resp => resp.url().includes('incidents'));
    }
  });

  test('should be able to search incidents', async ({ page }) => {
    await page.goto('/incidents');

    // Find search input
    const searchInput = page.getByPlaceholder(/search/i).or(
      page.getByRole('searchbox')
    );

    if (await searchInput.isVisible()) {
      await searchInput.fill('test');
      await searchInput.press('Enter');

      // Wait for search results
      await page.waitForResponse(resp => resp.url().includes('incidents'));
    }
  });
});
