package shopsavvy

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func errorServer(t *testing.T, status int, retryAfter, body string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if retryAfter != "" {
			w.Header().Set("Retry-After", retryAfter)
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return srv
}

const rateLimitBody = `{"success":false,"error":"Rate limit exceeded","error_detail":{"code":"ERR_RATE_LIMIT","message":"Rate limit exceeded","action":"Wait for Retry-After seconds, then retry.","docs_url":"https://shopsavvy.com/data/documentation/errors#err_rate_limit"},"meta":{"request_id":"req_rl_1"}}`
const creditsBody = `{"success":false,"error":"Credits exhausted","error_detail":{"code":"ERR_CREDITS_EXHAUSTED","message":"Credits exhausted","action":"Do not retry. Credits return when the billing period resets or the plan is upgraded.","docs_url":"https://shopsavvy.com/data/documentation/errors#err_credits_exhausted"},"meta":{"request_id":"req_cr_1"}}`

func TestRateLimitErrorCarriesDetailAndRetryAfter(t *testing.T) {
	srv := errorServer(t, 429, "7", rateLimitBody)
	c := newTestClient(t, srv.URL)
	_, err := c.GetProductDetails("611247373064")
	var rl *RateLimitError
	if !errors.As(err, &rl) {
		t.Fatalf("err = %T %v, want *RateLimitError", err, err)
	}
	if rl.Code != "ERR_RATE_LIMIT" || rl.StatusCode != 429 || rl.RequestID != "req_rl_1" {
		t.Errorf("code/status/request_id = %q/%d/%q", rl.Code, rl.StatusCode, rl.RequestID)
	}
	if rl.Action == "" || !strings.HasSuffix(rl.DocsURL, "#err_rate_limit") {
		t.Errorf("action/docs_url = %q/%q", rl.Action, rl.DocsURL)
	}
	if rl.RetryAfterSeconds == nil || *rl.RetryAfterSeconds != 7 {
		t.Errorf("RetryAfterSeconds = %v, want 7", rl.RetryAfterSeconds)
	}
}

func TestCreditsExhaustedIsDistinctFromRateLimit(t *testing.T) {
	srv := errorServer(t, 429, "", creditsBody)
	c := newTestClient(t, srv.URL)
	_, err := c.GetProductDetails("611247373064")
	var ce *CreditsExhaustedError
	if !errors.As(err, &ce) {
		t.Fatalf("err = %T %v, want *CreditsExhaustedError", err, err)
	}
	var rl *RateLimitError
	if errors.As(err, &rl) {
		t.Errorf("credits exhausted must not match *RateLimitError")
	}
	if ce.Code != "ERR_CREDITS_EXHAUSTED" || ce.StatusCode != 429 || ce.RequestID != "req_cr_1" {
		t.Errorf("code/status/request_id = %q/%d/%q", ce.Code, ce.StatusCode, ce.RequestID)
	}
	if ce.RetryAfterSeconds != nil {
		t.Errorf("RetryAfterSeconds = %v, want nil", *ce.RetryAfterSeconds)
	}
}

func TestPlain429WithoutCodeStaysRateLimit(t *testing.T) {
	srv := errorServer(t, 429, "", `{"success":false,"error":"Too many requests"}`)
	c := newTestClient(t, srv.URL)
	_, err := c.GetProductDetails("611247373064")
	var rl *RateLimitError
	if !errors.As(err, &rl) {
		t.Fatalf("err = %T %v, want *RateLimitError", err, err)
	}
	if rl.Code != "" || rl.RetryAfterSeconds != nil || rl.RequestID != "" {
		t.Errorf("expected empty details, got %+v", rl.ErrorDetails)
	}
}

func TestOtherErrorsCarryDetail(t *testing.T) {
	body := `{"success":false,"error":"Invalid API key","error_detail":{"code":"ERR_AUTH_INVALID","message":"Invalid API key","action":"Check your key.","docs_url":"https://shopsavvy.com/data/documentation/errors#err_auth_invalid"},"meta":{"request_id":"req_a_1"}}`
	srv := errorServer(t, 401, "", body)
	c := newTestClient(t, srv.URL)
	_, err := c.GetProductDetails("611247373064")
	var ae *AuthenticationError
	if !errors.As(err, &ae) {
		t.Fatalf("err = %T %v, want *AuthenticationError", err, err)
	}
	if ae.Code != "ERR_AUTH_INVALID" || ae.Action != "Check your key." || ae.RequestID != "req_a_1" {
		t.Errorf("details = %+v", ae.ErrorDetails)
	}

	srv2 := errorServer(t, 500, "", `{"success":false,"error":"boom","error_detail":{"code":"ERR_INTERNAL","message":"boom","action":"Retry later.","docs_url":"https://shopsavvy.com/data/documentation/errors#err_internal"}}`)
	c2 := newTestClient(t, srv2.URL)
	_, err = c2.GetProductDetails("611247373064")
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.Code != "ERR_INTERNAL" || apiErr.Message != "boom" {
		t.Fatalf("err = %#v", err)
	}
}

// Live check: a well-formed but bogus key is rejected for free with ERR_AUTH_REVOKED.
func TestLiveBogusKeyReturnsStructuredAuthError(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live check in -short mode")
	}
	c, err := NewClient("ss_live_" + strings.Repeat("a", 32))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	_, err = c.GetProductDetails("611247373064")
	var ae *AuthenticationError
	if !errors.As(err, &ae) {
		t.Fatalf("err = %T %v, want *AuthenticationError", err, err)
	}
	if ae.Code != "ERR_AUTH_REVOKED" || ae.DocsURL == "" || ae.StatusCode != 401 {
		t.Errorf("code/docs_url/status = %q/%q/%d", ae.Code, ae.DocsURL, ae.StatusCode)
	}
}
