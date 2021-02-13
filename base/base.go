package base

import (
	"io/ioutil"
	"path/filepath"

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
		if f.Name() == ".git" || f.Name() == ".ugit" {
			continue
		}
		if !f.IsDir() {
			dat, err := ioutil.ReadFile(filepath.Join(root, f.Name()))
			if err != nil {
				panic(err)
			}
			h := data.HashObject(dat, data.Blob)
			ents = append(ents, entry{name: f.Name(), oid: h})
		} else {
			h := WriteTree(filepath.Join(root, f.Name()))
			ents = append(ents, entry{name: f.Name(), oid: h})
		}
	}

	conts := make([]byte, 0)
	for _, ent := range ents {
		c := ent.oid
		c = append(c, 0)
		c = append(c, []byte(ent.name)...)
		c = append(c, 0)
		conts = append(conts, c...)
	}

	return data.HashObject(conts, data.Tree)
}
