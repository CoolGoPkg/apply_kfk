package consumer

import (
	"LittleBeeMark/CoolGoPkg/apply_kfk/conf"
	"LittleBeeMark/CoolGoPkg/apply_kfk/consumer/handler"
	"testing"
)

func TestNewKfkConsumerGroupStartASync(t *testing.T) {
	hold := make(chan struct{})
	conf.InitConfig("../conf/config.yaml")
	got := NewKfkConsumerGroup(&conf.Config.KafkaGroupConf)
	got.AddHandler(&DefaultMessageHandler{})
	got.StartASync()

	<-hold
}

func TestNewKfkConsumerGroupStartSync(t *testing.T) {
	conf.InitConfig("../conf/config.yaml")
	got := NewKfkConsumerGroup(&conf.Config.KafkaGroupConf)
	got.AddHandler(&handler.BinlogHandler{})
	got.StartSync()

}
