#!/bin/bash
# 依赖 x11-apps scrot
KILLER_TEST_SCREENSHOT=${KILLER_TEST_SCREENSHOT:-test-%d.jpg}
KILLER_TEST_STDIO=${KILLER_TEST_STDIO:-test-%d.log}
KILLER_TEST_SCREENSHOT_TOTAL=5
KILLER_TEST_STEP=${KILLER_TEST_STEP:-1}
KILLER_TEST_TIMEOUT=${KILLER_TEST_TIMEOUT:-30}
KILLER_TEST_QUIT_TIMEOUT=${KILLER_TEST_QUIT_TIMEOUT:-5}
cleanup() {
  kill -9 $APP_PID &>/dev/null
}
trap cleanup EXIT
function log() {
    echo "[PID ${APP_PID:-?}]" "$@"
}
function check_alive(){
    if ! kill -0 $APP_PID 2>/dev/null;then
        log "进程提前退出，测试失败"
        exit 1
    fi
}
function step_sleep(){
    sleep ${KILLER_TEST_STEP}
}
function step_check(){
    step_sleep
    check_alive
}
function take_shot(){
    name=$(printf "${KILLER_TEST_SCREENSHOT}" $SECONDS)
    log $(printf "正在截图 %d秒: %s" $SECONDS "$name")
    rm -f "${name}"
    scrot -z "${name}"
}


log "启动进程:" "$@"
$@ &>"$KILLER_TEST_STDIO" &
APP_PID=$!

log "进程已启动"
for((i=0;i<${KILLER_TEST_TIMEOUT};++i));do
    PID=$(xdotool search --onlyvisible ".*"  2>/dev/null | xargs -r -I{} xdotool getwindowpid {}  2>/dev/null|head -n1)
    if [ -n "$PID" ];then
        log "已检测到窗口:PID=${PID}"
        break
    fi
    step_check
done

SECONDS=0
step_time=0
while test $SECONDS -lt $KILLER_TEST_SCREENSHOT_TOTAL;do
    min_step=$((KILLER_TEST_SCREENSHOT_TOTAL-SECONDS))
    sleep $((min_step<step_time?min_step:step_time))
    take_shot
    step_time=$((step_time>5?step_time:step_time+1))
done
step_check
log "正在通知并等待进程退出"
kill $APP_PID 2>/dev/null
KILL_OK=0
for((i=0;i<${KILLER_TEST_QUIT_TIMEOUT};++i));do
    step_sleep
    if ! kill $APP_PID 2>/dev/null;then
        KILL_OK=1
        log "进程已退出"
        break
    fi
done
if [ "$KILL_OK" -eq "0" ];then
    log "进程仍在运行，将被强制退出"
fi
exit 0

