#!/bin/bash
source $(dirname $0)/env.sh

CLI="/usr/bin/ll-cli run ${LINGLONG_APPID} -- "
TRY_CLI='ll-cli'
REPL=(sed -i -E -e "s:^(\s*Exec\s*=):\1$CLI:g"
    -e "s:^(\s*ExecStart\s*=):\1$CLI:g"
    -e "/^\s*TryExec\s*=/c TryExec=${TRY_CLI}"
    -e "/\[Desktop Entry\]/a X-linglong=${LINGLONG_APPID}"
    '{}')

if [ -d "$PREFIX/share/applications" ]; then
    echo "*处理快捷方式和右键菜单"
    find "$PREFIX/share/applications" \( -name "*.desktop" -or -name "*.conf" \) -print -exec "${REPL[@]}" \;
fi

if [ -d "$PREFIX/share/dbus-1/services" ]; then
    echo "*处理D-Bus服务"
    find "$PREFIX/share/dbus-1/services" -name "*.service" -print -exec "${REPL[@]}" \;
fi

if [ -d "$PREFIX/share/systemd/user" ]; then
    echo "*处理Systemd服务"
    find "$PREFIX/share/systemd/user" -name "*.service" -print -exec "${REPL[@]}" \;
fi
