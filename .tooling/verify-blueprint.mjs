import { existsSync, mkdirSync, readFileSync, readdirSync, statSync, writeFileSync } from 'node:fs';
import { createHash } from 'node:crypto';
import { dirname, join, relative, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

// Documentation checks only. This is not application, legal or release validation.
export function verifyBlueprint(projectRoot) {
const root = resolve(projectRoot);
const blueprint = join(root, 'docs/blueprint');
const errors = [];
const warnings = [];
const walk = (dir) => readdirSync(dir).flatMap((name) => {
  const path = join(dir, name);
  return statSync(path).isDirectory() ? walk(path) : [path];
});
const files = walk(blueprint);
const markdown = files.filter((path) => path.endsWith('.md'));
const text = new Map(markdown.map((path) => [path, readFileSync(path, 'utf8')]));
const exists = (path) => files.includes(path);
const required = ['README.md', '01-requirements-ledger.md', '02-decision-register.md',
                  'specifications/20-functional-requirements.md',
                  'specifications/21-quality-attributes.md',
                  'specifications/22-accessibility-localization-and-calendar.md'];
for (const name of required) {
  if (!text.has(join(blueprint, name))) errors.push(`Missing source-of-truth file: ${name}`);
}
const readme = text.get(join(blueprint, 'README.md')) ?? '';
const catalogue = [...readme.matchAll(/^- `([^`]+\.md)`/gm)].map((m) => m[1]);
for (const path of catalogue) {
  if (!exists(resolve(blueprint, path))) errors.push(`Missing catalogue file: ${path}`);
}
const numbers = new Map();
for (const [path, content] of text) {
  const label = relative(root, path);
  const number = relative(blueprint, path).match(/(?:^|\/)(\d{2})-[^/]+\.md$/)?.[1];
  if (number) {
    const entries = numbers.get(number) ?? [];
    entries.push(label);
    numbers.set(number, entries);
    if (!/^Status:/m.test(content)) errors.push(`Missing status: ${label}`);
  }
  let fence = null;
  for (const line of content.split('\n')) {
    const opening = line.match(/^ {0,3}(`{3,}|~{3,})/);
    if (!opening) continue;
    if (!fence) fence = opening[1];
    else if (opening[1][0] === fence[0] && opening[1].length >= fence.length &&
             /^ {0,3}(`+|~+)\s*$/.test(line)) fence = null;
  }
  if (fence) errors.push(`Unclosed Markdown fence: ${label}`);
  // Diagram paths are written relative to the controlled blueprint root.
  for (const m of content.matchAll(/`(diagrams\/[^`]+\.mmd)`/g)) {
    if (!exists(resolve(blueprint, m[1]))) errors.push(`Missing diagram ${m[1]} in ${label}`);
  }
  for (const m of content.matchAll(/\]\(([^)]+\.md(?:#[^)]+)?)\)/g)) {
    const target = m[1].split('#')[0];
    if (/^(https?:|codex:)/.test(target)) continue;
    if (!existsSync(resolve(dirname(path), target))) errors.push(`Broken Markdown link ${target} in ${label}`);
  }
}
for (const [number, entries] of numbers) {
  if (entries.length > 1) errors.push(`Duplicate numbered document ${number}: ${entries.join(', ')}`);
}
const identifier = /\b[A-Z][A-Z0-9]*(?:-[A-Z][A-Z0-9]*)*-\d{3}\b/g;
const ledger = text.get(join(blueprint, '01-requirements-ledger.md')) ?? '';
const requirements = [...ledger.matchAll(/^\| ([A-Z][A-Z0-9-]*-\d{3}) \|/gm)].map((m) => m[1]);
const register = text.get(join(blueprint, '02-decision-register.md')) ?? '';
const decisions = [...register.matchAll(/^\| (DEC-\d{3}) \|/gm)].map((m) => m[1]);
for (const [name, ids] of [['requirement', requirements], ['decision', decisions]]) {
  const seen = new Set();
  for (const id of ids) {
    if (seen.has(id)) errors.push(`Duplicate ${name} definition: ${id}`);
    seen.add(id);
  }
}
const tracePath = join(blueprint, 'specifications/71-end-to-end-validation-and-traceability.md');
const known = new Set();
for (const source of ['01-requirements-ledger.md', 'specifications/20-functional-requirements.md',
                      'specifications/21-quality-attributes.md', 'specifications/22-accessibility-localization-and-calendar.md',
                      '02-decision-register.md']) {
  for (const id of (text.get(join(blueprint, source)) ?? '').match(identifier) ?? []) known.add(id);
}
let scenarioCount = 0;
if (text.has(tracePath)) {
  const trace = text.get(tracePath);
  const usedIds = new Set(trace.match(identifier) ?? []);
  for (const id of requirements) {
    if (!usedIds.has(id)) errors.push(`Requirement absent from Stage 9 crosswalk: ${id}`);
  }
  for (const id of trace.match(identifier) ?? []) {
    if (!id.startsWith('SCN-') && !known.has(id)) errors.push(`Unknown crosswalk identifier: ${id}`);
  }
  const scenarios = [...trace.matchAll(/^\| (SCN-\d{3}) \|/gm)].map((m) => m[1]);
  scenarioCount = scenarios.length;
  if (scenarios.length !== new Set(scenarios).size) errors.push('Duplicate scenario definition');
  for (const id of usedIds) {
    if (id.startsWith('SCN-') && !scenarios.includes(id)) errors.push(`Undefined scenario reference: ${id}`);
  }
} else warnings.push('Stage 9 crosswalk does not exist yet');
warnings.push('Mermaid syntax/rendering, application tests, devices, providers and legal approval are not checked');
const sourceEntries = files.filter((path) => /\.(md|mmd)$/.test(path)).sort().map((path) => [
  relative(blueprint, path), createHash('sha256').update(readFileSync(path)).digest('hex'),
]);
return {
  checkedAt: new Date().toISOString(), scope: 'Documentation structure and identifier coverage only',
  sourceDigest: createHash('sha256').update(JSON.stringify(sourceEntries)).digest('hex'),
  digestScope: 'All blueprint Markdown and Mermaid source paths/content, excluding generated report',
  markdownFiles: markdown.length, diagramSources: files.filter((p) => p.endsWith('.mmd')).length,
  catalogueEntries: catalogue.length, ledgerRequirements: requirements.length,
  decisionDefinitions: decisions.length, scenarioDefinitions: scenarioCount,
  status: errors.length ? 'FAIL' : 'PASS', errors: [...new Set(errors)], warnings,
};
}

if (process.argv[1] && resolve(process.argv[1]) === fileURLToPath(import.meta.url)) {
  const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
  const report = verifyBlueprint(root);
  const rendered = JSON.stringify(report, null, 2) + '\n';
  if (process.argv.includes('--report')) {
    const target = join(root, 'docs/blueprint/validation/document-structure-report.json');
    mkdirSync(dirname(target), { recursive: true });
    writeFileSync(target, rendered);
  }
  console.log(rendered.trimEnd());
  process.exitCode = report.status === 'PASS' ? 0 : 1;
}
