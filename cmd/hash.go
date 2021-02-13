package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/KoyamaSohei/ugit/data"
	"github.com/spf13/cobra"
)

var hashCmd = &cobra.Command{
	Use:   "hash-object",
	Short: "save object",
	Run:   hashHandler,
	Args:  cobra.ExactArgs(1),
}

func hashHandler(cmd *cobra.Command, args []string) {
	dat, err := ioutil.ReadFile(args[0])
	if err != nil {
		panic(err)
	}
	h := data.HashObject(dat, data.Blob)
	fmt.Printf("%x", h)
}
