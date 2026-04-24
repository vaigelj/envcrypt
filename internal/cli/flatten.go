package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewFlattenCmd returns the cobra command for the flatten sub-command.
func NewFlattenCmd() *cobra.Command {
	var (
		separator string
		uppercase bool
		prefix    string
		outFile   string
	)

	cmd := &cobra.Command{
		Use:   "flatten <file>",
		Short: "Flatten dot-notation or slash-notation keys into env-style keys",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := envfile.ParseFile(args[0])
			if err != nil {
				return fmt.Errorf("reading %s: %w", args[0], err)
			}

			opts := envfile.FlattenOptions{
				Separator: separator,
				Uppercase: uppercase,
				Prefix:    prefix,
			}

			flat := envfile.Flatten(entries, opts)

			if outFile != "" && outFile != "-" {
				if err := envfile.WriteFile(outFile, flat); err != nil {
					return fmt.Errorf("writing %s: %w", outFile, err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "wrote %d keys to %s\n", len(flat), outFile)
				return nil
			}

			// Default: print to stdout.
			for _, e := range flat {
				if e.Comment != "" {
					fmt.Fprintf(os.Stdout, "# %s\n", strings.TrimPrefix(e.Comment, "# "))
				}
				fmt.Fprintf(os.Stdout, "%s=%s\n", e.Key, e.Value)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&separator, "separator", "s", "_", "separator to use between key segments")
	cmd.Flags().BoolVarP(&uppercase, "uppercase", "u", false, "convert resulting keys to uppercase")
	cmd.Flags().StringVarP(&prefix, "prefix", "p", "", "prefix to prepend to every key")
	cmd.Flags().StringVarP(&outFile, "output", "o", "-", "output file (default: stdout)")

	return cmd
}

func init() {
	rootCmd.AddCommand(NewFlattenCmd())
}
