package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	base "github.com/KoyamaSohei/ugit/base"
	data "github.com/KoyamaSohei/ugit/data"
	diff "github.com/KoyamaSohei/ugit/diff"
	"github.com/spf13/cobra"
)

func initHandler(cmd *cobra.Command, args []string) {
	base.Init()
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
	oids2ref := map[string][]string{}
	names, refs, err := data.GetRefs("", true)
	if err != nil {
		panic(err)
	}
	for i, ref := range refs {
		oids := fmt.Sprintf("%x", ref.Value)
		oids2ref[oids] = append(oids2ref[oids], names[i])
	}
	oidset, err := base.GetCommitsAndParents([][]byte{oid})
	if err != nil {
		panic(err)
	}
	for _, o := range oidset {
		if err := base.PrintCommit(o, oids2ref[fmt.Sprintf("%x", o)]); err != nil {
			panic(err)
		}
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
	names, refs, err := data.GetRefs("", false)
	if err != nil {
		panic(err)
	}
	oidset := make([][]byte, 0)
	dot := "digraph commits {\n"
	for i, ref := range refs {
		if ref.Symblic {
			dot += fmt.Sprintf("\"%s\" [shape=note]\n", names[i])
			dot += fmt.Sprintf("\"%s\" -> \"%s\"", names[i], ref.Value)
			continue
		}
		dot += fmt.Sprintf("\"%s\" [shape=note]\n", names[i])
		dot += fmt.Sprintf("\"%s\" -> \"%x\"", names[i], ref.Value)
		oidset = append(oidset, ref.Value)
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
	if len(args) == 0 {
		c, err := base.GetBranchName()
		if err != nil {
			panic(err)
		}
		bs, err := base.GetBranchNames()
		if err != nil {
			panic(err)
		}
		for _, b := range bs {
			fmt.Printf("%t %s\n", b == c, b)
		}
		return
	}
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

func statusHandler(cmd *cobra.Command, args []string) {
	head, err := base.GetOid("@")
	if err != nil {
		panic(err)
	}
	b, err := base.GetBranchName()
	if err != nil {
		panic(err)
	}
	if len(b) > 0 {
		fmt.Printf("On branch %s\n", b)
	} else {
		fmt.Printf("HEAD detached at %x\n", head[:10])
	}
	fmt.Printf("Changes to be committed:\n\n")
	t, _, _, err := base.GetCommit(head)
	if err != nil {
		panic(err)
	}
	nt, err := base.WriteTree(".")
	if err != nil {
		panic(err)
	}
	out, err := diff.GetTreesDiff(t, nt)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", out)
}

func resetHandler(cmd *cobra.Command, args []string) {
	oid, err := base.GetOid(args[0])
	if err != nil {
		panic(err)
	}
	if err := data.UpdateRef("HEAD", data.RefValue{Symblic: false, Value: oid}, true); err != nil {
		panic(err)
	}
}

func showHandler(cmd *cobra.Command, args []string) {
	oid, err := base.GetOid(args[0])
	if err != nil {
		panic(err)
	}
	if err := base.PrintCommit(oid, nil); err != nil {
		panic(err)
	}
	ntoid, p, _, err := base.GetCommit(oid)
	if err != nil {
		panic(err)
	}
	if len(p) == 0 {
		return
	}
	ptoid, _, _, err := base.GetCommit(p)
	if err != nil {
		panic(err)
	}
	out, err := diff.GetTreesDiff(ptoid, ntoid)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", out)
}

func diffHandler(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		args = append(args, "@")
	}
	oid, err := base.GetOid(args[0])
	if err != nil {
		panic(err)
	}
	t, _, _, err := base.GetCommit(oid)
	if err != nil {
		panic(err)
	}
	nt, err := base.WriteTree(".")
	if err != nil {
		panic(err)
	}
	out, err := diff.GetTreesDiff(t, nt)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", out)
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "ugit",
		Short: "ugit is DIY Git",
		Long:  `golang version of https://www.leshenko.net/p/ugit/`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Create an empty ugit repository or reinitialize an existing one",
		Run:   initHandler,
		Args:  cobra.NoArgs,
	}
	hashCmd := &cobra.Command{
		Use:   "hash-object",
		Short: "Compute object ID and optionally creates a blob from a file",
		Run:   hashHandler,
		Args:  cobra.ExactArgs(1),
	}
	catCmd := &cobra.Command{
		Use:   "cat-file",
		Short: "Provide content or type and size information for repository objects",
		Run:   catHandler,
		Args:  cobra.ExactArgs(1),
	}
	writeCmd := &cobra.Command{
		Use:   "write-tree",
		Short: "Create a tree object from the current index",
		Run:   writeHandler,
		Args:  cobra.NoArgs,
	}
	readCmd := &cobra.Command{
		Use:   "read-tree",
		Short: "Reads tree information into the index",
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
		Short: "Show commit logs",
		Run:   logHandler,
		Args:  cobra.MaximumNArgs(1),
	}
	checkoutCmd := &cobra.Command{
		Use:   "checkout",
		Short: "Switch branches or restore working tree files",
		Run:   checkoutHandler,
		Args:  cobra.ExactArgs(1),
	}
	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "Create a tag object",
		Run:   tagHandler,
		Args:  cobra.RangeArgs(1, 2),
	}
	kCmd := &cobra.Command{
		Use:   "k",
		Short: "Visualize tool like gitk",
		Run:   kHandler,
		Args:  cobra.NoArgs,
	}
	branchCmd := &cobra.Command{
		Use:   "branch",
		Short: "List, create, or delete branches",
		Run:   branchHandler,
		Args:  cobra.MaximumNArgs(2),
	}
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show the working tree status",
		Run:   statusHandler,
		Args:  cobra.NoArgs,
	}
	resetCmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset current HEAD to the specified state",
		Run:   resetHandler,
		Args:  cobra.ExactArgs(1),
	}
	showCmd := &cobra.Command{
		Use:   "show",
		Short: "Show various types of objects",
		Run:   showHandler,
		Args:  cobra.ExactArgs(1),
	}
	diffCmd := &cobra.Command{
		Use:   "diff",
		Short: "Show changes between commits, commit and working tree, etc",
		Run:   diffHandler,
		Args:  cobra.MaximumNArgs(1),
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
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(diffCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
