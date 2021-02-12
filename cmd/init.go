package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// InitCmd is ugit init command.
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "ugit init",
	Run:   initHandler,
}

func initHandler(cmd *cobra.Command, args []string) {
	fmt.Printf("Hello,World\n")
}
