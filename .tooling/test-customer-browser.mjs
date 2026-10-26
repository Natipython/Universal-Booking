import assert from 'node:assert/strict';
import { createRequire } from 'node:module';
import { mkdirSync, writeFileSync } from 'node:fs';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const origin = process.argv[2];
assert.match(origin ?? '', /^http:\/\/127\.0\.0\.1:\d+$/);
const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const artifacts = join(root, 'prototypes/customer_web/artifacts');
mkdirSync(artifacts, { recursive: true });
const require = createRequire('/home/n-a-t/.cache/codex-runtimes/codex-primary-runtime/dependencies/node/package.json');
const { chromium } = require('playwright');
const browser = await chromium.launch({ headless: true, executablePath: '/usr/bin/google-chrome' });
const checks = [];
const outsideRequests = [];
const traffic = [];
let step = 'initialization';
const check = (name, condition) => { assert.ok(condition, name); checks.push(name); };
async function newContext(viewport) {
  const context = await browser.newContext({ javaScriptEnabled: false, viewport });
  context.on('request', request => {
    if (new URL(request.url()).origin !== origin) outsideRequests.push(new URL(request.url()).origin);
  });
  context.on('response', response => traffic.push({ route: new URL(response.url()).pathname.replace(/\/h\/[^/]+$/, '/h/[OPAQUE]'), status: response.status() }));
  return context;
}
async function click(page, name) {
  step = name;
  await Promise.all([page.waitForNavigation(), page.getByRole('button', { name, exact: true }).click()]);
}
async function layout(page, label) {
  const size = await page.evaluate(() => ({ width: innerWidth, content: document.documentElement.scrollWidth,
    targets: [...document.querySelectorAll('button,.button,input:not([type=hidden]),select')].map(e => e.getBoundingClientRect().height),
    markWidth: document.querySelector('.mark').getBoundingClientRect().width,
    markLine: parseFloat(getComputedStyle(document.querySelector('.mark')).lineHeight),
    markHeight: document.querySelector('.mark').getBoundingClientRect().height }));
  check(`${label}: no horizontal page overflow`, size.content <= size.width + 1);
  check(`${label}: primary controls at least 48 CSS px high`, size.targets.length > 0 && size.targets.every(h => h >= 48));
  check(`${label}: decorative logo not compressed or vertically clipped`, size.markWidth >= 36 && size.markLine <= size.markHeight);
}
try {
  const context = await newContext({ width: 390, height: 844 });
  const page = await context.newPage();
  await page.goto(origin);
  await page.keyboard.press('Tab');
  check('keyboard first focus is skip link', await page.locator('.skip').evaluate(e => e === document.activeElement));
  await layout(page, '390px public');
  await page.screenshot({ path: join(artifacts, 'public-390.png'), fullPage: true });
  await click(page, 'Try a prepared booking');
  check('preview excludes the private date and customer', !(await page.textContent('main')).includes('2030-01-15') && !(await page.textContent('main')).includes('Demo Guest 01'));
  await click(page, 'Continue to my details');
  check('conversation provenance is explicit', (await page.textContent('main')).includes('not an actual Telegram or WhatsApp chat'));
  await page.getByLabel('Preferred date · Gregorian calendar').fill('2030-01-16');
  await page.getByLabel('Preferred time · international 24-hour clock').fill('10:00');
  await page.getByLabel('Practice customer label').selectOption('Demo Guest 02');
  await layout(page, '390px details');
  await click(page, 'Review my details');
  const reviewText = await page.textContent('main');
  check('edited review values visible', ['10:00', 'Demo Guest 02', '60 minutes', '2000.00 birr'].every(text => reviewText.includes(text)));
  await Promise.all([page.waitForNavigation(), page.getByRole('link', { name: 'Change my details' }).click()]);
  check('back-to-edit retains input', await page.locator('#date').inputValue() === '2030-01-16' && await page.locator('#time').inputValue() === '10:00' && await page.locator('#customer').inputValue() === 'Demo Guest 02');
  await page.locator('#time').fill('11:00');
  await click(page, 'Review my details');
  await layout(page, '390px review');
  await page.screenshot({ path: join(artifacts, 'review-390.png'), fullPage: true });
  await click(page, 'Confirm demo booking');
  check('instant flow reaches synthetic confirmation', await page.locator('[data-state="confirmed_demo"]').count() === 1);
  check('no payment or real appointment claim', (await page.textContent('main')).includes('No payment was taken') && (await page.textContent('main')).includes('not a real business booking'));
  await page.screenshot({ path: join(artifacts, 'confirmed-390.png'), fullPage: true });
  await page.goto(origin);
  await click(page, 'Try a custom enquiry');
  await click(page, 'Continue to my details');
  await page.locator('#date').fill('2030-01-17');
  await page.locator('#time').fill('11:00');
  await page.getByLabel('Approximate event size').selectOption('small');
  await click(page, 'Review my details');
  check('custom details appear before submission', (await page.textContent('main')).includes('small') && (await page.textContent('main')).includes('No time is held'));
  await click(page, 'Send demo enquiry');
  check('custom flow is request, not confirmation', await page.locator('[data-state="request_saved_demo"]').count() === 1 && await page.locator('[data-state="confirmed_demo"]').count() === 0);
  await page.screenshot({ path: join(artifacts, 'request-390.png'), fullPage: true });
  const narrow = await newContext({ width: 320, height: 700 });
  const small = await narrow.newPage();
  await small.goto(origin);
  await small.evaluate(() => { document.documentElement.style.fontSize = '200%'; });
  check('200% root font simulation applied', await small.evaluate(() => getComputedStyle(document.documentElement).fontSize) === '32px');
  await layout(small, '320px public at simulated 200% root text');
  await small.screenshot({ path: join(artifacts, 'large-text-320.png'), fullPage: true });
  await click(small, 'Try a prepared booking');
  await click(small, 'Continue to my details');
  await small.evaluate(() => { document.documentElement.style.fontSize = '200%'; });
  await layout(small, '320px details at simulated 200% root text');
  check('no third-party network requests', outsideRequests.length === 0);
  const receipt = { status: 'PASS_LOCAL_BROWSER_COMPONENT_CHECKS', checkedAt: new Date().toISOString(),
    browser: browser.version(), javaScriptEnabled: false, checks, outsideRequests,
    viewports: [{ width: 390, height: 844 }, { width: 320, height: 700, rootFontScale: 2 }],
    notVerified: ['Physical Android or iOS device', 'TalkBack or screen reader', 'Actual OS font scale or browser zoom', 'Weak mobile network', 'Full WCAG conformance', 'Real customer usability'] };
  writeFileSync(join(artifacts, 'browser-receipt.json'), JSON.stringify(receipt, null, 2) + '\n');
  console.log(JSON.stringify({ status: receipt.status, checksPassed: checks.length, javaScriptEnabled: false, externalRequests: outsideRequests.length }));
} catch (error) {
  // Avoid retaining private loopback handoff URLs from Playwright's diagnostic call log.
  console.error(String(error.message).replace(/http:\/\/127\.0\.0\.1:\d+[^\s"']*/g, '[LOCAL_URL]').replace(/[A-Za-z0-9_-]{43}/g, '[DEMO_SECRET]'));
  console.error(JSON.stringify({ step, recentResponses: traffic.slice(-4) }));
  process.exitCode = 1;
} finally { await browser.close(); }
