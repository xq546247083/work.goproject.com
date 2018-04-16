package kafkaMgr

import "strings"

// Kafka配置对象
type KafkaConfig struct {
	// Brokers的地址
	Brokers string

	// 主题
	Topics string

	// 分区
	Partitions string

	// 分组Id
	GroupId string

	// 用户名
	UserName string

	// 密码
	Passward string

	// 需要的证书文件
	CertFile string
}

func (this *KafkaConfig) GetBrokerList() []string {
	if len(this.Brokers) <= 0 {
		return make([]string, 0)
	}

	return strings.Split(this.Brokers, ",")
}

func (this *KafkaConfig) GetTopicList() []string {
	if len(this.Topics) <= 0 {
		return make([]string, 0)
	}

	return strings.Split(this.Topics, ",")
}

func (this *KafkaConfig) GetPartitionList() []string {
	if len(this.Partitions) <= 0 {
		return make([]string, 0)
	}

	return strings.Split(this.Partitions, ",")
}
