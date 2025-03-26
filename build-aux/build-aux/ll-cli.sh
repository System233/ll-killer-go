#!/bin/bash

CMD="$1"
if [ "$CMD" == "--version" ];then
    echo "ll-cli.sh version 1.0.0"
    exit 0
fi
if [ "$CMD" == "--help" ];then
    echo "玲珑青春版 - 本命令使用ll-killer模拟一个玲珑环境，用于替换ll-cli以加速测试过程。"
    echo "使用LAYER,ROOTFS环境变量指定layer文件和base文件系统,RUNTIME指定runtime文件系统。"
    exit 0
fi
if [ "$CMD" != "run" ];then
    echo "$0 忽略命令：$@"
    exit 0
fi
PKG="$2"
SEARCHED_LAYER=($PKG*.layer)
LAYER=${LAYER:-$SEARCHED_LAYER}
ROOTFS=${ROOTFS}
RUNTIME=${RUNTIME}

echo "包名:$PKG"
echo "Layer文件:$LAYER"
echo "Base文件系统:$ROOTFS"
echo "Runtime文件系统:$RUNTIME"

if [ ! -e "$LAYER" ];then
    echo "未正确设置LAYER变量:$LAYER"
    exit 1
fi
if [ ! -d "$ROOTFS" ];then
    echo "未正确设置ROOTFS变量:$ROOTFS"
    exit 1
fi

while test "$#" -gt "0";do
    shift 1
    if [ "$1" == "--" ];then
        break
    fi
done
shift 1

TMP_DIR=$(mktemp -d)
MNT_DIR="${TMP_DIR}/mnt"
MERGED_DIR="${TMP_DIR}/merged"
OVERRIDE_DIR="${TMP_DIR}/override"
UPPER_DIR="${TMP_DIR}/upper"
WORK_DIR="${TMP_DIR}/work"
mkdir -p "${MNT_DIR}" "$MERGED_DIR" "$OVERRIDE_DIR" "$UPPER_DIR" "$WORK_DIR"

ll-killer layer mount "$LAYER" "$MNT_DIR"
cleanup() {
    ll-killer layer umount "$MNT_DIR"
    chmod 777 -R $WORK_DIR
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT
LINGLONG_APPID=$PKG ll-killer exec \
    --mount "overlay:$MERGED_DIR::overlay:lowerdir=$OVERRIDE_DIR:$ROOTFS,upperdir=$UPPER_DIR,workdir=$WORK_DIR" \
    --mount "$MNT_DIR/files:$MERGED_DIR/opt/apps/$PKG/files" \
    --mount "$RUNTIME:$MERGED_DIR/runtime" \
    --mount "tmpfs:$MERGED_DIR/run::tmpfs" \
    --mount "/run:$MERGED_DIR/run::merge: " \
    --mount "/:$MERGED_DIR/run/host/rootfs:rbind:rdonly+bind" \
    --mount "/etc/resolv.conf:$MERGED_DIR/etc/resolv.conf:rdonly+bind" \
    --mount "/etc/localtime:$MERGED_DIR/etc/localtime:rdonly+bind" \
    --mount "/etc/machine-id:$MERGED_DIR/etc/machine-id:rdonly+bind" \
    --mount "/etc/timezone:$MERGED_DIR/etc/timezone:rdonly+bind" \
    --mount "/etc/passwd:$MERGED_DIR/etc/passwd:rdonly+bind" \
    --mount "/etc/locale.conf:$MERGED_DIR/etc/locale.conf:rdonly+bind" \
    --mount "/etc/default/locale:$MERGED_DIR/etc/default/locale:rdonly+bind" \
    --mount "/usr/share/fonts:$MERGED_DIR/usr/share/fonts:rdonly+bind" \
    --mount "/usr/share/locale:$MERGED_DIR/usr/share/locale:rdonly+bind" \
    --mount "/usr/share/theme:$MERGED_DIR/usr/share/theme:rdonly+bind" \
    --mount "/usr/share/icons:$MERGED_DIR/usr/share/icons:rdonly+bind" \
    --mount "$MERGED_DIR:$MERGED_DIR:remount+rdonly" \
    --mount "/proc:$MERGED_DIR/proc:rbind" \
    --mount "/dev:$MERGED_DIR/dev:rbind" \
    --mount "/sys:$MERGED_DIR/sys:rbind" \
    --mount "/tmp:$MERGED_DIR/tmp:rbind" \
    --mount "/home:$MERGED_DIR/home:rbind" \
    --mount "/root:$MERGED_DIR/root:rbind" \
    --rootfs "$MERGED_DIR" \
    --wait \
    -- $@