#1 /bin/bash

cat gamename.ini | while read gamename
do
	killall -1 -v "ChatServerCenter_${gamename}"
	killall -1 -v "ChatServer_${gamename}"
done