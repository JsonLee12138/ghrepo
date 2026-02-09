package cli

import (
	"os"

	"github.com/spf13/cobra"

	clerrors "githubRAGCli/internal/exitcode"
	"githubRAGCli/internal/githubapi"
	"githubRAGCli/internal/output"
)

func newAuthCmd() *cobra.Command {
	auth := &cobra.Command{
		Use:   "auth",
		Short: "Authentication commands",
	}
	auth.AddCommand(newAuthCheckCmd())
	return auth
}

func newAuthCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Verify that the configured token is valid",
		RunE:  runAuthCheck,
	}
}

func runAuthCheck(cmd *cobra.Command, args []string) error {
	cfg := resolveConfig()

	if cfg.Token == "" {
		return clerrors.NewAuthFailure("no token provided: use --token, GITHUB_TOKEN, or GH_TOKEN", nil)
	}

	verboseLog(cfg, "api-base: %s", cfg.APIBase)
	verboseLog(cfg, "timeout: %s", cfg.Timeout)
	// Token is intentionally never logged.

	client := githubapi.NewClient(cfg.APIBase, cfg.Token, cfg.Timeout)
	result, err := client.GetAuthenticatedUser()
	if err != nil {
		return err
	}

	return output.PrintAuth(os.Stdout, output.AuthResult{
		Status:             "ok",
		User:               result.Login,
		RateLimitRemaining: result.RateLimitRemaining,
	}, cfg.JSON)
}
