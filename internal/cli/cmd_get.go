package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	clerrors "githubRAGCli/internal/exitcode"
)

func newGetCmd() *cobra.Command {
	var (
		flagRef       string
		flagOut       string
		flagOverwrite bool
	)

	cmd := &cobra.Command{
		Use:   "get <owner/repo> <path>",
		Short: "Download a file or directory to the local filesystem",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := resolveConfig()

			if err := requireToken(cfg); err != nil {
				return err
			}

			owner, repo, err := parseRepoArg(args)
			if err != nil {
				return err
			}
			path := args[1]

			if flagOut == "" {
				return clerrors.NewBadArgs("--out is required", nil)
			}

			verboseLog(cfg, "get %s/%s %s -> %s (ref=%s, overwrite=%v)", owner, repo, path, flagOut, flagRef, flagOverwrite)

			svc := newService(cfg, owner, repo)
			if err := svc.Download(flagRef, path, flagOut, flagOverwrite); err != nil {
				return err
			}

			fmt.Fprintf(os.Stderr, "downloaded to %s\n", flagOut)
			return nil
		},
	}

	cmd.Flags().StringVar(&flagRef, "ref", "", "Git ref (branch, tag, or SHA)")
	cmd.Flags().StringVar(&flagOut, "out", "", "Local output path (required)")
	cmd.Flags().BoolVar(&flagOverwrite, "overwrite", false, "Overwrite existing files")

	return cmd
}
