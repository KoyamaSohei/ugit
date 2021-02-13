package base

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteTree write tree
func WriteTree(dir string) {
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !info.IsDir() {
			fmt.Println(path)
		}
		return nil
	})
}
