package integrity

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
)

type Problem struct {
	Type          string `json:"type"`
	Title         string `json:"title"`
	Status        int    `json:"status"`
	Detail        string `json:"detail"`
	Code          string `json:"code"`
	CorrelationID string `json:"correlation_id"`
	Retryable     bool   `json:"retryable"`
	DemoOnly      bool   `json:"demo_only"`
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
func writeProblem(w http.ResponseWriter, err error) {
	var d *DomainError
	if !errors.As(translate(err), &d) {
		panic("untranslated error")
	}
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(d.Status)
	_ = json.NewEncoder(w).Encode(Problem{"urn:ubi:demo:problem:" + d.Code, "Synthetic request could not complete", d.Status, d.Detail, d.Code, w.Header().Get("X-Request-ID"), d.Retryable, true})
}
func decode(w http.ResponseWriter, r *http.Request, body any) error {
	typ, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || typ != "application/json" {
		return failure("unsupported_media_type", 415, "Send application/json for this synthetic command.")
	}
	r.Body = http.MaxBytesReader(w, r.Body, 32768)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err = dec.Decode(body); err != nil {
		var size *http.MaxBytesError
		if errors.As(err, &size) {
			return failure("body_too_large", 413, "Demo JSON is limited to 32 KiB.")
		}
		return failure("invalid_json", 400, "Use valid bounded JSON with declared fields and explicit-offset times.")
	}
	var rest any
	if err = dec.Decode(&rest); err != io.EOF {
		return failure("invalid_json", 400, "Only one JSON object is allowed.")
	}
	return nil
}

// TEST-ONLY. The caller must bind this fixture-selecting handler to loopback only.
func (s *Service) DemoHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v0/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"demo_only": true, "status": "synthetic_only"})
	})
	register := func(path string, fn func(http.ResponseWriter, *http.Request, Principal) (Receipt, error)) {
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			p, ok := DemoPrincipal(r.Header.Get("X-Ubi-Demo-Principal"))
			if !ok {
				writeProblem(w, failure("demo_principal_required", 401, "Select an invented demo principal; this is not a real login."))
				return
			}
			if !p.Active {
				writeProblem(w, failure("permission_denied", 403, "This demo principal is inactive."))
				return
			}
			result, err := fn(w, r, p)
			if err != nil {
				writeProblem(w, err)
				return
			}
			if result.Status == 201 && result.Result.State == "confirmed_demo" {
				w.Header().Set("Location", "/v0/bookings/"+result.Result.ID)
			}
			writeJSON(w, result.Status, result.Result)
		})
	}
	key := func(r *http.Request) string { return r.Header.Get("Idempotency-Key") }
	register("POST /v0/bookings", func(w http.ResponseWriter, r *http.Request, p Principal) (Receipt, error) {
		var in Intent
		if err := decode(w, r, &in); err != nil {
			return Receipt{}, err
		}
		return s.CreateBooking(r.Context(), p, key(r), in)
	})
	register("GET /v0/bookings/{id}", func(w http.ResponseWriter, r *http.Request, p Principal) (Receipt, error) {
		out, err := s.GetBooking(r.Context(), p, r.PathValue("id"))
		return Receipt{200, out}, err
	})
	register("POST /v0/bookings/{id}/reschedules", func(w http.ResponseWriter, r *http.Request, p Principal) (Receipt, error) {
		var in RescheduleInput
		if err := decode(w, r, &in); err != nil {
			return Receipt{}, err
		}
		return s.Reschedule(r.Context(), p, key(r), r.PathValue("id"), in)
	})
	register("POST /v0/bookings/{id}/cancellations", func(w http.ResponseWriter, r *http.Request, p Principal) (Receipt, error) {
		var in VersionInput
		if err := decode(w, r, &in); err != nil {
			return Receipt{}, err
		}
		return s.Cancel(r.Context(), p, key(r), r.PathValue("id"), in)
	})
	register("POST /v0/holds", func(w http.ResponseWriter, r *http.Request, p Principal) (Receipt, error) {
		var in HoldInput
		if err := decode(w, r, &in); err != nil {
			return Receipt{}, err
		}
		return s.CreateHold(r.Context(), p, key(r), in)
	})
	register("POST /v0/holds/{id}/bookings", func(w http.ResponseWriter, r *http.Request, p Principal) (Receipt, error) {
		var in VersionInput
		if err := decode(w, r, &in); err != nil {
			return Receipt{}, err
		}
		return s.ConsumeHold(r.Context(), p, key(r), r.PathValue("id"), in)
	})
	register("POST /v0/quotations", func(w http.ResponseWriter, r *http.Request, p Principal) (Receipt, error) {
		var in QuoteInput
		if err := decode(w, r, &in); err != nil {
			return Receipt{}, err
		}
		return s.IssueQuote(r.Context(), p, key(r), in)
	})
	register("POST /v0/quotations/{id}/acceptances", func(w http.ResponseWriter, r *http.Request, p Principal) (Receipt, error) {
		var in VersionInput
		if err := decode(w, r, &in); err != nil {
			return Receipt{}, err
		}
		return s.AcceptQuote(r.Context(), p, key(r), r.PathValue("id"), in)
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Request-ID", newID())
		mux.ServeHTTP(w, r)
	})
}
