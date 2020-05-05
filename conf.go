package main

import (
	"github.com/nixys/nxs-build-tools/fops"
	"github.com/nixys/nxs-go-conf"
)

type confOpts struct {
	ProjectName string      `conf:"name"`
	Version     versionConf `conf:"version"`
	Builds      []buildConf `conf:"builds" conf_extraopts:"required"`
}

type versionConf struct {
	Major int `conf:"major"`
	Minor int `conf:"minor"`
	Patch int `conf:"patch"`
}

type buildConf struct {
	Name   string            `conf:"name" conf_extraopts:"required"`
	Env    map[string]string `conf:"env"`
	Deb    *buildsDebConf    `conf:"deb"`
	Rpm    *buildsRpmConf    `conf:"rpm"`
	Docker buildsDockerConf  `conf:"docker"`
}

type buildsDebConf struct {
	DHMake           []string `conf:"dh_make"`
	DpkgBuildpackage []string `conf:"dpkg_buildpackage"`
}

type buildsRpmConf struct {
	CMake []string `conf:"cmake"`
	Make  []string `conf:"make"`
}

type buildsDockerConf struct {
	Img        string `conf:"img"`
	Dockerfile string `conf:"dockerfile"`
}

func confRead(confPath string) (confOpts, error) {

	var c confOpts

	p, err := fops.Normalize(confPath)
	if err != nil {
		return c, err
	}

	err = conf.Load(&c, conf.Settings{
		ConfPath:    p,
		ConfType:    conf.ConfigTypeYAML,
		UnknownDeny: true,
	})
	if err != nil {
		return c, err
	}

	return c, nil
}
