package ipConfig

import (
	"fmt"

	"work.goproject.com/Chat/ChatServer/src/dal"
	"work.goproject.com/Framework/ipMgr"
	"work.goproject.com/Framework/reloadMgr"
	"work.goproject.com/goutil/debugUtil"
)

func init() {
	if err := reload(); err != nil {
		panic(fmt.Errorf("初始化IP列表失败，错误信息为：%s", err))
	}

	// 注册重新加载的方法
	reloadMgr.RegisterReloadFunc("ipConfig.reload", reload)
}

// 重新加载IP列表
func reload() error {
	var ipList []*IpConfig
	if err := dal.GetAll(&ipList); err != nil {
		return err
	}

	tmpIpList := make([]string, 0, 32)
	for _, item := range ipList {
		tmpIpList = append(tmpIpList, item.IP)
		debugUtil.Println("ip:", item.IP)
	}

	// 初始化ipMgr的数据
	ipMgr.Init(tmpIpList)

	return nil
}
