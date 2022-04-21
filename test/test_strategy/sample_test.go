package teststrategy

import (
	"../../strategy"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	var mytest eliminationstrategy.EliminationStrategy = new(eliminationstrategy.SampleStrategy)
	mytest.Add("123", "123")
	t.Log("hello world")
}
