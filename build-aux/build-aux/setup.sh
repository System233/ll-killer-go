#!/bin/bash
source $(dirname $0)/env.sh

KILLER_EXEC=${KILLER_EXEC:-$(which ll-killer)}

if [ -z "$KILLER_EXEC" ];then
    echo "错误：未找到ll-killer，请确保通过ll-killer执行命令，或手动设置KILLER_EXEC环境变量。" >&2
    exit 1
fi

echo "[准备文件系统]"
setup-filesystem.sh

echo "[复制必要文件]"
test -e "build-aux/fuse-overlayfs" && cp -avf "build-aux/fuse-overlayfs" "$PREFIX"
cp -avf "$KILLER_EXEC" "$PREFIX/ll-killer"
cp -avf "build-aux/$ENTRYPOINT_NAME" "$PREFIX"

echo "[调整文件布局]"
if [ -e "$PREFIX/share" ]; then
    mkdir -pv "$PREFIX/usr"
    cp -arfT "$PREFIX/share" "$PREFIX/usr/share"
    rm -rf "$PREFIX/share"
fi
test -d "$PREFIX/usr/share" && mv -vf "$PREFIX/usr/share" "$PREFIX/share"
mkdir -pv "$PREFIX/usr/share"

if [ -d "$PREFIX/opt/apps/" ]; then
    echo "[合并share目录]"
    find "$PREFIX/opt/apps/" -type d \( -path "$PREFIX/opt/apps/*/entries" \
        -or -path "$PREFIX/opt/apps/*/files/share" \) \
        -exec "merge-share.sh" "{}" \;
fi

if [ "${KILLER_PACKER:-0}" == "0" ]; then
    echo "[修正符号链接]"
    echo "详细信息: https://github.com/OpenAtom-Linyaps/linyaps/issues/1039"
    find $PREFIX/share -xtype l -exec "relink.sh" "{}" \;
fi

echo "[配置快捷方式]"
find $PREFIX/share/applications -name "*.desktop" -exec "setup-desktop.sh" "{}" \;

if [ -d "$PREFIX/share/applications/context-menus" ]; then
    echo "[配置右键菜单]"
    find "$PREFIX/share/applications/context-menus" -name "*.conf" -exec "setup-desktop.sh" "{}" \;
fi

if [ -d "$PREFIX/etc/systemd" ]; then
    echo "[配置服务单元]"
    find "$PREFIX/share/systemd" -name "*.service" -type f -exec "setup-systemd.sh" "{}" \;
    find "$PREFIX/etc/systemd" -name "*.service" -type f -exec "setup-systemd.sh" "{}" \;
    find "$PREFIX/lib/systemd" -name "*.service" -type f -exec "setup-systemd.sh" "{}" \;
fi

if [ -d "$PREFIX/share/dbus-1/services" ]; then
    echo "[配置Dbus服务]"
    find "$PREFIX/share/dbus-1/services" -name "*.service" -type f -exec "setup-dbus.sh" "{}" \;
fi
