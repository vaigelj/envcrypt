package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewCloneCmd returns the cobra command for `envcrypt clone`.
func NewCloneCmd() *cobra.Command {
	var (
		keys      string
		overwrite bool
		strip     bool
	)

	cmd := &cobra.Command{
		Use:   "clone <src> <dst>",
		Short: "Clone an env file to a new destination",
		Long: `Clone copies a .env file to a new path.

Optionally restrict which keys are copied (--keys) and strip all values to
produce a template file (--strip).`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			src, dst := args[0], args[1]

			opts := envfile.CloneOptions{
				Overwrite:   overwrite,
				StripValues: strip,
			}

			if keys != "" {
				for _, k := range strings.Split(keys, ",") {
					k = strings.TrimSpace(k)
					if k != "" {
						opts.Keys = append(opts.Keys, k)
					}
				}
			}

			if err := envfile.CloneFile(src, dst, opts); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "cloned %q → %q\n", src, dst)
			return nil
		},
	}

	cmd.Flags().StringVar(&keys, "keys", "", "comma-separated list of keys to include (default: all)")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite destination if it already exists")
	cmd.Flags().BoolVar(&strip, "strip", false, "strip all values (produce a template)")

	return cmd
}
