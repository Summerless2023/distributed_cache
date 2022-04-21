package test

import (
	"main/strategy"
	"testing"
)

func TestSampleStrategy(t *testing.T) {
	t.Log("开始测试Sample淘汰策略")
	factory := new(strategy.StrategyFactory)
	var mytest strategy.EliminationStrategy = factory.CreateStrategy("sample")
	mytest.Add("123", "123")
	t.Log("测试完成")
}
