//package test
package main

import (
	"fmt"
)

var (
	con_player_tableName = "player"
)

func init() {
	registerSyncObj(con_player_tableName)
}

func insert(obj *player) {
	command := fmt.Sprintf("INSERT INTO `%s` (`Id`,`Name`) VALUES ('%v','%v') ", con_player_tableName, obj.Id, obj.Name)
	save(con_player_tableName, command)
}

func update(obj *player) {
	command := fmt.Sprintf("UPDATE `%s` SET  `Name` = '%v' WHERE `Id` = '%v';", con_player_tableName, obj.Name, obj.Id)
	save(con_player_tableName, command)
}

func clear(obj *player) {
	command := fmt.Sprintf("DELETE FROM %s where Id = '%v';", con_player_tableName, obj.Id)
	save(con_player_tableName, command)
}
