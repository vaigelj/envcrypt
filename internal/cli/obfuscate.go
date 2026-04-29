package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewObfuscateCmd returns the parent command for obfuscate/deobfuscate.
func NewObfuscateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "obfuscate",
		Short: "Base64-obfuscate or deobfuscate .env values",
	}
	cmd.AddCommand(newObfuscateEncodeCmd())
	cmd.AddCommand(newObfuscateDecodeCmd())
	return cmd
}

func newObfuscateEncodeCmd() *cobra.Command {
	var (
		keys    []string
		prefix  string
		padding int
		inPlace bool
	)
	cmd := &cobra.Command{
		Use:   "encode <file>",
		Short: "Obfuscate values in a .env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := envfile.ParseFile(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}
			var opts []envfile.ObfuscateOption
			if len(keys) > 0 {
				opts = append(opts, envfile.WithObfuscateKeys(keys))
			}
			if prefix != "" {
				opts = append(opts, envfile.WithObfuscatePrefix(prefix))
			}
			if padding > 0 {
				opts = append(opts, envfile.WithObfuscatePadding(padding))
			}
			out, err := envfile.Obfuscate(entries, opts...)
			if err != nil {
				return err
			}
			dest := os.Stdout
			if inPlace {
				f, err := os.Create(args[0])
				if err != nil {
					return err
				}
				defer f.Close()
				dest = f
			}
			for _, e := range out {
				if e.Blank {
					fmt.Fprintln(dest)
					continue
				}
				if e.Comment {
					fmt.Fprintln(dest, e.RawLine)
					continue
				}
				fmt.Fprintf(dest, "%s=%s\n", e.Key, e.Value)
			}
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&keys, "keys", nil, "comma-separated keys to obfuscate (default: all)")
	cmd.Flags().StringVar(&prefix, "prefix", "", "marker prefix for obfuscated values")
	cmd.Flags().IntVar(&padding, "padding", 0, "random padding bytes to prepend before encoding")
	cmd.Flags().BoolVarP(&inPlace, "in-place", "i", false, "write result back to the source file")
	return cmd
}

func newObfuscateDecodeCmd() *cobra.Command {
	var (
		keys    []string
		prefix  string
		padding int
	)
	cmd := &cobra.Command{
		Use:   "decode <file>",
		Short: "Deobfuscate values in a .env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := envfile.ParseFile(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}
			var opts []envfile.ObfuscateOption
			if len(keys) > 0 {
				opts = append(opts, envfile.WithObfuscateKeys(keys))
			}
			if prefix != "" {
				opts = append(opts, envfile.WithObfuscatePrefix(prefix))
			}
			if padding > 0 {
				opts = append(opts, envfile.WithObfuscatePadding(padding))
			}
			out, err := envfile.Deobfuscate(entries, opts...)
			if err != nil {
				return err
			}
			for _, e := range out {
				if e.Blank {
					fmt.Println()
					continue
				}
				if e.Comment {
					fmt.Println(strings.TrimRight(e.RawLine, "\n"))
					continue
				}
				fmt.Printf("%s=%s\n", e.Key, e.Value)
			}
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&keys, "keys", nil, "comma-separated keys to decode (default: all)")
	cmd.Flags().StringVar(&prefix, "prefix", "", "marker prefix used during encoding")
	cmd.Flags().IntVar(&padding, "padding", 0, "padding bytes that were prepended during encoding")
	return cmd
}
