package cli

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var keystorePath string

// Execute builds the root command and runs it.
func Execute() error {
	root := &cobra.Command{
		Use:   "envcrypt",
		Short: "Encrypt and manage .env files with team-friendly key rotation",
	}

	defaultStore := filepath.Join(os.Getenv("HOME"), ".envcrypt", "keys.json")
	root.PersistentFlags().StringVar(&keystorePath, "keystore", defaultStore, "path to keystore file")

	root.AddCommand(
		NewEncryptCmd(),
		NewDecryptCmd(),
		NewRotateCmd(),
		NewMergeCmd(),
		NewDiffCmd(),
		NewSchemaCmd(),
		NewProfileCmd(),
		NewHistoryCmd(),
		NewTagsCmd(),
		NewPinCmd(),
		NewCompareCmd(),
		NewCopyCmd(),
		NewRenameCmd(),
		NewSearchCmd(),
		NewImportCmd(),
	)

	return root.Execute()
}
