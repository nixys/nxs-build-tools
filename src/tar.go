package main

import (
	"fmt"
	"github.com/mholt/archiver"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	TAR_TYPE_UNKNOWN int = 0
	TAR_TYPE_GZ      int = 1
	TAR_TYPE_XZ      int = 2
)

func tarMakeGz(dstName, src string) result {
	var files []string

	p, err := filepath.Abs(src)
	if err != nil {

		fmt.Printf("Can't get absolute path: %s\n", err)
		return false
	}

	files = tarDirScan(p, files)
	if files == nil {

		return false
	}

	if err := archiver.TarGz.Make(dstName, files); err != nil {

		fmt.Printf("Can't create tar.gz archive: %s\n", err)
		return false
	}

	return true
}

func tarMakeXz(dstName, src string) result {
	var files []string

	p, err := filepath.Abs(src)
	if err != nil {

		fmt.Printf("Can't get absolute path: %s\n", err)
		return false
	}

	files = tarDirScan(p, files)
	if files == nil {

		return false
	}

	if err = archiver.TarXZ.Make(dstName, files); err != nil {

		fmt.Printf("Can't create tar.xz archive: %s\n", err)
		return false
	}

	return true
}

func tarOpenGz(archFile, dst string) result {

	p, err := filepath.Abs(dst)
	if err != nil {

		fmt.Printf("Can't get absolute path: %s\n", err)
		return false
	}

	if err = archiver.TarGz.Open(archFile, p); err != nil {

		fmt.Printf("Can't unpack tar.gz archive: %s\n", err)
		return false
	}

	return true
}

func tarOpenXz(archFile, dst string) result {

	p, err := filepath.Abs(dst)
	if err != nil {

		fmt.Printf("Can't get absolute path: %s\n", err)
		return false
	}

	if err = archiver.TarXZ.Open(archFile, p); err != nil {

		fmt.Printf("Can't unpack tar.xz archive: %s\n", err)
		return false
	}

	return true
}

func tarGetArchType(archFile string) int {

	if archiver.TarGz.Match(archFile) == true {

		return TAR_TYPE_GZ
	}

	if archiver.TarXZ.Match(archFile) == true {

		return TAR_TYPE_XZ
	}

	return TAR_TYPE_UNKNOWN
}

func tarDirScan(path string, files []string) []string {

	info, err := os.Stat(path)
	if err != nil {

		fmt.Printf("Stat error: %s (path: %s)\n", err, path)
		return nil
	}

	if info.IsDir() {

		infos, err := ioutil.ReadDir(path)
		if err != nil {

			fmt.Printf("Can't read source directory for tar: %s\n", err)
			return nil
		}

		for _, info := range infos {

			files = append(files, filepath.Join(path, info.Name()))
		}
	} else {

		files = append(files, path)
	}

	return files
}
