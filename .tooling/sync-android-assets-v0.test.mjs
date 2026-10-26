import { test } from 'node:test';
import assert from 'node:assert/strict';
import { existsSync, mkdirSync, mkdtempSync, readFileSync, rmSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { dirname, join } from 'node:path';
import { demoAssets, synchronizeDemoAssets } from './sync-android-assets-v0.mjs';

test('missing/stale generated assets fail read-only checks; sync copies canonical bytes', () => {
  const root = mkdtempSync(join(tmpdir(), 'ubi-asset-copy-test-'));
  try {
    for (const [i, asset] of demoAssets.entries()) {
      const path = join(root, asset.source);
      mkdirSync(dirname(path), { recursive: true });
      writeFileSync(path, JSON.stringify({ syntheticTest: i }));
    }
    assert.throws(() => synchronizeDemoAssets({ root, checkOnly: true }), /missing\/stale/);
    assert.equal(existsSync(join(root, demoAssets[0].target)), false);
    const result = synchronizeDemoAssets({ root });
    assert.equal(result.length, 2);
    for (const asset of demoAssets) {
      assert.deepEqual(readFileSync(join(root, asset.source)), readFileSync(join(root, asset.target)));
      assert.equal(asset.key.includes('..'), false);
    }
    assert.deepEqual(synchronizeDemoAssets({ root, checkOnly: true }), result);
    const stale = join(root, demoAssets[0].target);
    writeFileSync(stale, 'stale synthetic test data');
    assert.throws(() => synchronizeDemoAssets({ root, checkOnly: true }), /missing\/stale/);
    assert.equal(readFileSync(stale, 'utf8'), 'stale synthetic test data');
    synchronizeDemoAssets({ root });
    assert.deepEqual(readFileSync(stale), readFileSync(join(root, demoAssets[0].source)));
  } finally {
    rmSync(root, { recursive: true }); // Only this test's newly-created disposable directory.
  }
});
