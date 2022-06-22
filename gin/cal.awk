#!/bin/awk -f

# 在awk语言的任意动作区间，即在{}之内，均可随时定义变量，无须事先说明。但一般情况下是在BEGIN中定义变量并赋以初值，在动作区域内使用。
# BEGIN的作用是在处理第一条记录之前将BEGIN后面大括号之内的动作运行且只运行一次，也就是BEGIN匹配第一个输入记录之前。
# END的作用是在处理完最后一条记录之后将END后面大括号之内的动作运行且只运行一次，也就是说，END匹配最后一个输入记录之后。

#运行前
BEGIN {
    OFS="):"
    printf "-----------------awk脚本执行前-------------------\n"
}

#运行中

$1 ~ /6$/ , NR>1 {
    print NR,$1
    printf "\n"
    #app.tmp -id $id -key $key -p $port
}
#运行后
END {
    printf "---------------awk脚本执行后--------------------\n"
}
