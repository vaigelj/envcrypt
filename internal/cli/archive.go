package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewArchiveCmd returns the root archive command with sub-commands.
func NewArchiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archive",
		Short: "Save and restore named env file archives",
	}
	cmd.AddCommand(
		newArchiveSaveCmd(),
		newArchiveLoadCmd(),
		newArchiveListCmd(),
		newArchiveDeleteCmd(),
	)
	return cmd
}

func newArchiveSaveCmd() *cobra.Command {
	var file string
	cmd := &cobra.Command{
		Use:   "save <name>",
		Short: "Archive the current env file under a name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := envfile.ParseFile(file)
			if err != nil {
				return fmt.Errorf("parse env file: %w", err)
			}
			dir, _ := os.Getwd()
			if err := envfile.SaveArchive(dir, args[0], entries); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "archived %q from %s\n", args[0], file)
			return nil
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", ".env", "env file to archive")
	return cmd
}

func newArchiveLoadCmd() *cobra.Command {
	var out string
	cmd := &cobra.Command{
		Use:   "restore <name>",
		Short: "Restore an archive to an env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := os.Getwd()
			a, err := envfile.LoadArchive(dir, args[0])
			if err != nil {
				return err
			}
			if err := envfile.WriteFile(out, a.Entries); err != nil {
				return fmt.Errorf("write env file: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "restored archive %q to %s\n", args[0], out)
			return nil
		},
	}
	cmd.Flags().StringVarP(&out, "out", "o", ".env", "destination env file")
	return cmd
}

func newArchiveListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all saved archives",
		RunE: func(cmd *cobra.Command, _ []string) error {
			dir, _ := os.Getwd()
			names, err := envfile.ListArchives(dir)
			if err != nil {
				return err
			}
			if len(names) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no archives found")
				return nil
			}
			for _, n := range names {
				fmt.Fprintln(cmd.OutOrStdout(), n)
			}
			return nil
		},
	}
}

func newArchiveDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a named archive",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := os.Getwd()
			if err := envfile.DeleteArchive(dir, args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "deleted archive %q\n", args[0])
			return nil
		},
	}
}
