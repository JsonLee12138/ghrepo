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

const (
	ghrepoBlockStart = "<!-- GHREPO:START -->"
	ghrepoBlockEnd   = "<!-- GHREPO:END -->"
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

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init <owner/repo>",
		Short: "Initialize AGENTS.md with repository configuration",
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

			outPath := filepath.Join(cwd, agentsFileName)
			block := fmt.Sprintf(agentsTemplate, owner, repo)

			existing, err := os.ReadFile(outPath)
			if err != nil {
				// File does not exist — create with block.
				if err := os.WriteFile(outPath, []byte(block+"\n"), 0o644); err != nil {
					return clerrors.NewLocalWriteErr("failed to write "+agentsFileName, err)
				}
				fmt.Fprintf(os.Stderr, "created %s for %s/%s\n", agentsFileName, owner, repo)
				return nil
			}

			content := string(existing)

			// Already has ghrepo block — skip.
			if strings.Contains(content, ghrepoBlockStart) {
				fmt.Fprintf(os.Stderr, "%s already contains ghrepo configuration, skipping\n", agentsFileName)
				return nil
			}

			// Append block to existing file.
			if !strings.HasSuffix(content, "\n") {
				content += "\n"
			}
			content += "\n" + block + "\n"

			if err := os.WriteFile(outPath, []byte(content), 0o644); err != nil {
				return clerrors.NewLocalWriteErr("failed to update "+agentsFileName, err)
			}
			fmt.Fprintf(os.Stderr, "appended ghrepo configuration to %s for %s/%s\n", agentsFileName, owner, repo)
			return nil
		},
	}

	return cmd
}
