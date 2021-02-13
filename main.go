package main

import (
	"fmt"
	"io/ioutil"
	"os"

	base "github.com/KoyamaSohei/ugit/base"
	data "github.com/KoyamaSohei/ugit/data"
	"github.com/spf13/cobra"
)

func initHandler(cmd *cobra.Command, args []string) {
	fmt.Printf("Hello,World\n")
	data.Init()
	pwd, _ := os.Getwd()
	fmt.Printf("Initialized empty ugit repository in %s/%s\n", pwd, data.GITDIR)
}

func hashHandler(cmd *cobra.Command, args []string) {
	dat, err := ioutil.ReadFile(args[0])
	if err != nil {
		panic(err)
	}
	h := data.HashObject(dat, data.Blob)
	fmt.Printf("%x", h)
}

func catHandler(cmd *cobra.Command, args []string) {
	b := data.GetObject(args[0], data.None)
	fmt.Printf("%s", string(b))
}

func writeHandler(cmd *cobra.Command, args []string) {
	h := base.WriteTree(".")
	fmt.Printf("%x\n", h)
}

func readHandler(cmd *cobra.Command, args []string) {
	base.ClearDirectory(".")
	base.ReadTree(args[0])
}

func commitHandler(cmd *cobra.Command, args []string) {
	base.Commit(args[0])
}

func logHandler(cmd *cobra.Command, args []string) {
	oid := ""
	if len(args) == 1 {
		if data.GetType(args[0]) != data.Commit {
			panic(fmt.Errorf("hash type is %d,not Commit", data.GetType(args[0])))
		}
		oid = args[0]
	} else {
		oid = fmt.Sprintf("%x", data.GetRef("HEAD"))
	}
	for {
		t, p, m := base.GetCommit(oid)
		fmt.Printf("commit: %s\ntree: %x\nmessage: %s\n\n", oid, t, m)
		if len(p) == 0 {
			break
		}
		oid = fmt.Sprintf("%x", p)
	}
}

func checkoutHandler(cmd *cobra.Command, args []string) {
	if data.GetType(args[0]) != data.Commit {
		panic(fmt.Errorf("hash type is %d,not Commit", data.GetType(args[0])))
	}
	base.ClearDirectory(".")
	base.Checkout(args[0])
}

func tagHandler(cmd *cobra.Command, args []string) {
	base.CreateTag(args[0], args[1])
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "ugit",
		Short: "ugit is DIY Git",
		Long:  `golang ver of https://www.leshenko.net/p/ugit/`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "ugit init",
		Run:   initHandler,
		Args:  cobra.NoArgs,
	}
	hashCmd := &cobra.Command{
		Use:   "hash-object",
		Short: "save object",
		Run:   hashHandler,
		Args:  cobra.ExactArgs(1),
	}
	catCmd := &cobra.Command{
		Use:   "cat-file",
		Short: "ugit cat",
		Run:   catHandler,
		Args:  cobra.ExactArgs(1),
	}
	writeCmd := &cobra.Command{
		Use:   "write-tree",
		Short: "ugit write",
		Run:   writeHandler,
		Args:  cobra.NoArgs,
	}
	readCmd := &cobra.Command{
		Use:   "read-tree",
		Short: "read tree",
		Run:   readHandler,
		Args:  cobra.ExactArgs(1),
	}
	commitCmd := &cobra.Command{
		Use:   "commit",
		Short: "commit [commit message]",
		Run:   commitHandler,
		Args:  cobra.ExactArgs(1),
	}
	logCmd := &cobra.Command{
		Use:   "log",
		Short: "log",
		Run:   logHandler,
		Args:  cobra.MaximumNArgs(1),
	}
	checkoutCmd := &cobra.Command{
		Use:   "checkout",
		Short: "checkout",
		Run:   checkoutHandler,
		Args:  cobra.ExactArgs(1),
	}
	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "tag",
		Run:   tagHandler,
		Args:  cobra.ExactArgs(2),
	}

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(hashCmd)
	rootCmd.AddCommand(catCmd)
	rootCmd.AddCommand(writeCmd)
	rootCmd.AddCommand(readCmd)
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(checkoutCmd)
	rootCmd.AddCommand(tagCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
