package test

import (
	"main/src/models"
	"main/src/strategy"
	"main/src/utils"
	"strconv"
	"testing"

	"github.com/sirupsen/logrus"
)

// func addelement()

func TestFIFO(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Info("FIFOCache")
	fifotest := strategy.NewFIFOCache()
	for i := 0; i < 20; i++ {
		fifotest.Add(models.KeyType(strconv.Itoa(i)), models.ValueType(strconv.Itoa(i)))
	}
	utils.PrintList(*fifotest.GetCacheList())

}

func TestFactoryFIFO(t *testing.T) {
	factory := new(strategy.StrategyFactory)
	var mytest strategy.EliminationStrategy = factory.CreateStrategy("fifo")
	mytest.Add("123", "123")
	if _, ok := mytest.Get("123"); ok {
		t.Log("测试成功")
	} else {
		t.Errorf("LRU Get 函数错误")
	}
}
