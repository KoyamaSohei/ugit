package main

import (
	"fmt"
	"io/ioutil"
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
	h, err := hash.Output()
	assert.Nil(t, err)
	fmt.Printf("hash: %s\n", string(h))
	cat := exec.Command("./ugit", "cat-file", string(h))
	o, err := cat.Output()
	assert.Nil(t, err)
	m, err := ioutil.ReadFile("main.go")
	assert.Nil(t, err)
	assert.Equal(t, o, m)
}
