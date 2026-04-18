package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewDiffCmd returns a cobra command that compares two .env files and prints changes.
func NewDiffCmd() *cobra.Command {
	var redact bool

	cmd := &cobra.Command{
		Use:   "diff <base.env> <updated.env>",
		Short: "Show differences between two .env files",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			base, err := envfile.ParseFile(args[0])
			if err != nil {
				return fmt.Errorf("reading base: %w", err)
			}
			updated, err := envfile.ParseFile(args[1])
			if err != nil {
				return fmt.Errorf("reading updated: %w", err)
			}

			changes := envfile.Compare(base, updated)
			if len(changes) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No differences found.")
				return nil
			}

			for _, c := range changes {
				switch c.Type {
				case envfile.ChangeAdded:
					nv := maybeRedact(c.Key, c.NewVal, redact)
					fmt.Fprintf(cmd.OutOrStdout(), "+ %s=%s\n", c.Key, nv)
				case envfile.ChangeRemoved:
					ov := maybeRedact(c.Key, c.OldVal, redact)
					fmt.Fprintf(cmd.OutOrStdout(), "- %s=%s\n", c.Key, ov)
				case envfile.ChangeUpdated:
					ov := maybeRedact(c.Key, c.OldVal, redact)
					nv := maybeRedact(c.Key, c.NewVal, redact)
					fmt.Fprintf(cmd.OutOrStdout(), "~ %s: %s -> %s\n", c.Key, ov, nv)
				}
			}

			a, u, r := envfile.Summary(changes)
			fmt.Fprintf(cmd.OutOrStdout(), "\nSummary: +%d added, ~%d updated, -%d removed\n", a, u, r)
			return nil
		},
	}

	cmd.Flags().BoolVar(&redact, "redact", false, "Mask sensitive values in output")
	return cmd
}

func maybeRedact(key, val string, redact bool) string {
	if redact && envfile.IsSensitive(key) {
		return envfile.RedactString(val)
	}
	return val
}

func init() {
	_ = os.Stderr // ensure os import used
}
