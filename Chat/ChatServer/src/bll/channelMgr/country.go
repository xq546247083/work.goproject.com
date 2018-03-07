package channelMgr

import (
	. "work.goproject.com/Chat/ChatServerModel/src"
)

func init() {
	channelObj := newGroup("Country", "ChannelConfig_Country", "history_country")
	channelObj.setNotInGroupStatus(NotInCountry)
	channelObj.setNewGroupHistoryFunc(newCountryHistory)

	// 注册到集合中
	register(channelObj)
}
