package fops

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/monochromegane/go-gitignore"
)

const gitignoreName = ".gitignore"

var builtInIgnores = []string{
	"/.git",
	".gitignore",
}

// CopyWithIgnores copies `src` to `dst` considering `.gitignore` files in each subdirectory
func CopyWithIgnores(src, dst string) error {

	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("stat for opject error: %s (path: %s)", err, src)
	}

	gitignores := loadBuiltInIgnores()

	return copyExec(src, dst, info, gitignores)
}

// MkdirRecursive makes dir `path` with permissions `perm`.
// Also if needed parent directories will be created.
func MkdirRecursive(path string, perm os.FileMode) error {
	if err := os.MkdirAll(path, perm); err != nil {
		return fmt.Errorf("directory create error: %s (path: %s)", err, path)
	}
	return nil
}

// Remove removes object `path` from filesystem.
// Directories will be removed recursively.
func Remove(path string) error {

	p, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("get absolute path error: %s (path: %s)", err, path)
	}

	if err := os.RemoveAll(p); err != nil {
		return fmt.Errorf("delete path error: %s (path: %s)", err, p)
	}

	return nil
}

// Normalize normalizes path `path`.
// If `path` begins with `~/` it would been replaced
// by real home dir for user.
// If `path` is a directory, would be add `/` at the end of path.
func Normalize(path string) (string, error) {

	if strings.HasPrefix(path, "~/") == true {

		usr, err := user.Current()
		if err != nil {
			return "", err
		}

		path = usr.HomeDir + "/" + strings.TrimPrefix(path, "~/")
	}

	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if info.IsDir() {
		for strings.HasSuffix(path, "/") == true {
			path = strings.TrimSuffix(path, "/")
		}
		path += "/"
	}

	return path, nil
}

func loadBuiltInIgnores() []gitignore.IgnoreMatcher {

	var gitignores []gitignore.IgnoreMatcher

	p, _ := filepath.Abs("")

	eString := strings.Join(builtInIgnores, "\n")

	r := strings.NewReader(eString)
	g := gitignore.NewGitIgnoreFromReader(p, r)

	return append(gitignores, g)
}

// readGitignore reads .gitignore file in giving directory (if exists)
func readGitignore(path string) ([]gitignore.IgnoreMatcher, error) {

	var gitignores []gitignore.IgnoreMatcher

	p, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("get absolute path error: %s (path: %s)", err, path)
	}

	giPath := p + "/" + gitignoreName

	info, err := os.Stat(giPath)
	if err == nil && info.Mode().IsRegular() {
		// If .gitignore file exist
		gi, err := gitignore.NewGitIgnore(giPath)
		if err != nil {
			return nil, fmt.Errorf("read gitignore file error: %s (path: %s)", err, giPath)
		}
		gitignores = append(gitignores, gi)
	}

	return gitignores, nil
}

func matchGitignore(path string, isDir bool, gitignores []gitignore.IgnoreMatcher) (bool, error) {

	p, err := filepath.Abs(path)
	if err != nil {
		return false, fmt.Errorf("get absolute path error: %s (path: %s)", err, path)
	}

	/* Walk by nested gitignore files and find matches */
	for _, gi := range gitignores {
		if gi.Match(p, isDir) == true {
			/* Path matched on one of gitignore rules */
			return true, nil
		}
	}

	return false, nil
}

func copyExec(src, dst string, info os.FileInfo, gitignores []gitignore.IgnoreMatcher) error {
	if info.IsDir() {
		return copyDir(src, dst, info, gitignores)
	}
	return copyFile(src, dst, info)
}

func copyFile(src, dst string, info os.FileInfo) error {

	f, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create desination file error: %s (path: %s)", err, dst)
	}
	defer f.Close()

	if err = os.Chmod(dst, info.Mode()); err != nil {
		return fmt.Errorf("chmod to desination file error: %s (path: %s)", err, dst)
	}

	s, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source file error: %s (path: %s)", err, src)
	}
	defer s.Close()

	_, err = io.Copy(f, s)
	if err != nil {
		return fmt.Errorf("copy file error: %s (src: %s, dst: %s)", err, src, dst)
	}

	return nil
}

func copyDir(src, dst string, info os.FileInfo, gitignores []gitignore.IgnoreMatcher) error {

	if err := os.MkdirAll(dst, info.Mode()); err != nil {
		return fmt.Errorf("create destination directory error: %s (dir: %s)", err, dst)
	}

	if err := os.Chmod(dst, info.Mode()); err != nil {
		return fmt.Errorf("chmod to desination directory error: %s (dir: %s)", err, dst)
	}

	g, err := readGitignore(src)
	if err != nil {
		return err
	}
	gitignores = append(gitignores, g...)

	infos, err := ioutil.ReadDir(src)
	if err != nil {
		return fmt.Errorf("read source directory error: %s (path: %s)", err, src)
	}

	for _, i := range infos {

		m, err := matchGitignore(filepath.Join(src, i.Name()), i.IsDir(), gitignores)
		if err != nil {
			return err
		}

		if m == true {
			continue
		}

		if err := copyExec(
			filepath.Join(src, i.Name()),
			filepath.Join(dst, i.Name()),
			i,
			gitignores,
		); err != nil {
			return err
		}
	}

	return nil
}
