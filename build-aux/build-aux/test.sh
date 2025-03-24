#!/bin/bash
# 依赖 xwd scrot
set -e
if [ "$#" -lt "2" ];then
    echo "错误：无效参数"
    echo "用法：$0 <应用APPID> <layer文件> [空格分隔排除列表: NoDisplay Hidden Terminal]"
    exit 1
fi
APPID="$1"
LAYER="$2"
FILTER_LIST="$3"
NEEDS=()
function check_dep(){
    dep=$1
    need=${2:-$1}
    if ! which $dep >/dev/null;then
        echo "错误：${dep} 命令未找到" >&2
        NEEDS=("${NEEDS[@]}" "$need")
    fi
}
check_dep xvfb-run xvfb
check_dep xdotool xdotool
check_dep scrot scrot

if [ ! ${#NEEDS} -eq 0 ];then
    echo "请尝试安装这些库：${NEEDS[@]}" >&2
    exit 1
fi

echo "正在测试: APPID=${APPID} LAYER=${LAYER}"
LLCLI_VER=$(apt-cache show linglong-bin|grep -oP "^Version:\s*\K.*"||LANG=en ll-cli --version | grep -oP "version\s*\K.*"||echo "0.0.0-unknown")
SUDO=
if dpkg --compare-versions "$LLCLI_VER" ge "1.7.0"; then
    echo "提示：玲珑1.7.x需要root权限进行安装卸载，将使用sudo执行安装卸载命令。"
    SUDO=sudo
fi

$SUDO ll-cli uninstall "$APPID" &>/dev/null ||true
$SUDO ll-cli install "$LAYER"

cleanup() {
    $SUDO ll-cli uninstall "$APPID"
}
trap cleanup EXIT

echo "[正在测试快捷方式/服务单元]"
ll-cli run "${APPID}" -- "$(dirname $0)/test-desktop.sh" || true

i=0
echo "[正在测试启动项]"

mkdir -p tests
TASKLOG="$PWD/tests/task.log"
ll-cli run "${APPID}" -- "$(dirname $0)/test-extract-exec.sh" "$TASKLOG" "$FILTER_LIST"

while read args; do
    ARGS=($args)
    KILLER_TEST_SCREENSHOT="tests/screen$i-%d.jpg" \
        KILLER_TEST_STDIO="tests/output-$i.log" \
        xvfb-run -a "$(dirname $0)/test-display.sh" \
        "${ARGS[@]}" || true
    i=$((i + 1))
done <"$TASKLOG"

echo "测试完成，可在tests中查看测试输出。"
