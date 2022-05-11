package strategies

import (
	"main/conf"
	"main/src/models"

	"github.com/sirupsen/logrus"
)

// FIFO cache 结构体
type FIFOStrategy struct {
	*models.CacheStorage
	// optional and executed when an entry is purged.
	// OnEvicted func(key string, value Value)
}

// FIFO cache实例化函数
func NewFIFOStrategy() *FIFOStrategy {
	return &FIFOStrategy{
		CacheStorage: models.NewCacheStorage(conf.DEFAULT_MAX_BYTES),
	}
}

// 根据key获取value
func (fifo *FIFOStrategy) Get(key models.KeyType) (models.ValueType, bool) {
	logrus.Debug("调用FIFO的Get操作，key值为", key)
	if element, ok := fifo.GetCacheMap()[key]; ok {
		// fifo.GetCacheList().MoveToFront(element)
		kv := element.Value.(*models.Entry)
		var err error
		ok, err = fifo.JudgeKeyExpired(key)
		if err != nil { //该键不存在
			return "", false
		}
		if ok { //过期
			return "", false
		} else { //没过期
			return kv.GetValue(), true
		}
	}
	//查询到键不存在
	return "", false
}

//删除一个特定的键值对
func (fifo *FIFOStrategy) RemoveKey(key models.KeyType) bool {
	logrus.Debug("调用FIFO的RemoveKey方法")
	element, _ := fifo.GetCacheMap()[key]
	if element != nil {
		fifo.GetCacheList().Remove(element)
		kv := element.Value.(*models.Entry)
		var tmpBytes int64 = int64(len(kv.GetKey()) + len(kv.GetValue()))
		fifo.SubNBytes(tmpBytes)
		delete(fifo.GetCacheMap(), kv.GetKey())
		delete(fifo.GetExpiredTimeMap(), kv.GetKey())
	}
	return true
}

// 根据策略删除淘汰的Entry
func (fifo *FIFOStrategy) Remove() bool {
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
func (fifo *FIFOStrategy) Add(key models.KeyType, value models.ValueType, expiredTime int64) bool {
	logrus.Debug("调用FIFO的Add方法")
	// 如果值已存在，则更新
	if element, ok := fifo.GetCacheMap()[key]; ok {
		logrus.Debug("key已存在，更新value")
		kv := element.Value.(*models.Entry)
		var tmpBytes int64 = int64(len(value)) - int64(len(kv.GetValue()))
		for fifo.GetNbytes()+tmpBytes > fifo.GetMaxBytes() {
			fifo.Remove()
		}
		fifo.GetExpiredTimeMap()[key] = expiredTime
		kv.SetValue(value)
		fifo.AddNBytes(tmpBytes)
		return true
	} else { // 否则直接加
		logrus.Debug("key不存在，增加kv")
		element := fifo.GetCacheList().PushBack(models.NewEntry(key, value))
		fifo.GetCacheMap()[key] = element
		fifo.GetExpiredTimeMap()[key] = expiredTime
		var tmpBytes int64 = int64(len(key) + len(value))
		for fifo.GetNbytes()+tmpBytes > fifo.GetMaxBytes() {
			fifo.Remove()
		}
		fifo.AddNBytes(tmpBytes)
		return true
	}
}
