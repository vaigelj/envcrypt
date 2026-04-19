package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewCopyCmd returns the `envcrypt copy` subcommand.
func NewCopyCmd() *cobra.Command {
	var overwrite bool
	var exclude []string

	cmd := &cobra.Command{
		Use:   "copy <src> <dst>",
		Short: "Copy env vars from src file into dst file",
		Long: `Reads key/value pairs from <src> and writes them into <dst>.
Existing keys in <dst> are preserved unless --overwrite is set.
Use --exclude to skip specific keys.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			srcPath := args[0]
			dstPath := args[1]

			opts := envfile.CopyOptions{
				Overwrite: overwrite,
				Exclude:   exclude,
			}

			n, err := envfile.CopyFile(dstPath, srcPath, opts)
			if err != nil {
				return fmt.Errorf("copy failed: %w", err)
			}

			excludeNote := ""
			if len(exclude) > 0 {
				excludeNote = fmt.Sprintf(" (excluded: %s)", strings.Join(exclude, ", "))
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Copied %d key(s) from %s → %s%s\n", n, srcPath, dstPath, excludeNote)
			return nil
		},
	}

	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing keys in dst")
	cmd.Flags().StringSliceVar(&exclude, "exclude", nil, "Comma-separated list of keys to exclude")
	return cmd
}
