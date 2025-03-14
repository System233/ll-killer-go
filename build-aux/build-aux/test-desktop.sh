#!/bin/bash
source $(dirname $0)/env.sh
ARGS_OUTPUT="$1"
SHARE_DIR="/opt/apps/${LINGLONG_APPID}/files/share"
ERRORS=()
function log_error() {
    ERRORS=("${ERRORS[@]}" "$@")
    echo -e "\033[31m错误: $@ \033[0m" >&2
}
echo -n >"$ARGS_OUTPUT"
while read desktop; do
    while read icon; do
        found=$(find "${SHARE_DIR}/icons" "/usr/share/icons" "/usr/share/pixmaps" \( -name "${icon}.xpm" -o  -name "${icon}.png" -o -name "${icon}.svg" \) -print -quit)
        if [ ! -e "$found" ]; then
            log_error "$desktop:$icon: 找不到此图标"
        else
            echo "[Icon OK] $icon"
        fi
    done <<<$(grep -oP "^Icon.*?=\K.*$" "$desktop")

    while read args; do
        ARGS=($args)
        MAIN=${ARGS[0]}
        TEST_CMD=(which "${MAIN}")
        if [ "$MAIN" = "$ENTRYPOINT" ];then
            MAIN="${ARGS[1]}"
            TEST_CMD=("$ENTRYPOINT" which "${MAIN}")
        fi
        if ! "${TEST_CMD[@]}" >/dev/null; then
            log_error "$desktop:$MAIN: 找不到此启动文件"
        else
            echo "[Exec OK] $args"
        fi
        echo "$args" >>"$ARGS_OUTPUT"
    done <<<$(grep -oP "^Exec.*?=.*?--\s*\K.*$" "$desktop")
done <<<$(find "${SHARE_DIR}/applications/" -name "*.desktop")

exit ${#ERRORS[@]}