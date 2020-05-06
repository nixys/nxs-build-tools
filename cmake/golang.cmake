set(GO_VERSION_FILE_TPL "${CMAKE_CURRENT_SOURCE_DIR}/${PROJECT_SRC_DIR}version.go.in")
set(GO_VERSION_FILE "${CMAKE_CURRENT_SOURCE_DIR}/${PROJECT_SRC_DIR}version.go")

if(EXISTS ${GO_VERSION_FILE_TPL})
	configure_file("${GO_VERSION_FILE_TPL}" "${GO_VERSION_FILE}" @ONLY)
endif()

add_custom_target(${PROJECT_NAME}
	ALL
	COMMAND go build -o ${CMAKE_CURRENT_BINARY_DIR}/${PROJECT_BIN_DIR}${PROJECT_NAME}
	WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/${PROJECT_SRC_DIR})

install(PROGRAMS ${CMAKE_CURRENT_BINARY_DIR}/${PROJECT_BIN_DIR}${PROJECT_NAME} DESTINATION bin)
