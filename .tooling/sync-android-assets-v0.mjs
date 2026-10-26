// Mechanical generated copies keep shared fixtures/tokens canonical while
// giving Flutter safe project-local asset keys. No real customer data belongs here.
import { createHash } from 'node:crypto';
import { copyFileSync, existsSync, mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
const defaultRoot = resolve(dirname(fileURLToPath(import.meta.url)), '..');
export const demoAssets = [
  { source: 'prototypes/fixtures/booking-v0.json', target: 'prototypes/android_business/assets/demo/booking-v0.json', key: 'assets/demo/booking-v0.json' },
  { source: 'design/tokens.json', target: 'prototypes/android_business/assets/demo/tokens.json', key: 'assets/demo/tokens.json' }
];
const hash = path => createHash('sha256').update(readFileSync(path)).digest('hex');
export function synchronizeDemoAssets({ root = defaultRoot, checkOnly = false } = {}) {
  return demoAssets.map(({ source, target, key }) => {
    const from = join(root, source), to = join(root, target), expected = hash(from);
    if (!checkOnly) {
      mkdirSync(dirname(to), { recursive: true });
      if (!existsSync(to) || hash(to) !== expected) copyFileSync(from, to);
    }
    if (!existsSync(to) || hash(to) !== expected) throw Error(`Generated demo asset missing/stale: ${target}`);
    return { source, target, key, sha256: expected };
  });
}
if (process.argv[1] && resolve(process.argv[1]) === fileURLToPath(import.meta.url)) {
  const checkOnly = process.argv.includes('--check');
  const assets = synchronizeDemoAssets({ checkOnly });
  const receipt = { checkedAt: new Date().toISOString(), status: 'PASS_GENERATED_ASSET_MATCH', checkOnly, assets };
  if (!checkOnly) {
    const target = join(defaultRoot, 'prototypes/android_business/artifacts/generated-assets-receipt.json');
    mkdirSync(dirname(target), { recursive: true });
    writeFileSync(target, JSON.stringify(receipt, null, 2) + '\n');
  }
  console.log(JSON.stringify(receipt, null, 2));
}
