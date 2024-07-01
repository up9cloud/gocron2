#!/usr/bin/env bash

# 生成压缩包 xx.tar.gz或xx.zip
# 使用 ./package.sh -a amd64 -p linux -v v2.0.0

# 任何命令返回非0值退出
set -o errexit
# 使用未定义的变量退出
set -o nounset
# 管道中任一命令执行失败退出
set -o pipefail

eval $(go env)

# 提取git最新tag作为应用版本
VERSION=''
# 最新git commit id
GIT_COMMIT_ID=''

# 外部输入的系统
INPUT_OS=()
# 外部输入的架构
INPUT_ARCH=()
# 未指定OS，默认值
DEFAULT_OS=${GOHOSTOS}
# 未指定ARCH,默认值
DEFAULT_ARCH=${GOHOSTARCH}
# 支持的系统
SUPPORT_OS=(linux darwin windows)
# 支持的架构
SUPPORT_ARCH=(386 amd64 arm64)

# 编译参数
LDFLAGS=''
# 编译文件生成目录
BUILD_DIR='build'
# 打包文件生成目录
PACKAGE_DIR='packages'
# 需要額外打包的文件
INCLUDE_FILE=()

# 获取git 最新tag name
git_latest_tag() {
    local commit_id=""
    local tag_name=""
    commit_id=`git rev-list --tags --max-count=1`
    tag_name=`git describe --tags "${commit_id}"`
    echo ${tag_name}
}

# 获取git 最新commit id
git_latest_commit() {
    echo "$(git rev-parse --short HEAD)"
}

# 打印信息
print_message() {
    echo "$1"
}

# 打印信息后推出
print_message_and_exit() {
    if [[ -n $1 ]]; then
        print_message "$1"
    fi
    exit 1
}

# 初始化
init() {
    # 设置系统
    if [[ ${#INPUT_OS[@]} = 0 ]];then
        INPUT_OS=("${DEFAULT_OS}")
    fi
    for OS in "${INPUT_OS[@]}"; do
        if [[  ! "${SUPPORT_OS[*]}" =~ ${OS} ]]; then
            print_message_and_exit "不支持的系统${OS}"
        fi
    done
    # 设置CPU架构
    if [[ ${#INPUT_ARCH[@]} = 0 ]];then
        INPUT_ARCH=("${DEFAULT_ARCH}")
    fi
    for ARCH in "${INPUT_ARCH[@]}";do
        if [[ ! "${SUPPORT_ARCH[*]}" =~ ${ARCH} ]]; then
            print_message_and_exit "不支持的CPU架构${ARCH}"
        fi
    done

    if [[ -z "${VERSION}" ]];then
        VERSION=`git_latest_tag`
    fi
    GIT_COMMIT_ID=`git_latest_commit`
    LDFLAGS="-w -X 'main.AppVersion=${VERSION}' -X 'main.BuildDate=`date '+%Y-%m-%d %H:%M:%S'`' -X 'main.GitCommit=${GIT_COMMIT_ID}'"

    mkdir -p ${BUILD_DIR}
    mkdir -p ${PACKAGE_DIR}
}

# 编译
build() {
    # 二进制文件名
    local binary_name="$1"
    # main函数所在文件
    local main_file="./cmd/${binary_name}/main.go"
    local output_filename=''
    for os in "${INPUT_OS[@]}";do
        if [[ "${os}" = "windows"  ]];then
            output_filename=${binary_name}.exe
        else
            output_filename=${binary_name}
        fi
        for arch in "${INPUT_ARCH[@]}";do
            local build_target_dir="${binary_name}-${VERSION}-${os}-${arch}"
            env CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${build_target_dir}/${output_filename} ${main_file}
        done
    done
}

# 打包文件
package_file() {
    for item in "${INCLUDE_FILE[@]}"; do
        cp -r ../${item} $1
    done
}

# 打包
package_binary() {
    local binary_name="$1"
    cd ${BUILD_DIR}
    local root_dir=${OLDPWD}
    for os in "${INPUT_OS[@]}";do
        for arch in "${INPUT_ARCH[@]}";do
        local build_target_dir="${binary_name}-${VERSION}-${os}-${arch}"
        package_file ${build_target_dir}
        cd $build_target_dir
        if [[ "${os}" = "windows" ]];then
            local package_output_filepath="../../${PACKAGE_DIR}/${build_target_dir}.zip"
            # zip -rq ${package_output_filepath} .
            7z a -tzip -mx9 -r ${package_output_filepath} .
            7z l ${package_output_filepath}
        else
            local package_output_filepath="../../${PACKAGE_DIR}/${build_target_dir}.tar.gz"
            tar -I 'gzip -9' -cf ${package_output_filepath} .
            tar -ztvf ${package_output_filepath}
        fi
        cd ${OLDPWD}
        done
    done
    cd ${root_dir}
}

# 清理
clean() {
    if [[ -d ${BUILD_DIR} ]];then
        rm -rf ${BUILD_DIR}
    fi
    # if [[ -d ${PACKAGE_DIR} ]];then
    #     rm -rf ${PACKAGE_DIR}
    # fi
}

package_all_cmd() {
    init

    binary_name='gocron2'
    build ${binary_name}
    package_binary ${binary_name}

    binary_name='gocron2-node'
    build ${binary_name}
    package_binary ${binary_name}

    clean
}

# p 平台 linux darwin windows
# a 架构 386 amd64
# v 版本号  默认取git最新tag
while getopts "p:a:v:" OPT;
do
    case ${OPT} in
    p) IPS=',' read -r -a INPUT_OS <<< "${OPTARG}"
    ;;
    a) IPS=',' read -r -a INPUT_ARCH <<< "${OPTARG}"
    ;;
    v) VERSION=$OPTARG
    ;;
    *)
    ;;
    esac
done

package_all_cmd
