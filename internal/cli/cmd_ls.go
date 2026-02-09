package cli

import (
	"os"

	"github.com/spf13/cobra"

	"githubRAGCli/internal/output"
)

func newLsCmd() *cobra.Command {
	var (
		flagRef       string
		flagRecursive bool
	)

	cmd := &cobra.Command{
		Use:   "ls <owner/repo> <path>",
		Short: "List directory contents",
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

			verboseLog(cfg, "ls %s/%s %s (ref=%s, recursive=%v)", owner, repo, path, flagRef, flagRecursive)

			svc := newService(cfg, owner, repo)
			entries, err := svc.List(flagRef, path, flagRecursive)
			if err != nil {
				return err
			}

			outEntries := make([]output.EntryData, 0, len(entries))
			for i := range entries {
				outEntries = append(outEntries, serviceEntryToOutput(&entries[i]))
			}

			return output.PrintEntries(os.Stdout, outEntries, cfg.JSON)
		},
	}

	cmd.Flags().StringVar(&flagRef, "ref", "", "Git ref (branch, tag, or SHA)")
	cmd.Flags().BoolVar(&flagRecursive, "recursive", false, "List recursively using the Trees API")

	return cmd
}
