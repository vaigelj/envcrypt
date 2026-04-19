// Package cli wires together cobra commands and shared dependencies.
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envcrypt/internal/keystore"
	"envcrypt/internal/vault"
)

// Execute builds the root command tree and runs it.
func Execute() {
	ksPath := keystorePath()

	ks, err := keystore.New(ksPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening keystore: %v\n", err)
		os.Exit(1)
	}

	v := vault.New(ks)

	root := &cobra.Command{
		Use:   "envcrypt",
		Short: "Encrypt and manage .env files with team-friendly key rotation",
	}

	root.AddCommand(
		NewEncryptCmd(ks, v),
		NewDecryptCmd(ks, v),
		NewRotateCmd(ks, v),
	)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

// keystorePath returns the path to the keystore file. It checks the
// ENVCRYPT_KEYSTORE environment variable first, falling back to the default
// path ".envcrypt_keys" in the current working directory.
func keystorePath() string {
	if p := os.Getenv("ENVCRYPT_KEYSTORE"); p != "" {
		return p
	}
	return ".envcrypt_keys"
}
