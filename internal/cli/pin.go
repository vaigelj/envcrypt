package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewPinCmd returns the root pin command with subcommands.
func NewPinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pin",
		Short: "Manage named env file pins",
	}
	cmd.AddCommand(newPinSaveCmd(), newPinShowCmd(), newPinListCmd(), newPinDeleteCmd())
	return cmd
}

func newPinSaveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "save <name> <env-file>",
		Short: "Pin the current state of an env file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, file := args[0], args[1]
			vals, err := envfile.ParseFile(file)
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}
			dir := keystorePath()
			if err := envfile.SavePin(dir, name, vals); err != nil {
				return fmt.Errorf("save pin: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "pinned %q from %s\n", name, file)
			return nil
		},
	}
}

func newPinShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Show values stored in a pin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pin, err := envfile.LoadPin(keystorePath(), args[0])
			if err != nil {
				return err
			}
			if pin == nil {
				return fmt.Errorf("pin %q not found", args[0])
			}
			fmt.Fprintf(cmd.OutOrStdout(), "# pin: %s  created: %s\n", pin.Name, pin.CreatedAt.Format("2006-01-02 15:04:05"))
			for k, v := range pin.Values {
				fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
			}
			return nil
		},
	}
}

func newPinListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all pins",
		RunE: func(cmd *cobra.Command, _ []string) error {
			names, err := envfile.ListPins(keystorePath())
			if err != nil {
				return err
			}
			if len(names) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no pins found")
				return nil
			}
			for _, n := range names {
				fmt.Fprintln(cmd.OutOrStdout(), n)
			}
			return nil
		},
	}
}

func newPinDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a pin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := envfile.DeletePin(keystorePath(), args[0]); err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("pin %q not found", args[0])
				}
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "deleted pin %q\n", args[0])
			return nil
		},
	}
}
