cmake_minimum_required(VERSION 2.8)

# Set env variables
set(DEB $ENV{DEB})
set(RPM $ENV{RPM})
set(SRPM $ENV{SRPM})

set(PROJECT_NAME $ENV{PROJECT_NAME})
set(MAJOR_VERSION $ENV{PKG_VERSION_MAJOR})
set(MINOR_VERSION $ENV{PKG_VERSION_MINOR})
set(PATCH_VERSION $ENV{PKG_VERSION_PATCH})

project(${PROJECT_NAME} NONE)

set(PROJECT_BIN_DIR "objs/")
set(PROJECT_SRC_DIR "")

set(CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH} "${CMAKE_CURRENT_SOURCE_DIR}/cmake")

include(helpers_install)
include(general_install)

# Remove unnecessary modules
include(golang)
include(rpm-build)
