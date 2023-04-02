package producer

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
)

var (
	config   *sarama.Config
	producer sarama.SyncProducer
)

func init() {
	CreateConfig()
}

func CreateConfig() *sarama.Config {
	config = sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Compression = sarama.CompressionSnappy
	return config
}

func CreatedProducer(options []func(*sarama.Config)) {
	if len(options) > 0 {
		for _, option := range options {
			option(config)
		}
	}

	var err error
	producer, err = sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		panic(err)
	}
}

func SendToKfk(data interface{}) (int32, int64, error) {
	CreatedProducer(nil)
	body, _ := json.Marshal(data)
	fmt.Println("Send to kafka:", string(body))
	msg := &sarama.ProducerMessage{Topic: "maxwell", Key: nil, Value: sarama.StringEncoder(body)}
	return producer.SendMessage(msg)
}

func OptionSetNewManualPartitioner(c *sarama.Config) {
	config.Producer.Partitioner = sarama.NewManualPartitioner
}

func OptionSetNewHashPartitioner(c *sarama.Config) {
	config.Producer.Partitioner = sarama.NewHashPartitioner
}

func SendToKfkWithPartition(data interface{}, partition int) (int32, int64, error) {
	options := []func(*sarama.Config){
		OptionSetNewManualPartitioner,
	}
	CreatedProducer(options)
	body, _ := json.Marshal(data)
	fmt.Println("Send to kafka:", string(body))
	msg := &sarama.ProducerMessage{Topic: "test1", Key: nil, Value: sarama.StringEncoder(body), Partition: int32(partition)}
	return producer.SendMessage(msg)
}

// SendToKfkWithKey 哈希方式提交，key相同的数据会被分配到同一个分区
func SendToKfkWithKey(data interface{}, key string) (int32, int64, error) {
	options := []func(*sarama.Config){
		OptionSetNewHashPartitioner,
	}
	CreatedProducer(options)
	body, _ := json.Marshal(data)
	fmt.Println("Send to kafka:", string(body))
	msg := &sarama.ProducerMessage{Topic: "test1", Key: sarama.StringEncoder(key), Value: sarama.StringEncoder(body)}
	return producer.SendMessage(msg)
}
