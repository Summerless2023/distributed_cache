package test

import (
	"fmt"
	"main/conf"
	"main/src/models"
	"main/src/strategies"
	"main/src/utils"
	"strconv"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

// func addelement()

func TestFIFO(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Info("FIFOStrategy")
	conf.LoadConfig()
	//expiretime := conf.Config.ExpiredTime
	fifotest := strategies.NewFIFOStrategy()

	for i := 0; i < 20; i++ {
		fifotest.Add(models.KeyType(strconv.Itoa(i)), models.ValueType(strconv.Itoa(i)), time.Now().UnixNano()+3000000000*int64(i))
	}
}

func TestFactoryFIFO(t *testing.T) {
	factory := new(strategies.StrategyFactory)
	var mytest strategies.EliminationStrategy = factory.CreateStrategy("fifo")
	mytest.Add("123", "123", conf.Config.ExpiredTime)
	if _, ok := mytest.Get("123"); ok {
		t.Log("测试成功")
	} else {
		t.Errorf("LRU Get 函数错误")
	}
}

func TestRemoveExpiredKeyFIFO(t *testing.T) {
	conf.LoadConfig()
	expiretime := conf.Config.ExpiredTime
	fmt.Println("ut:", conf.Config.UpdateTime)
	fmt.Println("mb:", conf.Config.Max_Bytes)
	fmt.Println("s:", conf.Config.Strategy)
	fmt.Println("up:", conf.Config.UpdateTime)
	fifotest := strategies.NewFIFOStrategy()

	for i := 0; i < 20; i++ {
		fifotest.Add(models.KeyType(strconv.Itoa(i)), models.ValueType(strconv.Itoa(i)), time.Now().UnixNano()+3000000000*int64(i))
	}
	fmt.Println("conf:", conf.Config.UpdateTime, expiretime)
	var ch chan int
	ticker := time.NewTicker(time.Second * time.Duration(conf.Config.UpdateTime))
	//ticker := time.NewTicker(time.Second * time.Duration(7))
	go func() {
		for range ticker.C {
			fifotest.RemoveExpiredTimeKey()
		}
		ch <- 1
	}()
	<-ch
	utils.PrintList(fifotest.GetCacheList())
}
