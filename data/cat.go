package data

import (
	"fmt"
	"io/ioutil"
)

// GetObject get file from hash
func GetObject(oid string) []byte {
	path := fmt.Sprintf("%s/objects/%s", GITDIR, oid)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return b
}
