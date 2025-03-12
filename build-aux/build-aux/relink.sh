#!/bin/sh
DEST="$1"
test -d "$DEST" || exit
PTR=$(readlink "$DEST")
test "${PTR:0:1}" == "/" || exit
# SRC=$(realpath -m "$DEST" | sed -e "s:^/usr/share:/share:")
# ln -svf "${PREFIX}$SRC" "$DEST"
SRC=$(realpath -m "$DEST")
ln -svf "$SRC" "$DEST"
