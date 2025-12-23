import { test, expect } from '@playwright/test';
import { loginAs, type UserRole } from '../fixtures/test-utils';
import { MOBILE_VIEWPORTS, MOBILE_HEIGHT, ROLE_PAGES, OUTPUT_CONFIG } from './config';
import * as fs from 'fs';
import * as path from 'path';
import { fileURLToPath } from 'url';
import { dirname } from 'path';

// ES module compatible __dirname
const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// Ensure screenshot directory exists
const screenshotBaseDir = path.join(__dirname, '..', '..', OUTPUT_CONFIG.screenshotDir);

test.describe('Mobile Responsiveness Visual Tests', () => {
  // Increase timeout for visual testing
  test.setTimeout(300000); // 5 minutes total

  // Setup: create screenshot directories
  test.beforeAll(async () => {
    if (!fs.existsSync(screenshotBaseDir)) {
      fs.mkdirSync(screenshotBaseDir, { recursive: true });
    }
    for (const roleKey of Object.keys(ROLE_PAGES)) {
      const roleDir = path.join(screenshotBaseDir, roleKey);
      if (!fs.existsSync(roleDir)) {
        fs.mkdirSync(roleDir, { recursive: true });
      }
    }
  });

  // Test each role
  for (const [roleKey, roleConfig] of Object.entries(ROLE_PAGES)) {
    test.describe(`Role: ${roleConfig.label}`, () => {
      // Test each page for this role
      for (const pageConfig of roleConfig.pages) {
        test(`${pageConfig.name} - capture all viewports`, async ({ browser }) => {
          // Create a new browser context
          const context = await browser.newContext({
            viewport: { width: MOBILE_VIEWPORTS[0].width, height: MOBILE_HEIGHT },
          });
          const page = await context.newPage();

          try {
            // Login as this role
            await loginAs(page, roleKey as UserRole);

            // Navigate to the test page
            await page.goto(pageConfig.path);
            await page.waitForLoadState('networkidle');
            await page.waitForTimeout(OUTPUT_CONFIG.waitAfterNavigation);

            // Capture screenshot at each viewport width
            for (const viewport of MOBILE_VIEWPORTS) {
              // Set viewport size
              await page.setViewportSize({
                width: viewport.width,
                height: MOBILE_HEIGHT,
              });

              // Wait for responsive layout adjustments
              await page.waitForTimeout(OUTPUT_CONFIG.waitForViewportChange);

              // Close mobile menu if open (to get clean screenshot)
              const overlay = page.locator('.fixed.inset-0.bg-black\\/50');
              if (await overlay.isVisible({ timeout: 500 }).catch(() => false)) {
                await overlay.click();
                await page.waitForTimeout(200);
              }

              // Take screenshot
              const screenshotPath = path.join(
                screenshotBaseDir,
                roleKey,
                `${pageConfig.name}_${viewport.width}.png`
              );

              await page.screenshot({
                path: screenshotPath,
                fullPage: true,
              });

              console.log(`  Captured: ${roleKey}/${pageConfig.name}_${viewport.width}.png`);
            }
          } finally {
            await context.close();
          }
        });
      }
    });
  }

  // After all tests, generate the HTML report
  test.afterAll(async () => {
    await generateReport();
  });
});

/**
 * Generate HTML comparison report
 */
async function generateReport(): Promise<void> {
  const reportPath = path.join(screenshotBaseDir, OUTPUT_CONFIG.reportName);

  const html = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>ESSP Mobile Responsiveness Report</title>
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      background: #f5f5f5;
      padding: 20px;
    }
    .header {
      background: linear-gradient(135deg, #06b6d4, #0d9488);
      color: white;
      padding: 30px;
      border-radius: 12px;
      margin-bottom: 30px;
    }
    .header h1 { font-size: 2rem; margin-bottom: 10px; }
    .header p { opacity: 0.9; }
    .timestamp { font-size: 0.875rem; opacity: 0.8; margin-top: 10px; }

    .summary {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      gap: 20px;
      margin-bottom: 30px;
    }
    @media (max-width: 768px) {
      .summary { grid-template-columns: 1fr; }
    }
    .summary-card {
      background: white;
      border-radius: 12px;
      padding: 20px;
      text-align: center;
      box-shadow: 0 1px 3px rgba(0,0,0,0.1);
    }
    .summary-number {
      font-size: 2.5rem;
      font-weight: 700;
      color: #06b6d4;
    }
    .summary-label {
      color: #6b7280;
      margin-top: 5px;
    }

    .role-section {
      background: white;
      border-radius: 12px;
      padding: 24px;
      margin-bottom: 30px;
      box-shadow: 0 1px 3px rgba(0,0,0,0.1);
    }
    .role-title {
      font-size: 1.5rem;
      color: #1f2937;
      margin-bottom: 20px;
      padding-bottom: 10px;
      border-bottom: 2px solid #e5e7eb;
    }

    .page-group {
      margin-bottom: 30px;
    }
    .page-title {
      font-size: 1.125rem;
      color: #4b5563;
      margin-bottom: 15px;
      display: flex;
      align-items: center;
      gap: 8px;
    }
    .page-path {
      font-size: 0.875rem;
      color: #9ca3af;
      font-family: monospace;
    }

    .viewport-grid {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 16px;
    }
    @media (max-width: 1200px) {
      .viewport-grid { grid-template-columns: repeat(2, 1fr); }
    }
    @media (max-width: 600px) {
      .viewport-grid { grid-template-columns: 1fr; }
    }

    .viewport-card {
      border: 1px solid #e5e7eb;
      border-radius: 8px;
      overflow: hidden;
      transition: box-shadow 0.2s;
    }
    .viewport-card:hover {
      box-shadow: 0 4px 12px rgba(0,0,0,0.15);
    }
    .viewport-header {
      background: #f9fafb;
      padding: 10px 12px;
      border-bottom: 1px solid #e5e7eb;
    }
    .viewport-width {
      font-weight: 600;
      color: #1f2937;
    }
    .viewport-label {
      font-size: 0.75rem;
      color: #6b7280;
    }
    .viewport-image {
      width: 100%;
      height: auto;
      display: block;
      cursor: zoom-in;
    }
    .viewport-image:hover {
      opacity: 0.95;
    }

    .missing {
      background: #fef2f2;
      color: #dc2626;
      padding: 40px;
      text-align: center;
    }

    .lightbox {
      display: none;
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background: rgba(0,0,0,0.9);
      z-index: 1000;
      justify-content: center;
      align-items: center;
    }
    .lightbox.active { display: flex; }
    .lightbox img {
      max-width: 95%;
      max-height: 95%;
      object-fit: contain;
    }
    .lightbox-close {
      position: absolute;
      top: 20px;
      right: 30px;
      color: white;
      font-size: 2rem;
      cursor: pointer;
    }

    .nav-tabs {
      display: flex;
      gap: 8px;
      margin-bottom: 20px;
      flex-wrap: wrap;
    }
    .nav-tab {
      padding: 8px 16px;
      background: #e5e7eb;
      border-radius: 6px;
      cursor: pointer;
      font-weight: 500;
      transition: all 0.2s;
    }
    .nav-tab:hover {
      background: #d1d5db;
    }
    .nav-tab.active {
      background: #06b6d4;
      color: white;
    }
  </style>
</head>
<body>
  <div class="header">
    <h1>ESSP Mobile Responsiveness Report</h1>
    <p>Visual comparison of dashboard pages across mobile viewport widths</p>
    <div class="timestamp">Generated: ${new Date().toLocaleString()}</div>
  </div>

  <div class="summary">
    <div class="summary-card">
      <div class="summary-number">${Object.keys(ROLE_PAGES).length}</div>
      <div class="summary-label">Roles Tested</div>
    </div>
    <div class="summary-card">
      <div class="summary-number">${Object.values(ROLE_PAGES).reduce((sum, r) => sum + r.pages.length, 0)}</div>
      <div class="summary-label">Pages Tested</div>
    </div>
    <div class="summary-card">
      <div class="summary-number">${MOBILE_VIEWPORTS.length}</div>
      <div class="summary-label">Viewport Widths</div>
    </div>
  </div>

  <div class="nav-tabs">
    ${Object.entries(ROLE_PAGES).map(([roleKey, roleConfig], index) => `
      <div class="nav-tab ${index === 0 ? 'active' : ''}" onclick="showRole('${roleKey}')">${roleConfig.label}</div>
    `).join('')}
  </div>

  ${Object.entries(ROLE_PAGES).map(([roleKey, roleConfig], index) => `
    <div class="role-section" id="role-${roleKey}" style="${index !== 0 ? 'display:none' : ''}">
      <h2 class="role-title">${roleConfig.label}</h2>
      ${roleConfig.pages.map(pageConfig => `
        <div class="page-group">
          <h3 class="page-title">
            ${pageConfig.name.replace(/-/g, ' ').replace(/\\b\\w/g, l => l.toUpperCase())}
            <span class="page-path">${pageConfig.path}</span>
          </h3>
          <div class="viewport-grid">
            ${MOBILE_VIEWPORTS.map(vp => {
              const imgPath = `${roleKey}/${pageConfig.name}_${vp.width}.png`;
              return `
                <div class="viewport-card">
                  <div class="viewport-header">
                    <div class="viewport-width">${vp.width}px</div>
                    <div class="viewport-label">${vp.label}</div>
                  </div>
                  <img class="viewport-image"
                       src="${imgPath}"
                       alt="${pageConfig.name} at ${vp.width}px"
                       onclick="openLightbox(this.src)"
                       onerror="this.parentElement.innerHTML='<div class=\\'missing\\'>Screenshot not found</div>'" />
                </div>
              `;
            }).join('')}
          </div>
        </div>
      `).join('')}
    </div>
  `).join('')}

  <div class="lightbox" onclick="closeLightbox()">
    <span class="lightbox-close">&times;</span>
    <img id="lightbox-img" src="" alt="Zoomed screenshot">
  </div>

  <script>
    function openLightbox(src) {
      document.getElementById('lightbox-img').src = src;
      document.querySelector('.lightbox').classList.add('active');
    }
    function closeLightbox() {
      document.querySelector('.lightbox').classList.remove('active');
    }
    function showRole(roleKey) {
      // Hide all role sections
      document.querySelectorAll('.role-section').forEach(el => el.style.display = 'none');
      // Show selected role section
      document.getElementById('role-' + roleKey).style.display = 'block';
      // Update tab active state
      document.querySelectorAll('.nav-tab').forEach(el => el.classList.remove('active'));
      event.target.classList.add('active');
    }
    document.addEventListener('keydown', (e) => {
      if (e.key === 'Escape') closeLightbox();
    });
  </script>
</body>
</html>`;

  fs.writeFileSync(reportPath, html);
  console.log(`\\nReport generated: ${reportPath}`);
}
