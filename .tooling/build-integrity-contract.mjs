import { mkdirSync, writeFileSync } from 'node:fs';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const text = { type: 'string', minLength: 1, maxLength: 128 };
const integer = { type: 'integer', minimum: 1 };
const ref = (name) => ({ $ref: `#/components/schemas/${name}` });
const object = (properties, required) => ({ type: 'object', additionalProperties: false, properties, required });
const scope = { organization_id: text, operation_id: text };
const intent = { ...scope, service_id: text, expected_service_version: integer,
  start_at: { type: 'string', format: 'date-time' },
  customer_label: { type: 'string', enum: ['Demo Guest 01', 'Demo Guest 02', 'Demo Guest 03'] },
  source: { type: 'string', enum: ['Call', 'Walk-in', 'WhatsApp', 'Telegram', 'Other'] } };
const allIntent = Object.keys(intent);
const commonResult = { demo_only: { type: 'boolean', const: true }, id: text,
  organization_id: text, operation_id: text, version: integer, state: text,
  start_at: { type: 'string', format: 'date-time' }, end_at: { type: 'string', format: 'date-time' },
  service_id: text, service_version: integer, price_minor: { type: 'integer', minimum: 0 },
  currency: { type: 'string', const: 'ETB' }, customer_label: text, source: text,
  expires_at: { type: 'string', format: 'date-time' }, request_id: text, booking_id: text };
const schemas = {
  BookingInput: object(intent, allIntent),
  HoldInput: object({ ...intent, ttl_seconds: { type: 'integer', minimum: 1, maximum: 900 } }, [...allIntent, 'ttl_seconds']),
  QuoteInput: object({ ...intent, request_id: { type: 'string', pattern: '^demo-request-[a-z0-9-]+$', maxLength: 64 },
    demo_price_minor: { type: 'integer', minimum: 0, maximum: 10000000 } }, [...allIntent, 'request_id']),
  VersionInput: object({ expected_version: integer }, ['expected_version']),
  RescheduleInput: object({ expected_version: integer, start_at: { type: 'string', format: 'date-time' } }, ['expected_version', 'start_at']),
  Result: object(commonResult, ['demo_only', 'id', 'organization_id', 'operation_id', 'version', 'state',
    'start_at', 'end_at', 'service_id', 'service_version', 'price_minor', 'currency', 'customer_label', 'source']),
  Health: object({ demo_only: { type: 'boolean', const: true }, status: { type: 'string', const: 'synthetic_only' } }, ['demo_only', 'status']),
  Problem: object({ type: { type: 'string' }, title: text, status: { type: 'integer', minimum: 400, maximum: 599 },
    detail: text, code: text, correlation_id: text, retryable: { type: 'boolean' }, demo_only: { type: 'boolean', const: true } },
    ['type', 'title', 'status', 'detail', 'code', 'correlation_id', 'retryable', 'demo_only']),
};
const problem = { description: 'Safe synthetic Problem Details; no internal SQL or foreign-tenant details',
  content: { 'application/problem+json': { schema: ref('Problem') } } };
const response = (schema, description) => ({ description,
  headers: { 'Cache-Control': { schema: { type: 'string', const: 'no-store' } }, 'X-Request-ID': { schema: text } },
  content: { 'application/json': { schema: ref(schema) } } });
const key = { name: 'Idempotency-Key', in: 'header', required: true,
  schema: { type: 'string', pattern: '^[A-Za-z0-9_.:-]{8,128}$' }, description: 'Stable successful-command retry key; not authentication' };
const pathId = { name: 'id', in: 'path', required: true, schema: text };
const operation = (id, input, status, description, withId = false) => ({
  operationId: id, summary: description, security: [{ demoPrincipal: [] }],
  parameters: [...(input ? [key] : []), ...(withId ? [pathId] : [])],
  ...(input ? { requestBody: { required: true, content: { 'application/json': { schema: ref(input) } } } } : {}),
  responses: { [status]: response('Result', description), '400': problem, '401': problem,
    '403': problem, '404': problem, '409': problem, '413': problem, '415': problem, '422': problem, '500': problem },
});
const spec = {
  openapi: '3.1.0', info: { title: 'Ubi synthetic booking-integrity V0', version: '0.0.1',
    description: 'LOCAL SYNTHETIC EXPERIMENT ONLY. Demo principal is a fixture selector, not authentication. No real bookings, notices or money. Production needs HTTPS and real authorization.' },
  servers: [{ url: 'http://127.0.0.1:8080', description: 'Illustrative test-only loopback; tests use a random local port' }],
  paths: {
    '/v0/health': { get: { operationId: 'getSyntheticHealth', security: [], responses: { '200': response('Health', 'Synthetic marker') } } },
    '/v0/bookings': { post: operation('createDemoBooking', 'BookingInput', '201', 'Commit an invented exclusive-resource booking') },
    '/v0/bookings/{id}': { get: operation('getDemoBooking', null, '200', 'Read authorized synthetic booking without mutation', true) },
    '/v0/bookings/{id}/reschedules': { post: operation('rescheduleDemoBooking', 'RescheduleInput', '200', 'Atomically move a synthetic booking', true) },
    '/v0/bookings/{id}/cancellations': { post: operation('cancelDemoBooking', 'VersionInput', '200', 'Cancel locally in the disposable database', true) },
    '/v0/holds': { post: operation('createDemoHold', 'HoldInput', '201', 'Hold invented capacity; no confirmation') },
    '/v0/holds/{id}/bookings': { post: operation('consumeDemoHold', 'VersionInput', '201', 'Convert the same held reservation to a synthetic booking', true) },
    '/v0/quotations': { post: operation('issueDemoQuote', 'QuoteInput', '201', 'Issue immutable synthetic quotation version') },
    '/v0/quotations/{id}/acceptances': { post: operation('acceptDemoQuote', 'VersionInput', '201', 'Atomically convert an exact valid quotation', true) },
  },
  components: { securitySchemes: { demoPrincipal: { type: 'apiKey', in: 'header', name: 'X-Ubi-Demo-Principal',
    description: 'Not authentication. Hard-coded invented fixture selector; loopback tests only.' } }, schemas },
};
const target = resolve(root, 'prototypes/booking_integrity/openapi.json');
mkdirSync(dirname(target), { recursive: true });
writeFileSync(target, JSON.stringify(spec, null, 2) + '\n');
console.log('Generated synthetic OpenAPI contract: prototypes/booking_integrity/openapi.json');
