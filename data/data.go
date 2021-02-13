package data

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
)

// DataType is data type
type DataType byte

const (
	// None type
	None DataType = iota
	// Blob type
	Blob
	// Tree type
	Tree
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

// HashObject gen hash from data and save data.
func HashObject(data []byte, dtype DataType) []byte {
	data = append([]byte{byte(dtype)}, data...)
	h := sha1.New()
	if _, err := h.Write(data); err != nil {
		panic(err)
	}
	bs := h.Sum(nil)
	p := fmt.Sprintf("%s/objects/%x", GITDIR, bs)
	if err := ioutil.WriteFile(p, data, 0755); err != nil {
		panic(err)
	}
	return bs
}

// GetObject get file from hash
func GetObject(oid string, expected DataType) []byte {
	path := fmt.Sprintf("%s/objects/%s", GITDIR, oid)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if t := DataType(b[0]); expected != None && expected != t {
		panic(fmt.Errorf("data type is invalid"))
	}
	b = b[1:]
	return b
}
