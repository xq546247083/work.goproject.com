package zooKeeperMgr

import (
	"strings"
	"time"
)

// ZooKeeper配置对象
type ZooKeeperConfig struct {
	// ZooKeeper的地址(如果有多个地址，则用,分隔)
	Address string

	// 父节点名称
	ParentPath string

	// 会话超时时间（单位：秒）
	SessionTimeout time.Duration
}

// 获取所有的ZooKeeper服务器列表
func (this *ZooKeeperConfig) GetAddressList() []string {
	return strings.Split(this.Address, ",")
}
