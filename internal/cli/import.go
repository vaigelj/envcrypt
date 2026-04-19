package cli

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewImportCmd returns the cobra command for importing env vars from external files.
func NewImportCmd() *cobra.Command {
	var format string
	var output string
	var overwrite bool

	cmd := &cobra.Command{
		Use:   "import <src>",
		Short: "Import environment variables from a file into a .env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			src := args[0]
			imported, err := envfile.Import(src, format)
			if err != nil {
				return err
			}

			if output == "" {
				keys := make([]string, 0, len(imported))
				for k := range imported {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, imported[k])
				}
				return nil
			}

			existing := map[string]string{}
			if _, err := os.Stat(output); err == nil {
				existing, err = envfile.ParseFile(output)
				if err != nil {
					return fmt.Errorf("read output file: %w", err)
				}
			}

			merged, _ := envfile.Merge(existing, imported, overwrite)
			return envfile.WriteFile(output, merged)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "dotenv", "source format: dotenv, json, shell")
	cmd.Flags().StringVarP(&output, "output", "o", "", "target .env file to merge into (omit to print)")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing keys in target file")
	return cmd
}
