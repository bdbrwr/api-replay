package cliutils

import (
	"fmt"

	"github.com/spf13/cobra"
)

func GetArgOrFlag(cmd *cobra.Command, args []string, flagName string, argIndex int, what string) (string, error) {
	val, _ := cmd.Flags().GetString(flagName)
	if val == "" {
		if len(args) <= argIndex {
			return "", fmt.Errorf("missing %s: use --%s or pass as positional argument (index %d)", what, flagName, argIndex)
		}
		val = args[argIndex]
	}
	return val, nil
}
