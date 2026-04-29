package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewSignCmd returns the parent 'sign' command with sub-commands.
func NewSignCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign",
		Short: "Sign and verify .env files with an HMAC signature",
	}
	cmd.AddCommand(newSignCreateCmd(), newSignVerifyCmd())
	return cmd
}

func newSignCreateCmd() *cobra.Command {
	var keyHex string
	var outPath string

	cmd := &cobra.Command{
		Use:   "create <envfile>",
		Short: "Create a signed envelope from a .env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := resolveSignKey(keyHex)
			if err != nil {
				return err
			}
			out := outPath
			if out == "" {
				out = args[0] + ".sig"
			}
			if err := envfile.SignFile(args[0], out, key); err != nil {
				return fmt.Errorf("sign create: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "signed envelope written to %s\n", out)
			return nil
		},
	}
	cmd.Flags().StringVar(&keyHex, "key", "", "HMAC key (raw string; prefer ENVCRYPT_SIGN_KEY env var)")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output path for the signed envelope (default: <envfile>.sig)")
	return cmd
}

func newSignVerifyCmd() *cobra.Command {
	var keyHex string

	cmd := &cobra.Command{
		Use:   "verify <envelope>",
		Short: "Verify a signed envelope",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := resolveSignKey(keyHex)
			if err != nil {
				return err
			}
			env, err := envfile.VerifyFile(args[0], key)
			if err != nil {
				return fmt.Errorf("sign verify: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "OK — signature valid (signed at %s, %d entries)\n",
				env.SignedAt.Format("2006-01-02T15:04:05Z"), len(env.Entries))
			return nil
		},
	}
	cmd.Flags().StringVar(&keyHex, "key", "", "HMAC key (raw string; prefer ENVCRYPT_SIGN_KEY env var)")
	return cmd
}

func resolveSignKey(flag string) ([]byte, error) {
	if flag != "" {
		return []byte(flag), nil
	}
	if v := os.Getenv("ENVCRYPT_SIGN_KEY"); v != "" {
		return []byte(v), nil
	}
	return nil, fmt.Errorf("HMAC key required: use --key flag or set ENVCRYPT_SIGN_KEY")
}
