#!/bin/bash
#将ChatServerCenter、ChatServer上传到10.1.0.21上
ip=10.1.0.21
user=root
password=MOQIKAKA_redis

scp ChatServerCenter $user@$ip:/home/Chat/
echo "Upload ChatServerCenter successfully"

scp ChatServer $user@$ip:/home/Chat/
echo "Upload ChatServer successfully"