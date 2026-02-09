## MODIFIED Requirements

### Requirement: Initialize Repository Configuration
The system SHALL create or update initialization files to configure AI assistant behaviors and repository contexts.

#### Scenario: Create both AGENTS.md and CLAUDE.md
- **WHEN** `init <owner/repo>` is executed in a directory with no configuration files
- **THEN** both `AGENTS.md` (with ghrepo commands) and `CLAUDE.md` (with OpenSpec instructions) are created

#### Scenario: Update existing files
- **WHEN** `init <owner/repo>` is executed in a directory with existing AGENTS.md or CLAUDE.md
- **THEN** the command appends ghrepo/OpenSpec configuration if not already present, or skips if already configured

#### Scenario: Idempotent execution
- **WHEN** `init <owner/repo>` is executed multiple times on the same directory
- **THEN** the command succeeds without duplicating content
