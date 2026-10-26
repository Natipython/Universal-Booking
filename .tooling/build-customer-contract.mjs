import { mkdirSync, writeFileSync } from 'node:fs';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const project = join(root, 'prototypes/customer_web');
const fields = {
  csrf: { type: 'string', pattern: '^[A-Za-z0-9_-]{43}$' },
  token: { type: 'string', pattern: '^[A-Za-z0-9_-]{43}$' },
  scenario: { type: 'string', enum: ['instant', 'custom', 'river'] },
  expected_draft_version: { type: 'integer', minimum: 1 },
  service_id: { type: 'string', enum: ['portrait', 'consultation'] },
  date: { type: 'string', format: 'date' },
  time: { type: 'string', pattern: '^([01][0-9]|2[0-3]):[0-5][0-9]$' },
  calendar: { type: 'string', enum: ['Gregorian'] },
  clock: { type: 'string', enum: ['international-24h'] },
  customer_label: { type: 'string', enum: ['Demo Guest 01', 'Demo Guest 02', 'Demo Guest 03'] },
  event_size: { type: 'string', enum: ['not-specified', 'small', 'medium', 'large'] },
  review_id: { type: 'string', pattern: '^[A-Za-z0-9_-]{43}$' },
};
const spec = { openapi: '3.1.0', info: { title: 'Ubi synthetic customer handoff V0', version: '0.0.1',
  description: 'Loopback-only HTML/form experiment. Dummy records and fixture grants are NOT verified customer identity. No public/live API approval.' },
  servers: [{ url: 'http://127.0.0.1:1', description: 'Replace only with the temporary loopback test origin.' }],
  paths: {}, components: { securitySchemes: {
    BootstrapCookie: { type: 'apiKey', in: 'cookie', name: 'ubi_demo_bootstrap', description: 'Demo browser binding only; NOT customer authentication.' },
    DemoSessionCookie: { type: 'apiKey', in: 'cookie', name: 'ubi_demo_session', description: 'Ephemeral scoped demo session; no real identity, history or manage permission.' },
  } } };
function response(code) {
  return { description: code === 303 ? 'Same-origin redirect after a bounded form action; no secret in destination except initial opaque handoff.' : code < 400 ? 'Synthetic HTML page, no booking mutation on GET/HEAD.' : 'Safe HTML error with stable data-problem code; no private scope disclosure.',
    headers: { 'Cache-Control': { schema: { type: 'string', const: 'no-store' } },
      'Referrer-Policy': { schema: { type: 'string', const: 'no-referrer' } },
      ...(code === 303 ? { Location: { schema: { type: 'string' } } } : {}) },
    ...(code !== 303 ? { content: { 'text/html': { schema: { type: 'string' } } } } : {}) };
}
function op(path, method, id, inputs = [], auth = '') {
  const success = method === 'post' ? 303 : 200;
  const codes = [...new Set([success, 303, 400, 401, 403, 404, 409, 410, 413, 415, 422, 500, 503])];
  const value = { operationId: id, security: auth ? [{ [auth]: [] }] : [],
    responses: Object.fromEntries(codes.map(c => [String(c), response(c)])) };
  if (inputs.length) value.requestBody = { required: true, content: { 'application/x-www-form-urlencoded': {
    schema: { type: 'object', additionalProperties: false, required: inputs, properties: Object.fromEntries(inputs.map(k => [k, fields[k]])) },
  } } };
  if (path.includes('{token}')) value.parameters = [{ name: 'token', in: 'path', required: true, schema: fields.token }];
  (spec.paths[path] ??= {})[method] = value;
}
op('/', 'get', 'publicStudio');
op('/', 'head', 'publicStudioHeaders');
op('/h/{token}', 'get', 'safeHandoffPreview');
op('/h/{token}', 'head', 'safeHandoffPreviewHeaders');
op('/demo/starts', 'post', 'startInventedScenario', ['csrf', 'scenario'], 'BootstrapCookie');
op('/v0/handoff-sessions', 'post', 'redeemScopedHandoff', ['csrf', 'token'], 'BootstrapCookie');
for (const path of ['/continue', '/review', '/outcome']) {
  op(path, 'get', path.slice(1), [], 'DemoSessionCookie');
  op(path, 'head', path.slice(1)+'Headers', [], 'DemoSessionCookie');
}
op('/v0/reviews', 'post', 'reviewEditedIntent', ['csrf', 'expected_draft_version', 'service_id', 'date', 'time', 'calendar', 'clock', 'customer_label', 'event_size'], 'DemoSessionCookie');
op('/v0/submissions', 'post', 'submitReviewedIntent', ['csrf', 'review_id'], 'DemoSessionCookie');
mkdirSync(project, { recursive: true });
writeFileSync(join(project, 'openapi.json'), JSON.stringify(spec, null, 2) + '\n');

const collection = join(project, 'postman/collections/ubi-customer-v0');
mkdirSync(join(collection, '.resources'), { recursive: true });
writeFileSync(join(collection, '.resources/definition.yaml'), "$kind: collection\nname: 'Ubi local customer handoff checks'\ndescription: 'Synthetic scoped forms only; no account/cloud, real contact or token exports.'\nvariables:\n  - key: base_url\n    value: 'http://127.0.0.1:1'\n");
const quote = s => `'${s.replaceAll("'", "''")}'`;
const cases = [];
const add = (name, method, path, status, form = null, scripts = [], session = false) => cases.push({ name, method, path, status, form, scripts, session });
const csrf = ["const csrf=body.match(/name=\"csrf\" value=\"([A-Za-z0-9_-]{43})\"/);", "pm.test('CSRF form present',()=>pm.expect(Boolean(csrf)).to.eql(true));", "if(csrf)pm.collectionVariables.set('csrf',csrf[1]);"];
const cookie = name => [`const cookie=pm.response.headers.get('Set-Cookie')?.match(/${name}=([A-Za-z0-9_-]{43})/);`, "pm.test('scoped cookie issued',()=>pm.expect(Boolean(cookie)).to.eql(true));", `if(cookie)pm.collectionVariables.set('${name === 'ubi_demo_bootstrap' ? 'bootstrap' : 'session'}',cookie[1]);`];
const redirectToken = ["const location=pm.response.headers.get('Location');", "pm.test('opaque link only',()=>pm.expect(location).to.match(/^\\/h\\/[A-Za-z0-9_-]{43}$/));", "pm.collectionVariables.set('token',location.split('/').pop());"];
add('public-page', 'GET', '/', 200, null, [...csrf, ...cookie('ubi_demo_bootstrap')]);
add('start-instant', 'POST', '/demo/starts', 303, { csrf: '{{csrf}}', scenario: 'instant' }, redirectToken);
add('safe-preview', 'GET', '/h/{{token}}', 200, null, [...csrf, "pm.test('no private draft in preview',()=>pm.expect(body).not.to.include('2030-01-15'));", "pm.test('no contact prefill',()=>pm.expect(body).not.to.include('Demo Guest 01'));"]);
add('redeem', 'POST', '/v0/handoff-sessions', 303, { csrf: '{{csrf}}', token: '{{token}}' }, cookie('ubi_demo_session'));
add('edit-context', 'GET', '/continue', 200, null, [...csrf, "pm.test('conversation provenance',()=>pm.expect(body).to.include('From your conversation'));"], true);
const reviewForm = () => ({ csrf: '{{csrf}}', expected_draft_version: '1', service_id: 'portrait', date: '2030-01-15', time: '10:00', calendar: 'Gregorian', clock: 'international-24h', customer_label: 'Demo Guest 02', event_size: 'not-specified' });
add('save-edited-review', 'POST', '/v0/reviews', 303, reviewForm(), [], true);
const reviewID = [...csrf, "const review=body.match(/name=\"review_id\" value=\"([A-Za-z0-9_-]{43})\"/);", "pm.test('owned review present',()=>pm.expect(Boolean(review)).to.eql(true));", "if(review)pm.collectionVariables.set('review',review[1]);", "pm.test('edited time retained',()=>pm.expect(body).to.include('10:00'));"];
add('review-before-commit', 'GET', '/review', 200, null, reviewID, true);
add('reject-forged-CSRF', 'POST', '/v0/submissions', 403, { csrf: 'wrong', review_id: '{{review}}' }, ["pm.test('stable rejection',()=>pm.expect(body).to.include('data-problem=\"csrf_invalid\"'));"], true);
add('commit-reviewed-demo', 'POST', '/v0/submissions', 303, { csrf: '{{csrf}}', review_id: '{{review}}' }, [], true);
add('confirmed-outcome', 'GET', '/outcome', 200, null, ["pm.test('synthetic confirmation only',()=>pm.expect(body).to.include('data-state=\"confirmed_demo\"'));", "pm.test('no payment claim',()=>pm.expect(body).to.include('No payment was taken'));"] , true);
add('retry-same-review', 'POST', '/v0/submissions', 303, { csrf: '{{csrf}}', review_id: '{{review}}' }, [], true);
add('invalid-session-cookie', 'GET', '/continue', 401);
add('public-again', 'GET', '/', 200, null, csrf);
add('start-custom', 'POST', '/demo/starts', 303, { csrf: '{{csrf}}', scenario: 'custom' }, redirectToken);
add('custom-preview', 'GET', '/h/{{token}}', 200, null, csrf);
add('redeem-custom', 'POST', '/v0/handoff-sessions', 303, { csrf: '{{csrf}}', token: '{{token}}' }, cookie('ubi_demo_session'));
add('edit-custom', 'GET', '/continue', 200, null, csrf, true);
add('review-custom', 'POST', '/v0/reviews', 303, { ...reviewForm(), service_id: 'consultation', time: '11:00', event_size: 'small' }, [], true);
add('custom-summary', 'GET', '/review', 200, null, [...reviewID.filter(s => !s.includes('edited time retained')), "pm.test('request expectation',()=>pm.expect(body).to.include('No time is held'));"] , true);
add('send-custom-request', 'POST', '/v0/submissions', 303, { csrf: '{{csrf}}', review_id: '{{review}}' }, [], true);
add('request-not-confirmation', 'GET', '/outcome', 200, null, ["pm.test('request saved only',()=>pm.expect(body).to.include('data-state=\"request_saved_demo\"'));", "pm.test('not booked',()=>pm.expect(body).to.include('No time is held'));"], true);
for (const [i, c] of cases.entries()) {
  const template = c.path.replace(/\{\{token\}\}/, '{token}');
  if (!spec.paths[template]?.[c.method.toLowerCase()]?.responses[c.status]) throw Error('Undeclared contract case');
  const headers = [['Origin', '{{base_url}}']];
  // Let the runtime cookie jar carry its issued cookie; an explicit copy causes duplicate-cookie rejection.
  if(c.name==='invalid-session-cookie') headers.push(['Cookie','ubi_demo_session=invalid-demo-session']);
  if (c.form) headers.push(['Content-Type', 'application/x-www-form-urlencoded']);
  const assertions = [`pm.test('status ${c.status}',()=>pm.response.to.have.status(${c.status}));`,
    'const body=pm.response.text();',
    "pm.test('private cache policy',()=>pm.expect(pm.response.headers.get('Cache-Control')).to.eql('no-store'));",
    "pm.test('no referrer leakage',()=>pm.expect(pm.response.headers.get('Referrer-Policy')).to.eql('no-referrer'));",
    "pm.test('no framing',()=>pm.expect(pm.response.headers.get('X-Frame-Options')).to.eql('DENY'));",
    ...c.scripts].join('\n');
  const value = `$kind: http-request\nmethod: ${c.method}\nurl: ${quote('{{base_url}}' + c.path)}\norder: ${(i+1)*1000}\nsettings:\n  followRedirects: false\n` +
    'headers:\n' + headers.map(([k, v]) => `  - key: ${quote(k)}\n    value: ${quote(v)}\n`).join('') +
    (c.form ? 'body:\n  type: urlencoded\n  content:\n' + Object.entries(c.form).map(([k,v]) => `    - key: ${quote(k)}\n      value: ${quote(v)}\n`).join('') : '') +
    'scripts:\n  - type: beforeRequest\n    language: text/javascript\n    code: |-\n' +
    "      if(!/^http:\\/\\/127\\.0\\.0\\.1:\\d+$/.test(pm.variables.get('base_url'))){pm.execution.skipRequest();throw new Error('Refusing non-loopback demo target');}\n" +
    '  - type: afterResponse\n    language: text/javascript\n    code: |-\n' + assertions.split('\n').map(s => '      ' + s).join('\n') + '\n';
  writeFileSync(join(collection, c.name + '.request.yaml'), value);
}
console.log(`Generated OpenAPI and ${cases.length} saved local form-contract requests.`);
