package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "envcrypt",
	Short: "Encrypt and manage .env files with team-friendly key rotation",
}

func init() {
	rootCmd.AddCommand(NewEncryptCmd())
	rootCmd.AddCommand(NewDecryptCmd())
	rootCmd.AddCommand(NewRotateCmd())
	rootCmd.AddCommand(NewMergeCmd())
	rootCmd.AddCommand(NewDiffCmd())
	rootCmd.AddCommand(NewSchemaCmd())
	rootCmd.AddCommand(NewProfileCmd())
	rootCmd.AddCommand(NewHistoryCmd())
	rootCmd.AddCommand(NewTagsCmd())
	rootCmd.AddCommand(NewPinCmd())
	rootCmd.AddCommand(NewCompareCmd())
	rootCmd.AddCommand(NewCopyCmd())
	rootCmd.AddCommand(NewRenameCmd())
	rootCmd.AddCommand(NewSearchCmd())
	rootCmd.AddCommand(NewImportCmd())
	rootCmd.AddCommand(NewGroupCmd())
	rootCmd.AddCommand(NewTransformCmd())
	rootCmd.AddCommand(NewGenerateCmd())
	rootCmd.AddCommand(NewInjectCmd())
	rootCmd.AddCommand(NewEncryptFieldCmd())
	rootCmd.AddCommand(newDecryptFieldCmd())
	rootCmd.AddCommand(NewChainCmd())
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
