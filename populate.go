package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/nixys/nxs-build-tools/fops"
)

func populateProject(projectRoot string) error {

	files, err := populateWalkTemplates(populateTpls, "")
	if err != nil {
		return err
	}

	// Check no files from templates directory
	// exists in `projectRoot` dir
	for _, s := range files {

		p := fmt.Sprintf("%s/%s", projectRoot, s)

		if _, err := os.Stat(p); err == nil {
			return fmt.Errorf("file `%s` already exist", p)
		}
	}

	// Copy template files into project root directory
	if err := fops.CopyWithIgnores(populateTpls, projectRoot); err != nil {
		return err
	}

	// Helper message
	fmt.Println(
		`Project has been successfully populated!

Make sure your .gitignore file contains necessary excludes.
If not you may just use next command:

cat <<EOF >> ` + projectRoot + `/.gitignore
/builds
/objs
_CPack_Packages
CMakeCache.txt
CMakeFiles/
Makefile
CPackConfig.cmake
CPackSourceConfig.cmake
cmake_install.cmake
EOF`)

	return nil
}

// populateWalkTemplates returns all files contains in templates directory
func populateWalkTemplates(path, relativePath string) ([]string, error) {

	var files []string

	info, err := os.Stat(path)
	if err != nil {
		return []string{}, err
	}

	if info.IsDir() {

		infos, err := ioutil.ReadDir(path)
		if err != nil {
			return []string{}, fmt.Errorf("read template directory error: %s (path: %s)", err, path)
		}

		if len(relativePath) > 0 {
			relativePath += "/"
		}

		for _, i := range infos {

			f, err := populateWalkTemplates(filepath.Join(path, i.Name()), relativePath+i.Name())
			if err != nil {
				return []string{}, err
			}
			files = append(files, f...)
		}
	} else {

		files = append(files, relativePath)
	}

	return files, nil
}
