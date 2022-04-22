package main

import (
	"main/strategy"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.Info("分布式缓存服务已启动...")
	factory := new(strategy.StrategyFactory)
	var mytest strategy.EliminationStrategy = factory.CreateStrategy("lru")
	mytest.Add("123", "123")
}
