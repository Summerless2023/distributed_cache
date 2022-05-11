package conf

import (
	"sync"

	"github.com/go-ini/ini"
)

type ConfigStruct struct {
	Max_Bytes   int64  `ini:"Max_Bytes"`
	Strategy    string `ini:"Strategy"`
	UpdateTime  int64  `ini:"UpdateTime"`
	ExpiredTime int64  `ini:"ExpiredTime"`
}

var Config = &ConfigStruct{} //定义存储配置信息的全局变量
var once sync.Once

func LoadConfig() {
	once.Do(func() {

		errConfMsg := ini.MapTo(Config, "./conf/conf.ini") //先读第一个文件
		if errConfMsg == nil {
			CheckParameter()
		}

		var ConfigBack = &ConfigStruct{}
		errBackMsg := ini.MapTo(ConfigBack, "./conf/back_conf.ini") //再读第二个文件
		if errBackMsg == nil {
			ModifyParameter(ConfigBack)
			CheckParameter()
		}

		//读第三个文件，一定存在
		var ConfigDefault = &ConfigStruct{}
		ConfigDefault.Max_Bytes = DEFAULT_MAX_BYTES
		ConfigDefault.Strategy = DEFAULT_STRATEGY
		ModifyParameter(ConfigDefault)
		CheckParameter()
	})
}

func CheckParameter() {
	if Config.Max_Bytes > DEFAULT_MAX_BYTES || Config.Max_Bytes <= 0 {
		Config.Max_Bytes = 0
	}
	if Config.Strategy != "lru" && Config.Strategy != "lfu" && Config.Strategy != "fifo" {
		Config.Strategy = ""
	}
}

func ModifyParameter(ConfigTemp *ConfigStruct) {
	if Config.Max_Bytes == 0 {
		Config.Max_Bytes = ConfigTemp.Max_Bytes
	}
	if Config.Strategy == "" {
		Config.Strategy = ConfigTemp.Strategy
	}
}
