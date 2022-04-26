package concurrency

import (
	"main/src/models"
	"main/src/strategies"
	"sync"

	"github.com/sirupsen/logrus"
)

type ConcurrencyCache struct {
	lock  sync.Mutex
	cache strategies.EliminationStrategy
}

func (c *ConcurrencyCache) getCache() strategies.EliminationStrategy {
	return c.cache
}

func (c *ConcurrencyCache) Add(key models.KeyType, value models.ValueType) bool {
	logrus.Debug("Add 加锁 " + key)
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.cache == nil {
		factory := new(strategies.StrategyFactory)
		c.cache = factory.CreateStrategy("fifo")
	}
	c.getCache().Add(key, value)
	logrus.Debug("Add 放锁" + key)
	return true
}

func (c *ConcurrencyCache) Get(key models.KeyType) (value models.ValueType, ok bool) {
	logrus.Debug("Get 加锁" + key)
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.getCache() == nil {
		return "", false
	}

	if v, ok := c.getCache().Get(key); ok {
		return v, ok
	}
	logrus.Debug("Get 放锁" + key)
	return "", false
}
