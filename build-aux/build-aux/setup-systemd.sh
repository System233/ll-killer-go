#!/bin/bash
if [[ -z "$1" ]]; then
    echo "导出指定包的systemd服务"
    echo "用法: $0 <deb包名>"
    exit 1
fi

PACKAGE_NAME="$1"

set -e
source $(dirname $0)/env.sh

DSTDIR="/usr/share/systemd/user"
dpkg -S ".service" 2>/dev/null  | grep '.service$' | grep -oP "^$PACKAGE_NAME:\K.*$"| while read -r FILE; do
    SRC="$FILE"
    DST="$DSTDIR/$(basename $SRC)"
    mkdir -p "$DSTDIR"
    sed -i -E -e "s:^\s*ExecStart\s*=:ExecStart=$ENTRYPOINT :g" -e '/^User=/d' -e '/WantedBy/ s/multi-user.target/default.target/' "$SRC"
    if mv -Tv "$SRC" "$DST";then
        RVL=$(realpath --relative-to="$(dirname "$SRC")" "$DST")
        ln -svTf "$RVL" "$SRC"
    fi
done