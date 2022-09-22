package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

type partitionRetry struct {
	topic     string
	partition int32
}

type MessageHandler interface {
	HandleMessage(*sarama.ConsumerMessage) error
}

type KfkConsumer struct {
	brokerList     []string
	consumer       sarama.Consumer
	pcList         []sarama.PartitionConsumer
	errorRetry     chan *partitionRetry
	messageHandler MessageHandler
}

func (self *KfkConsumer) AddHandler(handler MessageHandler) {
	self.messageHandler = handler
}

func (self *KfkConsumer) Init(brokerList []string) {
	self.errorRetry = make(chan *partitionRetry, 100)
	self.brokerList = brokerList
	go self.checkTopic()
}

func (self *KfkConsumer) checkTopic() {
	for {
		select {
		case retry := <-self.errorRetry:
			fmt.Println("CheckTopic retry", retry)
			self.eventLoop(retry.topic, retry.partition)
		}
		time.Sleep(time.Second)
	}
}

func (self *KfkConsumer) eventLoop(topic string, partition int32) {
	pc, err := self.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		fmt.Println("eventLoop Error:", err)
		self.errorRetry <- &partitionRetry{
			topic:     topic,
			partition: partition,
		}
	} else {
		self.pcList = append(self.pcList, pc)
		go self.Handle(pc)
	}
}

func (self *KfkConsumer) Handle(pc sarama.PartitionConsumer) {
	fmt.Println("KfkConsumer start Handle:")
	for msg := range pc.Messages() {
		if self.messageHandler != nil {
			self.messageHandler.HandleMessage(msg)
		}
	}
}

func (self *KfkConsumer) Close() {
	for _, pc := range self.pcList {
		pc.Close()
	}
}

func (self *KfkConsumer) Start(topic string) {
	self.pcList = make([]sarama.PartitionConsumer, 0)
	config := sarama.NewConfig()
	var err error
	self.consumer, err = sarama.NewConsumer(self.brokerList, config)
	if err != nil {
		fmt.Println("Failed to start Sarama NewConsumer:", err)
		panic(err)
	} else {
		fmt.Println("KfkConsumer success:", self.brokerList, topic)
	}
	if partitionList, err1 := self.consumer.Partitions(topic); err1 == nil {
		fmt.Println("KfkConsumer Partitions success:", partitionList, self.consumer, topic)
		for _, partition := range partitionList {
			self.eventLoop(topic, partition)
		}
	} else {
		fmt.Println("Failed to start Sarama Partitions:", partitionList, err1)
		panic(err)
	}
}



type DefaultConsumer struct {}


func(dc DefaultConsumer)HandleMessage(msg *sarama.ConsumerMessage)error{
	fmt.Printf("Partition:%d Offset:%d Key:%s Value:%s \n", msg.Partition, msg.Offset, msg.Key, msg.Value)
	return nil
}

func Consume(){
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions("web_log") // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println(partitionList)
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition("web_log", int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		defer pc.AsyncClose()
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v \n", msg.Partition, msg.Offset, msg.Key, msg.Value)
			}
		}(pc)
	}
}