package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewInterpolateCmd returns the cobra command for in-place variable interpolation.
func NewInterpolateCmd() *cobra.Command {
	var strict bool
	var useEnviron bool
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "interpolate <file>",
		Short: "Expand $VAR and ${VAR} references inside an env file",
		Long: `Reads the given .env file and expands variable references of the form
$VAR or ${VAR} using values defined within the same file.

Use --environ to also fall back to OS environment variables.
Use --strict to treat any unresolved reference as an error.
Use --dry-run to print the result without modifying the file.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			opts := envfile.InterpolateOptions{
				Strict:  strict,
				Environ: useEnviron,
			}

			entries, err := envfile.ParseFile(path)
			if err != nil {
				return fmt.Errorf("reading %s: %w", path, err)
			}

			resolved, err := envfile.Interpolate(entries, opts)
			if err != nil {
				return fmt.Errorf("interpolation failed: %w", err)
			}

			if dryRun {
				for _, e := range resolved {
					fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", e.Key, e.Value)
				}
				return nil
			}

			if err := envfile.WriteFile(path, resolved); err != nil {
				return fmt.Errorf("writing %s: %w", path, err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "interpolated %d entries in %s\n", len(resolved), path)
			return nil
		},
	}

	cmd.Flags().BoolVar(&strict, "strict", false, "error on unresolved variable references")
	cmd.Flags().BoolVar(&useEnviron, "environ", false, "fall back to OS environment variables")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print result without modifying the file")
	return cmd
}
