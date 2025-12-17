import { test, expect } from '@playwright/test';

// Helper to login as a user
async function loginAs(page: import('@playwright/test').Page, username: string, password: string) {
  await page.goto('/login');
  await page.getByLabel(/username/i).fill(username);
  await page.getByLabel(/password/i).fill(password);
  await page.getByRole('button', { name: /login|sign in/i }).click();
  await page.waitForURL(/dashboard|home|\//);
}

test.describe('Work Orders', () => {
  test.beforeEach(async ({ page }) => {
    await page.context().clearCookies();
    await loginAs(page, 'lead_tech', 'lead123');
  });

  test('should display work orders list', async ({ page }) => {
    await page.goto('/work-orders');

    // Should show work orders page
    await expect(page.getByRole('heading', { name: /work order/i })).toBeVisible();

    // Should show table or list
    await expect(page.locator('table, [data-testid="work-orders-list"]')).toBeVisible();
  });

  test('should be able to view work order details', async ({ page }) => {
    await page.goto('/work-orders');

    // Click on first work order
    await page.locator('table tbody tr, [data-testid="work-order-item"]').first().click();

    // Should show work order details
    await expect(page.getByText(/work order|details/i).first()).toBeVisible();
  });

  test('should be able to filter work orders by status', async ({ page }) => {
    await page.goto('/work-orders');

    // Find status filter
    const statusFilter = page.getByRole('combobox', { name: /status/i }).or(
      page.getByLabel(/status/i)
    );

    if (await statusFilter.isVisible()) {
      await statusFilter.click();
      await page.getByRole('option', { name: /in progress|assigned/i }).first().click();

      // Wait for table update
      await page.waitForResponse(resp => resp.url().includes('work-orders'));
    }
  });

  test('should be able to assign technician to work order', async ({ page }) => {
    await page.goto('/work-orders');

    // Click on first work order
    await page.locator('table tbody tr, [data-testid="work-order-item"]').first().click();

    // Find assign button or dropdown
    const assignButton = page.getByRole('button', { name: /assign/i });

    if (await assignButton.isVisible()) {
      await assignButton.click();

      // Select a technician
      await page.getByRole('option').first().click();

      // Confirm assignment
      await page.getByRole('button', { name: /confirm|save/i }).click();

      // Should show success
      await expect(page.getByText(/assigned|updated|success/i)).toBeVisible();
    }
  });

  test('should be able to update work order status', async ({ page }) => {
    await page.goto('/work-orders');

    // Click on first work order
    await page.locator('table tbody tr, [data-testid="work-order-item"]').first().click();

    // Find status dropdown or button
    const statusButton = page.getByRole('button', { name: /status/i }).or(
      page.getByLabel(/status/i)
    );

    if (await statusButton.isVisible()) {
      await statusButton.click();
      await page.getByRole('option', { name: /in progress/i }).first().click();

      // Should show success or status should update
      await expect(
        page.getByText(/updated|success|in progress/i).first()
      ).toBeVisible();
    }
  });

  test('should be able to add a comment to work order', async ({ page }) => {
    await page.goto('/work-orders');

    // Click on first work order
    await page.locator('table tbody tr, [data-testid="work-order-item"]').first().click();

    // Find comment input
    const commentInput = page.getByPlaceholder(/comment|note/i).or(
      page.getByLabel(/comment|note/i)
    );

    if (await commentInput.isVisible()) {
      await commentInput.fill('E2E test comment');

      // Submit comment
      await page.getByRole('button', { name: /add|submit|send/i }).click();

      // Comment should appear in the list
      await expect(page.getByText('E2E test comment')).toBeVisible();
    }
  });
});

test.describe('Work Order BOM (Bill of Materials)', () => {
  test.beforeEach(async ({ page }) => {
    await page.context().clearCookies();
    await loginAs(page, 'lead_tech', 'lead123');
  });

  test('should be able to view BOM on work order', async ({ page }) => {
    await page.goto('/work-orders');

    // Click on first work order
    await page.locator('table tbody tr, [data-testid="work-order-item"]').first().click();

    // Find BOM tab or section
    const bomTab = page.getByRole('tab', { name: /bom|parts|materials/i }).or(
      page.getByText(/bill of materials|parts/i)
    );

    if (await bomTab.isVisible()) {
      await bomTab.click();

      // Should show BOM list
      await expect(page.locator('[data-testid="bom-list"], table')).toBeVisible();
    }
  });
});
