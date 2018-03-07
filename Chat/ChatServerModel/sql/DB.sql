/*
SQLyog 企业版 - MySQL GUI v7.14 
MySQL - 5.5.27 : Database - chatserver_v5.0
*********************************************************************
*/


/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;

/*Table structure for table `config` */

DROP TABLE IF EXISTS `config`;

CREATE TABLE `config` (
  `ConfigKey` varchar(64) NOT NULL COMMENT '配置Key',
  `ConfigValue` varchar(1024) NOT NULL COMMENT '配置值',
  `ConfigDesc` varchar(64) DEFAULT NULL COMMENT '配置描述',
  PRIMARY KEY (`ConfigKey`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Table structure for table `config_ip` */

DROP TABLE IF EXISTS `config_ip`;

CREATE TABLE `config_ip` (
  `IP` varchar(128) NOT NULL COMMENT 'IP',
  PRIMARY KEY (`IP`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Table structure for table `config_word_forbid` */

DROP TABLE IF EXISTS `config_word_forbid`;

CREATE TABLE `config_word_forbid` (
  `Word` varchar(32) NOT NULL COMMENT '单词',
  PRIMARY KEY (`Word`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Table structure for table `config_word_sensitive` */

DROP TABLE IF EXISTS `config_word_sensitive`;

CREATE TABLE `config_word_sensitive` (
  `Word` varchar(100) NOT NULL COMMENT '单词',
  PRIMARY KEY (`Word`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Table structure for table `history_country` */

DROP TABLE IF EXISTS `history_country`;

CREATE TABLE `history_country` (
  `Id` int(11) NOT NULL,
  `Identifier` varchar(64) NOT NULL COMMENT '分组消息的唯一标识',
  `Channel` varchar(32) NOT NULL COMMENT '聊天频道',
  `Message` varchar(128) NOT NULL COMMENT '聊天消息',
  `Voice` varchar(512) NOT NULL COMMENT '语音信息',
  `FromPlayer` varchar(5120) NOT NULL COMMENT '说话的源玩家',
  `FromPlayerId` varchar(64) NOT NULL COMMENT '说话的源玩家Id',
  `Crtime` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`Id`),
  KEY `IX_ServerGroupId_GroupName` (`Identifier`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Table structure for table `history_crossserver` */

DROP TABLE IF EXISTS `history_crossserver`;

CREATE TABLE `history_crossserver` (
  `Id` int(11) NOT NULL,
  `Channel` varchar(32) NOT NULL COMMENT '聊天渠道',
  `Message` varchar(128) NOT NULL COMMENT '聊天消息',
  `Voice` varchar(512) NOT NULL COMMENT '语音信息',
  `FromPlayer` varchar(5120) NOT NULL COMMENT '说话的源玩家',
  `FromPlayerId` varchar(64) NOT NULL COMMENT '说话的源玩家Id',
  `Crtime` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Table structure for table `history_private` */

DROP TABLE IF EXISTS `history_private`;

CREATE TABLE `history_private` (
  `Id` int(11) NOT NULL,
  `PlayerId` varchar(64) NOT NULL COMMENT '私聊消息的接收者Id',
  `Channel` varchar(32) NOT NULL COMMENT '聊天频道',
  `Message` varchar(128) NOT NULL COMMENT '聊天消息',
  `Voice` varchar(512) NOT NULL COMMENT '语音信息',
  `FromPlayer` varchar(5120) NOT NULL COMMENT '源玩家',
  `FromPlayerId` varchar(64) NOT NULL COMMENT '源玩家Id',
  `ToPlayer` varchar(5120) NOT NULL COMMENT '目标玩家对象',
  `ToPlayerId` varchar(64) NOT NULL COMMENT '目标玩家Id',
  `Crtime` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`Id`,`PlayerId`),
  KEY `IX_PlayerId` (`PlayerId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Table structure for table `history_world` */

DROP TABLE IF EXISTS `history_world`;

CREATE TABLE `history_world` (
  `Id` int(11) NOT NULL,
  `ServerGroupId` int(11) NOT NULL COMMENT '服务器组Id',
  `Channel` varchar(32) NOT NULL COMMENT '聊天渠道',
  `Message` varchar(128) NOT NULL COMMENT '聊天消息',
  `Voice` varchar(512) NOT NULL COMMENT '语音信息',
  `FromPlayer` varchar(5120) NOT NULL COMMENT '说话的源玩家',
  `FromPlayerId` varchar(64) NOT NULL COMMENT '说话的源玩家Id',
  `Crtime` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`Id`),
  KEY `IX_ServerGroupId` (`ServerGroupId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Table structure for table `log_message` */

DROP TABLE IF EXISTS `log_message`;

CREATE TABLE `log_message` (
  `Id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `PlayerId` varchar(36) NOT NULL COMMENT '玩家Id',
  `PlayerName` varchar(32) NOT NULL COMMENT '玩家名称',
  `PartnerId` int(11) NOT NULL COMMENT '合作商Id',
  `ServerId` int(11) NOT NULL COMMENT '服务器Id',
  `ServerGroupId` int(11) NOT NULL COMMENT '服务器组Id',
  `Message` varchar(128) NOT NULL COMMENT '聊天内容',
  `Voice` varchar(512) NOT NULL DEFAULT '' COMMENT '语音信息',
  `Channel` varchar(32) NOT NULL COMMENT '聊天频道',
  `ToPlayerId` varchar(36) DEFAULT NULL COMMENT '如果是私聊，表示目标玩家',
  `Crtime` datetime NOT NULL COMMENT '发送时间',
  PRIMARY KEY (`Id`),
  KEY `IX_ServerGroupId_Crtime` (`ServerGroupId`,`Crtime`),
  KEY `IX_PlayerId_Crtime` (`PlayerId`,`Crtime`),
  KEY `IX_Name_Crtime` (`PlayerName`,`Crtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Table structure for table `log_online` */

DROP TABLE IF EXISTS `log_online`;

CREATE TABLE `log_online` (
  `OnlineTime` datetime NOT NULL COMMENT '在线时间',
  `Sid` int(11) NOT NULL COMMENT '序号',
  `ServerAddress` varchar(64) NOT NULL COMMENT '服务器Id',
  `ClientCount` int(11) NOT NULL COMMENT '客户端数量',
  `PlayerCount` int(11) NOT NULL COMMENT '玩家数量',
  `TotalCount` int(11) NOT NULL COMMENT '所有服务器的总数量',
  PRIMARY KEY (`OnlineTime`,`Sid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Table structure for table `player` */

DROP TABLE IF EXISTS `player`;

CREATE TABLE `player` (
  `Id` varchar(64) NOT NULL COMMENT '玩家Id',
  `PartnerId` int(11) NOT NULL COMMENT '合作商Id',
  `ServerId` int(11) NOT NULL COMMENT '服务器Id',
  `Name` varchar(64) NOT NULL COMMENT '玩家名称',
  `Lv` int(11) NOT NULL COMMENT '玩家等级',
  `Vip` int(11) NOT NULL COMMENT '玩家Vip等级',
  `Token` varchar(64) NOT NULL COMMENT '玩家登录令牌',
  `ExtendInfo` varchar(1024) NOT NULL COMMENT '扩展信息',
  `RegisterTime` datetime NOT NULL COMMENT '注册时间',
  `SilentEndTime` datetime NOT NULL COMMENT '禁言结束时间',
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
