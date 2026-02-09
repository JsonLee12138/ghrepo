## 1. Core Infrastructure

- [x] 1.1 Add `ExitUserAbort = 17` and `CatUserAbort` to `internal/exitcode/codes.go`
- [x] 1.2 Add `NewUserAbort` constructor to `internal/exitcode/codes.go`
- [x] 1.3 Add `doPut` and `doDelete` methods to `internal/githubapi/client.go` (authenticated PUT/DELETE with JSON body)
- [x] 1.4 Add `CreateOrUpdateFile` and `DeleteFile` request/response types to `internal/githubapi/client.go`

## 2. Service Layer

- [x] 2.1 Add `CreateOrUpdateFile(branch, path, message string, content []byte) (*MutationResult, error)` to `internal/service/repo_service.go`
- [x] 2.2 Add `DeleteFile(branch, path, message string) (*MutationResult, error)` to `internal/service/repo_service.go`
- [x] 2.3 Add `MutationResult` struct (Action, Path, SHA, Branch) to service package

## 3. Output Layer

- [x] 3.1 Add `MutationResultData` and `PrintMutationResult` to `internal/output/print.go`

## 4. CLI Layer

- [x] 4.1 Add `confirmPrompt` helper to `internal/cli/helpers.go` (stderr prompt, stdin read, --yes bypass, non-TTY detection)
- [x] 4.2 Create `internal/cli/cmd_put.go` with `put` command
- [x] 4.3 Create `internal/cli/cmd_rm.go` with `rm` command
- [x] 4.4 Register `put` and `rm` in `internal/cli/root.go`, update root Short description

## 5. Tests

- [x] 5.1 Add tests for exit code user abort
- [x] 5.2 Add tests for `CreateOrUpdateFile` and `DeleteFile` service methods (using httptest)
- [x] 5.3 Add tests for `PrintMutationResult` output formatting
