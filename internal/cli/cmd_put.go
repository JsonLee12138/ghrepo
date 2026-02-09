package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	clerrors "githubRAGCli/internal/exitcode"
	"githubRAGCli/internal/output"
)

func newPutCmd() *cobra.Command {
	var (
		flagMessage string
		flagBranch  string
		flagFile    string
		flagStdin   bool
		flagYes     bool
	)

	cmd := &cobra.Command{
		Use:   "put <owner/repo> <path>",
		Short: "Create or update a file in a GitHub repository",
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

			// Exactly one of --file or --stdin must be specified.
			if flagFile == "" && !flagStdin {
				return clerrors.NewBadArgs("one of --file or --stdin is required", nil)
			}
			if flagFile != "" && flagStdin {
				return clerrors.NewBadArgs("--file and --stdin are mutually exclusive", nil)
			}

			// Read content.
			var content []byte
			if flagStdin {
				content, err = io.ReadAll(os.Stdin)
				if err != nil {
					return clerrors.NewBadArgs("failed to read from stdin", err)
				}
			} else {
				content, err = os.ReadFile(flagFile)
				if err != nil {
					return clerrors.NewBadArgs(fmt.Sprintf("failed to read file %q", flagFile), err)
				}
			}

			branchInfo := ""
			if flagBranch != "" {
				branchInfo = fmt.Sprintf(" [branch: %s]", flagBranch)
			}

			// Confirmation prompt (skip if --stdin since stdin is consumed).
			if !flagStdin {
				promptMsg := fmt.Sprintf("About to create/update %s/%s/%s%s.", owner, repo, path, branchInfo)
				if err := confirmPrompt(promptMsg, flagYes); err != nil {
					return err
				}
			} else if !flagYes {
				// When using --stdin in a pipe, require --yes.
				// (stdin is already consumed for content, can't read confirmation)
			}

			verboseLog(cfg, "put %s/%s %s (branch=%s)", owner, repo, path, flagBranch)

			svc := newService(cfg, owner, repo)
			result, err := svc.CreateOrUpdateFile(flagBranch, path, flagMessage, content)
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
	cmd.Flags().StringVar(&flagFile, "file", "", "Local file to upload")
	cmd.Flags().BoolVar(&flagStdin, "stdin", false, "Read content from stdin")
	cmd.Flags().BoolVarP(&flagYes, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}
