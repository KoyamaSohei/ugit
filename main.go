package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

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
	h, err := data.HashObject(dat, data.Blob)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x", h)
}

func catHandler(cmd *cobra.Command, args []string) {
	oid, err := base.GetOid(args[0])
	if err != nil {
		panic(err)
	}
	b, err := data.GetObject(oid, data.None)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", string(b))
}

func writeHandler(cmd *cobra.Command, args []string) {
	h, err := base.WriteTree(".")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", h)
}

func readHandler(cmd *cobra.Command, args []string) {
	oid, err := base.GetOid(args[0])
	if err != nil {
		panic(err)
	}
	base.ClearDirectory(".")
	base.ReadTree(oid)
}

func commitHandler(cmd *cobra.Command, args []string) {
	base.Commit(args[0])
}

func logHandler(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		args = append(args, "@")
	}
	oid, err := base.GetOid(args[0])
	if err != nil {
		panic(err)
	}
	t, err := data.GetType(oid)
	if err != nil {
		panic(err)
	}
	if t != data.Commit {
		panic(fmt.Errorf("hash type is %d,not Commit", t))
	}
	oidset, err := base.GetCommitsAndParents([][]byte{oid})
	if err != nil {
		panic(err)
	}
	for _, o := range oidset {
		t, _, m, err := base.GetCommit(o)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Commit  %x\ntree    %x\nmessage %s\n\n", o, t, m)
	}
}

func checkoutHandler(cmd *cobra.Command, args []string) {
	if err := base.Checkout(args[0]); err != nil {
		panic(err)
	}
}

func tagHandler(cmd *cobra.Command, args []string) {
	if len(args) == 1 {
		args = append(args, "@")
	}
	oid, err := base.GetOid(args[1])
	if err != nil {
		panic(err)
	}
	base.CreateTag(args[0], oid)
}

func kHandler(cmd *cobra.Command, args []string) {
	refs, err := data.GetRefs()
	if err != nil {
		panic(err)
	}
	oidset := make([][]byte, 0)
	dot := "digraph commits {\n"
	for _, ref := range refs {
		r, err := data.GetRef(ref, false)
		fmt.Printf("ref %s Symblic %t %x\n", ref, r.Symblic, r.Value)
		if err != nil {
			panic(err)
		}
		if r.Symblic {
			dot += fmt.Sprintf("\"%s\" [shape=note]\n", ref)
			dot += fmt.Sprintf("\"%s\" -> \"%s\"", ref, r.Value)
			continue
		}
		dot += fmt.Sprintf("\"%s\" [shape=note]\n", ref)
		dot += fmt.Sprintf("\"%s\" -> \"%x\"", ref, r.Value)
		oidset = append(oidset, r.Value)
	}

	if oidset, err = base.GetCommitsAndParents(oidset); err != nil {
		panic(err)
	}

	for _, oid := range oidset {
		_, p, _, err := base.GetCommit(oid)
		if err != nil {
			panic(err)
		}
		dot += fmt.Sprintf("\"%x\" [shape=box style=filled label=\"%x\"]\n", oid, oid[:10])
		if len(p) > 0 {
			dot += fmt.Sprintf("\"%x\" -> \"%x\"\n", oid, p)
		}
	}

	dot += "}\n"

	viz := exec.Command("dot", "-Tgtk", "/dev/stdin")
	wc, err := viz.StdinPipe()
	if err != nil {
		panic(err)
	}

	viz.Start()
	if _, err := wc.Write([]byte(dot)); err != nil {
		panic(err)
	}
	if err := wc.Close(); err != nil {
		panic(err)
	}
	if err := viz.Wait(); err != nil {
		panic(err)
	}
}

func branchHandler(cmd *cobra.Command, args []string) {
	if len(args) == 1 {
		args = append(args, "@")
	}
	oid, err := base.GetOid(args[1])
	if err != nil {
		panic(err)
	}
	if err := base.CreateBranch(args[0], oid); err != nil {
		panic(err)
	}
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
		Args:  cobra.RangeArgs(1, 2),
	}
	kCmd := &cobra.Command{
		Use:   "k",
		Short: "k",
		Run:   kHandler,
		Args:  cobra.NoArgs,
	}
	branchCmd := &cobra.Command{
		Use:   "branch",
		Short: "branch",
		Run:   branchHandler,
		Args:  cobra.RangeArgs(1, 2),
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
	rootCmd.AddCommand(kCmd)
	rootCmd.AddCommand(branchCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
