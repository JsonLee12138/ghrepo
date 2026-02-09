package githubapi

import (
	"testing"

	clerrors "githubRAGCli/internal/exitcode"
)

func TestParseRepo_Valid(t *testing.T) {
	tests := []struct {
		input     string
		wantOwner string
		wantRepo  string
	}{
		{"octocat/hello-world", "octocat", "hello-world"},
		{"org/repo", "org", "repo"},
		{"a/b", "a", "b"},
	}
	for _, tt := range tests {
		owner, repo, err := ParseRepo(tt.input)
		if err != nil {
			t.Errorf("ParseRepo(%q): unexpected error: %v", tt.input, err)
			continue
		}
		if owner != tt.wantOwner || repo != tt.wantRepo {
			t.Errorf("ParseRepo(%q): got (%q, %q), want (%q, %q)", tt.input, owner, repo, tt.wantOwner, tt.wantRepo)
		}
	}
}

func TestParseRepo_Invalid(t *testing.T) {
	tests := []string{
		"",
		"noslash",
		"too/many/slashes",
		"/leadingslash",
		"trailingslash/",
		"/",
	}
	for _, input := range tests {
		_, _, err := ParseRepo(input)
		if err == nil {
			t.Errorf("ParseRepo(%q): expected error, got nil", input)
			continue
		}
		ce, ok := err.(*clerrors.CLIError)
		if !ok {
			t.Errorf("ParseRepo(%q): expected CLIError, got %T", input, err)
			continue
		}
		if ce.ExitCode() != clerrors.ExitBadArgs {
			t.Errorf("ParseRepo(%q): exit code got %d, want %d", input, ce.ExitCode(), clerrors.ExitBadArgs)
		}
	}
}
