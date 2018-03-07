package messageLog

import (
	. "work.goproject.com/Chat/ChatServerModel/src"
)

func Save(playerId, playerName string, partnerId, serverId, serverGroupId int32, message, voice, channel, toPlayerId string) (id int, err error) {
	messageLogObj := NewMessageLog(playerId, playerName, partnerId, serverId, serverGroupId, message, voice, channel, toPlayerId)
	id, err = insert(messageLogObj)
	return
}
