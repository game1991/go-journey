#!/bin/bash

set -e
set -o pipefail

# 这个PS4变量方便我们在调试bash 的时候可以看到行号
PS4='+${BASH_SOURCE}:${LINENO}:{${FUNCNAME[0]}}: '

wait_and_echo() {
    sleep 3
    PID=$1
    echo "waiting PID $PID to terminate"
    wait "$PID"
    CODE=$?

    echo "PID $PID terminated with exit code $CODE"
    return $CODE
}

function exitWithCode(){
    if [ "$CODE" -gt 0 ];then
    echo "wait_and_echo执行的对象返回值为:$CODE" 
    exit "$CODE" 
    fi
}

killapp(){
    killall gin-test &
    wait_and_echo $!
    exitWithCode
}

startapp() {
    ./app.tt &
    wait_and_echo $! 
    exitWithCode
}

killapp

while IFS=" " read -r -a line
do
# 因为使用了-e 意味着只要遇到非零返回码就会退出，那么这个表达式((i+))返回码正好是从0到1，故需要特殊处理，不能直接用((i++))，一种方案是使用:这个表达式，意思是这个操作不做任何处理，一种是利用 || true来使整个表达式返回码为0
    : $((i++))
    echo "我是第${i}行"
    echo "第一列${line[0]}"
    echo "第二列${line[1]}"

    if [ "$i" -eq 2 ];then
        ./app.tt &
        wait_and_echo $! 
        exitWithCode
    fi

    startapp
done <<< "1 2 3 4 
2 2 3 4
3 2 3 4"

echo '我结束了'
