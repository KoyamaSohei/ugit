package data

import "os"

// GITDIR is git directory
const GITDIR = ".ugit"

// Init initialize .ugit
func Init() {
	os.Mkdir(GITDIR, os.FileMode(os.O_RDWR))
}
