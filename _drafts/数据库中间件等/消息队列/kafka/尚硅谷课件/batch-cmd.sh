#!/bin/bash

case $1 in
"start")
    for i in host1 host2 host3            # 前提，配置免密登陆
    do
        ssh $i "/path/to/cmd -c xxx, yyy" # 绝对路径
    done
;;

"stop")

;;
esac