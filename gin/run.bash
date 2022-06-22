#!/bin/bash

#set -x

killall app.tmp

list=("ac172f739b3542248e862488b5956046,4109a3bd03394763a92865eddcbb6083,8880" "5a09f25d4dff4ddf856ca2316da39b81,e0d1be4248a14f14b5dd7cf05a1e0045,8881" "0742e54e93f5480185a8f5897c379a8c,a026dd49284e48e3bf9226e696020e0c,8882" "c49c303bf6d84c08bc2473d6ed3fc0d4,066f32c5e87e482aa43f7e314c20fd52,8883" "721f158dd0d443efb64cfc631ca6983f,be8a43be5bb4466a9d0573c63606524b,8884")


for i in "${list[@]}";do 
id=$(echo "$i" | awk -F ',' '{print $1}')
key=$(echo "$i" | awk -F ',' '{print $2}')
port=$(echo "$i" | awk -F ',' '{print $3}')
    echo "$id"
    ./app.tmp -id "$id" -key "$key" -p "$port" &
done;