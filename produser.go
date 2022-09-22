package main

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
)

var (
	producer sarama.SyncProducer
)

func CreatedProducer() {
	config := sarama.NewConfig()
	// 发送成功配置
	config.Producer.Return.Successes = true
	config.Producer.Compression = sarama.CompressionSnappy
	var err error
	producer, err = sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		panic(err)
	}
}

func SendToKfk(data interface{}) (int32, int64, error) {
	CreatedProducer()
	body, _ := json.Marshal(data)
	fmt.Println("Send to kafka:", string(body))
	msg := &sarama.ProducerMessage{Topic: "test1", Key: nil, Value: sarama.StringEncoder(body)}
	return producer.SendMessage(msg)
}
//config := sarama.NewConfig()
//config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
//config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
//config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
//
//// 构造一个消息
//msg := &sarama.ProducerMessage{}
//msg.Topic = "web_log"
//msg.Value = sarama.StringEncoder("this is a test log")
//// 连接kafka
//client, err := sarama.NewSyncProducer([]string{"192.168.1.7:9092"}, config)
//if err != nil {
//fmt.Println("producer closed, err:", err)
//return
//}
//defer client.Close()
//// 发送消息
//pid, offset, err := client.SendMessage(msg)
//if err != nil {
//fmt.Println("send msg failed, err:", err)
//return
//}
//fmt.Printf("pid:%v offset:%v\n", pid, offset)