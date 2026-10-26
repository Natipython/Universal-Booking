import { createHash } from 'node:crypto';
import { mkdirSync, readFileSync, readdirSync, statSync, writeFileSync } from 'node:fs';
import { dirname, join, relative, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { spawn } from 'node:child_process';
import { synchronizeDemoAssets } from './sync-android-assets-v0.mjs';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const project = join(root, 'prototypes/android_business');
const flutter = join(root, '.tooling/flutter/bin/flutter');
const walk = (dir) => readdirSync(dir).flatMap((name) => {
  const path = join(dir, name);
  return statSync(path).isDirectory() ? walk(path) : [path];
});
const generatedDemoAssets = synchronizeDemoAssets({ root });
const sourceFiles = ['lib', 'test', 'tool', 'assets'].flatMap((dir) => walk(join(project, dir)))
  .concat(['pubspec.yaml', 'pubspec.lock', 'analysis_options.yaml',
    'android/app/src/main/AndroidManifest.xml', 'android/app/build.gradle.kts'].map((p) => join(project, p)))
  .concat([fileURLToPath(import.meta.url), join(root, '.tooling/sync-android-assets-v0.mjs')])
  .concat([join(root, 'prototypes/fixtures/booking-v0.json'), join(root, 'design/tokens.json')]).sort();
const digest = () => createHash('sha256').update(JSON.stringify(sourceFiles.map((p) => [
  relative(root, p), createHash('sha256').update(readFileSync(p)).digest('hex'),
]))).digest('hex');
const initialDigest = digest();
async function run(args) {
  const cliArgs = ['--suppress-analytics', '--no-version-check', ...args];
  process.stdout.write(`Running Flutter ${cliArgs.join(' ')}\n`);
  const child = spawn(flutter, cliArgs, { cwd: project, stdio: ['ignore', 'pipe', 'pipe'] });
  let output = '', stderr = '';
  child.stdout.on('data', (data) => { output += data; });
  child.stderr.on('data', (data) => { stderr += data; });
  const exitCode = await new Promise((ok, fail) => {
    child.on('error', fail); child.on('close', ok);
  });
  return { args: cliArgs, exitCode, output, stderr };
}
const version = await run(['--version', '--machine']);
const analyze = await run(['analyze', '--no-pub']);
const tests = await run(['test', '--no-pub', '--machine']);
const events = tests.output.split('\n').flatMap((line) => {
  try { return [JSON.parse(line)]; } catch { return []; }
});
const testNames = new Map(events.filter((e) => e.type === 'testStart').map((e) => [e.test.id, e.test.name]));
const finished = events.filter((e) => e.type === 'testDone' && !e.hidden).map((e) => ({
  name: testNames.get(e.testID), result: e.result, skipped: e.skipped === true,
}));
const completed = events.findLast((e) => e.type === 'done');
let toolchain = null;
try { toolchain = JSON.parse(version.output); } catch { /* preserve failure in receipt */ }
const finalDigest = digest();
const androidFloor = /\bminSdk\s*=\s*30\b/.test(readFileSync(join(project, 'android/app/build.gradle.kts'), 'utf8'));
const passed = version.exitCode === 0 && toolchain !== null && androidFloor && analyze.exitCode === 0 && tests.exitCode === 0 &&
  completed?.success === true && finished.length > 0 &&
  finished.every((t) => t.result === 'success' && !t.skipped) && initialDigest === finalDigest;
const report = {
  checkedAt: new Date().toISOString(), status: passed ? 'PASS' : 'FAIL',
  scope: 'Synthetic V0 local SQLite/domain/UI tests on Linux Flutter host',
  sourceDigest: finalDigest, sourcesChangedDuringRun: initialDigest !== finalDigest, generatedDemoAssets,
  supportedBaseline: { android: '11+', minSdk: 30, ramGb: 4, configurationChecked: androidFloor,
    ramEnforcedAtRuntime: false, physicalDeviceVerified: false },
  toolchain, analyzerExitCode: analyze.exitCode, testExitCode: tests.exitCode,
  commands: [version.args, analyze.args, tests.args],
  testsRun: finished.length, passed: finished.filter((t) => t.result === 'success' && !t.skipped).length,
  tests: finished,
  notVerified: ['APK build or Android native SQLite adapter', 'Physical-device startup/frames/battery/network',
    'Actual Android process-death/background/low-storage behavior', 'Real auth/server sync or Go/PostgreSQL integrity',
    'Field comprehension, TalkBack, Amharic or Ethiopian-calendar entry', 'Legal/provider/production readiness'],
};
const target = join(project, 'artifacts');
mkdirSync(target, { recursive: true });
writeFileSync(join(target, 'validation-receipt.json'), JSON.stringify(report, null, 2) + '\n');
writeFileSync(join(target, 'analyzer.txt'), analyze.output + analyze.stderr);
writeFileSync(join(target, 'test-events.jsonl'), tests.output);
writeFileSync(join(target, 'test-stderr.txt'), tests.stderr);
process.stdout.write(JSON.stringify({status: report.status, testsRun: report.testsRun,
  passed: report.passed, analyzerExitCode: report.analyzerExitCode,
  receipt: relative(root, join(target, 'validation-receipt.json'))}, null, 2) + '\n');
if (!passed) process.stdout.write((analyze.stderr + tests.stderr).slice(-6000));
process.exitCode = passed ? 0 : 1;
