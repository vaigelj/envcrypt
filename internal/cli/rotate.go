package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"envcrypt/internal/keystore"
	"envcrypt/internal/vault"
)

// NewRotateCmd returns the cobra command for rotating the encryption key.
func NewRotateCmd(ks *keystore.Store, v *vault.Vault) *cobra.Command {
	return &cobra.Command{
		Use:   "rotate <env-file>",
		Short: "Re-encrypt an .env file under a freshly generated key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			envFile := args[0]

			newKeyID, err := v.RotateKey(envFile)
			if err != nil {
				return fmt.Errorf("rotate: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Key rotated successfully. New key ID: %s\n", newKeyID)
			return nil
		},
	}
}
