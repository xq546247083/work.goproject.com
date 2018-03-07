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
}

func (this *KafkaConfig) GetBrokerList() []string {
	return strings.Split(this.Brokers, ",")
}

func (this *KafkaConfig) GetTopicList() []string {
	return strings.Split(this.Topics, ",")
}

func (this *KafkaConfig) GetPartitionList() []string {
	return strings.Split(this.Partitions, ",")
}
