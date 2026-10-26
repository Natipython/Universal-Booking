package customerweb

import (
	"context"
	"encoding/json"
	"errors"
	"html"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	integrity "ubi.local/booking-integrity-v0"
)

var testDB *Disposable

func TestMain(m *testing.M) {
	db, err := StartDisposable()
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(2)
	}
	testDB = db
	code := m.Run()
	if err = db.Close(); err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		code = 1
	}
	os.Exit(code)
}

type fixture struct {
	app    *App
	server *httptest.Server
	client *http.Client
}

func reset(t *testing.T, gateway Gateway) *fixture {
	return resetWithHandler(t, gateway, nil)
}
func resetWithHandler(t *testing.T, gateway Gateway, wrap func(http.Handler) http.Handler) *fixture {
	t.Helper()
	if err := testDB.Reset(); err != nil {
		t.Fatal(err)
	}
	if gateway == nil {
		gateway = SQLGateway{Pool: testDB.Runtime}
	}
	server := httptest.NewUnstartedServer(nil)
	a, err := New("http://"+server.Listener.Addr().String(), gateway)
	if err != nil {
		server.Close()
		t.Fatal(err)
	}
	server.Config.Handler = a
	if wrap != nil {
		server.Config.Handler = wrap(a)
	}
	server.Start()
	jar, _ := cookiejar.New(nil)
	f := &fixture{a, server, &http.Client{Jar: jar, Timeout: 10 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}
	t.Cleanup(server.Close)
	return f
}

type response struct {
	status int
	body   string
	header http.Header
}

func (f *fixture) request(t *testing.T, method, path string, form url.Values, change func(*http.Request)) response {
	t.Helper()
	var body io.Reader
	if form != nil {
		body = strings.NewReader(form.Encode())
	}
	r, err := http.NewRequest(method, f.server.URL+path, body)
	if err != nil {
		t.Fatal(err)
	}
	if method == "POST" {
		r.Header.Set("Origin", f.server.URL)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if change != nil {
		change(r)
	}
	out, err := f.client.Do(r)
	if err != nil {
		t.Fatal("local HTTP request failed")
	}
	defer out.Body.Close()
	raw, err := io.ReadAll(out.Body)
	if err != nil {
		t.Fatal("local response read failed")
	}
	return response{out.StatusCode, string(raw), out.Header}
}
func requireStatus(t *testing.T, r response, status int) {
	t.Helper()
	if r.status != status {
		t.Fatalf("wanted HTTP %d, got %d (response body omitted to protect demo grants)", status, r.status)
	}
}
func hidden(t *testing.T, body, name string) string {
	t.Helper()
	m := regexp.MustCompile(`name="` + regexp.QuoteMeta(name) + `" value="([^"]+)"`).FindStringSubmatch(body)
	if len(m) != 2 {
		t.Fatalf("missing hidden field %s", name)
	}
	return m[1]
}
func (f *fixture) start(t *testing.T, name string) string {
	r := f.request(t, "GET", "/", nil, nil)
	requireStatus(t, r, 200)
	r = f.request(t, "POST", "/demo/starts", url.Values{"csrf": {hidden(t, r.body, "csrf")}, "scenario": {name}}, nil)
	requireStatus(t, r, 303)
	path := r.header.Get("Location")
	if !strings.HasPrefix(path, "/h/") || !secretShape(strings.TrimPrefix(path, "/h/")) {
		t.Fatal("non-opaque redirect")
	}
	return path
}
func (f *fixture) redeem(t *testing.T, path string) response {
	r := f.request(t, "GET", path, nil, nil)
	requireStatus(t, r, 200)
	return f.request(t, "POST", "/v0/handoff-sessions", url.Values{"csrf": {hidden(t, r.body, "csrf")}, "token": {strings.TrimPrefix(path, "/h/")}}, nil)
}
func (f *fixture) enter(t *testing.T, name string) string {
	path := f.start(t, name)
	requireStatus(t, f.redeem(t, path), 303)
	return path
}
func (f *fixture) reviewForm(t *testing.T) url.Values {
	r := f.request(t, "GET", "/continue", nil, nil)
	requireStatus(t, r, 200)
	return url.Values{"csrf": {hidden(t, r.body, "csrf")}, "expected_draft_version": {hidden(t, r.body, "expected_draft_version")}, "service_id": {hidden(t, r.body, "service_id")}, "date": {"2030-01-15"}, "time": {"10:00"}, "calendar": {"Gregorian"}, "clock": {"international-24h"}, "customer_label": {"Demo Guest 02"}, "event_size": {"not-specified"}}
}
func (f *fixture) reviewed(t *testing.T) url.Values {
	requireStatus(t, f.request(t, "POST", "/v0/reviews", f.reviewForm(t), nil), 303)
	r := f.request(t, "GET", "/review", nil, nil)
	requireStatus(t, r, 200)
	return url.Values{"csrf": {hidden(t, r.body, "csrf")}, "review_id": {hidden(t, r.body, "review_id")}}
}
func count(t *testing.T, table string) int {
	t.Helper()
	allowed := map[string]bool{"bookings": true, "commands": true, "outbox": true, "audit": true, "resource_reservations": true}
	if !allowed[table] {
		t.Fatal("unapproved count table")
	}
	var n int
	if err := testDB.Admin.QueryRow(context.Background(), "SELECT count(*) FROM ubi_v0."+table).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}
func safeHeaders(t *testing.T, r response) {
	t.Helper()
	for k, v := range map[string]string{"Cache-Control": "no-store", "Referrer-Policy": "no-referrer", "X-Frame-Options": "DENY", "X-Content-Type-Options": "nosniff"} {
		if r.header.Get(k) != v {
			t.Errorf("missing header %s", k)
		}
	}
	if !strings.Contains(r.header.Get("Content-Security-Policy"), "default-src 'none'") {
		t.Fatal("missing restrictive CSP")
	}
}

func TestLoopbackAndDisposableBoundary(t *testing.T) {
	f := reset(t, nil)
	for _, origin := range []string{"https://example.com:443", "http://127.0.0.1", "http://user@127.0.0.1:1", "http://127.0.0.1:1/path", "http://127.0.0.1:1?x=1"} {
		if _, e := New(origin, f.app.gateway); e == nil {
			t.Fatal("accepted out-of-scope origin")
		}
	}
	var listen string
	if e := testDB.Admin.QueryRow(context.Background(), "SHOW listen_addresses").Scan(&listen); e != nil || listen != "" {
		t.Fatal("database TCP listener enabled")
	}
	var super, bypass bool
	if e := testDB.Runtime.QueryRow(context.Background(), "SELECT rolsuper,rolbypassrls FROM pg_roles WHERE rolname=current_user").Scan(&super, &bypass); e != nil || super || bypass {
		t.Fatal("runtime database role too privileged")
	}
	if filepath.Dir(testDB.dir) != os.TempDir() || !strings.HasPrefix(filepath.Base(testDB.dir), "ubi-wb04-") {
		t.Fatal("not a newly owned temporary directory")
	}
}
func TestCrawlerPreviewAndHeadDoNotRedeem(t *testing.T) {
	f := reset(t, nil)
	path := f.start(t, "instant")
	gh := digest(strings.TrimPrefix(path, "/h/"))
	for _, method := range []string{"GET", "GET", "HEAD"} {
		r := f.request(t, method, path, nil, func(r *http.Request) { r.Header.Set("User-Agent", "TelegramBot (link preview)") })
		requireStatus(t, r, 200)
		safeHeaders(t, r)
		if strings.Contains(r.body, "2030-01-15") || strings.Contains(r.body, "Demo Guest 01") {
			t.Fatal("private draft exposed in preview")
		}
		if method == "HEAD" && r.body != "" {
			t.Fatal("HEAD emitted a body")
		}
	}
	if f.app.grants[gh].BoundBrowser != "" || len(f.app.sessions) != 0 || count(t, "bookings") != 0 {
		t.Fatal("preview performed an irreversible action")
	}
	requireStatus(t, f.redeem(t, path), 303)
}
func TestGrantDigestsRedemptionRetryAndForwarding(t *testing.T) {
	f := reset(t, nil)
	path := f.start(t, "instant")
	token := strings.TrimPrefix(path, "/h/")
	if f.app.grants[token] != nil || f.app.grants[digest(token)] == nil {
		t.Fatal("raw grant secret retained")
	}
	r := f.redeem(t, path)
	requireStatus(t, r, 303)
	first := r.header.Get("Set-Cookie")
	r = f.redeem(t, path)
	requireStatus(t, r, 303)
	if first != r.header.Get("Set-Cookie") || len(f.app.sessions) != 1 {
		t.Fatal("redemption retry changed the session")
	}
	jar, _ := cookiejar.New(nil)
	old := f.client.Jar
	f.client.Jar = jar
	r = f.redeem(t, path)
	requireStatus(t, r, 410)
	if strings.Contains(r.body, "2030-01-15") {
		t.Fatal("forwarded link exposed draft")
	}
	f.client.Jar = old
	requireStatus(t, f.request(t, "GET", "/continue", nil, nil), 200)
}
func TestUnusedForwardedLinkIsNotIdentityVerification(t *testing.T) {
	f := reset(t, nil)
	path := f.start(t, "instant")
	original := f.client.Jar
	other, _ := cookiejar.New(nil)
	f.client.Jar = other
	requireStatus(t, f.redeem(t, path), 303)
	f.client.Jar = original
	requireStatus(t, f.redeem(t, path), 410)
	if count(t, "bookings") != 0 {
		t.Fatal("redemption alone booked capacity")
	}
}
func TestUnavailableGrantsHaveGenericPreviewAndDenial(t *testing.T) {
	for _, kind := range []string{"unknown", "expired", "revoked", "purpose", "scope", "mode"} {
		t.Run(kind, func(t *testing.T) {
			f := reset(t, nil)
			path := f.start(t, "instant")
			g := f.app.grants[digest(strings.TrimPrefix(path, "/h/"))]
			switch kind {
			case "unknown":
				delete(f.app.grants, digest(strings.TrimPrefix(path, "/h/")))
			case "expired":
				g.Expires = time.Now().Add(-time.Minute)
			case "revoked":
				g.Revoked = true
			case "purpose":
				g.Purpose = "manage_booking"
			case "scope":
				g.Scope.Operation = "foreign"
			case "mode":
				g.Mode = "custom"
			}
			r := f.request(t, "GET", path, nil, nil)
			requireStatus(t, r, 200)
			if !strings.Contains(r.body, "Continue your booking") {
				t.Fatal("preview disclosed validity")
			}
			r = f.redeem(t, path)
			requireStatus(t, r, 410)
			if !strings.Contains(r.body, `data-problem="link_unavailable"`) {
				t.Fatal("non-generic denial")
			}
			if len(f.app.sessions) != 0 || count(t, "bookings") != 0 {
				t.Fatal("invalid grant mutated business state")
			}
		})
	}
}
func TestRevocationAndExpiryCheckedAfterRedemption(t *testing.T) {
	for _, kind := range []string{"revoked", "expired", "purpose", "scope", "mode"} {
		t.Run(kind, func(t *testing.T) {
			f := reset(t, nil)
			path := f.enter(t, "instant")
			submit := f.reviewed(t)
			g := f.app.grants[digest(strings.TrimPrefix(path, "/h/"))]
			switch kind {
			case "revoked":
				g.Revoked = true
			case "expired":
				g.Expires = time.Now().Add(-time.Second)
			case "purpose":
				g.Purpose = "export"
			case "scope":
				g.Scope.Service = "birthday"
			case "mode":
				g.Mode = "custom"
			}
			for _, route := range []string{"/continue", "/review", "/outcome"} {
				requireStatus(t, f.request(t, "GET", route, nil, nil), 410)
			}
			requireStatus(t, f.request(t, "POST", "/v0/submissions", submit, nil), 410)
			if count(t, "bookings") != 0 {
				t.Fatal("revoked grant booked capacity")
			}
		})
	}
}
func TestFormProtocolAndCSRFNegatives(t *testing.T) {
	for _, kind := range []string{"csrf", "origin", "missing-origin", "duplicate-origin", "null-origin-missing-metadata", "null-origin-cross-site", "null-origin-same-site", "null-origin-wrong-mode", "null-origin-wrong-destination", "null-origin-bad-csrf", "fetch-site", "host", "duplicate", "extra-authority", "missing", "media", "query", "oversized", "duplicate-cookie"} {
		t.Run(kind, func(t *testing.T) {
			f := reset(t, nil)
			f.enter(t, "instant")
			form := f.reviewForm(t)
			path := "/v0/reviews"
			status := 400
			change := func(*http.Request) {}
			switch kind {
			case "csrf":
				form.Set("csrf", randomSecret())
				status = 403
			case "origin":
				change = func(r *http.Request) { r.Header.Set("Origin", "https://evil.example") }
				status = 403
			case "missing-origin":
				change = func(r *http.Request) { r.Header.Del("Origin") }
				status = 403
			case "duplicate-origin":
				change = func(r *http.Request) { r.Header.Add("Origin", "https://evil.example") }
				status = 403
			case "null-origin-missing-metadata", "null-origin-cross-site", "null-origin-same-site", "null-origin-wrong-mode", "null-origin-wrong-destination", "null-origin-bad-csrf":
				change = func(r *http.Request) {
					r.Header.Set("Origin", "null")
					if kind != "null-origin-missing-metadata" {
						r.Header.Set("Sec-Fetch-Site", "same-origin")
						r.Header.Set("Sec-Fetch-Mode", "navigate")
						r.Header.Set("Sec-Fetch-Dest", "document")
					}
					if kind == "null-origin-cross-site" {
						r.Header.Set("Sec-Fetch-Site", "cross-site")
					}
					if kind == "null-origin-same-site" {
						r.Header.Set("Sec-Fetch-Site", "same-site")
					}
					if kind == "null-origin-wrong-mode" {
						r.Header.Set("Sec-Fetch-Mode", "cors")
					}
					if kind == "null-origin-wrong-destination" {
						r.Header.Set("Sec-Fetch-Dest", "iframe")
					}
				}
				if kind == "null-origin-bad-csrf" {
					form.Set("csrf", randomSecret())
				}
				status = 403
			case "fetch-site":
				change = func(r *http.Request) { r.Header.Set("Sec-Fetch-Site", "cross-site") }
				status = 403
			case "host":
				change = func(r *http.Request) { r.Host = "evil.example" }
				status = 403
			case "duplicate":
				form.Add("date", "2030-01-16")
			case "extra-authority":
				form.Set("organization_id", "demo-org-b")
				form.Set("price_minor", "1")
			case "missing":
				form.Del("clock")
			case "media":
				change = func(r *http.Request) { r.Header.Set("Content-Type", "text/plain") }
				status = 415
			case "query":
				path += "?customer=private"
			case "oversized":
				form.Set("date", strings.Repeat("x", 9000))
				status = 413
			case "duplicate-cookie":
				change = func(r *http.Request) {
					r.AddCookie(&http.Cookie{Name: "ubi_demo_session", Value: randomSecret()})
					r.AddCookie(&http.Cookie{Name: "ubi_demo_session", Value: randomSecret()})
				}
				status = 401
			}
			r := f.request(t, "POST", path, form, change)
			requireStatus(t, r, status)
			safeHeaders(t, r)
			if count(t, "bookings") != 0 || count(t, "commands") != 0 {
				t.Fatal("invalid form changed booking state")
			}
		})
	}
}
func TestPrivateNoReferrerNavigationStillRequiresBoundCSRF(t *testing.T) {
	f := reset(t, nil)
	f.enter(t, "instant")
	change := func(r *http.Request) {
		r.Header.Set("Origin", "null")
		r.Header.Set("Sec-Fetch-Site", "same-origin")
		r.Header.Set("Sec-Fetch-Mode", "navigate")
		r.Header.Set("Sec-Fetch-Dest", "document")
	}
	requireStatus(t, f.request(t, "POST", "/v0/reviews", f.reviewForm(t), change), 303)
	if count(t, "bookings") != 0 {
		t.Fatal("private navigation bypassed explicit submit")
	}
}
func TestEditsReviewAndInstantOutcome(t *testing.T) {
	f := reset(t, nil)
	f.enter(t, "instant")
	submit := f.reviewed(t)
	r := f.request(t, "GET", "/review", nil, nil)
	for _, value := range []string{"10:00", "Demo Guest 02", "60 minutes", "2000.00 birr", "Gregorian", "UTC+03:00"} {
		if !strings.Contains(html.UnescapeString(r.body), value) {
			t.Errorf("missing review term %s", value)
		}
	}
	if count(t, "bookings") != 0 {
		t.Fatal("review reserved capacity")
	}
	requireStatus(t, f.request(t, "POST", "/v0/submissions", submit, nil), 303)
	r = f.request(t, "GET", "/outcome", nil, nil)
	requireStatus(t, r, 200)
	safeHeaders(t, r)
	if !strings.Contains(r.body, `data-state="confirmed_demo"`) || !strings.Contains(r.body, "No payment was taken") {
		t.Fatal("false/incomplete outcome")
	}
	for i := 0; i < 3; i++ {
		requireStatus(t, f.request(t, "POST", "/v0/submissions", submit, nil), 303)
	}
	for _, table := range []string{"bookings", "commands", "outbox", "audit", "resource_reservations"} {
		if count(t, table) != 1 {
			t.Fatalf("duplicate/partial %s", table)
		}
	}
	for _, route := range []string{"/continue", "/review"} {
		r = f.request(t, "GET", route, nil, nil)
		requireStatus(t, r, 303)
		if r.header.Get("Location") != "/outcome" {
			t.Fatal("completed attempt showed another commit form")
		}
	}
}
func TestCustomRequestIsNotReservationAndRetainsDetails(t *testing.T) {
	f := reset(t, nil)
	f.enter(t, "custom")
	form := f.reviewForm(t)
	form.Set("event_size", "small")
	requireStatus(t, f.request(t, "POST", "/v0/reviews", form, nil), 303)
	r := f.request(t, "GET", "/review", nil, nil)
	if !strings.Contains(r.body, "No time is held") || !strings.Contains(r.body, "<dd>small</dd>") {
		t.Fatal("custom expectation/edited scope missing")
	}
	submit := url.Values{"csrf": {hidden(t, r.body, "csrf")}, "review_id": {hidden(t, r.body, "review_id")}}
	for i := 0; i < 3; i++ {
		requireStatus(t, f.request(t, "POST", "/v0/submissions", submit, nil), 303)
	}
	r = f.request(t, "GET", "/outcome", nil, nil)
	requireStatus(t, r, 200)
	if !strings.Contains(r.body, `data-state="request_saved_demo"`) || strings.Contains(r.body, `class="card result success"`) || f.app.Requests != 1 || count(t, "bookings") != 0 || count(t, "resource_reservations") != 0 {
		t.Fatal("enquiry falsely confirmed/reserved/duplicated")
	}
}
func TestInvalidDateClockAndEscapingRetainSafeInput(t *testing.T) {
	for _, kind := range []string{"calendar", "clock", "date", "time", "customer", "xss"} {
		t.Run(kind, func(t *testing.T) {
			f := reset(t, nil)
			f.enter(t, "instant")
			form := f.reviewForm(t)
			switch kind {
			case "calendar":
				form.Set("calendar", "Ethiopian")
			case "clock":
				form.Set("clock", "Ethiopian")
			case "date":
				form.Set("date", "2030-02-30")
			case "time":
				form.Set("time", "2:00")
			case "customer":
				form.Set("customer_label", "Real Person")
			case "xss":
				form.Set("date", `"><script>alert(1)</script>`)
			}
			r := f.request(t, "POST", "/v0/reviews", form, nil)
			requireStatus(t, r, 422)
			if strings.Contains(r.body, "<script>") || strings.Contains(r.body, `value=""><script`) {
				t.Fatal("unescaped HTML input")
			}
			if !strings.Contains(r.body, `value="10:00"`) && kind != "time" {
				t.Fatal("safe edit lost during validation")
			}
			if count(t, "bookings") != 0 {
				t.Fatal("invalid details booked")
			}
		})
	}
}
func TestStaleForeignAndExpiredReviewCannotCommit(t *testing.T) {
	for _, kind := range []string{"stale-draft", "old-review", "foreign-review", "expired"} {
		t.Run(kind, func(t *testing.T) {
			f := reset(t, nil)
			f.enter(t, "instant")
			form := f.reviewForm(t)
			requireStatus(t, f.request(t, "POST", "/v0/reviews", form, nil), 303)
			r := f.request(t, "GET", "/review", nil, nil)
			submit := url.Values{"csrf": {hidden(t, r.body, "csrf")}, "review_id": {hidden(t, r.body, "review_id")}}
			switch kind {
			case "stale-draft":
				requireStatus(t, f.request(t, "POST", "/v0/reviews", form, nil), 409)
				return
			case "old-review":
				form = f.reviewForm(t)
				form.Set("time", "11:00")
				requireStatus(t, f.request(t, "POST", "/v0/reviews", form, nil), 303)
			case "foreign-review":
				submit.Set("review_id", randomSecret())
			case "expired":
				f.app.now = func() time.Time { return time.Now().Add(11 * time.Minute) }
			}
			requireStatus(t, f.request(t, "POST", "/v0/submissions", submit, nil), 409)
			if count(t, "bookings") != 0 {
				t.Fatal("invalid review committed")
			}
		})
	}
}
func TestCapacityAndChangedServicePreserveDraft(t *testing.T) {
	for _, kind := range []string{"capacity", "price"} {
		t.Run(kind, func(t *testing.T) {
			f := reset(t, nil)
			f.enter(t, "instant")
			submit := f.reviewed(t)
			if kind == "price" {
				_, e := testDB.Admin.Exec(context.Background(), "UPDATE ubi_v0.services SET version=2,price_minor=250000 WHERE organization_id='demo-org-a' AND operation_id='studio' AND service_id='portrait'")
				if e != nil {
					t.Fatal(e)
				}
			} else {
				var intent integrity.Intent
				for _, s := range f.app.sessions {
					intent = s.Review.Intent
				}
				_, e := f.app.gateway.Book(context.Background(), Scope{"demo-org-a", "studio", "portrait"}, intent, "competing-review-001")
				if e != nil {
					t.Fatal(e)
				}
			}
			r := f.request(t, "POST", "/v0/submissions", submit, nil)
			requireStatus(t, r, 409)
			if !strings.Contains(r.body, `value="10:00"`) || strings.Contains(r.body, `data-state="confirmed_demo"`) {
				t.Fatal("failure lost draft or claimed success")
			}
			if kind == "price" {
				if count(t, "bookings") != 0 || !strings.Contains(r.body, "2500.00 birr") {
					t.Fatal("stale price was silently accepted")
				}
			}
		})
	}
}

type faultGateway struct {
	Gateway
	afterCommit, beforeCommit, wrongPrice bool
}

func (g *faultGateway) Book(ctx context.Context, s Scope, in integrity.Intent, key string) (integrity.Result, error) {
	if g.beforeCommit {
		g.beforeCommit = false
		return integrity.Result{}, errors.New("synthetic precommit outage")
	}
	r, e := g.Gateway.Book(ctx, s, in, key)
	if e != nil {
		return r, e
	}
	if g.afterCommit {
		g.afterCommit = false
		return integrity.Result{}, errors.New("synthetic committed response lost")
	}
	if g.wrongPrice {
		g.wrongPrice = false
		r.PriceMinor++
	}
	return r, nil
}
func TestUnknownOutcomeRecoveryNeverCreatesAfterReviewExpiry(t *testing.T) {
	for _, kind := range []string{"committed", "not-committed", "wrong-price"} {
		t.Run(kind, func(t *testing.T) {
			gateway := &faultGateway{Gateway: SQLGateway{Pool: testDB.Runtime}, afterCommit: kind == "committed", beforeCommit: kind == "not-committed", wrongPrice: kind == "wrong-price"}
			f := reset(t, gateway)
			f.enter(t, "instant")
			submit := f.reviewed(t)
			r := f.request(t, "POST", "/v0/submissions", submit, nil)
			requireStatus(t, r, 503)
			if !strings.Contains(r.body, `data-state="checking_demo"`) {
				t.Fatal("unknown outcome called success")
			}
			r = f.request(t, "GET", "/continue", nil, nil)
			requireStatus(t, r, 303)
			if r.header.Get("Location") != "/review" {
				t.Fatal("unknown attempt offered new entry")
			}
			f.app.now = func() time.Time { return time.Now().Add(11 * time.Minute) }
			r = f.request(t, "POST", "/v0/submissions", submit, nil)
			if kind == "not-committed" {
				requireStatus(t, r, 503)
				if count(t, "bookings") != 0 {
					t.Fatal("expired unknown review created a new booking")
				}
			} else {
				requireStatus(t, r, 303)
				if count(t, "bookings") != 1 || count(t, "outbox") != 1 {
					t.Fatal("receipt recovery duplicated committed effects")
				}
			}
		})
	}
}
func TestConcurrentSubmissionReplaysOneBooking(t *testing.T) {
	f := reset(t, nil)
	f.enter(t, "instant")
	submit := f.reviewed(t)
	gate := make(chan struct{})
	codes := make(chan int, 8)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-gate
			r := f.request(t, "POST", "/v0/submissions", submit, nil)
			codes <- r.status
		}()
	}
	close(gate)
	wg.Wait()
	close(codes)
	for code := range codes {
		if code != 303 {
			t.Fatal("concurrent duplicate failed recovery")
		}
	}
	if count(t, "bookings") != 1 || count(t, "outbox") != 1 || count(t, "commands") != 1 {
		t.Fatal("concurrent HTTP submit duplicated effect")
	}
}
func TestActualLostWebResponseReplaysOneBooking(t *testing.T) {
	f := resetWithHandler(t, nil, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/v0/submissions" && r.Header.Get("X-Demo-Drop-Response") == "1" {
				rr := httptest.NewRecorder()
				next.ServeHTTP(rr, r) // Run the entire commit before dropping the actual HTTP connection.
				if rr.Code != 303 {
					t.Error("drop fixture did not reach committed redirect")
				}
				connection, _, err := w.(http.Hijacker).Hijack()
				if err != nil {
					t.Error("drop fixture could not close its own test connection")
					return
				}
				_ = connection.Close()
				return
			}
			next.ServeHTTP(w, r)
		})
	})
	f.enter(t, "instant")
	submit := f.reviewed(t)
	r, _ := http.NewRequest("POST", f.server.URL+"/v0/submissions", strings.NewReader(submit.Encode()))
	r.Header.Set("Origin", f.server.URL)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("X-Demo-Drop-Response", "1")
	if out, err := f.client.Do(r); err == nil {
		out.Body.Close()
		t.Fatal("response was not actually lost")
	}
	if count(t, "bookings") != 1 {
		t.Fatal("lost response fixture did not commit exactly once")
	}
	requireStatus(t, f.request(t, "POST", "/v0/submissions", submit, nil), 303)
	if count(t, "bookings") != 1 || count(t, "outbox") != 1 || count(t, "commands") != 1 {
		t.Fatal("HTTP retry duplicated a committed booking")
	}
}
func TestTenantBoundReviewsAndActualSQLReceiptScope(t *testing.T) {
	f := reset(t, nil)
	f.enter(t, "river")
	submit := f.reviewed(t)
	r := f.request(t, "GET", "/review", nil, nil)
	if !strings.Contains(r.body, "Demo River Studio") || !strings.Contains(r.body, "1200.00 birr") {
		t.Fatal("wrong tenant catalogue")
	}
	requireStatus(t, f.request(t, "POST", "/v0/submissions", submit, nil), 303)
	var org string
	if e := testDB.Admin.QueryRow(context.Background(), "SELECT organization_id FROM ubi_v0.bookings").Scan(&org); e != nil || org != "demo-org-b" {
		t.Fatal("booking escaped grant scope")
	}
	var review *Review
	for _, s := range f.app.sessions {
		review = s.Review
	}
	in := review.Intent
	in.OrganizationID = "demo-org-a"
	if _, e := f.app.gateway.Recover(context.Background(), Scope{"demo-org-b", "studio", "portrait"}, in, "web-"+review.ID); e == nil {
		t.Fatal("recovery accepted forged scope")
	}
	in = review.Intent
	in.CustomerLabel = "Demo Guest 03"
	if _, e := f.app.gateway.Recover(context.Background(), Scope{"demo-org-b", "studio", "portrait"}, in, "web-"+review.ID); e == nil {
		t.Fatal("recovery accepted changed command input")
	}
}
func TestCookiesSecurityHeadersAndNoScriptAssets(t *testing.T) {
	f := reset(t, nil)
	r := f.request(t, "GET", "/", nil, nil)
	safeHeaders(t, r)
	if strings.Contains(r.body, "<script") || strings.Contains(r.body, "https://") {
		t.Fatal("third-party/script dependency in core page")
	}
	cookies := (&http.Response{Header: r.header}).Cookies()
	if len(cookies) != 1 || !cookies[0].HttpOnly || cookies[0].Secure || cookies[0].Domain != "" || cookies[0].SameSite != http.SameSiteLaxMode || cookies[0].MaxAge != 1800 {
		t.Fatal("loopback demo cookie policy wrong")
	}
	a, e := New("https://127.0.0.1:9999", f.app.gateway)
	if e != nil {
		t.Fatal(e)
	}
	rr := httptest.NewRecorder()
	a.cookie(rr, "ubi_demo_session", randomSecret())
	if !rr.Result().Cookies()[0].Secure {
		t.Fatal("HTTPS mode did not set Secure")
	}
	for _, method := range []string{"GET", "HEAD"} {
		r = f.request(t, method, "/assets/ui.css", nil, nil)
		requireStatus(t, r, 200)
		safeHeaders(t, r)
		if method == "HEAD" && r.body != "" {
			t.Fatal("CSS HEAD body")
		}
	}
}
func TestProcessRestartFailsClosedRatherThanRestoringIdentity(t *testing.T) {
	f := reset(t, nil)
	f.enter(t, "instant")
	next, e := New(f.server.URL, f.app.gateway)
	if e != nil {
		t.Fatal(e)
	}
	r := httptest.NewRequest("GET", f.server.URL+"/continue", nil)
	u, _ := url.Parse(f.server.URL)
	for _, c := range f.client.Jar.Cookies(u) {
		r.AddCookie(c)
	}
	rr := httptest.NewRecorder()
	next.ServeHTTP(rr, r)
	if rr.Code != 401 {
		t.Fatal("restart restored an ephemeral session")
	}
	if count(t, "bookings") != 0 {
		t.Fatal("restart performed booking")
	}
}
func redact(output []byte) []byte {
	return regexp.MustCompile(`[A-Za-z0-9_-]{43}`).ReplaceAll(output, []byte("REDACTED_DEMO_SECRET"))
}
func TestSavedPostmanCustomerContract(t *testing.T) {
	f := resetWithHandler(t, nil, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				cookieCount := 0
				for _, c := range r.Cookies() {
					if c.Name == "ubi_demo_bootstrap" {
						cookieCount++
					}
				}
				t.Logf("Local request diagnostic: bootstrap_cookie_count=%d origin_matches=%t fetch_site=%q", cookieCount, r.Header.Get("Origin") == "http://"+r.Host, r.Header.Get("Sec-Fetch-Site"))
			}
			next.ServeHTTP(w, r)
		})
	})
	cli, e := filepath.Abs("../../.tooling/postman-runtime/node_modules/.bin/postman")
	if e != nil {
		t.Fatal(e)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, cli, "collection", "run", "postman/collections/ubi-customer-v0", "--env-var", "base_url="+f.server.URL, "--no-report-events", "--timeout", "30000", "--timeout-request", "5000", "--color", "off")
	out, e := cmd.CombinedOutput()
	clean := redact(out)
	_ = os.MkdirAll("artifacts", 0755)
	_ = os.WriteFile("artifacts/postman-run.txt", clean, 0644)
	if e != nil {
		t.Fatalf("Postman local collection failed: %v\n%s", e, clean)
	}
	t.Logf("Local Postman contract result (secrets redacted):\n%s", clean)
	if count(t, "bookings") != 1 || count(t, "outbox") != 1 || f.app.Requests != 1 {
		t.Fatal("contract run effects do not match instant/custom distinction")
	}
}
func TestBrowserNoJavaScriptAndNarrowLargeTextFlow(t *testing.T) {
	f := resetWithHandler(t, nil, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				cookieCount := 0
				for _, c := range r.Cookies() {
					if c.Name == "ubi_demo_bootstrap" {
						cookieCount++
					}
				}
				t.Logf("Browser form diagnostic: bootstrap_cookie_count=%d origin_matches=%t origin_is_null=%t fetch_site=%q", cookieCount, r.Header.Get("Origin") == "http://"+r.Host, r.Header.Get("Origin") == "null", r.Header.Get("Sec-Fetch-Site"))
			}
			next.ServeHTTP(w, r)
		})
	})
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "node", "../../.tooling/test-customer-browser.mjs", f.server.URL)
	out, e := cmd.CombinedOutput()
	if e != nil {
		t.Fatalf("browser checks failed: %v\n%s", e, redact(out))
	}
	t.Logf("%s", redact(out))
	if count(t, "bookings") != 1 || count(t, "outbox") != 1 || f.app.Requests != 1 {
		t.Fatal("browser paths did not match database outcomes")
	}
}
func TestPayloadBudgets(t *testing.T) {
	f := reset(t, nil)
	samples := map[string]int{}
	samples["public_html"] = len(f.request(t, "GET", "/", nil, nil).body)
	path := f.start(t, "instant")
	samples["preview_html"] = len(f.request(t, "GET", path, nil, nil).body)
	requireStatus(t, f.redeem(t, path), 303)
	samples["details_html"] = len(f.request(t, "GET", "/continue", nil, nil).body)
	submit := f.reviewed(t)
	samples["review_html"] = len(f.request(t, "GET", "/review", nil, nil).body)
	requireStatus(t, f.request(t, "POST", "/v0/submissions", submit, nil), 303)
	samples["outcome_html"] = len(f.request(t, "GET", "/outcome", nil, nil).body)
	samples["css"] = len(f.request(t, "GET", "/assets/ui.css", nil, nil).body)
	for name, bytes := range samples {
		if bytes > 16384 {
			t.Errorf("local uncompressed %s exceeds 16 KiB experiment budget: %d", name, bytes)
		}
	}
	raw, _ := json.MarshalIndent(map[string]any{"scope": "Uncompressed loopback body bytes only; not network-transfer size or speed", "bytes": samples, "javascript_bytes": 0, "external_asset_requests_required": 0}, "", "  ")
	_ = os.MkdirAll("artifacts", 0755)
	if e := os.WriteFile("artifacts/payload-receipt.json", raw, 0644); e != nil {
		t.Fatal(e)
	}
}
