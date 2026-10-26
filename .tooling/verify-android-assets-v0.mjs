// Inspect the actual APK, not just Flutter's intermediate asset manifest.
import { createHash } from 'node:crypto';
import { execFileSync } from 'node:child_process';
import { synchronizeDemoAssets } from './sync-android-assets-v0.mjs';

export function verifyBundledDemoAssets({ apkPath, root }) {
  const expected = synchronizeDemoAssets({ root, checkOnly: true });
  const entries = execFileSync('unzip', ['-Z1', apkPath], { encoding: 'utf8', maxBuffer: 1024 * 1024 }).split('\n');
  return expected.map(asset => {
    const apkEntry = `assets/flutter_assets/${asset.key}`;
    if (entries.filter(entry => entry === apkEntry).length !== 1) {
      throw Error(`APK demo asset missing or duplicated: ${apkEntry}`);
    }
    const bytes = execFileSync('unzip', ['-p', apkPath, apkEntry], { maxBuffer: 1024 * 1024 });
    const sha256 = createHash('sha256').update(bytes).digest('hex');
    if (sha256 !== asset.sha256) throw Error(`APK demo asset checksum mismatch: ${apkEntry}`);
    return { source: asset.source, key: asset.key, apkEntry, bytes: bytes.length, sha256 };
  });
}
