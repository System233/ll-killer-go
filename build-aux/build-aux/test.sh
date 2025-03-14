#!/bin/bash
# 依赖 xwd scrot
set -e
if [ "$#" != "2" ];then
    echo "错误：无效参数"
    echo "用法：$0 <应用APPID> <layer文件>"
    exit 1
fi
DEPS=("xvfb-run" "xdotool" scrot)
function check_dep(){
    for dep in "${DEPS[@]}";do
        if ! which $dep >/dev/null;then
            echo "错误：${dep}命令未找到，请尝试安装这些依赖: xvfb xdotool scrot" >&2
            exit 1
        fi
    done
}
check_dep
source $(dirname $0)/env.sh
APPID="$1"
LAYER="$2"
mkdir -p tests
TMPLOG="tests/task.log"
echo "正在测试: APPID=${APPID} LAYER=${LAYER}"
LLCLI_VER=$(LANG=en ll-cli --version | grep -oP "version\s*\K.*")
SUDO=
if dpkg --compare-versions "$LLCLI_VER" \>= "1.7.0"; then
    echo "提示：玲珑1.7.x需要root权限进行安装卸载，将使用sudo执行安装卸载命令。"
    SUDO=sudo
fi

$SUDO ll-cli uninstall "$APPID"||true
$SUDO ll-cli install "$LAYER"
echo "[正在测试快捷方式]"
ll-cli run "${APPID}" -- "$(dirname $0)/test-desktop.sh" "$TMPLOG"
i=0

# TODO: 检测并排除KDE Plasma Panels
echo "[正在测试启动命令]"
while read args; do
    ARGS=($args)
    KILLER_TEST_SCREENSHOT="tests/screen$i-%d.jpg" \
        KILLER_TEST_STDIO="tests/output-$i.log" \
        xvfb-run -a "$(dirname $0)/test-display.sh" \
        ll-cli run "${APPID}" -- "${ARGS[@]}"
    i=$((i + 1))
done <"$TMPLOG"
$SUDO ll-cli uninstall "$1"

echo "全部测试成功，你可以在tests中找到测试输出。"
