package strategies

import (
	"main/conf"
	"main/src/models"

	"github.com/sirupsen/logrus"
)

type LRUStrategy struct {
	*models.StorageCache
	// optional and executed when an entry is purged.
	// OnEvicted func(key string, value Value)
}

//根据key获取value
func (lru *LRUStrategy) Get(key models.KeyType) (models.ValueType, bool) {
	logrus.Debug("调用LRU的Get操作，key值为", key)
	if element, ok := lru.GetCacheMap()[key]; ok {
		lru.GetCacheList().MoveToFront(element)
		kv := element.Value.(*models.Entry)
		return kv.GetValue(), true
	}
	return "", false
}

//根据key和value增加一个value，如果已经存在则更新value
func (lru *LRUStrategy) Add(key models.KeyType, value models.ValueType) bool {
	//如果值已经存在，则更新
	if element, ok := lru.GetCacheMap()[key]; ok {
		lru.GetCacheList().MoveToFront(element)
		kv := element.Value.(*models.Entry)
		var tmpBytes int64 = int64(len(kv.GetValue())) - int64(len(value)) //如果值为正，则表示占用byte增加，值为负则表示占用减少
		kv.SetValue(value)
		for lru.GetNbytes()+tmpBytes > lru.GetMaxBytes() {
			lru.Remove()
		}
		lru.AddNBytes(tmpBytes)
		return true
	} else { //否则直接增加
		logrus.Debug("LRU Add (", key, ",", value, ")")
		element := lru.GetCacheList().PushFront(models.NewEntry(key, value))
		lru.GetCacheMap()[key] = element
		var tmpBytes int64 = int64(len(key) + len(value))
		for lru.GetNbytes()+tmpBytes > lru.GetMaxBytes() {
			lru.Remove()
		}
		lru.AddNBytes(tmpBytes)
		return true
	}
}

//根据key删除对应的Entry
func (lru *LRUStrategy) Remove() bool {
	logrus.Debug("调用LRU的Remove方法")
	element := lru.GetCacheList().Back()
	if element != nil {
		lru.GetCacheList().Remove(element)
		kv := element.Value.(*models.Entry)
		var tmpBytes int64 = int64(len(kv.GetKey()) + len(kv.GetValue()))
		lru.SubNBytes(tmpBytes)
		delete(lru.GetCacheMap(), kv.GetKey())
	}
	return true
}

func NewLRUStrategy() *LRUStrategy {
	return &LRUStrategy{
		StorageCache: models.NewStorageCache(conf.Default_Max_Bytes),
	}
}
