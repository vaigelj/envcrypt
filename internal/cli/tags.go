package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envcrypt/internal/envfile"
)

// NewTagsCmd returns the root tags command.
func NewTagsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "Manage key tags for grouping env variables",
	}
	cmd.AddCommand(newTagsSetCmd(), newTagsListCmd(), newTagsRemoveCmd())
	return cmd
}

func newTagsSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <name> <KEY1,KEY2,...>",
		Short: "Create or replace a tag",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			ts, err := envfile.LoadTags(dir)
			if err != nil {
				return err
			}
			keys := strings.Split(args[1], ",")
			ts.AddTag(args[0], keys)
			if err := envfile.SaveTags(dir, ts); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "tag %q set with %d key(s)\n", args[0], len(keys))
			return nil
		},
	}
}

func newTagsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all tags",
		RunE: func(cmd *cobra.Command, args []string) error {
			ts, err := envfile.LoadTags(".")
			if err != nil {
				return err
			}
			if len(ts.Tags) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no tags defined")
				return nil
			}
			for _, t := range ts.Tags {
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", t.Name, strings.Join(t.Keys, ", "))
			}
			return nil
		},
	}
}

func newTagsRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a tag by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			ts, err := envfile.LoadTags(dir)
			if err != nil {
				return err
			}
			if !ts.RemoveTag(args[0]) {
				return fmt.Errorf("tag %q not found", args[0])
			}
			if err := envfile.SaveTags(dir, ts); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "tag %q removed\n", args[0])
			return nil
		},
	}
}
