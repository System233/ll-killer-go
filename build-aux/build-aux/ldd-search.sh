#!/bin/bash
FOUND=${2:-/dev/stdout}
MISSING=${3:-/dev/stderr}
INPUT=${1:--}
if [ "$INPUT" == "-" ];then
    TMP_DIR=$(mktemp -d ll-killer.XXXXXX -p /tmp)
    TMP_FILE="$TMP_DIR/soname.list"
    cat $INPUT>"$TMP_FILE"
    INPUT="$TMP_FILE"
fi
SEARCHED=$(cat "$INPUT" | xargs -P$(nproc) -I{} sh -c 'apt-file find -x "{}$"| grep -P "^lib|/usr/lib/x86_64-linux-gnu/" | sort -urk2 | head -n1')
grep -Fvf <(echo "$SEARCHED" | awk -F/ '{print $NF}') "$INPUT" >$MISSING
echo "$SEARCHED" |grep -vP "^\s*$"| cut -d: -f1 | sort -u >$FOUND

if [ -n "$TMP_DIR" ];then
    rm -rf "$TMP_DIR"
fi