---
name: ghrepo
description: >
  CLI tool for GitHub repository contents - browse, inspect, download, create, update, and
  delete files in any GitHub repository without cloning. Use when the agent needs to:
  (1) List directory contents of a remote GitHub repo,
  (2) Read/cat file contents from GitHub,
  (3) Download files or directories from GitHub to local filesystem,
  (4) Check file/directory metadata (type, SHA, size),
  (5) Create or update files in a GitHub repo (with commit message),
  (6) Delete files from a GitHub repo (with commit message),
  (7) Verify GitHub token authentication,
  (8) Initialize a project with AGENTS.md for a target repo.
  Triggers: "github repo", "remote file", "ghrepo", "browse repo", "download from github",
  "upload to github", "delete from github", "repo contents", "init repo".
---

# ghrepo - GitHub Repository Contents CLI

## Prerequisites

- `ghrepo` binary in PATH (install via `brew tap JsonLee12138/ghrepo && brew install ghrepo`, or `go install`)
- A GitHub token set via `GITHUB_TOKEN`, `GH_TOKEN`, or `--token` flag
- For write/delete operations: token needs `Contents: Read and Write` scope

## Quick Reference

```bash
ghrepo init <owner/repo>                                    # create AGENTS.md
ghrepo auth check                                          # verify token
ghrepo ls <owner/repo> <path> [--ref <ref>] [--recursive]  # list directory
ghrepo stat <owner/repo> <path> [--ref <ref>]               # file/dir metadata
ghrepo cat <owner/repo> <path> [--ref <ref>]                # output file content
ghrepo get <owner/repo> <path> --out <local> [--overwrite]  # download to local
ghrepo put <owner/repo> <path> -m <msg> --file <local> [-b <branch>] [-y] # create/update
ghrepo rm <owner/repo> <path> -m <msg> [-b <branch>] [-y]  # delete file
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--token <t>` | GitHub PAT (overrides env vars) |
| `--api-base <url>` | Custom API URL (for GHES) |
| `--timeout <dur>` | HTTP timeout (default `15s`) |
| `--json` | Structured JSON output |
| `--verbose` | Debug logging (never prints token) |

## Commands

### Initialize Project

```bash
ghrepo init owner/repo
```

Creates `AGENTS.md` in the current directory with the target repo and command reference. No token required. Fails if `AGENTS.md` already exists.

### Read Operations

**List directory:**
```bash
ghrepo ls octocat/Hello-World src/
ghrepo ls octocat/Hello-World src/ --recursive --json
```

**Read file to stdout:**
```bash
ghrepo cat octocat/Hello-World README.md
ghrepo cat octocat/Hello-World README.md --ref v1.0 > local.md
```

**Get metadata:**
```bash
ghrepo stat octocat/Hello-World README.md --json
```

**Download files/directories:**
```bash
ghrepo get octocat/Hello-World README.md --out ./README.md
ghrepo get octocat/Hello-World docs/ --out ./local-docs --overwrite
```

### Write Operations (require confirmation)

**Create or update a file:**
```bash
# From local file (interactive confirmation)
ghrepo put owner/repo path/to/file.txt -m "add file" --file ./local.txt

# From stdin (requires --yes since stdin is consumed)
echo "content" | ghrepo put owner/repo file.txt -m "create" --stdin --yes

# Skip confirmation + specify branch
ghrepo put owner/repo config.yml -m "update" --file ./config.yml -b develop --yes
```

**Delete a file:**
```bash
ghrepo rm owner/repo old-file.txt -m "remove file"
ghrepo rm owner/repo old-file.txt -m "cleanup" -b main --yes
```

### Auth

```bash
ghrepo auth check
ghrepo auth check --json
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 10 | Auth failure (missing/invalid token) |
| 11 | Permission denied |
| 12 | Not found |
| 13 | Bad arguments |
| 14 | Network/timeout error |
| 15 | Rate limited |
| 16 | Local write failure |
| 17 | User cancelled operation |

## Key Behaviors

- `put` auto-detects create vs update (no need to specify SHA)
- `put` and `rm` show confirmation prompt by default; use `--yes`/`-y` to skip
- `put` requires exactly one of `--file` or `--stdin` for content source
- `put` and `rm` require `-m`/`--message` for commit message
- `--ref` works on all read commands (branch, tag, or SHA)
- `get` preserves directory structure when downloading folders
- All commands respect `--json` for machine-readable output

For detailed usage docs, see [references/commands.md](references/commands.md).
