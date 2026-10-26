// Offline provenance check for the bounded Home-action fix. No ADB/device calls.
// Geometry/semantics are exercised by Flutter tests, not inferred from PNGs.
import assert from 'node:assert/strict';
import { createHash } from 'node:crypto';
import { readFileSync, readdirSync, statSync, writeFileSync } from 'node:fs';
import { dirname, join, relative, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const project = join(root, 'prototypes/android_business');
const artifacts = join(project, 'artifacts');
const sha = path => createHash('sha256').update(readFileSync(path)).digest('hex');
const json = path => JSON.parse(readFileSync(path));
const walk = dir => readdirSync(dir).flatMap(name => {
  const path = join(dir, name);
  return statSync(path).isDirectory() ? walk(path) : [path];
});
const hostPath = join(artifacts, 'validation-receipt.json');
const buildPath = join(artifacts, 'android-build/build-receipt.json');
const checkpointPath = join(artifacts, 'checkpoints/2026-10-15-pre-layout/checkpoint.json');
const nativePath = join(artifacts, 'android-native/2026-10-15-sm-g975w/validation-receipt.json');
const host = json(hostPath), build = json(buildPath), checkpoint = json(checkpointPath), native = json(nativePath);
const frozenFiles = new Map(checkpoint.files.map(file => [file.originalPath, file]));
const oldBuild = json(join(root, frozenFiles.get(checkpoint.buildReceiptPath).snapshotPath));
const redPath = join(artifacts, 'layout/before-fix-events.jsonl');
const red = readFileSync(redPath, 'utf8').split('\n').flatMap(line => {
  try { return [JSON.parse(line)]; } catch { return []; }
});
const checks = [];
function check(name, verify) {
  try { verify(); checks.push({ name, result: 'PASS' }); }
  catch (error) { checks.push({ name, result: 'FAIL', reason: error.message }); }
}
check('Original 82-file checkpoint and phone-tested APK are intact', () => {
  assert.equal(checkpoint.status, 'PASS_ARCHIVED_BASELINE');
  assert.equal(checkpoint.files.length, 82);
  for (const file of checkpoint.files) assert.equal(sha(join(root, file.snapshotPath)), file.sha256, file.originalPath);
  assert.equal(oldBuild.apk.sha256, 'd42e02e366b255213b54cb1866b9127f8a4171f8d09f3c48e19dcd95d47b07f4');
});
check('The new regression really failed on the old overlapping layout', () => {
  assert.equal(red.findLast(event => event.type === 'done')?.success, false);
  assert.ok(red.some(event => event.type === 'error' &&
    event.error.includes('never cover booking content')));
});
check('Current host receipt passes all 41 tests with clean analyzer and unchanged covered sources', () => {
  assert.equal(host.status, 'PASS'); assert.equal(host.testsRun, 41); assert.equal(host.passed, 41);
  assert.equal(host.analyzerExitCode, 0); assert.equal(host.sourcesChangedDuringRun, false);
  assert.ok(host.tests.every(test => test.result === 'success' && !test.skipped));
  const files = ['lib', 'test', 'tool', 'assets'].flatMap(dir => walk(join(project, dir)))
    .concat(['pubspec.yaml', 'pubspec.lock', 'analysis_options.yaml',
      'android/app/src/main/AndroidManifest.xml', 'android/app/build.gradle.kts'].map(path => join(project, path)))
    .concat(['.tooling/validate-v0.mjs', '.tooling/sync-android-assets-v0.mjs',
      'prototypes/fixtures/booking-v0.json', 'design/tokens.json'].map(path => join(root, path))).sort();
  const digest = createHash('sha256').update(JSON.stringify(files.map(path => [relative(root, path), sha(path)]))).digest('hex');
  assert.equal(host.sourceDigest, digest);
});
check('All 12 viewport/state checks and the labelled action/navigation check passed', () => {
  const layouts = host.tests.filter(test => test.name.startsWith('reserved Add area / '));
  assert.equal(layouts.length, 12);
  for (const name of ['phone-sized', 'portrait', 'small-portrait', 'large-text', 'landscape', 'landscape-large-text']) {
    for (const state of ['empty', 'seeded']) assert.ok(layouts.some(test => test.name === `reserved Add area / ${name} / ${state}`));
  }
  assert.equal(host.tests.filter(test => test.name.startsWith('reserved Add action is labelled')).length, 1);
});
check('Current APK is the separately inspected 0.1.1 (2) debug candidate', () => {
  assert.equal(build.status, 'PASS_LOCAL_DEBUG_APK_BUILD'); assert.equal(build.failure, null);
  assert.equal(sha(join(root, build.apk.path)), build.apk.sha256);
  assert.notEqual(build.apk.sha256, oldBuild.apk.sha256);
  for (const source of build.sources) assert.equal(sha(join(root, source.path)), source.sha256, source.path);
  assert.ok(build.steps.every(step => step.exitCode === 0));
  assert.match(build.manifestSummary.package, /versionCode='2' versionName='0\.1\.1'/);
  assert.equal(build.manifestSummary.minSdk, 30); assert.equal(build.manifestSummary.targetSdk, 36);
  assert.equal(build.manifestSummary.allowBackup, false); assert.equal(build.manifestSummary.debuggable, true);
  assert.equal(build.nativeSqlitePluginIncluded, true);
  assert.equal(build.deviceInstalled, false); assert.equal(build.nativeSqliteExecuted, false);
  assert.equal(build.performanceMeasured, false); assert.equal(build.productionSigned, false);
});
const changedBuildInputs = oldBuild.sources.filter(source => sha(join(root, source.path)) !== source.sha256).map(source => source.path);
check('Only Home presentation and demo version changed among the prior covered build inputs', () => {
  assert.deepEqual([...changedBuildInputs].sort(), [
    'prototypes/android_business/lib/features/bookings/presentation/business_home.dart',
    'prototypes/android_business/pubspec.yaml',
  ]);
});
check('Phone session stays closed and native evidence still belongs only to the old APK', () => {
  const auth = json(join(artifacts, 'android-native/2026-10-15-sm-g975w/session-authorization.json'));
  assert.equal(auth.sessionStatus, 'CLOSED'); assert.equal(auth.testingConnectionStopped, true);
  assert.equal(native.status, 'PASS_BOUNDED_NATIVE_SMOKE'); assert.equal(native.testedApk.sha256, oldBuild.apk.sha256);
  assert.equal(native.workingSourcesStillMatchTestedBuild, false);
  assert.equal(native.testedBuildCheckpoint, relative(root, checkpointPath));
});
const captures = ['phone-sized', 'portrait', 'small-portrait', 'large-text', 'landscape', 'landscape-large-text']
  .map(name => join(artifacts, `layout/${name}-seeded.png`));
check('Six generated host renders exist as PNGs; not golden or native evidence', () => {
  for (const path of captures) assert.equal(readFileSync(path).subarray(0, 8).toString('hex'), '89504e470d0a1a0a');
});
const evidenceFiles = [hostPath, buildPath, checkpointPath, nativePath, redPath,
  join(artifacts, 'layout/apk-entry-comparison.json'), join(artifacts, 'layout/tooling-tests.txt'),
  fileURLToPath(import.meta.url), ...captures].map(path => ({ path: relative(root, path), sha256: sha(path) }));
const failed = checks.some(item => item.result !== 'PASS');
const report = {
  checkedAt: new Date().toISOString(), status: failed ? 'FAIL' : 'PASS_LOCAL_LAYOUT_FIX_NATIVE_RECHECK_PENDING',
  scope: 'Offline cross-check of the Home action fix, 41 host tests and new static debug build; not a new phone run',
  checks, changedBuildInputs, hostSourceDigest: host.sourceDigest, buildSourceDigest: build.sourceDigest,
  testedOldApkSha256: oldBuild.apk.sha256, candidateApk: build.apk,
  evidenceFiles, installedOnPhone: false, performanceMeasured: false, fullAcceptanceScenariosCompleted: 0,
  finding: { id: 'NATIVE-UX-01', status: 'FIXED_SOURCE_HOST_PENDING_NATIVE_RECHECK' },
  limitations: ['Synthetic safe-area/text-scale simulations, not actual Android accessibility or cutouts',
    'PNG existence/hashes do not automatically judge visual quality; no independent human sign-off',
    'Initial demo content still requires scrolling; production usability/comprehension is not validated',
    'No new phone installation or permission, no minimum-device or release/profile speed measurements',
    'WB-02/WB-05, all 28 full scenarios and G7/G8/G9 remain open'],
};
const target = join(artifacts, 'layout/validation-receipt.json');
writeFileSync(target, JSON.stringify(report, null, 2) + '\n');
console.log(JSON.stringify({ status: report.status, checks: checks.length,
  failed: checks.filter(item => item.result !== 'PASS'), receipt: relative(root, target) }, null, 2));
process.exitCode = failed ? 1 : 0;
