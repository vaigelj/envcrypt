package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yourorg/envcrypt/internal/envfile"
)

// NewRenameCmd returns a cobra command that renames a key inside a .env file.
func NewRenameCmd() *cobra.Command {
	var overwrite bool

	cmd := &cobra.Command{
		Use:   "rename <file> <old-key> <new-key>",
		Short: "Rename a key in a .env file",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, oldKey, newKey := args[0], args[1], args[2]

			if err := envfile.RenameFile(path, oldKey, newKey, overwrite); err != nil {
				return fmt.Errorf("rename: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Renamed %q to %q in %s\n", oldKey, newKey, path)
			return nil
		},
	}

	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite new-key if it already exists")
	return cmd
}
