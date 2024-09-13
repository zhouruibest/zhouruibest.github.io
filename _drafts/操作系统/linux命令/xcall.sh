#!/bin/bash
# 前提 免密登陆
# 用法 xcall.sh "echo 123"
#####ssh到别的主机执行命令。
command=$1
for host_ip in 192.168.216.{128,129,130}
do
        echo "#########$host_ip##############"
        ssh ubuntu@${host_ip} "$command"
done