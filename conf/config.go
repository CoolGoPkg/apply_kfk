package conf

import (
	"fmt"
	"github.com/jinzhu/configor"
	"log"
)

var Config = &CustomerConfig{}

// InitConfig doc
func InitConfig(file string) *CustomerConfig {
	LoadConfig(Config, file)
	fmt.Println(file)
	return Config
}

func LoadConfig(dest interface{}, path string) {
	// 	AutoReload         bool 是否热更新
	//	AutoReloadInterval time.Duration 几秒进行一次热更新

	// 初始化默认 config 的代码
	//if os.Getenv("CONFIGOR_DEBUG_MODE") != "" {
	//	config.Debug = true
	//}
	//
	//if os.Getenv("CONFIGOR_VERBOSE_MODE") != "" {
	//	config.Verbose = true
	//}
	//
	//if os.Getenv("CONFIGOR_SILENT_MODE") != "" {
	//	config.Silent = true
	//}
	//
	//if config.AutoReload && config.AutoReloadInterval == 0 {
	//	config.AutoReloadInterval = time.Second
	//}
	// 支持 tag = default 赋值默认值

	// 首先使用config.Environment 不去赋值config.Environment 那么取环境变量 CONFIGOR_ENV,否则取获取第一个命令，否则就返回 development
	//var testRegexp = regexp.MustCompile("_test|(\\.test$)")
	//
	//// GetEnvironment get environment
	//func (configor *Configor) GetEnvironment() string {
	//	if configor.Environment == "" {
	//	if env := os.Getenv("CONFIGOR_ENV"); env != "" {
	//	return env
	//}
	//
	//	if testRegexp.MatchString(os.Args[0]) {
	//	return "test"
	//}
	//
	//	return "development"
	//}
	//	return configor.Environment
	//}

	// 假设是 如果什么都没有设置那么 那么就会去找 a.development.yaml
	//func getConfigurationFileWithENVPrefix(file, env string) (string, time.Time, error) {
	//	var (
	//		envFile string
	//		extname = path.Ext(file)
	//	)
	//
	//	if extname == "" {
	//		envFile = fmt.Sprintf("%v.%v", file, env)
	//	} else {
	//		envFile = fmt.Sprintf("%v.%v%v", strings.TrimSuffix(file, extname), env, extname)
	//	}
	//
	//	if fileInfo, err := os.Stat(envFile); err == nil && fileInfo.Mode().IsRegular() {
	//		return envFile, fileInfo.ModTime(), nil
	//	}
	//	return "", time.Now(), fmt.Errorf("failed to find file %v", file)
	//}

	// 如果CONFIGOR_ENV_PREFIX 为空那么就会去找环境变量 CONFIGOR_ENV_PREFIX，否则就是 Configor
	//func (configor *Configor) getENVPrefix(config interface{}) string {
	//if configor.Config.ENVPrefix == "" {
	//if prefix := os.Getenv("CONFIGOR_ENV_PREFIX"); prefix != "" {
	//return prefix
	//}
	//return "Configor"
	//}
	//return configor.Config.ENVPrefix
	//}

	// 尝试取tag为env的key的值，没有就取（前缀_字段名）的环境变量，取到了就去赋值给字段
	//if prefix := configor.getENVPrefix(config); prefix == "-" {
	//	err = configor.processTags(config)
	//} else {
	//	err = configor.processTags(config, prefix)
	//}
	if err := configor.Load(dest, path); err != nil {
		log.Panicf("failed to load local config file: %v", err)
	}
}

type KafkaConfig struct {
	Brokers     []string `yaml:"brokers"`
	Version     string   `yaml:"version"`
	TickerTopic string   `yaml:"ticker_topic"`
}

type KafkaGroupConfig struct {
	Brokers  []string `yaml:"brokers"`
	Group    string   `yaml:"group"`
	Version  string   `yaml:"version"`
	Assignor string   `yaml:"assignor"`
	Oldest   bool     `yaml:"oldest"`
	Verbose  bool     `yaml:"verbose"`
	Topics   []string `yaml:"topics"`
}

type CustomerConfig struct {
	KafkaConf      KafkaConfig      `yaml:"kafka"`
	KafkaGroupConf KafkaGroupConfig `yaml:"kafka_group"`
}
