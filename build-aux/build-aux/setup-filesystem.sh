#!/bin/bash

test -d "linglong/filesystem/diff" && cp -arfT "linglong/filesystem/diff" "$PREFIX"
find $PREFIX \( -type c,b,p,s -or -name ".wh.*" \) -exec rm -rf {} \;
chmod a+Xr -R "$PREFIX"
rm -fv "$PREFIX/etc/resolv.conf" \
    "$PREFIX/etc/localtime" \
    "$PREFIX/etc/timezone" \
    "$PREFIX/etc/machine-id"
