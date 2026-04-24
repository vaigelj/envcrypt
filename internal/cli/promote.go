package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewPromoteCmd returns the cobra command for the promote sub-command.
func NewPromoteCmd() *cobra.Command {
	var overwrite bool
	var exclude []string

	cmd := &cobra.Command{
		Use:   "promote <src> <dst>",
		Short: "Promote env vars from one file to another",
		Long: `Copy environment variables from a source .env file into a destination
.env file. Keys that already exist in the destination are reported as
conflicts unless --overwrite is set.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			srcPath := args[0]
			dstPath := args[1]

			var opts []envfile.PromoteOption
			if overwrite {
				opts = append(opts, envfile.WithPromoteOverwrite())
			}
			if len(exclude) > 0 {
				opts = append(opts, envfile.WithPromoteExclude(exclude...))
			}

			res, err := envfile.PromoteFile(srcPath, dstPath, opts...)
			if err != nil {
				return err
			}

			if len(res.Promoted) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "promoted: %s\n", strings.Join(res.Promoted, ", "))
			}
			if len(res.Skipped) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "skipped:  %s\n", strings.Join(res.Skipped, ", "))
			}
			if len(res.Conflict) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "conflict: %s (use --overwrite to replace)\n",
					strings.Join(res.Conflict, ", "))
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing keys in destination")
	cmd.Flags().StringSliceVar(&exclude, "exclude", nil, "comma-separated list of keys to skip")
	return cmd
}

func init() {
	rootCmd.AddCommand(NewPromoteCmd())
}
