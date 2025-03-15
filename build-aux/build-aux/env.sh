#!/bin/bash
if [ -z "$ENV_SETUPED" ];then
    CWD=$(dirname $(readlink -f "$0"))
    export ENV_SETUPED=1
    export PATH=$PATH:$CWD
    export ENTRYPOINT_NAME=${ENTRYPOINT_NAME:-entrypoint.sh}
    export ENTRYPOINT=${ENTRYPOINT:-/opt/apps/$LINGLONG_APPID/files/$ENTRYPOINT_NAME}
fi