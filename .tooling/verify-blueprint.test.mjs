import assert from 'node:assert/strict';
import { mkdirSync, mkdtempSync, rmSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { dirname, join } from 'node:path';
import test from 'node:test';
import { verifyBlueprint } from './verify-blueprint.mjs';

function fixture(run) {
  const root = mkdtempSync(join(tmpdir(), 'ubi-blueprint-check-'));
  const put = (name, content) => {
    const path = join(root, 'docs/blueprint', name);
    mkdirSync(dirname(path), { recursive: true });
    writeFileSync(path, content);
  };
  put('README.md', '# Test\n- `01-requirements-ledger.md`\n');
  put('01-requirements-ledger.md', 'Status: Draft\n| PRD-001 | Test requirement |\n');
  put('02-decision-register.md', 'Status: Draft\n| DEC-001 | Test decision |\n');
  put('specifications/20-functional-requirements.md', 'Status: Draft\n| FR-BKG-001 | Test |\n');
  put('specifications/21-quality-attributes.md', 'Status: Draft\n');
  put('specifications/22-accessibility-localization-and-calendar.md', 'Status: Draft\n');
  put('specifications/71-end-to-end-validation-and-traceability.md',
      'Status: Draft\nPRD-001, FR-BKG-001 → SCN-001\n| SCN-001 | Defined scenario |\n');
  try { run({ root, put }); }
  finally { rmSync(root, { recursive: true, force: true }); }
}

test('valid structure passes, without implying executed tests', () => fixture(({ root }) => {
  const report = verifyBlueprint(root);
  assert.equal(report.status, 'PASS');
  assert.equal(report.scenarioDefinitions, 1);
  assert.match(report.warnings.join(' '), /not checked/);
}));

test('missing catalogue file fails', () => fixture(({ root, put }) => {
  put('README.md', '- `specifications/23-missing.md`\n');
  assert.match(verifyBlueprint(root).errors.join(' '), /Missing catalogue file/);
}));

test('duplicate document number and missing status fail', () => fixture(({ root, put }) => {
  put('specifications/20-duplicate.md', '# Missing status\n');
  assert.match(verifyBlueprint(root).errors.join(' '), /Duplicate numbered document/);
  assert.match(verifyBlueprint(root).errors.join(' '), /Missing status/);
}));

test('unclosed fence and missing diagram fail', () => fixture(({ root, put }) => {
  put('specifications/23-invalid.md', 'Status: Draft\n`diagrams/absent.mmd`\n```text\nunclosed\n');
  assert.match(verifyBlueprint(root).errors.join(' '), /Unclosed Markdown fence/);
  assert.match(verifyBlueprint(root).errors.join(' '), /Missing diagram/);
}));

test('longer identifier cannot falsely cover a requirement', () => fixture(({ root, put }) => {
  put('specifications/71-end-to-end-validation-and-traceability.md', 'Status: Draft\nPRD-0011\n');
  assert.match(verifyBlueprint(root).errors.join(' '), /Requirement absent.*PRD-001/);
}));

test('unknown requirement and undefined scenario fail', () => fixture(({ root, put }) => {
  put('specifications/71-end-to-end-validation-and-traceability.md',
      'Status: Draft\nPRD-001, FR-BKG-999, SCN-999\n');
  assert.match(verifyBlueprint(root).errors.join(' '), /Unknown crosswalk identifier/);
  assert.match(verifyBlueprint(root).errors.join(' '), /Undefined scenario reference/);
}));

test('duplicate definitions and broken local links fail', () => fixture(({ root, put }) => {
  put('01-requirements-ledger.md',
      'Status: Draft\n| PRD-001 | A |\n| PRD-001 | B |\n[missing](missing.md)\n');
  put('02-decision-register.md', 'Status: Draft\n| DEC-001 | A |\n| DEC-001 | B |\n');
  assert.match(verifyBlueprint(root).errors.join(' '), /Duplicate requirement definition/);
  assert.match(verifyBlueprint(root).errors.join(' '), /Duplicate decision definition/);
  assert.match(verifyBlueprint(root).errors.join(' '), /Broken Markdown link/);
}));

test('mixed fences inside a valid fenced block do not terminate it', () => fixture(({ root, put }) => {
  put('specifications/23-fences.md', 'Status: Draft\n````text\n```\n~~~\n````\n');
  assert.equal(verifyBlueprint(root).status, 'PASS');
}));

test('existing prototype Markdown outside blueprint passes; missing target fails', () => fixture(({ root, put }) => {
  put('../../prototypes/README.md', '# Synthetic prototype\n');
  put('validation/receipt.md', '[prototype](../../../prototypes/README.md)\n');
  assert.equal(verifyBlueprint(root).status, 'PASS');
  put('validation/receipt.md', '[missing](../../../prototypes/missing.md)\n');
  assert.match(verifyBlueprint(root).errors.join(' '), /Broken Markdown link/);
}));
