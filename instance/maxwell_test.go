package instance

import (
	"LittleBeeMark/CoolGoPkg/apply_kfk/conf"
	"LittleBeeMark/CoolGoPkg/apply_kfk/consumer"
	"LittleBeeMark/CoolGoPkg/apply_kfk/consumer/handler"
	"LittleBeeMark/CoolGoPkg/apply_kfk/db"
	"testing"
)

func TestMaxWell_Start(t *testing.T) {
	// 初始化配置
	conf.InitConfig("../conf/config.yaml")
	db.InitDB()

	mw := &MaxWell{}
	// 初始化kafka消费者
	mw.KafkaConsumer = consumer.NewKfkConsumerGroup(&conf.Config.KafkaGroupConf)
	mw.KafkaConsumer.AddHandler(&handler.BinlogHandler{})
	mw.Start()
}
