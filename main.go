package main

import (
	"log"
	"main/strategy"
)

func main() {
	log.Println("分布式缓存服务已启动...")
	factory := new(strategy.StrategyFactory)
	var mytest strategy.EliminationStrategy = factory.CreateStrategy("sample")
	mytest.Add("123", "123")
}
