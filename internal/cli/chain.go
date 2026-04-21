package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewChainCmd returns the cobra command for the `chain` subcommand.
func NewChainCmd() *cobra.Command {
	var redact bool
	var format string

	cmd := &cobra.Command{
		Use:   "chain <file1> [file2 ...]",
		Short: "Merge and display layered .env files in order",
		Long: `chain merges multiple .env files left-to-right.
Entries from later files override entries from earlier files.
Useful for composing base + environment-specific + local overrides.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := envfile.NewChain(args...)
			entries, err := c.Resolve()
			if err != nil {
				return err
			}

			if redact {
				entries = envfile.Redact(entries, "", envfile.IsSensitive)
			}

			switch strings.ToLower(format) {
			case "dotenv", "":
				for _, e := range entries {
					fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", e.Key, e.Value)
				}
			case "keys":
				for _, e := range entries {
					fmt.Fprintln(cmd.OutOrStdout(), e.Key)
				}
			case "sources":
				for i, s := range c.Sources() {
					fmt.Fprintf(cmd.OutOrStdout(), "[%d] %s\n", i+1, s)
				}
			default:
				return fmt.Errorf("unknown format %q (dotenv, keys, sources)", format)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&redact, "redact", false, "mask sensitive values")
	cmd.Flags().StringVar(&format, "format", "dotenv", "output format: dotenv | keys | sources")
	return cmd
}
