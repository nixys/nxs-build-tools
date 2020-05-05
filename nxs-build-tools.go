package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	populateTpls     = "/usr/share/nxs-build-tools/template"
	settingsFile     = ".proj-settings.yml"
	defaultTargetDir = "builds"
)

const (
	commandBuild      = "build"
	commandMakeOrig   = "make-orig"
	commandPopulate   = "populate"
	commandSettingGet = "setting-get"
)

const (
	settingGetProjectName  = "project-name"
	settingGetVersionMajor = "version-major"
	settingGetVersionMinor = "version-minor"
	settingGetVersionPatch = "version-patch"
)

func main() {

	opts := argsParse()

	if opts.cmd == commandPopulate {
		projectRoot, err := filepath.Abs(opts.projectRoot)
		if err != nil {
			fmt.Printf("Can't get absolute path: %s \n", err)
			os.Exit(1)
		}

		if err := populateProject(projectRoot); err != nil {
			fmt.Printf("Can't populate project: %s\n", err)
			os.Exit(1)
		}

		os.Exit(0)
	}

	ctx, err := contextInit(opts)
	if err != nil {
		fmt.Printf("Context init error: %s\n", err)
		os.Exit(1)
	}

	switch opts.cmd {
	case commandBuild:

		/* Check build name specified */
		if len(opts.buildName) == 0 {
			fmt.Printf("Build name must be specified\n")
			os.Exit(1)
		}

		if err := buildPackage(ctx); err != nil {
			fmt.Printf("Build package error: %s\n", err)
			os.Exit(1)
		}

	case commandMakeOrig:

		err = buildMakeOrig(ctx)
		if err != nil {
			fmt.Printf("Make orig error error: %s\n", err)
			os.Exit(1)
		}

	case commandSettingGet:

		/* Check setting specified */
		if len(opts.setting) == 0 {
			fmt.Printf("No setting specified\n")
			os.Exit(1)
		}

		settingVal, err := settingsGet(ctx, opts.setting)
		if err != nil {
			fmt.Printf("Setting get error: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("%s", settingVal)

	default:

		fmt.Printf("Unknown command: %s\n", opts.cmd)
		os.Exit(1)
	}

	os.Exit(0)
}

func settingsGet(ctx selfContext, setting string) (string, error) {

	switch setting {
	case settingGetProjectName:
		return fmt.Sprintf("%s", ctx.conf.ProjectName), nil
	case settingGetVersionMajor:
		return fmt.Sprintf("%d", ctx.conf.Version.Major), nil
	case settingGetVersionMinor:
		return fmt.Sprintf("%d", ctx.conf.Version.Minor), nil
	case settingGetVersionPatch:
		return fmt.Sprintf("%d", ctx.conf.Version.Patch), nil
	}

	return "", fmt.Errorf("unknown setting: %s", setting)
}
