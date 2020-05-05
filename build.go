package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/nixys/nxs-build-tools/arch"
	"github.com/nixys/nxs-build-tools/fops"
)

func buildMakeOrig(ctx selfContext) error {

	origDir := buildSourceDirPath(
		ctx.targetDir,
		"orig",
		ctx.conf.ProjectName,
		ctx.conf.Version.Major,
		ctx.conf.Version.Minor,
		ctx.conf.Version.Patch)

	origFileTgz := buildOrigFilePath(
		ctx.targetDir,
		ctx.conf.ProjectName,
		ctx.conf.Version.Major,
		ctx.conf.Version.Minor,
		ctx.conf.Version.Patch,
		"tar.gz")

	origFileTxz := buildOrigFilePath(
		ctx.targetDir,
		ctx.conf.ProjectName,
		ctx.conf.Version.Major,
		ctx.conf.Version.Minor,
		ctx.conf.Version.Patch,
		"tar.xz")

	if err := fops.CopyWithIgnores(ctx.projectRoot, origDir); err != nil {
		return err
	}

	/* Make tar.gz */
	if err := arch.Make(origDir, origFileTgz); err != nil {
		return err
	}

	/* Make tar.xz */
	if err := arch.Make(origDir, origFileTxz); err != nil {
		return err
	}

	return fops.Remove(origDir)
}

func buildPackage(ctx selfContext) error {

	for _, b := range ctx.conf.Builds {

		if b.Name == ctx.buildName {

			/* Set env variables for specified build */
			for key, value := range b.Env {
				os.Setenv(key, value)
			}

			if b.Deb != nil {
				return buildPackageDeb(ctx, b)
			}

			if b.Rpm != nil {
				return buildPackageRpm(ctx, b)
			}

			return fmt.Errorf("no rules specified for build (build name: %s)", ctx.buildName)
		}
	}

	return fmt.Errorf("build not exist (build: %s)", ctx.buildName)
}

func buildPackageDeb(ctx selfContext, build buildConf) error {

	sourceDir, err := buildPrepSourceDir(ctx, build)
	if err != nil {
		return err
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
		return fmt.Errorf("deb build command `dh_make` error: %s", err)
	}

	if err := cmdBuildpkg.Run(); err != nil {
		return fmt.Errorf("deb build command `dpkg-buildpackage` error: %s", err)
	}

	return nil
}

func buildPackageRpm(ctx selfContext, build buildConf) error {

	sourceDir, err := buildPrepSourceDir(ctx, build)
	if err != nil {
		return err
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
		return fmt.Errorf("rpm build command `cmake` error: %s", err)
	}

	if err := cmdMake.Run(); err != nil {
		return fmt.Errorf("rpm build command `make` error: %s", err)
	}

	return nil
}

func buildPrepSourceDir(ctx selfContext, build buildConf) (string, error) {

	sourceDir := buildSourceDirPath(
		ctx.targetDir,
		build.Name,
		ctx.conf.ProjectName,
		ctx.conf.Version.Major,
		ctx.conf.Version.Minor,
		ctx.conf.Version.Patch)

	buildDir := buildDirPath(ctx.targetDir, build.Name)

	if len(ctx.origFile) > 0 {

		dstOrig := buildDir + "/" + filepath.Base(ctx.origFile)

		if err := buildCheckOrigVersion(ctx); err != nil {
			return "", err
		}

		if err := fops.MkdirRecursive(buildDir, 0755); err != nil {
			return "", err
		}

		if err := fops.CopyWithIgnores(ctx.origFile, dstOrig); err != nil {
			return "", err
		}

		if err := arch.Open(dstOrig, sourceDir); err != nil {
			return "", err
		}
	} else {
		if err := fops.CopyWithIgnores(ctx.projectRoot, sourceDir); err != nil {
			return "", err
		}
	}

	return sourceDir, nil
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

func buildCheckOrigVersion(ctx selfContext) error {

	regexString := `^.*` + ctx.conf.ProjectName + `_([\d]+.[\d]+.[\d]+).orig.*$`

	rgx := regexp.MustCompile(regexString)
	s := rgx.FindStringSubmatch(ctx.origFile)

	if len(s) == 0 {
		return fmt.Errorf("wrong orig file name format: does not satisfy regex (orig file: \"%s\", regex: \"%s\")", ctx.origFile, regexString)
	}

	oVersion := s[1]
	pVersion := fmt.Sprintf("%d.%d.%d",
		ctx.conf.Version.Major,
		ctx.conf.Version.Minor,
		ctx.conf.Version.Patch)

	if oVersion != pVersion {
		return fmt.Errorf("mismatch project and orig file versions (project version: \"%s\", orig file version: \"%s\")", pVersion, oVersion)
	}

	return nil
}
