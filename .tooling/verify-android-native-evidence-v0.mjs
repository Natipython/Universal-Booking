// Offline cross-check of this dated, manually driven physical-phone session.
// This script NEVER connects to a phone or changes any Android data/settings.
import assert from 'node:assert/strict';
import { createHash } from 'node:crypto';
import { existsSync, readFileSync, readdirSync, writeFileSync } from 'node:fs';
import { dirname, join, relative, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { DatabaseSync } from 'node:sqlite';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const evidence = join(root, 'prototypes/android_business/artifacts/android-native/2026-10-15-sm-g975w');
const hash = path => createHash('sha256').update(readFileSync(path)).digest('hex');
const json = name => JSON.parse(readFileSync(join(evidence, name)));
// The working UI/build may evolve after this phone session. Validate the frozen
// tested inputs, not whichever APK/current source happens to exist now.
const checkpointPath = 'prototypes/android_business/artifacts/checkpoints/2026-10-15-pre-layout/checkpoint.json';
const checkpoint = JSON.parse(readFileSync(join(root, checkpointPath)));
const archived = new Map(checkpoint.files.map(file => [file.originalPath, file]));
function frozen(originalPath) {
  const file = archived.get(originalPath);
  if (!file || hash(join(root, file.snapshotPath)) !== file.sha256) throw Error(`Frozen input absent/changed: ${originalPath}`);
  return join(root, file.snapshotPath);
}
const build = JSON.parse(readFileSync(frozen(checkpoint.buildReceiptPath)));
const events = json('session-events.json');
const device = json('device-info.json');
const checks = [];
function check(name, verify) {
  try { verify(); checks.push({ name, result: 'PASS' }); }
  catch (error) { checks.push({ name, result: 'FAIL', reason: error.message }); }
}
function snapshot(label) {
  const report = json(`${label}-database.json`);
  const path = join(evidence, `${label}-synthetic.sqlite`);
  assert.equal(hash(path), report.sha256, `Snapshot checksum: ${label}`);
  const db = new DatabaseSync(path, { readOnly: true });
  try {
    assert.equal(db.prepare('PRAGMA quick_check').get().quick_check, 'ok');
    const tables = Object.fromEntries(['bookings', 'outbox', 'local_audit', 'drafts'].map(table =>
      [table, db.prepare(`SELECT * FROM ${table} ORDER BY rowid`).all()]));
    assert.equal(tables.bookings.length, report.bookings.length);
    assert.equal(tables.outbox.length, report.outbox.length);
    assert.equal(tables.local_audit.length, report.audit.length);
    assert.equal(tables.drafts.length, report.drafts.length);
    return { ...tables, sha256: report.sha256 };
  } finally { db.close(); }
}
let fatal = null;
try {
  check('Installed APK and frozen covered application sources match the tested build', () => {
    const authorization = json('session-authorization.json');
    assert.equal(authorization.sessionStatus, 'CLOSED');
    assert.equal(authorization.testingConnectionStopped, true);
    assert.equal(build.status, 'PASS_LOCAL_DEBUG_APK_BUILD');
    assert.equal(checkpoint.status, 'PASS_ARCHIVED_BASELINE');
    assert.equal(hash(frozen(build.apk.path)), build.apk.sha256);
    for (const source of build.sources) assert.equal(hash(frozen(source.path)), source.sha256, source.path);
    const installed = events.find(event => event.action === 'verify-install');
    assert.equal(installed.installedApkSha256, build.apk.sha256);
    assert.equal(installed.matchesLocalVerifiedArtifact, true);
    assert.equal(device.model, 'SM-G975W'); assert.equal(device.api, 31);
  });
  const empty = snapshot('01-empty-store');
  const review = snapshot('03-review-no-booking');
  const saved = snapshot('04-saved-once');
  const conflict = snapshot('06-conflict-retained');
  const revised = snapshot('09-revised-booking');
  const beforeStop = snapshot('11-before-force-stop');
  const afterStop = snapshot('12-after-force-stop');
  const overlap = snapshot('13-no-overlap-duplicate');
  const kept = snapshot('14-kept-item');
  const cancelled = snapshot('15-cancelled-history-retained');
  const otherScope = snapshot('17-events-scope');
  const returned = snapshot('19-returned-scope');
  const final = snapshot('20-new-item-after-cancellation');
  check('Native store is empty on fresh launch; review does not create a booking', () => {
    for (const table of ['bookings', 'outbox', 'local_audit', 'drafts']) assert.equal(empty[table].length, 0);
    assert.equal(empty.sha256, review.sha256);
  });
  check('Two consecutive save taps produce one scoped pending item, command and history event', () => {
    assert.equal(saved.bookings.length, 1); assert.equal(saved.outbox.length, 1); assert.equal(saved.local_audit.length, 1);
    const item = saved.bookings[0];
    assert.equal(item.status, 'pending'); assert.equal(item.organization_id, 'demo-org-a');
    assert.equal(item.operation_id, 'studio'); assert.equal(item.customer_label, 'Demo Guest 01');
    assert.equal(item.price_minor, 200000); assert.equal(item.start_ms, 1894687200000);
    assert.equal(item.end_ms - item.start_ms, 3600000);
  });
  check('Conflict simulation preserves details and appends history', () => {
    assert.equal(conflict.bookings[0].id, saved.bookings[0].id);
    assert.equal(conflict.bookings[0].status, 'simulatedConflict');
    assert.equal(conflict.bookings[0].start_ms, saved.bookings[0].start_ms);
    assert.equal(conflict.local_audit.length, 2);
    assert.equal(conflict.local_audit[1].action, 'simulated_conflict');
  });
  check('Reviewed revision keeps the booking ID, original command payload and prior history', () => {
    const item = revised.bookings[0]; assert.equal(item.id, saved.bookings[0].id);
    assert.equal(item.status, 'pending'); assert.equal(item.revision, 3);
    assert.equal(item.start_ms, saved.bookings[0].start_ms + 3600000);
    assert.equal(revised.outbox.length, 2); assert.equal(revised.local_audit.length, 3);
    assert.equal(revised.outbox.find(row => row.command_id === saved.outbox[0].command_id).payload, saved.outbox[0].payload);
    assert.equal(revised.local_audit[2].action, 'local_intent_revised');
  });
  check('One authorized force-stop actually stops Ubi; all stored bytes and draft survive reopening', () => {
    const stops = events.filter(event => event.action === 'force-stop-once');
    assert.equal(stops.length, 1); assert.equal(stops[0].processStopped, true);
    assert.equal(stops[0].afterPid, null); assert.equal(stops[0].dataCleared, false);
    assert.equal(beforeStop.sha256, afterStop.sha256);
    const draft = JSON.parse(afterStop.drafts[0].payload);
    assert.equal(draft.customerLabel, 'Demo Guest 02'); assert.equal(draft.startUtcMs, 1894690800000);
    assert.ok(existsSync(join(evidence, '12-recovered-draft-screen.png')));
  });
  check('Rejected overlapping save and declined cancellation leave stored data unchanged', () => {
    assert.equal(afterStop.sha256, overlap.sha256); assert.equal(overlap.sha256, kept.sha256);
  });
  check('Confirmed local cancellation retains booking/history and cancels queued commands', () => {
    assert.equal(cancelled.bookings.length, 1); assert.equal(cancelled.bookings[0].status, 'cancelled');
    assert.equal(cancelled.bookings[0].revision, 4); assert.equal(cancelled.local_audit.length, 4);
    assert.equal(cancelled.local_audit[3].action, 'local_intent_cancelled');
    assert.ok(cancelled.outbox.every(row => row.state === 'cancelled'));
  });
  check('Scope switching does not mutate the original draft; UI evidence for all three scopes exists', () => {
    assert.equal(cancelled.sha256, otherScope.sha256); assert.equal(otherScope.sha256, returned.sha256);
    for (const name of ['17-events-no-foreign-draft', '18-other-business-no-foreign-draft', '19-returned-own-draft']) {
      assert.ok(existsSync(join(evidence, `${name}-screen.png`)));
    }
    assert.equal(JSON.parse(returned.drafts[0].payload).customerLabel, 'Demo Guest 02');
  });
  check('Cancelled time can be reused; original history remains and completed draft is cleared', () => {
    assert.equal(final.bookings.length, 2); assert.equal(final.outbox.length, 3); assert.equal(final.local_audit.length, 5);
    const pending = final.bookings.filter(item => item.status === 'pending');
    assert.equal(pending.length, 1); assert.equal(pending[0].customer_label, 'Demo Guest 02');
    assert.equal(final.bookings.filter(item => item.status === 'cancelled').length, 1);
    assert.equal(pending[0].start_ms, cancelled.bookings[0].start_ms); assert.equal(final.drafts.length, 0);
    assert.equal(final.local_audit.filter(row => row.booking_id === cancelled.bookings[0].id).length, 4);
  });
  check('Screenshots are Ubi-focused and retained PID logs omit local debug-service addresses', () => {
    const captures = events.filter(event => event.action === 'snapshot');
    assert.ok(captures.length >= 18); assert.ok(captures.every(event => event.focusedPackageCheckedBeforeAndAfter === true));
    for (const name of ['01-first-launch-ubi-log.txt', '20-final-ubi-log.txt']) {
      const text = readFileSync(join(evidence, name), 'utf8');
      assert.equal(/https?:\/\/(?:127\.0\.0\.1|localhost)[^\s]+/.test(text), false);
      assert.equal(/FATAL EXCEPTION|E\/flutter|SQLiteException/.test(text), false);
    }
  });
} catch (error) { fatal = error.message; }
const evidenceFiles = readdirSync(evidence).filter(name => name !== 'validation-receipt.json').sort()
  .map(name => ({ path: relative(root, join(evidence, name)), sha256: hash(join(evidence, name)) }));
const report = {
  checkedAt: new Date().toISOString(), status: fatal || checks.some(check => check.result !== 'PASS') ? 'FAIL' : 'PASS_BOUNDED_NATIVE_SMOKE',
  scope: 'Offline cross-check of one manually driven Android 12 physical-phone debug session; partial WB-02/EV-03/05 only',
  device, testedApk: { ...build.apk, originalBuildPath: build.apk.path,
    path: relative(root, frozen(build.apk.path)) }, buildSourceDigest: build.sourceDigest,
  testedBuildCheckpoint: checkpointPath,
  workingSourcesStillMatchTestedBuild: build.sources.every(source =>
    existsSync(join(root, source.path)) && hash(join(root, source.path)) === source.sha256),
  verificationNote: 'Offline recheck of frozen historical phone evidence, not a new device run. Current source/build can differ.',
  verificationSources: ['.tooling/android-device-v0.mjs', '.tooling/verify-android-native-evidence-v0.mjs']
    .map(path => ({ path, sha256: hash(join(root, path)) })),
  checks, fatal, evidenceFiles,
  nativeSqliteExecuted: true, installedDeviceCount: 1, performanceMeasured: false, fullAcceptanceScenariosCompleted: 0,
  manualVisualEvidence: 'Agent inspected app-focused screenshots; existence/hashes do not automatically validate their UI semantics',
  findings: [{ id: 'NATIVE-UX-01', status: 'OPEN_NOT_FIXED',
    detail: 'At the initial Home scroll position, the floating Add button overlaps part of the booking card. Scrolling reveals the card; simplify the initial schedule layout and keep controls clear in a separately approved UI change.' }],
  notVerified: ['Android 11/4 GB baseline or other/newer devices', 'Release/profile startup, frames, memory, battery or performance targets',
    'Low-memory kill, reboot/power loss, background restrictions or low-storage failure', 'Airplane-mode/weak-network transitions',
    'TalkBack, OS text scaling, language/calendar localization and participant comprehension',
    'Production identity/permissions, server synchronization, payments, connectors, native messaging/calling', 'All 28 full scenarios and G7/G8/G9'],
};
writeFileSync(join(evidence, 'validation-receipt.json'), JSON.stringify(report, null, 2) + '\n');
console.log(JSON.stringify({ status: report.status, checksPassed: checks.filter(check => check.result === 'PASS').length,
  checksFailed: checks.filter(check => check.result !== 'PASS'), fatal, nativeSqliteExecuted: true,
  receipt: relative(root, join(evidence, 'validation-receipt.json')) }, null, 2));
process.exitCode = report.status === 'FAIL' ? 1 : 0;
