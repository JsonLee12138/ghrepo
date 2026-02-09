package config

import (
	"os"
	"testing"
)

func TestResolveToken_FlagTakesPrecedence(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "env-github")
	t.Setenv("GH_TOKEN", "env-gh")

	got := ResolveToken("flag-token")
	if got != "flag-token" {
		t.Errorf("expected flag-token, got %q", got)
	}
}

func TestResolveToken_GitHubTokenOverGHToken(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "env-github")
	t.Setenv("GH_TOKEN", "env-gh")

	got := ResolveToken("")
	if got != "env-github" {
		t.Errorf("expected env-github, got %q", got)
	}
}

func TestResolveToken_FallbackToGHToken(t *testing.T) {
	os.Unsetenv("GITHUB_TOKEN")
	t.Setenv("GH_TOKEN", "env-gh")

	got := ResolveToken("")
	if got != "env-gh" {
		t.Errorf("expected env-gh, got %q", got)
	}
}

func TestResolveToken_EmptyWhenNoneSet(t *testing.T) {
	os.Unsetenv("GITHUB_TOKEN")
	os.Unsetenv("GH_TOKEN")

	got := ResolveToken("")
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}
