package producer

import (
	"fmt"
	"testing"
)

func TestSendToKfk(t *testing.T) {
	fmt.Println(13 / 2)
	a, b, err := SendToKfk("hello my name is Kafka")
	if err != nil {
		fmt.Println("has an error", err)
		return
	}
	fmt.Printf("partition is %v, offset is %v", a, b)

}
