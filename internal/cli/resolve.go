package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewResolveCmd returns the cobra command for the `resolve` subcommand.
func NewResolveCmd() *cobra.Command {
	var (
		strict  bool
		useEnv  bool
		outFmt  string
	)

	cmd := &cobra.Command{
		Use:   "resolve <file>",
		Short: "Expand variable references in a .env file",
		Long: `Reads a .env file and expands $VAR / ${VAR} references found in values.

By default unresolved references are left unchanged (loose mode).
Pass --strict to treat any unresolved reference as a fatal error.
Pass --env to also consult the current process environment as a fallback.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			mode := envfile.ResolveModeLoose
			if strict {
				mode = envfile.ResolveModeStrict
			}

			opts := envfile.ResolveOptions{
				Mode:    mode,
				Environ: useEnv,
			}

			entries, err := envfile.ResolveFile(path, opts)
			if err != nil {
				return fmt.Errorf("resolve: %w", err)
			}

			switch outFmt {
			case "dotenv", "":
				for _, e := range entries {
					fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", e.Key, e.Value)
				}
			case "export":
				for _, e := range entries {
					fmt.Fprintf(cmd.OutOrStdout(), "export %s=%q\n", e.Key, e.Value)
				}
			default:
				return fmt.Errorf("unknown format %q (supported: dotenv, export)", outFmt)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&strict, "strict", false, "fail on unresolved variable references")
	cmd.Flags().BoolVar(&useEnv, "env", false, "fall back to OS environment variables")
	cmd.Flags().StringVar(&outFmt, "format", "dotenv", "output format: dotenv|export")

	return cmd
}

func init() {
	// Register with root in root.go via rootCmd.AddCommand.
	_ = os.Getenv // suppress unused import lint
}
