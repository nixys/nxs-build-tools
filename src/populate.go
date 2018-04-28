package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func populateProject(projectRoot string) result {
	var files []string

	files, res := populateWalkTemplates(populateTpls, "", files)
	if res == false {

		return false
	}

	for _, s := range files {

		p := fmt.Sprintf("%s/%s", projectRoot, s)

		if _, err := os.Stat(p); err == nil {

			fmt.Printf("Can't populate your project, file `%s` already exist\n", p)
			return false
		}
	}

	if res := fopsCopy(populateTpls, projectRoot); res == false {

		return false
	}

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

	return true
}

func populateWalkTemplates(path, relativePath string, files []string) ([]string, result) {
	var res result

	info, err := os.Stat(path)
	if err != nil {

		return nil, false
	}

	if info.IsDir() {

		infos, err := ioutil.ReadDir(path)
		if err != nil {

			fmt.Printf("Can't read template directory: %s\n", err)
			return nil, false
		}

		if len(relativePath) > 0 {

			relativePath += "/"
		}

		for _, i := range infos {

			files, res = populateWalkTemplates(filepath.Join(path, i.Name()), relativePath+i.Name(), files)
			if res == false {

				return nil, false
			}
		}
	} else {

		files = append(files, relativePath)
	}

	return files, true
}
