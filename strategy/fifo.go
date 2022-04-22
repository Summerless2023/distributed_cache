package strategy

import (
	"main/conf"
	"main/models"

	"github.com/sirupsen/logrus"
)

// FIFO cache 结构体
type FIFOCache struct {
	*models.StorageCache
	// optional and executed when an entry is purged.
	// OnEvicted func(key string, value Value)
}

// FIFO cache实例化函数
func NewFIFOCache() *FIFOCache {
	return &FIFOCache{
		StorageCache: models.NewStorageCache(conf.Default_Max_Bytes),
	}
}

// 根据key获取value
func (fifo *FIFOCache) Get(key models.KeyType) (models.ValueType, bool) {
	logrus.Debug("调用FIFO的Get操作，key值为", key)
	if element, ok := fifo.GetCacheMap()[key]; ok {
		// fifo.GetCacheList().MoveToFront(element)
		kv := element.Value.(*models.Entry)
		return kv.GetValue(), true
	}
	return "", false
}

// 根据策略删除淘汰的Entry
func (fifo *FIFOCache) Remove() bool {
	logrus.Debug("调用FIFO的Remove方法")
	element := fifo.GetCacheList().Front()
	if element != nil {
		fifo.GetCacheList().Remove(element)
		kv := element.Value.(*models.Entry)
		var tmpBytes int64 = int64(len(kv.GetKey()) + len(kv.GetValue()))
		fifo.SubNBytes(tmpBytes)
		delete(fifo.GetCacheMap(), kv.GetKey())
	}
	return true
}

// 根据key和value增加一个value，如果已经存在则更新value
func (fifo *FIFOCache) Add(key models.KeyType, value models.ValueType) bool {
	logrus.Debug("调用FIFO的Add方法")
	// 如果值已存在，则更新
	if element, ok := fifo.GetCacheMap()[key]; ok {
		logrus.Debug("key已存在，更新value")
		kv := element.Value.(*models.Entry)
		var tmpBytes int64 = int64(len(value)) - int64(len(kv.GetValue()))
		for fifo.GetNbytes()+tmpBytes > fifo.GetMaxBytes() {
			fifo.Remove()
		}
		kv.SetValue(value)
		fifo.AddNBytes(tmpBytes)
		return true
	} else { // 否则直接加
		logrus.Debug("key不存在，增加kv")
		element := fifo.GetCacheList().PushBack(models.NewEntry(key, value))
		fifo.GetCacheMap()[key] = element
		var tmpBytes int64 = int64(len(key) + len(value))
		for fifo.GetNbytes()+tmpBytes > fifo.GetMaxBytes() {
			fifo.Remove()
		}
		fifo.AddNBytes(tmpBytes)
		return true
	}
}
