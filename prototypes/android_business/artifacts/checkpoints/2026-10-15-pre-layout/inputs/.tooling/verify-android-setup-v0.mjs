// Read-only tool checks plus generated evidence; no installation, license acceptance or device commands.
import { createHash } from 'node:crypto';
import { existsSync, mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import { dirname, join, relative, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { spawnSync } from 'node:child_process';
const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const pins = JSON.parse(readFileSync(join(root, '.tooling/android-toolchain-pins.json')));
const sdk = join(root, '.tooling/android-sdk'), jdk = join(root, '.tooling/jdk-17');
const target = join(root, 'prototypes/android_business/artifacts/android-build');
mkdirSync(target, { recursive: true });
const sha = path => createHash('sha256').update(readFileSync(path)).digest('hex');
const env = { ...process.env, JAVA_HOME: jdk, ANDROID_HOME: sdk, ANDROID_SDK_ROOT: sdk,
  ANDROID_USER_HOME: join(root, '.tooling/android-user-home'), TMPDIR: join(root, '.tooling/android-downloads/tmp') };
const checkVersion = (name, path, args) => {
  const result = spawnSync(path, args, { env, encoding: 'utf8', timeout: 30000 });
  const output = (result.stdout ?? '') + (result.stderr ?? '');
  writeFileSync(join(target, `${name}.txt`), output);
  return { exitCode: result.status, output, error: result.error?.message ?? null };
};
const java = checkVersion('setup-java-version', join(jdk, 'bin/java'), ['-version']);
const sdkmanager = checkVersion('setup-sdkmanager-version', join(sdk, 'cmdline-tools/latest/bin/sdkmanager'), ['--version']);
const downloads = [pins.jdk, pins.commandLineTools].map(pin => {
  const path = join(root, '.tooling/android-downloads', pin.archive);
  const actual = sha(path);
  return { archive: relative(root, path), source: pin.url, expectedSha256: pin.sha256,
    actualSha256: actual, verified: actual === pin.sha256 };
});
const wrapper = sha(join(root, 'prototypes/android_business/android/gradle/wrapper/gradle-wrapper.jar'));
const requirements = [
  { path: 'platform-tools', expected: {} },
  { path: 'platforms/android-36', expected: { 'AndroidVersion.ApiLevel': '36' } },
  { path: 'build-tools/36.0.0', expected: { 'Pkg.Revision': '36.0.0' } },
  { path: 'ndk/28.2.13676358', expected: { 'Pkg.Revision': '28.2.13676358' } }
];
const packages = requirements.map(({ path, expected }) => {
  const source = join(sdk, path, 'source.properties'), installed = existsSync(source);
  const properties = installed ? Object.fromEntries(readFileSync(source, 'utf8').split('\n')
    .filter(line => line.includes('=') && !line.trimStart().startsWith('#'))
    .map(line => { const at = line.indexOf('='); return [line.slice(0, at).trim(), line.slice(at + 1).trim()]; })) : null;
  return { path, installed, properties, expected,
    versionMatches: installed && Object.entries(expected).every(([key, value]) => properties[key] === value) };
});
const authorized = pins.sdkAuthorization.stablePackageAgreement.status === 'EXPLICITLY_AUTHORIZED';
const checksPass = downloads.every(d => d.verified) && wrapper === pins.gradle.wrapperSha256 &&
  java.exitCode === 0 && java.output.includes('17.0.20.1') && sdkmanager.exitCode === 0 && /\b22\.0\b/.test(sdkmanager.output);
const receipt = {
  checkedAt: new Date().toISOString(),
  status: !checksPass ? 'FAIL_SETUP_CHECKS' : !authorized ? 'PARTIAL_STABLE_PACKAGE_AGREEMENT_PENDING' :
    packages.every(p => p.versionMatches) ? 'PASS_REQUIRED_TOOL_SETUP' :
      packages.every(p => p.installed) ? 'FAIL_PACKAGE_VERSION' : 'PARTIAL_REQUIRED_PACKAGES_PENDING',
  scope: 'Project-local JDK/Android command-line tools, pinned downloads and setup boundary only',
  checksPass, downloads, java, sdkmanager, packages,
  wrapper: { sha256: wrapper, verified: wrapper === pins.gradle.wrapperSha256, version: pins.gradle.version },
  sdkAuthorization: pins.sdkAuthorization,
  sdkManagerLicenseFile: existsSync(join(sdk, 'licenses/android-sdk-license')) ? {
    path: relative(root, join(sdk, 'licenses/android-sdk-license')),
    sha256: sha(join(sdk, 'licenses/android-sdk-license')) } : null,
  stableAgreementAcceptedByThisRunner: false,
  buildAttemptedByThisRunner: false, deviceAccessAttemptedByThisRunner: false,
  apkPresent: existsSync(join(root, 'prototypes/android_business/build/app/outputs/flutter-apk/app-debug.apk')),
  hostReceipt: JSON.parse(readFileSync(join(root, 'prototypes/android_business/artifacts/validation-receipt.json'))),
  notVerified: ['APK compilation and static inspection', 'Any Android native runtime or physical-device performance',
    'Production signing, identity/server integration, complete acceptance scenarios or release gates']
};
writeFileSync(join(target, 'setup-receipt.json'), JSON.stringify(receipt, null, 2) + '\n');
console.log(JSON.stringify({ status: receipt.status, checksPass, java: java.exitCode, sdkmanager: sdkmanager.exitCode,
  packages, apkPresent: receipt.apkPresent, receipt: relative(root, join(target, 'setup-receipt.json')) }, null, 2));
process.exitCode = checksPass && receipt.status !== 'FAIL_PACKAGE_VERSION' ? 0 : 1;
