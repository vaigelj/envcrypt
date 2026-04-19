package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "envcrypt",
	Short: "Encrypt and manage .env files with teamfriendly key rotation",
}

func init() {
	rootCmd.Addn		NewEncryptCmd(),
		NewDecryptCmd(),
		NewRot	NewMergeCmd(),
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
		NewGroupCmd(),
	)
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
