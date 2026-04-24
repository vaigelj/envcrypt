package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewPatchCmd returns the cobra command for the patch sub-command.
func NewPatchCmd() *cobra.Command {
	var (
		setPairs    []string
		deleteKeys  []string
		renameExprs []string
	)

	cmd := &cobra.Command{
		Use:   "patch <file>",
		Short: "Apply set/delete/rename operations to an env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			var instructions []envfile.PatchInstruction

			for _, pair := range setPairs {
				k, v, ok := strings.Cut(pair, "=")
				if !ok {
					return fmt.Errorf("--set: invalid format %q, expected KEY=VALUE", pair)
				}
				instructions = append(instructions, envfile.PatchInstruction{
					Op: envfile.PatchSet, Key: k, Value: v,
				})
			}

			for _, key := range deleteKeys {
				instructions = append(instructions, envfile.PatchInstruction{
					Op: envfile.PatchDelete, Key: strings.TrimSpace(key),
				})
			}

			for _, expr := range renameExprs {
				old, newk, ok := strings.Cut(expr, ":")
				if !ok {
					return fmt.Errorf("--rename: invalid format %q, expected OLD:NEW", expr)
				}
				instructions = append(instructions, envfile.PatchInstruction{
					Op: envfile.PatchRename, Key: old, NewKey: newk,
				})
			}

			if len(instructions) == 0 {
				return fmt.Errorf("no patch operations specified; use --set, --delete, or --rename")
			}

			if err := envfile.PatchFile(path, instructions); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "patched %s (%d operation(s))\n", path, len(instructions))
			return nil
		},
	}

	cmd.Flags().StringArrayVar(&setPairs, "set", nil, "Set KEY=VALUE (repeatable)")
	cmd.Flags().StringArrayVar(&deleteKeys, "delete", nil, "Delete KEY (repeatable)")
	cmd.Flags().StringArrayVar(&renameExprs, "rename", nil, "Rename OLD:NEW (repeatable)")
	return cmd
}

func init() {
	rootCmd.AddCommand(NewPatchCmd())
}
