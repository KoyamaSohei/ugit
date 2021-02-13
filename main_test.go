package main

import (
	"os"
	"os/exec"
	"testing"

	data "github.com/KoyamaSohei/ugit/data"
	"github.com/stretchr/testify/assert"
)

func cleanUp() {
	os.RemoveAll(data.GITDIR)
}

func Test_main(t *testing.T) {
	t.Cleanup(cleanUp)
	init := exec.Command("./ugit", "init")
	assert.Nil(t, init.Run())
	hash := exec.Command("./ugit", "hash-object", "main.go")
	assert.Nil(t, hash.Run())
}
