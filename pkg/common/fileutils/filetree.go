package fileutils

import (
	"github.com/shirou/gopsutil/disk"
	"os"
	"path/filepath"
	"runtime"
)

type File struct {
	AllPath string
	Path    string
	IsDir   bool
}

// root 表示文件根目录
func FileList(root string) []File {
	var ret []File

	if root == "" && runtime.GOOS == "windows" {
		partitions, err := disk.Partitions(true)
		if err != nil {
			return ret
		}
		for _, partition := range partitions {
			dir := filepath.Join(partition.Device, "/")
			ret = append(ret, File{AllPath: dir, Path: dir, IsDir: true})
		}
		return ret
	}

	if root == "" {
		root = "/"
	}

	dir, err := os.ReadDir(root)

	if err != nil {
		return ret
	}

	for _, d := range dir {
		ret = append(ret, File{
			AllPath: filepath.Join(root, d.Name()),
			Path:    d.Name(),
			IsDir:   d.IsDir(),
		})
	}
	return ret
}
