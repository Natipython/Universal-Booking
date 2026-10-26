package customerweb

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	integrity "ubi.local/booking-integrity-v0"
)

type Disposable struct {
	dir            string
	Admin, Runtime *pgxpool.Pool
	running        bool
}

func pgRun(name string, args ...string) error {
	out, err := exec.Command(filepath.Join("/usr/lib/postgresql/18/bin", name), args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("private synthetic %s failed: %w (%d diagnostic bytes)", name, err, len(out))
	}
	return nil
}

// StartDisposable accepts no connection string. Only its own new directory is a removal target.
func StartDisposable() (*Disposable, error) {
	if os.Getenv("UBI_RUN_SYNTHETIC_WEB") != "1" {
		return nil, fmt.Errorf("set UBI_RUN_SYNTHETIC_WEB=1 for invented local data only; no existing database is accepted")
	}
	dir, err := os.MkdirTemp("", "ubi-wb04-")
	if err != nil {
		return nil, err
	}
	d := &Disposable{dir: dir}
	data := filepath.Join(dir, "data")
	failed := true
	defer func() {
		if failed {
			_ = d.Close()
		}
	}()
	if err = pgRun("initdb", "-D", data, "-U", "ubi_v0_admin", "--auth-local=trust", "--auth-host=reject", "--encoding=UTF8", "--no-locale"); err != nil {
		return nil, err
	}
	d.running = true
	if err = pgRun("pg_ctl", "-D", data, "-l", filepath.Join(dir, "postgres.log"), "-o", fmt.Sprintf("-h '' -k %s -p 5432 -c max_connections=32", dir), "-w", "start"); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	dsn := fmt.Sprintf("host=%s port=5432 dbname=postgres user=ubi_v0_admin sslmode=disable timezone=UTC", dir)
	d.Admin, err = pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	schema, err := os.ReadFile("../booking_integrity/schema.sql")
	if err != nil {
		return nil, err
	}
	if _, err = d.Admin.Exec(ctx, string(schema)); err != nil {
		return nil, fmt.Errorf("synthetic schema setup failed: %w", err)
	}
	d.Runtime, err = pgxpool.New(ctx, strings.Replace(dsn, "user=ubi_v0_admin", "user=ubi_v0_runtime", 1))
	if err != nil {
		return nil, err
	}
	if err = d.Runtime.Ping(ctx); err != nil {
		return nil, err
	}
	if err = d.Reset(); err != nil {
		return nil, err
	}
	failed = false
	return d, nil
}
func (d *Disposable) Close() error {
	if d.Runtime != nil {
		d.Runtime.Close()
		d.Runtime = nil
	}
	if d.Admin != nil {
		d.Admin.Close()
		d.Admin = nil
	}
	stopped := !d.running
	if d.running {
		if err := pgRun("pg_ctl", "-D", filepath.Join(d.dir, "data"), "-m", "fast", "-w", "stop"); err == nil {
			d.running = false
			stopped = true
		}
	}
	removed := false
	if stopped && filepath.Dir(d.dir) == os.TempDir() && strings.HasPrefix(filepath.Base(d.dir), "ubi-wb04-") {
		removed = os.RemoveAll(d.dir) == nil
	}
	_ = os.MkdirAll("artifacts", 0755)
	raw, _ := json.Marshal(map[string]bool{"synthetic_only": true, "new_private_temporary_cluster": true, "existing_database_accepted": false, "tcp_enabled": false, "server_stopped": stopped, "temporary_files_removed": removed})
	_ = os.WriteFile("artifacts/cluster-receipt.json", raw, 0644)
	if !stopped || !removed {
		return fmt.Errorf("private temporary demo files retained or cleanup failed; no broad removal attempted")
	}
	return nil
}
func (d *Disposable) Reset() error {
	if _, err := d.Admin.Exec(context.Background(), "TRUNCATE ubi_v0.audit,ubi_v0.outbox,ubi_v0.commands,ubi_v0.quotations,ubi_v0.resource_reservations,ubi_v0.holds,ubi_v0.bookings,ubi_v0.services,ubi_v0.resources,ubi_v0.operations,ubi_v0.organizations"); err != nil {
		return err
	}
	var f struct {
		Synthetic     bool                  `json:"syntheticOnly"`
		Exponent      int                   `json:"currencyExponent"`
		Organizations []struct{ ID string } `json:"organizations"`
		Operations    []struct {
			Organization string `json:"organizationId"`
			ID           string
		} `json:"operations"`
		Services []struct {
			Organization string `json:"organizationId"`
			Operation    string `json:"operationId"`
			ID, Name     string
			Resource     string `json:"resourceId"`
			Duration     int    `json:"durationMinutes"`
			Price        int64  `json:"priceMinor"`
		} `json:"services"`
	}
	raw, err := os.ReadFile("../fixtures/booking-v0.json")
	if err != nil {
		return err
	}
	if err = json.Unmarshal(raw, &f); err != nil {
		return err
	}
	if !f.Synthetic || f.Exponent != 2 || len(f.Organizations) != 2 || len(f.Operations) != 3 || len(f.Services) != 4 {
		return fmt.Errorf("refusing non-synthetic fixture shape")
	}
	for _, o := range f.Organizations {
		if _, err = d.Admin.Exec(context.Background(), "INSERT INTO ubi_v0.organizations VALUES($1)", o.ID); err != nil {
			return err
		}
	}
	for _, o := range f.Operations {
		if _, err = d.Admin.Exec(context.Background(), "INSERT INTO ubi_v0.operations VALUES($1,$2)", o.Organization, o.ID); err != nil {
			return err
		}
	}
	for _, s := range f.Services {
		if _, err = d.Admin.Exec(context.Background(), "INSERT INTO ubi_v0.resources VALUES($1,$2) ON CONFLICT DO NOTHING", s.Organization, s.Resource); err != nil {
			return err
		}
		if _, err = d.Admin.Exec(context.Background(), "INSERT INTO ubi_v0.services(organization_id,operation_id,service_id,resource_id,version,duration_minutes,price_minor,currency) VALUES($1,$2,$3,$4,1,$5,$6,'ETB')", s.Organization, s.Operation, s.ID, s.Resource, s.Duration, s.Price); err != nil {
			return err
		}
	}
	return nil
}

type SQLGateway struct{ Pool *pgxpool.Pool }

func (g SQLGateway) Snapshot(ctx context.Context, scope Scope) (Catalog, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if !validScope(scope) {
		return Catalog{}, fmt.Errorf("unsupported synthetic scope")
	}
	tx, err := g.Pool.Begin(ctx)
	if err != nil {
		return Catalog{}, err
	}
	defer tx.Rollback(context.Background())
	if _, err = tx.Exec(ctx, "SELECT set_config('ubi_v0.organization_id',$1,true)", scope.Organization); err != nil {
		return Catalog{}, err
	}
	var c Catalog
	err = tx.QueryRow(ctx, "SELECT version,duration_minutes,price_minor FROM ubi_v0.services WHERE organization_id=$1 AND operation_id=$2 AND service_id=$3", scope.Organization, scope.Operation, scope.Service).Scan(&c.Version, &c.Duration, &c.Price)
	if err != nil {
		return c, err
	}
	c.Name = "Portrait session"
	if scope.Service == "consultation" {
		c.Name = "Event consultation"
	}
	return c, tx.Commit(ctx)
}
func (g SQLGateway) Book(ctx context.Context, scope Scope, in integrity.Intent, key string) (integrity.Result, error) {
	if !validScope(scope) || in.OrganizationID != scope.Organization || in.OperationID != scope.Operation || in.ServiceID != scope.Service {
		return integrity.Result{}, &integrity.DomainError{Code: "permission_denied", Status: 403, Detail: "Unsupported demo grant scope."}
	}
	name := "demo-owner-a"
	if scope.Organization == "demo-org-b" {
		name = "demo-owner-b"
	}
	p, _ := integrity.DemoPrincipal(name) // Server-selected fixture only; NOT a real customer's identity or owner role.
	r, err := (&integrity.Service{Pool: g.Pool}).CreateBooking(ctx, p, key, in)
	return r.Result, err
}

// Recover is read-only: no command insert, reservation or automatic new attempt.
func (g SQLGateway) Recover(ctx context.Context, scope Scope, in integrity.Intent, key string) (integrity.Result, error) {
	if !validScope(scope) || in.OrganizationID != scope.Organization || in.OperationID != scope.Operation || in.ServiceID != scope.Service {
		return integrity.Result{}, &integrity.DomainError{Code: "permission_denied", Status: 403, Detail: "Unsupported demo grant scope."}
	}
	tx, err := g.Pool.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
	if err != nil {
		return integrity.Result{}, err
	}
	defer tx.Rollback(context.Background())
	if _, err = tx.Exec(ctx, "SELECT set_config('ubi_v0.organization_id',$1,true)", scope.Organization); err != nil {
		return integrity.Result{}, err
	}
	name := "demo-owner-a"
	if scope.Organization == "demo-org-b" {
		name = "demo-owner-b"
	}
	var raw []byte
	var storedHash string
	err = tx.QueryRow(ctx, "SELECT request_hash,result FROM ubi_v0.commands WHERE organization_id=$1 AND principal_id=$2 AND action='create_booking' AND key=$3", scope.Organization, name, key).Scan(&storedHash, &raw)
	if err != nil || raw == nil {
		return integrity.Result{}, &integrity.DomainError{Code: "outcome_unresolved", Status: 503, Detail: "No committed receipt is visible; do not create a new booking automatically.", Retryable: true}
	}
	in.StartAt = in.StartAt.UTC()
	canonical, err := json.Marshal(in)
	if err != nil {
		return integrity.Result{}, err
	}
	sum := sha256.Sum256(canonical)
	if storedHash != fmt.Sprintf("%x", sum[:]) {
		return integrity.Result{}, fmt.Errorf("receipt input mismatch")
	}
	var receipt integrity.Receipt
	if err = json.Unmarshal(raw, &receipt); err != nil {
		return integrity.Result{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return integrity.Result{}, err
	}
	return receipt.Result, nil
}
