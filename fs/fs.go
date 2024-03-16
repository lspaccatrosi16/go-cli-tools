package fs

import (
	"io"
	"os"
	"path/filepath"

	"github.com/mandelsoft/vfs/pkg/vfs"
)

func VfsToDisk(diskRoot string, fsRoot string, fs vfs.FileSystem) error {
	// fmt.Printf("translate: from %s to %s\n", fsRoot, diskRoot)
	files, err := vfs.ReadDir(fs, fsRoot)
	if err != nil {
		return err
	}
	for _, file := range files {
		// fmt.Printf("translate %s from %s to %s\n", file.Name(), fsRoot, diskRoot)
		newDiskFile := filepath.Join(diskRoot, file.Name())
		newFsFile := filepath.Join(fsRoot, file.Name())
		if file.IsDir() {
			err = os.MkdirAll(newDiskFile, 0o755)
			if err != nil {
				return err
			}
			err = VfsToDisk(newDiskFile, newFsFile, fs)
			if err != nil {
				return err
			}
		} else {
			src, err := fs.Open(newFsFile)
			if err != nil {
				return err
			}
			defer src.Close()
			dst, err := os.Create(newDiskFile)
			if err != nil {
				return err
			}
			defer dst.Close()
			err = dst.Chmod(0o755)
			if err != nil {
				return err
			}
			_, err = io.Copy(dst, src)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
