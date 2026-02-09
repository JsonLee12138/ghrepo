## Context
The CLI currently supports only read operations (ls, stat, cat, get). This change adds write (create/update) and delete capabilities using the GitHub Contents API. Both operations are destructive and require confirmation prompts.

## Goals / Non-Goals
- Goals:
  - Support creating new files and updating existing files in repositories
  - Support deleting files from repositories
  - Interactive confirmation for all write/delete operations
  - `--yes` flag to skip confirmation for scripting/CI use
  - Consistent error handling and exit codes
  - Text and JSON dual output for mutation results
- Non-Goals:
  - Batch operations (multi-file create/delete in one command)
  - Directory creation/deletion (GitHub API is file-level only)
  - Branch creation (use existing branch or default branch)
  - Pull request creation

## Decisions

### Command naming: `put` and `rm`
- `put` — mirrors HTTP PUT semantics; handles both create and update
- `rm` — familiar Unix-style delete command
- Alternatives: `write`/`delete`, `create`+`update`/`remove` — rejected for verbosity

### Confirmation prompt
- Default: interactive prompt on stderr (`⚠ About to [action]. Continue? [y/N]`)
- Default answer is No (safe default)
- `--yes` / `-y` flag skips confirmation
- When stdin is not a terminal (pipe/script), require `--yes` or abort with exit code 17
- Confirmation message includes: action type, owner/repo, path, branch

### File content input for `put`
- `--file` flag: read content from a local file path
- `--stdin` flag: read content from stdin
- Exactly one of `--file` or `--stdin` must be specified
- Content is base64-encoded before sending to API

### Auto-detect create vs update for `put`
- Internally call `GET /repos/{owner}/{repo}/contents/{path}` first
- If file exists, get its SHA and perform update
- If 404, perform create
- User does not need to know or specify the SHA

### Commit message
- Required via `--message` / `-m` flag
- No default message — must be explicit

### Branch
- Optional `--branch` / `-b` flag
- If not specified, API defaults to the repo's default branch

## Risks / Trade-offs
- Race condition: file may change between SHA lookup and PUT — acceptable for CLI tool
- Non-terminal stdin detection: use `os.Stdin.Stat()` to check for pipe mode

## Open Questions
- None
