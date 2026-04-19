package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewProfileCmd returns the profile management command.
func NewProfileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage environment profiles (.env.<name>)",
	}
	cmd.AddCommand(newProfileListCmd(), newProfileShowCmd())
	return cmd
}

func newProfileListCmd() *cobra.Command {
	var dir string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available profiles in a directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			profiles, err := envfile.ListProfiles(dir)
			if err != nil {
				return err
			}
			if len(profiles) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no profiles found")
				return nil
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "PROFILE\tPATH")
			for _, p := range profiles {
				fmt.Fprintf(w, "%s\t%s\n", p.Name, p.Path)
			}
			return w.Flush()
		},
	}
	cmd.Flags().StringVarP(&dir, "dir", "d", ".", "directory to scan for profiles")
	return cmd
}

func newProfileShowCmd() *cobra.Command {
	var dir string
	var redact bool
	cmd := &cobra.Command{
		Use:   "show <profile>",
		Short: "Show key/value pairs for a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := envfile.LoadProfile(dir, args[0])
			if err != nil {
				return err
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			for k, v := range m {
				if redact && envfile.IsSensitive(k) {
					v = "********"
				}
				fmt.Fprintf(w, "%s\t%s\n", k, v)
			}
			if err := w.Flush(); err != nil {
				return err
			}
			os.Stdout.Sync() //nolint
			return nil
		},
	}
	cmd.Flags().StringVarP(&dir, "dir", "d", ".", "directory containing the profile")
	cmd.Flags().BoolVar(&redact, "redact", false, "mask sensitive values")
	return cmd
}
