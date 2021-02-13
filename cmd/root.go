package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ugit",
	Short: "ugit is DIY Git",
	Long:  `golang ver of https://www.leshenko.net/p/ugit/`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(hashCmd)
}

// Execute executes command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
