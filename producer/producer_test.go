package producer

import (
	"fmt"
	"testing"
)

func TestSendToKfk(t *testing.T) {
	a, b, err := SendToKfk("hello my name is mark14")
	if err != nil {
		fmt.Println("has an error", err)
		return
	}
	fmt.Printf("partition is %v, offset is %v", a, b)

}
