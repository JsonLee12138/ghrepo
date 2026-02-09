# Change: Add CLI Foundation and Auth Check

## Why
The repository currently has planning documents but no executable baseline.  
We need a minimal, testable vertical slice that establishes runtime configuration, authentication validation, and stable exit behavior before adding content commands (`ls`, `cat`, `get`, `stat`).

## What Changes
- Add the initial CLI foundation contract for global flags: `--token`, `--api-base`, `--timeout`, `--json`, `--verbose`.
- Add token resolution rules with deterministic precedence: flag, then `GITHUB_TOKEN`, then `GH_TOKEN`.
- Add `ghrepo auth check` behavior contract for text and JSON output.
- Add failure classification and standardized exit-code mapping for authentication and transport failures in `auth check`.
- Explicitly defer repository content operations (`ls`, `cat`, `get`, `stat`) to follow-up changes.

## Impact
- Affected specs: `cli-auth-bootstrap` (new capability)
- Affected code (planned): `cmd/ghrepo/main.go`, `internal/cli/root.go`, `internal/cli/cmd_auth.go`, `internal/config/config.go`, `internal/githubapi/client.go`, `internal/output/print.go`, `internal/errors/codes.go`
- User-visible impact: users can verify token validity and environment setup through one stable command.
