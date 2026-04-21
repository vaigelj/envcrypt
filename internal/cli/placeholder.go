package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewPlaceholderCmd returns the `envcrypt placeholder` command which resolves
// {{KEY}} tokens inside a .env file using sibling values or an optional
// override map supplied via --set KEY=VALUE flags.
func NewPlaceholderCmd() *cobra.Command {
	var (
		setFlags []string
		strict   bool
		outFmt   string
	)

	cmd := &cobra.Command{
		Use:   "placeholder <file>",
		Short: "Resolve {{PLACEHOLDER}} tokens in a .env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := envfile.ParseFile(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			overrides := parseCSVSet(setFlags)

			resolved, err := envfile.ResolvePlaceholders(entries, overrides, strict)
			if err != nil {
				return err
			}

			switch outFmt {
			case "json":
				m := make(map[string]string, len(resolved))
				for _, e := range resolved {
					m[e.Key] = e.Value
				}
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(m)
			default:
				for _, e := range resolved {
					fmt.Fprintf(os.Stdout, "%s=%s\n", e.Key, e.Value)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringArrayVar(&setFlags, "set", nil, "Override values: KEY=VALUE (repeatable)")
	cmd.Flags().BoolVar(&strict, "strict", false, "Fail if any placeholder cannot be resolved")
	cmd.Flags().StringVar(&outFmt, "format", "dotenv", "Output format: dotenv|json")
	return cmd
}
