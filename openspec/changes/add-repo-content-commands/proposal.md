# Change: Add Repository Content Commands

## Why
The CLI foundation and auth check are complete, but the four core read-only commands (`stat`, `ls`, `cat`, `get`) that define the tool's primary value are missing.
Users currently cannot browse, read, or download repository content.

## What Changes
- Add `owner/repo` argument parsing with validation (exit code `13` on failure).
- Add GitHub API client methods for `GET /repos/{owner}/{repo}/contents/{path}` and `GET /repos/{owner}/{repo}/git/trees/{sha}`.
- Add a service layer (`internal/service/repo_service.go`) to orchestrate content lookup, recursive tree walking, file reading, and download with write strategy.
- Add `ghrepo stat <owner/repo> <path>` command with text and JSON output.
- Add `ghrepo ls <owner/repo> <path>` command with `--ref`, `--recursive`, text and JSON output.
- Add `ghrepo cat <owner/repo> <path>` command with `--ref`, outputting file content to stdout.
- Add `ghrepo get <owner/repo> <path> --out <local-path>` command with `--ref`, `--overwrite`, supporting single file and directory download.
- Add unified `Entry` model (`type/path/sha/size/download_url`) for output.

## Impact
- Affected specs: `repo-content-commands` (new capability)
- Affected code: `internal/githubapi/client.go` (new methods), `internal/service/repo_service.go` (new), `internal/cli/cmd_stat.go` (new), `internal/cli/cmd_ls.go` (new), `internal/cli/cmd_cat.go` (new), `internal/cli/cmd_get.go` (new), `internal/output/print.go` (new formatters), `internal/cli/root.go` (register commands)
- User-visible impact: users can browse, read, and download repository content through four new commands.
