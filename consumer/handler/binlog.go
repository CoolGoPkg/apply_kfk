package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
)

type BinLogItem struct {
	Database string      `json:"database"`
	Table    string      `json:"table"`
	Type     string      `json:"type"` //insert update delete
	Ts       int64       `json:"ts"`
	Xid      int64       `json:"xid"`
	Xoffset  int64       `json:"xoffset"`
	Commit   bool        `json:"commit"`
	Data     interface{} `json:"data"`
	Old      interface{} `json:"old"`
}

func ToObject(obj interface{}, str string) {
	json.Unmarshal([]byte(str), obj)
}

type BinlogHandler struct {
}

func (b *BinlogHandler) HandleMessage(msg *sarama.ConsumerMessage) error {
	var binlogItem BinLogItem
	ToObject(&binlogItem, string(msg.Value))
	fmt.Printf("BinlogHandler HandleMessage  obj:%#v \n", binlogItem)
	return nil
}
