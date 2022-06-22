#!/bin/bash

# 检测错误，遇到返回码为非零值，直接退出进程
set -e

# 检测未定义的变量，如果存在则退出
#set -u
set -o pipefail

# 这个PS4变量方便我们在调试bash 的时候可以看到行号
PS4='+[${BASH_SOURCE}]第${LINENO}行:{${FUNCNAME[0]}}: '

# CMD_WAIT 本命令函数针对的是单一后台进程执行失败进行的wait并退出父进程的功能
CMD_WAIT(){
    PID=$1

    wait "$PID"

    CODE=$?

    echo "当前返回值 ${CODE}"
}


wait_echo(){
    wait "$1"
    CODE=$?

    if [[ $CODE -gt 0 ]];then
        echo "等待进程退出,退出码$CODE"
    fi
}

startapp(){
    CMD_WAIT $!
}

exitWithCode(){
    if [[ $CODE -gt 0 ]];then
        echo "等待进程退出,退出码$CODE"
        exit $CODE 
    fi
}

FILENAME="./data/gin-test.csv"

killall "$APPNAME" || true

while read -r -a line 
do
    : $((i++))

    if [ "$i" -eq 1 ];then continue ;fi
    echo "读取文件第${i}行"
    port=$(echo "${line[0]}" | tr -d '" \r\t\n')
    echo "port:$port"

    ./"${APPNAME}" -p "${port}"
    
    CMD_WAIT $!

    exitWithCode

done < $FILENAME




