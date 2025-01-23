cmake_minimum_required(VERSION 3.16)

project(fe VERSION 0.1 LANGUAGES CXX)

set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(Qt6 6.5 REQUIRED COMPONENTS Quick WebSockets)

qt_standard_project_setup(REQUIRES 6.5)

qt_add_executable(appfe
    main.cpp
)

qt_add_qml_module(appfe
    URI fe
    VERSION 1.0
    QML_FILES
        Main.qml
        SOURCES datamodel.h datamodel.cpp
        SOURCES datatype.h
        SOURCES datamodel_rep.h datamodel_rep.cpp
        SOURCES config.h config.cpp
        RESOURCES settings.ini
        QML_FILES Kanban.qml
        QML_FILES Card.qml
)

file (COPY ${CMAKE_CURRENT_SOURCE_DIR}/settings.ini DESTINATION ${CMAKE_CURRENT_BINARY_DIR})

# Qt for iOS sets MACOSX_BUNDLE_GUI_IDENTIFIER automatically since Qt 6.1.
# If you are developing for iOS or macOS you should consider setting an
# explicit, fixed bundle identifier manually though.
set_target_properties(appfe PROPERTIES
#    MACOSX_BUNDLE_GUI_IDENTIFIER com.example.appfe
    MACOSX_BUNDLE_BUNDLE_VERSION ${PROJECT_VERSION}
    MACOSX_BUNDLE_SHORT_VERSION_STRING ${PROJECT_VERSION_MAJOR}.${PROJECT_VERSION_MINOR}
    MACOSX_BUNDLE TRUE
    WIN32_EXECUTABLE TRUE
)

target_link_libraries(appfe
    PRIVATE Qt6::Quick Qt6::WebSockets
)

include(GNUInstallDirs)
install(TARGETS appfe
    BUNDLE DESTINATION .
    LIBRARY DESTINATION ${CMAKE_INSTALL_LIBDIR}
    RUNTIME DESTINATION ${CMAKE_INSTALL_BINDIR}
)
