# Change: Add write and delete commands for repository files

## Why
The CLI is currently read-only. Users need the ability to create/update files and delete files in GitHub repositories directly from the command line, with safety confirmation prompts to prevent accidental modifications.

## What Changes
- Add `put` command: create or update a file in a GitHub repository via Contents API (`PUT /repos/{owner}/{repo}/contents/{path}`)
- Add `rm` command: delete a file from a GitHub repository via Contents API (`DELETE /repos/{owner}/{repo}/contents/{path}`)
- Both commands require interactive confirmation before executing (with `--yes` / `-y` flag to skip for scripting)
- Add `doJSON` (PUT/DELETE with JSON body) method to `githubapi.Client`
- Add `CreateOrUpdateFile` and `DeleteFile` methods to `service.RepoService`
- Add `MutationResult` output formatting to `output` package
- Add new exit code `ExitUserAbort = 17` for user-cancelled confirmation
- Update root command description from "Read-only CLI" to "CLI for GitHub repository contents"

## Impact
- Affected specs: repo-write-commands (new capability)
- Affected code:
  - `internal/githubapi/client.go` — add PUT/DELETE HTTP methods
  - `internal/service/repo_service.go` — add write/delete business logic
  - `internal/cli/cmd_put.go` — new command
  - `internal/cli/cmd_rm.go` — new command
  - `internal/cli/root.go` — register new commands, update description
  - `internal/cli/helpers.go` — add confirmation prompt helper
  - `internal/output/print.go` — add mutation result output
  - `internal/exitcode/codes.go` — add user-abort exit code
