package integrity

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct{ Pool *pgxpool.Pool }
type serviceSnapshot struct {
	ResourceID                       string
	Version, Duration, Before, After int
	Price                            int64
}

func translate(err error) error {
	if err == nil {
		return nil
	}
	var d *DomainError
	if errors.As(err, &d) {
		return err
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return failure("not_found", 404, "That synthetic item is not available in this scope.")
	}
	var pg *pgconn.PgError
	if errors.As(err, &pg) && pg.Code == "23P01" {
		return failure("capacity_conflict", 409, "That demo time is occupied. Keep your details and choose another time.")
	}
	if errors.As(err, &pg) && (pg.Code == "40P01" || pg.Code == "40001") {
		return &DomainError{"retry_transaction", "Retry the same command key; no new outcome is acknowledged.", 409, true}
	}
	return &DomainError{"storage_failure", "The synthetic transaction could not finish. Retry the same key to recover a committed outcome.", 500, true}
}
func (s *Service) begin(ctx context.Context, p Principal) (pgx.Tx, error) {
	if !p.Active {
		return nil, failure("permission_denied", 403, "This demo principal is inactive.")
	}
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(ctx, "SELECT set_config('ubi_v0.organization_id',$1,true)", p.OrganizationID)
	if err != nil {
		_ = tx.Rollback(context.Background())
		return nil, err
	}
	return tx, nil
}
func lock(ctx context.Context, tx pgx.Tx, value string) error {
	_, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock(hashtextextended($1,0))", value)
	return err
}

// Only successful committed effects are recorded. Lost responses replay exactly.
func (s *Service) command(ctx context.Context, p Principal, action, key string, input any, fn func(pgx.Tx) (Receipt, error)) (Receipt, error) {
	if !keyPattern.MatchString(key) {
		return Receipt{}, failure("invalid_idempotency_key", 400, "Supply a stable 8–128 character demo command key.")
	}
	sum := sha256.Sum256(jsonBytes(input))
	hash := hex.EncodeToString(sum[:])
	tx, err := s.begin(ctx, p)
	if err != nil {
		return Receipt{}, translate(err)
	}
	defer tx.Rollback(context.Background())
	_, err = tx.Exec(ctx, "INSERT INTO ubi_v0.commands(organization_id,principal_id,action,key,request_hash) VALUES($1,$2,$3,$4,$5) ON CONFLICT DO NOTHING", p.OrganizationID, p.ID, action, key, hash)
	if err != nil {
		return Receipt{}, translate(err)
	}
	var oldHash string
	var stored []byte
	err = tx.QueryRow(ctx, "SELECT request_hash,result FROM ubi_v0.commands WHERE organization_id=$1 AND principal_id=$2 AND action=$3 AND key=$4 FOR UPDATE", p.OrganizationID, p.ID, action, key).Scan(&oldHash, &stored)
	if err != nil {
		return Receipt{}, translate(err)
	}
	if oldHash != hash {
		return Receipt{}, failure("idempotency_mismatch", 409, "This successful command key belongs to different details.")
	}
	if stored != nil {
		var r Receipt
		if err = json.Unmarshal(stored, &r); err != nil {
			return Receipt{}, translate(err)
		}
		if err = tx.Commit(ctx); err != nil {
			return Receipt{}, translate(err)
		}
		return r, nil
	}
	r, err := fn(tx)
	if err != nil {
		return Receipt{}, translate(err)
	}
	_, err = tx.Exec(ctx, "UPDATE ubi_v0.commands SET result=$5 WHERE organization_id=$1 AND principal_id=$2 AND action=$3 AND key=$4", p.OrganizationID, p.ID, action, key, jsonBytes(r))
	if err != nil {
		return Receipt{}, translate(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return Receipt{}, translate(err)
	}
	return r, nil
}

func snapshot(ctx context.Context, tx pgx.Tx, in Intent) (serviceSnapshot, error) {
	var v serviceSnapshot
	err := tx.QueryRow(ctx, "SELECT resource_id,version,duration_minutes,buffer_before,buffer_after,price_minor FROM ubi_v0.services WHERE organization_id=$1 AND operation_id=$2 AND service_id=$3", in.OrganizationID, in.OperationID, in.ServiceID).Scan(&v.ResourceID, &v.Version, &v.Duration, &v.Before, &v.After, &v.Price)
	if err != nil {
		return v, err
	}
	if v.Version != in.ExpectedServiceVersion {
		return v, failure("service_version_conflict", 409, "The demo service changed. Review its new terms before retrying.")
	}
	return v, nil
}
func fromIntent(in Intent, v serviceSnapshot, state string) Result {
	return Result{DemoOnly: true, ID: newID(), Scope: in.Scope, Version: 1, State: state, StartAt: in.StartAt.UTC(), EndAt: in.StartAt.UTC().Add(time.Duration(v.Duration) * time.Minute), ServiceID: in.ServiceID, ServiceVersion: v.Version, PriceMinor: v.Price, Currency: "ETB", CustomerLabel: in.CustomerLabel, Source: in.Source}
}

const columns = "organization_id,id,operation_id,service_id,service_version,start_at,end_at,price_minor,currency,customer_label,source,state,version"

func resultRow(ctx context.Context, tx pgx.Tx, p Principal, table, id string, forUpdate bool) (Result, error) {
	// table is an internal closed enum, never client-provided SQL.
	extra := ""
	switch table {
	case "bookings":
	case "holds":
		extra = ",expires_at"
	case "quotations":
		extra = ",expires_at,request_id,COALESCE(booking_id,'')"
	default:
		panic("invalid table enum")
	}
	suffix := ""
	if forUpdate {
		suffix = " FOR UPDATE"
	}
	r := Result{DemoOnly: true}
	args := []any{&r.OrganizationID, &r.ID, &r.OperationID, &r.ServiceID, &r.ServiceVersion, &r.StartAt, &r.EndAt, &r.PriceMinor, &r.Currency, &r.CustomerLabel, &r.Source, &r.State, &r.Version}
	if table != "bookings" {
		args = append(args, &r.ExpiresAt)
	}
	if table == "quotations" {
		args = append(args, &r.RequestID, &r.BookingID)
	}
	err := tx.QueryRow(ctx, "SELECT "+columns+extra+" FROM ubi_v0."+table+" WHERE organization_id=$1 AND id=$2"+suffix, p.OrganizationID, id).Scan(args...)
	if err != nil {
		return r, err
	}
	if !p.Operations[r.OperationID] {
		return r, pgx.ErrNoRows
	}
	r.StartAt = r.StartAt.UTC()
	r.EndAt = r.EndAt.UTC()
	return r, nil
}
func insertBooking(ctx context.Context, tx pgx.Tx, r Result) error {
	_, err := tx.Exec(ctx, "INSERT INTO ubi_v0.bookings("+columns+") VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)", r.OrganizationID, r.ID, r.OperationID, r.ServiceID, r.ServiceVersion, r.StartAt, r.EndAt, r.PriceMinor, r.Currency, r.CustomerLabel, r.Source, r.State, r.Version)
	return err
}
func record(ctx context.Context, tx pgx.Tx, p Principal, action string, before *Result, after Result) error {
	var previous any
	if before != nil {
		previous = jsonBytes(before)
	}
	_, err := tx.Exec(ctx, "INSERT INTO ubi_v0.audit(organization_id,id,action,aggregate_id,principal_id,before_data,after_data) VALUES($1,$2,$3,$4,$5,$6,$7)", p.OrganizationID, newID(), action, after.ID, p.ID, previous, jsonBytes(after))
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "INSERT INTO ubi_v0.outbox(organization_id,id,event_type,aggregate_id,payload) VALUES($1,$2,$3,$4,$5)", p.OrganizationID, newID(), action, after.ID, jsonBytes(after))
	return err
}
func allocate(ctx context.Context, tx pgx.Tx, r Result, v serviceSnapshot, hold bool) error {
	var holdID, bookingID any
	state := "committed"
	if hold {
		holdID = r.ID
		state = "held"
	} else {
		bookingID = r.ID
	}
	_, err := tx.Exec(ctx, "INSERT INTO ubi_v0.resource_reservations(organization_id,id,resource_id,occupied_interval,state,hold_id,booking_id) VALUES($1,$2,$3,tstzrange($4,$5,'[)'),$6,$7,$8)", r.OrganizationID, newID(), v.ResourceID, r.StartAt.Add(-time.Duration(v.Before)*time.Minute), r.EndAt.Add(time.Duration(v.After)*time.Minute), state, holdID, bookingID)
	return err
}
func expireHolds(ctx context.Context, tx pgx.Tx, p Principal, resource string) error {
	rows, err := tx.Query(ctx, "SELECT h.id FROM ubi_v0.holds h JOIN ubi_v0.resource_reservations r ON r.organization_id=h.organization_id AND r.hold_id=h.id WHERE h.organization_id=$1 AND r.resource_id=$2 AND r.state='held' AND h.state='active' AND h.expires_at<=clock_timestamp()", p.OrganizationID, resource)
	if err != nil {
		return err
	}
	var ids []string
	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			rows.Close()
			return err
		}
		ids = append(ids, id)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return err
	}
	// Resource reconciliation is tenant-wide, including another operation sharing a resource.
	tenantPrincipal := p
	tenantPrincipal.Operations = map[string]bool{"studio": true, "events": true}
	for _, id := range ids {
		old, err := resultRow(ctx, tx, tenantPrincipal, "holds", id, true)
		if err != nil {
			return err
		}
		next := old
		next.State = "expired"
		next.Version++
		if _, err = tx.Exec(ctx, "UPDATE ubi_v0.holds SET state='expired',version=version+1 WHERE organization_id=$1 AND id=$2", p.OrganizationID, id); err != nil {
			return err
		}
		if _, err = tx.Exec(ctx, "UPDATE ubi_v0.resource_reservations SET state='released' WHERE organization_id=$1 AND hold_id=$2 AND state='held'", p.OrganizationID, id); err != nil {
			return err
		}
		if err = record(ctx, tx, p, "hold_expired", &old, next); err != nil {
			return err
		}
	}
	return nil
}
func (s *Service) CreateBooking(ctx context.Context, p Principal, key string, in Intent) (Receipt, error) {
	if err := authorize(p, in.Scope, false); err != nil {
		return Receipt{}, err
	}
	if err := in.validate(); err != nil {
		return Receipt{}, err
	}
	in.StartAt = in.StartAt.UTC()
	return s.command(ctx, p, "create_booking", key, in, func(tx pgx.Tx) (Receipt, error) {
		v, err := snapshot(ctx, tx, in)
		if err != nil {
			return Receipt{}, err
		}
		if err = lock(ctx, tx, p.OrganizationID+":resource:"+v.ResourceID); err != nil {
			return Receipt{}, err
		}
		if err = expireHolds(ctx, tx, p, v.ResourceID); err != nil {
			return Receipt{}, err
		}
		r := fromIntent(in, v, "confirmed_demo")
		if err = insertBooking(ctx, tx, r); err != nil {
			return Receipt{}, err
		}
		if err = allocate(ctx, tx, r, v, false); err != nil {
			return Receipt{}, err
		}
		if err = record(ctx, tx, p, "booking_created", nil, r); err != nil {
			return Receipt{}, err
		}
		return Receipt{201, r}, nil
	})
}
func (s *Service) CreateHold(ctx context.Context, p Principal, key string, in HoldInput) (Receipt, error) {
	if err := authorize(p, in.Scope, true); err != nil {
		return Receipt{}, err
	}
	if err := in.Intent.validate(); err != nil {
		return Receipt{}, err
	}
	if in.TTLSeconds < 1 || in.TTLSeconds > 900 {
		return Receipt{}, failure("invalid_input", 422, "Demo hold TTL must be 1–900 seconds.")
	}
	in.StartAt = in.StartAt.UTC()
	return s.command(ctx, p, "create_hold", key, in, func(tx pgx.Tx) (Receipt, error) {
		v, err := snapshot(ctx, tx, in.Intent)
		if err != nil {
			return Receipt{}, err
		}
		if err = lock(ctx, tx, p.OrganizationID+":resource:"+v.ResourceID); err != nil {
			return Receipt{}, err
		}
		if err = expireHolds(ctx, tx, p, v.ResourceID); err != nil {
			return Receipt{}, err
		}
		r := fromIntent(in.Intent, v, "active")
		var expires time.Time
		err = tx.QueryRow(ctx, "INSERT INTO ubi_v0.holds("+columns+",expires_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,clock_timestamp()+$14*interval '1 second') RETURNING expires_at", r.OrganizationID, r.ID, r.OperationID, r.ServiceID, r.ServiceVersion, r.StartAt, r.EndAt, r.PriceMinor, r.Currency, r.CustomerLabel, r.Source, r.State, r.Version, in.TTLSeconds).Scan(&expires)
		if err != nil {
			return Receipt{}, err
		}
		r.ExpiresAt = &expires
		if err = allocate(ctx, tx, r, v, true); err != nil {
			return Receipt{}, err
		}
		if err = record(ctx, tx, p, "hold_created", nil, r); err != nil {
			return Receipt{}, err
		}
		return Receipt{201, r}, nil
	})
}
func reservation(ctx context.Context, tx pgx.Tx, p Principal, id string, held bool) (string, string, time.Time, time.Time, error) {
	column, state := "booking_id", "committed"
	if held {
		column, state = "hold_id", "held"
	}
	var rid, resource string
	var lo, hi time.Time
	err := tx.QueryRow(ctx, "SELECT id,resource_id,lower(occupied_interval),upper(occupied_interval) FROM ubi_v0.resource_reservations WHERE organization_id=$1 AND "+column+"=$2 AND state=$3", p.OrganizationID, id, state).Scan(&rid, &resource, &lo, &hi)
	return rid, resource, lo, hi, err
}
func (s *Service) ConsumeHold(ctx context.Context, p Principal, key, id string, in VersionInput) (Receipt, error) {
	if !p.Owner {
		return Receipt{}, failure("permission_denied", 403, "Only an invented owner may convert a hold.")
	}
	if in.ExpectedVersion < 1 {
		return Receipt{}, failure("invalid_input", 422, "Provide the reviewed hold version.")
	}
	return s.command(ctx, p, "consume_hold", key, struct {
		ID    string
		Input VersionInput
	}{id, in}, func(tx pgx.Tx) (Receipt, error) {
		pre, err := resultRow(ctx, tx, p, "holds", id, false)
		if err != nil {
			return Receipt{}, err
		}
		if pre.State != "active" || pre.Version != in.ExpectedVersion {
			return Receipt{}, failure("hold_not_valid", 409, "This demo hold changed, expired or was consumed.")
		}
		rid, resource, _, _, err := reservation(ctx, tx, p, id, true)
		if err != nil {
			return Receipt{}, err
		}
		if err = lock(ctx, tx, p.OrganizationID+":resource:"+resource); err != nil {
			return Receipt{}, err
		}
		old, err := resultRow(ctx, tx, p, "holds", id, true)
		if err != nil {
			return Receipt{}, err
		}
		var live bool
		if err = tx.QueryRow(ctx, "SELECT expires_at>clock_timestamp() FROM ubi_v0.holds WHERE organization_id=$1 AND id=$2", p.OrganizationID, id).Scan(&live); err != nil {
			return Receipt{}, err
		}
		if old.State != "active" || old.Version != in.ExpectedVersion || !live {
			return Receipt{}, failure("hold_not_valid", 409, "This demo hold expired, changed or was consumed. Nothing is newly confirmed.")
		}
		r := old
		r.ID = newID()
		r.State = "confirmed_demo"
		r.Version = 1
		r.ExpiresAt = nil
		if err = insertBooking(ctx, tx, r); err != nil {
			return Receipt{}, err
		}
		tag, err := tx.Exec(ctx, "UPDATE ubi_v0.resource_reservations SET state='committed',hold_id=NULL,booking_id=$3 WHERE organization_id=$1 AND id=$2 AND state='held' AND hold_id=$4", p.OrganizationID, rid, r.ID, id)
		if err != nil {
			return Receipt{}, err
		}
		if tag.RowsAffected() != 1 {
			return Receipt{}, failure("hold_not_valid", 409, "This allocation is no longer held.")
		}
		if _, err = tx.Exec(ctx, "UPDATE ubi_v0.holds SET state='consumed',version=version+1 WHERE organization_id=$1 AND id=$2", p.OrganizationID, id); err != nil {
			return Receipt{}, err
		}
		next := old
		next.State = "consumed"
		next.Version++
		if err = record(ctx, tx, p, "hold_consumed", &old, next); err != nil {
			return Receipt{}, err
		}
		if err = record(ctx, tx, p, "booking_created", nil, r); err != nil {
			return Receipt{}, err
		}
		return Receipt{201, r}, nil
	})
}
func (s *Service) GetBooking(ctx context.Context, p Principal, id string) (Result, error) {
	tx, err := s.begin(ctx, p)
	if err != nil {
		return Result{}, translate(err)
	}
	defer tx.Rollback(context.Background())
	r, err := resultRow(ctx, tx, p, "bookings", id, false)
	if err != nil {
		return r, translate(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return r, translate(err)
	}
	return r, nil
}
func (s *Service) Reschedule(ctx context.Context, p Principal, key, id string, in RescheduleInput) (Receipt, error) {
	if in.ExpectedVersion < 1 || !validTime(in.StartAt) {
		return Receipt{}, failure("invalid_input", 422, "Review the version and explicit-offset new demo time.")
	}
	in.StartAt = in.StartAt.UTC()
	return s.command(ctx, p, "reschedule", key, struct {
		ID    string
		Input RescheduleInput
	}{id, in}, func(tx pgx.Tx) (Receipt, error) {
		pre, err := resultRow(ctx, tx, p, "bookings", id, false)
		if err != nil {
			return Receipt{}, err
		}
		if pre.State != "confirmed_demo" || pre.Version != in.ExpectedVersion {
			return Receipt{}, failure("booking_version_conflict", 409, "Reload the changed or cancelled demo booking before editing.")
		}
		rid, resource, lo, hi, err := reservation(ctx, tx, p, id, false)
		if err != nil {
			return Receipt{}, err
		}
		if err = lock(ctx, tx, p.OrganizationID+":resource:"+resource); err != nil {
			return Receipt{}, err
		}
		old, err := resultRow(ctx, tx, p, "bookings", id, true)
		if err != nil {
			return Receipt{}, err
		}
		if old.State != "confirmed_demo" || old.Version != in.ExpectedVersion {
			return Receipt{}, failure("booking_version_conflict", 409, "The demo booking changed or was cancelled. Reload before editing.")
		}
		// Re-read allocation after acquiring the lock; never calculate from a stale range.
		rid, _, lo, hi, err = reservation(ctx, tx, p, id, false)
		if err != nil {
			return Receipt{}, err
		}
		if err = expireHolds(ctx, tx, p, resource); err != nil {
			return Receipt{}, err
		}
		r := old
		r.StartAt = in.StartAt
		r.EndAt = in.StartAt.Add(old.EndAt.Sub(old.StartAt))
		r.Version++
		if _, err = tx.Exec(ctx, "UPDATE ubi_v0.resource_reservations SET occupied_interval=tstzrange($3,$4,'[)') WHERE organization_id=$1 AND id=$2", p.OrganizationID, rid, r.StartAt.Add(-old.StartAt.Sub(lo)), r.EndAt.Add(hi.Sub(old.EndAt))); err != nil {
			return Receipt{}, err
		}
		if _, err = tx.Exec(ctx, "UPDATE ubi_v0.bookings SET start_at=$3,end_at=$4,version=$5 WHERE organization_id=$1 AND id=$2", p.OrganizationID, id, r.StartAt, r.EndAt, r.Version); err != nil {
			return Receipt{}, err
		}
		if err = record(ctx, tx, p, "booking_rescheduled", &old, r); err != nil {
			return Receipt{}, err
		}
		return Receipt{200, r}, nil
	})
}
func (s *Service) Cancel(ctx context.Context, p Principal, key, id string, in VersionInput) (Receipt, error) {
	if in.ExpectedVersion < 1 {
		return Receipt{}, failure("invalid_input", 422, "Provide the reviewed booking version.")
	}
	return s.command(ctx, p, "cancel", key, struct {
		ID    string
		Input VersionInput
	}{id, in}, func(tx pgx.Tx) (Receipt, error) {
		pre, err := resultRow(ctx, tx, p, "bookings", id, false)
		if err != nil {
			return Receipt{}, err
		}
		if pre.State != "confirmed_demo" || pre.Version != in.ExpectedVersion {
			return Receipt{}, failure("booking_version_conflict", 409, "Reload the changed or cancelled demo booking before acting.")
		}
		_, resource, _, _, err := reservation(ctx, tx, p, id, false)
		if err != nil {
			return Receipt{}, err
		}
		if err = lock(ctx, tx, p.OrganizationID+":resource:"+resource); err != nil {
			return Receipt{}, err
		}
		old, err := resultRow(ctx, tx, p, "bookings", id, true)
		if err != nil {
			return Receipt{}, err
		}
		if old.State != "confirmed_demo" || old.Version != in.ExpectedVersion {
			return Receipt{}, failure("booking_version_conflict", 409, "Reload the changed or cancelled demo booking before acting.")
		}
		r := old
		r.State = "cancelled_demo"
		r.Version++
		if _, err = tx.Exec(ctx, "UPDATE ubi_v0.bookings SET state=$3,version=$4 WHERE organization_id=$1 AND id=$2", p.OrganizationID, id, r.State, r.Version); err != nil {
			return Receipt{}, err
		}
		if _, err = tx.Exec(ctx, "UPDATE ubi_v0.resource_reservations SET state='released' WHERE organization_id=$1 AND booking_id=$2 AND state='committed'", p.OrganizationID, id); err != nil {
			return Receipt{}, err
		}
		if err = record(ctx, tx, p, "booking_cancelled", &old, r); err != nil {
			return Receipt{}, err
		}
		return Receipt{200, r}, nil
	})
}
func (s *Service) IssueQuote(ctx context.Context, p Principal, key string, in QuoteInput) (Receipt, error) {
	if err := authorize(p, in.Scope, true); err != nil {
		return Receipt{}, err
	}
	if err := in.Intent.validate(); err != nil {
		return Receipt{}, err
	}
	if !requestPattern.MatchString(in.RequestID) || (in.DemoPriceMinor != nil && (*in.DemoPriceMinor < 0 || *in.DemoPriceMinor > 10000000)) {
		return Receipt{}, failure("invalid_input", 422, "Use a dummy request ID and bounded integer demo price.")
	}
	in.StartAt = in.StartAt.UTC()
	return s.command(ctx, p, "issue_quote", key, in, func(tx pgx.Tx) (Receipt, error) {
		if err := lock(ctx, tx, p.OrganizationID+":request:"+in.RequestID); err != nil {
			return Receipt{}, err
		}
		v, err := snapshot(ctx, tx, in.Intent)
		if err != nil {
			return Receipt{}, err
		}
		var version, accepted int
		var existingOperation *string
		err = tx.QueryRow(ctx, "SELECT COALESCE(max(version),0),count(*) FILTER(WHERE state='accepted'),min(operation_id) FROM ubi_v0.quotations WHERE organization_id=$1 AND request_id=$2", p.OrganizationID, in.RequestID).Scan(&version, &accepted, &existingOperation)
		if err != nil {
			return Receipt{}, err
		}
		if accepted > 0 {
			return Receipt{}, failure("quote_already_converted", 409, "This demo request already became a booking.")
		}
		if existingOperation != nil && *existingOperation != in.OperationID {
			return Receipt{}, failure("request_scope_conflict", 409, "This demo request belongs to another operation.")
		}
		if _, err = tx.Exec(ctx, "UPDATE ubi_v0.quotations SET state='superseded' WHERE organization_id=$1 AND request_id=$2 AND state='issued'", p.OrganizationID, in.RequestID); err != nil {
			return Receipt{}, err
		}
		r := fromIntent(in.Intent, v, "issued")
		r.Version = version + 1
		r.RequestID = in.RequestID
		if in.DemoPriceMinor != nil {
			r.PriceMinor = *in.DemoPriceMinor
		}
		var expires time.Time
		err = tx.QueryRow(ctx, "INSERT INTO ubi_v0.quotations("+columns+",request_id,expires_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,clock_timestamp()+interval '15 minutes') RETURNING expires_at", r.OrganizationID, r.ID, r.OperationID, r.ServiceID, r.ServiceVersion, r.StartAt, r.EndAt, r.PriceMinor, r.Currency, r.CustomerLabel, r.Source, r.State, r.Version, r.RequestID).Scan(&expires)
		if err != nil {
			return Receipt{}, err
		}
		r.ExpiresAt = &expires
		if err = record(ctx, tx, p, "quote_issued", nil, r); err != nil {
			return Receipt{}, err
		}
		return Receipt{201, r}, nil
	})
}
func (s *Service) AcceptQuote(ctx context.Context, p Principal, key, id string, in VersionInput) (Receipt, error) {
	if !p.Owner {
		return Receipt{}, failure("permission_denied", 403, "Only an invented owner may exercise this conversion experiment.")
	}
	if in.ExpectedVersion < 1 {
		return Receipt{}, failure("invalid_input", 422, "Provide the reviewed quotation version.")
	}
	return s.command(ctx, p, "accept_quote", key, struct {
		ID    string
		Input VersionInput
	}{id, in}, func(tx pgx.Tx) (Receipt, error) {
		pre, err := resultRow(ctx, tx, p, "quotations", id, false)
		if err != nil {
			return Receipt{}, err
		}
		if err = lock(ctx, tx, p.OrganizationID+":request:"+pre.RequestID); err != nil {
			return Receipt{}, err
		}
		old, err := resultRow(ctx, tx, p, "quotations", id, true)
		if err != nil {
			return Receipt{}, err
		}
		var live bool
		if err = tx.QueryRow(ctx, "SELECT expires_at>clock_timestamp() FROM ubi_v0.quotations WHERE organization_id=$1 AND id=$2", p.OrganizationID, id).Scan(&live); err != nil {
			return Receipt{}, err
		}
		if old.State != "issued" || old.Version != in.ExpectedVersion || !live {
			return Receipt{}, failure("quote_not_valid", 409, "This demo quotation changed, expired or was already accepted. Review the current offer.")
		}
		v, err := snapshot(ctx, tx, Intent{Scope: old.Scope, ServiceID: old.ServiceID, ExpectedServiceVersion: old.ServiceVersion})
		if err != nil {
			return Receipt{}, err
		}
		if err = lock(ctx, tx, p.OrganizationID+":resource:"+v.ResourceID); err != nil {
			return Receipt{}, err
		}
		// Waiting for capacity must not extend a quotation's expiry.
		if err = tx.QueryRow(ctx, "SELECT expires_at>clock_timestamp() FROM ubi_v0.quotations WHERE organization_id=$1 AND id=$2", p.OrganizationID, id).Scan(&live); err != nil {
			return Receipt{}, err
		}
		if !live {
			return Receipt{}, failure("quote_not_valid", 409, "This demo quotation expired while waiting. Nothing was converted.")
		}
		if err = expireHolds(ctx, tx, p, v.ResourceID); err != nil {
			return Receipt{}, err
		}
		r := old
		r.ID = newID()
		r.State = "confirmed_demo"
		r.Version = 1
		r.RequestID = ""
		r.ExpiresAt = nil
		r.BookingID = ""
		if err = insertBooking(ctx, tx, r); err != nil {
			return Receipt{}, err
		}
		if err = allocate(ctx, tx, r, v, false); err != nil {
			return Receipt{}, err
		}
		if _, err = tx.Exec(ctx, "UPDATE ubi_v0.quotations SET state='accepted',booking_id=$3 WHERE organization_id=$1 AND id=$2", p.OrganizationID, id, r.ID); err != nil {
			return Receipt{}, err
		}
		next := old
		next.State = "accepted"
		next.BookingID = r.ID
		if err = record(ctx, tx, p, "quote_accepted", &old, next); err != nil {
			return Receipt{}, err
		}
		if err = record(ctx, tx, p, "booking_created", nil, r); err != nil {
			return Receipt{}, err
		}
		return Receipt{201, r}, nil
	})
}
