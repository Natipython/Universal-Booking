import { createHash } from 'node:crypto';
import { mkdirSync, readFileSync, readdirSync, statSync, writeFileSync } from 'node:fs';
import { dirname, join, relative, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { createRequire } from 'node:module';
import { spawn } from 'node:child_process';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const project = join(root, 'prototypes/customer_web');
const artifacts = join(project, 'artifacts');
const go = join(root, '.tooling/go/bin/go');
const postman = join(root, '.tooling/postman-runtime/node_modules/.bin/postman');
const repeatCount = 3;
mkdirSync(artifacts, { recursive: true });
const walk = dir => readdirSync(dir).filter(name => name !== 'artifacts').flatMap(name => {
  const path = join(dir, name); return statSync(path).isDirectory() ? walk(path) : [path];
});
const sources = walk(project).concat([
  '.tooling/build-customer-contract.mjs', '.tooling/test-customer-browser.mjs', '.tooling/validate-customer-v0.mjs',
  '.tooling/postman-runtime/package-lock.json', '.tooling/api-validation/package-lock.json',
  'prototypes/fixtures/booking-v0.json', 'prototypes/booking_integrity/schema.sql',
  'prototypes/booking_integrity/go.mod', 'prototypes/booking_integrity/go.sum',
  'prototypes/booking_integrity/model.go', 'prototypes/booking_integrity/service.go', 'prototypes/booking_integrity/http.go',
].map(path => join(root, path))).sort();
const digest = () => createHash('sha256').update(JSON.stringify(sources.map(path => [relative(root, path), createHash('sha256').update(readFileSync(path)).digest('hex')]))).digest('hex');
const initialDigest = digest();
async function run(executable, args, extraEnv = {}) {
  console.log(`Running ${relative(root, executable)} ${args.join(' ')}`);
  const child = spawn(executable, args, { cwd: project, env: { ...process.env, ...extraEnv }, stdio: ['ignore', 'pipe', 'pipe'] });
  let output = '', stderr = '';
  child.stdout.on('data', data => { output += data; });
  child.stderr.on('data', data => { stderr += data; });
  const exitCode = await new Promise((ok, fail) => { child.on('error', fail); child.on('close', ok); });
  return { args, exitCode, output, stderr };
}
const save = (name, value) => writeFileSync(join(artifacts, name), typeof value === 'string' ? value : JSON.stringify(value, null, 2) + '\n');
const cliJSON = output => { try { return JSON.parse(output.slice(output.indexOf('{'))); } catch { return null; } };
const [goVersion, pgVersion, postmanVersion, format, vet, modules] = await Promise.all([
  run(go, ['version']), run('/usr/lib/postgresql/18/bin/postgres', ['--version']), run(postman, ['--version']),
  run(join(root, '.tooling/go/bin/gofmt'), ['-l', 'app.go', 'database.go', 'app_test.go', 'cmd/demo/main.go']),
  run(go, ['vet', './...']), run(go, ['mod', 'verify']),
]);
save('static-checks.json', { format, vet, modules });
let independentValidation;
try {
  const require = createRequire(join(root, '.tooling/api-validation/package.json'));
  const checked = await require('@apidevtools/swagger-parser').validate(join(project, 'openapi.json'));
  independentValidation = { status: 'PASS', validator: '@apidevtools/swagger-parser', openapi: checked.openapi };
} catch (error) { independentValidation = { status: 'FAIL', error: String(error) }; }
save('independent-openapi-validation.json', independentValidation);
const specLint = await run(postman, ['spec', 'lint', 'openapi.json', '--no-report-events', '-o', 'json']);
save('postman-spec-lint.stdout.txt', specLint.output); save('postman-spec-lint.stderr.txt', specLint.stderr);
const collectionLint = await run(postman, ['collection', 'lint', 'postman/collections/ubi-customer-v0', '--reporter', 'json']);
save('postman-collection-lint.stdout.txt', collectionLint.output); save('postman-collection-lint.stderr.txt', collectionLint.stderr);
const tests = await run(go, ['test', '-race', '-json', `-count=${repeatCount}`, '-timeout=180s', './...'], { UBI_RUN_SYNTHETIC_WEB: '1' });
save('go-test-events.jsonl', tests.output); save('go-test-stderr.txt', tests.stderr);
const events = tests.output.split('\n').flatMap(line => { try { return [JSON.parse(line)]; } catch { return []; } });
const cases = events.filter(e => e.Test && ['pass', 'fail', 'skip'].includes(e.Action)).map(e => ({ name: e.Test, result: e.Action }));
const names = [...new Set(cases.map(e => e.name))];
const postmanOutput = events.filter(e => e.Test === 'TestSavedPostmanCustomerContract' && e.Action === 'output').map(e => e.Output ?? '').join('');
const metric = label => [...postmanOutput.matchAll(new RegExp(`\\|\\s*${label}\\s*\\|\\s*(\\d+)\\s*\\|\\s*(\\d+)\\s*\\|`, 'g'))].map(m => ({ executed: Number(m[1]), failed: Number(m[2]) }));
const requests = metric('requests'), assertions = metric('assertions'), guards = metric('prerequest-scripts');
const readReceipt = name => { try { return JSON.parse(readFileSync(join(artifacts, name), 'utf8')); } catch { return null; } };
const cluster = readReceipt('cluster-receipt.json'), browser = readReceipt('browser-receipt.json'), payload = readReceipt('payload-receipt.json');
const collectionDiagnostics = cliJSON(collectionLint.output), specDiagnostics = cliJSON(specLint.output);
const finalDigest = digest();
const expectedTopLevel = (readFileSync(join(project, 'app_test.go'), 'utf8').match(/^func Test(?!Main\b)\w+\(t \*testing\.T\)/gm) ?? []).length;
const testsPassed = tests.exitCode === 0 && events.some(e => e.Action === 'pass' && !e.Test) &&
  names.filter(name => !name.includes('/')).length === expectedTopLevel && expectedTopLevel >= 21 &&
  cases.every(e => e.result === 'pass') && names.every(name => cases.filter(e => e.name === name).length === repeatCount) &&
  requests.length === repeatCount && requests.every(m => m.executed === 21 && m.failed === 0) &&
  assertions.length === repeatCount && assertions.every(m => m.executed >= 100 && m.failed === 0) &&
  guards.length === repeatCount && guards.every(m => m.executed === 21 && m.failed === 0);
const passed = [goVersion, pgVersion, postmanVersion, format, vet, modules].every(r => r.exitCode === 0) &&
  format.output.trim() === '' && independentValidation.status === 'PASS' && independentValidation.openapi === '3.1.0' &&
  specLint.exitCode === 0 && specDiagnostics?.violations?.length === 0 &&
  collectionLint.exitCode === 0 && collectionDiagnostics?.errorCount === 0 && collectionDiagnostics?.warningCount === 0 && collectionDiagnostics?.filesScanned === 21 &&
  testsPassed && initialDigest === finalDigest && cluster?.server_stopped === true && cluster?.temporary_files_removed === true &&
  cluster?.existing_database_accepted === false && cluster?.tcp_enabled === false &&
  browser?.status === 'PASS_LOCAL_BROWSER_COMPONENT_CHECKS' && browser?.javaScriptEnabled === false && browser?.outsideRequests?.length === 0 &&
  payload?.javascript_bytes === 0 && Object.values(payload?.bytes ?? {}).length === 6 && Object.values(payload?.bytes ?? {}).every(bytes => bytes < 16384);
const receipt = {
  status: passed ? 'PASS_LOCAL_COMPONENT_CHECKS' : 'FAIL_OR_INCOMPLETE', checkedAt: new Date().toISOString(),
  scope: 'WB-04 local synthetic component only; not acceptance-scenario or release approval',
  sourceDigest: finalDigest, sourceUnchangedDuringValidation: initialDigest === finalDigest,
  digestScope: sources.map(path => relative(root, path)),
  tools: { go: goVersion.output.trim(), postgres: pgVersion.output.trim(), postman: postmanVersion.output.trim(), chromium: browser?.browser ?? null },
  staticChecksPassed: [format, vet, modules].every(r => r.exitCode === 0) && format.output.trim() === '',
  authenticatedPostmanGovernanceReview: false,
  openapi: independentValidation, postman: { specLintExitCode: specLint.exitCode, specViolations: specDiagnostics?.violations?.length ?? null,
    collectionLintExitCode: collectionLint.exitCode, collectionDiagnostics, requests, assertions, guards },
  tests: { exitCode: tests.exitCode, repeatCount, distinctTopLevel: names.filter(name => !name.includes('/')).length,
    distinctCasesIncludingSubtests: names.length, caseResults: cases.length, passed: cases.filter(e => e.result === 'pass').length,
    failed: cases.filter(e => e.result === 'fail').length, skipped: cases.filter(e => e.result === 'skip').length, cases },
  cluster, browser, payload,
  notVerified: ['All 28 full acceptance scenarios', 'G7/G8/G9', 'Android/server integration or APK', 'Physical Android 11+/4 GB speed/offline behavior',
    'Weak-network N1 or Core Web Vitals', 'Real-user research and localization/accessibility acceptance', 'Persistent production grant/session/recovery storage',
    'Customer identity, public TLS/domain/native deep links', 'Real social inboxes, messages, calls, notifications or money',
    'Production schema migrations, restore, sustained load or independent review', 'Existing WB-03 application/problem+json spec-linter disagreement'],
};
save('validation-receipt.json', receipt);
console.log(JSON.stringify({ status: receipt.status, topLevel: receipt.tests.distinctTopLevel, distinctCases: names.length, caseResults: cases.length,
  failed: receipt.tests.failed, skipped: receipt.tests.skipped, requests, assertions, specLintExitCode: specLint.exitCode, sourceUnchanged: initialDigest === finalDigest }, null, 2));
process.exitCode = passed ? 0 : 1;
