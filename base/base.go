package base

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	data "github.com/KoyamaSohei/ugit/data"
)

// WriteTree write tree
func WriteTree(root string) ([]byte, error) {
	ents := make([]data.Entry, 0)
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		p := filepath.Join(root, f.Name())
		if isIgnored(p) {
			continue
		}
		if !f.IsDir() {
			dat, err := ioutil.ReadFile(p)
			if err != nil {
				return nil, err
			}
			h, err := data.HashObject(dat, data.Blob)
			if err != nil {
				return nil, err
			}
			ents = append(ents, data.Entry{Name: p, Oid: h})
		} else {
			h, err := WriteTree(p)
			if err != nil {
				return nil, err
			}
			ents = append(ents, data.Entry{Name: p, Oid: h})
		}
	}

	return data.HashTreeEntries(ents)
}

// ClearDirectory clear dir
func ClearDirectory(root string) error {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}

	for _, f := range files {
		p := filepath.Join(root, f.Name())
		if isIgnored(p) {
			continue
		}
		if !f.IsDir() {
			fmt.Printf("remove %s\n", p)
			if err := os.Remove(p); err != nil {
				return err
			}
		} else {
			if err := ClearDirectory(p); err != nil {
				return err
			}
		}
	}
	if err := os.Remove(root); err != nil {
		fmt.Printf("warn: not empty dir %s\n", root)
	}
	return nil
}

// ReadTree read tree
func ReadTree(oid []byte) error {
	if t, err := data.GetType(oid); err != nil || t != data.Tree {
		return fmt.Errorf("this object is not tree")
	}
	ents, err := data.GetTreeEntries(oid)
	if err != nil {
		return err
	}
	for _, e := range ents {
		t, err := data.GetType(e.Oid)
		if err != nil {
			return err
		}
		switch t {
		case data.Tree:
			if err := os.MkdirAll(e.Name, 0755); err != nil {
				return err
			}
			ReadTree(e.Oid)
		case data.Blob:
			b, err := data.GetObject(e.Oid, data.Blob)
			if err != nil {
				return err
			}
			if err := os.MkdirAll(filepath.Dir(e.Name), 0755); err != nil {
				return err
			}
			if err := ioutil.WriteFile(e.Name, b, 0644); err != nil {
				return err
			}
			fmt.Printf("%s: %x\n", e.Name, e.Oid)
		}
	}
	return nil
}

func isIgnored(path string) bool {
	return strings.Contains(path, ".git") || strings.Contains(path, ".ugit") || strings.Contains(path, "ugit")
}

// Commit commit
func Commit(mes string) error {
	dat, err := WriteTree(".")
	if err != nil {
		return err
	}
	dat = append(dat, []byte{0, 0}...)
	parent, _ := data.GetRef("HEAD", true)
	dat = append(dat, parent.Value...)
	dat = append(dat, []byte{0, 0}...)
	dat = append(dat, []byte(mes)...)
	h, err := data.HashObject(dat, data.Commit)
	if err != nil {
		return err
	}
	if err := data.UpdateRef("HEAD", data.RefValue{Symblic: false, Value: h}, true); err != nil {
		return err
	}
	return nil
}

// GetCommit get commit
func GetCommit(oid []byte) ([]byte, []byte, string, error) {
	b, err := data.GetObject(oid, data.Commit)
	if err != nil {
		return nil, nil, "", err
	}
	prop := bytes.Split(b, []byte{0, 0})
	if len(prop) != 3 {
		return nil, nil, "", fmt.Errorf("invalid commit")
	}
	return prop[0], prop[1], string(prop[2]), nil
}

// Checkout checkout
func Checkout(name string) error {
	oid, err := GetOid(name)
	if err != nil {
		return err
	}
	t, _, _, err := GetCommit(oid)
	if err != nil {
		return err
	}
	if err := ClearDirectory("."); err != nil {
		panic(err)
	}
	if err := ReadTree(t); err != nil {
		return err
	}
	head := data.RefValue{Symblic: false, Value: oid}
	if isBranch(name) {
		head = data.RefValue{Symblic: true, Value: []byte(fmt.Sprintf("refs/heads/%s", name))}
	}
	return data.UpdateRef("HEAD", head, false)
}

// CreateTag create tag
func CreateTag(name string, oid []byte) error {
	path := fmt.Sprintf("refs/tags/%s", name)
	if err := data.UpdateRef(path, data.RefValue{Symblic: false, Value: oid}, true); err != nil {
		return err
	}
	return nil
}

// GetOid get oid
func GetOid(oids string) ([]byte, error) {
	if oids == "@" {
		b, err := data.GetRef("HEAD", true)
		if err != nil {
			return nil, err
		}
		return b.Value, nil
	}
	prefixs := []string{
		"",
		"refs/",
		"refs/tags/",
		"refs/heads/",
	}
	for _, p := range prefixs {
		path := fmt.Sprintf("%s%s", p, oids)
		b, err := data.GetRef(path, true)
		if err != nil {
			continue
		}
		return b.Value, nil
	}
	if len(oids) != 40 {
		return nil, fmt.Errorf("Unknown name %s", oids)
	}
	b, err := hex.DecodeString(oids)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GetCommitsAndParents get commits and parents
func GetCommitsAndParents(oidset [][]byte) ([][]byte, error) {
	used := map[string]int{}
	resset := make([][]byte, 0)

	for len(oidset) > 0 {
		oid := oidset[0]
		oidset = oidset[1:]
		oids := fmt.Sprintf("%x", oid)
		if _, ok := used[oids]; ok {
			continue
		}
		used[oids] = 0
		resset = append(resset, oid)
		_, p, _, err := GetCommit(oid)
		if err != nil {
			return nil, err
		}
		if len(p) > 0 {
			oidset = append(oidset, p)
		}
	}

	return resset, nil
}

// CreateBranch create branch
func CreateBranch(name string, oid []byte) error {
	path := fmt.Sprintf("refs/heads/%s", name)
	return data.UpdateRef(path, data.RefValue{Symblic: false, Value: oid}, true)
}

func isBranch(branch string) bool {
	path := fmt.Sprintf("refs/heads/%s", branch)
	_, err := data.GetRef(path, true)
	if err != nil {
		return false
	}
	return true
}
