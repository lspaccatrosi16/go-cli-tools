package zip

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mandelsoft/vfs/pkg/vfs"
)

func crawlAndAdd(base string, zipBase string, fs vfs.FileSystem, zipWriter *zip.Writer) error {
	files, err := vfs.ReadDir(fs, base)
	if err != nil {
		return err
	}
	for _, file := range files {
		// fmt.Printf("file: %s from %s (ref %s)\n", file.Name(), base, zipBase)
		if file.IsDir() {
			newBase := filepath.Join(base, file.Name())
			newZipBase := filepath.Join(zipBase, file.Name())
			err = crawlAndAdd(newBase, newZipBase, fs, zipWriter)
			if err != nil {
				return err
			}
		} else {
			srcPath := filepath.Join(base, file.Name())
			src, err := fs.Open(srcPath)
			if err != nil {
				return err
			}
			defer src.Close()
			dstPath := filepath.Join(zipBase, file.Name())
			dst, err := zipWriter.Create(dstPath)
			if err != nil {
				return err
			}
			io.Copy(dst, src)
		}
	}
	return nil
}

func ZipFolder(startDir string, fs vfs.FileSystem) (*[]byte, error) {
	if exists, err := vfs.DirExists(fs, startDir); err != nil || !exists {
		return nil, fmt.Errorf("path %s does not exist or is not a directory", startDir)
	}
	output := bytes.NewBuffer([]byte{})
	zipWriter := zip.NewWriter(output)
	err := crawlAndAdd(startDir, "", fs, zipWriter)
	if err != nil {
		return nil, err
	}
	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}
	outBytes := output.Bytes()
	return &outBytes, nil
}

func UnzipFolder(data *[]byte, outpath string, fs vfs.FileSystem) (string, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(*data), int64(len(*data)))
	if err != nil {
		return "", err
	}

	baseFolderName := ""
	for _, zipFile := range zipReader.File {
		if baseFolderName == "" {
			baseFolderName = zipFile.Name
		}
		inFsName := filepath.Join(outpath, zipFile.Name)
		// fmt.Printf("Reading %s as %s\n", zipFile.Name, inFsName)
		if zipFile.FileInfo().IsDir() {
			// fmt.Printf("is dir %s\n", inFsName)
			err = fs.MkdirAll(inFsName, os.ModeDir)
			if err != nil {
				return "", err
			}
			continue
		}
		exists, err := vfs.Exists(fs, filepath.Dir(inFsName))
		if err != nil {
			return "", err
		}
		if !exists {
			fs.MkdirAll(filepath.Dir(inFsName), os.ModeDir)
		}
		// fmt.Printf("processing file at path %s\n", inFsName)
		dstBuf := bytes.NewBuffer(nil)
		if err != nil {
			return "", err
		}
		archivedFile, err := zipFile.Open()
		if err != nil {
			return "", err
		}
		defer archivedFile.Close()
		_, err = io.Copy(dstBuf, archivedFile)
		if err != nil {
			return "", err
		}
		err = vfs.WriteFile(fs, inFsName, dstBuf.Bytes(), 0o755)
		if err != nil {
			return "", err
		}
	}
	return baseFolderName, nil
}
