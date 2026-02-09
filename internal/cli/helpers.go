package cli

import (
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
