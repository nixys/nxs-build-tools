package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

func buildMakeOrig(ctx context) result {

	sourceDir := buildSourceDirPath(
		ctx.targetDir,
		"orig",
		ctx.pSettings.ProjectName,
		ctx.pSettings.Version.Major,
		ctx.pSettings.Version.Minor,
		ctx.pSettings.Version.Patch)
	origNameTgz := buildOrigFilePath(
		ctx.targetDir,
		ctx.pSettings.ProjectName,
		ctx.pSettings.Version.Major,
		ctx.pSettings.Version.Minor,
		ctx.pSettings.Version.Patch,
		"tar.gz")
	origNameTxz := buildOrigFilePath(
		ctx.targetDir,
		ctx.pSettings.ProjectName,
		ctx.pSettings.Version.Major,
		ctx.pSettings.Version.Minor,
		ctx.pSettings.Version.Patch,
		"tar.xz")

	if res := fopsCopy(ctx.projectRoot, sourceDir); res == false {

		return false
	}

	/* Make tar.gz */
	if res := tarMakeGz(origNameTgz, sourceDir); res == false {

		return false
	}

	/* Make tar.xz */
	if res := tarMakeXz(origNameTxz, sourceDir); res == false {

		return false
	}

	if res := fopsRemove(sourceDir); res == false {

		return false
	}

	return true
}

func buildPackage(ctx context) result {

	for _, b := range ctx.pSettings.Builds {

		if b.Name == ctx.buildName {

			if b.Deb != nil {

				return buildPackageDeb(ctx, b)
			}

			if b.Rpm != nil {

				return buildPackageRpm(ctx, b)
			}
		}
	}

	fmt.Printf("No such build: \"%s\"\n", ctx.buildName)
	return false
}

func buildPackageDeb(ctx context, build PSettingsBuild) result {

	sourceDir, res := buildPrepSourceDir(ctx, build)
	if res == false {

		return false
	}

	cmdDhmake := exec.Command("dh_make", build.Deb.DHMake...)
	cmdBuildpkg := exec.Command("dpkg-buildpackage", build.Deb.DpkgBuildpackage...)

	cmdDhmake.Dir = sourceDir
	cmdDhmake.Env = append(os.Environ(), "PWD="+sourceDir)
	cmdDhmake.Stdout = os.Stdout
	cmdDhmake.Stderr = os.Stderr

	cmdBuildpkg.Dir = sourceDir
	cmdBuildpkg.Env = append(os.Environ(), "PWD="+sourceDir)
	cmdBuildpkg.Stdout = os.Stdout
	cmdBuildpkg.Stderr = os.Stderr

	if err := cmdDhmake.Run(); err != nil {

		fmt.Printf("Deb build command `dh_make` error: %s\n", err)
		return false
	}

	if err := cmdBuildpkg.Run(); err != nil {

		fmt.Printf("Deb build command `dpkg-buildpackage` error: %s\n", err)
		return false
	}

	return true
}

func buildPackageRpm(ctx context, build PSettingsBuild) result {

	sourceDir, res := buildPrepSourceDir(ctx, build)
	if res == false {

		return false
	}

	cmdCmake := exec.Command("cmake", build.Rpm.CMake...)
	cmdMake := exec.Command("make", build.Rpm.Make...)

	cmdCmake.Dir = sourceDir
	cmdCmake.Env = append(os.Environ(), "PWD="+sourceDir)
	cmdCmake.Stdout = os.Stdout
	cmdCmake.Stderr = os.Stderr

	cmdMake.Dir = sourceDir
	cmdMake.Env = append(os.Environ(), "PWD="+sourceDir)
	cmdMake.Stdout = os.Stdout
	cmdMake.Stderr = os.Stderr

	if err := cmdCmake.Run(); err != nil {

		fmt.Printf("Rpm build command `cmake` error: %s\n", err)
		return false
	}

	if err := cmdMake.Run(); err != nil {

		fmt.Printf("Rpm build command `make` error: %s\n", err)
		return false
	}

	return true
}

func buildPrepSourceDir(ctx context, build PSettingsBuild) (string, result) {

	sourceDir := buildSourceDirPath(
		ctx.targetDir,
		build.Name,
		ctx.pSettings.ProjectName,
		ctx.pSettings.Version.Major,
		ctx.pSettings.Version.Minor,
		ctx.pSettings.Version.Patch)
	buildDir := buildDirPath(ctx.targetDir, build.Name)

	if len(ctx.origFile) > 0 {

		archType := tarGetArchType(ctx.origFile)
		dstOrig := buildDir + "/" + filepath.Base(ctx.origFile)

		if res := buildCheckOrigVersion(ctx); res == false {

			return "", false
		}

		if res := fopsMkdir(buildDir, 0755); res == false {

			return "", false
		}

		if res := fopsCopy(ctx.origFile, dstOrig); res == false {

			return "", false
		}

		switch archType {

		case TAR_TYPE_GZ:

			if res := tarOpenGz(dstOrig, sourceDir); res == false {

				return "", false
			}

		case TAR_TYPE_XZ:

			if res := tarOpenXz(dstOrig, sourceDir); res == false {

				return "", false
			}

		default:

			fmt.Printf("Wrong orig file archive, only `tar.gz` or `tar.xz` formats allowed\n")
			return "", false
		}

	} else {

		if res := fopsCopy(ctx.projectRoot, sourceDir); res == false {

			return "", false
		}
	}

	return sourceDir, true
}

func buildSourceDirPath(targetDir, buildName, projectName string, versionMajor, versionMinor, versionPatch int) string {

	return fmt.Sprintf("%s/%s/%s-%d.%d.%d",
		targetDir,
		buildName,
		projectName,
		versionMajor,
		versionMinor,
		versionPatch)
}

func buildDirPath(targetDir string, buildName string) string {

	return fmt.Sprintf("%s/%s",
		targetDir,
		buildName)
}

func buildOrigFilePath(targetDir, projectName string, versionMajor, versionMinor, versionPatch int, archExt string) string {

	return fmt.Sprintf("%s/orig/%s_%d.%d.%d.orig.%s",
		targetDir,
		projectName,
		versionMajor,
		versionMinor,
		versionPatch,
		archExt)
}

func buildCheckOrigVersion(ctx context) result {

	rgx := regexp.MustCompile(`^.*` + ctx.pSettings.ProjectName + `_([\d]+.[\d]+.[\d]+).orig.*$`)
	s := rgx.FindStringSubmatch(ctx.origFile)

	if len(s) == 0 {

		fmt.Printf("Wrong orig file name format\n")
		return false
	}

	oVersion := s[1]
	pVersion := fmt.Sprintf("%d.%d.%d",
		ctx.pSettings.Version.Major,
		ctx.pSettings.Version.Minor,
		ctx.pSettings.Version.Patch)

	if oVersion != pVersion {

		fmt.Printf("Mismatch of project and orig file versions\n")
		return false
	}

	return true
}
