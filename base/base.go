package base

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	data "github.com/KoyamaSohei/ugit/data"
)

type entry struct {
	oid  []byte
	name string
}

// WriteTree write tree
func WriteTree(root string) []byte {
	ents := make([]entry, 0)
	files, err := ioutil.ReadDir(root)

	if err != nil {
		panic(err)
	}

	for _, f := range files {
		p := filepath.Join(root, f.Name())
		if isIgnored(p) {
			continue
		}
		if !f.IsDir() {
			dat, err := ioutil.ReadFile(p)
			if err != nil {
				panic(err)
			}
			h := data.HashObject(dat, data.Blob)
			ents = append(ents, entry{name: p, oid: h})
		} else {
			h := WriteTree(p)
			ents = append(ents, entry{name: p, oid: h})
		}
	}

	conts := make([]byte, 0)
	for _, ent := range ents {
		c := ent.oid
		c = append(c, []byte{0, 0}...)
		c = append(c, []byte(ent.name)...)
		c = append(c, []byte{0, 0}...)
		conts = append(conts, c...)
	}

	return data.HashObject(conts, data.Tree)
}

func iterTreeEntries(oid string) []entry {
	h := data.GetObject(oid, data.Tree)
	ents := make([]entry, 0)
	o := make([]byte, 0)
	for k, b := range bytes.Split(h, []byte{0, 0}) {
		if k%2 == 0 {
			o = b
			continue
		}
		ents = append(ents, entry{oid: o, name: string(b)})
	}
	return ents
}

func emptyCurrentDirectory() {
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == "." {
			return filepath.SkipDir
		}
		if isIgnored(path) && info.IsDir() {
			return filepath.SkipDir
		}
		if isIgnored(path) {
			return nil
		}
		fmt.Printf("remove %s\n", path)
		if info.IsDir() {
			if err := os.Remove(path); err != nil {
				fmt.Printf("warn: %s is not empty\n", path)
			}
			return nil
		}
		return os.Remove(path)
	})
}

// ReadTree read tree
func ReadTree(oid string) {
	emptyCurrentDirectory()
	if data.GetDataType(oid) != data.Tree {
		panic(fmt.Errorf("this object is not tree"))
	}
	for _, e := range iterTreeEntries(oid) {
		oids := fmt.Sprintf("%x", e.oid)
		switch data.GetDataType(oids) {
		case data.Tree:
			if err := os.MkdirAll(e.name, 0755); err != nil {
				panic(err)
			}
			ReadTree(oids)
		case data.Blob:
			b := data.GetObject(oids, data.Blob)
			if err := os.MkdirAll(filepath.Dir(e.name), 0755); err != nil {
				panic(err)
			}
			if err := ioutil.WriteFile(e.name, b, 0755); err != nil {
				panic(err)
			}
			fmt.Printf("%s: %x\n", e.name, e.oid)
		}
	}
}

func isIgnored(path string) bool {
	return strings.Contains(path, ".git") || strings.Contains(path, ".ugit") || strings.Contains(path, "ugit")
}
