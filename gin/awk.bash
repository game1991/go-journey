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
i=1
#APPNAME=app.tmp
if [ "$1" = 1 ] ;then awk -F, -f cal.awk $FILENAME;
else 
    grep -v '^ *#' < $FILENAME | while IFS= read -r _
    do
        #echo "$line"
        NUM=$i
        [ $NUM -eq 1 ] && { : $((i++)) ; continue ; }

        app_id=$(awk -F, 'NR>1 && NR=='$NUM' {print $1}' "$FILENAME")
        app_key=$(awk -F, 'NR>1 && NR=='$NUM' {print $2}' "$FILENAME")
        app_port=$(awk -F, 'NR>1 && NR=='$NUM' {print $NF}' "$FILENAME")
        
        echo 应用id="$app_id" 应用key="$app_key" 应用端口port="$app_port"

        #./"$APPNAME" -id "$app_id" -key "$app_key" -p "$app_port" &

        : $((i++))
    done
fi

# ./$APPNAME &