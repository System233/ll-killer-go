#!/bin/bash

# 检查是否提供了包名
if [[ -z "$1" ]]; then
    echo "删除不属于指定包的desktop文件"
    echo "用法: $0 <deb包名>"
    exit 1
fi

PACKAGE_NAME="$1"

# 定义要检查的目录
DESKTOP_DIRS=(
    "/usr/share/applications"
    "/usr/local/share/applications"
)

# 遍历 /opt 目录，查找可能包含 .desktop 文件的子目录
for opt_dir in /opt/apps/*/files/share/applications /opt/apps/*/entires/applications; do
    [[ -d "$opt_dir" ]] && DESKTOP_DIRS+=("$opt_dir")
done


# 处理每个目录
for DIR in "${DESKTOP_DIRS[@]}"; do
    [[ -d "$DIR" ]] || continue  # 如果目录不存在则跳过

    find "$DIR" -maxdepth 1 -type f -name "*.desktop" | while read -r FILE; do
        # 检查文件是否属于指定包
        if ! dpkg -S "$FILE" 2>/dev/null | grep -q "^$PACKAGE_NAME:"; then
            rm -fv "$FILE"
        fi
    done
done
