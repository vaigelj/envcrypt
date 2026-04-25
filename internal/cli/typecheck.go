package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/envcrypt/internal/envfile"
)

// NewTypecheckCmd returns the cobra command for type-checking an env file
// against a JSON rules file.
func NewTypecheckCmd() *cobra.Command {
	var rulesFile string
	var outputJSON bool

	cmd := &cobra.Command{
		Use:   "typecheck <env-file>",
		Short: "Validate env value types against a rules file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := envfile.ParseFile(args[0])
			if err != nil {
				return fmt.Errorf("reading env file: %w", err)
			}

			rules, err := loadRules(rulesFile)
			if err != nil {
				return fmt.Errorf("reading rules file: %w", err)
			}

			violations := envfile.TypeCheck(entries, rules)

			if outputJSON {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(violations)
			}

			if len(violations) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "all type checks passed")
				return nil
			}

			for _, v := range violations {
				fmt.Fprintf(cmd.OutOrStdout(), "FAIL  %s\n", v.Error())
			}
			return fmt.Errorf("%d type violation(s) found", len(violations))
		},
	}

	cmd.Flags().StringVarP(&rulesFile, "rules", "r", ".envcrypt-types.json", "path to JSON rules file")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "output violations as JSON")
	return cmd
}

func loadRules(path string) ([]envfile.TypeRule, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var rules []envfile.TypeRule
	if err := json.NewDecoder(f).Decode(&rules); err != nil {
		return nil, err
	}
	return rules, nil
}

func init() {
	rootCmd.AddCommand(NewTypecheckCmd())
}
