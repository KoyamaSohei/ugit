package data

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
)

// HashObject gen hash from data and save data.
func HashObject(data []byte, dtype dataType) []byte {
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
