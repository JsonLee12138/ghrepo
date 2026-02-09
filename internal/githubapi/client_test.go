package githubapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	clerrors "githubRAGCli/internal/exitcode"
)

func TestGetAuthenticatedUser_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/user" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("unexpected auth header: %s", got)
		}
		w.Header().Set("X-RateLimit-Remaining", "4999")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]string{"login": "octocat"})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "test-token", 5*time.Second)
	result, err := c.GetAuthenticatedUser()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Login != "octocat" {
		t.Errorf("login: got %q, want octocat", result.Login)
	}
	if result.RateLimitRemaining != 4999 {
		t.Errorf("rate limit remaining: got %d, want 4999", result.RateLimitRemaining)
	}
}

func TestGetAuthenticatedUser_401(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		w.Write([]byte(`{"message":"Bad credentials"}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "bad-token", 5*time.Second)
	_, err := c.GetAuthenticatedUser()
	if err == nil {
		t.Fatal("expected error")
	}
	ce, ok := err.(*clerrors.CLIError)
	if !ok {
		t.Fatalf("expected CLIError, got %T", err)
	}
	if ce.ExitCode() != clerrors.ExitAuthFailure {
		t.Errorf("exit code: got %d, want %d", ce.ExitCode(), clerrors.ExitAuthFailure)
	}
}

func TestGetAuthenticatedUser_403_Permission(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Remaining", "4000")
		w.WriteHeader(403)
		w.Write([]byte(`{"message":"Resource not accessible by integration"}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "scoped-token", 5*time.Second)
	_, err := c.GetAuthenticatedUser()
	if err == nil {
		t.Fatal("expected error")
	}
	ce, ok := err.(*clerrors.CLIError)
	if !ok {
		t.Fatalf("expected CLIError, got %T", err)
	}
	if ce.ExitCode() != clerrors.ExitPermission {
		t.Errorf("exit code: got %d, want %d", ce.ExitCode(), clerrors.ExitPermission)
	}
}

func TestGetAuthenticatedUser_403_RateLimit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Remaining", "0")
		w.WriteHeader(403)
		w.Write([]byte(`{"message":"API rate limit exceeded"}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "token", 5*time.Second)
	_, err := c.GetAuthenticatedUser()
	if err == nil {
		t.Fatal("expected error")
	}
	ce, ok := err.(*clerrors.CLIError)
	if !ok {
		t.Fatalf("expected CLIError, got %T", err)
	}
	if ce.ExitCode() != clerrors.ExitRateLimit {
		t.Errorf("exit code: got %d, want %d", ce.ExitCode(), clerrors.ExitRateLimit)
	}
}

func TestGetAuthenticatedUser_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(200)
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "token", 50*time.Millisecond)
	_, err := c.GetAuthenticatedUser()
	if err == nil {
		t.Fatal("expected error")
	}
	ce, ok := err.(*clerrors.CLIError)
	if !ok {
		t.Fatalf("expected CLIError, got %T", err)
	}
	if ce.ExitCode() != clerrors.ExitTransport {
		t.Errorf("exit code: got %d, want %d", ce.ExitCode(), clerrors.ExitTransport)
	}
}
