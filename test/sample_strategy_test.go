package test

import (
	"main/src/strategies"
	"testing"
)

func TestSampleStrategy(t *testing.T) {
	t.Log("开始测试Sample淘汰策略")
	factory := new(strategies.StrategyFactory)
	var mytest strategies.EliminationStrategy = factory.CreateStrategy("sample")
	mytest.Add("123", "123")
	t.Log("测试完成")
}
