package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewTrimCmd returns the cobra command for the trim sub-command.
func NewTrimCmd() *cobra.Command {
	var (
		keys    []string
		exclude []string
		cutset  string
		inPlace bool
		output  string
	)

	cmd := &cobra.Command{
		Use:   "trim <file>",
		Short: "Trim leading/trailing whitespace from .env values",
		Long: `Trim removes leading and trailing whitespace (or a custom cutset)
from the values of all entries in a .env file.

Use --keys to limit trimming to specific keys, or --exclude to skip them.
Use --cutset to trim a custom set of characters instead of whitespace.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			entries, err := envfile.ParseFile(path)
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			var opts []envfile.TrimOption
			if len(keys) > 0 {
				opts = append(opts, envfile.WithTrimKeys(keys...))
			}
			if len(exclude) > 0 {
				opts = append(opts, envfile.WithTrimExclude(exclude...))
			}
			if cutset != "" {
				opts = append(opts, envfile.WithTrimCutset(cutset))
			}

			trimmed := envfile.Trim(entries, opts...)

			switch {
			case inPlace:
				if err := envfile.WriteFile(path, trimmed); err != nil {
					return fmt.Errorf("write: %w", err)
				}
			case output != "":
				if err := envfile.WriteFile(output, trimmed); err != nil {
					return fmt.Errorf("write: %w", err)
				}
			default:
				for _, e := range trimmed {
					if e.Comment != "" {
						fmt.Fprintf(os.Stdout, "# %s\n", e.Comment)
					}
					fmt.Fprintf(os.Stdout, "%s=%s\n", e.Key, e.Value)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&keys, "keys", "k", nil, "only trim these keys (comma-separated)")
	cmd.Flags().StringSliceVarP(&exclude, "exclude", "e", nil, "skip these keys (comma-separated)")
	cmd.Flags().StringVar(&cutset, "cutset", "", "characters to trim instead of whitespace")
	cmd.Flags().BoolVarP(&inPlace, "in-place", "i", false, "overwrite the source file")
	cmd.Flags().StringVarP(&output, "output", "o", "", "write result to this file")

	return cmd
}
