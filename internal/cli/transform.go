package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewTransformCmd returns the cobra command for env value transformation.
func NewTransformCmd() *cobra.Command {
	var (
		keys    []string
		exclude []string
		op      string
		prefix  string
		output  string
	)

	cmd := &cobra.Command{
		Use:   "transform <file>",
		Short: "Bulk-transform values in a .env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pairs, err := envfile.ParseFile(args[0])
			if err != nil {
				return err
			}

			var fn envfile.TransformFunc
			switch strings.ToLower(op) {
			case "uppercase":
				fn = envfile.UppercaseValues()
			case "trim":
				fn = envfile.TrimValues()
			case "prefix":
				if prefix == "" {
					return fmt.Errorf("--prefix is required for prefix operation")
				}
				fn = envfile.PrefixValues(prefix)
			default:
				return fmt.Errorf("unknown operation %q; choose uppercase|trim|prefix", op)
			}

			opts := envfile.TransformOptions{Keys: keys, Exclude: exclude}
			result := envfile.Transform(pairs, fn, opts)

			dest := args[0]
			if output != "" {
				dest = output
			}
			if err := envfile.WriteFile(dest, result); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "transformed %d entries -> %s\n", len(result), dest)
			return nil
		},
	}

	cmd.Flags().StringVarP(&op, "op", "o", "trim", "transformation: uppercase|trim|prefix")
	cmd.Flags().StringVar(&prefix, "prefix", "", "prefix string (used with --op=prefix)")
	cmd.Flags().StringSliceVarP(&keys, "keys", "k", nil, "limit to these keys")
	cmd.Flags().StringSliceVarP(&exclude, "exclude", "e", nil, "exclude these keys")
	cmd.Flags().StringVar(&output, "output", "", "write result to this file instead of overwriting input")
	return cmd
}
