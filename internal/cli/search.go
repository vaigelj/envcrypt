package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewSearchCmd returns a cobra command for searching keys/values in an env file.
func NewSearchCmd() *cobra.Command {
	var keyPattern, valuePattern string
	var useRegex, caseSensitive bool

	cmd := &cobra.Command{
		Use:   "search <file>",
		Short: "Search for keys or values in an env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			opts := envfile.SearchOptions{
				KeyPattern:    keyPattern,
				ValuePattern:  valuePattern,
				UseRegex:      useRegex,
				CaseSensitive: caseSensitive,
			}
			results, err := envfile.SearchFile(path, opts)
			if err != nil {
				return fmt.Errorf("search: %w", err)
			}
			if len(results) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no matches found")
				return nil
			}
			for _, r := range results {
				fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", r.Key, r.Value)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&keyPattern, "key", "k", "", "pattern to match against keys")
	cmd.Flags().StringVarP(&valuePattern, "value", "v", "", "pattern to match against values")
	cmd.Flags().BoolVarP(&useRegex, "regex", "r", false, "treat patterns as regular expressions")
	cmd.Flags().BoolVarP(&caseSensitive, "case-sensitive", "s", false, "use case-sensitive matching")

	return cmd
}
