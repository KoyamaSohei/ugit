package diff

import (
	"bytes"
	"fmt"

	"github.com/KoyamaSohei/ugit/data"
)

// GetTreesDiff return tree's diff
func GetTreesDiff(ptoid, ntoid []byte) (string, error) {
	out := ""
	pent, err := data.GetTreeEntries(ptoid)
	if err != nil {
		return "", err
	}
	nent, err := data.GetTreeEntries(ntoid)
	if err != nil {
		return "", err
	}
	pmap := map[string][]byte{}
	used := map[string]bool{}
	for _, e := range pent {
		pmap[e.Name] = e.Oid
		used[e.Name] = false
	}
	for _, e := range nent {
		po, ok := pmap[e.Name]
		if !ok {
			out += fmt.Sprintf("new file %s\n", e.Name)
			continue
		}
		used[e.Name] = true
		if bytes.Equal(e.Oid, po) {
			continue
		}
		pt, err := data.GetType(po)
		if err != nil {
			return "", err
		}
		nt, err := data.GetType(e.Oid)
		if err != nil {
			return "", err
		}
		if pt != data.Tree || nt != data.Tree {
			out += fmt.Sprintf("mod file %s\n", e.Name)
			continue
		}
		cout, err := GetTreesDiff(po, e.Oid)
		if err != nil {
			return "", err
		}
		out += cout
	}
	for name, u := range used {
		if !u {
			out += fmt.Sprintf("del file %s\n", name)
		}
	}
	return out, nil
}
