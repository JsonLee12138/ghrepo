## 1. CLI Foundation
- [x] 1.1 Scaffold `ghrepo` entrypoint and root command registration.
- [x] 1.2 Add global flags `--token`, `--api-base`, `--timeout`, `--json`, `--verbose` with documented defaults.
- [x] 1.3 Implement config resolution with token precedence `--token > GITHUB_TOKEN > GH_TOKEN`.

## 2. Auth Check Slice
- [x] 2.1 Add GitHub API client method for authenticated user lookup (`GET /user`) with configurable base URL and timeout.
- [x] 2.2 Implement `ghrepo auth check` command using service/client abstraction (no direct HTTP in CLI layer).
- [x] 2.3 Add text and JSON output formatting for auth status, username, and rate-limit remaining.

## 3. Error Handling
- [x] 3.1 Introduce typed error categories and map to exit codes (`10`, `11`, `14`, `15` for this slice).
- [x] 3.2 Ensure a single process exit path in `main` that converts typed errors to `os.Exit(code)`.
- [x] 3.3 Ensure logs and verbose output never include token values.

## 4. Validation
- [x] 4.1 Unit tests: token resolution precedence, output serialization, and error-code mapping.
- [x] 4.2 Integration tests (`httptest`): auth success, 401 invalid token, 403 insufficient permission, 403 rate limit.
- [x] 4.3 Run `go test ./...`.
- [x] 4.4 Run `openspec validate add-cli-foundation-auth-check --strict`.
