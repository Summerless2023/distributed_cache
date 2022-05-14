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

func (c *ConcurrencyCache) Add(key models.KeyType, value models.ValueType, expiredTime int64) bool {
	logrus.Debug("Add 加锁 " + key)
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.cache == nil {
		factory := new(strategies.StrategyFactory)
		c.cache = factory.CreateStrategy("lru")
	}
	c.getCache().Add(key, value, expiredTime)
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

	val, ok := c.getCache().Get(key)
	if ok { //查询到键且不过期
		return val, true
	} else { //查询到不存在或者过期
		c.getCache().RemoveKey(key)
	}
	logrus.Debug("Get 放锁" + key)
	return "", false
}

func (c *ConcurrencyCache) DeleteRegulary() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	//锁住Cache
	if c.getCache() == nil {
		return false
	}
	c.getCache().DeleteRegulary()
	return true
}
func NewConcurrencyCache() *ConcurrencyCache {
	return &ConcurrencyCache{}
}
