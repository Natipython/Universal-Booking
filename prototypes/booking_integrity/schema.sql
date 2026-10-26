-- DISPOSABLE SYNTHETIC CLUSTER ONLY. Not a production migration.
CREATE EXTENSION btree_gist;
REVOKE CREATE ON SCHEMA public FROM PUBLIC;
CREATE ROLE ubi_v0_runtime LOGIN NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT NOBYPASSRLS;
CREATE SCHEMA ubi_v0;
GRANT USAGE ON SCHEMA ubi_v0 TO ubi_v0_runtime;

CREATE TABLE ubi_v0.organizations (organization_id text PRIMARY KEY);
CREATE TABLE ubi_v0.operations (
  organization_id text NOT NULL REFERENCES ubi_v0.organizations,
  operation_id text NOT NULL, PRIMARY KEY (organization_id, operation_id)
);
CREATE TABLE ubi_v0.resources (
  organization_id text NOT NULL REFERENCES ubi_v0.organizations,
  resource_id text NOT NULL, PRIMARY KEY (organization_id, resource_id)
);
CREATE TABLE ubi_v0.services (
  organization_id text NOT NULL, operation_id text NOT NULL, service_id text NOT NULL,
  resource_id text NOT NULL, version integer NOT NULL CHECK (version > 0),
  duration_minutes integer NOT NULL CHECK (duration_minutes BETWEEN 1 AND 1440),
  buffer_before integer NOT NULL DEFAULT 0 CHECK (buffer_before BETWEEN 0 AND 120),
  buffer_after integer NOT NULL DEFAULT 0 CHECK (buffer_after BETWEEN 0 AND 120),
  price_minor bigint NOT NULL CHECK (price_minor BETWEEN 0 AND 10000000),
  currency text NOT NULL CHECK (currency = 'ETB'),
  PRIMARY KEY (organization_id, operation_id, service_id),
  FOREIGN KEY (organization_id, operation_id) REFERENCES ubi_v0.operations,
  FOREIGN KEY (organization_id, resource_id) REFERENCES ubi_v0.resources
);
CREATE TABLE ubi_v0.bookings (
  organization_id text NOT NULL, id text NOT NULL, operation_id text NOT NULL, service_id text NOT NULL,
  service_version integer NOT NULL CHECK (service_version > 0),
  start_at timestamptz NOT NULL, end_at timestamptz NOT NULL,
  price_minor bigint NOT NULL CHECK (price_minor BETWEEN 0 AND 10000000), currency text NOT NULL CHECK (currency = 'ETB'),
  customer_label text NOT NULL CHECK (customer_label IN ('Demo Guest 01', 'Demo Guest 02', 'Demo Guest 03')),
  source text NOT NULL CHECK (source IN ('Call','Walk-in','WhatsApp','Telegram','Other')),
  state text NOT NULL CHECK (state IN ('confirmed_demo','cancelled_demo')), version integer NOT NULL CHECK (version > 0),
  PRIMARY KEY (organization_id, id),
  FOREIGN KEY (organization_id, operation_id, service_id) REFERENCES ubi_v0.services,
  CHECK (end_at > start_at AND isfinite(start_at) AND isfinite(end_at))
);
CREATE TABLE ubi_v0.holds (
  organization_id text NOT NULL, id text NOT NULL, operation_id text NOT NULL, service_id text NOT NULL,
  service_version integer NOT NULL CHECK (service_version > 0),
  start_at timestamptz NOT NULL, end_at timestamptz NOT NULL,
  price_minor bigint NOT NULL CHECK (price_minor BETWEEN 0 AND 10000000), currency text NOT NULL CHECK (currency = 'ETB'),
  customer_label text NOT NULL CHECK (customer_label IN ('Demo Guest 01', 'Demo Guest 02', 'Demo Guest 03')),
  source text NOT NULL CHECK (source IN ('Call','Walk-in','WhatsApp','Telegram','Other')),
  state text NOT NULL CHECK (state IN ('active','consumed','expired')), version integer NOT NULL CHECK (version > 0),
  created_at timestamptz NOT NULL DEFAULT clock_timestamp(), expires_at timestamptz NOT NULL,
  PRIMARY KEY (organization_id, id),
  FOREIGN KEY (organization_id, operation_id, service_id) REFERENCES ubi_v0.services,
  CHECK (end_at > start_at AND isfinite(start_at) AND isfinite(end_at)),
  CHECK (expires_at > created_at AND isfinite(expires_at))
);
CREATE TABLE ubi_v0.resource_reservations (
  organization_id text NOT NULL, id text NOT NULL, resource_id text NOT NULL,
  occupied_interval tstzrange NOT NULL,
  state text NOT NULL CHECK (state IN ('held','committed','released')),
  hold_id text, booking_id text,
  PRIMARY KEY (organization_id, id),
  FOREIGN KEY (organization_id, resource_id) REFERENCES ubi_v0.resources,
  FOREIGN KEY (organization_id, hold_id) REFERENCES ubi_v0.holds,
  FOREIGN KEY (organization_id, booking_id) REFERENCES ubi_v0.bookings,
  CHECK ((hold_id IS NULL) <> (booking_id IS NULL)),
  CHECK (state <> 'held' OR hold_id IS NOT NULL),
  CHECK (state <> 'committed' OR booking_id IS NOT NULL),
  CHECK (NOT isempty(occupied_interval) AND NOT lower_inf(occupied_interval) AND NOT upper_inf(occupied_interval)
    AND isfinite(lower(occupied_interval)) AND isfinite(upper(occupied_interval))
    AND lower_inc(occupied_interval) AND NOT upper_inc(occupied_interval)),
  EXCLUDE USING gist (organization_id WITH =, resource_id WITH =, occupied_interval WITH &&)
    WHERE (state IN ('held','committed'))
);
CREATE UNIQUE INDEX one_active_booking_resource ON ubi_v0.resource_reservations (organization_id, booking_id, resource_id) WHERE state = 'committed';
CREATE UNIQUE INDEX one_active_hold_resource ON ubi_v0.resource_reservations (organization_id, hold_id, resource_id) WHERE state = 'held';
CREATE TABLE ubi_v0.quotations (
  organization_id text NOT NULL, id text NOT NULL, operation_id text NOT NULL, service_id text NOT NULL,
  service_version integer NOT NULL CHECK (service_version > 0), request_id text NOT NULL,
  start_at timestamptz NOT NULL, end_at timestamptz NOT NULL,
  price_minor bigint NOT NULL CHECK (price_minor BETWEEN 0 AND 10000000), currency text NOT NULL CHECK (currency = 'ETB'),
  customer_label text NOT NULL CHECK (customer_label IN ('Demo Guest 01', 'Demo Guest 02', 'Demo Guest 03')),
  source text NOT NULL CHECK (source IN ('Call','Walk-in','WhatsApp','Telegram','Other')),
  state text NOT NULL CHECK (state IN ('issued','superseded','accepted')),
  version integer NOT NULL CHECK (version > 0), expires_at timestamptz NOT NULL, booking_id text,
  PRIMARY KEY (organization_id, id), UNIQUE (organization_id, request_id, version),
  FOREIGN KEY (organization_id, operation_id, service_id) REFERENCES ubi_v0.services,
  FOREIGN KEY (organization_id, booking_id) REFERENCES ubi_v0.bookings,
  CHECK ((state = 'accepted') = (booking_id IS NOT NULL)),
  CHECK (end_at > start_at AND isfinite(start_at) AND isfinite(end_at) AND isfinite(expires_at))
);
CREATE UNIQUE INDEX one_quote_conversion ON ubi_v0.quotations (organization_id, request_id) WHERE state = 'accepted';
CREATE TABLE ubi_v0.commands (
  organization_id text NOT NULL, principal_id text NOT NULL, action text NOT NULL, key text NOT NULL,
  request_hash text NOT NULL, result jsonb,
  PRIMARY KEY (organization_id, principal_id, action, key),
  CHECK (length(request_hash) = 64), CHECK (result IS NULL OR jsonb_typeof(result) = 'object')
);
CREATE TABLE ubi_v0.outbox (
  organization_id text NOT NULL, id text NOT NULL, event_type text NOT NULL, aggregate_id text NOT NULL,
  payload jsonb NOT NULL CHECK (jsonb_typeof(payload) = 'object'),
  PRIMARY KEY (organization_id, id)
);
CREATE TABLE ubi_v0.audit (
  organization_id text NOT NULL, id text NOT NULL, action text NOT NULL, aggregate_id text NOT NULL,
  principal_id text NOT NULL, before_data jsonb, after_data jsonb NOT NULL,
  PRIMARY KEY (organization_id, id), CHECK (jsonb_typeof(after_data) = 'object')
);

-- Every tenant table fails closed without transaction-local context.
DO $$ DECLARE tab text; BEGIN
  FOREACH tab IN ARRAY ARRAY['organizations','operations','resources','services','bookings','holds','resource_reservations','quotations','commands','outbox','audit'] LOOP
    EXECUTE format('ALTER TABLE ubi_v0.%I ENABLE ROW LEVEL SECURITY', tab);
    EXECUTE format('ALTER TABLE ubi_v0.%I FORCE ROW LEVEL SECURITY', tab);
    EXECUTE format('CREATE POLICY tenant_scope ON ubi_v0.%I USING (organization_id = current_setting(''ubi_v0.organization_id'', true)) WITH CHECK (organization_id = current_setting(''ubi_v0.organization_id'', true))', tab);
    EXECUTE format('GRANT SELECT ON ubi_v0.%I TO ubi_v0_runtime', tab);
  END LOOP;
END $$;
GRANT INSERT, UPDATE ON ubi_v0.bookings, ubi_v0.holds, ubi_v0.resource_reservations, ubi_v0.quotations, ubi_v0.commands TO ubi_v0_runtime;
GRANT INSERT ON ubi_v0.audit, ubi_v0.outbox TO ubi_v0_runtime;

CREATE FUNCTION ubi_v0.protect_keys() RETURNS trigger LANGUAGE plpgsql SET search_path = pg_catalog AS $$
BEGIN
  IF OLD.organization_id <> NEW.organization_id THEN RAISE EXCEPTION 'immutable tenant' USING ERRCODE = '23514'; END IF;
  IF to_jsonb(OLD) ? 'id' AND to_jsonb(OLD)->'id' <> to_jsonb(NEW)->'id' THEN RAISE EXCEPTION 'immutable id' USING ERRCODE = '23514'; END IF;
  RETURN NEW;
END $$;
DO $$ DECLARE tab text; BEGIN
  FOREACH tab IN ARRAY ARRAY['bookings','holds','resource_reservations','quotations','commands'] LOOP
    EXECUTE format('CREATE TRIGGER protect_keys BEFORE UPDATE ON ubi_v0.%I FOR EACH ROW EXECUTE FUNCTION ubi_v0.protect_keys()', tab);
  END LOOP;
END $$;
CREATE FUNCTION ubi_v0.protect_quote_terms() RETURNS trigger LANGUAGE plpgsql SET search_path = pg_catalog AS $$
BEGIN
  IF (to_jsonb(OLD) - ARRAY['state','booking_id']) <> (to_jsonb(NEW) - ARRAY['state','booking_id']) THEN
    RAISE EXCEPTION 'immutable quotation terms' USING ERRCODE = '23514';
  END IF;
  RETURN NEW;
END $$;
CREATE TRIGGER immutable_quote_terms BEFORE UPDATE ON ubi_v0.quotations FOR EACH ROW EXECUTE FUNCTION ubi_v0.protect_quote_terms();
CREATE FUNCTION ubi_v0.protect_command() RETURNS trigger LANGUAGE plpgsql SET search_path = pg_catalog AS $$
BEGIN
  IF (to_jsonb(OLD) - 'result') <> (to_jsonb(NEW) - 'result') OR OLD.result IS NOT NULL THEN
    RAISE EXCEPTION 'immutable successful command' USING ERRCODE = '23514';
  END IF;
  RETURN NEW;
END $$;
CREATE TRIGGER immutable_command BEFORE UPDATE ON ubi_v0.commands FOR EACH ROW EXECUTE FUNCTION ubi_v0.protect_command();
