import { mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
const root=resolve(dirname(fileURLToPath(import.meta.url)),'..');
const project=join(root,'prototypes/booking_integrity');
const spec=JSON.parse(readFileSync(join(project,'openapi.json')));
const dir=join(project,'postman/collections/ubi-synthetic-v0');
mkdirSync(join(dir,'.resources'),{recursive:true});
const quote = (s) => `'${s.replaceAll("'","''")}'`;
writeFileSync(join(dir,'.resources/definition.yaml'),`$kind: collection\nname: 'Ubi synthetic V0 contract checks'\ndescription: 'Local synthetic fixtures only; no real authentication, messages or money.'\nvariables:\n  - key: base_url\n    value: 'http://127.0.0.1:1'\nscripts:\n  - type: 'http:beforeRequest'\n    language: 'text/javascript'\n    code: |-\n      if (!/^http:\\/\\/127\\.0\\.0\\.1:\\d+$/.test(pm.variables.get('base_url'))) throw new Error('Refusing non-loopback synthetic target');\n`);
const intent=(hour=9)=>({organization_id:'demo-org-a',operation_id:'studio',service_id:'portrait',expected_service_version:1,
 start_at:`2030-01-15T${String(hour).padStart(2,'0')}:00:00+03:00`,customer_label:'Demo Guest 01',source:'Call'});
const cases=[];
const add=(name,method,path,status,body,key,extra={})=>cases.push({name,method,path,status,body,key,...extra});
add('health','GET','/v0/health',200,null,null);
add('create-booking','POST','/v0/bookings',201,intent(),'postman-create-001',{save:'booking_id'});
add('retry-booking','POST','/v0/bookings',201,intent(),'postman-create-001',{same:'booking_id'});
add('read-own-booking','GET','/v0/bookings/{{booking_id}}',200,null,null,{same:'booking_id'});
add('read-foreign-booking','GET','/v0/bookings/{{booking_id}}',404,null,null,{principal:'demo-owner-b',code:'not_found'});
add('capacity-conflict','POST','/v0/bookings',409,intent(),'postman-overlap-001',{code:'capacity_conflict'});
add('changed-input-same-key','POST','/v0/bookings',409,intent(10),'postman-create-001',{code:'idempotency_mismatch'});
add('reschedule','POST','/v0/bookings/{{booking_id}}/reschedules',200,{expected_version:1,start_at:intent(10).start_at},'postman-reschedule-001',{same:'booking_id',version:2});
add('stale-reschedule','POST','/v0/bookings/{{booking_id}}/reschedules',409,{expected_version:1,start_at:intent(11).start_at},'postman-stale-001',{code:'booking_version_conflict'});
add('cancel','POST','/v0/bookings/{{booking_id}}/cancellations',200,{expected_version:2},'postman-cancel-001',{state:'cancelled_demo'});
add('book-released-capacity','POST','/v0/bookings',201,intent(),'postman-replacement-001');
add('create-hold','POST','/v0/holds',201,{...intent(11),ttl_seconds:300},'postman-hold-001',{save:'hold_id',state:'active'});
add('consume-hold','POST','/v0/holds/{{hold_id}}/bookings',201,{expected_version:1},'postman-consume-001',{save:'held_booking_id'});
add('retry-hold-conversion','POST','/v0/holds/{{hold_id}}/bookings',201,{expected_version:1},'postman-consume-001',{same:'held_booking_id'});
add('issue-quote','POST','/v0/quotations',201,{...intent(14),request_id:'demo-request-postman',demo_price_minor:223459},'postman-quote-001',{save:'quote_id',state:'issued'});
add('accept-quote','POST','/v0/quotations/{{quote_id}}/acceptances',201,{expected_version:1},'postman-accept-001',{save:'quoted_booking_id',price:223459});
add('retry-quote-conversion','POST','/v0/quotations/{{quote_id}}/acceptances',201,{expected_version:1},'postman-accept-001',{same:'quoted_booking_id'});
add('revoked-fixture','POST','/v0/bookings',403,intent(16),'postman-revoked-001',{principal:'demo-revoked',code:'permission_denied'});
add('unknown-field','POST','/v0/bookings',400,{...intent(16),confirmed:true},'postman-unknown-001',{code:'invalid_json'});
add('fractional-money','POST','/v0/quotations',400,{...intent(16),request_id:'demo-request-fraction',demo_price_minor:2.9},'postman-fractional-001',{code:'invalid_json'});
add('missing-key','POST','/v0/bookings',400,intent(16),null,{code:'invalid_idempotency_key'});
add('wrong-tenant','POST','/v0/bookings',403,{...intent(16),organization_id:'demo-org-b'},'postman-wrong-tenant-001',{code:'permission_denied'});
add('wrong-media-type','POST','/v0/bookings',415,intent(16),'postman-media-001',{contentType:'text/plain',code:'unsupported_media_type'});
add('missing-principal','POST','/v0/bookings',401,intent(16),'postman-no-principal-001',{principal:null,code:'demo_principal_required'});
add('oversized-body','POST','/v0/bookings',413,{...intent(16),customer_label:'X'.repeat(40000)},'postman-large-001',{code:'body_too_large'});

for(const [i,c] of cases.entries()) {
 const template=c.path.replaceAll(/\{\{[^}]+\}\}/g,'{id}');
 const op=spec.paths[template]?.[c.method.toLowerCase()];
 if(!op?.responses[String(c.status)])throw Error(`Case not covered by declared contract: ${c.name}`);
 const schema=c.status>=400?spec.components.schemas.Problem:c.path==='/v0/health'?spec.components.schemas.Health:spec.components.schemas.Result;
 const script=[`pm.test('status ${c.status}',()=>pm.response.to.have.status(${c.status}));`,
  "const body=pm.response.json();",
  "pm.test('synthetic marker',()=>pm.expect(body.demo_only).to.eql(true));",
  "pm.test('no-store',()=>pm.expect(pm.response.headers.get('Cache-Control')).to.eql('no-store'));",
  "pm.test('correlation ID',()=>pm.expect(pm.response.headers.get('X-Request-ID')).to.match(/^[a-f0-9-]{36}$/));",
  `pm.test('declared response schema',()=>pm.response.to.have.jsonSchema(${JSON.stringify(schema)}));`,
  `pm.test('declared content type',()=>pm.expect(pm.response.headers.get('Content-Type')).to.include('${c.status>=400?'application/problem+json':'application/json'}'));`,
  ...(c.save?[`pm.collectionVariables.set('${c.save}',body.id);`]:[]),
  ...(c.same?[`pm.test('same committed identity',()=>pm.expect(body.id).to.eql(pm.collectionVariables.get('${c.same}')));`]:[]),
  ...(c.code?[`pm.test('safe domain code',()=>pm.expect(body.code).to.eql('${c.code}'));`]:[]),
  ...(c.version?[`pm.test('version',()=>pm.expect(body.version).to.eql(${c.version}));`]:[]),
  ...(c.state?[`pm.test('truthful state',()=>pm.expect(body.state).to.eql('${c.state}'));`]:[]),
  ...(c.price?[`pm.test('integer quotation snapshot',()=>pm.expect(body.price_minor).to.eql(${c.price}));`]:[]),
 ].join('\n');
 const headers=[];if(c.body)headers.push(['Content-Type',c.contentType??'application/json']);
 const principal=Object.hasOwn(c,'principal')?c.principal:'demo-owner-a';if(principal)headers.push(['X-Ubi-Demo-Principal',principal]);if(c.key)headers.push(['Idempotency-Key',c.key]);
 const value=`$kind: http-request\nmethod: ${c.method}\nurl: ${quote('{{base_url}}'+c.path)}\norder: ${(i+1)*1000}\n`+
  (headers.length?'headers:\n'+headers.map(([key,value])=>`  - key: ${quote(key)}\n    value: ${quote(value)}\n`).join(''):'')+
  (c.body?'body:\n  type: json\n  content: |-\n'+JSON.stringify(c.body,null,2).split('\n').map(l=>'    '+l).join('\n')+'\n':'')+
  'scripts:\n  - type: beforeRequest\n    language: text/javascript\n    code: |-\n'+
  "      if (!/^http:\\/\\/127\\.0\\.0\\.1:\\d+$/.test(pm.variables.get('base_url'))) { pm.execution.skipRequest(); throw new Error('Refusing non-loopback synthetic target'); }\n"+
  '  - type: afterResponse\n    language: text/javascript\n    code: |-\n'+script.split('\n').map(l=>'      '+l).join('\n')+'\n';
 writeFileSync(join(dir,c.name+'.request.yaml'),value);
}
console.log(`Generated ${cases.length} local-only v3 contract requests with persistent assertions.`);
