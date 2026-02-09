## 1. Implementation

### 1.1 Update cmd_init.go
- [x] Add `CLAUDE.md` template constant
- [x] Add logic to create/append to `CLAUDE.md` with OpenSpec instructions
- [x] Update `init` command short description to reflect dual-file behavior
- [x] Ensure proper handling of both AGENTS.md and CLAUDE.md in the same invocation

### 1.2 Testing
- [x] Test init with no existing files (both AGENTS.md and CLAUDE.md should be created)
- [x] Test init with existing AGENTS.md (should append/skip appropriately)
- [x] Test init with existing CLAUDE.md (should append/skip appropriately)
- [x] Test init with both files already present

### 1.3 Verification
- [x] Verify CLAUDE.md content includes OpenSpec instructions
- [x] Verify CLAUDE.md content includes ghrepo configuration context
- [x] Verify backward compatibility (existing workflows still work)
