package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envcrypt/internal/envfile"
	"github.com/user/envcrypt/internal/keystore"
)

// NewEncryptFieldCmd returns a command that encrypts specific fields inside a
// plain .env file in-place, using the active key from the keystore.
func NewEncryptFieldCmd(ks *keystore.Store) *cobra.Command {
	var (
		filePath string
		keyName  string
		fields   []string
	)
	cmd := &cobra.Command{
		Use:   "encrypt-field",
		Short: "Encrypt specific fields inside a .env file",
		Example: `  envcrypt encrypt-field --file .env --key mykey --fields DB_PASS,API_SECRET`,
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := ks.Get(keyName)
			if err != nil {
				return fmt.Errorf("keystore get %q: %w", keyName, err)
			}
			entries, err := envfile.ParseFile(filePath)
			if err != nil {
				return fmt.Errorf("parse %s: %w", filePath, err)
			}
			encrypted, err := envfile.EncryptFields(entries, key, fields...)
			if err != nil {
				return fmt.Errorf("encrypt fields: %w", err)
			}
			if err := envfile.WriteFile(filePath, encrypted); err != nil {
				return fmt.Errorf("write %s: %w", filePath, err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Encrypted fields [%s] in %s\n",
				strings.Join(fields, ", "), filePath)
			return nil
		},
	}
	cmd.Flags().StringVarP(&filePath, "file", "f", ".env", "path to .env file")
	cmd.Flags().StringVarP(&keyName, "key", "k", "", "key name in keystore (required)")
	cmd.Flags().StringSliceVar(&fields, "fields", nil, "comma-separated list of fields to encrypt")
	_ = cmd.MarkFlagRequired("key")

	cmd.AddCommand(newDecryptFieldCmd(ks))
	return cmd
}

func newDecryptFieldCmd(ks *keystore.Store) *cobra.Command {
	var (
		filePath string
		keyName  string
	)
	return &cobra.Command{
		Use:   "decrypt",
		Short: "Decrypt all encrypted fields in a .env file",
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := ks.Get(keyName)
			if err != nil {
				return fmt.Errorf("keystore get %q: %w", keyName, err)
			}
			entries, err := envfile.ParseFile(filePath)
			if err != nil {
				return fmt.Errorf("parse %s: %w", filePath, err)
			}
			decrypted, err := envfile.DecryptFields(entries, key)
			if err != nil {
				return fmt.Errorf("decrypt fields: %w", err)
			}
			if err := envfile.WriteFile(filePath, decrypted); err != nil {
				return fmt.Errorf("write %s: %w", filePath, err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Decrypted all encrypted fields in %s\n", filePath)
			return nil
		},
	}
}
