package main

import (
	"flag"
	"main/conf"
	"main/src/models"
	"main/src/strategy"
	"strconv"

	"github.com/sirupsen/logrus"
)

var ToUseStrategy = flag.String("s", "lru", "调度策略")

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Info("分布式缓存服务已启动...")
	// 解析参数
	flag.Parse()
	// 输出无用参数：若命令行指定未定义参数为 -b
	logrus.Info("无用参数：", flag.Args())
	if !conf.Strategy_Map[*ToUseStrategy] {
		logrus.Fatal("不存在该策略：", *ToUseStrategy)
	}
	factory := new(strategy.StrategyFactory)
	var mytest strategy.EliminationStrategy = factory.CreateStrategy(*ToUseStrategy)
	// mytest := strategy.NewLRUCache()
	for i := 0; i < 20; i++ {
		mytest.Add(models.KeyType(strconv.Itoa(i)), models.ValueType(strconv.Itoa(i)))
	}
	// utils.PrintList(*mytest.GetCacheList())
}
