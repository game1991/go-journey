#!/bin/bash
# FS:    输入记录字段间的分隔符 field separate
# RS：   输入记录的分隔符 row separate
# OFS：  输出记录字段间的分隔符 out field separate
# ORS：  输出记录的分隔符  out row separate
# 
# NR：   当前行数 num row
# NF：   当前记录字段数  num field
# ARGC： 命令行变元数 argc
PS4='+${BASH_SOURCE}:${LINENO}:${FUNCNAME[0]}: '
set -ex
set -o pipefail

#APPNAME=app.tp
echo "$?"

[ -z "${APPNAME}" ] && { echo "程序名称不能为空" ;exit 99 ; }

# 启动程序前先杀掉所有已启动的程序，避免端口占用服务无法启动的问题，如果命令失败可以忽略
#killall "${APPNAME}"

readonly FILENAME="./data/data_copy.csv"
i=0

OLD_IFS=${IFS}
IFS=','

[ ! -f ${FILENAME} ] && { echo "${FILENAME} file not found " ;exit 99 ; }

if [ "$1" = 1 ] ;then awk -F, -f cal.awk ${FILENAME}
else 
    grep -v '^ *#' ${FILENAME} | while read -r line && [ -n "${line}" ]
    do 
        # 计数当前第几行
        : $((i++))
        NUM=${i}
        echo "${line}"

        [ "${NUM}" -eq 1 ] && { echo "当前第[${NUM}]行是csv标题行不做处理" ;continue ; }
     
        app_id=$(echo "${line}" | awk -F, '{print $1}' | tr -d '" \r\t\n')
        app_key=$(echo "${line}" | awk -F, '{print $2}' | tr -d '"\r\t\n')
        app_port=$(echo "${line}" | awk -F, '{print $NF}' | tr -d '"\r\t\n')
        
        echo 应用id="${app_id}" 应用key="${app_key}" 应用端口port="${app_port}"

        sleep 3

        if [[ ${app_id} = "" || ${app_key} = "" || ${app_port} = "" ]]
        then echo "${APPNAME} 当前第[${NUM}]行 启动命令存在空值:[app_id]${app_id};[app_key]${app_key};[app_port]${app_port}" 
        exit 99
        fi 
        
        #start=$(./"${APPNAME}" -id "${app_id}" -key "${app_key}" -p "${app_port}" &)


        
        # else 
        #     ./"${APPNAME}" -id "${app_id}" -key "${app_key}" -p "${app_port}" &
        # fi
        done 
fi

# while命令允许在while语句行定义多条命令，但是只有最后一条测试命令的退出状态 决定循环是否停止。

echo "数据读取完毕"

#./"${APPNAME}" &

IFS=${OLD_IFS}

#|| {  }

#if ! 
#then
#echo "${APPNAME} 当前第[${NUM}]行 启动程序失败:[app_id]${app_id};[app_key]${app_key};[app_port]${app_port}" ;exit 99
#fi