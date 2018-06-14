package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type PSettingsBuildsRpm struct {
	CMake []string `yaml:"cmake"`
	Make  []string `yaml:"make"`
}

type PSettingsBuildsDeb struct {
	DHMake           []string `yaml:"dh_make"`
	DpkgBuildpackage []string `yaml:"dpkg_buildpackage"`
}

type PSettingsBuildsDocker struct {
	Img        string `yaml:"img"`
	Dockerfile string `yaml:"dockerfile"`
}

type PSettingsBuild struct {
	Name   string                `yaml:"name"`
	Env    map[string]string     `yaml:"env"`
	Deb    *PSettingsBuildsDeb   `yaml:"deb,omitempty"`
	Rpm    *PSettingsBuildsRpm   `yaml:"rpm,omitempty"`
	Docker PSettingsBuildsDocker `yaml:"docker"`
}

type PSettingsVersion struct {
	Major int `yaml:"major"`
	Minor int `yaml:"minor"`
	Patch int `yaml:"patch"`
}

type PSettings struct {
	ProjectName string           `yaml:"name"`
	Version     PSettingsVersion `yaml:"version"`
	Builds      []PSettingsBuild `yaml:"builds"`
}

func psettingsLoad(projectRoot string) (PSettings, result) {
	var c PSettings

	settingsFile, err := ioutil.ReadFile(projectRoot + "/" + pSettingsFile)
	if err != nil {

		fmt.Printf("Project settings file read error: %s\n", err.Error())
		return c, false
	}

	if err := yaml.Unmarshal(settingsFile, &c); err != nil {

		fmt.Printf("Project settings parse error: %s\n", err)
		return c, false
	}

	/*  */
	for i, b := range c.Builds {

		if len(b.Name) == 0 {

			fmt.Printf("Empty name for build[%d]\n", i)
			return c, false
		}

		if b.Deb == nil && b.Rpm == nil {

			fmt.Printf("Build \"%s\" must contain `deb` or `rpm` block\n", b.Name)
			return c, false
		}

		if b.Deb != nil && b.Rpm != nil {

			fmt.Printf("Build \"%s\" can't contain `deb` or `rpm` blocks together\n", b.Name)
			return c, false
		}
	}

	return c, true
}

func psettingsGet(ctx context, setting string) (string, result) {
	var value string

	switch setting {

	case settingGetProjectName:

		value = fmt.Sprintf("%s", ctx.pSettings.ProjectName)

	case settingGetVersionMajor:

		value = fmt.Sprintf("%d", ctx.pSettings.Version.Major)

	case settingGetVersionMinor:

		value = fmt.Sprintf("%d", ctx.pSettings.Version.Minor)

	case settingGetVersionPatch:

		value = fmt.Sprintf("%d", ctx.pSettings.Version.Patch)

	default:

		fmt.Printf("Unknown setting\n", setting)
		return "", false
	}

	return value, true
}
