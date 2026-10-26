package integrity

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var admin *pgxpool.Pool
var runtimePool *pgxpool.Pool
var privateCluster string
var testContext = context.Background()

func TestMain(m *testing.M) {
	if os.Getenv("UBI_RUN_SYNTHETIC_PG") != "1" {
		fmt.Fprintln(os.Stderr, "Refusing database tests without UBI_RUN_SYNTHETIC_PG=1. No existing database is accepted.")
		os.Exit(2)
	}
	code := func() int {
		dir, err := os.MkdirTemp("", "ubi-wb03-")
		if err != nil {
			panic(err)
		}
		privateCluster = dir
		// Only this newly created, owned temporary directory can be removed.
		serverRunning := false
		defer func() {
			removed := false
			if !serverRunning {
				removed = os.RemoveAll(dir) == nil
			} else {
				fmt.Fprintln(os.Stderr, "Private synthetic cluster retained because shutdown failed:", dir)
			}
			_ = os.MkdirAll("artifacts", 0755)
			_ = os.WriteFile("artifacts/cluster-receipt.json", jsonBytes(map[string]any{
				"synthetic_only": true, "new_private_temporary_cluster": true,
				"existing_database_accepted": false, "tcp_enabled": false,
				"server_stopped": !serverRunning, "temporary_files_removed": removed,
			}), 0644)
		}()
		bin := "/usr/lib/postgresql/18/bin"
		data := filepath.Join(dir, "data")
		run := func(name string, args ...string) error {
			out, e := exec.Command(filepath.Join(bin, name), args...).CombinedOutput()
			if e != nil {
				return fmt.Errorf("%s: %w: %s", name, e, out)
			}
			return nil
		}
		if err = run("initdb", "-D", data, "-U", "ubi_v0_admin", "--auth-local=trust", "--auth-host=reject", "--encoding=UTF8", "--no-locale"); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 2
		}
		serverRunning = true
		defer func() {
			if err := run("pg_ctl", "-D", data, "-m", "fast", "-w", "stop"); err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				serverRunning = false
			}
		}()
		if err = run("pg_ctl", "-D", data, "-l", filepath.Join(dir, "postgres.log"), "-o", fmt.Sprintf("-h '' -k %s -p 5432 -c max_connections=120", dir), "-w", "start"); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 2
		}
		dsn := fmt.Sprintf("host=%s port=5432 dbname=postgres user=ubi_v0_admin sslmode=disable timezone=UTC", dir)
		admin, err = pgxpool.New(testContext, dsn)
		if err != nil {
			panic(err)
		}
		defer admin.Close()
		sql, err := os.ReadFile("schema.sql")
		if err != nil {
			panic(err)
		}
		if _, err = admin.Exec(testContext, string(sql)); err != nil {
			fmt.Fprintln(os.Stderr, "Synthetic schema failed:", err)
			return 2
		}
		cfg, err := pgxpool.ParseConfig(strings.Replace(dsn, "user=ubi_v0_admin", "user=ubi_v0_runtime", 1))
		if err != nil {
			panic(err)
		}
		cfg.MaxConns = 32
		runtimePool, err = pgxpool.NewWithConfig(testContext, cfg)
		if err != nil {
			panic(err)
		}
		defer runtimePool.Close()
		if err = runtimePool.Ping(testContext); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 2
		}
		return m.Run()
	}()
	os.Exit(code)
}

func reset(t *testing.T) *Service {
	t.Helper()
	_, err := admin.Exec(testContext, "TRUNCATE ubi_v0.audit,ubi_v0.outbox,ubi_v0.commands,ubi_v0.quotations,ubi_v0.resource_reservations,ubi_v0.holds,ubi_v0.bookings,ubi_v0.services,ubi_v0.resources,ubi_v0.operations,ubi_v0.organizations")
	if err != nil {
		t.Fatal(err)
	}
	var f struct {
		SyntheticOnly    bool `json:"syntheticOnly"`
		CurrencyExponent int  `json:"currencyExponent"`
		Organizations    []struct {
			ID string `json:"id"`
		} `json:"organizations"`
		Operations []struct {
			OrganizationID string `json:"organizationId"`
			ID             string `json:"id"`
		} `json:"operations"`
		Services []struct {
			OrganizationID string `json:"organizationId"`
			OperationID    string `json:"operationId"`
			ID             string `json:"id"`
			ResourceID     string `json:"resourceId"`
			Duration       int    `json:"durationMinutes"`
			Price          int64  `json:"priceMinor"`
		} `json:"services"`
	}
	raw, err := os.ReadFile("../fixtures/booking-v0.json")
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(raw, &f); err != nil {
		t.Fatal(err)
	}
	if !f.SyntheticOnly || f.CurrencyExponent != 2 || len(f.Organizations) != 2 || len(f.Operations) != 3 {
		t.Fatal("invalid synthetic fixture boundary")
	}
	for _, o := range f.Organizations {
		if _, err = admin.Exec(testContext, "INSERT INTO ubi_v0.organizations VALUES($1)", o.ID); err != nil {
			t.Fatal(err)
		}
	}
	for _, o := range f.Operations {
		if _, err = admin.Exec(testContext, "INSERT INTO ubi_v0.operations VALUES($1,$2)", o.OrganizationID, o.ID); err != nil {
			t.Fatal(err)
		}
	}
	for _, s := range f.Services {
		if _, err = admin.Exec(testContext, "INSERT INTO ubi_v0.resources VALUES($1,$2) ON CONFLICT DO NOTHING", s.OrganizationID, s.ResourceID); err != nil {
			t.Fatal(err)
		}
		if _, err = admin.Exec(testContext, "INSERT INTO ubi_v0.services(organization_id,operation_id,service_id,resource_id,version,duration_minutes,price_minor,currency) VALUES($1,$2,$3,$4,1,$5,$6,'ETB')", s.OrganizationID, s.OperationID, s.ID, s.ResourceID, s.Duration, s.Price); err != nil {
			t.Fatal(err)
		}
	}
	return &Service{runtimePool}
}
func ownerA() Principal { p, _ := DemoPrincipal("demo-owner-a"); return p }
func ownerB() Principal { p, _ := DemoPrincipal("demo-owner-b"); return p }
func intentAt(hour int) Intent {
	return Intent{Scope: Scope{"demo-org-a", "studio"}, ServiceID: "portrait", ExpectedServiceVersion: 1, StartAt: time.Date(2030, 1, 15, hour, 0, 0, 0, time.FixedZone("Addis fixture", 10800)), CustomerLabel: "Demo Guest 01", Source: "Call"}
}
func success(t *testing.T, r Receipt, e error) Receipt {
	t.Helper()
	if e != nil {
		t.Fatal(e)
	}
	if !r.Result.DemoOnly {
		t.Fatal("missing synthetic marker")
	}
	return r
}
func requireError(t *testing.T, err error, code string) {
	t.Helper()
	var d *DomainError
	if !errors.As(err, &d) || d.Code != code {
		t.Fatalf("wanted %s, got %v", code, err)
	}
}
func count(t *testing.T, table string) int {
	t.Helper()
	var n int
	if err := admin.QueryRow(testContext, "SELECT count(*) FROM ubi_v0."+table).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}
func create(t *testing.T, s *Service, key string, hour int) Receipt {
	t.Helper()
	r, e := s.CreateBooking(testContext, ownerA(), key, intentAt(hour))
	return success(t, r, e)
}

func TestRuntimeRoleAndSyntheticBoundary(t *testing.T) {
	reset(t)
	var super, bypass, owns bool
	err := runtimePool.QueryRow(testContext, "SELECT rolsuper,rolbypassrls,EXISTS(SELECT 1 FROM pg_tables WHERE schemaname='ubi_v0' AND tableowner=current_user) FROM pg_roles WHERE rolname=current_user").Scan(&super, &bypass, &owns)
	if err != nil {
		t.Fatal(err)
	}
	if super || bypass || owns {
		t.Fatal("runtime role can bypass tenant enforcement")
	}
	var hosts string
	if err = admin.QueryRow(testContext, "SHOW listen_addresses").Scan(&hosts); err != nil || hosts != "" {
		t.Fatal("PostgreSQL must have no TCP listener")
	}
	if !strings.HasPrefix(filepath.Base(privateCluster), "ubi-wb03-") {
		t.Fatal("not a disposable owned cluster")
	}
}
func TestAtomicBookingRetryAndLostOutcome(t *testing.T) {
	s := reset(t)
	first := create(t, s, "local-intent-0001", 9)
	r, e := s.CreateBooking(testContext, ownerA(), "local-intent-0001", intentAt(9))
	again := success(t, r, e)
	if !bytes.Equal(jsonBytes(first), jsonBytes(again)) {
		t.Fatal("replay changed first result")
	}
	changed := intentAt(10)
	_, e = s.CreateBooking(testContext, ownerA(), "local-intent-0001", changed)
	requireError(t, e, "idempotency_mismatch")
	for _, tab := range []string{"bookings", "resource_reservations", "commands", "audit", "outbox"} {
		if count(t, tab) != 1 {
			t.Fatalf("duplicate/missing %s", tab)
		}
	}
	if first.Result.State != "confirmed_demo" || first.Result.PriceMinor != 200000 {
		t.Fatal("snapshot or truth boundary wrong")
	}
}
func TestOneHundredConcurrentContenders(t *testing.T) {
	s := reset(t)
	ready := make(chan struct{})
	out := make(chan error, 100)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-ready
			_, err := s.CreateBooking(testContext, ownerA(), fmt.Sprintf("contender-%03d", i), intentAt(9))
			out <- err
		}(i)
	}
	close(ready)
	wg.Wait()
	close(out)
	wins, conflicts := 0, 0
	for err := range out {
		if err == nil {
			wins++
		} else {
			requireError(t, err, "capacity_conflict")
			conflicts++
		}
	}
	if wins != 1 || conflicts != 99 {
		t.Fatalf("wins=%d conflicts=%d", wins, conflicts)
	}
	for _, tab := range []string{"bookings", "resource_reservations", "commands", "audit", "outbox"} {
		if count(t, tab) != 1 {
			t.Fatalf("partial contention effect: %s", tab)
		}
	}
}
func TestConcurrentDuplicateCommandHasOneEffect(t *testing.T) {
	s := reset(t)
	gate := make(chan struct{})
	type answer struct {
		r Receipt
		e error
	}
	ch := make(chan answer, 12)
	for i := 0; i < 12; i++ {
		go func() {
			<-gate
			r, e := s.CreateBooking(testContext, ownerA(), "duplicate-0001", intentAt(9))
			ch <- answer{r, e}
		}()
	}
	close(gate)
	id := ""
	for i := 0; i < 12; i++ {
		a := <-ch
		r := success(t, a.r, a.e)
		if id == "" {
			id = r.Result.ID
		}
		if r.Result.ID != id {
			t.Fatal("duplicate effect")
		}
	}
	if count(t, "bookings") != 1 || count(t, "outbox") != 1 {
		t.Fatal("duplicate booking/notice work")
	}
}
func TestHeldAndCommittedCapacityShareOneBoundary(t *testing.T) {
	s := reset(t)
	r, e := s.CreateHold(testContext, ownerA(), "hold-create-001", HoldInput{intentAt(9), 300})
	h := success(t, r, e)
	_, e = s.CreateBooking(testContext, ownerA(), "book-blocked-001", intentAt(9))
	requireError(t, e, "capacity_conflict")
	var reservationID string
	if err := admin.QueryRow(testContext, "SELECT id FROM ubi_v0.resource_reservations WHERE hold_id=$1", h.Result.ID).Scan(&reservationID); err != nil {
		t.Fatal(err)
	}
	r, e = s.ConsumeHold(testContext, ownerA(), "consume-hold-001", h.Result.ID, VersionInput{1})
	b := success(t, r, e)
	var afterID, state string
	var oldHold *string
	if err := admin.QueryRow(testContext, "SELECT id,state,hold_id FROM ubi_v0.resource_reservations WHERE booking_id=$1", b.Result.ID).Scan(&afterID, &state, &oldHold); err != nil {
		t.Fatal(err)
	}
	if reservationID != afterID || state != "committed" || oldHold != nil || count(t, "resource_reservations") != 1 {
		t.Fatal("hold was released/reallocated instead of atomically converted")
	}
	r, e = s.ConsumeHold(testContext, ownerA(), "consume-hold-001", h.Result.ID, VersionInput{1})
	if success(t, r, e).Result.ID != b.Result.ID {
		t.Fatal("duplicate hold conversion")
	}
}
func TestExpiredHoldCannotConfirmAndReconciles(t *testing.T) {
	s := reset(t)
	r, e := s.CreateHold(testContext, ownerA(), "expiry-hold-001", HoldInput{intentAt(9), 300})
	h := success(t, r, e)
	if _, e = admin.Exec(testContext, "UPDATE ubi_v0.holds SET created_at=clock_timestamp()-interval '10 minutes',expires_at=clock_timestamp()-interval '1 minute' WHERE id=$1", h.Result.ID); e != nil {
		t.Fatal(e)
	}
	_, e = s.ConsumeHold(testContext, ownerA(), "expired-consume-001", h.Result.ID, VersionInput{1})
	requireError(t, e, "hold_not_valid")
	if count(t, "bookings") != 0 {
		t.Fatal("expired hold confirmed")
	}
	create(t, s, "after-expiry-001", 9)
	var state string
	if err := admin.QueryRow(testContext, "SELECT state FROM ubi_v0.holds WHERE id=$1", h.Result.ID).Scan(&state); err != nil || state != "expired" {
		t.Fatal("expired hold not reconciled")
	}
}
func TestAdjacencyAndBuffers(t *testing.T) {
	s := reset(t)
	create(t, s, "adjacent-first-001", 9)
	create(t, s, "adjacent-second-001", 10)
	s = reset(t)
	if _, err := admin.Exec(testContext, "UPDATE ubi_v0.services SET buffer_after=15 WHERE organization_id='demo-org-a' AND service_id='portrait'"); err != nil {
		t.Fatal(err)
	}
	create(t, s, "buffer-first-001", 9)
	_, e := s.CreateBooking(testContext, ownerA(), "buffer-next-001", intentAt(10))
	requireError(t, e, "capacity_conflict")
	in := intentAt(10)
	in.StartAt = in.StartAt.Add(15 * time.Minute)
	r, e := s.CreateBooking(testContext, ownerA(), "buffer-adjacent-001", in)
	success(t, r, e)
}
func TestTenantAndOperationIsolation(t *testing.T) {
	s := reset(t)
	a := create(t, s, "tenant-a-0001", 9)
	in := intentAt(9)
	in.OrganizationID = "demo-org-b"
	r, e := s.CreateBooking(testContext, ownerB(), "tenant-b-0001", in)
	b := success(t, r, e)
	if a.Result.ID == b.Result.ID || b.Result.PriceMinor != 120000 {
		t.Fatal("tenant snapshot mixed")
	}
	_, e = s.GetBooking(testContext, ownerB(), a.Result.ID)
	requireError(t, e, "not_found")
	_, e = s.CreateBooking(testContext, ownerA(), "wrong-tenant-0001", in)
	requireError(t, e, "permission_denied")
	event := intentAt(9)
	event.OperationID = "events"
	event.ServiceID = "consultation"
	r, e = s.CreateBooking(testContext, ownerA(), "event-operation-001", event)
	evt := success(t, r, e)
	staff, _ := DemoPrincipal("demo-staff-a")
	_, e = s.GetBooking(testContext, staff, evt.Result.ID)
	requireError(t, e, "not_found")
	_, e = s.CreateHold(testContext, staff, "staff-no-hold-001", HoldInput{intentAt(11), 300})
	requireError(t, e, "permission_denied")
	revoked, _ := DemoPrincipal("demo-revoked")
	_, e = s.CreateBooking(testContext, revoked, "revoked-demo-001", intentAt(11))
	requireError(t, e, "permission_denied")
}
func TestRLSFailsClosedAndContextDoesNotLeak(t *testing.T) {
	s := reset(t)
	create(t, s, "rls-seed-0001", 9)
	cfg := runtimePool.Config().Copy()
	cfg.MaxConns = 1
	one, err := pgxpool.NewWithConfig(testContext, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer one.Close()
	var n int
	if err = one.QueryRow(testContext, "SELECT count(*) FROM ubi_v0.bookings").Scan(&n); err != nil || n != 0 {
		t.Fatal("missing tenant context did not deny reads")
	}
	svc := &Service{one}
	tx, err := svc.begin(testContext, ownerA())
	if err != nil {
		t.Fatal(err)
	}
	if err = tx.QueryRow(testContext, "SELECT count(*) FROM ubi_v0.bookings").Scan(&n); err != nil || n != 1 {
		t.Fatal("own tenant not visible")
	}
	if err = tx.Commit(testContext); err != nil {
		t.Fatal(err)
	}
	if err = one.QueryRow(testContext, "SELECT count(*) FROM ubi_v0.bookings").Scan(&n); err != nil || n != 0 {
		t.Fatal("committed context leaked into next checkout")
	}
	tx, err = svc.begin(testContext, ownerB())
	if err != nil {
		t.Fatal(err)
	}
	if err = tx.QueryRow(testContext, "SELECT count(*) FROM ubi_v0.bookings").Scan(&n); err != nil || n != 0 {
		t.Fatal("unfiltered query saw foreign tenant")
	}
	_ = tx.Rollback(testContext)
	tx, err = svc.begin(testContext, ownerA())
	if err != nil {
		t.Fatal(err)
	}
	_ = tx.Rollback(testContext)
	if err = one.QueryRow(testContext, "SELECT count(*) FROM ubi_v0.bookings").Scan(&n); err != nil || n != 0 {
		t.Fatal("rolled-back context leaked")
	}
}
func TestCompositeKeysAndAppendOnlyPrivileges(t *testing.T) {
	s := reset(t)
	create(t, s, "constraints-seed-001", 9)
	tx, err := s.begin(testContext, ownerA())
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec(testContext, "UPDATE ubi_v0.bookings SET operation_id='events' WHERE organization_id='demo-org-a'")
	var pg *pgconn.PgError
	if !errors.As(err, &pg) || pg.Code != "23503" {
		t.Fatalf("missing composite FK: %v", err)
	}
	_ = tx.Rollback(testContext)
	for _, q := range []string{"UPDATE ubi_v0.audit SET action='rewritten'", "DELETE FROM ubi_v0.outbox", "ALTER TABLE ubi_v0.bookings DISABLE ROW LEVEL SECURITY", "CREATE TABLE ubi_v0.unapproved(id int)"} {
		tx, err = s.begin(testContext, ownerA())
		if err != nil {
			t.Fatal(err)
		}
		_, err = tx.Exec(testContext, q)
		if err == nil {
			t.Fatalf("runtime privilege too broad: %s", q)
		}
		_ = tx.Rollback(testContext)
	}
}
func TestAtomicRescheduleAndStaleEdit(t *testing.T) {
	s := reset(t)
	first := create(t, s, "move-first-0001", 9)
	create(t, s, "move-blocker-0001", 11)
	_, e := s.Reschedule(testContext, ownerA(), "move-failed-0001", first.Result.ID, RescheduleInput{VersionInput{1}, intentAt(11).StartAt})
	requireError(t, e, "capacity_conflict")
	old, e := s.GetBooking(testContext, ownerA(), first.Result.ID)
	if e != nil || !old.StartAt.Equal(first.Result.StartAt) || old.Version != 1 {
		t.Fatal("failed move lost old booking")
	}
	r, e := s.Reschedule(testContext, ownerA(), "move-success-0001", first.Result.ID, RescheduleInput{VersionInput{1}, intentAt(10).StartAt})
	next := success(t, r, e)
	if next.Result.ID != first.Result.ID || next.Result.Version != 2 {
		t.Fatal("move replaced booking identity")
	}
	_, e = s.Reschedule(testContext, ownerA(), "move-stale-0001", first.Result.ID, RescheduleInput{VersionInput{1}, intentAt(12).StartAt})
	requireError(t, e, "booking_version_conflict")
	if count(t, "bookings") != 2 || count(t, "outbox") != 3 || count(t, "audit") != 3 {
		t.Fatal("partial/duplicate move history")
	}
}
func TestCancellationReleasesCapacityAndStaleEditCannotResurrect(t *testing.T) {
	s := reset(t)
	first := create(t, s, "cancel-first-001", 9)
	r, e := s.Cancel(testContext, ownerA(), "cancel-command-001", first.Result.ID, VersionInput{1})
	next := success(t, r, e)
	if next.Result.State != "cancelled_demo" {
		t.Fatal("wrong cancellation")
	}
	_, e = s.Reschedule(testContext, ownerA(), "cancel-stale-move", first.Result.ID, RescheduleInput{VersionInput{1}, intentAt(10).StartAt})
	if e == nil {
		t.Fatal("cancelled booking resurrected")
	}
	create(t, s, "cancel-replacement", 9)
	if count(t, "bookings") != 2 || count(t, "audit") != 3 {
		t.Fatal("cancellation erased history")
	}
}
func TestForcedOutboxFailureRollsBackEveryEffect(t *testing.T) {
	s := reset(t)
	_, err := admin.Exec(testContext, "CREATE FUNCTION ubi_v0.reject_outbox() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'synthetic failure'; END $$; CREATE TRIGGER reject_outbox BEFORE INSERT ON ubi_v0.outbox FOR EACH ROW EXECUTE FUNCTION ubi_v0.reject_outbox()")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, e := admin.Exec(testContext, "DROP TRIGGER reject_outbox ON ubi_v0.outbox; DROP FUNCTION ubi_v0.reject_outbox()")
		if e != nil {
			t.Error(e)
		}
	})
	_, err = s.CreateBooking(testContext, ownerA(), "forced-rollback-001", intentAt(9))
	requireError(t, err, "storage_failure")
	for _, tab := range []string{"bookings", "resource_reservations", "commands", "audit", "outbox"} {
		if count(t, tab) != 0 {
			t.Fatalf("rollback left %s", tab)
		}
	}
}
func TestQuoteVersionAndAtomicConversion(t *testing.T) {
	s := reset(t)
	price := int64(123457)
	in := QuoteInput{Intent: intentAt(9), RequestID: "demo-request-wedding", DemoPriceMinor: &price}
	r, e := s.IssueQuote(testContext, ownerA(), "quote-version-one", in)
	v1 := success(t, r, e)
	price2 := int64(223459)
	in.DemoPriceMinor = &price2
	r, e = s.IssueQuote(testContext, ownerA(), "quote-version-two", in)
	v2 := success(t, r, e)
	if v2.Result.Version != 2 || count(t, "bookings") != 0 || count(t, "resource_reservations") != 0 {
		t.Fatal("quote issuance allocated capacity or lost version")
	}
	_, e = s.AcceptQuote(testContext, ownerA(), "quote-accept-stale", v1.Result.ID, VersionInput{1})
	requireError(t, e, "quote_not_valid")
	r, e = s.AcceptQuote(testContext, ownerA(), "quote-accept-current", v2.Result.ID, VersionInput{2})
	b := success(t, r, e)
	r, e = s.AcceptQuote(testContext, ownerA(), "quote-accept-current", v2.Result.ID, VersionInput{2})
	if success(t, r, e).Result.ID != b.Result.ID {
		t.Fatal("duplicate quote conversion")
	}
	_, e = s.AcceptQuote(testContext, ownerA(), "quote-accept-new-key", v2.Result.ID, VersionInput{2})
	requireError(t, e, "quote_not_valid")
	if b.Result.PriceMinor != 223459 || count(t, "bookings") != 1 || count(t, "resource_reservations") != 1 {
		t.Fatal("price/version conversion mismatch")
	}
	tx, err := s.begin(testContext, ownerA())
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec(testContext, "UPDATE ubi_v0.quotations SET price_minor=0 WHERE id=$1", v2.Result.ID)
	if err == nil {
		t.Fatal("quote terms could be rewritten")
	}
	_ = tx.Rollback(testContext)
}
func TestConflictingQuoteAcceptanceLeavesOfferAndCapacityIntact(t *testing.T) {
	s := reset(t)
	r, e := s.IssueQuote(testContext, ownerA(), "blocked-quote-001", QuoteInput{Intent: intentAt(9), RequestID: "demo-request-blocked"})
	q := success(t, r, e)
	create(t, s, "block-quote-slot-001", 9)
	_, e = s.AcceptQuote(testContext, ownerA(), "blocked-accept-001", q.Result.ID, VersionInput{1})
	requireError(t, e, "capacity_conflict")
	var state string
	var linked *string
	if err := admin.QueryRow(testContext, "SELECT state,booking_id FROM ubi_v0.quotations WHERE id=$1", q.Result.ID).Scan(&state, &linked); err != nil || state != "issued" || linked != nil {
		t.Fatal("failed acceptance partially converted")
	}
	if count(t, "bookings") != 1 || count(t, "commands") != 2 || count(t, "outbox") != 2 {
		t.Fatal("partial quote command retained")
	}
}
func TestInvalidInputServiceVersionAndExactOffset(t *testing.T) {
	s := reset(t)
	in := intentAt(9)
	in.CustomerLabel = "Real-looking customer"
	_, e := s.CreateBooking(testContext, ownerA(), "invalid-label-001", in)
	requireError(t, e, "invalid_input")
	in = intentAt(9)
	in.ExpectedServiceVersion = 2
	_, e = s.CreateBooking(testContext, ownerA(), "stale-service-001", in)
	requireError(t, e, "service_version_conflict")
	in = intentAt(9)
	r, e := s.CreateBooking(testContext, ownerA(), "offset-first-001", in)
	first := success(t, r, e)
	in.StartAt = in.StartAt.UTC()
	r, e = s.CreateBooking(testContext, ownerA(), "offset-first-001", in)
	if success(t, r, e).Result.ID != first.Result.ID {
		t.Fatal("equivalent offset created another command")
	}
	if first.Result.StartAt.Hour() != 6 || first.Result.PriceMinor != 200000 {
		t.Fatal("offset or integer money changed")
	}
}
