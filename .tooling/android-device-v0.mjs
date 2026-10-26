// Bounded, synthetic Ubi-only USB smoke-test helper. Never installs/uninstalls,
// clears data, changes phone settings, or reads another app. A single Ubi-only
// force-stop needs the separately recorded explicit session authorization.
// Session-specific model/transport must be rechecked if USB reconnects.
import { spawnSync } from 'node:child_process';
import { createHash } from 'node:crypto';
import { existsSync, mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { DatabaseSync } from 'node:sqlite';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const adb = join(root, '.tooling/android-sdk/platform-tools/adb');
const pkg = 'dev.ubi.validation.ubi_business_v0';
const model = 'SM-G975W';
const transport = '1';
const out = join(root, 'prototypes/android_business/artifacts/android-native/2026-10-15-sm-g975w');
const [action, ...args] = process.argv.slice(2);
const allowed = new Set(['info', 'launch', 'snapshot', 'tap', 'doubletap', 'swipe', 'back', 'database', 'logs', 'force-stop-once', 'verify-install']);
if (!allowed.has(action)) throw Error('Supported: info, launch, verify-install, snapshot LABEL, tap/doubletap X Y, swipe X1 Y1 X2 Y2 MS, back, database LABEL, logs LABEL, force-stop-once');
const authorization = JSON.parse(readFileSync(join(out, 'session-authorization.json')));
if (authorization.sessionStatus !== 'ACTIVE') {
  throw Error('Device test session is closed. Renewed human authorization and target checks are required before any ADB access.');
}

function run(command, { binary = false, optional = false, timeout = 20000 } = {}) {
  const result = spawnSync(adb, ['-P', '5038', '-t', transport, ...command], {
    encoding: binary ? null : 'utf8', timeout, maxBuffer: 24 * 1024 * 1024,
  });
  if ((result.status !== 0 || result.error) && !optional) {
    throw Error(`ADB ${command.slice(0, 4).join(' ')} failed: ${result.error?.message ?? result.stderr}`);
  }
  return binary ? result.stdout : (result.stdout ?? '').trim();
}
if (run(['shell', 'getprop', 'ro.product.model']) !== model) throw Error('Unexpected model/transport. No action taken.');
const build = JSON.parse(readFileSync(join(root, 'prototypes/android_business/artifacts/android-build/build-receipt.json')));
const apkHash = createHash('sha256').update(readFileSync(join(root, build.apk.path))).digest('hex');
if (build.status !== 'PASS_LOCAL_DEBUG_APK_BUILD' || apkHash !== build.apk.sha256) throw Error('Unverified/changed APK');
if (!run(['shell', 'pm', 'path', pkg], { optional: true }).startsWith('package:')) throw Error('Ubi package not installed');
mkdirSync(out, { recursive: true });
const eventsPath = join(out, 'session-events.json');
function record(details) {
  const events = existsSync(eventsPath) ? JSON.parse(readFileSync(eventsPath)) : [];
  events.push({ checkedAt: new Date().toISOString(), action, package: pkg, model, apkSha256: apkHash, ...details });
  writeFileSync(eventsPath, JSON.stringify(events, null, 2) + '\n');
}
function focusedUbi() {
  // Samsung's windows-only dump omits the global mCurrentFocus field.
  const windows = run(['shell', 'dumpsys', 'window']);
  const focus = windows.split('\n').find(line => line.includes('mCurrentFocus='));
  if (!focus?.includes(`${pkg}/`)) throw Error('Ubi is not the focused window. No screenshot/input was taken.');
  return true; // Do not persist unrelated window/process information.
}
function label() {
  if (!/^[a-z0-9-]{1,64}$/.test(args[0] ?? '')) throw Error('Use a short lowercase evidence label');
  return args[0];
}
function save(name, bytes) {
  const path = join(out, name);
  if (existsSync(path)) throw Error(`Evidence already exists; choose a new label: ${name}`);
  writeFileSync(path, bytes);
  return path;
}
let result;
if (action === 'info') {
  result = {
    manufacturer: run(['shell', 'getprop', 'ro.product.manufacturer']), model,
    androidVersion: run(['shell', 'getprop', 'ro.build.version.release']),
    api: Number(run(['shell', 'getprop', 'ro.build.version.sdk'])),
    abis: run(['shell', 'getprop', 'ro.product.cpu.abilist']),
    memory: run(['shell', 'cat', '/proc/meminfo']).split('\n').filter(line => /^MemTotal:|^MemAvailable:/.test(line)),
    dataStorage: run(['shell', 'df', '-k', '/data']),
    screenSize: run(['shell', 'wm', 'size']), density: run(['shell', 'wm', 'density']),
    ramNote: 'MemTotal is OS-usable memory, not an independently verified physical-RAM specification',
    serialRetained: false, permissionsAutoGranted: false,
  };
  writeFileSync(join(out, 'device-info.json'), JSON.stringify({ checkedAt: new Date().toISOString(), ...result }, null, 2) + '\n');
} else if (action === 'verify-install') {
  const packagePath = run(['shell', 'pm', 'path', pkg]).replace(/^package:/, '');
  if (!/^\/data\/app\/[A-Za-z0-9_~.+/=-]+\/base\.apk$/.test(packagePath) || !packagePath.includes(pkg)) {
    throw Error('Unexpected installed APK location');
  }
  const installedSha256 = run(['shell', 'sha256sum', packagePath]).split(/\s+/)[0];
  if (installedSha256 !== apkHash) throw Error('Installed APK hash differs from verified artifact');
  const metadata = run(['shell', 'dumpsys', 'package', pkg]);
  result = { installedApkSha256: installedSha256, matchesLocalVerifiedArtifact: true,
    packageMetadata: metadata.split('\n').map(line => line.trim()).filter(line =>
      /^versionCode=|^versionName=|^pkgFlags=/.test(line)),
    installWasFresh: true, autoGrantedPermissions: false };
} else if (action === 'launch') {
  result = { output: run(['shell', 'am', 'start', '-W', '-n', `${pkg}/${pkg}.MainActivity`]),
    note: 'Launch diagnostic only; not a release/profile startup benchmark' };
} else if (action === 'snapshot') {
  const tag = label(); focusedUbi();
  // Capture only with Ubi focused. Never open the notification shade or another
  // app; system status/navigation bars remain part of this screen capture.
  const png = run(['exec-out', 'screencap', '-p'], { binary: true });
  focusedUbi();
  const path = save(`${tag}-screen.png`, png);
  result = { screenshot: path, focusedPackageCheckedBeforeAndAfter: true,
    note: 'Full app-focused screen capture; status/navigation bars may contain generic system indicators' };
} else if (action === 'tap' || action === 'doubletap' || action === 'swipe') {
  focusedUbi();
  const count = action === 'swipe' ? 5 : 2;
  if (args.length !== count || args.some(value => !/^\d{1,4}$/.test(value))) throw Error('Invalid coordinates');
  result = { coordinates: args, output: action === 'doubletap'
    ? run(['shell', `input tap ${args.join(' ')}; input tap ${args.join(' ')}`])
    : run(['shell', 'input', action, ...args]),
    note: action === 'doubletap' ? 'Two consecutive device input commands; not a measured inter-tap interval or retry stress test' : undefined };
} else if (action === 'back') {
  focusedUbi(); result = { output: run(['shell', 'input', 'keyevent', 'KEYCODE_BACK']) };
} else if (action === 'force-stop-once') {
  const permission = authorization;
  const events = existsSync(eventsPath) ? JSON.parse(readFileSync(eventsPath)) : [];
  if (permission.forceStopAuthorization !== 'Yes, force-stop Ubi only' ||
      permission.forceStopMaximumCount !== 1 || events.some(event => event.action === 'force-stop-once')) {
    throw Error('No unused one-time Ubi-only force-stop authorization');
  }
  focusedUbi();
  const beforePid = run(['shell', 'pidof', pkg]);
  if (!/^\d+$/.test(beforePid)) throw Error('Expected a single running Ubi process');
  const output = run(['shell', 'am', 'force-stop', pkg]);
  const afterPid = run(['shell', 'pidof', pkg], { optional: true });
  result = { output, beforePid, afterPid: afterPid || null, processStopped: afterPid === '',
    dataCleared: false, note: 'Explicitly authorized once. Force-stop is not a low-memory kill, reboot or power-loss test.' };
} else if (action === 'database') {
  const tag = label(); focusedUbi();
  // Read only this newly installed synthetic app's known database; no shared
  // storage, external files, contacts, messages or other application data.
  const bytes = run(['exec-out', 'run-as', pkg, 'cat', 'databases/ubi_synthetic_v0.sqlite'], { binary: true });
  const path = save(`${tag}-synthetic.sqlite`, bytes);
  const db = new DatabaseSync(path, { readOnly: true });
  try {
    result = { database: path, sha256: createHash('sha256').update(bytes).digest('hex'),
      quickCheck: db.prepare('PRAGMA quick_check').all(),
      bookings: db.prepare('SELECT organization_id,operation_id,id,service_id,customer_label,source,start_ms,end_ms,price_minor,status,revision FROM bookings ORDER BY organization_id,operation_id,start_ms').all(),
      outbox: db.prepare('SELECT organization_id,operation_id,command_id,booking_id,state FROM outbox ORDER BY command_id').all(),
      audit: db.prepare('SELECT event_id,organization_id,operation_id,booking_id,action FROM local_audit ORDER BY event_id').all(),
      drafts: db.prepare('SELECT organization_id,operation_id,payload FROM drafts ORDER BY organization_id,operation_id').all(),
      note: 'Snapshot while app idle; no power-loss or concurrent-write durability guarantee',
    };
  } finally { db.close(); }
  save(`${tag}-database.json`, JSON.stringify(result, null, 2) + '\n');
} else if (action === 'logs') {
  const tag = label();
  const pid = run(['shell', 'pidof', pkg]);
  if (!/^\d+$/.test(pid)) throw Error('Expected a single running Ubi process');
  const logs = run(['logcat', '-d', `--pid=${pid}`, '-t', '250', '-v', 'brief'])
    .replace(/https?:\/\/(?:127\.0\.0\.1|localhost)[^\s]+/g, '[LOCAL_DEBUG_ADDRESS_REDACTED]');
  result = { path: save(`${tag}-ubi-log.txt`, logs + '\n'), scope: 'Only the current Ubi PID; no global log dump',
    suspectedErrors: logs.split('\n').filter(line => /FATAL EXCEPTION|E\/flutter|SQLiteException|AndroidRuntime/.test(line)) };
}
record(result);
console.log(JSON.stringify(result, null, 2));
