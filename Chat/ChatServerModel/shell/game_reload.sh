#! /bin/bash

name1=`echo $PWD | awk -F/ '{print $5}'`
name2=`echo $PWD | awk -F/ '{print $4}'`
gamename="${name1}_${name2}"

killall -1 -v $gamename
