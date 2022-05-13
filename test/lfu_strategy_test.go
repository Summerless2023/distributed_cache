package test

import (
	"main/conf"
	"main/src/models"
	"main/src/strategies"
	"time"

	//"main/src/utils"
	"fmt"
	"strconv"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLFU(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Info("LFUStrategy")
	lfutest := strategies.NewLFUStrategy()
	conf.LoadConfig()

	for j := 0; j < 3; j++ {
		for i := 5; i < 10; i++ {
			lfutest.Add(models.KeyType(strconv.Itoa(i)), models.ValueType(strconv.Itoa(2*i)), conf.Config.ExpiredTime)
		}
	}
	for i := 5; i < 15; i++ {
		lfutest.Add(models.KeyType(strconv.Itoa(i)), models.ValueType(strconv.Itoa(3*i)), conf.Config.ExpiredTime)
	}
	logrus.Debug("打印当前缓存中频率对应元素")
	for _, value := range lfutest.GetCacheMap() {
		lfuentry := value.Value.(*models.Entry)
		fmt.Println(lfuentry.GetKey(), ":", lfuentry.GetValue(), ":", lfutest.GetFreqKeyMap()[lfuentry.GetKey()])
	}
	//utils.PrintList(lfutest.GetCacheList())
}

// func addelement()
func TestFactoryLFU(t *testing.T) {
	factory := new(strategies.StrategyFactory)
	var lfutest strategies.EliminationStrategy = factory.CreateStrategy("lfu")
	lfutest.Add("123", "123", time.Now().UnixNano()+10000000000)
	if _, ok := lfutest.Get("123"); ok {
		t.Log("测试成功")
	} else {
		t.Errorf("LFU Get 函数错误")
	}
}
