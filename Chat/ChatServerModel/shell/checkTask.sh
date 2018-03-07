#! /bin/bash

set -e
#任务URL
token="353522F275963BF726B7F594A90D2D01"
taskUrl="http://10.255.0.1/Manage/ChatManage.ashx?ip=1.1.1.1&token=$token"

#curl 获取任务 JSON格式
content=`curl $taskUrl 2>/dev/null | jq .` || echo "curl error"
echo $content
#jq 获取 JSON 数据
Code=`echo $content | jq '.Code'`
Message=`echo $content | jq '.Message'`
Data=`echo $content | jq '.Data'`

#获取任务失败
if [[ $Code != 0 ]]; then
	echo "Error :$Message"
	exit 1
fi

#没有任务,正常退出
if [[ $Data == "" || $Data == "null" ]]; then
	exit 0
fi

#获取任务数据
TaskId=`echo $content | jq '.Data.TaskId'`
TaskName=`echo $content | jq '.Data.TaskName'`
TaskType=`echo $content | jq '.Data.TaskType'`
TaskStatus=`echo $content | jq '.Data.TaskStatus'`
GroupId=`echo $content | jq '.Data.GroupId'`
SourceUrl=`echo $content | jq '.Data.SourceUrl'`

now=$(date +"%m-%d-%Y")
echo $now
taskStatus="3"
result=""
#读取服务器信息passwd.txt(IP USER PWD)
#TaskType==1 新建任务 执行build.exp
#TaskType==2 更新任务 执行update.exp
cat passwd.txt | while read ip user passwd
do
	if [[ $TaskType == "1" ]]; then
		result=`expect build.exp $ip $user $passwd $GroupId $SourceUrl`
		if [[ $? != 0 ]]; then
			taskStatus="4"
			echo "build fault(ip GroupId) $ip $GroupId"
		fi
	elif [[ $TaskType == "2" ]]; then
		result=`expect update.exp $ip $user $passwd $GroupId $SourceUrl`
		if [[ $? != 0 ]]; then
			taskStatus="4"
			echo "update fault(ip GroupId) $ip $GroupId"
		fi
	else
		echo "TaskType error"
		exit 1
	fi
done 

#返回任务结果
#TaskStatus（3：完成 4：失败，1：重置）
postUrl="{ "Token":$token, "TaskId":$TaskId, "TaskStatus":$taskStatus, "Result":$result }"
echo $postUrl
curl -X POST -H "Content-Type: application/json" --data $postUrl $taskUrl

