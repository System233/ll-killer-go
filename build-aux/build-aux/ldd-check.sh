#!/bin/bash
TMP_DIR=$(mktemp -d ll-killer.XXXXXX -p /tmp)
TMP_FILE="$TMP_DIR/soname.list"
DIR_LIST="/opt /usr /lib /bin /runtime"
find $DIR_LIST -name "*.so*" | xargs -n1 basename | sort -u >$TMP_FILE
find $DIR_LIST '(' -name "*.so" -or -executable ')' | xargs ldd 2>/dev/null | grep -F "=> not found" | sort -u | grep -oP '^\s*\K\S+' | grep -vFxf "$TMP_FILE"
rm -rf $TMP_DIR
