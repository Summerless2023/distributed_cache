package strategy

import (
	"container/list"
	"main/models"

	"github.com/sirupsen/logrus"
)

type LRUCache struct {
	*models.StorageCache
	cacheMap map[models.KeyType]*list.Element
	// optional and executed when an entry is purged.
	// OnEvicted func(key string, value Value)
}

func (lru LRUCache) Get(key models.KeyType) (models.ValueType, bool) {
	logrus.Info("调用LRU的Get操作，key值为", key)
	if element, ok := lru.cacheMap[key]; ok {
		kv := element.Value.(*models.Entry)
		return kv.Value, true
	}
	return "", false
}

func (lru LRUCache) Add(key models.KeyType, value models.ValueType) bool {
	element := lru.CacheList.PushFront(&models.Entry{Key: key, Value: value})
	lru.cacheMap[key] = element
	return true
}

func (lru LRUCache) Len() int {
	return 1
}

func NewLRUCache() *LRUCache {
	return &LRUCache{
		StorageCache: models.NewStorageCache(1000),
		cacheMap:     make(map[models.KeyType]*list.Element),
	}
}
