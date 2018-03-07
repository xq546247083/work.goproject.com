#!/bin/bash
#将ChatServerCenter、ChatServer上传到120.131.9.117上
ip=120.131.9.117
user=root
password=MOQIkaka$#@!1234

scp ChatServerCenter $user@$ip:/home/Chat/
echo "Upload ChatServerCenter successfully"

scp ChatServer $user@$ip:/home/Chat/
echo "Upload ChatServer successfully"