package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	populateTpls           string = "/usr/share/nxs-build-tools/template"
	gitignoreName          string = ".gitignore"
	pSettingsFile          string = ".proj-settings.yml"
	defaultTargetDir       string = "builds"
	commandBuild           string = "build"
	commandMakeOrig        string = "make-orig"
	commandPopulate        string = "populate"
	commandSettingGet      string = "setting-get"
	settingGetProjectName  string = "project-name"
	settingGetVersionMajor string = "version-major"
	settingGetVersionMinor string = "version-minor"
	settingGetVersionPatch string = "version-patch"
)

type result bool

func main() {

	opts := argsParse()

	switch opts.cmd {

	case commandBuild:

		/* Check build name specified */
		if len(opts.buildName) == 0 {

			fmt.Printf("Build name must be specified\n")
			os.Exit(1)
		}

		ctx, res := contextInit(opts)
		if res == false {

			os.Exit(1)
		}

		res = buildPackage(ctx)
		if res == false {

			os.Exit(1)
		}

	case commandMakeOrig:

		ctx, res := contextInit(opts)
		if res == false {

			os.Exit(1)
		}

		res = buildMakeOrig(ctx)
		if res == false {

			os.Exit(1)
		}

	case commandSettingGet:

		/* Check setting specified */
		if len(opts.setting) == 0 {

			fmt.Printf("No setting specified\n")
			os.Exit(1)
		}

		ctx, res := contextInit(opts)
		if res == false {

			os.Exit(1)
		}

		settingVal, res := psettingsGet(ctx, opts.setting)
		if res == false {

			os.Exit(1)
		}

		fmt.Printf("%s", settingVal)

	case commandPopulate:

		projectRoot, err := filepath.Abs(opts.projectRoot)
		if err != nil {

			fmt.Printf("Can't get absolute path: %s\n", err)
			os.Exit(1)
		}

		if res := populateProject(projectRoot); res == false {

			os.Exit(1)
		}

	default:

		fmt.Printf("Unknown command: %s\n", opts.cmd)
		os.Exit(1)
	}

	os.Exit(0)
}
