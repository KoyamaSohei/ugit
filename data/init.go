package data

import (
	"fmt"
	"os"
)

// GITDIR is git directory
const GITDIR = ".ugit"

// Init initialize .ugit
func Init() {
	if err := os.MkdirAll(GITDIR, 0755); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(fmt.Sprintf("%s/objects", GITDIR), 0755); err != nil {
		panic(err)
	}
}
