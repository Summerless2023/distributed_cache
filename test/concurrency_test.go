package test

import (
	"main/src/concurrency"
	"main/src/models"
	"strconv"
	"sync"
	"testing"
)

func add(concurrencyCache *concurrency.ConcurrencyCache, key string, value string, wg *sync.WaitGroup) {
	// concurrencyCache.Add(models.KeyType(key), models.ValueType(value))
	wg.Done()
}

func TestCon(t *testing.T) {
	concurrencyCache := &concurrency.ConcurrencyCache{}
	var waitGroup sync.WaitGroup
	for i := 0; i < 100; i++ {
		waitGroup.Add(1)
		go add(concurrencyCache, strconv.Itoa(i), strconv.Itoa(i), &waitGroup)
	}
	waitGroup.Wait()
	cnt := 0
	for i := 0; i < 100; i++ {
		_, ok := concurrencyCache.Get(models.KeyType(strconv.Itoa(i)))
		if ok {
			cnt++
		}
	}
	t.Log(cnt)
}
