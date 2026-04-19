package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewGenerateCmd returns the 'generate' subcommand.
func NewGenerateCmd() *cobra.Command {
	var length int
	var noSymbols bool
	var numeric bool
	var keys []string
	var output string

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate random values for env keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := envfile.GenerateOptions{
				Length:    length,
				NoSymbols: noSymbols,
				Numeric:   numeric,
			}

			if len(keys) == 0 {
				// generate a single value and print it
				v, err := envfile.GenerateValue(opts)
				if err != nil {
					return err
				}
				fmt.Println(v)
				return nil
			}

			generated, err := envfile.GenerateForKeys(keys, opts)
			if err != nil {
				return err
			}

			if output != "" {
				existing := map[string]string{}
				_ = func() {
					parsed, e := envfile.ParseFile(output)
					if e == nil {
						existing = parsed
					}
				}
				for k, v := range generated {
					existing[k] = v
				}
				if err := envfile.WriteFile(output, existing); err != nil {
					return fmt.Errorf("write %s: %w", output, err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "wrote %d key(s) to %s\n", len(generated), output)
				return nil
			}

			for _, k := range keys {
				fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, generated[k])
			}
			_ = strings.Join // satisfy import
			return nil
		},
	}

	cmd.Flags().IntVarP(&length, "length", "l", 32, "length of generated value")
	cmd.Flags().BoolVar(&noSymbols, "no-symbols", false, "exclude symbol characters")
	cmd.Flags().BoolVar(&numeric, "numeric", false, "digits only")
	cmd.Flags().StringArrayVarP(&keys, "key", "k", nil, "key name(s) to generate values for")
	cmd.Flags().StringVarP(&output, "output", "o", "", "write generated values to this .env file")
	return cmd
}
