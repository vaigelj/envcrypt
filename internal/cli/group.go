package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewGroupCmd returns the root group command.
func NewGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group",
		Short: "Manage named groups of env var keys",
	}
	cmd.AddCommand(newGroupSetCmd(), newGroupListCmd(), newGroupRemoveCmd(), newGroupShowCmd())
	return cmd
}

func newGroupSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <name> <KEY1,KEY2,...>",
		Short: "Create or update a group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			keys := strings.Split(args[1], ",")
			if err := envfile.AddGroup(".", args[0], keys); err != nil {
				return err
			}
			fmt.Printf("Group %q saved with %d key(s).\n", args[0], len(keys))
			return nil
		},
	}
}

func newGroupListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			groups, err := envfile.LoadGroups(".")
			if err != nil {
				return err
			}
			if len(groups) == 0 {
				fmt.Println("No groups defined.")
				return nil
			}
			for _, g := range groups {
				fmt.Printf("%-20s %s\n", g.Name, strings.Join(g.Keys, ", "))
			}
			return nil
		},
	}
}

func newGroupShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Show keys in a group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			g, err := envfile.GetGroup(".", args[0])
			if err != nil {
				return err
			}
			for _, k := range g.Keys {
				fmt.Println(k)
			}
			return nil
		},
	}
}

func newGroupRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := envfile.RemoveGroup(".", args[0]); err != nil {
				return err
			}
			fmt.Printf("Group %q removed.\n", args[0])
			return nil
		},
	}
}
