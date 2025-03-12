#!/bin/bash
if [ -z "$ENV_SETUPED" ];then
    export ENV_SETUPED=1
    export PATH=$PATH:$(dirname $0):$PWD
    export ENTRYPOINT_NAME=${ENTRYPOINT_NAME:-entrypoint.sh}
    export ENTRYPOINT=${ENTRYPOINT:-/opt/apps/$LINGLONG_APPID/files/$ENTRYPOINT_NAME}
    export KILLER_EXEC=${KILLER_EXEC:-$(which ll-killer)}
    
    if [ -z "$KILLER_EXEC" ];then
        echo "错误：未找到ll-killer，请确保ll-killer在当前或build-aux目录中，或使用'll-killer script $*'重新执行命令。" >&2
        exit 1
    fi
fi