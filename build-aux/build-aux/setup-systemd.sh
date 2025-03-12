#!/bin/bash
set -e
source $(dirname $0)/env.sh
SRC="$1"
DSTDIR="$PREFIX/share/systemd/user"
DST="$DSTDIR/$(basename $SRC)"
mkdir -p "$DSTDIR"
sed -i -E -e "s:^\s*ExecStart\s*=:ExecStart=$ENTRYPOINT :g" -e '/^User=/d' -e '/WantedBy/ s/multi-user.target/default.target/' "$SRC"
if mv -Tv "$SRC" "$DST";then
    RVL=$(realpath --relative-to="$(dirname "$SRC")" "$DST")
    ln -svTf "$RVL" "$SRC"
fi