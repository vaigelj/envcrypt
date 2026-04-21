package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewInjectCmd returns a cobra command that injects variables from a .env file
// into the current shell session by printing export statements.
func NewInjectCmd() *cobra.Command {
	var (
		overwrite bool
		prefix    string
		only      string
		shell     bool
	)

	cmd := &cobra.Command{
		Use:   "inject [file]",
		Short: "Inject variables from a .env file into the environment",
		Long: `Reads a .env file and prints export statements suitable for eval.

Example:
  eval $(envcrypt inject .env)
  eval $(envcrypt inject --prefix APP_ .env.production)`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			entries, err := envfile.ParseFile(path)
			if err != nil {
				return fmt.Errorf("inject: %w", err)
			}

			onlySet := parseCSVSet(only)

			for _, e := range entries {
				key := e.Key
				if onlySet != nil && !onlySet[key] {
					continue
				}
				if prefix != "" {
					key = prefix + key
				}
				if !overwrite {
					// Emit conditional export (no-clobber)
					if shell {
						fmt.Fprintf(cmd.OutOrStdout(), "[ -z \"${%s+x}\" ] && export %s=%q\n", key, key, e.Value)
					} else {
						fmt.Fprintf(cmd.OutOrStdout(), "export %s=%q\n", key, e.Value)
					}
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "export %s=%q\n", key, e.Value)
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing environment variables")
	cmd.Flags().StringVar(&prefix, "prefix", "", "Prefix to prepend to all keys")
	cmd.Flags().StringVar(&only, "only", "", "Comma-separated list of keys to inject")
	cmd.Flags().BoolVar(&shell, "no-clobber", false, "Emit conditional exports (do not overwrite existing vars)")

	return cmd
}

// parseCSVSet converts a comma-separated string to a set map.
func parseCSVSet(csv string) map[string]bool {
	if csv == "" {
		return nil
	}
	m := make(map[string]bool)
	start := 0
	for i := 0; i <= len(csv); i++ {
		if i == len(csv) || csv[i] == ',' {
			k := trimSpace(csv[start:i])
			if k != "" {
				m[k] = true
			}
			start = i + 1
		}
	}
	return m
}

func trimSpace(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t') {
		s = s[:len(s)-1]
	}
	return s
}
