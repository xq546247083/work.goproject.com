#! /bin/bash

name1=`echo $PWD | awk -F/ '{print $5}'`
name2=`echo $PWD | awk -F/ '{print $4}'`
gamename="${name1}_${name2}"

kill $(ps aux | grep ${gamename}$ | grep daemon.sh | grep -v grep | awk '{print $2}') 2>/dev/null
killall -15 $gamename