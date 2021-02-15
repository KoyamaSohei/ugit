package data

import (
	"bytes"
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

// Entry is dir's content
type Entry struct {
	Oid  []byte
	Name string
}

// RefValue is ref container
type RefValue struct {
	Symblic bool
	Value   []byte
}

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
func UpdateRef(name string, ref RefValue) error {
	if ref.Symblic {
		return fmt.Errorf("cannnot update Symblic ref")
	}
	path := fmt.Sprintf("%s/%s", GITDIR, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, ref.Value, 0644); err != nil {
		return err
	}
	return nil
}

// GetRef get ref
func GetRef(name string) (RefValue, error) {
	path := fmt.Sprintf("%s/%s", GITDIR, name)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return RefValue{}, err
	}
	r := []byte("ref:")
	if bytes.HasPrefix(b, r) {
		s := bytes.Split(b, r)
		if len(s) != 2 {
			return RefValue{}, fmt.Errorf("invalid format")
		}
		return GetRef(fmt.Sprintf("%s", s[1]))
	}
	return RefValue{false, b}, nil
}

// GetTreeEntries get entries
func GetTreeEntries(oid []byte) ([]Entry, error) {
	h, err := GetObject(oid, Tree)
	if err != nil {
		return nil, err
	}
	ents := make([]Entry, 0)
	o := make([]byte, 0)
	for k, b := range bytes.Split(h, []byte{0, 0}) {
		if k%2 == 0 {
			o = b
			continue
		}
		ents = append(ents, Entry{Oid: o, Name: string(b)})
	}
	return ents, nil
}

// HashTreeEntries set entries
func HashTreeEntries(ents []Entry) ([]byte, error) {
	conts := make([]byte, 0)
	for _, ent := range ents {
		c := ent.Oid
		c = append(c, []byte{0, 0}...)
		c = append(c, []byte(ent.Name)...)
		c = append(c, []byte{0, 0}...)
		conts = append(conts, c...)
	}
	return HashObject(conts, Tree)
}

// GetRefs get refs
func GetRefs() ([]string, error) {
	refs := []string{"HEAD"}
	err := filepath.Walk(fmt.Sprintf("%s/refs", GITDIR), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			kind := filepath.Base(filepath.Dir(path))
			refs = append(refs, fmt.Sprintf("%s/%s", kind, info.Name()))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return refs, nil
}
