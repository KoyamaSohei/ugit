package cmd

import (
	"fmt"

	"github.com/KoyamaSohei/ugit/data"
	"github.com/spf13/cobra"
)

var catCmd = &cobra.Command{
	Use:   "cat-file",
	Short: "ugit cat",
	Run:   catHandler,
	Args:  cobra.ExactArgs(1),
}

func catHandler(cmd *cobra.Command, args []string) {
	b := data.GetObject(args[0])
	fmt.Printf("%s", string(b))
}
