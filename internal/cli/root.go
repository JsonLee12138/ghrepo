package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"githubRAGCli/internal/config"
)

var (
	flagToken   string
	flagAPIBase string
	flagTimeout time.Duration
	flagJSON    bool
	flagVerbose bool
)

// NewRootCmd creates the top-level ghrepo command.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "ghrepo",
		Short:         "Read-only CLI for GitHub repository contents",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.PersistentFlags().StringVar(&flagToken, "token", "", "GitHub personal access token (overrides GITHUB_TOKEN / GH_TOKEN)")
	root.PersistentFlags().StringVar(&flagAPIBase, "api-base", config.DefaultAPIBase, "GitHub API base URL")
	root.PersistentFlags().DurationVar(&flagTimeout, "timeout", config.DefaultTimeout, "HTTP request timeout")
	root.PersistentFlags().BoolVar(&flagJSON, "json", false, "Output in JSON format")
	root.PersistentFlags().BoolVar(&flagVerbose, "verbose", false, "Enable verbose output (never prints token)")

	root.AddCommand(newAuthCmd())
	root.AddCommand(newStatCmd())
	root.AddCommand(newLsCmd())
	root.AddCommand(newCatCmd())
	root.AddCommand(newGetCmd())

	return root
}

// resolveConfig builds a Config from flags and environment.
func resolveConfig() config.Config {
	return config.Config{
		Token:   config.ResolveToken(flagToken),
		APIBase: flagAPIBase,
		Timeout: flagTimeout,
		JSON:    flagJSON,
		Verbose: flagVerbose,
	}
}

// JSONFlag returns the current value of --json. Used by main for error output.
func JSONFlag() bool {
	return flagJSON
}

// verboseLog prints a message only when --verbose is set.
// It must never include token values.
func verboseLog(cfg config.Config, format string, args ...any) {
	if cfg.Verbose {
		fmt.Printf("[verbose] "+format+"\n", args...)
	}
}
