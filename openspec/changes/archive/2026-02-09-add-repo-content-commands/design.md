## Context
This change builds on the CLI foundation (`add-cli-foundation-auth-check`) to deliver all four read-only content commands defined in `docs/USAGE.md`.
The existing `githubapi.Client`, `exitcode`, `config`, and `output` packages are reused and extended.

## Goals / Non-Goals
- Goals:
  - Implement `stat`, `ls`, `cat`, `get` per USAGE.md contract.
  - Introduce a service layer that isolates business logic (recursive listing, type checking, download orchestration) from the CLI and API layers.
  - Reuse the existing error classification and exit-code mapping for new HTTP failure modes (`404` → `12`, `13` for argument errors, `16` for local write failures).
- Non-Goals:
  - Concurrent downloads (deferred to v0.2).
  - Caching or retry logic beyond what the foundation already provides.
  - `include/exclude` filtering (v0.2).

## Decisions
- Decision: Add `owner/repo` parsing as a shared utility, not per-command.
  - Rationale: All four commands need it; single validation point avoids divergence.
- Decision: Service layer sits between CLI and API (`cli` → `service` → `githubapi`).
  - Rationale: Matches DESIGN.md layering. Commands stay thin; business rules (e.g. "ls on a file is an error") live in the service.
- Decision: `get` writes files sequentially in v0.1.
  - Rationale: Simplicity first; concurrent download is a v0.2 feature.
- Decision: `cat` decodes base64 content from the Contents API response.
  - Rationale: Avoids a second HTTP call to `download_url`; works for files up to 1 MB (GitHub API limit).

## Risks / Trade-offs
- Risk: Contents API returns base64 content only for files ≤ 1 MB.
  - Mitigation: For v0.1 this is acceptable; document the limit. v0.2 can fall back to Blob API for larger files.
- Risk: Recursive tree listing via Trees API may hit size limits for very large repos.
  - Mitigation: Trees API returns `truncated: true`; surface a warning and still return partial results.

## Open Questions
- None for this slice.
