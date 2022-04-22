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

//根据key获取value
func (lru LRUCache) Get(key models.KeyType) (models.ValueType, bool) {
	logrus.Info("调用LRU的Get操作，key值为", key)
	if element, ok := lru.cacheMap[key]; ok {
		lru.CacheList.MoveToFront(element)
		kv := element.Value.(*models.Entry)
		return kv.Value, true
	}
	return "", false
}

//根据key和value增加一个value，如果已经存在则更新value
func (lru LRUCache) Add(key models.KeyType, value models.ValueType) bool {
	//如果值已经存在，则更新
	if element, ok := lru.cacheMap[key]; ok {
		lru.CacheList.MoveToFront(element)
		kv := element.Value.(*models.Entry)
		kv.Value = value
		return true
	} else {
		element := lru.CacheList.PushFront(&models.Entry{Key: key, Value: value})
		lru.cacheMap[key] = element
		return true
	}
}

//根据key删除对应的Entry
func (lru LRUCache) Remove() bool {
	element := lru.CacheList.Back()
	if element != nil {
		lru.CacheList.Remove(element)
		kv := element.Value.(*models.Entry)
		delete(lru.cacheMap, kv.Key)
	}
	return true
}

func NewLRUCache() *LRUCache {
	return &LRUCache{
		StorageCache: models.NewStorageCache(1000),
		cacheMap:     make(map[models.KeyType]*list.Element),
	}
}
