package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"shamir/pkg/version"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.SetOut(os.Stdout)
	cmd.SetIn(os.Stdin)
	cmd.Use = "shamir"
	cmd.Short = "Command line for Shamir"
	cmd.Long =
		`Command line for Shamir

"Shamir“ be used for (k, n) encrypt. You can use it to encrypt a string or a file.
It will be encrypted as n keys which contains (x, y) and one necessary key.
Any k keys can restore the secret.
For help with any of those, simply call them with --help.`
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.CompletionOptions = cobra.CompletionOptions{DisableDefaultCmd: true}

	cmd.SetVersionTemplate("{{.Version}}\n")
	cmd.Version = version.Version
	return cmd
}
