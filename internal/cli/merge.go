package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewMergeCmd returns a cobra command that merges two plain .env files.
func NewMergeCmd() *cobra.Command {
	var preferOverride bool
	var outputPath string

	cmd := &cobra.Command{
		Use:   "merge <base.env> <override.env>",
		Short: "Merge two .env files, writing the result to stdout or a file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			base, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("reading base: %w", err)
			}
			over, err := envfile.Parse(args[1])
			if err != nil {
				return fmt.Errorf("reading override: %w", err)
			}

			strategy := envfile.PreferBase
			if preferOverride {
				strategy = envfile.PreferOverride
			}
			merged := envfile.Merge(base, over, strategy)

			w := cmd.OutOrStdout()
			if outputPath != "" {
				f, err := os.Create(outputPath)
				if err != nil {
					return fmt.Errorf("creating output file: %w", err)
				}
				defer f.Close()
				w = f
			}
			return envfile.Write(w, merged)
		},
	}

	cmd.Flags().BoolVarP(&preferOverride, "override", "O", false,
		"override values from the second file take precedence on conflict")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "",
		"write merged result to file instead of stdout")
	return cmd
}
