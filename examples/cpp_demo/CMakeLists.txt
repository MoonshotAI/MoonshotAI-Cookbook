cmake_minimum_required(VERSION 3.16...3.28 FATAL_ERROR)
# 设置项目名称和使用的语言
project(main LANGUAGES CXX)

# Setup vcpkg script with CMake (note: should be placed before project() call)
#set(CMAKE_TOOLCHAIN_FILE "C:/Users/25181/Documents/GitHub/buildscripts/windows/vcpkg/scripts/buildsystems/vcpkg.cmake" CACHE STRING "Vcpkg toolchain file") # ${CMAKE_CURRENT_SOURCE_DIR}
# C++17 is required
set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS ON)

# 指定要编译的源文件
add_executable(${PROJECT_NAME} main.cpp)

find_package(CURL REQUIRED)
if (CURL_FOUND)
    include_directories(${CURL_INCLUDE_DIRS})
    target_link_libraries(${PROJECT_NAME} ${CURL_LIBRARIES})
endif()