package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envcrypt/internal/envfile"
)

// NewReorderCmd returns the cobra command for the `reorder` sub-command.
func NewReorderCmd() *cobra.Command {
	var (
		inPlace    bool
		missingOk  bool
		outputPath string
	)

	cmd := &cobra.Command{
		Use:   "reorder <file> <KEY,KEY,...>",
		Short: "Reorder keys in a .env file",
		Long: `Reorder rearranges the entries in a .env file so that the keys
listed in the second argument appear first (in that order). Any keys
not mentioned keep their original relative position after the pinned ones.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			order := splitCSV(args[1])

			var opts []envfile.ReorderOption
			if missingOk {
				opts = append(opts, envfile.WithReorderMissingOk())
			}

			if inPlace {
				return envfile.ReorderFile(path, order, opts...)
			}

			entries, err := envfile.ParseFile(path)
			if err != nil {
				return err
			}
			reordered, err := envfile.Reorder(entries, order, opts...)
			if err != nil {
				return err
			}

			dest := outputPath
			if dest == "" {
				dest = "-"
			}
			if dest == "-" {
				for _, e := range reordered {
					cmd.Printf("%s=%s\n", e.Key, e.Value)
				}
				return nil
			}
			return envfile.WriteFile(dest, reordered)
		},
	}

	cmd.Flags().BoolVarP(&inPlace, "in-place", "i", false, "Rewrite the source file in place")
	cmd.Flags().BoolVar(&missingOk, "missing-ok", false, "Silently skip order keys absent from the file")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Write result to this file (default: stdout)")

	return cmd
}

// splitCSV splits a comma-separated string and trims whitespace.
func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
