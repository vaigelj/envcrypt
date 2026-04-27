package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewFreezeCmd returns the top-level `freeze` command with subcommands.
func NewFreezeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "freeze",
		Short: "Freeze and verify immutable env snapshots",
	}
	cmd.AddCommand(newFreezeCreateCmd(), newFreezeShowCmd(), newFreezeDeleteCmd(), newFreezeVerifyCmd())
	return cmd
}

func newFreezeCreateCmd() *cobra.Command {
	var dir string
	cmd := &cobra.Command{
		Use:   "create <name> <file>",
		Short: "Freeze an env file under a given name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, file := args[0], args[1]
			entries, err := envfile.ParseFile(file)
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}
			f, err := envfile.Freeze(dir, name, file, entries)
			if err != nil {
				return err
			}
			fmt.Printf("Frozen %q at %s (checksum: %s)\n", name, f.FrozenAt.Format("2006-01-02T15:04:05Z"), f.Checksum[:12])
			return nil
		},
	}
	cmd.Flags().StringVar(&dir, "dir", ".", "base directory for frozen store")
	return cmd
}

func newFreezeShowCmd() *cobra.Command {
	var dir string
	cmd := &cobra.Command{
		Use:   "show <name>",
		Short: "Show a frozen env snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := envfile.LoadFrozen(dir, args[0])
			if err != nil {
				return err
			}
			fmt.Printf("Name:     %s\nSource:   %s\nFrozen:   %s\nChecksum: %s\n\n",
				args[0], f.Source, f.FrozenAt.Format("2006-01-02T15:04:05Z"), f.Checksum)
			tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			for _, e := range f.Entries {
				fmt.Fprintf(tw, "%s\t%s\n", e.Key, e.Value)
			}
			return tw.Flush()
		},
	}
	cmd.Flags().StringVar(&dir, "dir", ".", "base directory for frozen store")
	return cmd
}

func newFreezeDeleteCmd() *cobra.Command {
	var dir string
	cmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a frozen env snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := envfile.DeleteFrozen(dir, args[0]); err != nil {
				return err
			}
			fmt.Printf("Deleted frozen snapshot %q\n", args[0])
			return nil
		},
	}
	cmd.Flags().StringVar(&dir, "dir", ".", "base directory for frozen store")
	return cmd
}

func newFreezeVerifyCmd() *cobra.Command {
	var dir string
	cmd := &cobra.Command{
		Use:   "verify <name>",
		Short: "Verify a frozen snapshot has not been tampered with",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := envfile.LoadFrozen(dir, args[0])
			if err != nil {
				return err
			}
			if f.IsTampered() {
				fmt.Fprintf(os.Stderr, "TAMPERED: checksum mismatch for %q\n", args[0])
				os.Exit(1)
			}
			fmt.Printf("OK: snapshot %q is intact\n", args[0])
			return nil
		},
	}
	cmd.Flags().StringVar(&dir, "dir", ".", "base directory for frozen store")
	return cmd
}
