package githubapi

import (
	"strings"

	clerrors "githubRAGCli/internal/exitcode"
)

// ParseRepo splits an "owner/repo" string into its components.
// Returns a CLIError with CatBadArgs if the format is invalid.
func ParseRepo(ownerRepo string) (owner, repo string, err error) {
	parts := strings.SplitN(ownerRepo, "/", 3)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", clerrors.NewBadArgs("invalid owner/repo format: "+ownerRepo, nil)
	}
	return parts[0], parts[1], nil
}
