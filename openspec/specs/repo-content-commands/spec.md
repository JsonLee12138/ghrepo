# repo-content-commands Specification

## Purpose
TBD - created by archiving change add-repo-content-commands. Update Purpose after archive.
## Requirements
### Requirement: Repository Argument Parsing
The CLI SHALL parse `<owner/repo>` as a positional argument for all content commands and fail with exit code `13` when the format is invalid.

#### Scenario: Valid owner/repo is accepted
- **WHEN** the user provides `owner/repo` as the first positional argument
- **THEN** the CLI extracts `owner` and `repo` correctly

#### Scenario: Invalid format fails with exit code 13
- **WHEN** the user provides a string without exactly one `/` separator
- **THEN** the command fails with exit code `13`

### Requirement: Stat Command
The system SHALL provide `ghrepo stat <owner/repo> <path>` to query path metadata from the GitHub Contents API and return type, path, SHA, size, and download URL.

#### Scenario: Stat a file returns metadata
- **WHEN** the user runs `ghrepo stat owner/repo README.md`
- **AND** the path exists and is a file
- **THEN** the command exits with code `0`
- **AND** output includes `type`, `path`, `sha`, `size`, and `download_url`

#### Scenario: Stat a directory returns metadata
- **WHEN** the user runs `ghrepo stat owner/repo docs`
- **AND** the path exists and is a directory
- **THEN** the command exits with code `0`
- **AND** output includes `type: dir`, `path`, and `sha`

#### Scenario: Stat a non-existent path fails
- **WHEN** the path does not exist in the repository
- **THEN** the command exits with code `12`

#### Scenario: JSON output uses stable keys
- **WHEN** the user passes `--json`
- **THEN** the command emits valid JSON with stable fields `type`, `path`, `sha`, `size`, `download_url`

### Requirement: List Command
The system SHALL provide `ghrepo ls <owner/repo> <path>` to list directory contents, with optional `--recursive` and `--ref` flags.

#### Scenario: Non-recursive listing returns immediate children
- **WHEN** the user runs `ghrepo ls owner/repo docs`
- **AND** `docs` is a directory
- **THEN** the command exits with code `0`
- **AND** output lists each immediate child entry

#### Scenario: Recursive listing returns full subtree
- **WHEN** the user runs `ghrepo ls owner/repo docs --recursive`
- **THEN** the command uses the Trees API to return the complete subtree

#### Scenario: Listing a file path fails
- **WHEN** the target path is a file, not a directory
- **THEN** the command fails with exit code `13`

#### Scenario: Ref flag selects branch or tag
- **WHEN** the user passes `--ref main`
- **THEN** the API request includes the specified ref

### Requirement: Cat Command
The system SHALL provide `ghrepo cat <owner/repo> <path>` to read file content and write it to stdout.

#### Scenario: Cat a file outputs content to stdout
- **WHEN** the user runs `ghrepo cat owner/repo README.md`
- **AND** the path is a file
- **THEN** the decoded file content is written to stdout
- **AND** the command exits with code `0`

#### Scenario: Cat a directory fails
- **WHEN** the target path is a directory
- **THEN** the command fails with exit code `13`

#### Scenario: Ref flag selects branch or tag
- **WHEN** the user passes `--ref v1.0`
- **THEN** the API request includes the specified ref

### Requirement: Get Command
The system SHALL provide `ghrepo get <owner/repo> <path> --out <local-path>` to download files or directories to the local filesystem.

#### Scenario: Download a single file
- **WHEN** the user runs `ghrepo get owner/repo README.md --out ./README.md`
- **AND** the path is a file
- **THEN** the file is written to the specified local path
- **AND** the command exits with code `0`

#### Scenario: Download a directory preserves structure
- **WHEN** the user runs `ghrepo get owner/repo docs --out ./local-docs`
- **AND** the path is a directory
- **THEN** all files in the directory are downloaded sequentially
- **AND** the repository-relative directory structure is preserved under `--out`

#### Scenario: Overwrite is blocked by default
- **WHEN** the local file already exists
- **AND** `--overwrite` is not set
- **THEN** the command fails with exit code `16`

#### Scenario: Overwrite flag allows replacement
- **WHEN** `--overwrite` is set
- **THEN** existing local files are replaced

#### Scenario: Local write failure maps to exit code 16
- **WHEN** writing to the local filesystem fails
- **THEN** the command exits with code `16`

