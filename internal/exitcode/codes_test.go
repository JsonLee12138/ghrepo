package exitcode

import (
	"testing"
)

func TestExitCode_Mapping(t *testing.T) {
	tests := []struct {
		cat  Category
		want int
	}{
		{CatAuthFailure, 10},
		{CatPermission, 11},
		{CatNotFound, 12},
		{CatBadArgs, 13},
		{CatTransport, 14},
		{CatRateLimit, 15},
		{CatLocalWriteErr, 16},
	}
	for _, tt := range tests {
		e := &CLIError{Cat: tt.cat, Message: "test"}
		if got := e.ExitCode(); got != tt.want {
			t.Errorf("Category %d: got %d, want %d", tt.cat, got, tt.want)
		}
	}
}

func TestClassifyHTTP(t *testing.T) {
	tests := []struct {
		status      int
		rateLimited bool
		wantCode    int
	}{
		{401, false, ExitAuthFailure},
		{403, false, ExitPermission},
		{403, true, ExitRateLimit},
		{404, false, ExitNotFound},
		{500, false, ExitTransport},
	}
	for _, tt := range tests {
		e := ClassifyHTTP(tt.status, tt.rateLimited, "body")
		if e.ExitCode() != tt.wantCode {
			t.Errorf("HTTP %d (rateLimit=%v): got exit %d, want %d",
				tt.status, tt.rateLimited, e.ExitCode(), tt.wantCode)
		}
	}
}

func TestCLIError_Error(t *testing.T) {
	e := NewAuthFailure("bad token", nil)
	if e.Error() != "bad token" {
		t.Errorf("got %q", e.Error())
	}
}
