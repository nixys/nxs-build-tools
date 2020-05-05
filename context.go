package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type context struct {
	projectRoot string    /* Root directory for current project */
	targetDir   string    /* Path to directory where build packages for project */
	buildName   string    /* Name of the build to create package for */
	origFile    string    /* Path to orig archive */
	pSettings   PSettings /* Project settings (from project settings file) */
}

func contextInit(opts optArgs) (context, result) {
	var ctx context
	var res result

	/* Build path */

	ctx.projectRoot, res = contextInitProjectRoot(opts.projectRoot)
	if res == false {

		return ctx, false
	}

	ctx.targetDir, res = contextInitTargetDir(ctx.projectRoot, opts.targetDir)
	if res == false {

		return ctx, false
	}

	ctx.buildName = opts.buildName
	ctx.origFile = opts.origFile

	ctx.pSettings, res = psettingsLoad(ctx.projectRoot)
	if res == false {

		return ctx, false
	}

	return ctx, true
}

func contextInitProjectRoot(optsProjectRoot string) (string, result) {
	var projectRoot string
	var res result
	var err error

	if len(optsProjectRoot) == 0 {

		projectRoot, res = contextProjectRootGet(optsProjectRoot)
		if res == false {

			return projectRoot, false
		}
	} else {

		projectRoot, err = filepath.Abs(optsProjectRoot)
		if err != nil {

			fmt.Printf("Can't get absolute path: %s\n", err)
			return projectRoot, false
		}
	}

	return projectRoot, true
}

func contextInitTargetDir(ctxProjectRoot, optsTargetDir string) (string, result) {
	var targetDir string
	var err error

	if len(optsTargetDir) == 0 {

		targetDir = ctxProjectRoot + "/" + defaultTargetDir
	} else {

		targetDir, err = filepath.Abs(optsTargetDir)
		if err != nil {

			fmt.Printf("Can't get absolute path: %s\n", err)
			return targetDir, false
		}
	}

	return targetDir, true
}

/*
 Get project root

 The project root is a directory that contains .nxs-settings.proj file.
 Looking for in parent directory if the current directory doesn't contain .nxs-settings.proj file.
*/
func contextProjectRootGet(path string) (string, result) {

	p, err := filepath.Abs(path)
	if err != nil {

		fmt.Printf("Can't get absolute path: %s\n", err)
		return "", false
	}

	pSPath := p + "/" + pSettingsFile

	info, err := os.Stat(pSPath)
	if err == nil && info.Mode().IsRegular() {

		/* If project settings file found */

		return p, true
	}

	if p == "/" {

		fmt.Printf("Can't find project settings file %s\n", pSettingsFile)
		return "", false
	}

	return contextProjectRootGet(filepath.Dir(p))
}
