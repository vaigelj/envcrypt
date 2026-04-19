package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewCompareCmd returns a cobra command to compare two .env files as versions.
func NewCompareCmd() *cobra.Command {
	var fromLabel string
	var toLabel string
	var redact bool

	cmd := &cobra.Command{
		Use:   "compare <from-file> <to-file>",
		Short: "Compare two .env files and show a version diff",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fromVars, err := envfile.ParseFile(args[0])
			if err != nil {
				return fmt.Errorf("reading %s: %w", args[0], err)
			}
			toVars, err := envfile.ParseFile(args[1])
			if err != nil {
				return fmt.Errorf("reading %s: %w", args[1], err)
			}

			if redact {
				fromVars = envfile.Redact(fromVars, "***")
				toVars = envfile.Redact(toVars, "***")
			}

			if fromLabel == "" {
				fromLabel = args[0]
			}
			if toLabel == "" {
				toLabel = args[1]
			}

			fromV := envfile.Version{Name: fromLabel, Vars: fromVars}
			toV := envfile.Version{Name: toLabel, Vars: toVars}
			diff := envfile.CompareVersions(fromV, toV)
			fmt.Fprint(os.Stdout, envfile.FormatVersionDiff(diff))
			return nil
		},
	}

	cmd.Flags().StringVar(&fromLabel, "from-label", "", "Label for the from file")
	cmd.Flags().StringVar(&toLabel, "to-label", "", "Label for the to file")
	cmd.Flags().BoolVar(&redact, "redact", false, "Redact sensitive values before comparing")
	return cmd
}
