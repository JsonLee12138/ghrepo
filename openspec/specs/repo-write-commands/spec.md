# repo-write-commands Specification

## Purpose
TBD - created by archiving change add-write-delete-commands. Update Purpose after archive.
## Requirements
### Requirement: Put Command
The CLI SHALL provide a `put` command that creates or updates a single file in a GitHub repository.

The command signature SHALL be: `ghrepo put <owner/repo> <path> --message <msg> (--file <local-path> | --stdin) [--branch <branch>] [--yes]`

The command SHALL read file content from `--file` or `--stdin` (exactly one required), base64-encode it, and send a `PUT /repos/{owner}/{repo}/contents/{path}` request.

The command SHALL auto-detect whether the file already exists (create vs update) by first calling `GET /repos/{owner}/{repo}/contents/{path}`. If the file exists, the current SHA SHALL be included in the PUT request body.

The `--message` / `-m` flag SHALL be required.

The `--branch` / `-b` flag SHALL be optional (defaults to repo default branch).

#### Scenario: Create a new file with confirmation
- **WHEN** user runs `ghrepo put owner/repo path/to/file.txt -m "add file" --file ./local.txt`
- **AND** the file does not exist in the repository
- **AND** user confirms the prompt with "y"
- **THEN** the file SHALL be created in the repository
- **AND** the commit SHA and file path SHALL be printed to stdout

#### Scenario: Update an existing file with confirmation
- **WHEN** user runs `ghrepo put owner/repo existing.txt -m "update" --file ./new.txt`
- **AND** the file already exists in the repository
- **AND** user confirms the prompt with "y"
- **THEN** the file SHALL be updated with the new content
- **AND** the commit SHA and file path SHALL be printed to stdout

#### Scenario: Read content from stdin
- **WHEN** user runs `echo "hello" | ghrepo put owner/repo file.txt -m "msg" --stdin --yes`
- **THEN** the content from stdin SHALL be used as the file content

#### Scenario: User cancels put operation
- **WHEN** user runs `ghrepo put owner/repo file.txt -m "msg" --file ./f.txt`
- **AND** user responds "n" to the confirmation prompt
- **THEN** no API write call SHALL be made
- **AND** the CLI SHALL exit with code 17 (user abort)

#### Scenario: Non-interactive without --yes
- **WHEN** stdin is not a terminal (piped input)
- **AND** `--yes` flag is not provided
- **AND** `--stdin` is not used
- **THEN** the CLI SHALL abort with exit code 17

### Requirement: Rm Command
The CLI SHALL provide an `rm` command that deletes a single file from a GitHub repository.

The command signature SHALL be: `ghrepo rm <owner/repo> <path> --message <msg> [--branch <branch>] [--yes]`

The command SHALL first retrieve the file's current SHA via `GET /repos/{owner}/{repo}/contents/{path}`, then send a `DELETE /repos/{owner}/{repo}/contents/{path}` request with the SHA.

The `--message` / `-m` flag SHALL be required.

#### Scenario: Delete a file with confirmation
- **WHEN** user runs `ghrepo rm owner/repo path/to/file.txt -m "remove file"`
- **AND** the file exists in the repository
- **AND** user confirms the prompt with "y"
- **THEN** the file SHALL be deleted from the repository
- **AND** the commit SHA SHALL be printed to stdout

#### Scenario: User cancels delete operation
- **WHEN** user runs `ghrepo rm owner/repo file.txt -m "remove"`
- **AND** user responds "n" to the confirmation prompt
- **THEN** no API delete call SHALL be made
- **AND** the CLI SHALL exit with code 17 (user abort)

#### Scenario: Delete non-existent file
- **WHEN** user runs `ghrepo rm owner/repo nonexistent.txt -m "remove" --yes`
- **AND** the file does not exist
- **THEN** the CLI SHALL exit with code 12 (not found)

### Requirement: Confirmation Prompt
All write and delete operations SHALL display an interactive confirmation prompt before executing.

The prompt SHALL be printed to stderr in the format: `⚠ About to <action> <owner/repo>/<path> [branch: <branch>]. Continue? [y/N] `

The default answer SHALL be No (safe default).

The `--yes` / `-y` flag SHALL skip the confirmation prompt.

#### Scenario: Confirmation shown for put
- **WHEN** user runs `ghrepo put owner/repo f.txt -m "msg" --file ./f.txt`
- **THEN** a confirmation prompt SHALL be displayed on stderr before the API call

#### Scenario: Confirmation shown for rm
- **WHEN** user runs `ghrepo rm owner/repo f.txt -m "msg"`
- **THEN** a confirmation prompt SHALL be displayed on stderr before the API call

#### Scenario: Skip confirmation with --yes
- **WHEN** user provides the `--yes` or `-y` flag
- **THEN** no confirmation prompt SHALL be displayed
- **AND** the operation SHALL proceed immediately

### Requirement: User Abort Exit Code
The CLI SHALL define exit code 17 for user-aborted operations.

When the user declines a confirmation prompt or when a non-interactive session lacks the `--yes` flag, the CLI SHALL exit with code 17.

#### Scenario: Exit code on user abort
- **WHEN** user declines the confirmation prompt
- **THEN** the process SHALL exit with code 17

### Requirement: Mutation Result Output
Write and delete operations SHALL output results in both text and JSON formats, controlled by the `--json` flag.

The result SHALL include: action (created/updated/deleted), path, commit SHA, and branch.

#### Scenario: Text output for put
- **WHEN** `--json` is not set
- **THEN** output SHALL be human-readable text lines (action, path, sha, branch)

#### Scenario: JSON output for rm
- **WHEN** `--json` is set
- **THEN** output SHALL be a JSON object with action, path, sha, branch fields

