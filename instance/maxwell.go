package instance

import (
	"LittleBeeMark/CoolGoPkg/apply_kfk/consumer"
	"LittleBeeMark/CoolGoPkg/apply_kfk/db"
	"LittleBeeMark/CoolGoPkg/apply_kfk/model"
	"fmt"
	"time"
)

type MaxWell struct {
	KafkaConsumer *consumer.KfkConsumerGroup
}

func (mw *MaxWell) Start() {
	// 启动kafka消费者
	mw.KafkaConsumer.StartASync()

	// 模拟插入数据
	mw.KafkaConsumer.Wg.Add(1)
	go func() {
		defer mw.KafkaConsumer.Wg.Done()
		for {
			if mw.KafkaConsumer.StopCtx.Err() != nil {
				fmt.Println("insert db stop ctx err : ", mw.KafkaConsumer.StopCtx.Err())
				break
			}
			time.Sleep(500 * time.Millisecond)
			err := db.InsertStockInfo(&model.Stock{
				ProdCode:    "000001.SS",
				ProdName:    "上证指数",
				TradeStatus: "Trade",
				Last:        3000.00,
			})
			if err != nil {
				fmt.Println("insert db err : ", err)
			}
		}
	}()

	// 5秒后优雅退出
	time.Sleep(5 * time.Second)
	mw.KafkaConsumer.Quit()
}
