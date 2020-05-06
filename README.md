# Nxs-build-tools

Nxs-build-tools provides tools to create deb and rpm packages for your projects.

## Getting started

The nxs-build-tools is built on top of CMake project build system. So if you are already using CMake you may need to make additional settings to resolve possible conflicts.

To be able to use this tools for your project, please follow these recommendations:
1.  If you need to place the project source code into separate directory within the project root (e.g. `src/`), 
    set variable `PROJECT_SRC_DIR` to appropriate value (see details below).
2.  For projects based on compiled programming languages  make sure the binaries located within `objs/` directory.
    Otherwise set variable `PROJECT_BIN_DIR` to appropriate value (see details below).

## Preparing the project environment

First you need to populate your project with the template files from nxs-build-tools package. 
To do this execute following command from project root dir or use flag `--project-root` with specified directory:
```
nxs-build-tools --command=populate
```
Notice: To be sure that nxs-build-tools does not brake your project, the populate command will fail if any of template files already exist in your project root.

Command creates new files and directories in your project root:
```
.
├── build-scope
│   ├── pkg
│   │   ├── general
│   │   └── os
│   │       ├── centos-6
│   │       ├── centos-7
│   │       ├── debian-7
│   │       ├── debian-8
│   │       └── debian-9
│   └── tpls
│       ├── centos
│       │   ├── description
│       │   ├── postinstall
│       │   ├── postuninstall
│       │   ├── preinstall
│       │   └── preuninstall
│       └── debian
│           ├── changelog
│           ├── compat
│           ├── control
│           ├── postinst
│           ├── postrm
│           ├── preinst
│           ├── prerm
│           ├── README.Debian
│           └── README.source
├── cmake
│   ├── app-python.spec
│   ├── general_install.cmake
│   ├── golang.cmake
│   ├── helpers_install.cmake
│   ├── python.cmake
│   └── rpm-build.cmake
├── CMakeLists.txt
└── .proj-settings.yml
```

* `build-scope`: directory contains directories and files used for build deb or rpm packages
  * `pkg`: directory with parts of packages content, i.e. /etc/, /usr/share/ and other files and directories. For example, if your application uses the configuration file (e.g. /etc/your-app/your-app.conf) place this file to build-scope/pkg/general/etc/your-app/your-app.conf.
      * `general`: directory with the OS independent packages content.
      * `os`: directory with the OS specific packages content. If your project has files with same names but different content for different OS (e.g. /etc/init.d/your-app) - place them to corresponding `os` subdirectory.
    * `tpls`: directory contains the configuration files and scripts to create deb or rpm packages. With the nxs-build-tools you can create separate packages for specific releases to each OS. So if you need to have different packages for Debian 7, 8 and 9 (as example) - create corresponding subdirectories within the `tpls` directory (e.g. 'debian-7', 'debian-8' and 'debian-9') with appropriate content. Later you will be able to use this configurations in `.proj-settings.yml` file to create specific packages build.
* `cmake`: directory contains CMake modules. Please also read the [CMake documentation](https://cmake.org/documentation/).
  * `general_install.cmake`: module describes the instalations of your packages. See [this](https://github.com/nixys/nxs-build-tools/blob/master/cmake/general_install.cmake) for example.
  * `golang.cmake`: module for build Golang projects. If your project is not the Golang project - you need to exclude this module in main `CMakeLists.txt` file.
    Module provides following variables:
      * `GO_VERSION_FILE_TPL`: path to CMake template file that contains the Go code with CMake variables to specify the version for your apllication.
      * `GO_VERSION_FILE`: path to destination version file (after CMake substitutes). You may include this file to your source code to display the apllication version (e.g. with '--version' arg).
      * The `install` command defines application binary file destination path within the packages. You may change the source and/or destination path.
  * `python.cmake`: module for build Python projects. If your project is the Python project - you need to uncomment this module in main `CMakeLists.txt` file.
    Module provides following variables:
      * `PYTHON_VERSION_FILE_TPL`: path to CMake template file that contains the Python code with CMake variables to specify the version for your apllication.
      * `PYTHON_VERSION_FILE`: path to destination version file (after CMake substitutes). You may include this file to your source code to display the apllication version (e.g. with '--version' arg).
      * `PYTHON_SPEC_FILE`: Python spec file for PyInstaller to build your application.
      * `PYTHON_MODULES`: list of the Python modules used in your project. Specified modules will be downloaded and installed during project build.
      * The `install` command defines application binary file destination path within the packages. You may change the source and/or destination path.
  * `helpers_install.cmake`: module contains helpers to use them with CMake in your project.
  * `rpm-build.cmake`: module contains the settings for an rpm packages builds. You may need to consult CMake documentation to tuning this file for your project.
* `CMakeLists.txt`: it is main CMake file that contains a specific settings for your project such as `PROJECT_BIN_DIR` and `PROJECT_SRC_DIR` or CMake modules includes.
  Module provides follow variables:
    * `PROJECT_BIN_DIR`: the name of directory within the project root with the application binaries, uses for projects based on compiled programming languages (e.g. Go projects).
    * `PROJECT_SRC_DIR`: the name of directory within the project root with the application source code. Empty string by default. Set this variable if you need to store you source code in separate dir inside the project root (e.g. `src/`).
* `.proj-settings.yml`: file contains the settings for your project, such as packages build configurations. See details of this file below.

### .proj-settings.yml file

Every project uses nxs-build-tools to build packages need the file `.proj-settings.yml` within project directory. It is a file in yaml format and consits of following fields:
* `name` (optional): the name of your project. This value used as name for packages. You can override this option with the `--package-name` command line argument. Useful for CI.
* `version` (optional): defines the packages version. You can override this option with the `--package-version` command line argument (see [semantic versioning](https://semver.org/)). Useful for CI.
* `builds`: array of builds description. Each element of this array describes a specific options to build either `deb` or `rpm` packages.
  * `name`: the name of package build. This value uses to specify the name of package build by nxs-build-tools `--build-name` arg.
  * `env` (optional): environment variables list specified in `VARIABLE_NAME: VARIABLE_VALUE` format. This may useful for CMake build process for deb and rpm packages.
  * `deb`: block disribes options to build 'deb' packages.
      * `dh_make`: array with an options for 'dh_make'. See dh_make man for details.
      It is important to note that argument `--templates`, specifies the template directory in "build-scope/tpls/"" with the configuration files to build deb package.
      * `dpkg_buildpackage`: array with an options for 'dpkg-buildpackage'. See dpkg-buildpackage man for details.
  * `rpm`: block describes options to build 'rpm' packages.
      * `cmake`: array with an options for 'cmake' which are used to prepare project for building rpm packages. See cmake documentation for details.
      In addition to `env` section you may specify flags for CMake such as `-DRPM=on` or `-DSRPM=on` to define a rpm build process. For example, with flag `-DSRPM=on` specified you get an rpm source package (srpm).
      * `make`: array with an options for 'make'. See make man for details.

### .gitignore file

It is important to note that nxs-build-tools environment contains the files and directories to be excluded from package build process. After your project has been populated by nxs-build-tools templates files the list with recommended excludes will be offered to appended into your .gitignore file.

## Build packages

After the nxs-build-tools environment is prepared for the project you may build the packages.

There is two ways to make package:
* Create the package directly from source code
* Create the archive from source first and then build the deb or rpm package from this archive (useful for CI).

### Direct package build

To make packages directly use the following command either from your project root or using `--project-root` arg:
```
nxs-build-tools --build-name=debian --package-name=some-project --package-version=0.0.1
```
where `--build-name` arg defines the appropriate build name from `.proj-settings.yml`.

This will give you a deb package in the "builds/debian" directory within your project root directory.

### Build package via original file

This is the more correct way and may be used for following purposes:
* If your project has a separatid builds for different releases of the same OS (e.g. deb packages for Debian 8, Debian 9, Debian 10) and you have original files with same names (e.g. "some-project_0.0.1.orig.tar.xz"), but different md5 hash sum, you will not be able to upload your packages into deb repository. In this case you need to use same original file for every Debian build.
* If your project uses CI/CD process (e.g. Gitlab CI).
  In this case (in addition to case above) you'll be able to make CI/CD process more effective and optimal due to separation into different stages.

First you need to prepare the original files:
```
nxs-build-tools --command=make-orig --package-name=some-project --package-version=0.0.1
```

It makes the .tar.gz and .tar.xz files with your source code. You can find the result original files in "builds/orig" directory within your project root (e.g. some-project_0.0.1.orig.tar.gz and some-project_0.0.1.orig.tar.xz).

Now you can build the deb or rpm packages from these origs. For example:
```
nxs-build-tools --orig-file=builds/orig/some-project_0.0.1.orig.tar.xz --build-name=debian --package-name=some-project --package-version=0.0.1
```

As in the previous case you'll get a deb package in the "builds/debian" directory within your project root.

## Example of usage 

The simple example of nxs-build-tools usage.

Suppose you have a simple Go project with following structure:
```
.
└── main.go
```

The main.go file content has something like:
```
package main

import "fmt"

func main() {
        fmt.Println("Hello!")
}
```

To create the deb package you should do:

1.  Change your directory the your project root and execute command:
    ```
    nxs-build-tools --command=populate
    ```

2.  In accordance with previous command output create the .gitignore file:
    ```
    cat <<EOF >> /tmp/some-project/.gitignore
    /builds
    /objs
    _CPack_Packages
    CMakeCache.txt
    CMakeFiles/
    Makefile
    CPackConfig.cmake
    CPackSourceConfig.cmake
    cmake_install.cmake
    EOF
    ```

3.  Build the package:
    ```
    nxs-build-tools --build-name=debian --package-name=some-project --package-version=0.0.1
    ```

After that you may observe your deb package in the "builds/debian" directory that can be installed in the Debian system.

## Install nxs-build-tools

### Debian

1.  Add Nixys repository key:

    ```
    apt-key adv --fetch-keys http://packages.nixys.ru/debian/repository.gpg.key
    ```

2.  Add the repository. Currently Debian wheezy, jessie, stretch and buster are available:

    ```
    echo "deb [arch=amd64] http://packages.nixys.ru/debian/ wheezy main" > /etc/apt/sources.list.d/packages.nixys.ru.list
    ```

    ```
    echo "deb [arch=amd64] http://packages.nixys.ru/debian/ jessie main" > /etc/apt/sources.list.d/packages.nixys.ru.list
    ```

    ```
    echo "deb [arch=amd64] http://packages.nixys.ru/debian/ stretch main" > /etc/apt/sources.list.d/packages.nixys.ru.list
    ```

    ```
    echo "deb [arch=amd64] http://packages.nixys.ru/debian/ buster main" > /etc/apt/sources.list.d/packages.nixys.ru.list
    ```

3.  Make an update:

    ```
    apt-get update
    ```

4.  Install nxs-build-tools:

    ```
    apt-get install nxs-build-tools
    ```

### CentOS

1.  Add Nixys repository key:

    ```
    rpm --import http://packages.nixys.ru/centos/repository.gpg.key
    ```

2.  Add the repository. Currently CentOS 6 and 7 are available:

    ```
    cat <<EOF > /etc/yum.repos.d/packages.nixys.ru.repo
    [packages.nixys.ru]
    name=Nixys Packages for CentOS \$releasever - \$basearch
    baseurl=http://packages.nixys.ru/centos/\$releasever/\$basearch
    enabled=1
    gpgcheck=1
    gpgkey=http://packages.nixys.ru/centos/repository.gpg.key
    EOF
    ```

3.  Install nxs-build-tools:

    ```
    yum install nxs-build-tools
    ```
