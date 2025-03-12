#!/bin/bash
set -e
source $(dirname $0)/env.sh
echo "$1"
sed -i -E -e "s:^\s*Exec\s*=:Exec=$ENTRYPOINT :g" -e '/^User=/d' "$1"