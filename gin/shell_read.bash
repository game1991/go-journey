#!/bin/bash
# FS:    输入记录字段间的分隔符 field separate
# RS：   输入记录的分隔符 row separate
# OFS：  输出记录字段间的分隔符 out field separate
# ORS：  输出记录的分隔符  out row separate
# 
# NR：   当前行数 num row
# NF：   当前记录字段数  num field
# ARGC： 命令行变元数 argc

FILENAME="./data/data.csv"
EXEC="$APP_NAME"
i=1
OLD_IFS="$IFS"
IFS=","

if [ "$1" = 1 ];then awk -F, -f cal.awk $FILENAME;fi

[ ! -f $FILENAME ] && { echo "$FILENAME 数据文件file not found" ;exit 99; }

[ ! -e "$EXEC" ] && { echo "$EXEC 执行程序不存在" ;exit 99; }

#-r 如果给出此选项，则反斜杠不会充当转义字符。
while read -r -a array || { echo "$FILENAME 内容读取失败" ;exit ; } 
do
    echo "${array[@]}"
    
    id=${array[0]}

    key=${array[1]}

    port=${array[2]}

    if ((i > 1));then
    echo 应用id="$id" 应用key="$key" 应用端口port="$port"

    { ./"$EXEC" -id "$id" -key "$key" -p "$port" & }

    fi
    # Alternative POSIX version that does not preserve the exit code
    : $((i++))
done < $FILENAME

#./"$EXEC" ||exit 0 ;

# 还原分隔符
IFS="$OLD_IFS"
