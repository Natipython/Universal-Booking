package customerweb

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	integrity "ubi.local/booking-integrity-v0"
)

//go:embed views.html ui.css
var assets embed.FS

type Scope struct{ Organization, Operation, Service string }
type Catalog struct {
	Name              string
	Version, Duration int
	Price             int64
}
type Gateway interface {
	Snapshot(context.Context, Scope) (Catalog, error)
	Book(context.Context, Scope, integrity.Intent, string) (integrity.Result, error)
	Recover(context.Context, Scope, integrity.Intent, string) (integrity.Result, error)
}
type Draft struct{ Date, Time, Customer, EventSize, Calendar, Clock string }
type Grant struct {
	Scope                                            Scope
	Mode, Purpose, Source, BoundBrowser, SessionHash string
	Expires                                          time.Time
	Revoked                                          bool
	Initial                                          Draft
}
type Review struct {
	ID      string
	Draft   Draft
	Version int
	Catalog Catalog
	Expires time.Time
	Intent  integrity.Intent
}
type Outcome struct{ State, ID string }
type Session struct {
	GrantHash string
	Draft     Draft
	Version   int
	Review    *Review
	Outcome   *Outcome
	Checking  bool
}
type App struct {
	mu       sync.Mutex // Deliberately coarse demo lock; not a scalable production session store.
	origin   string
	secure   bool
	key      []byte
	grants   map[string]*Grant
	sessions map[string]*Session
	gateway  Gateway
	now      func() time.Time
	views    *template.Template
	Requests int
}

func randomSecret() string {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b[:])
}
func digest(v string) string { h := sha256.Sum256([]byte(v)); return fmt.Sprintf("%x", h[:]) }
func secretShape(v string) bool {
	b, e := base64.RawURLEncoding.DecodeString(v)
	return e == nil && len(b) == 32 && len(v) == 43
}
func (a *App) mac(parts ...string) string {
	h := hmac.New(sha256.New, a.key)
	for _, p := range parts {
		_, _ = h.Write([]byte(p))
		_, _ = h.Write([]byte{0})
	}
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
func New(origin string, gateway Gateway) (*App, error) {
	u, err := url.Parse(origin)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || (u.Hostname() != "127.0.0.1" && u.Hostname() != "localhost" && u.Hostname() != "::1") || u.Port() == "" || u.User != nil || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || u.Path != "" || gateway == nil {
		return nil, fmt.Errorf("only an explicit loopback demo origin and gateway are accepted")
	}
	key := make([]byte, 32)
	if _, err = rand.Read(key); err != nil {
		return nil, err
	}
	text, err := assets.ReadFile("views.html")
	if err != nil {
		return nil, err
	}
	views, err := template.New("page").Parse(string(text))
	if err != nil {
		return nil, err
	}
	return &App{origin: origin, secure: u.Scheme == "https", key: key, grants: map[string]*Grant{}, sessions: map[string]*Session{}, gateway: gateway, now: time.Now, views: views}, nil
}
func (a *App) mint(g Grant) string {
	token := randomSecret()
	a.grants[digest(token)] = &g
	return token
}
func initialDraft() Draft {
	return Draft{"2030-01-15", "09:00", "Demo Guest 01", "not-specified", "Gregorian", "international-24h"}
}
func scenario(name string) (Grant, bool) {
	g := Grant{Scope: Scope{"demo-org-a", "studio", "portrait"}, Mode: "instant", Purpose: "prepare_booking", Source: "Demo staff-entered conversation", Initial: initialDraft()}
	switch name {
	case "instant":
	case "custom":
		g.Scope.Operation = "events"
		g.Scope.Service = "consultation"
		g.Mode = "custom"
	case "river":
		g.Scope.Organization = "demo-org-b"
	default:
		return Grant{}, false
	}
	return g, true
}
func (a *App) cookie(w http.ResponseWriter, name, value string) {
	http.SetCookie(w, &http.Cookie{Name: name, Value: value, Path: "/", MaxAge: 1800, HttpOnly: true, Secure: a.secure, SameSite: http.SameSiteLaxMode})
}
func uniqueCookie(r *http.Request, name string) string {
	value, count := "", 0
	for _, c := range r.Cookies() {
		if c.Name == name {
			value = c.Value
			count++
		}
	}
	if count != 1 || !secretShape(value) {
		return ""
	}
	return value
}
func (a *App) bootstrap(w http.ResponseWriter, r *http.Request) string {
	b := uniqueCookie(r, "ubi_demo_bootstrap")
	if b == "" {
		b = randomSecret()
		a.cookie(w, "ubi_demo_bootstrap", b)
	}
	return a.mac("csrf-bootstrap", b)
}

type page struct {
	Kind, Title, Problem, Message, CSRF, Token, Business, Mode, Service, Price, When, ReviewID, State string
	Draft                                                                                             Draft
	DraftVersion                                                                                      int
	Duration                                                                                          int
}

func (a *App) render(w http.ResponseWriter, status int, p page) {
	var body bytes.Buffer
	if err := a.views.Execute(&body, p); err != nil {
		http.Error(w, "The local demo page could not be prepared. No success is confirmed by this response.", 500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(body.Bytes())
}
func (a *App) problem(w http.ResponseWriter, status int, code, message string) {
	a.render(w, status, page{Kind: "error", Title: "Let's try again", Problem: code, Message: message})
}
func (a *App) form(w http.ResponseWriter, r *http.Request, names ...string) (url.Values, bool) {
	// A no-referrer HTML navigation can send Origin: null (Fetch Standard 3.2).
	// Accept that browser case only with same-origin navigation metadata, in addition
	// to the exact Host, issued cookie and secret CSRF checks. Missing metadata fails closed.
	origin, site := r.Header.Get("Origin"), r.Header.Get("Sec-Fetch-Site")
	privateNavigation := origin == "null" && site == "same-origin" && r.Header.Get("Sec-Fetch-Mode") == "navigate" && r.Header.Get("Sec-Fetch-Dest") == "document"
	if len(r.Header.Values("Origin")) != 1 || (origin != a.origin && !privateNavigation) || (site != "" && site != "same-origin") {
		a.problem(w, 403, "origin_denied", "Open the original demo page and try again.")
		return nil, false
	}
	if strings.Split(r.Header.Get("Content-Type"), ";")[0] != "application/x-www-form-urlencoded" {
		a.problem(w, 415, "unsupported_media_type", "Use the form on this demo page.")
		return nil, false
	}
	if r.URL.RawQuery != "" {
		a.problem(w, 400, "invalid_form", "Form details cannot be supplied in the URL.")
		return nil, false
	}
	r.Body = http.MaxBytesReader(w, r.Body, 8192)
	if err := r.ParseForm(); err != nil {
		var limit *http.MaxBytesError
		if errors.As(err, &limit) {
			a.problem(w, 413, "body_too_large", "The demo form is too large.")
		} else {
			a.problem(w, 400, "invalid_form", "Use the original form and try again.")
		}
		return nil, false
	}
	allowed := map[string]bool{}
	for _, n := range names {
		allowed[n] = true
	}
	for k, v := range r.PostForm {
		if !allowed[k] || len(v) != 1 || len(v[0]) > 120 {
			a.problem(w, 400, "invalid_form", "Use only the fields on this demo form.")
			return nil, false
		}
	}
	for _, n := range names {
		if _, ok := r.PostForm[n]; !ok {
			a.problem(w, 400, "invalid_form", "A required demo field is missing.")
			return nil, false
		}
	}
	return r.PostForm, true
}
func (a *App) csrf(w http.ResponseWriter, value, expected string) bool {
	if !secretShape(value) || !hmac.Equal([]byte(value), []byte(expected)) {
		a.problem(w, 403, "csrf_invalid", "This page is no longer ready to submit. Reopen the original demo page.")
		return false
	}
	return true
}
func (a *App) session(w http.ResponseWriter, r *http.Request) (*Session, *Grant, string, bool) {
	value := uniqueCookie(r, "ubi_demo_session")
	s := a.sessions[digest(value)]
	if value == "" || s == nil {
		a.problem(w, 401, "session_required", "Continue with the original demo link. No booking was made by opening this page.")
		return nil, nil, "", false
	}
	g := a.grants[s.GrantHash]
	if g == nil || g.Revoked || !a.now().Before(g.Expires) || !validGrant(g) || g.SessionHash != digest(value) {
		a.problem(w, 410, "link_unavailable", "This demo link is unavailable. Return to the person who shared it. Do not assume a new booking is needed.")
		return nil, nil, "", false
	}
	return s, g, a.mac("csrf-session", value), true
}
func validScope(s Scope) bool {
	return ((s.Organization == "demo-org-a" || s.Organization == "demo-org-b") && s.Operation == "studio" && s.Service == "portrait") || (s.Organization == "demo-org-a" && s.Operation == "events" && s.Service == "consultation")
}
func validGrant(g *Grant) bool {
	return g.Purpose == "prepare_booking" && validScope(g.Scope) &&
		((g.Mode == "instant" && g.Scope.Service == "portrait") || (g.Mode == "custom" && g.Scope.Service == "consultation"))
}
func (a *App) details(w http.ResponseWriter, status int, s *Session, g *Grant, csrf, message string, draft Draft) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c, err := a.gateway.Snapshot(ctx, g.Scope)
	if err != nil {
		a.problem(w, 503, "catalog_unavailable", "The demo service details cannot be checked. Your saved draft remains unchanged.")
		return
	}
	p := page{Kind: "edit", Title: "Check your details", CSRF: csrf, Business: business(g.Scope), Service: c.Name, Price: money(c.Price), Mode: g.Mode, Draft: draft, DraftVersion: s.Version, Message: message}
	a.render(w, status, p)
}
func business(s Scope) string {
	if s.Organization == "demo-org-b" {
		return "Demo River Studio"
	}
	return "Demo Light Studio"
}
func money(minor int64) string { return fmt.Sprintf("%d.%02d birr (ETB)", minor/100, minor%100) }
func draftIntent(d Draft, scope Scope, version int) (integrity.Intent, bool) {
	t, err := time.Parse(time.RFC3339, d.Date+"T"+d.Time+":00+03:00")
	label := d.Customer == "Demo Guest 01" || d.Customer == "Demo Guest 02" || d.Customer == "Demo Guest 03"
	size := d.EventSize == "not-specified" || d.EventSize == "small" || d.EventSize == "medium" || d.EventSize == "large"
	valid := err == nil && t.Year() >= 2029 && t.Year() <= 2031 && t.Format("2006-01-02") == d.Date && t.Format("15:04") == d.Time && d.Calendar == "Gregorian" && d.Clock == "international-24h" && label && size
	return integrity.Intent{Scope: integrity.Scope{OrganizationID: scope.Organization, OperationID: scope.Operation}, ServiceID: scope.Service, ExpectedServiceVersion: version, StartAt: t, CustomerLabel: d.Customer, Source: "Other"}, valid
}

type headWriter struct{ http.ResponseWriter }

func (w headWriter) Write(p []byte) (int, error) { return len(p), nil }

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodHead {
		w = headWriter{w}
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Robots-Tag", "noindex, nofollow, noarchive")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; form-action 'self'; base-uri 'none'; frame-ancestors 'none'")
	u, _ := url.Parse(a.origin)
	if r.Host != u.Host {
		a.problem(w, 403, "host_denied", "Use the configured local demo address.")
		return
	}
	if r.URL.RawQuery != "" || r.URL.ForceQuery {
		a.problem(w, 400, "invalid_url", "Private demo details cannot be supplied in the URL. Use the original link.")
		return
	}
	if r.Method != "GET" && r.Method != "HEAD" && r.Method != "POST" {
		a.problem(w, 405, "method_not_allowed", "Use the demo page controls.")
		return
	}
	path := r.URL.Path
	if r.Method == "GET" || r.Method == "HEAD" {
		if path == "/assets/ui.css" {
			css, _ := assets.ReadFile("ui.css")
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
			if r.Method != "HEAD" {
				_, _ = w.Write(css)
			}
			return
		}
		if path == "/" {
			a.render(w, 200, page{Kind: "public", Title: "Book with a familiar business", CSRF: a.bootstrap(w, r)})
			return
		}
		if strings.HasPrefix(path, "/h/") {
			token := strings.TrimPrefix(path, "/h/")
			if !secretShape(token) || r.URL.RawQuery != "" {
				a.problem(w, 404, "not_found", "That demo page is unavailable.")
				return
			}
			a.render(w, 200, page{Kind: "gate", Title: "Continue your booking", CSRF: a.bootstrap(w, r), Token: token})
			return
		}
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if r.Method == "POST" && (path == "/demo/starts" || path == "/v0/handoff-sessions") {
		field := "scenario"
		if path == "/v0/handoff-sessions" {
			field = "token"
		}
		f, ok := a.form(w, r, "csrf", field)
		if !ok {
			return
		}
		b := uniqueCookie(r, "ubi_demo_bootstrap")
		if b == "" || !a.csrf(w, f.Get("csrf"), a.mac("csrf-bootstrap", b)) {
			if b == "" {
				a.problem(w, 403, "csrf_invalid", "Open the original demo page first.")
			}
			return
		}
		if path == "/demo/starts" {
			g, ok := scenario(f.Get("scenario"))
			if !ok {
				a.problem(w, 422, "invalid_scenario", "Choose one of the invented demo scenarios.")
				return
			}
			if len(a.grants) >= 200 {
				a.problem(w, 503, "demo_limit", "Restart the disposable demo to reset its bounded data.")
				return
			}
			g.Expires = a.now().Add(30 * time.Minute)
			token := a.mint(g)
			http.Redirect(w, r, "/h/"+token, 303)
			return
		}
		token := f.Get("token")
		gh := digest(token)
		g := a.grants[gh]
		if !secretShape(token) || g == nil || g.Revoked || !a.now().Before(g.Expires) || !validGrant(g) || (g.BoundBrowser != "" && g.BoundBrowser != digest(b)) {
			a.problem(w, 410, "link_unavailable", "This demo link is unavailable. Ask the person who shared it for help.")
			return
		}
		value := a.mac("session", gh, b)
		if g.BoundBrowser == "" {
			g.BoundBrowser = digest(b)
			g.SessionHash = digest(value)
			a.sessions[g.SessionHash] = &Session{GrantHash: gh, Draft: g.Initial, Version: 1}
		}
		a.cookie(w, "ubi_demo_session", value)
		http.Redirect(w, r, "/continue", 303)
		return
	}
	if path != "/continue" && path != "/review" && path != "/outcome" && path != "/v0/reviews" && path != "/v0/submissions" {
		a.problem(w, 404, "not_found", "That demo page is unavailable.")
		return
	}
	s, g, csrf, ok := a.session(w, r)
	if !ok {
		return
	}
	if r.Method != "POST" {
		switch path {
		case "/continue":
			if s.Outcome != nil {
				http.Redirect(w, r, "/outcome", 303)
				return
			}
			if s.Checking {
				http.Redirect(w, r, "/review", 303)
				return
			}
			a.details(w, 200, s, g, csrf, "", s.Draft)
		case "/review":
			if s.Outcome != nil {
				http.Redirect(w, r, "/outcome", 303)
				return
			}
			if s.Review == nil {
				a.problem(w, 409, "review_required", "Review your details before submitting. Nothing new was booked.")
				return
			}
			a.showReview(w, s, g, csrf, "")
		case "/outcome":
			if s.Outcome == nil {
				a.problem(w, 409, "outcome_unknown", "No completed demo outcome is available. Return to the review and retry the same attempt.")
				return
			}
			a.render(w, 200, page{Kind: "outcome", Title: "Your demo result", State: s.Outcome.State, Business: business(g.Scope), Service: s.Review.Catalog.Name, When: when(s.Review.Intent.StartAt), Price: money(s.Review.Catalog.Price), Draft: s.Review.Draft})
		default:
			a.problem(w, 405, "method_not_allowed", "Submit using the form, not a page link.")
		}
		return
	}
	if path == "/v0/reviews" {
		f, ok := a.form(w, r, "csrf", "expected_draft_version", "service_id", "date", "time", "calendar", "clock", "customer_label", "event_size")
		if !ok || !a.csrf(w, f.Get("csrf"), csrf) {
			return
		}
		if s.Outcome != nil || s.Checking {
			a.problem(w, 409, "attempt_in_progress", "Recover the existing demo attempt; do not start another booking.")
			return
		}
		version, _ := strconv.Atoi(f.Get("expected_draft_version"))
		if version != s.Version || f.Get("service_id") != g.Scope.Service {
			a.problem(w, 409, "draft_changed", "This form belongs to different or older details. Reopen your current draft.")
			return
		}
		d := Draft{f.Get("date"), f.Get("time"), f.Get("customer_label"), f.Get("event_size"), f.Get("calendar"), f.Get("clock")}
		c, err := a.gateway.Snapshot(r.Context(), g.Scope)
		if err != nil {
			a.problem(w, 503, "catalog_unavailable", "Service terms cannot be checked. Your saved draft remains unchanged.")
			return
		}
		intent, valid := draftIntent(d, g.Scope, c.Version)
		if !valid {
			a.details(w, 422, s, g, csrf, "Check the demo date, time and choices. Use Gregorian dates and the international clock.", d)
			return
		}
		s.Draft = d
		s.Version++
		s.Review = &Review{ID: randomSecret(), Draft: d, Version: s.Version, Catalog: c, Intent: intent, Expires: a.now().Add(10 * time.Minute)}
		http.Redirect(w, r, "/review", 303)
		return
	}
	if path != "/v0/submissions" {
		a.problem(w, 405, "method_not_allowed", "Use the form controls for changes.")
		return
	}
	f, ok := a.form(w, r, "csrf", "review_id")
	if !ok || !a.csrf(w, f.Get("csrf"), csrf) {
		return
	}
	v := s.Review
	if v == nil || f.Get("review_id") != v.ID || v.Version != s.Version {
		a.problem(w, 409, "review_changed", "These are not your current reviewed details. Reopen the draft and review again.")
		return
	}
	if s.Outcome != nil {
		http.Redirect(w, r, "/outcome", 303)
		return
	}
	// An unknown committed outcome can be reconciled with the same key after review expiry,
	// while the grant/session remains valid. It must never trigger a newly reviewed attempt.
	if !s.Checking && !a.now().Before(v.Expires) {
		s.Review = nil
		a.details(w, 409, s, g, csrf, "Your review expired. Review the current details again; nothing new was confirmed.", s.Draft)
		return
	}
	if g.Mode == "custom" {
		a.Requests++
		s.Outcome = &Outcome{State: "request_saved_demo", ID: randomSecret()}
		http.Redirect(w, r, "/outcome", 303)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	var result integrity.Result
	var err error
	if s.Checking && !a.now().Before(v.Expires) {
		// Lookup only: an expired review must never authorize a new commitment.
		result, err = a.gateway.Recover(ctx, g.Scope, v.Intent, "web-"+v.ID)
	} else {
		result, err = a.gateway.Book(ctx, g.Scope, v.Intent, "web-"+v.ID)
	}
	if err != nil {
		var d *integrity.DomainError
		if errors.As(err, &d) && !d.Retryable && d.Status < 500 {
			s.Checking = false
			s.Review = nil
			a.details(w, d.Status, s, g, csrf, "That demo time or service changed. Your details are kept; review again before submitting.", s.Draft)
			return
		}
		s.Checking = true
		a.showReview(w, s, g, csrf, "The demo booking outcome is still being checked. Retry this same attempt; do not start another booking.")
		return
	}
	if !result.DemoOnly || result.State != "confirmed_demo" || result.ID == "" || result.OrganizationID != g.Scope.Organization || result.OperationID != g.Scope.Operation || result.ServiceID != g.Scope.Service || result.ServiceVersion != v.Catalog.Version || result.PriceMinor != v.Catalog.Price || result.Currency != "ETB" || result.CustomerLabel != v.Intent.CustomerLabel || !result.StartAt.Equal(v.Intent.StartAt) || !result.EndAt.Equal(v.Intent.StartAt.Add(time.Duration(v.Catalog.Duration)*time.Minute)) {
		s.Checking = true
		a.showReview(w, s, g, csrf, "The demo result cannot be verified. Retry this same attempt or stop for review.")
		return
	}
	s.Checking = false
	s.Outcome = &Outcome{State: result.State, ID: result.ID}
	http.Redirect(w, r, "/outcome", 303)
}
func when(t time.Time) string {
	return t.In(time.FixedZone("Addis Ababa", 3*3600)).Format("Monday, 02 January 2006 · 15:04 (3:04 PM)") + " · Gregorian · international clock · UTC+03:00"
}
func (a *App) showReview(w http.ResponseWriter, s *Session, g *Grant, csrf, message string) {
	status := 200
	if s.Checking {
		status = 503
	}
	a.render(w, status, page{Kind: "review", Title: "Review before sending", CSRF: csrf, Business: business(g.Scope), Service: s.Review.Catalog.Name, Price: money(s.Review.Catalog.Price), When: when(s.Review.Intent.StartAt), Mode: g.Mode, ReviewID: s.Review.ID, Message: message, Draft: s.Review.Draft, Duration: s.Review.Catalog.Duration, State: func() string {
		if s.Checking {
			return "checking_demo"
		}
		return "review_demo"
	}()})
}
