#!/bin/bash
# xcall.sh "echo 123"
#####ssh到别的主机执行命令。
command=$1
thisip=$(hostname -I | awk  '{print $1}')
for host_ip in 192.168.216.{128,129,130}
do
        
        if [ "$thisip" = "$host_ip" ]; then
        echo "#########this host $host_ip##############"
        $command
        else
        echo "#########$host_ip##############"
        ssh ${host_ip} "$command"
        fi
done
