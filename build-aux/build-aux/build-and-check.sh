#!/bin/bash

source $(dirname $0)/env.sh

$KILLER_EXEC build $KILLER_BUILD_ARGS -- build-aux/ldd-check.sh > missing.log
if [ ! -s missing.log ];then
    echo "没有找到缺失文件，跳过。">&2
    exit 0
fi
$KILLER_EXEC apt -- apt-file update
$KILLER_EXEC apt -- build-aux/ldd-search.sh missing.log found.log notfound.log
if [ ! -s found.log ];then
    echo "错误：没有找到任何库。">&2
    exit 1
fi
$KILLER_EXEC build $KILLER_BUILD_ARGS -- sh -c 'apt install -y $(cat found.log)'