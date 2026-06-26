import { fileURLToPath } from 'node:url';
import { mkdirSync, statSync } from 'node:fs';
import path from 'node:path';
import { chromium } from 'playwright';

const miniprogramRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
const outputDir = path.join(miniprogramRoot, 'docs/validation/visual-parity/phase-1/screenshots');

const viewports = [
  { name: '390x844', width: 390, height: 844 },
  { name: '519x927', width: 519, height: 927 }
];

const pagePairs = [
  {
    key: 'home',
    originalUrl: 'http://127.0.0.1:8000/index.html',
    migratedUrl: 'http://127.0.0.1:5174/#/pages/home/index',
    originalPage: 'home'
  },
  {
    key: 'profile',
    originalUrl: 'http://127.0.0.1:8000/index.html',
    migratedUrl: 'http://127.0.0.1:5174/#/pages/profile/index',
    originalPage: 'mine'
  },
  {
    key: 'orders',
    originalUrl: 'http://127.0.0.1:8000/index.html',
    migratedUrl: 'http://127.0.0.1:5174/#/pages/orders/index',
    originalPage: 'orders'
  },
  {
    key: 'ticket',
    originalUrl: 'http://127.0.0.1:8000/ticket.html',
    migratedUrl: 'http://127.0.0.1:5174/#/pages/ticket/index'
  },
  {
    key: 'island-detail',
    originalUrl: 'http://127.0.0.1:8000/island-cruise.html',
    migratedUrl: 'http://127.0.0.1:5174/#/pages/island-cruise/index?step=detail'
  },
  {
    key: 'island-booking',
    originalUrl: 'http://127.0.0.1:8000/island-cruise-booking.html',
    migratedUrl: 'http://127.0.0.1:5174/#/pages/island-cruise/index?step=detail'
  },
  {
    key: 'driver',
    originalUrl: 'http://127.0.0.1:8000/driver.html',
    migratedUrl: 'http://127.0.0.1:5174/#/subpackages/driver/pages/home/index'
  }
];

async function prepareOriginal(page, pageName) {
  if (pageName && pageName !== 'home') {
    await page.evaluate((targetPageName) => {
      if (typeof window.switchPage === 'function') {
        window.switchPage(targetPageName);
      }
      window.scrollTo(0, 0);
    }, pageName);
    await page.waitForTimeout(350);
    return;
  }

  await page.evaluate(() => window.scrollTo(0, 0));
}

async function capture(page, kind, pagePair, viewport) {
  await page.setViewportSize({ width: viewport.width, height: viewport.height });
  await page.goto(kind === 'original' ? pagePair.originalUrl : pagePair.migratedUrl, {
    waitUntil: 'networkidle',
    timeout: 30000
  });

  if (kind === 'original') {
    await prepareOriginal(page, pagePair.originalPage);
  }

  await page.waitForTimeout(1000);
  await page.evaluate(() => window.scrollTo(0, 0));

  const filename = `${kind}-${pagePair.key}-${viewport.name}.png`;
  const screenshotPath = path.join(outputDir, filename);
  await page.screenshot({ path: screenshotPath, fullPage: false });
  const { size } = statSync(screenshotPath);
  console.log(`${filename}\t${size}`);
}

mkdirSync(outputDir, { recursive: true });

const browser = await chromium.launch({ headless: true });
const page = await browser.newPage();

try {
  for (const viewport of viewports) {
    for (const pagePair of pagePairs) {
      await capture(page, 'original', pagePair, viewport);
      await capture(page, 'migrated', pagePair, viewport);
    }
  }
} finally {
  await browser.close();
}
