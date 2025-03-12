#!/bin/bash

if [ -n "$KILLER_ENTRYPOINT" ];then
    exec "${@:-bash}"
    exit 1
fi

if [ -e ".killer-debug" ];then
    export KILLER_DEBUG=1
fi
export KILLER_ENTRYPOINT="/entrypoint.sh"
APP_DIR=${APP_DIR:-"/opt/apps/$LINGLONG_APPID/files"}
OVERLAY_EXEC_PATH=$APP_DIR/fuse-overlayfs
if [ ! -e "$OVERLAY_EXEC_PATH" ];then
    OVERLAY_EXEC_PATH=$(which "fuse-overlayfs")
fi
OVERLAY_EXEC=${OVERLAY_EXEC:-"$OVERLAY_EXEC_PATH"}
KILLER_EXEC=${KILLER_EXEC:-"$APP_DIR/ll-killer"}
if (test -z "$NO_OVERLAYFS" && test -e /run/host/rootfs/dev/fuse ) 2>/dev/null >/dev/null; then
    if ! "$OVERLAY_EXEC" --version 2>/dev/null >/dev/null;then
        OVERLAY_EXEC=""
    fi
    exec $APP_DIR/ll-killer exec \
        --fuse-overlayfs "$OVERLAY_EXEC" \
        --fuse-overlayfs-args "$OVERLAY_ARGS" \
        --mount "/dev:/run/app.dev:rbind" \
        --mount "/run/host/rootfs/dev:/dev:rbind" \
        --mount "/:/run/app.oldfs:rbind" \
        --mount "overlay:$APP_DIR/usr/share::fuse-overlayfs:lowerdir=$APP_DIR/share:/run/app.oldfs/usr/share,squash_to_root,static_nlink" \
        --mount "overlay:/run/app.rootfs::fuse-overlayfs:lowerdir=$APP_DIR:/run/app.oldfs,squash_to_root" \
        --mount "/proc:/run/app.rootfs/proc:rbind" \
        --mount "/run/host/rootfs/dev:/run/app.rootfs/dev:rbind" \
        --mount "/run:/run/app.rootfs/run:rbind" \
        --mount "/sys:/run/app.rootfs/sys:rbind" \
        --mount "/tmp:/run/app.rootfs/tmp:rbind" \
        --mount "/home:/run/app.rootfs/home:rbind" \
        --mount "/root:/run/app.rootfs/root:rbind" \
        --mount "/run/app.dev/pts:/run/app.rootfs/dev/pts:rbind" \
        --mount "/run/app.dev/shm:/run/app.rootfs/dev/shm:rbind" \
        --mount "/run/app.dev/mqueue:/run/app.rootfs/dev/mqueue:rbind" \
        --mount "/opt/apps/$LINGLONG_APPID:/run/app.rootfs/opt/apps/$LINGLONG_APPID:rbind" \
        --rootfs /run/app.rootfs \
	    --no-bind-rootfs \
        --socket=/run/app.unix \
        --wait \
        -- "$KILLER_ENTRYPOINT" "${@:-bash}"
fi
exec $APP_DIR/ll-killer exec \
    --mount "$APP_DIR/share:$APP_DIR/usr/share:rbind" \
    --mount "rootfs:/run/app.rootfs::tmpfs" \
    --mount "/+$APP_DIR:/run/app.rootfs::merge" \
    --mount "/run/host/rootfs/dev:/run/app.rootfs/dev:rbind" \
    --mount "/dev/pts:/run/app.rootfs/dev/pts:rbind" \
    --mount "/dev/shm:/run/app.rootfs/dev/shm:rbind" \
    --mount "/dev/mqueue:/run/app.rootfs/dev/mqueue:rbind" \
    --mount "/proc:/run/app.rootfs/proc:rbind" \
    --mount "/run:/run/app.rootfs/run:rbind" \
    --mount "/sys:/run/app.rootfs/sys:rbind" \
    --mount "/tmp:/run/app.rootfs/tmp:rbind" \
    --mount "/home:/run/app.rootfs/home:rbind" \
    --mount "/root:/run/app.rootfs/root:rbind" \
    --socket=/run/app.unix \
    --rootfs=/run/app.rootfs \
    --no-bind-rootfs \
    --wait \
    -- "$KILLER_ENTRYPOINT" "${@:-bash}"