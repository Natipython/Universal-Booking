package integrity

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestExclusionConstraintBackstopWithoutServiceLocks(t *testing.T) {
	s := reset(t)
	create(t, s, "direct-first-001", 9)
	second := create(t, s, "direct-second-001", 10)
	tx, err := s.begin(testContext, ownerA())
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(testContext)
	_, err = tx.Exec(testContext, "UPDATE ubi_v0.resource_reservations SET occupied_interval=tstzrange($3,$4,'[)') WHERE organization_id=$1 AND booking_id=$2", ownerA().OrganizationID, second.Result.ID, intentAt(9).StartAt, intentAt(10).StartAt)
	var pg *pgconn.PgError
	if !errors.As(err, &pg) || pg.Code != "23P01" {
		t.Fatalf("database itself allowed overlap: %v", err)
	}
}

func TestHeldConversionAndBookingRace(t *testing.T) {
	for _, expired := range []bool{false, true} {
		t.Run(fmt.Sprintf("expired=%t", expired), func(t *testing.T) {
			s := reset(t)
			r, e := s.CreateHold(testContext, ownerA(), "racing-hold-001", HoldInput{intentAt(9), 300})
			h := success(t, r, e)
			if expired {
				if _, e = admin.Exec(testContext, "UPDATE ubi_v0.holds SET created_at=clock_timestamp()-interval '10 minutes',expires_at=clock_timestamp()-interval '1 minute' WHERE id=$1", h.Result.ID); e != nil {
					t.Fatal(e)
				}
			}
			gate := make(chan struct{})
			out := make(chan error, 2)
			go func() {
				<-gate
				_, e := s.ConsumeHold(testContext, ownerA(), "racing-consume-001", h.Result.ID, VersionInput{1})
				out <- e
			}()
			go func() {
				<-gate
				_, e := s.CreateBooking(testContext, ownerA(), "racing-booking-001", intentAt(9))
				out <- e
			}()
			close(gate)
			wins := 0
			for i := 0; i < 2; i++ {
				if e := <-out; e == nil {
					wins++
				} else {
					var d *DomainError
					if !errors.As(e, &d) || d.Status < 400 || d.Status >= 500 {
						t.Fatalf("unexpected race failure: %v", e)
					}
				}
			}
			if wins != 1 || count(t, "bookings") != 1 {
				t.Fatalf("hold race made %d winners", wins)
			}
			var active int
			if e = admin.QueryRow(testContext, "SELECT count(*) FROM ubi_v0.resource_reservations WHERE state IN ('held','committed')").Scan(&active); e != nil || active != 1 {
				t.Fatal("capacity leak/duplicate")
			}
		})
	}
}

func TestConcurrentQuotationConversionHasOneBooking(t *testing.T) {
	s := reset(t)
	r, e := s.IssueQuote(testContext, ownerA(), "racing-quote-001", QuoteInput{Intent: intentAt(9), RequestID: "demo-request-race"})
	q := success(t, r, e)
	gate := make(chan struct{})
	out := make(chan error, 8)
	for i := 0; i < 8; i++ {
		go func(i int) {
			<-gate
			_, e := s.AcceptQuote(testContext, ownerA(), fmt.Sprintf("quote-race-key-%02d", i), q.Result.ID, VersionInput{1})
			out <- e
		}(i)
	}
	close(gate)
	wins := 0
	for i := 0; i < 8; i++ {
		e := <-out
		if e == nil {
			wins++
		} else {
			requireError(t, e, "quote_not_valid")
		}
	}
	if wins != 1 || count(t, "bookings") != 1 || count(t, "resource_reservations") != 1 || count(t, "outbox") != 3 {
		t.Fatal("quotation race duplicated booking/effects")
	}
}

func TestActualHTTPResponseLossReplaysCommittedBooking(t *testing.T) {
	s := reset(t)
	handler := s.DemoHandler()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Ubi-Demo-Drop-Response") == "1" {
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, r)
			if recorder.Code != 201 {
				t.Errorf("fault injection must follow a committed response, got %d", recorder.Code)
			}
			conn, _, err := w.(http.Hijacker).Hijack()
			if err != nil {
				t.Error(err)
				return
			}
			_ = conn.Close()
			return
		}
		handler.ServeHTTP(w, r)
	}))
	defer server.Close()
	request := func(drop bool) (*http.Response, error) {
		req, e := http.NewRequest("POST", server.URL+"/v0/bookings", bytes.NewReader(jsonBytes(intentAt(9))))
		if e != nil {
			return nil, e
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Ubi-Demo-Principal", "demo-owner-a")
		req.Header.Set("Idempotency-Key", "dropped-response-001")
		if drop {
			req.Header.Set("X-Ubi-Demo-Drop-Response", "1")
		}
		return server.Client().Do(req)
	}
	response, err := request(true)
	if response != nil {
		_ = response.Body.Close()
	}
	if err == nil {
		t.Fatal("client did not observe dropped response")
	}
	if count(t, "bookings") != 1 || count(t, "outbox") != 1 {
		t.Fatal("outcome not committed before loss")
	}
	response, err = request(false)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil || response.StatusCode != 201 || !bytes.Contains(body, []byte(`"confirmed_demo"`)) {
		t.Fatalf("committed response not recovered: %s %v", body, err)
	}
	if count(t, "bookings") != 1 || count(t, "audit") != 1 || count(t, "outbox") != 1 {
		t.Fatal("lost-response retry duplicated an effect")
	}
}

func TestPostmanV3ContractCollection(t *testing.T) {
	s := reset(t)
	server := httptest.NewServer(s.DemoHandler())
	defer server.Close()
	cli, err := filepath.Abs("../../.tooling/postman-runtime/node_modules/.bin/postman")
	if err != nil {
		t.Fatal(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(cli, "collection", "run", "postman/collections/ubi-synthetic-v0", "--env-var", "base_url="+server.URL,
		"--no-report-events", "--working-dir", cwd, "--no-insecure-file-read", "--timeout", "30000", "--timeout-request", "5000", "--color", "off")
	output, err := cmd.CombinedOutput()
	if e := os.MkdirAll("artifacts", 0755); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile("artifacts/postman-run.txt", output, 0644); e != nil {
		t.Fatal(e)
	}
	if err != nil {
		t.Fatalf("Postman collection failed: %v\n%s", err, output)
	}
	t.Logf("Actual Postman local collection result:\n%s", output)
	if count(t, "bookings") != 4 || count(t, "holds") != 1 || count(t, "quotations") != 1 || count(t, "outbox") != 10 {
		t.Fatal("HTTP contract run left unexpected durable effects")
	}
}

func TestCatalogueChangeDoesNotSilentlyRewriteQuote(t *testing.T) {
	s := reset(t)
	r, e := s.IssueQuote(testContext, ownerA(), "catalogue-quote-001", QuoteInput{Intent: intentAt(9), RequestID: "demo-request-catalogue"})
	q := success(t, r, e)
	if _, e = admin.Exec(testContext, "UPDATE ubi_v0.services SET version=2,price_minor=777777 WHERE organization_id='demo-org-a' AND operation_id='studio' AND service_id='portrait'"); e != nil {
		t.Fatal(e)
	}
	_, e = s.AcceptQuote(testContext, ownerA(), "catalogue-accept-001", q.Result.ID, VersionInput{1})
	requireError(t, e, "service_version_conflict")
	var price int64
	if e = admin.QueryRow(testContext, "SELECT price_minor FROM ubi_v0.quotations WHERE id=$1", q.Result.ID).Scan(&price); e != nil || price != 200000 || count(t, "bookings") != 0 {
		t.Fatal("old quote was silently rewritten/converted")
	}
}

func TestOutboxFailureDuringMovePreservesOriginalAllocation(t *testing.T) {
	s := reset(t)
	first := create(t, s, "outbox-move-first", 9)
	_, err := admin.Exec(testContext, "CREATE FUNCTION ubi_v0.reject_move_outbox() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'synthetic move failure'; END $$; CREATE TRIGGER reject_move_outbox BEFORE INSERT ON ubi_v0.outbox FOR EACH ROW EXECUTE FUNCTION ubi_v0.reject_move_outbox()")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, e := admin.Exec(testContext, "DROP TRIGGER reject_move_outbox ON ubi_v0.outbox; DROP FUNCTION ubi_v0.reject_move_outbox()")
		if e != nil {
			t.Error(e)
		}
	})
	_, err = s.Reschedule(testContext, ownerA(), "outbox-move-failed", first.Result.ID, RescheduleInput{VersionInput{1}, intentAt(10).StartAt})
	requireError(t, err, "storage_failure")
	old, err := s.GetBooking(testContext, ownerA(), first.Result.ID)
	if err != nil || old.Version != 1 || !old.StartAt.Equal(first.Result.StartAt) {
		t.Fatal("failed move lost original booking")
	}
	var lower time.Time
	if e := admin.QueryRow(testContext, "SELECT lower(occupied_interval) FROM ubi_v0.resource_reservations WHERE booking_id=$1", first.Result.ID).Scan(&lower); e != nil || !lower.Equal(first.Result.StartAt) {
		t.Fatal("allocation partially moved")
	}
	if count(t, "commands") != 1 || count(t, "audit") != 1 || count(t, "outbox") != 1 {
		t.Fatal("failed move left partial history")
	}
}
