package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewBackupCmd returns the root backup command with subcommands.
func NewBackupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Manage env file backups",
	}
	cmd.AddCommand(newBackupCreateCmd())
	cmd.AddCommand(newBackupListCmd())
	cmd.AddCommand(newBackupRestoreCmd())
	cmd.AddCommand(newBackupDeleteCmd())
	return cmd
}

func newBackupCreateCmd() *cobra.Command {
	var label string
	cmd := &cobra.Command{
		Use:   "create <file>",
		Short: "Create a backup of an env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := envfile.ParseFile(args[0])
			if err != nil {
				return err
			}
			dir := "."
			b, err := envfile.CreateBackup(dir, entries, label)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "backup created: %s\n", b.ID)
			return nil
		},
	}
	cmd.Flags().StringVarP(&label, "label", "l", "", "optional label for the backup")
	return cmd
}

func newBackupListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all backups",
		RunE: func(cmd *cobra.Command, args []string) error {
			list, err := envfile.ListBackups(".")
			if err != nil {
				return err
			}
			if len(list) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no backups found")
				return nil
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tCREATED\tLABEL")
			for _, b := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\n", b.ID, b.CreatedAt.Format("2006-01-02 15:04:05"), b.Label)
			}
			return w.Flush()
		},
	}
}

func newBackupRestoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restore <id> <dest-file>",
		Short: "Restore a backup to a file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := envfile.LoadBackup(".", args[0])
			if err != nil {
				return err
			}
			if err := envfile.WriteFile(args[1], b.Entries); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "restored backup %s → %s\n", args[0], args[1])
			return nil
		},
	}
}

func newBackupDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a backup by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := envfile.DeleteBackup(".", args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "deleted backup %s\n", args[0])
			return nil
		},
	}
}

func init() {
	_ = os.Stderr // satisfy import
}
