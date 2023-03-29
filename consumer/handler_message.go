package consumer

import (
	"fmt"
	"github.com/Shopify/sarama"
)

type MessageHandler interface {
	HandleMessage(*sarama.ConsumerMessage) error
}

type DefaultMessageHandler struct {
	Name string
}

func NewDefaultMessageHandler() *DefaultMessageHandler {
	return &DefaultMessageHandler{}
}

func (dfh *DefaultMessageHandler) HandleMessage(msg *sarama.ConsumerMessage) error {
	fmt.Printf("DefaultMessageHandler HandleMessage key: %s , val: %s\n", string(msg.Key), string(msg.Value))
	return nil
}
