## Context
`ghrepo` is planned as a read-only GitHub repository CLI.  
This first slice establishes the minimum operational baseline so later content commands can reuse shared runtime, API, output, and error patterns.

## Goals / Non-Goals
- Goals:
- Define a stable global runtime contract (flags + config precedence).
- Define a reliable `auth check` command contract.
- Define deterministic exit behavior for authentication-related failures.
- Non-Goals:
- Implement `ls`, `cat`, `get`, or `stat`.
- Add parallel download, caching, or release automation.

## Decisions
- Decision: Keep strict layering (`cli` -> `service` -> `githubapi` -> HTTP transport).
  - Rationale: Prevent command handlers from hardcoding request behavior and keep later commands consistent.
- Decision: Token source precedence is fixed as `flag > GITHUB_TOKEN > GH_TOKEN`.
  - Rationale: Matches usage and development docs; deterministic and script-friendly.
- Decision: Use one centralized exit-code mapping path.
  - Rationale: Ensures command behavior is scriptable and avoids per-command divergence.
- Decision: Keep output dual-mode (human text and stable JSON keys) from day one.
  - Rationale: Enables immediate shell and CI integration.

## Risks / Trade-offs
- Risk: Over-scoping the first slice can delay delivery.
  - Mitigation: Limit this change to foundation + auth check only.
- Risk: GitHub `403` can represent different error classes.
  - Mitigation: classify by response context (rate-limit headers vs permission semantics).

## Migration Plan
1. Land this change as the foundation.
2. Build `stat` and `ls` as the next change on top of shared client and error mapping.
3. Add `cat` and `get` after content metadata and type checks are stable.

## Open Questions
- None for this slice.
