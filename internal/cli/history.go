package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envcrypt/internal/envfile"
)

// NewHistoryCmd returns the history subcommand.
func NewHistoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history",
		Short: "Manage env file history",
	}
	cmd.AddCommand(newHistoryListCmd())
	cmd.AddCommand(newHistoryRecordCmd())
	cmd.AddCommand(newHistoryClearCmd())
	return cmd
}

func newHistoryListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list <env-file>",
		Short: "List recorded history entries for an env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			hf, err := envfile.LoadHistory(args[0])
			if err != nil {
				return fmt.Errorf("no history found: %w", err)
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "#\tTIMESTAMP\tLABEL\tKEYS")
			for i, e := range hf.Entries {
				fmt.Fprintf(w, "%d\t%s\t%s\t%d\n",
					i+1,
					e.Timestamp.Format("2006-01-02 15:04:05"),
					e.Label,
					len(e.Values),
				)
			}
			return w.Flush()
		},
	}
}

func newHistoryRecordCmd() *cobra.Command {
	var label string
	cmd := &cobra.Command{
		Use:   "record <env-file>",
		Short: "Record current state of an env file into history",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			vals, err := envfile.ParseFile(args[0])
			if err != nil {
				return err
			}
			if label == "" {
				label = "manual"
			}
			if err := envfile.AppendHistory(args[0], label, vals); err != nil {
				return err
			}
			fmt.Printf("Recorded history entry %q for %s\n", label, args[0])
			return nil
		},
	}
	cmd.Flags().StringVarP(&label, "label", "l", "", "Label for this history entry")
	return cmd
}

func newHistoryClearCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clear <env-file>",
		Short: "Clear all history for an env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := envfile.ClearHistory(args[0]); err != nil {
				return err
			}
			fmt.Printf("History cleared for %s\n", args[0])
			return nil
		},
	}
}
