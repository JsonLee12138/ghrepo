## 1. API Client Extension
- [x] 1.1 Add `GetContents(owner, repo, path, ref)` method returning raw JSON (supports file and directory responses).
- [x] 1.2 Add `GetTree(owner, repo, sha, recursive)` method for `GET /repos/{owner}/{repo}/git/trees/{sha}`.
- [x] 1.3 Add `owner/repo` parsing utility with exit code `13` on invalid format.

## 2. Service Layer
- [x] 2.1 Create `internal/service/repo_service.go` with `Entry` model and `RepoRef` struct.
- [x] 2.2 Implement `Stat(ref, path)` — call Contents API, return single `Entry`.
- [x] 2.3 Implement `List(ref, path, recursive)` — non-recursive via Contents API; recursive via Trees API.
- [x] 2.4 Implement `ReadFile(ref, path)` — call Contents API, validate type is file, decode base64 content.
- [x] 2.5 Implement `Download(ref, path, outPath, overwrite)` — file: read + write; directory: recursive list + sequential download.

## 3. Output Formatting
- [x] 3.1 Add `PrintEntries(entries, asJSON)` for `stat` (single) and `ls` (list) output.

## 4. CLI Commands
- [x] 4.1 Implement `ghrepo stat <owner/repo> <path>` with `--ref` and `--json`.
- [x] 4.2 Implement `ghrepo ls <owner/repo> <path>` with `--ref`, `--recursive`, and `--json`.
- [x] 4.3 Implement `ghrepo cat <owner/repo> <path>` with `--ref`.
- [x] 4.4 Implement `ghrepo get <owner/repo> <path> --out <local-path>` with `--ref` and `--overwrite`.
- [x] 4.5 Register all new commands in `root.go`.

## 5. Validation
- [x] 5.1 Unit tests: `owner/repo` parsing, `Entry` output serialization, service-layer type checks.
- [x] 5.2 Integration tests (`httptest`): `stat` success/404, `ls` flat/recursive, `cat` file/directory error, `get` file write/overwrite/directory.
- [x] 5.3 Run `CGO_ENABLED=0 go test ./...`.
- [x] 5.4 Run `openspec validate add-repo-content-commands --strict`.
