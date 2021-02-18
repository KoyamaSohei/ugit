package main

import (
	"os"
	"os/exec"
	"testing"

	"gotest.tools/v3/assert"
)

func setup() {
	init := exec.Command("./ugit", "init")
	if err := init.Run(); err != nil {
		panic(err)
	}
}

func teardown() {
	if err := os.RemoveAll(".ugit"); err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	teardown()
	os.Exit(ret)
}

func TestCommit(t *testing.T) {
	commit := exec.Command("./ugit", "commit", "hello-ugit")
	err := commit.Run()
	assert.NilError(t, err)
}

func TestCommitAgain(t *testing.T) {
	commit := exec.Command("./ugit", "commit", "hello-ugit")
	commit2 := exec.Command("./ugit", "commit", "hello-ugit")
	err := commit.Run()
	assert.NilError(t, err)
	err = commit2.Run()
	assert.NilError(t, err)
}

func TestLog(t *testing.T) {
	commit := exec.Command("./ugit", "commit", "hello-ugit")
	log := exec.Command("./ugit", "log")
	err := commit.Run()
	assert.NilError(t, err)
	err = log.Run()
	assert.NilError(t, err)
}

func TestCheckout(t *testing.T) {
	commit := exec.Command("./ugit", "commit", "hello-ugit")
	branch := exec.Command("./ugit", "branch", "main")
	checkout := exec.Command("./ugit", "checkout", "main")
	err := commit.Run()
	assert.NilError(t, err)
	err = branch.Run()
	assert.NilError(t, err)
	err = checkout.Run()
	assert.NilError(t, err)
}
