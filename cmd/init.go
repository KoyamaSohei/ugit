package cmd

import (
	"fmt"
	"os"

	data "github.com/KoyamaSohei/ugit/data"

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
	data.Init()
	pwd, _ := os.Getwd()
	fmt.Printf("Initialized empty ugit repository in %s/%s\n", pwd, data.GITDIR)
}
