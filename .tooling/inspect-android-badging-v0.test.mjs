import { test } from 'node:test';
import assert from 'node:assert/strict';
import { inspectApkBadging } from './inspect-android-badging-v0.mjs';
const fixture = "package: name='dev.ubi.validation.ubi_business_v0' versionCode='1'\nminSdkVersion:'30'\ntargetSdkVersion:'36'\napplication-debuggable\nuses-permission: name='android.permission.INTERNET'\nnative-code: 'arm64-v8a' 'armeabi-v7a' 'x86_64'\n";
test('actual aapt2 minimum API label parses without a false failure', () => {
  const result = inspectApkBadging(fixture);
  assert.equal(result.packageName, 'dev.ubi.validation.ubi_business_v0');
  assert.equal(result.minSdk, 30);
  assert.equal(result.targetSdk, 36);
  assert.equal(result.debuggable, true);
  assert.deepEqual(result.permissions, ['android.permission.INTERNET']);
});
test('badging parser preserves a changed floor rather than assuming 30', () => {
  assert.equal(inspectApkBadging(fixture.replace("minSdkVersion:'30'", "minSdkVersion:'31'")).minSdk, 31);
});
test('badging parser rejects missing minimum API rather than guessing it', () => {
  assert.throws(() => inspectApkBadging(fixture.replace("minSdkVersion:'30'", "sdkVersion:'30'")), /lacks/);
});
