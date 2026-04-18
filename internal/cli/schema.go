package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewSchemaCmd returns a command that validates an env file against a simple
// inline schema defined via flags.
func NewSchemaCmd() *cobra.Command {
	var required []string
	var patterns []string

	cmd := &cobra.Command{
		Use:   "schema [env-file]",
		Short: "Validate an env file against a schema",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.ParseFile(args[0])
			if err != nil {
				return fmt.Errorf("reading env file: %w", err)
			}

			fields := make([]envfile.SchemaField, 0, len(required)+len(patterns))
			for _, k := range required {
				fields = append(fields, envfile.SchemaField{Key: k, Required: true})
			}
			for _, kp := range patterns {
				var key, pat string
				fmt.Sscanf(kp, "%s %s", &key, &pat)
				fields = append(fields, envfile.SchemaField{Key: key, Pattern: pat})
			}

			schema := &envfile.Schema{Fields: fields}
			errors := schema.Validate(env)
			if len(errors) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "schema validation passed")
				return nil
			}
			for _, e := range errors {
				fmt.Fprintln(os.Stderr, e.Error())
			}
			return fmt.Errorf("%d schema violation(s) found", len(errors))
		},
	}

	cmd.Flags().StringArrayVar(&required, "require", nil, "required key (repeatable)")
	cmd.Flags().StringArrayVar(&patterns, "pattern", nil, "KEY PATTERN pair to validate value (repeatable)")
	return cmd
}
