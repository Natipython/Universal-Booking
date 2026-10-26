// One-time recoverable checkpoint before the approved NATIVE-UX-01 layout fix.
// Copies known synthetic build/test inputs and artifacts only; no device access.
import { createHash } from 'node:crypto';
import { copyFileSync, existsSync, mkdirSync, readFileSync, readdirSync, statSync, writeFileSync } from 'node:fs';
import { dirname, join, relative, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const project = join(root, 'prototypes/android_business');
const checkpoint = join(project, 'artifacts/checkpoints/2026-10-15-pre-layout');
if (existsSync(checkpoint)) throw Error('Checkpoint already exists; refusing to overwrite historical evidence');
const sha = path => createHash('sha256').update(readFileSync(path)).digest('hex');
const buildReceiptPath = 'prototypes/android_business/artifacts/android-build/build-receipt.json';
const build = JSON.parse(readFileSync(join(root, buildReceiptPath)));
if (build.status !== 'PASS_LOCAL_DEBUG_APK_BUILD' || sha(join(root, build.apk.path)) !== build.apk.sha256) throw Error('Baseline build/APK not verified');
for (const source of build.sources) if (sha(join(root, source.path)) !== source.sha256) throw Error(`Baseline input drift: ${source.path}`);
const walk = dir => readdirSync(dir).flatMap(name => statSync(join(dir, name)).isDirectory() ? walk(join(dir, name)) : [join(dir, name)]);
const hostSources = ['lib', 'test', 'tool', 'assets'].flatMap(dir => walk(join(project, dir)))
  .concat(['pubspec.yaml', 'pubspec.lock', 'analysis_options.yaml', 'android/app/src/main/AndroidManifest.xml', 'android/app/build.gradle.kts'].map(path => join(project, path)))
  .concat(['.tooling/validate-v0.mjs', '.tooling/sync-android-assets-v0.mjs', 'prototypes/fixtures/booking-v0.json', 'design/tokens.json'].map(path => join(root, path))).sort();
const host = JSON.parse(readFileSync(join(project, 'artifacts/validation-receipt.json')));
const hostDigest = createHash('sha256').update(JSON.stringify(hostSources.map(path => [relative(root, path), sha(path)]))).digest('hex');
if (host.status !== 'PASS' || host.sourceDigest !== hostDigest) throw Error('Host baseline receipt/source mismatch');
const buildLogs = readdirSync(join(project, 'artifacts/android-build')).filter(name => statSync(join(project, 'artifacts/android-build', name)).isFile())
  .map(name => `prototypes/android_business/artifacts/android-build/${name}`);
const originals = [...new Set(build.sources.map(source => source.path).concat(hostSources.map(path => relative(root, path)), buildLogs,
  ['validation-receipt.json', 'analyzer.txt', 'test-events.jsonl', 'test-stderr.txt', 'home.png', 'add-booking.png', 'review.png'].map(name => `prototypes/android_business/artifacts/${name}`),
  ['.tooling/verify-android-native-evidence-v0.mjs', 'prototypes/android_business/artifacts/android-native/2026-10-15-sm-g975w/validation-receipt.json', build.apk.path]))].sort();
const files = originals.map(originalPath => {
  const source = join(root, originalPath), target = join(checkpoint, 'inputs', originalPath);
  mkdirSync(dirname(target), { recursive: true }); copyFileSync(source, target);
  const sha256 = sha(source);
  if (sha(target) !== sha256) throw Error(`Checkpoint copy mismatch: ${originalPath}`);
  return { originalPath, snapshotPath: relative(root, target), sha256 };
});
const receipt = { createdAt: new Date().toISOString(), status: 'PASS_ARCHIVED_BASELINE',
  scope: 'Immutable-intent snapshot of the pre-layout debug APK, covered inputs and dated host/native/build receipts; not a new test run',
  buildReceiptPath, originalBuildSourceDigest: build.sourceDigest, originalHostSourceDigest: hostDigest,
  testedApkSha256: build.apk.sha256, files };
writeFileSync(join(checkpoint, 'checkpoint.json'), JSON.stringify(receipt, null, 2) + '\n');
console.log(JSON.stringify({ status: receipt.status, files: files.length, apkSha256: build.apk.sha256,
  checkpoint: relative(root, join(checkpoint, 'checkpoint.json')) }, null, 2));
