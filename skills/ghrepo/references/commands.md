# ghrepo Command Reference

## Token Configuration

Priority order:
1. `--token` flag (highest)
2. `GITHUB_TOKEN` environment variable
3. `GH_TOKEN` environment variable

```bash
# Verify before running commands
ghrepo auth check
```

## ls - List Directory

```bash
ghrepo ls <owner/repo> <path> [flags]
```

| Flag | Description |
|------|-------------|
| `--ref <ref>` | Git ref (branch/tag/SHA) |
| `--recursive` | List full subtree via Git Trees API |

- Path `.` or empty = repo root
- Errors if path is a file (use `cat` or `stat` instead)
- `--recursive` may be truncated by GitHub API for very large trees

## stat - File/Directory Metadata

```bash
ghrepo stat <owner/repo> <path> [--ref <ref>]
```

Returns: `type`, `path`, `sha`, `size`, `download_url` (files only).

## cat - Read File Content

```bash
ghrepo cat <owner/repo> <path> [--ref <ref>]
```

- Outputs raw file content to stdout
- Errors on directory paths
- Pipe to file: `ghrepo cat owner/repo file > local`

## get - Download

```bash
ghrepo get <owner/repo> <path> --out <local-path> [flags]
```

| Flag | Description |
|------|-------------|
| `--ref <ref>` | Git ref |
| `--out <path>` | Local output path (**required**) |
| `--overwrite` | Overwrite existing local files |

- Single file: downloads to exact `--out` path
- Directory: preserves internal structure under `--out`
- Without `--overwrite`, exits with code 16 if file exists

## put - Create or Update File

```bash
ghrepo put <owner/repo> <path> -m <msg> (--file <path> | --stdin) [flags]
```

| Flag | Description |
|------|-------------|
| `-m`, `--message` | Commit message (**required**) |
| `--file <path>` | Read content from local file |
| `--stdin` | Read content from stdin |
| `-b`, `--branch` | Target branch (optional) |
| `-y`, `--yes` | Skip confirmation prompt |

- Auto-detects create vs update (fetches current SHA internally)
- Exactly one of `--file` or `--stdin` required
- Shows confirmation prompt unless `--yes` is set
- Non-interactive sessions (piped input without `--stdin`) require `--yes`
- Output: `action` (created/updated), `path`, `sha` (commit), `branch`

## rm - Delete File

```bash
ghrepo rm <owner/repo> <path> -m <msg> [flags]
```

| Flag | Description |
|------|-------------|
| `-m`, `--message` | Commit message (**required**) |
| `-b`, `--branch` | Target branch (optional) |
| `-y`, `--yes` | Skip confirmation prompt |

- Fetches file SHA internally before deleting
- Cannot delete directories (GitHub API is file-level only)
- Shows confirmation prompt unless `--yes` is set
- Output: `action` (deleted), `path`, `sha` (commit), `branch`

## Common Patterns

### Browse then download
```bash
ghrepo ls owner/repo src/ --recursive
ghrepo cat owner/repo src/main.go
ghrepo get owner/repo src/ --out ./local-src
```

### Create a file from generated content
```bash
echo '{"key": "value"}' | ghrepo put owner/repo config.json -m "add config" --stdin --yes
```

### Update and verify
```bash
ghrepo put owner/repo README.md -m "update readme" --file ./README.md --yes
ghrepo cat owner/repo README.md
```

### Scripted batch operations with JSON output
```bash
ghrepo ls owner/repo docs/ --json | jq -r '.[].path'
ghrepo put owner/repo docs/new.md -m "add doc" --file ./new.md --yes --json
```
