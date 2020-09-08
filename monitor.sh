#!/bin/bash

while [ "1" = "1" ]
do 
    ps -efww |grep -v grep |grep vpn/chimney-go
    ch=$?
    if [ $ch -ne 0 ]; then
      nohup /home/evan/works/vpn/chimney-go -s &      
    fi

    sleep 30
done
