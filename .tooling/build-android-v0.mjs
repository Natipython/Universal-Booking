// Local synthetic APK only. This runner never accepts licenses or installs onto a device.
import { createHash } from 'node:crypto';
import { existsSync, mkdirSync, readFileSync, readdirSync, statSync, writeFileSync } from 'node:fs';
import { dirname, join, relative, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { spawn } from 'node:child_process';
import { synchronizeDemoAssets } from './sync-android-assets-v0.mjs';
import { verifyBundledDemoAssets } from './verify-android-assets-v0.mjs';
import { inspectApkBadging } from './inspect-android-badging-v0.mjs';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const project = join(root, 'prototypes/android_business');
const artifacts = join(project, 'artifacts/android-build');
const pins = JSON.parse(readFileSync(join(root, '.tooling/android-toolchain-pins.json')));
const sdk = join(root, '.tooling/android-sdk');
const jdk = join(root, '.tooling/jdk-17');
const flutter = join(root, '.tooling/flutter/bin/flutter');
const env = { ...process.env, JAVA_HOME: jdk, ANDROID_HOME: sdk, ANDROID_SDK_ROOT: sdk,
  ANDROID_USER_HOME: join(root, '.tooling/android-user-home'),
  GRADLE_USER_HOME: join(root, '.tooling/gradle-home'),
  TMPDIR: join(root, '.tooling/android-downloads/tmp'),
  PATH: `${join(jdk, 'bin')}:${join(sdk, 'platform-tools')}:${process.env.PATH}` };
mkdirSync(artifacts, { recursive: true });
const sha = path => createHash('sha256').update(readFileSync(path)).digest('hex');
const excluded = new Set(['.gradle', '.kotlin', 'build', 'local.properties']);
const walk = path => readdirSync(path).sort().flatMap(name => excluded.has(name) ? [] :
  statSync(join(path, name)).isDirectory() ? walk(join(path, name)) : [join(path, name)]);
const sourcePaths = () => ['lib', 'android', 'assets'].flatMap(p => walk(join(project, p)))
  .concat(['pubspec.yaml', 'pubspec.lock', '.flutter-plugins-dependencies'].map(p => join(project, p)))
  .concat(['prototypes/fixtures/booking-v0.json', 'design/tokens.json',
    '.tooling/android-toolchain-pins.json', '.tooling/gradle-home/gradle.properties',
    '.tooling/build-android-v0.mjs', '.tooling/verify-android-setup-v0.mjs',
    '.tooling/sync-android-assets-v0.mjs', '.tooling/verify-android-assets-v0.mjs',
    '.tooling/inspect-android-badging-v0.mjs'].map(p => join(root, p))).sort();
const snapshot = () => sourcePaths().map(path => ({ path: relative(root, path), sha256: sha(path) }));
const digest = sources => createHash('sha256').update(JSON.stringify(sources)).digest('hex');
const steps = [];
async function run(name, command, args, cwd = project, stream = true) {
  process.stdout.write(`${name}: running\n`);
  const started = Date.now();
  const child = spawn(command, args, { cwd, env, stdio: ['ignore', 'pipe', 'pipe'] });
  let stdout = '', stderr = '';
  child.stdout.on('data', data => { stdout += data; if (stream) process.stdout.write(data); });
  child.stderr.on('data', data => { stderr += data; if (stream) process.stderr.write(data); });
  const exitCode = await new Promise((ok, fail) => { child.on('error', fail); child.on('close', ok); });
  writeFileSync(join(artifacts, `${name}.txt`), stdout + stderr);
  steps.push({ name, command: relative(root, command), args, exitCode, elapsedMs: Date.now() - started });
  if (exitCode !== 0) throw Error(`${name} failed with exit ${exitCode}; see saved log`);
  return stdout + stderr;
}
let sources = [], initialDigest = null, apk = null, failure = null, manifest = null, signing = null;
let versions = null, sdkPackages = null, sqlitePluginIncluded = false, manifestSummary = null;
let generatedDemoAssets = null, bundledDemoAssets = null;
try {
  if (pins.sdkAuthorization.stablePackageAgreement.status !== 'EXPLICITLY_AUTHORIZED') {
    throw Error('Stable SDK package agreement requires explicit human authorization before build');
  }
  for (const [archive, expected] of [[pins.jdk.archive, pins.jdk.sha256],
    [pins.commandLineTools.archive, pins.commandLineTools.sha256]]) {
    if (sha(join(root, '.tooling/android-downloads', archive)) !== expected) throw Error(`Checksum mismatch: ${archive}`);
  }
  if (sha(join(project, 'android/gradle/wrapper/gradle-wrapper.jar')) !== pins.gradle.wrapperSha256) throw Error('Wrapper checksum mismatch');
  if (!readFileSync(join(project, 'android/gradle/wrapper/gradle-wrapper.properties'), 'utf8').includes(`distributionSha256Sum=${pins.gradle.distributionSha256}`)) throw Error('Distribution checksum missing');
  for (const path of ['platform-tools', 'platforms/android-36', 'build-tools/36.0.0', 'ndk/28.2.13676358']) {
    if (!existsSync(join(sdk, path, 'source.properties'))) throw Error(`Required SDK package absent: ${path}`);
  }
  await run('tool-setup', process.execPath, [join(root, '.tooling/verify-android-setup-v0.mjs')], root);
  if (JSON.parse(readFileSync(join(artifacts, 'setup-receipt.json'))).status !== 'PASS_REQUIRED_TOOL_SETUP') {
    throw Error('Required tool/package setup checks did not pass');
  }
  versions = await run('java-version', join(jdk, 'bin/java'), ['-version']);
  await run('gradle-version', join(project, 'android/gradlew'), ['--version'], join(project, 'android'));
  await run('flutter-version', flutter, ['--suppress-analytics', '--no-version-check', '--version', '--machine']);
  sdkPackages = await run('sdk-packages', join(sdk, 'cmdline-tools/latest/bin/sdkmanager'), [`--sdk_root=${sdk}`, '--list_installed']);
  generatedDemoAssets = synchronizeDemoAssets({ root });
  await run('locked-dependencies', flutter, ['--suppress-analytics', '--no-version-check', 'pub', 'get', '--enforce-lockfile']);
  sources = snapshot(); initialDigest = digest(sources);
  await run('build-debug', flutter, ['--suppress-analytics', '--no-version-check', 'build', 'apk', '--debug', '--no-pub']);
  const path = join(project, 'build/app/outputs/flutter-apk/app-debug.apk');
  apk = { path: relative(root, path), bytes: statSync(path).size, sha256: sha(path) };
  const bt = join(sdk, 'build-tools/36.0.0');
  // The supplied apkanalyzer shell launcher loses quoting for toolsdir when
  // SDK paths contain spaces. Invoke its unchanged official Java entry point
  // with an argument array rather than modifying checksum-verified SDK files.
  const analyzerPrefix = [`-Dcom.android.sdklib.toolsdir=${join(sdk, 'cmdline-tools/latest')}`,
    '-classpath', join(sdk, 'cmdline-tools/latest/lib/apkanalyzer-classpath.jar'),
    'com.android.tools.apk.analyzer.ApkAnalyzerCli'];
  manifest = await run('apk-badging', join(bt, 'aapt2'), ['dump', 'badging', path], project, false);
  const xml = await run('apk-manifest', join(jdk, 'bin/java'), [...analyzerPrefix, 'manifest', 'print', path], project, false);
  signing = await run('apk-signature', join(bt, 'apksigner'), ['verify', '--verbose', '--print-certs', path]);
  const dex = await run('apk-dex', join(jdk, 'bin/java'), [...analyzerPrefix, 'dex', 'packages', '--defined-only', path], project, false);
  sqlitePluginIncluded = dex.includes('com.tekartik.sqflite');
  manifestSummary = inspectApkBadging(manifest);
  if (manifestSummary.packageName !== 'dev.ubi.validation.ubi_business_v0' || manifestSummary.minSdk !== 30 ||
    manifestSummary.targetSdk !== 36 || !sqlitePluginIncluded ||
    !/android:allowBackup\s*=\s*"false"/.test(xml)) throw Error('APK identity/API floor/backup/native-plugin inspection failed');
  manifestSummary.allowBackup = false;
  bundledDemoAssets = verifyBundledDemoAssets({ apkPath: path, root });
  if (digest(snapshot()) !== initialDigest) throw Error('Covered build sources changed during build');
} catch (error) { failure = error.message; process.stderr.write(`${failure}\n`); }
const receipt = {
  checkedAt: new Date().toISOString(), status: failure ? 'FAIL' : 'PASS_LOCAL_DEBUG_APK_BUILD',
  scope: 'Synthetic Ubi Business V0 debug APK compilation and static artifact inspection only',
  sourceDigest: initialDigest, sources, steps, pins, javaVersion: versions, sdkPackages,
  apk, manifestSummary, generatedDemoAssets, bundledDemoAssets,
  signaturePublicIdentity: signing, nativeSqlitePluginIncluded: sqlitePluginIncluded, nativeSqliteExecuted: false,
  failure, productionSigned: false, deviceInstalled: false, performanceMeasured: false,
  notVerified: ['Physical Android 11+/4 GB or newer device behavior', 'Native SQLite execution, Android process death and background/low-storage recovery',
    'Startup, frames, battery, accessibility, localization and field comprehension', 'Real identity, Android/server integration, payments or social connectors',
    'All 28 full acceptance scenarios and G7/G8/G9', 'Production signing, deployment or launch approval']
};
writeFileSync(join(artifacts, 'build-receipt.json'), JSON.stringify(receipt, null, 2) + '\n');
process.stdout.write(JSON.stringify({ status: receipt.status, apk, failure, receipt: relative(root, join(artifacts, 'build-receipt.json')) }, null, 2) + '\n');
process.exitCode = failure ? 1 : 0;
