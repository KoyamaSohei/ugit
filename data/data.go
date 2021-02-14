package data

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Type is data type
type Type byte

const (
	// None type
	None Type = iota
	// Blob type
	Blob
	// Tree type
	Tree
	// Commit type
	Commit
)

// GITDIR is git directory
const GITDIR = ".ugit"

// Init initialize .ugit
func Init() error {
	if err := os.MkdirAll(GITDIR, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(fmt.Sprintf("%s/objects", GITDIR), 0755); err != nil {
		return err
	}
	return nil
}

// HashObject gen hash from data and save data.
func HashObject(data []byte, dtype Type) ([]byte, error) {
	data = append([]byte{byte(dtype)}, data...)
	h := sha1.New()
	if _, err := h.Write(data); err != nil {
		return []byte{}, err
	}
	bs := h.Sum(nil)
	p := fmt.Sprintf("%s/objects/%x", GITDIR, bs)
	if err := ioutil.WriteFile(p, data, 0755); err != nil {
		return []byte{}, err
	}
	return bs, nil
}

// GetObject get file from hash
func GetObject(oid []byte, expected Type) ([]byte, error) {
	path := fmt.Sprintf("%s/objects/%x", GITDIR, oid)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}
	if t := Type(b[0]); expected != None && expected != t {
		return []byte{}, fmt.Errorf("data type is invalid")
	}
	b = b[1:]
	return b, nil
}

// GetType get data type
func GetType(oid []byte) (Type, error) {
	path := fmt.Sprintf("%s/objects/%x", GITDIR, oid)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return None, err
	}
	return Type(b[0]), nil
}

// UpdateRef update ref
func UpdateRef(name string, oid []byte) error {
	path := fmt.Sprintf("%s/%s", GITDIR, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, oid, 0644); err != nil {
		return err
	}
	return nil
}

// GetRef get ref
func GetRef(name string) ([]byte, error) {
	path := fmt.Sprintf("%s/%s", GITDIR, name)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}
