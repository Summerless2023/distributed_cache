package main

import (
	"main/src/concurrency"
	"main/src/models"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func add(concurrencyCache *concurrency.ConcurrencyCache, key string, value string, wg *sync.WaitGroup) {
	concurrencyCache.Add(models.KeyType(key), models.ValueType(value))
	wg.Done()
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Info("分布式缓存服务已启动...")
	// // factory := new(strategy.StrategyFactory)
	// // var mytest strategy.EliminationStrategy = factory.CreateStrategy("lru")
	// mytest := strategy.NewFIFOStrategy()
	// for i := 0; i < 100; i++ {
	// 	mytest.Add(models.KeyType(strconv.Itoa(i)), models.ValueType(strconv.Itoa(i)))
	// }
	// utils.PrintList(*mytest.GetCacheList())
	concurrencyCache := &concurrency.ConcurrencyCache{}
	var waitGroup sync.WaitGroup
	for i := 0; i < 1000; i++ {
		waitGroup.Add(1)
		go add(concurrencyCache, strconv.Itoa(i), strconv.Itoa(i), &waitGroup)
	}
	waitGroup.Wait()
	time.Sleep(3 * time.Second)
	// logrus.Info(concurrencyCache.Get(models.KeyType(strconv.Itoa(99))))
	cnt := 0
	for i := 0; i < 1000; i++ {
		_, ok := concurrencyCache.Get(models.KeyType(strconv.Itoa(i)))
		if ok {
			cnt++
		}
	}
	logrus.Debug(cnt)
}
