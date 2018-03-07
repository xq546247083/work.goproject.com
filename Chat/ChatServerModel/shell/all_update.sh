#!/bin/bash
#所有聊天服务器的更新脚本

# 1、检测更新文件是否存在
if [ ! -x "ChatServerCenter" ];then
	echo "ChatServerCenter不存在或不具有执行权限"
	exit 1
fi

if [ ! -x "ChatServer" ];then
	echo "ChatServer不存在或不具有执行权限"
	exit 1
fi


# 2、对每一个游戏进行更新处理
#ls -l | grep "^d" | awk '{print $9}' | while read gamename
cat gamename.ini | while read gamename
do
	#Step1:Stop
	#进入ChatServerCenter目录
	cd $PWD/$gamename/ChatServerCenter
	./stop.sh
	
	#进入ChatServer目录
	cd ../ChatServer
	./stop.sh

	#回到最外层目录
	cd ../../

	#Step2:Copy File
	cp ChatServerCenter $PWD/$gamename/ChatServerCenter
	cp ChatServer $PWD/$gamename/ChatServer

	#Step3:Rename File
	#进入ChatServerCenter目录
	cd $PWD/$gamename/ChatServerCenter
	mv ChatServerCenter ChatServerCenter_${gamename}

	#进入ChatServer目录
	cd ../ChatServer
	mv ChatServer ChatServer_${gamename}

	#回到最外层目录
	cd ../../

	#Step4:Start
	#进入ChatServerCenter目录
	cd $PWD/$gamename/ChatServerCenter
	./start.sh
	sleep 1

	#进入ChatServer目录
	cd ../ChatServer
	./start.sh
	sleep 1

	#回到最外层目录
	cd ../../	

	echo $gamename "更新完成"
done