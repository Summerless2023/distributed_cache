package test

import (
	"main/strategy"
	"testing"
)

func TestLRUStrategy(t *testing.T) {
	factory := new(strategy.StrategyFactory)
	var mytest strategy.EliminationStrategy = factory.CreateStrategy("lru")
	mytest.Add("123", "123")
	t.Log("测试完成")
}

func TestLRUStrategy1(t *testing.T) {
	factory := new(strategy.StrategyFactory)
	var mytest strategy.EliminationStrategy = factory.CreateStrategy("lru")
	mytest.Add("123", "123")
	if _, ok := mytest.Get("123"); ok {
		t.Log("测试成功")
	} else {
		t.Errorf("LRU Get 函数错误")
	}

}
