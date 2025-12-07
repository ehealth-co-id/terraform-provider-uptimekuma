/**
 * Playwright script to automate Uptime Kuma admin account setup.
 * This script is used in CI to configure the initial admin account
 * since Uptime Kuma requires manual setup on first run.
 * 
 * Usage:
 *   node setup-uptime-kuma.js
 *   
 * Environment variables:
 *   UPTIME_KUMA_URL - The URL of the Uptime Kuma instance (default: http://localhost:3001)
 *   UPTIME_KUMA_ADMIN_USERNAME - The admin username to create (default: admin)
 *   UPTIME_KUMA_ADMIN_PASSWORD - The admin password to create (default: admin123)
 */

const { chromium } = require('playwright');

const UPTIME_KUMA_URL = process.env.UPTIME_KUMA_URL || 'http://localhost:3001';
const ADMIN_USERNAME = process.env.UPTIME_KUMA_ADMIN_USERNAME || 'admin';
const ADMIN_PASSWORD = process.env.UPTIME_KUMA_ADMIN_PASSWORD || 'admin123';

async function waitForUptimeKuma(maxRetries = 30, retryDelay = 2000) {
  console.log(`Waiting for Uptime Kuma at ${UPTIME_KUMA_URL}...`);
  
  for (let i = 0; i < maxRetries; i++) {
    try {
      const response = await fetch(UPTIME_KUMA_URL);
      if (response.ok) {
        console.log('Uptime Kuma is ready!');
        return true;
      }
    } catch (error) {
      // Server not ready yet
    }
    console.log(`Retry ${i + 1}/${maxRetries}...`);
    await new Promise(resolve => setTimeout(resolve, retryDelay));
  }
  
  throw new Error(`Uptime Kuma did not become ready within ${maxRetries * retryDelay / 1000} seconds`);
}

async function setupAdmin() {
  console.log('Starting Uptime Kuma admin setup...');
  
  // Wait for Uptime Kuma to be ready
  await waitForUptimeKuma();
  
  const browser = await chromium.launch({
    headless: true,
  });
  
  try {
    const context = await browser.newContext();
    const page = await context.newPage();
    
    console.log(`Navigating to ${UPTIME_KUMA_URL}...`);
    await page.goto(UPTIME_KUMA_URL, { waitUntil: 'networkidle', timeout: 60000 });
    
    // Check if we're on the setup page or login page
    const url = page.url();
    console.log(`Current URL: ${url}`);
    
    // If not on setup page, already configured
    if (!url.includes('/setup')) {
      console.log('Uptime Kuma is already configured (not on setup page).');
      return;
    }
    
    // Wait for the page to load
    await page.waitForLoadState('domcontentloaded');
    
    // Look for the username field using ID selector
    const usernameField = await page.$('#floatingInput');
    if (!usernameField) {
      console.log('Could not find username input. Taking screenshot for debugging...');
      await page.screenshot({ path: 'uptime-kuma-debug.png', fullPage: true });
      throw new Error('Username input not found');
    }
    
    console.log('Found setup form. Creating admin account...');
    
    // Fill in the username
    await page.fill('#floatingInput', ADMIN_USERNAME);
    console.log(`Entered username: ${ADMIN_USERNAME}`);
    
    // Fill in the password
    await page.fill('#floatingPassword', ADMIN_PASSWORD);
    console.log('Entered password');
    
    // Fill in the repeat password
    await page.fill('#repeat', ADMIN_PASSWORD);
    console.log('Entered repeat password');
    
    // Submit the form by clicking the submit button
    await page.click('button[type="submit"]');
    console.log('Clicked Create button');
    
    // Wait for navigation
    await page.waitForURL('**/dashboard**', { timeout: 30000 });
    
    // Verify we're on the dashboard
    const newUrl = page.url();
    console.log(`After submit URL: ${newUrl}`);
    
    if (newUrl.includes('/dashboard')) {
      console.log('Admin account created successfully!');
    } else {
      console.log('Setup may have failed. Current URL:', newUrl);
      await page.screenshot({ path: 'uptime-kuma-after-setup.png' });
    }
    
  } finally {
    await browser.close();
  }
}

// Run the setup
setupAdmin()
  .then(() => {
    console.log('Setup completed successfully');
    process.exit(0);
  })
  .catch((error) => {
    console.error('Setup failed:', error);
    process.exit(1);
  });
