package main

import (
	"fmt"
	"main/strategy"
	"main/utils"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.Info("分布式缓存服务已启动...")
	// factory := new(strategy.StrategyFactory)
	// var mytest strategy.EliminationStrategy = factory.CreateStrategy("lru")
	mytest := strategy.NewLRUCache()
	mytest.Add("123", "123")
	mytest.Add("1", "1")
	mytest.Add("2", "2")
	mytest.Add("3", "3")
	utils.PrintList(*mytest.CacheList)
	fmt.Println(mytest.Get("123"))
	utils.PrintList(*mytest.CacheList)
	fmt.Println(mytest.Get("4"))
	mytest.Remove()
	mytest.Remove()
	mytest.Remove()
	mytest.Remove()
	mytest.Remove()
	utils.PrintList(*mytest.CacheList)
	mytest.Add("3", "3")
	utils.PrintList(*mytest.CacheList)
}
