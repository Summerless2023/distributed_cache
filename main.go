package main

import (
	"log"

	"./strategy"
)

func main() {
	log.Println("分布式缓存服务已启动...")
	var mytest *eliminationstrategy.SampleStrategy
	mytest = eliminationstrategy.New()
	mytest.Add("123", "123")
}
