package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewNormalizeCmd returns the cobra command for normalizing .env files.
func NewNormalizeCmd() *cobra.Command {
	var (
		upperKeys   bool
		trimValues  bool
		quoteValues bool
		removeEmpty bool
		output      string
	)

	cmd := &cobra.Command{
		Use:   "normalize <file>",
		Short: "Normalize a .env file (uppercase keys, trim values, etc.)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			entries, err := envfile.ParseFile(path)
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			var opts []envfile.NormalizeOption
			if upperKeys {
				opts = append(opts, envfile.WithUpperKeys())
			}
			if trimValues {
				opts = append(opts, envfile.WithTrimValues())
			}
			if quoteValues {
				opts = append(opts, envfile.WithQuoteValues())
			}
			if removeEmpty {
				opts = append(opts, envfile.WithRemoveEmpty())
			}

			normalized := envfile.Normalize(entries, opts...)

			dest := path
			if output != "" {
				dest = output
			}

			if err := envfile.WriteFile(dest, normalized); err != nil {
				return fmt.Errorf("write: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "normalized %d entries → %s\n", len(normalized), dest)
			return nil
		},
	}

	cmd.Flags().BoolVar(&upperKeys, "upper-keys", false, "Convert all keys to uppercase")
	cmd.Flags().BoolVar(&trimValues, "trim-values", true, "Trim whitespace from values")
	cmd.Flags().BoolVar(&quoteValues, "quote-values", false, "Quote values that contain spaces")
	cmd.Flags().BoolVar(&removeEmpty, "remove-empty", false, "Remove entries with empty values")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Write result to a different file (default: overwrite input)")

	return cmd
}
