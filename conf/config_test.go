package conf

import (
	"fmt"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	LoadConfig(Config, "./config.yaml")
	fmt.Printf("conf.Config is %#v", Config)
}
