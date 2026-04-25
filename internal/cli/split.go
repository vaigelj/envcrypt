package cli

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/yourorg/envcrypt/internal/envfile"
)

// NewSplitCmd returns the cobra command for the split sub-command.
func NewSplitCmd() *cobra.Command {
	var (
		outDir    string
		sep       string
		overwrite bool
	)

	cmd := &cobra.Command{
		Use:   "split <file>",
		Short: "Split an env file into per-prefix files",
		Long: `Split reads a .env file and writes one output file per key-prefix group.

The prefix is the portion of the key before the first separator (default "_").
Keys without a separator are collected into a file named _default.env.

Example:
  envcrypt split .env --out ./envs
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			src := args[0]

			var opts []envfile.SplitOption
			if sep != "" {
				opts = append(opts, envfile.WithSplitSeparator(sep))
			}
			if overwrite {
				opts = append(opts, envfile.WithSplitOverwrite())
			}

			written, err := envfile.SplitFile(src, outDir, opts...)
			if err != nil {
				return err
			}

			sort.Strings(written)
			fmt.Fprintf(cmd.OutOrStdout(), "wrote %d file(s):\n", len(written))
			for _, f := range written {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", f)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&outDir, "out", "o", ".", "output directory for split files")
	cmd.Flags().StringVar(&sep, "sep", "_", "key prefix separator")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing output files")

	return cmd
}
