package ipUtil

import (
	"encoding/json"
	"fmt"

	"work.goproject.com/goutil/securityUtil"
	"work.goproject.com/goutil/webUtil"
)

// Ip信息
type IpInfo struct {
	// 所属大陆块
	Continent string

	// 国家
	Country string

	// 省份
	Region string

	// 城市
	City string

	// 网络服务提供商
	Isp string
}

var address = "ip.work.goproject.com/query"

// 设置请求地址
// _address:请求地址
func SetAddress(_address string) {
	address = _address
}

// 查询IP信息
// IP:待查询的IP
// appId:AppId
// authCode:授权码
// 返回值:
// ipInfo:查询的IP结果
// err:错误信息
func QueryIpInfo(ip string, appId string, authCode string) (ipInfo *IpInfo, err error) {
	return QueryIpInfo2(ip, appId, authCode, -1)
}

// 查询IP信息
// IP:待查询的IP
// appId:App的Id
// appKey:应用密钥
// timeSeconds:请求的超时秒数，非正数则使用默认超时。
// 返回值:
// ipInfo:查询的IP结果
// err:错误信息
func QueryIpInfo2(ip string, appId string, appKey string, timeSeconds int) (ipInfo *IpInfo, err error) {
	signVal := securityUtil.Md5String(fmt.Sprintf("%v%v%v", appId, appKey, ip), true)
	requestUrl := fmt.Sprintf("http://%v?appid=%v&ip=%v&sign=%v", address, appId, ip, signVal)

	var resultBytes []byte
	var statusCode int = 200
	if timeSeconds > 0 {
		// 增加超时方式的请求
		transport := webUtil.NewTransport()
		transport = webUtil.GetTimeoutTransport(transport, timeSeconds)
		statusCode, resultBytes, err = webUtil.GetWebData2(requestUrl, nil, nil, transport)
	} else {
		// 直接请求
		resultBytes, err = webUtil.GetWebData(requestUrl, nil)
	}

	if err != nil {
		// logUtil.ErrorLog("获取IP信息失败，错误信息:%v", err.Error())
		return
	}

	if statusCode != 200 {
		err = fmt.Errorf("http response error, code:%v", statusCode)
		return
	}

	// 反序列化
	mapData := make(map[string]interface{})
	if err = json.Unmarshal(resultBytes, &mapData); err != nil {
		err = fmt.Errorf("数据反序列化失败，url:%v data:%v error:%v", requestUrl, string(resultBytes), err.Error())
		return
	}

	// 检查结果
	if code, exist := mapData["Code"]; exist == false {
		err = fmt.Errorf("应答的数据格式不正确: url:%v data:%v", requestUrl, string(resultBytes))
		return
	} else if tmpCode, ok := code.(float64); ok == false || tmpCode != 0 {
		err = fmt.Errorf("应答的数据结果不正确: url:%v data:%v", requestUrl, string(resultBytes))
		return
	}

	ipInfo = &IpInfo{}
	resultData := mapData["Data"].(map[string]interface{})
	ipInfo.Continent = resultData["Continent"].(string)
	ipInfo.Country = resultData["Country"].(string)
	ipInfo.Region = resultData["Region"].(string)
	ipInfo.City = resultData["City"].(string)
	ipInfo.Isp = resultData["Isp"].(string)

	return
}
