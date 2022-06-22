#!/bin/bash

PS4='+${BASH_SOURCE}:${LINENO}:${FUNCNAME[0]}: '
set -x
set -o pipefail

wait_and_echo() {
    for PID in "${!PIDS[@]}"; do
        PIDS[$PID]=1
    done
    while [ ${#PIDS[@]} -ne 0 ]; do
        wait -n -p PID "${!PIDS[@]}"
        CODE=$?
        echo PID $PID terminated with exit code $CODE
        unset 'PIDS[$PID]'
    done
}

handle_sigchld() {
    for PID in "${!PIDS[@]}"; do
        if [ ! -d "/proc/$PID" ]; then
            wait $PID
            CODE=$?
            echo PID $PID terminated with exit code $CODE
            unset 'PIDS[$PID]'
        fi
    done
}

PIDS=()

trap handle_sigchld SIGCHLD


show

function killapp(){
    killall "${APPNAME}" &
    PIDS[$!]=1
}

killapp

# wait_and_echo $!

# echo $?

# killall "${APPNAME}" 

command=./"${APPNAME}"

function start(){
    sleep 3
    echo "当前循环第$1次"
    local num="$1"
    [ "$num" -eq 4 ] && { echo "当前是第4任务，测试错误场景" ; ./app & }

    $command &
}


echo "执行bash结束了"

echo $?


