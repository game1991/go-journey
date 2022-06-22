#!/bin/bash
read -r -p "请输入任意数字: " val
real_val=66
if [ "$val" -gt "$real_val" ]
then
   echo "输入值大于等于预设值"
else
   echo "输入值比预设值小"
fi