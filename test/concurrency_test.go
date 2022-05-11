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

func TestConFIFOGet(t *testing.T) {
	c := concurrency.NewConcurrencyCache() //通过更改concurrencycache里的参数验证不同的策略模式
	c.Add("1", "1", 1652188799104820800)   //这是一个过期的键
	c.Add("2", "2", 3000000000000000000)   //这是一个没有过期的键

	if _, ok := c.Get("1"); ok {
		t.Log("测试失败")
	} else {
		t.Log("测试成功")
	}

	if _, ok := c.Get("2"); ok {
		t.Log("测试成功")
	} else {
		t.Log("测试失败")
	}
}

func TestConLFUGet(t *testing.T) {
	c := concurrency.NewConcurrencyCache() //通过更改concurrencycache里的参数验证不同的策略模式
	c.Add("1", "1", 1652188799104820800)   //这是一个过期的键
	c.Add("2", "2", 3000000000000000000)   //这是一个没有过期的键

	if _, ok := c.Get("1"); ok {
		t.Log("测试失败")
	} else {
		t.Log("测试成功")
	}

	if _, ok := c.Get("2"); ok {
		t.Log("测试成功")
	} else {
		t.Log("测试失败")
	}
}

func TestConLRUGet(t *testing.T) {
	c := concurrency.NewConcurrencyCache() //通过更改concurrencycache里的参数验证不同的策略模式
	c.Add("1", "1", 1652188799104820800)   //这是一个过期的键
	c.Add("2", "2", 3000000000000000000)   //这是一个没有过期的键

	if _, ok := c.Get("1"); ok {
		t.Log("测试失败")
	} else {
		t.Log("测试成功")
	}

	if _, ok := c.Get("2"); ok {
		t.Log("测试成功")
	} else {
		t.Log("测试失败")
	}
}
