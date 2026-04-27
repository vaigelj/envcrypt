package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewCompactCmd returns the cobra command for the compact sub-command.
func NewCompactCmd() *cobra.Command {
	var (
		inPlace        bool
		removeComments bool
		removeBlanks   bool
		dedupeKeys     bool
		output         string
	)

	cmd := &cobra.Command{
		Use:   "compact <file>",
		Short: "Remove comments, blank lines, and duplicate keys from an env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			var opts []envfile.CompactOption
			if removeComments {
				opts = append(opts, envfile.WithRemoveComments())
			}
			if removeBlanks {
				opts = append(opts, envfile.WithRemoveBlanks())
			}
			if dedupeKeys {
				opts = append(opts, envfile.WithDedupeKeys())
			}

			if inPlace {
				if err := envfile.CompactFile(path, opts...); err != nil {
					return fmt.Errorf("compact: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), "compacted:", path)
				return nil
			}

			entries, err := envfile.ParseFile(path)
			if err != nil {
				return fmt.Errorf("compact: %w", err)
			}
			compacted := envfile.Compact(entries, opts...)

			dest := output
			if dest == "" {
				for _, e := range compacted {
					if e.Key == "" {
						continue
					}
					fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", e.Key, e.Value)
				}
				return nil
			}
			if err := envfile.WriteFile(dest, compacted); err != nil {
				return fmt.Errorf("compact: write %s: %w", dest, err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "written:", dest)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&inPlace, "in-place", "i", false, "overwrite the source file")
	cmd.Flags().BoolVar(&removeComments, "remove-comments", true, "strip comment lines")
	cmd.Flags().BoolVar(&removeBlanks, "remove-blanks", true, "strip blank lines")
	cmd.Flags().BoolVar(&dedupeKeys, "dedupe", true, "remove duplicate keys (keep last)")
	cmd.Flags().StringVarP(&output, "output", "o", "", "write result to this file instead of stdout")

	return cmd
}
