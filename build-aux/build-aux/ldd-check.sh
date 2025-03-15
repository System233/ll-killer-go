#!/bin/bash
MODE=$1
TMP_DIR=$(mktemp -d ll-killer.XXXXXX -p /tmp)
TMP_FILE="$TMP_DIR/soname.list"
DIR_LIST="/opt /usr /lib /bin /runtime"
if [ "$MODE" = "fast" ];then
    DIR_LIST="/opt/apps/${LINGLONG_APPID}/files"
fi
find $DIR_LIST -name "*.so*" | xargs -rn1 basename | sort -u >$TMP_FILE
find $DIR_LIST '(' -name "*.so" -or -executable ')' | xargs -r ldd 2>/dev/null | grep -F "=> not found" | sort -u | grep -oP '^\s*\K\S+' | grep -vFxf "$TMP_FILE"
rm -rf $TMP_DIR
