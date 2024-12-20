package main

import (
	"fmt"
	"testing"
)

func TestEmbedFile(t *testing.T) {
	dir, err := fs.ReadDir("data")
	check(err)

	for _, dentry := range dir {
		fmt.Println(dentry.Name())
	}
}
