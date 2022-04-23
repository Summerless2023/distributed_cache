package main

import (
	"main/src/models"
	"main/src/strategy"
	"main/src/utils"
	"strconv"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Info("分布式缓存服务已启动...")
	// factory := new(strategy.StrategyFactory)
	// var mytest strategy.EliminationStrategy = factory.CreateStrategy("lru")
	mytest := strategy.NewLRUCache()
	for i := 0; i < 20; i++ {
		mytest.Add(models.KeyType(strconv.Itoa(i)), models.ValueType(strconv.Itoa(i)))
	}
	utils.PrintList(*mytest.GetCacheList())
}
