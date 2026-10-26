package integrity

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"regexp"
	"time"
)

type Scope struct {
	OrganizationID string `json:"organization_id"`
	OperationID    string `json:"operation_id"`
}
type Intent struct {
	Scope
	ServiceID              string    `json:"service_id"`
	ExpectedServiceVersion int       `json:"expected_service_version"`
	StartAt                time.Time `json:"start_at"`
	CustomerLabel          string    `json:"customer_label"`
	Source                 string    `json:"source"`
}
type HoldInput struct {
	Intent
	TTLSeconds int `json:"ttl_seconds"`
}
type QuoteInput struct {
	Intent
	RequestID      string `json:"request_id"`
	DemoPriceMinor *int64 `json:"demo_price_minor,omitempty"`
}
type VersionInput struct {
	ExpectedVersion int `json:"expected_version"`
}
type RescheduleInput struct {
	VersionInput
	StartAt time.Time `json:"start_at"`
}
type Result struct {
	DemoOnly bool   `json:"demo_only"`
	ID       string `json:"id"`
	Scope
	Version        int        `json:"version"`
	State          string     `json:"state"`
	StartAt        time.Time  `json:"start_at"`
	EndAt          time.Time  `json:"end_at"`
	ServiceID      string     `json:"service_id"`
	ServiceVersion int        `json:"service_version"`
	PriceMinor     int64      `json:"price_minor"`
	Currency       string     `json:"currency"`
	CustomerLabel  string     `json:"customer_label"`
	Source         string     `json:"source"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
	RequestID      string     `json:"request_id,omitempty"`
	BookingID      string     `json:"booking_id,omitempty"`
}
type Receipt struct {
	Status int    `json:"status"`
	Result Result `json:"result"`
}
type Principal struct {
	ID, OrganizationID string
	Operations         map[string]bool
	Owner, Active      bool
}

func DemoPrincipal(id string) (Principal, bool) {
	switch id {
	case "demo-owner-a":
		return Principal{id, "demo-org-a", map[string]bool{"studio": true, "events": true}, true, true}, true
	case "demo-owner-b":
		return Principal{id, "demo-org-b", map[string]bool{"studio": true}, true, true}, true
	case "demo-staff-a":
		return Principal{id, "demo-org-a", map[string]bool{"studio": true}, false, true}, true
	case "demo-revoked":
		return Principal{id, "demo-org-a", map[string]bool{"studio": true}, false, false}, true
	default:
		return Principal{}, false
	}
}

type DomainError struct {
	Code, Detail string
	Status       int
	Retryable    bool
}

func (e *DomainError) Error() string { return e.Code }
func failure(code string, status int, detail string) error {
	return &DomainError{Code: code, Status: status, Detail: detail}
}
func authorize(p Principal, scope Scope, owner bool) error {
	if !p.Active || p.OrganizationID != scope.OrganizationID || !p.Operations[scope.OperationID] || (owner && !p.Owner) {
		return failure("permission_denied", 403, "This demo fixture cannot perform that action.")
	}
	return nil
}

var keyPattern = regexp.MustCompile(`^[A-Za-z0-9_.:-]{8,128}$`)
var requestPattern = regexp.MustCompile(`^demo-request-[a-z0-9-]{1,50}$`)

func validTime(t time.Time) bool {
	return !t.IsZero() && t.Year() >= 2029 && t.Year() <= 2031 && t.Second() == 0 && t.Nanosecond() == 0
}
func (in Intent) validate() error {
	labels := map[string]bool{"Demo Guest 01": true, "Demo Guest 02": true, "Demo Guest 03": true}
	sources := map[string]bool{"Call": true, "Walk-in": true, "WhatsApp": true, "Telegram": true, "Other": true}
	if !validTime(in.StartAt) || !labels[in.CustomerLabel] || !sources[in.Source] || in.ExpectedServiceVersion < 1 {
		return failure("invalid_input", 422, "Use the declared synthetic labels, service version and explicit-offset fixture time.")
	}
	return nil
}
func newID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	b[6] = (b[6] & 15) | 64
	b[8] = (b[8] & 63) | 128
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}
func jsonBytes(v any) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
