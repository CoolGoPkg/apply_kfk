package main

import "LittleBeeMark/CoolGoPkg/apply_kfk/consumer"

func main() {
	holdChan := make(chan struct{})
	for i := 0; i < 4; i++ {
		go func() {
			kafkaConsumer := new(consumer.KfkConsumer)
			kafkaConsumer.Init([]string{"127.0.0.1:9092"})
			cs := consumer.DefaultConsumer{}
			kafkaConsumer.AddHandler(cs)
			kafkaConsumer.Start("test1")
		}()

	}
	<-holdChan
}
