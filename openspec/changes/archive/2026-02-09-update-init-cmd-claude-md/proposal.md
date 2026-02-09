# Change: Extend init Command to Configure CLAUDE.md

## Why
The `init` command currently only configures `AGENTS.md` with ghrepo repository configuration. Many projects also use `CLAUDE.md` for AI assistant instructions. By extending `init` to also create/update `CLAUDE.md`, we provide a complete initialization experience for projects that need both files.

## What Changes
- Extend `init` command to create or append to `CLAUDE.md` file
- `CLAUDE.md` will contain OpenSpec instructions and ghrepo configuration context
- Maintain backward compatibility with existing `init` behavior (still creates/updates `AGENTS.md`)
- Both files should be created/updated in a single `init` invocation

## Impact
- Affected specs: `init-command` capability
- Affected code: `internal/cli/cmd_init.go`
- No breaking changes; purely additive enhancement
