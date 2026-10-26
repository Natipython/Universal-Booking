import { createHash } from 'node:crypto';
import { mkdirSync, readFileSync, readdirSync, statSync, writeFileSync } from 'node:fs';
import { dirname, join, relative, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { createRequire } from 'node:module';
import { spawn } from 'node:child_process';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const project = join(root, 'prototypes/booking_integrity');
const artifacts = join(project, 'artifacts');
const go = join(root, '.tooling/go/bin/go');
const postman = join(root, '.tooling/postman-runtime/node_modules/.bin/postman');
const pg = '/usr/lib/postgresql/18/bin/postgres';
const repeatCount = 3;
mkdirSync(artifacts, { recursive: true });
const walk = (dir) => readdirSync(dir).flatMap((name) => {
  const path = join(dir, name);
  return statSync(path).isDirectory() ? walk(path) : [path];
});
const sources = ['README.md', 'CONTRACT.md', 'AGENTS.md', 'go.mod', 'go.sum', 'schema.sql',
  'model.go', 'service.go', 'http.go', 'service_test.go', 'contract_test.go', 'openapi.json',
  '.postman/resources.yaml', 'postman/skills/api-builder/SKILL.md'].map((p) => join(project, p))
  .concat(walk(join(project, 'postman/collections/ubi-synthetic-v0')))
  .concat(['.tooling/build-integrity-contract.mjs', '.tooling/build-integrity-collection.mjs',
    '.tooling/validate-integrity-v0.mjs', '.tooling/postman-runtime/package-lock.json',
    '.tooling/api-validation/package-lock.json', 'prototypes/fixtures/booking-v0.json'].map((p) => join(root, p)))
  .sort();
const digest = () => createHash('sha256').update(JSON.stringify(sources.map((p) => [
  relative(root, p), createHash('sha256').update(readFileSync(p)).digest('hex'),
]))).digest('hex');
const initialDigest = digest();
async function run(executable, args, extraEnv = {}) {
  process.stdout.write(`Running ${relative(root, executable)} ${args.join(' ')}\n`);
  const child = spawn(executable, args, { cwd: project, env: { ...process.env, ...extraEnv },
    stdio: ['ignore', 'pipe', 'pipe'] });
  let output = '', stderr = '';
  child.stdout.on('data', (data) => { output += data; });
  child.stderr.on('data', (data) => { stderr += data; });
  const exitCode = await new Promise((ok, fail) => { child.on('error', fail); child.on('close', ok); });
  return { args, exitCode, output, stderr };
}
function save(name, value) {
  writeFileSync(join(artifacts, name), typeof value === 'string' ? value : JSON.stringify(value, null, 2) + '\n');
}
function parseCLIJSON(output) {
  try { return JSON.parse(output.slice(output.indexOf('{'))); } catch { return null; }
}
const [goVersion, pgVersion, postmanVersion, format, vet, modules] = await Promise.all([
  run(go, ['version']), run(pg, ['--version']), run(postman, ['--version']),
  run(join(root, '.tooling/go/bin/gofmt'), ['-l', 'model.go', 'service.go', 'http.go', 'service_test.go', 'contract_test.go']),
  run(go, ['vet', './...']), run(go, ['mod', 'verify']),
]);
save('static-checks.json', { format, vet, modules });
const spec = JSON.parse(readFileSync(join(project, 'openapi.json'), 'utf8'));
let independentValidation;
try {
  const require = createRequire(join(root, '.tooling/api-validation/package.json'));
  const validator = require('@apidevtools/swagger-parser');
  const checked = await validator.validate(join(project, 'openapi.json'));
  independentValidation = { status: 'PASS', validator: '@apidevtools/swagger-parser', openapi: checked.openapi };
} catch (error) {
  independentValidation = { status: 'FAIL', error: String(error) };
}
save('independent-openapi-validation.json', independentValidation);
const specLint = await run(postman, ['spec', 'lint', 'openapi.json', '--no-report-events', '-o', 'json']);
save('postman-spec-lint.stdout.txt', specLint.output);
save('postman-spec-lint.stderr.txt', specLint.stderr);
const diagnostics = parseCLIJSON(specLint.output);
const knownMediaRangeOnly = specLint.exitCode === 1 && diagnostics?.violations?.length > 0 &&
  diagnostics.violations.every((v) => {
    const path = v.path?.match(/^paths\.(.*?)\.(get|post)\.responses\.(\d{3})\.content$/);
    const content = path && spec.paths[path[1]]?.[path[2]]?.responses[path[3]]?.content;
    return v.severity === 'ERROR' && ['must match format "media-range"', 'property name must be valid'].includes(v.issue) &&
      content && Object.keys(content).length === 1 && Object.hasOwn(content, 'application/problem+json');
  });
const collectionLint = await run(postman, ['collection', 'lint', 'postman/collections/ubi-synthetic-v0', '--reporter', 'json']);
save('postman-collection-lint.stdout.txt', collectionLint.output);
save('postman-collection-lint.stderr.txt', collectionLint.stderr);
const collectionDiagnostics = parseCLIJSON(collectionLint.output);
const tests = await run(go, ['test', '-race', '-json', `-count=${repeatCount}`, '-timeout=180s', './...'],
  { UBI_RUN_SYNTHETIC_PG: '1' });
save('go-test-events.jsonl', tests.output);
save('go-test-stderr.txt', tests.stderr);
const events = tests.output.split('\n').flatMap((line) => {
  try { return [JSON.parse(line)]; } catch { return []; }
});
const cases = events.filter((e) => e.Test && ['pass', 'fail', 'skip'].includes(e.Action))
  .map((e) => ({ name: e.Test, result: e.Action }));
const passedCases = cases.filter((e) => e.result === 'pass');
const suitePass = events.some((e) => e.Action === 'pass' && !e.Test);
const collectionOutput = events.filter((e) => e.Test === 'TestPostmanV3ContractCollection' && e.Action === 'output')
  .map((e) => e.Output ?? '').join('');
const metric = (label) => [...collectionOutput.matchAll(new RegExp(`\\|\\s*${label}\\s*\\|\\s*(\\d+)\\s*\\|\\s*(\\d+)\\s*\\|`, 'g'))]
  .map((m) => ({ executed: Number(m[1]), failed: Number(m[2]) }));
const requests = metric('requests'), assertions = metric('assertions'), guards = metric('prerequest-scripts');
let cluster = null;
try { cluster = JSON.parse(readFileSync(join(artifacts, 'cluster-receipt.json'), 'utf8')); } catch { /* fail closed */ }
const finalDigest = digest();
const testsPassed = tests.exitCode === 0 && suitePass && passedCases.length === 25 * repeatCount &&
  cases.every((e) => e.result === 'pass') && new Set(passedCases.map((e) => e.name)).size === 25 &&
  requests.length === repeatCount && requests.every((m) => m.executed === 25 && m.failed === 0) &&
  assertions.length === repeatCount && assertions.every((m) => m.executed === 172 && m.failed === 0) &&
  guards.length === repeatCount && guards.every((m) => m.executed === 25 && m.failed === 0);
const coreChecksPassed = [goVersion, pgVersion, postmanVersion, format, vet, modules].every((r) => r.exitCode === 0) &&
  format.output.trim() === '' && independentValidation.status === 'PASS' && independentValidation.openapi === '3.1.0' &&
  collectionLint.exitCode === 0 && collectionDiagnostics?.errorCount === 0 && collectionDiagnostics?.warningCount === 0 &&
  collectionDiagnostics?.filesScanned === 25 && testsPassed && initialDigest === finalDigest &&
  cluster?.synthetic_only === true && cluster?.new_private_temporary_cluster === true &&
  cluster?.existing_database_accepted === false && cluster?.tcp_enabled === false &&
  cluster?.server_stopped === true && cluster?.temporary_files_removed === true;
const boundedResult = coreChecksPassed && (specLint.exitCode === 0 || knownMediaRangeOnly);
const report = {
  checkedAt: new Date().toISOString(),
  status: boundedResult ? (specLint.exitCode === 0 ? 'PASS_LOCAL_CHECKS' : 'PARTIAL_TOOLING_LIMITATION') : 'FAIL',
  coreChecksPassed, allRequestedLintersPassed: specLint.exitCode === 0 && collectionLint.exitCode === 0,
  scope: 'WB-03 synthetic loopback Go/PostgreSQL transaction and HTTP component checks on Linux; not production',
  sourceDigest: finalDigest, sourcesChangedDuringRun: initialDigest !== finalDigest,
  digestScope: sources.map((p) => relative(root, p)),
  toolchain: { go: goVersion.output.trim(), postgres: pgVersion.output.trim(), postman: postmanVersion.output.trim() },
  independentOpenAPI: independentValidation,
  postmanSpecLint: { status: specLint.exitCode === 0 ? 'PASS' : 'FAIL', exitCode: specLint.exitCode,
    violations: diagnostics?.violations?.length ?? null, knownMediaRangeOnly,
    limitation: specLint.exitCode === 0 ? null : 'Recorded disagreement over RFC 9457 application/problem+json media type; NOT a lint pass or approved waiver. No authenticated governance review.',
    diagnosticFiles: ['postman-spec-lint.stdout.txt', 'postman-spec-lint.stderr.txt'] },
  collectionLint: { status: collectionLint.exitCode === 0 ? 'PASS' : 'FAIL', ...collectionDiagnostics },
  raceDetector: true, repeatCount, testExitCode: tests.exitCode,
  uniqueTopLevelTests: new Set(passedCases.filter((e) => !e.name.includes('/')).map((e) => e.name)).size,
  uniqueCaseResultsIncludingSubtests: new Set(passedCases.map((e) => e.name)).size,
  caseResults: cases.length, passedCaseResults: passedCases.length, tests: cases,
  postmanRuns: { requests, assertions, loopbackGuardScripts: guards,
    totalRequests: requests.reduce((sum, m) => sum + m.executed, 0),
    totalAssertions: assertions.reduce((sum, m) => sum + m.executed, 0) },
  databaseSafety: cluster,
  releaseGatesPassed: false, physicalDeviceVerified: false,
  notVerified: ['Real identity or current staff-membership authorization', 'Android/server integration, APK or physical phone',
    'Network/device latency, sustained throughput, process-crash/power-loss recovery or restore',
    'Multi-resource, pooled or nightly/hotel inventory', 'Customer link/guest-token security or actual quotation UI',
    'Notice delivery, payments, social connectors, in-app messaging or calls',
    'Production migration/backup, independent reviewer, licensed counsel/provider/pilot/release approval',
    'Complete SCN scenario suite or G7/G8/G9 passage'],
};
save('validation-receipt.json', report);
process.stdout.write(JSON.stringify({ status: report.status, coreChecksPassed,
  passedCaseResults: report.passedCaseResults, repeatCount, requests: report.postmanRuns.totalRequests,
  assertions: report.postmanRuns.totalAssertions, postmanSpecLint: report.postmanSpecLint.status,
  receipt: relative(root, join(artifacts, 'validation-receipt.json')) }, null, 2) + '\n');
// Nonzero for the known linter limitation keeps automation from interpreting this as a full pass.
process.exitCode = boundedResult && specLint.exitCode === 0 ? 0 : boundedResult ? 2 : 1;
