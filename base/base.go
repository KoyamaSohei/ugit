package base

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/KoyamaSohei/ugit/data"
)

// WriteTree write tree
func WriteTree() {
	_ = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if path == ".git" || path == ".ugit" {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			dat, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			h := data.HashObject(dat, data.Blob)
			fmt.Printf("%x\n", h)
		}
		return nil
	})
}
