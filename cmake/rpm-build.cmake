if(RPM OR SRPM)

	if(EXISTS "${CMAKE_ROOT}/Modules/CPack.cmake")
		include(InstallRequiredSystemLibraries)

		set(CPACK_GENERATOR "RPM")

		#set(CPACK_RPM_PACKAGE_DEBUG "ON")

		# Set the common CentOS directories to exclude them from rpm package
		set(CPACK_RPM_EXCLUDE_FROM_AUTO_FILELIST_ADDITION
			"/etc/cron.d"
			"/usr/sbin"
			"/etc/cron.d"
			"/etc/logrotate.d"
			"/etc/sudoers.d"
			"/var"
			"/var/log")

		if(SRPM)
			set(CPACK_RPM_PACKAGE_SOURCES "ON")
		endif(SRPM)

		set(CPACK_PACKAGE_DESCRIPTION_FILE "${CMAKE_CURRENT_SOURCE_DIR}/build-scope/tpls/centos/description")
		set(CPACK_PACKAGE_DESCRIPTION_SUMMARY "Nixys build tools")
		set(CPACK_PACKAGE_VENDOR $ENV{RPM_VENDOR})
		set(CPACK_PACKAGE_CONTACT $ENV{RPM_CONTACT})
		set(CPACK_RPM_PACKAGE_LICENSE $ENV{RPM_LICENSE})
		set(CPACK_PACKAGE_VERSION_MAJOR "${MAJOR_VERSION}")
		set(CPACK_PACKAGE_VERSION_MINOR "${MINOR_VERSION}")
		set(CPACK_PACKAGE_VERSION_PATCH "${PATCH_VERSION}")
		set(CPACK_PACKAGE_FILE_NAME "${CMAKE_PROJECT_NAME}_${MAJOR_VERSION}.${MINOR_VERSION}.${CPACK_PACKAGE_VERSION_PATCH}")
		set(CPACK_SOURCE_PACKAGE_FILE_NAME "${CMAKE_PROJECT_NAME}_${MAJOR_VERSION}.${MINOR_VERSION}.${CPACK_PACKAGE_VERSION_PATCH}")
		set(CPACK_RPM_PACKAGE_REQUIRES "cmake >= 2.8, rpm-build, rpmdevtools, make")

		set(CPACK_RPM_PRE_INSTALL_SCRIPT_FILE "${CMAKE_CURRENT_SOURCE_DIR}/build-scope/tpls/centos/preinstall")
		set(CPACK_RPM_POST_INSTALL_SCRIPT_FILE "${CMAKE_CURRENT_SOURCE_DIR}/build-scope/tpls/centos/postinstall")
		set(CPACK_RPM_PRE_UNINSTALL_SCRIPT_FILE "${CMAKE_CURRENT_SOURCE_DIR}/build-scope/tpls/centos/preuninstall")
		set(CPACK_RPM_POST_UNINSTALL_SCRIPT_FILE "${CMAKE_CURRENT_SOURCE_DIR}/build-scope/tpls/centos/postuninstall")

		set(CPACK_COMPONENTS_ALL Libraries ApplicationData)
		include(CPack)
	endif(EXISTS "${CMAKE_ROOT}/Modules/CPack.cmake")
endif(RPM OR SRPM)
