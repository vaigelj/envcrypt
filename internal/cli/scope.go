package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewScopeCmd returns the root "scope" command with subcommands.
func NewScopeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scope",
		Short: "Manage key scopes for filtered env views",
	}
	cmd.AddCommand(newScopeSetCmd())
	cmd.AddCommand(newScopeListCmd())
	cmd.AddCommand(newScopeShowCmd())
	cmd.AddCommand(newScopeRemoveCmd())
	return cmd
}

func newScopeSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <name> <KEY1,KEY2,...>",
		Short: "Create or update a scope with a comma-separated list of keys",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			keys := strings.Split(args[1], ",")
			for i, k := range keys {
				keys[i] = strings.TrimSpace(k)
			}
			if err := envfile.AddScope(".", args[0], keys); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "scope %q saved with %d key(s)\n", args[0], len(keys))
			return nil
		},
	}
}

func newScopeListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all defined scopes",
		RunE: func(cmd *cobra.Command, _ []string) error {
			scopes, err := envfile.LoadScopes(".")
			if err != nil {
				return err
			}
			if len(scopes) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no scopes defined")
				return nil
			}
			for _, s := range scopes {
				fmt.Fprintf(cmd.OutOrStdout(), "%s (%d keys)\n", s.Name, len(s.Keys))
			}
			return nil
		},
	}
}

func newScopeShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Show keys belonging to a scope",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			scopes, err := envfile.LoadScopes(".")
			if err != nil {
				return err
			}
			for _, s := range scopes {
				if s.Name == args[0] {
					fmt.Fprintln(cmd.OutOrStdout(), strings.Join(s.Keys, "\n"))
					return nil
				}
			}
			return fmt.Errorf("scope %q not found", args[0])
		},
	}
}

func newScopeRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name>",
		Short: "Delete a scope",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := envfile.RemoveScope(".", args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "scope %q removed\n", args[0])
			return nil
		},
	}
}
