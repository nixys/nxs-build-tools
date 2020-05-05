package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type selfContext struct {

	// Project root directory
	projectRoot string

	// Directory path for project build
	targetDir string

	// Name of the build to create package for
	buildName string

	// Path to orig archive file
	origFile string

	// Project settings
	conf confOpts
}

func contextInit(args argsOpts) (selfContext, error) {

	var (
		ctx selfContext
		err error
	)

	ctx.projectRoot, err = contextInitProjectRoot(args.projectRoot)
	if err != nil {
		return ctx, err
	}

	ctx.targetDir, err = contextInitTargetDir(ctx.projectRoot, args.targetDir)
	if err != nil {
		return ctx, err
	}

	ctx.buildName = args.buildName
	ctx.origFile = args.origFile

	ctx.conf, err = confRead(ctx.projectRoot + "/" + settingsFile)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

// contextInitProjectRoot gets project root directory
func contextInitProjectRoot(projectRoot string) (string, error) {

	var (
		pr  string
		err error
	)

	if len(projectRoot) == 0 {
		pr, err = contextProjectRootLookup(projectRoot)
		if err != nil {
			return pr, err
		}
	} else {
		pr, err = filepath.Abs(projectRoot)
		if err != nil {
			return pr, fmt.Errorf("can't get absolute path: %s (path: %s)", err, projectRoot)
		}
	}

	return pr, nil
}

func contextInitTargetDir(projectRoot, targetDir string) (string, error) {

	var (
		td  string
		err error
	)

	if len(targetDir) == 0 {
		td = projectRoot + "/" + defaultTargetDir
	} else {
		td, err = filepath.Abs(targetDir)
		if err != nil {
			return td, fmt.Errorf("can't get absolute path: %s (path: %s)", err, targetDir)
		}
	}

	return td, nil
}

// contextProjectRootLookup lookups the directory from `path` till `/` contains settings file.
func contextProjectRootLookup(path string) (string, error) {

	p, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("can't get absolute path: %s (path: %s)", err, path)
	}

	pSPath := p + "/" + settingsFile

	info, err := os.Stat(pSPath)
	if err == nil && info.Mode().IsRegular() {
		// If project settings file found
		return p, nil
	}

	if p == "/" {
		return "", fmt.Errorf("can't find project settings file %s", settingsFile)
	}

	return contextProjectRootLookup(filepath.Dir(p))
}
