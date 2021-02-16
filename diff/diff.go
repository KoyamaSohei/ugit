package diff

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	data "github.com/KoyamaSohei/ugit/data"
)

func getBlobsDiff(poid, noid []byte, name string) (string, error) {
	po, err := data.GetObject(poid, data.Blob)
	if err != nil {
		return "", err
	}
	no, err := data.GetObject(noid, data.Blob)
	if err != nil {
		return "", err
	}
	pname := fmt.Sprintf("%s.prev", filepath.Base(name))
	ptmp, err := ioutil.TempFile(filepath.Dir(name), pname)
	if err != nil {
		return "", err
	}
	defer os.Remove(ptmp.Name())
	nname := fmt.Sprintf("%s.next", filepath.Base(name))
	ntmp, err := ioutil.TempFile(filepath.Dir(name), nname)
	if err != nil {
		return "", err
	}
	defer os.Remove(ntmp.Name())
	if _, err := ptmp.Write(po); err != nil {
		return "", err
	}
	if _, err := ntmp.Write(no); err != nil {
		return "", err
	}
	diff := exec.Command("diff", "--unified", "--show-c-function", "--label", pname, ptmp.Name(), "--label", nname, ntmp.Name())
	out, _ := diff.Output()
	return string(out), nil
}

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
		if pt == data.Blob && nt == data.Blob {
			cout, err := getBlobsDiff(po, e.Oid, e.Name)
			if err != nil {
				return "", err
			}
			out += cout
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
