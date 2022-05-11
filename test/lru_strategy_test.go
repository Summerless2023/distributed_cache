package test

import (
	"main/conf"
	"main/src/models"
	"main/src/strategies"
	"main/src/utils"
	"strconv"
	"testing"
	"time"
)

func TestLRUStrategy(t *testing.T) {
	factory := new(strategies.StrategyFactory)
	var mytest strategies.EliminationStrategy = factory.CreateStrategy("lru")
	conf.LoadConfig()
	expiretime := conf.Config.ExpiredTime
	mytest.Add("1", "1", expiretime)
	mytest.Add("2", "2", expiretime)
	mytest.Add("3", "3", expiretime)
	// for i := mytest.CacheList.Front(); i != nil; i = mytest.CacheList.Next() {
	// 	t.Log(i.Value)
	// }

	t.Log("测试完成")
}

func TestLRUStrategy1(t *testing.T) {
	factory := new(strategies.StrategyFactory)
	var mytest strategies.EliminationStrategy = factory.CreateStrategy("lru")
	mytest.Add("123", "123", 3)
	mytest.Add("123", "1", 3)

	value, ok := mytest.Get("123")
	if ok && value == "1" {
		t.Log("测试成功")
	} else {
		t.Error("测试失败")
	}
}

func TestRemoveExpiredKeyLRU(t *testing.T) {
	conf.LoadConfig()
	//expiretime := conf.Config.ExpiredTime
	lrutest := strategies.NewLRUStrategy()

	for i := 0; i < 10; i++ {
		lrutest.Add(models.KeyType(strconv.Itoa(i)), models.ValueType(strconv.Itoa(i)), time.Now().UnixNano()+3000000000*int64(i))
	}
	var ch chan int
	ticker := time.NewTicker(time.Second * time.Duration(conf.Config.UpdateTime))
	//ticker := time.NewTicker(time.Second * time.Duration(7))
	go func() {
		for range ticker.C {
			lrutest.RemoveExpiredTimeKey()
		}
		ch <- 1
	}()
	<-ch
	utils.PrintList(lrutest.GetCacheList())
}
