package arch

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mholt/archiver"
)

// Make makes an archive from `src` and saves result to file with `dst` name
func Make(src, dst string) error {

	p, err := filepath.Abs(src)
	if err != nil {
		return fmt.Errorf("can't get absolute path: %s (file: %s)", err, src)
	}

	files, err := prepare(p)
	if err != nil {
		return err
	}

	if err = archiver.Archive(files, dst); err != nil {
		return fmt.Errorf("can't create archive: %s (archive: %s)", err, dst)
	}

	return nil
}

// Open opens archive `arch` to `dst`
func Open(arch, dst string) error {

	p, err := filepath.Abs(dst)
	if err != nil {
		return fmt.Errorf("can't get absolute path: %s (file: %s)", err, dst)
	}

	if err = archiver.Unarchive(arch, p); err != nil {
		return fmt.Errorf("can't unpack archive: %s (file: %s)", err, arch)
	}

	return nil
}

func prepare(path string) ([]string, error) {

	var files []string

	info, err := os.Stat(path)
	if err != nil {
		return []string{}, fmt.Errorf("can't prepare file to archive, stat error: %s (path: %s)", err, path)
	}

	if info.IsDir() {

		infos, err := ioutil.ReadDir(path)
		if err != nil {
			return []string{}, fmt.Errorf("can't prepare file to archive, read source directory error: %s (path: %s)", err, path)
		}

		for _, info := range infos {
			files = append(files, filepath.Join(path, info.Name()))
		}
	} else {
		files = append(files, path)
	}

	return files, nil
}
