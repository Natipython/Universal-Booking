import { test } from 'node:test';
import assert from 'node:assert/strict';
import { execFileSync } from 'node:child_process';
import { copyFileSync, mkdirSync, mkdtempSync, rmSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { dirname, join } from 'node:path';
import { demoAssets, synchronizeDemoAssets } from './sync-android-assets-v0.mjs';
import { verifyBundledDemoAssets } from './verify-android-assets-v0.mjs';

function withArchive(mode, check) {
  const root = mkdtempSync(join(tmpdir(), 'ubi apk asset test '));
  try {
    for (const [i, asset] of demoAssets.entries()) {
      const path = join(root, asset.source);
      mkdirSync(dirname(path), { recursive: true });
      writeFileSync(path, JSON.stringify({ syntheticTest: i }));
    }
    synchronizeDemoAssets({ root });
    const staged = join(root, 'zip staging');
    const entries = [];
    for (const [i, asset] of demoAssets.entries()) {
      if (mode === 'missing' && i === 0) continue;
      const entry = `assets/flutter_assets/${asset.key}`;
      const target = join(staged, entry);
      mkdirSync(dirname(target), { recursive: true });
      copyFileSync(join(root, asset.target), target);
      if (mode === 'stale' && i === 0) writeFileSync(target, 'stale fixture');
      entries.push(entry);
    }
    const apkPath = join(root, 'synthetic zip.apk');
    execFileSync('zip', ['-q', apkPath, ...entries], { cwd: staged });
    check({ root, apkPath });
  } finally {
    rmSync(root, { recursive: true }); // Only this test's newly-created disposable directory.
  }
}

test('APK asset check verifies actual archived bytes and handles spaces', () => {
  withArchive('valid', input => {
    const result = verifyBundledDemoAssets(input);
    assert.equal(result.length, 2);
    assert.ok(result.every(asset => asset.bytes > 0 && asset.sha256.length === 64));
  });
});
test('APK asset check rejects missing startup data', () => {
  withArchive('missing', input => assert.throws(() => verifyBundledDemoAssets(input), /missing or duplicated/));
});
test('APK asset check rejects stale startup data', () => {
  withArchive('stale', input => assert.throws(() => verifyBundledDemoAssets(input), /checksum mismatch/));
});
