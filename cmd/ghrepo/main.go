package main

import (
	"fmt"
	"os"

	"githubRAGCli/internal/cli"
	clerrors "githubRAGCli/internal/exitcode"
	"githubRAGCli/internal/output"
)

func main() {
	root := cli.NewRootCmd()
	if err := root.Execute(); err != nil {
		// Print the error via the unified output path.
		asJSON := cli.JSONFlag()
		output.PrintError(os.Stderr, err.Error(), asJSON)

		// Map CLIError to its exit code; default to 1 for unknown errors.
		if ce, ok := err.(*clerrors.CLIError); ok {
			os.Exit(ce.ExitCode())
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
