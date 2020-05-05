package main

import (
	"fmt"
	"github.com/pborman/getopt/v2"
	"os"
	"strings"
)

type optArgs struct {
	projectRoot string
	targetDir   string
	buildName   string
	cmd         string
	setting     string
	origFile    string
}

var argsCommands = []string{commandBuild, commandMakeOrig, commandPopulate, commandSettingGet}
var argsSettings = []string{settingGetProjectName,
	settingGetVersionMajor,
	settingGetVersionMinor,
	settingGetVersionPatch}

func argsParse() optArgs {
	var opts optArgs

	args := getopt.New()

	helpFlag := args.BoolLong(
		"help",
		'h',
		"Show help")

	versionFlag := args.BoolLong(
		"version",
		'v',
		"Show program version")

	projectRoot := args.StringLong(
		"project-root",
		'p',
		"",
		"Project root. The project root it is a directory that contains '"+pSettingsFile+"' file. If specified directory does not contain settings file then search in the parent directory will be continued")

	cmd := args.EnumLong(
		"command",
		'c',
		argsCommands,
		"",
		"Available commands:\n- "+strings.Join(argsCommands, "\n- "))

	targetDir := args.StringLong(
		"target-dir",
		't',
		"",
		"Target directory for build packages. By default will be used directory '"+defaultTargetDir+"' in project root directory.")

	buildName := args.StringLong(
		"build-name",
		'b',
		"",
		"Name of build to make package. All available builds specified in the project settings file '"+pSettingsFile+"'")

	setting := args.EnumLong(
		"setting",
		's',
		argsSettings,
		"",
		"Get project setting from settings file. Available settings:\n- "+strings.Join(argsSettings, "\n- "))

	origFile := args.StringLong(
		"orig-file",
		'o',
		"",
		"If specified the package will be created from a orig archive file instead of source code from project root directory. Available either '.tar.gz' or '.tar.xz' files.")

	args.Parse(os.Args)

	/* Show help */
	if *helpFlag == true {

		argsHelp(args)
		os.Exit(0)
	}

	/* Show version */
	if *versionFlag == true {

		argsVersion()
		os.Exit(0)
	}

	/* Project root */
	opts.projectRoot = *projectRoot

	/* Target dir */
	opts.targetDir = *targetDir

	/* Build name */
	opts.buildName = *buildName

	/* Command */
	if len(*cmd) == 0 {

		opts.cmd = commandBuild
	} else {

		opts.cmd = *cmd
	}

	/* Get option value */
	opts.setting = *setting

	/* Orig file */
	opts.origFile = *origFile

	return opts
}

func argsHelp(args *getopt.Set) {

	additionalDescription := `
	
Additional description

  Each launch of nxs-build-tools makes some actions in accordance with of specified command (option '--command'). Available commands:

    '` + commandBuild + `': It is a default command. Create either 'deb' or 'rpm' package from the project source code or orig archive file. Command usage:
        nxs-build-tools --command=` + commandBuild + ` --build-name=BUILD_NAME [--project-root=PROJECT_ROOT] [--target-dir=TARGET_DIR] [--orig-file=ORIG_FILE_PATH]

    '` + commandMakeOrig + `': Create a source code orig archives. Two files ('.tar.gz' and '.tar.xz') will be created as a result from execution of this command. Command usage:
        nxs-build-tools --command=` + commandMakeOrig + ` [--project-root=PROJECT_ROOT] [--target-dir=TARGET_DIR]

    '` + commandPopulate + `': Populate specified directory (project root) with the necessary files to allows the project to use nxs-build-tools for build packages. Command usage:
        nxs-build-tools --command=` + commandPopulate + ` [--project-root=PROJECT_ROOT]

    '` + commandSettingGet + `': Get project settings from the '` + pSettingsFile + `' file. Basically this command is used by CMake, but it also can by used to automate your build processes. Command usage:
        nxs-build-tools --command=` + commandSettingGet + ` --setting=SETTING_NAME [--project-root=PROJECT_ROOT]
`

	args.PrintUsage(os.Stdout)

	fmt.Println(additionalDescription)
}

func argsVersion() {

	fmt.Println(VERSION)
}
