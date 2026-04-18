package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"envcrypt/internal/keystore"
	"envcrypt/internal/vault"
)

// NewEncryptCmd returns the cobra command for encrypting an .env file.
func NewEncryptCmd(ks *keystore.Store, v *vault.Vault) *cobra.Command {
	return &cobra.Command{
		Use:   "encrypt <env-file>",
		Short: "Encrypt an .env file and store the key in the keystore",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			envFile := args[0]

			keyID, err := v.EncryptFile(envFile)
			if err != nil {
				return fmt.Errorf("encrypt: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Encrypted. Key ID: %s\n", keyID)
			return nil
		},
	}
}

// NewDecryptCmd returns the cobra command for decrypting an .env file.
func NewDecryptCmd(ks *keystore.Store, v *vault.Vault) *cobra.Command {
	var keyID string

	cmd := &cobra.Command{
		Use:   "decrypt <env-file>",
		Short: "Decrypt an .env file using the stored key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			envFile := args[0]

			if err := v.DecryptFile(envFile, keyID); err != nil {
				return fmt.Errorf("decrypt: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Decrypted successfully.\n")
			return nil
		},
	}

	cmd.Flags().StringVar(&keyID, "key-id", "", "Key ID to use for decryption (defaults to latest)")
	return cmd
}
