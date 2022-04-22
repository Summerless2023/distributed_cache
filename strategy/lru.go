package strategy

import (
	"main/conf"
	"main/models"

	"github.com/sirupsen/logrus"
)

type LRUCache struct {
	*models.StorageCache
	// optional and executed when an entry is purged.
	// OnEvicted func(key string, value Value)
}

//根据key获取value
func (lru *LRUCache) Get(key models.KeyType) (models.ValueType, bool) {
	logrus.Debug("调用LRU的Get操作，key值为", key)
	if element, ok := lru.CacheMap[key]; ok {
		lru.CacheList.MoveToFront(element)
		kv := element.Value.(*models.Entry)
		return kv.GetValue(), true
	}
	return "", false
}

//根据key和value增加一个value，如果已经存在则更新value
func (lru *LRUCache) Add(key models.KeyType, value models.ValueType) bool {
	//如果值已经存在，则更新
	if element, ok := lru.CacheMap[key]; ok {
		lru.CacheList.MoveToFront(element)
		kv := element.Value.(*models.Entry)
		kv.SetValue(value)
		return true
	} else { //否则直接增加
		logrus.Debug("LRU Add (", key, ",", value, ")")
		element := lru.CacheList.PushFront(models.NewEntry(key, value))
		lru.CacheMap[key] = element
		var tmpBytes int64 = int64(len(key) + len(value))
		for lru.GetNbytes()+tmpBytes > lru.GetMaxBytes() {
			lru.Remove()
		}
		lru.AddBytes(tmpBytes)
		return true
	}
}

//根据key删除对应的Entry
func (lru *LRUCache) Remove() bool {
	logrus.Debug("调用LRU的Remove方法")
	element := lru.CacheList.Back()
	if element != nil {
		lru.CacheList.Remove(element)
		kv := element.Value.(*models.Entry)
		var tmpBytes int64 = int64(len(kv.GetKey()) + len(kv.GetValue()))
		lru.SubBytes(tmpBytes)
		delete(lru.CacheMap, kv.GetKey())
	}
	return true
}

func NewLRUCache() *LRUCache {
	return &LRUCache{
		StorageCache: models.NewStorageCache(conf.Default_Max_Bytes),
	}
}
