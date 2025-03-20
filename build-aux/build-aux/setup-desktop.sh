#!/bin/bash
set -e
source $(dirname $0)/env.sh

DESKTOP="$1"
echo "$DESKTOP"
sed -E -i -e "s:^\s*Exec\s*=:Exec=$ENTRYPOINT :g" "$DESKTOP"
while read icon; do
    if [[ $icon == /* ]]; then
        icon_mapped=$(echo $icon | sed -e "s:^/usr/share:/share:")
        real_path="$PREFIX/$icon_mapped"
        icon_name_ext=$(basename "$icon")
        extension=$(echo "${icon_name_ext##*.}" | tr '[:upper:]' '[:lower:]')
        icon_name="${LINGLONG_APPID}-${icon_name_ext%.*}"
        if [ "$extension" = "xpm" ];then
            mv "$real_path" "$PREFIX/share/pixmaps/${icon_name}.xpm"
        else
            setup-icon.sh "$real_path" "$icon_name"
        fi
        sed -E -i -e "s:$icon:$icon_name:g" "$DESKTOP"
    fi
done <<<$(grep -oP "^\s*Icon\s*=\s*\K.*$" "$DESKTOP")
