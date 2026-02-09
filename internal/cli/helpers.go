package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"githubRAGCli/internal/config"
	clerrors "githubRAGCli/internal/exitcode"
	"githubRAGCli/internal/githubapi"
	"githubRAGCli/internal/output"
	"githubRAGCli/internal/service"
)

// requireToken validates that a token is available and returns an error if not.
func requireToken(cfg config.Config) error {
	if cfg.Token == "" {
		return clerrors.NewAuthFailure("no token provided: use --token, GITHUB_TOKEN, or GH_TOKEN", nil)
	}
	return nil
}

// parseRepoArg parses and validates the owner/repo positional argument.
func parseRepoArg(args []string) (owner, repo string, err error) {
	if len(args) < 1 {
		return "", "", clerrors.NewBadArgs("missing required argument: <owner/repo>", nil)
	}
	return githubapi.ParseRepo(args[0])
}

// newService creates a RepoService from resolved config and parsed owner/repo.
func newService(cfg config.Config, owner, repo string) *service.RepoService {
	return service.NewRepoService(cfg.APIBase, cfg.Token, cfg.Timeout, owner, repo)
}

// serviceEntryToOutput converts a service.Entry to an output.EntryData.
func serviceEntryToOutput(e *service.Entry) output.EntryData {
	return output.EntryData{
		Type:        e.Type,
		Path:        e.Path,
		SHA:         e.SHA,
		Size:        e.Size,
		DownloadURL: e.DownloadURL,
	}
}

// confirmPrompt displays a confirmation prompt on stderr and reads user input.
// If skipConfirm is true, the prompt is skipped and the operation proceeds.
// If stdin is not a terminal and skipConfirm is false, it returns a user abort error.
func confirmPrompt(message string, skipConfirm bool) error {
	if skipConfirm {
		return nil
	}

	// Check if stdin is a terminal.
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return clerrors.NewUserAbort("non-interactive session requires --yes flag", nil)
	}

	fmt.Fprintf(os.Stderr, "⚠ %s Continue? [y/N] ", message)

	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer != "y" && answer != "yes" {
		return clerrors.NewUserAbort("operation cancelled by user", nil)
	}
	return nil
}
