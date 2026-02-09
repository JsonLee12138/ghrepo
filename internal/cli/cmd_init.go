package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	clerrors "githubRAGCli/internal/exitcode"
	"githubRAGCli/internal/githubapi"
)

const agentsFileName = "AGENTS.md"
const claudeFileName = "CLAUDE.md"

const (
	ghrepoBlockStart = "<!-- GHREPO:START -->"
	ghrepoBlockEnd   = "<!-- GHREPO:END -->"
	openspecBlockStart = "<!-- OPENSPEC:START -->"
	openspecBlockEnd   = "<!-- OPENSPEC:END -->"
)

const agentsTemplate = `<!-- GHREPO:START -->
# ghrepo

This project is configured to work with the GitHub repository: %s/%s

## Commands

` + "```" + `bash
# List directory contents
ghrepo ls %[1]s/%[2]s .

# Read file content
ghrepo cat %[1]s/%[2]s <path>

# Show file metadata
ghrepo stat %[1]s/%[2]s <path>

# Download files
ghrepo get %[1]s/%[2]s <path> --out <local-path>

# Create or update a file
ghrepo put %[1]s/%[2]s <path> -m "message" --file <local-file>

# Delete a file
ghrepo rm %[1]s/%[2]s <path> -m "message"
` + "```" + `
<!-- GHREPO:END -->`

const claudeTemplate = `<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open ` + "`@/openspec/AGENTS.md`" + ` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use ` + "`@/openspec/AGENTS.md`" + ` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->`

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init <owner/repo>",
		Short: "Initialize AGENTS.md and CLAUDE.md with repository and project configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			owner, repo, err := githubapi.ParseRepo(args[0])
			if err != nil {
				return err
			}

			cwd, err := os.Getwd()
			if err != nil {
				return clerrors.NewLocalWriteErr("failed to get current directory", err)
			}

			// Handle AGENTS.md
			agentsPath := filepath.Join(cwd, agentsFileName)
			agentsBlock := fmt.Sprintf(agentsTemplate, owner, repo)

			existing, err := os.ReadFile(agentsPath)
			if err != nil {
				// File does not exist — create with block.
				if err := os.WriteFile(agentsPath, []byte(agentsBlock+"\n"), 0o644); err != nil {
					return clerrors.NewLocalWriteErr("failed to write "+agentsFileName, err)
				}
				fmt.Fprintf(os.Stderr, "created %s for %s/%s\n", agentsFileName, owner, repo)
			} else {
				content := string(existing)
				// Already has ghrepo block — skip.
				if strings.Contains(content, ghrepoBlockStart) {
					fmt.Fprintf(os.Stderr, "%s already contains ghrepo configuration, skipping\n", agentsFileName)
				} else {
					// Append block to existing file.
					if !strings.HasSuffix(content, "\n") {
						content += "\n"
					}
					content += "\n" + agentsBlock + "\n"

					if err := os.WriteFile(agentsPath, []byte(content), 0o644); err != nil {
						return clerrors.NewLocalWriteErr("failed to update "+agentsFileName, err)
					}
					fmt.Fprintf(os.Stderr, "appended ghrepo configuration to %s for %s/%s\n", agentsFileName, owner, repo)
				}
			}

			// Handle CLAUDE.md
			claudePath := filepath.Join(cwd, claudeFileName)

			existing, err = os.ReadFile(claudePath)
			if err != nil {
				// File does not exist — create with block.
				if err := os.WriteFile(claudePath, []byte(claudeTemplate+"\n"), 0o644); err != nil {
					return clerrors.NewLocalWriteErr("failed to write "+claudeFileName, err)
				}
				fmt.Fprintf(os.Stderr, "created %s for AI assistant instructions\n", claudeFileName)
			} else {
				content := string(existing)
				// Already has openspec block — skip.
				if strings.Contains(content, openspecBlockStart) {
					fmt.Fprintf(os.Stderr, "%s already contains OpenSpec instructions, skipping\n", claudeFileName)
				} else {
					// Append block to existing file.
					if !strings.HasSuffix(content, "\n") {
						content += "\n"
					}
					content += "\n" + claudeTemplate + "\n"

					if err := os.WriteFile(claudePath, []byte(content), 0o644); err != nil {
						return clerrors.NewLocalWriteErr("failed to update "+claudeFileName, err)
					}
					fmt.Fprintf(os.Stderr, "appended OpenSpec instructions to %s\n", claudeFileName)
				}
			}

			return nil
		},
	}

	return cmd
}
