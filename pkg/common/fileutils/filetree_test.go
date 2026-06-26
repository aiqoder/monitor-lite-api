package fileutils

import (
	"fmt"
	"testing"
)

func TestFileTree(t *testing.T) {
	list := FileList("")
	for _, v := range list {
		fmt.Println(FileList(v.AllPath))
	}
}
