package cli

import (
	"os"

	"github.com/spf13/cobra"
)

func newCatCmd() *cobra.Command {
	var flagRef string

	cmd := &cobra.Command{
		Use:   "cat <owner/repo> <path>",
		Short: "Output file content to stdout",
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

			verboseLog(cfg, "cat %s/%s %s (ref=%s)", owner, repo, path, flagRef)

			svc := newService(cfg, owner, repo)
			data, err := svc.ReadFile(flagRef, path)
			if err != nil {
				return err
			}

			_, err = os.Stdout.Write(data)
			return err
		},
	}

	cmd.Flags().StringVar(&flagRef, "ref", "", "Git ref (branch, tag, or SHA)")

	return cmd
}
