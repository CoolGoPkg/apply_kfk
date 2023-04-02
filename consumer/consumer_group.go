package consumer

import (
	"LittleBeeMark/CoolGoPkg/apply_kfk/conf"
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"sync"
)

// KfkConsumerGroup doc
type KfkConsumerGroup struct {
	Brokers  []string
	Topics   []string
	Version  string
	Group    string
	Assignor string
	Oldest   bool
	Verbose  bool

	ready       chan bool
	StopCtx     context.Context
	StopCancel  context.CancelFunc
	Wg          *sync.WaitGroup
	groupClient sarama.ConsumerGroup

	Handler MessageHandler
}

func NewKfkConsumerGroup(conf *conf.KafkaGroupConfig) *KfkConsumerGroup {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	return &KfkConsumerGroup{
		Brokers:    conf.Brokers,
		Topics:     conf.Topics,
		Version:    conf.Version,
		Group:      conf.Group,
		Assignor:   conf.Assignor,
		Oldest:     conf.Oldest,
		Verbose:    conf.Verbose,
		ready:      make(chan bool),
		StopCtx:    ctx,
		StopCancel: cancel,
		Wg:         wg,
	}
}

// AddHandler 注入处理者
func (cg *KfkConsumerGroup) AddHandler(handler MessageHandler) {
	cg.Handler = handler
}

// AddClient 注入Client
func (cg *KfkConsumerGroup) AddClient(client sarama.ConsumerGroup) {
	cg.groupClient = client
}

func (cg *KfkConsumerGroup) PrepareConfig() {
	log.Println("Starting a new Sarama consumer")

	if cg.Verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	version, err := sarama.ParseKafkaVersion(cg.Version)
	if err != nil {
		panic("Error parsing Kafka version: " + err.Error())
		return
	}
	fmt.Println("parse version : ", version)

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = false
	//config.Consumer.Offsets.AutoCommit.Interval = 3*time.Second
	//config.Consumer.Offsets.AutoCommit.Enable = false
	config.Version = version
	switch cg.Assignor {
	// 尝试将分区保持在同一消费者上
	// Example with topic T with six partitions (0..5) and two members (M1, M2):
	//
	//	M1: {T: [0, 2, 4]}
	//	M2: {T: [1, 3, 5]}
	//
	// On reassignment with an additional consumer, you might get an assignment plan like:
	//
	//	M1: {T: [0, 2]}
	//	M2: {T: [1, 3]}
	//	M3: {T: [4, 5]}
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategySticky}
	// 交替分配分区
	// For example, there are two topics (t0, t1) and two consumer (m0, m1), and each topic has three partitions (p0, p1, p2):
	// M0: [t0p0, t0p2, t1p1]
	// M1: [t0p1, t1p0, t1p2]
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin}
	// 循环分配分区
	// Example with two topics T1 and T2 with six partitions each (0..5) and two members (M1, M2):
	//
	//	M1: {T1: [0, 1, 2], T2: [0, 1, 2]}
	//	M2: {T2: [3, 4, 5], T2: [3, 4, 5]}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRange}
	default:
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRange}
	}

	if cg.Oldest {
		// OffsetNewest stands for the log head offset, i.e. the offset that will be
		// assigned to the next message that will be produced to the partition. You
		// can send this to a client's GetOffset method to get this offset, or when
		// calling ConsumePartition to start consuming new messages.
		//OffsetNewest int64 = -1
		// OffsetOldest stands for the oldest offset available on the broker for a
		// partition. You can send this to a client's GetOffset method to get this
		// offset, or when calling ConsumePartition to start consuming from the
		// oldest offset that is still available on the broker.
		//	OffsetOldest int64 = -2
		// 使用offsetNewest，从最新的offset开始消费，就是使用最新消息的偏移量，这个抛弃了进入消息队列但未消费的消息，不推荐使用，系统挂了就会丢失消息
		// 使用offsetOldest，从最旧的offset开始消费，就是使还未消费的消息的偏移量
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	// 创建client
	newClient, err := sarama.NewClient(cg.Brokers, config)
	if err != nil {
		panic("consumer group newClient err : " + err.Error())
		return
	}
	// 获取所有的topic
	//topics, err := newClient.Topics()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("topics", topics)

	client, err := sarama.NewConsumerGroupFromClient(cg.Group, newClient)
	if err != nil {
		panic("Error creating consumer group client: " + err.Error())
		return
	}
	cg.AddClient(client)
}

// StartASync 异步启动
func (cg *KfkConsumerGroup) StartASync() {
	cg.PrepareConfig()
	cg.Wg.Add(1)
	go func() {
		defer cg.Wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := cg.groupClient.Consume(cg.StopCtx, cg.Topics, cg); err != nil {
				fmt.Println("Error from consumer: ", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if cg.StopCtx.Err() != nil {
				log.Println("ctx error", cg.StopCtx.Err())
				return
			}
			cg.ready = make(chan bool)
		}
	}()

	<-cg.ready
	fmt.Println("Sarama consumer up and running!...")
}

// StartSync 同步启动
func (cg *KfkConsumerGroup) StartSync() {
	cg.PrepareConfig()
	for {
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		if err := cg.groupClient.Consume(cg.StopCtx, cg.Topics, cg); err != nil {
			fmt.Println("Error from consumer: ", err)
		}
		// check if context was cancelled, signaling that the consumer should stop
		if cg.StopCtx.Err() != nil {
			log.Println("ctx error", cg.StopCtx.Err())
			return
		}
	}
}

// Quit 退出
func (cg *KfkConsumerGroup) Quit() {
	cg.StopCancel()
	cg.Wg.Wait()
	err := cg.groupClient.Close()
	if err != nil {
		fmt.Println("kfk consumer group Error closing client: ", err)
		return
	}
	fmt.Printf("kfk topic:%v consumer has quit  ..... \n", cg.Topics)
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (cg *KfkConsumerGroup) Setup(sarama.ConsumerGroupSession) error {
	fmt.Println("Consumer Group Setup .....", cg.Topics)
	close(cg.ready)
	// Mark the consumer as ready
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (cg *KfkConsumerGroup) Cleanup(sarama.ConsumerGroupSession) error {
	fmt.Println("Consumer Group Cleanup .....", cg.Topics)
	return nil
}

//var consumerCount int

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (cg *KfkConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message := <-claim.Messages():
			cg.Handler.HandleMessage(message)
			// 更新Offset
			session.MarkMessage(message, "")

			//consumerCount++
			//if consumerCount%3 == 0 { // 假设每消费 3 条数据 commit 一次
			//	// 手动提交，不能频繁调用
			//	t1 := time.Now().Nanosecond()
			//	session.Commit()
			//	t2 := time.Now().Nanosecond()
			//	fmt.Println("commit cost:", (t2-t1)/(1000*1000), "ms")
			//}
		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}
