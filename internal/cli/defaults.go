package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewDefaultsCmd returns the cobra command for the `defaults` sub-command.
func NewDefaultsCmd() *cobra.Command {
	var (
		overwrite  bool
		schemaFile string
		outFile    string
	)

	cmd := &cobra.Command{
		Use:   "defaults <env-file>",
		Short: "Apply default values to a .env file",
		Long: `Apply a set of default key=value pairs to an existing .env file.
Keys that already exist are left untouched unless --overwrite is set.

Defaults can be supplied inline via --set KEY=VALUE flags or loaded from
a JSON schema file via --schema.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			var defs []envfile.DefaultEntry

			// Load from schema file if provided.
			if schemaFile != "" {
				data, err := os.ReadFile(schemaFile)
				if err != nil {
					return fmt.Errorf("defaults: reading schema: %w", err)
				}
				if err := json.Unmarshal(data, &defs); err != nil {
					return fmt.Errorf("defaults: parsing schema: %w", err)
				}
			}

			// Inline --set KEY=VALUE flags.
			sets, _ := cmd.Flags().GetStringArray("set")
			for _, s := range sets {
				parts := strings.SplitN(s, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("defaults: invalid --set value %q (expected KEY=VALUE)", s)
				}
				defs = append(defs, envfile.DefaultEntry{Key: parts[0], Default: parts[1]})
			}

			if len(defs) == 0 {
				return fmt.Errorf("defaults: no defaults specified; use --set or --schema")
			}

			var opts []envfile.ApplyDefaultsOption
			if overwrite {
				opts = append(opts, envfile.WithDefaultsOverwrite())
			}

			target := path
			if outFile != "" {
				target = outFile
				// Copy source to target first so ApplyDefaultsFile can read it.
				data, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				if err := os.WriteFile(target, data, 0o600); err != nil {
					return err
				}
			}

			applied, err := envfile.ApplyDefaultsFile(target, defs, opts...)
			if err != nil {
				return err
			}

			if len(applied) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no defaults applied (all keys already present)")
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "applied defaults: %s\n", strings.Join(applied, ", "))
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing keys with default values")
	cmd.Flags().StringVar(&schemaFile, "schema", "", "JSON file containing default entries")
	cmd.Flags().StringVar(&outFile, "out", "", "write result to a different file (default: in-place)")
	cmd.Flags().StringArray("set", nil, "inline default as KEY=VALUE (repeatable)")

	return cmd
}

func init() {
	rootCmd.AddCommand(NewDefaultsCmd())
}
