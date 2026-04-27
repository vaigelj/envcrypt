package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewAnnotateCmd returns the root annotate command with subcommands.
func NewAnnotateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "annotate",
		Short: "Manage annotations (description, owner, tags) for env keys",
	}
	cmd.AddCommand(newAnnotateSetCmd())
	cmd.AddCommand(newAnnotateShowCmd())
	cmd.AddCommand(newAnnotateRemoveCmd())
	return cmd
}

func newAnnotateSetCmd() *cobra.Command {
	var desc, owner string
	var deprecated bool
	var tags []string

	cmd := &cobra.Command{
		Use:   "set <key>",
		Short: "Set annotation for an env key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := os.Getwd()
			ann := envfile.Annotation{
				Description: desc,
				Owner:       owner,
				Deprecated:  deprecated,
				Tags:        tags,
			}
			if err := envfile.SetAnnotation(dir, args[0], ann); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "annotation set for %s\n", args[0])
			return nil
		},
	}
	cmd.Flags().StringVar(&desc, "desc", "", "Description of the key")
	cmd.Flags().StringVar(&owner, "owner", "", "Owning team or person")
	cmd.Flags().BoolVar(&deprecated, "deprecated", false, "Mark key as deprecated")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "Comma-separated tags")
	return cmd
}

func newAnnotateShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <key>",
		Short: "Show annotation for an env key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := os.Getwd()
			ann, ok, err := envfile.GetAnnotation(dir, args[0])
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("no annotation found for %s", args[0])
			}
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "key:         %s\n", args[0])
			fmt.Fprintf(out, "description: %s\n", ann.Description)
			fmt.Fprintf(out, "owner:       %s\n", ann.Owner)
			fmt.Fprintf(out, "deprecated:  %v\n", ann.Deprecated)
			if len(ann.Tags) > 0 {
				fmt.Fprintf(out, "tags:        %s\n", strings.Join(ann.Tags, ", "))
			}
			return nil
		},
	}
}

func newAnnotateRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <key>",
		Short: "Remove annotation for an env key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := os.Getwd()
			if err := envfile.RemoveAnnotation(dir, args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "annotation removed for %s\n", args[0])
			return nil
		},
	}
}
