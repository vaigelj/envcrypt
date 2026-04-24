package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewConvertCmd returns the cobra command for format conversion.
func NewConvertCmd() *cobra.Command {
	var from, to, output string

	cmd := &cobra.Command{
		Use:   "convert <file>",
		Short: "Convert an env file between formats (dotenv, json, shell)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			src := args[0]

			from = strings.ToLower(from)
			to = strings.ToLower(to)

			if output == "" {
				// Print to stdout
				entries, err := envfile.Import(src, from)
				if err != nil {
					return fmt.Errorf("read %s: %w", src, err)
				}
				out, err := envfile.ConvertFormat(entries, from, to)
				if err != nil {
					return err
				}
				fmt.Fprint(cmd.OutOrStdout(), out)
				return nil
			}

			if err := envfile.ConvertFile(src, output, from, to); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "converted %s (%s) → %s (%s)\n", src, from, output, to)
			return nil
		},
	}

	cmd.Flags().StringVar(&from, "from", "dotenv", "source format: dotenv, json, shell")
	cmd.Flags().StringVar(&to, "to", "json", "target format: dotenv, json, shell")
	cmd.Flags().StringVarP(&output, "output", "o", "", "write output to file instead of stdout")

	return cmd
}
