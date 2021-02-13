package data

import (
	"fmt"
	"io/ioutil"
)

// GetObject get file from hash
func GetObject(oid string, expected dataType) []byte {
	path := fmt.Sprintf("%s/objects/%s", GITDIR, oid)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if t := dataType(b[0]); expected != None && expected != t {
		panic(fmt.Errorf("data type is invalid"))
	}
	b = b[1:]
	return b
}
