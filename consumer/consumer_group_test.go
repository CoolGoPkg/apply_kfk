package consumer

import (
	"LittleBeeMark/CoolGoPkg/apply_kfk/conf"
	"testing"
)

func TestNewKfkConsumerGroup(t *testing.T) {
	conf.InitConfig("")
	got := NewKfkConsumerGroup(&conf.Config.KafkaGroupConf)
	got.Start()
}
