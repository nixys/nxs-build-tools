package main

import (
	"fmt"
	"github.com/monochromegane/go-gitignore"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

type giSlice []gitignore.IgnoreMatcher

var excludes = []string{
	"/.git",
	".gitignore",
}

func fopsCopy(src, dst string) result {
	var gitignores giSlice
	var res result

	info, err := os.Stat(src)
	if err != nil {

		fmt.Printf("Copy error, can't exec stat for opject: %s (opject path: %s)\n", err, src)
		return false
	}

	gitignores, res = fopsLoadExcludes(gitignores, "")
	if res == false {

		return false
	}

	return fopsCopyLaunch(src, dst, info, gitignores)
}

func fopsMkdir(path string, perm os.FileMode) result {

	if err := os.MkdirAll(path, perm); err != nil {

		fmt.Printf("Can't create directory: %s (path: %s)\n", err, path)
		return false
	}

	return true
}

func fopsRemove(path string) result {

	p, err := filepath.Abs(path)
	if err != nil {

		fmt.Printf("Can't get absolute path: %s\n", err)
		return false
	}

	if err := os.RemoveAll(p); err != nil {

		fmt.Printf("Can't delete path: %s\n", err)
		return false
	}

	return true
}

func fopsLoadExcludes(gitignores giSlice, path string) (giSlice, result) {

	p, err := filepath.Abs(path)
	if err != nil {

		fmt.Printf("Can't get absolute path: %s\n", err)
		return nil, false
	}

	eString := strings.Join(excludes, "\n")

	r := strings.NewReader(eString)
	g := gitignore.NewGitIgnoreFromReader(p, r)

	return append(gitignores, g), true
}

/* Read .gitignore file in giving directory (if exists) */
func fopsReadGitignore(path string, gitignores giSlice) (giSlice, result) {

	p, err := filepath.Abs(path)
	if err != nil {

		fmt.Printf("Can't get absolute path: %s\n", err)
		return nil, false
	}

	giPath := p + "/" + gitignoreName

	info, err := os.Stat(giPath)
	if err == nil && info.Mode().IsRegular() {

		gi, err := gitignore.NewGitIgnore(giPath)
		if err != nil {

			fmt.Printf("Read gitignore file error: %s\n", err)
			return nil, false
		}

		gitignores = append(gitignores, gi)
	}

	return gitignores, true
}

func fopsMatchGitignore(path string, isDir bool, gitignores giSlice) (bool, result) {

	p, err := filepath.Abs(path)
	if err != nil {

		fmt.Printf("Can't get absolute path: %s\n", err)
		return false, false
	}

	/* Walk by nested gitignore files and find matches */
	for _, gi := range gitignores {

		if gi.Match(p, isDir) == true {

			/* Path matched on one of gitignore rules */

			return true, true
		}
	}

	return false, true
}

func fopsCopyLaunch(src, dst string, info os.FileInfo, gitignores giSlice) result {

	if info.IsDir() {

		return fopsCopyLaunchFileDir(src, dst, info, gitignores)
	}

	return fopsCopyLaunchFile(src, dst, info)
}

func fopsCopyLaunchFile(src, dst string, info os.FileInfo) result {

	f, err := os.Create(dst)
	if err != nil {

		fmt.Printf("Can't create desination file: %s\n", err)
		return false
	}
	defer f.Close()

	uid := int(info.Sys().(*syscall.Stat_t).Uid)
	gid := int(info.Sys().(*syscall.Stat_t).Gid)

	if err = os.Chmod(dst, info.Mode()); err != nil {

		fmt.Printf("Can't chmod to desination file: %s\n", err)
		return false
	}

	if err = os.Chown(dst, uid, gid); err != nil {

		fmt.Printf("Can't chown to desination file: %s\n", err)
		return false
	}

	s, err := os.Open(src)
	if err != nil {

		fmt.Printf("Can't open source file: %s\n", err)
		return false
	}
	defer s.Close()

	_, err = io.Copy(f, s)
	if err != nil {

		fmt.Printf("Can't copy fie: %s (src: %s, dst: %s)\n", err, src, dst)
		return false
	}

	return true
}

func fopsCopyLaunchFileDir(src, dst string, info os.FileInfo, gitignores giSlice) result {

	if err := os.MkdirAll(dst, info.Mode()); err != nil {

		fmt.Printf("Can't create destination directory: %s\n", err)
		return false
	}

	uid := int(info.Sys().(*syscall.Stat_t).Uid)
	gid := int(info.Sys().(*syscall.Stat_t).Gid)

	if err := os.Chmod(dst, info.Mode()); err != nil {

		fmt.Printf("Can't chmod to desination directory: %s\n", err)
		return false
	}

	if err := os.Chown(dst, uid, gid); err != nil {

		fmt.Printf("Can't chown to desination directory: %s\n", err)
		return false
	}

	gitignores, res := fopsReadGitignore(src, gitignores)
	if res == false {

		return false
	}

	infos, err := ioutil.ReadDir(src)
	if err != nil {

		fmt.Printf("Can't read source directory: %s\n", err)
		return false
	}

	for _, info := range infos {

		m, res := fopsMatchGitignore(filepath.Join(src, info.Name()), info.IsDir(), gitignores)
		if res == false {

			return false
		}

		if m {

			continue
		}

		if res := fopsCopyLaunch(
			filepath.Join(src, info.Name()),
			filepath.Join(dst, info.Name()),
			info,
			gitignores,
		); res == false {

			return false
		}
	}

	return true
}
