#!/bin/bash
# NoDisplay Hidden Terminal
SHARE_DIR="/opt/apps/${LINGLONG_APPID}/files/share"
OUTPUT="$1"
FILTER_LIST=($2)
if [ ${#FILTER_LIST[@]} -eq 0 ]; then
  FILTER_LIST=(NoDisplay Hidden)
fi
echo "启动项排除列表：${FILTER_LIST[@]}"
FILTER=$(printf "%s\s*=\s*true\n" "${FILTER_LIST[@]}"|paste -sd '|')
awk -v RS='\n\\[Desktop Entry\\]' "
    /${FILTER}/ { next }
    match(\$0, /\nExec=([^\n]*)/, exec) { print exec[1] }
" $(find "${SHARE_DIR}/applications" -name "*.desktop" -not -path "*/screensavers/*") | sed -E -e "s:['\"]?%[fFuUdDNcik]['\"]?::g" >"$OUTPUT"

