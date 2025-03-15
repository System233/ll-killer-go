#!/bin/bash
SHARE_DIR="/opt/apps/${LINGLONG_APPID}/files/share"
awk -v RS='\n\\[Desktop Entry\\]' '
    /NoDisplay\s*=\s*true|Hidden\s*=\s*true|Terminal\s*=\s*true/ { next }
    match($0, /\nExec=([^\n]*)/, exec) { print exec[1] }
' $(find "${SHARE_DIR}/applications" -name "*.desktop") >"$1"
