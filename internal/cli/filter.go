package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewFilterCmd returns the cobra command for the filter sub-command.
func NewFilterCmd() *cobra.Command {
	var (
		prefix   string
		suffix   string
		pattern  string
		keys     []string
		exclude  bool
		inPlace  bool
		output   string
	)

	cmd := &cobra.Command{
		Use:   "filter <file>",
		Short: "Filter entries in a .env file by key prefix, suffix, pattern, or name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			entries, err := envfile.ParseFile(path)
			if err != nil {
				return fmt.Errorf("reading %s: %w", path, err)
			}

			var opts []envfile.FilterOption
			if len(keys) > 0 {
				opts = append(opts, envfile.WithFilterKeys(keys...))
			}
			if prefix != "" {
				opts = append(opts, envfile.WithFilterPrefix(prefix))
			}
			if suffix != "" {
				opts = append(opts, envfile.WithFilterSuffix(suffix))
			}
			if pattern != "" {
				opts = append(opts, envfile.WithFilterPattern(pattern))
			}
			if exclude {
				opts = append(opts, envfile.WithFilterExclude())
			}

			result, err := envfile.Filter(entries, opts...)
			if err != nil {
				return err
			}

			if inPlace {
				return envfile.WriteFile(path, result)
			}

			dest := os.Stdout
			if output != "" {
				f, err := os.Create(output)
				if err != nil {
					return err
				}
				defer f.Close()
				dest = f
			}

			for _, e := range result {
				fmt.Fprintf(dest, "%s=%s\n", e.Key, e.Value)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&prefix, "prefix", "", "keep keys with this prefix")
	cmd.Flags().StringVar(&suffix, "suffix", "", "keep keys with this suffix")
	cmd.Flags().StringVar(&pattern, "pattern", "", "keep keys matching this regex")
	cmd.Flags().StringArrayVar(&keys, "key", nil, "keep only these keys (repeatable)")
	cmd.Flags().BoolVar(&exclude, "exclude", false, "invert: remove matched entries")
	cmd.Flags().BoolVarP(&inPlace, "in-place", "i", false, "overwrite source file")
	cmd.Flags().StringVarP(&output, "output", "o", "", "write result to file instead of stdout")

	return cmd
}
