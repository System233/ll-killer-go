#!/bin/bash
source $(dirname $0)/env.sh
SHARE_DIR="/opt/apps/${LINGLONG_APPID}/files/share"
ERRORS=()
function log_error() {
    ERRORS+=("$@")
    echo -e "\033[31m错误: $@ \033[0m" >&2
}
DIR_LIST=()
function check_dir() {
    if [ -d "$1" ]; then
        DIR_LIST+=("$1")
    fi
}
check_dir "${SHARE_DIR}/systemd/user"
check_dir "${SHARE_DIR}/applications"
if [ "${#DIR_LIST}" -eq "0" ]; then
    echo "无可测试项"
    exit 0
fi
while read desktop; do
    while read icon; do
        if [ -z "$icon" ];then
            continue
        fi
        found=$(find "${SHARE_DIR}/icons" "/usr/share/icons" "/usr/share/pixmaps" \( -name "${icon}.xpm" -o -name "${icon}.png" -o -name "${icon}.svg" \) -print -quit)
        if [ ! -e "$found" ]; then
            log_error "[失败] $desktop:$icon: 找不到此图标"
        else
            echo "[成功] $desktop: $icon"
        fi
    done <<<$(grep -oP "^Icon.*?=\K.*$" "$desktop")

    while read args; do
        ARGS=($args)
        MAIN=${ARGS[0]}
        TEST_CMD=(which "${MAIN}")
        if [ "$MAIN" = "$ENTRYPOINT" ]; then
            MAIN="${ARGS[1]}"
            TEST_CMD=("$ENTRYPOINT" which "${MAIN}")
        fi
        if ! "${TEST_CMD[@]}" >/dev/null; then
            log_error "[失败] $desktop:$MAIN: 找不到此启动文件"
        else
            echo "[成功] $desktop: $args"
        fi
    done <<<$(grep -oP "^Exec(Start)?.*?=.*?--\s*\K.*$" "$desktop")
done <<<$(find "${DIR_LIST[@]}" \( -name "*.desktop" -o -name "*.service" -o -name "*.conf" \))

exit ${#ERRORS[@]}
