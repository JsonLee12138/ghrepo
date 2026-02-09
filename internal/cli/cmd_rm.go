package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	clerrors "githubRAGCli/internal/exitcode"
	"githubRAGCli/internal/output"
)

func newRmCmd() *cobra.Command {
	var (
		flagMessage string
		flagBranch  string
		flagYes     bool
	)

	cmd := &cobra.Command{
		Use:   "rm <owner/repo> <path>",
		Short: "Delete a file from a GitHub repository",
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

			if flagMessage == "" {
				return clerrors.NewBadArgs("--message / -m is required", nil)
			}

			branchInfo := ""
			if flagBranch != "" {
				branchInfo = fmt.Sprintf(" [branch: %s]", flagBranch)
			}

			promptMsg := fmt.Sprintf("About to delete %s/%s/%s%s.", owner, repo, path, branchInfo)
			if err := confirmPrompt(promptMsg, flagYes); err != nil {
				return err
			}

			verboseLog(cfg, "rm %s/%s %s (branch=%s)", owner, repo, path, flagBranch)

			svc := newService(cfg, owner, repo)
			result, err := svc.DeleteFile(flagBranch, path, flagMessage)
			if err != nil {
				return err
			}

			return output.PrintMutationResult(os.Stdout, output.MutationResultData{
				Action: result.Action,
				Path:   result.Path,
				SHA:    result.SHA,
				Branch: result.Branch,
			}, cfg.JSON)
		},
	}

	cmd.Flags().StringVarP(&flagMessage, "message", "m", "", "Commit message (required)")
	cmd.Flags().StringVarP(&flagBranch, "branch", "b", "", "Target branch (defaults to repo default branch)")
	cmd.Flags().BoolVarP(&flagYes, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}
