## ADDED Requirements

### Requirement: Global CLI Runtime Options
The `ghrepo` CLI SHALL expose global runtime flags for authentication and transport behavior: `--token`, `--api-base`, `--timeout`, `--json`, and `--verbose`.

#### Scenario: Default runtime values are applied
- **WHEN** the user runs `ghrepo auth check` without global flags
- **THEN** the CLI uses `https://api.github.com` as `api-base`
- **AND** the CLI uses a default HTTP timeout of `15s`

#### Scenario: User overrides runtime values
- **WHEN** the user passes `--api-base` or `--timeout`
- **THEN** the CLI applies those values to subsequent API requests in that invocation

### Requirement: Token Resolution Precedence
The runtime configuration resolver SHALL resolve the token in this precedence order: `--token`, then `GITHUB_TOKEN`, then `GH_TOKEN`.

#### Scenario: Flag token overrides environment variables
- **GIVEN** both `GITHUB_TOKEN` and `GH_TOKEN` are set
- **WHEN** the user passes `--token`
- **THEN** the CLI uses the `--token` value for authentication

#### Scenario: Environment fallback is deterministic
- **GIVEN** `--token` is not provided and `GITHUB_TOKEN` is set
- **WHEN** the user runs `ghrepo auth check`
- **THEN** the CLI uses `GITHUB_TOKEN` and ignores `GH_TOKEN`

#### Scenario: Missing token is an authentication failure
- **GIVEN** no token is provided from flag or environment
- **WHEN** the user runs `ghrepo auth check`
- **THEN** the command fails with exit code `10`

### Requirement: Auth Check Command Contract
The system SHALL provide `ghrepo auth check` to validate the configured token against the GitHub API and report identity and rate-limit status.

#### Scenario: Valid credentials return success output
- **WHEN** `GET /user` returns success
- **THEN** the command exits with code `0`
- **AND** text output includes authentication success, username, and rate-limit remaining count

#### Scenario: JSON output uses stable keys
- **WHEN** the user runs `ghrepo auth check --json`
- **THEN** the command emits valid JSON
- **AND** JSON includes stable fields for auth status, username, and rate-limit remaining

### Requirement: Auth Failure and Transport Exit Codes
For `ghrepo auth check`, the CLI SHALL map failures to standardized exit codes defined in project documentation.

#### Scenario: Invalid token maps to authentication failure
- **WHEN** the GitHub API returns HTTP `401`
- **THEN** the command exits with code `10`

#### Scenario: Permission denied maps to authorization failure
- **WHEN** the GitHub API returns HTTP `403` due to insufficient permission
- **THEN** the command exits with code `11`

#### Scenario: Rate limiting maps to rate-limit failure
- **WHEN** the GitHub API returns HTTP `403` and indicates rate limiting
- **THEN** the command exits with code `15`

#### Scenario: Network or timeout error maps to transport failure
- **WHEN** the request fails because of timeout or network transport error
- **THEN** the command exits with code `14`
