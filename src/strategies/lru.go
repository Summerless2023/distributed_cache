package strategies

import (
	"main/conf"
	"main/src/models"
	"time"

	"github.com/sirupsen/logrus"
)

type LRUStrategy struct {
	*models.CacheStorage
	// optional and executed when an entry is purged.
	// OnEvicted func(key string, value Value)
}

//根据key获取value
func (lru *LRUStrategy) Get(key models.KeyType) (models.ValueType, bool) {
	logrus.Debug("调用LRU的Get操作，key值为", key)
	if element, ok := lru.GetCacheMap()[key]; ok {
		lru.GetCacheList().MoveToFront(element)
		kv := element.Value.(*models.Entry)
		var err error
		ok, err = lru.JudgeKeyExpired(key)
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

//根据key和value增加一个value，如果已经存在则更新value
func (lru *LRUStrategy) Add(key models.KeyType, value models.ValueType, expiredTime int64) bool {
	//如果值已经存在，则更新
	if element, ok := lru.GetCacheMap()[key]; ok {
		lru.GetCacheList().MoveToFront(element)
		kv := element.Value.(*models.Entry)
		lru.GetExpiredTimeMap()[key] = expiredTime
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
		lru.GetExpiredTimeMap()[key] = expiredTime
		var tmpBytes int64 = int64(len(key) + len(value))
		for lru.GetNbytes()+tmpBytes > lru.GetMaxBytes() {
			lru.Remove()
		}
		lru.AddNBytes(tmpBytes)
		return true
	}
}

//删除一个特定的键值对
func (lru *LRUStrategy) RemoveKey(key models.KeyType) bool {
	logrus.Debug("调用LRU的RemoveKey方法")
	element, _ := lru.GetCacheMap()[key]
	if element != nil {
		lru.GetCacheList().Remove(element)
		kv := element.Value.(*models.Entry)
		var tmpBytes int64 = int64(len(kv.GetKey()) + len(kv.GetValue()))
		lru.SubNBytes(tmpBytes)
		delete(lru.GetCacheMap(), kv.GetKey())
		delete(lru.GetExpiredTimeMap(), kv.GetKey())
	}
	return true
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

func (lru *LRUStrategy) DeleteRegulary() bool {
	logrus.Info("current time: ", time.Now().UnixNano(), "执行一次定期删除")
	//任务B
	for i := lru.CacheStorage.GetCacheList().Front(); i != nil; {
		kv := i.Value.(*models.Entry)
		//logrus.Info(kv.GetKey())
		isExpired, _ := lru.CacheStorage.JudgeKeyExpired(kv.GetKey())
		//logrus.Info(isExpired)
		//logrus.Info(cacheStorage.expiredTimeMap[kv.GetKey()])
		j := i
		i = i.Next()
		if isExpired {
			//从list,cacheMap,expiredTimeMap中移除
			logrus.Info("移除Element,expiredTime:")
			logrus.Info(kv.GetKey(), "--", lru.CacheStorage.GetExpiredTimeMap()[kv.GetKey()])

			logrus.Info("remove -element", j.Value.(*models.Entry).GetKey())
			lru.CacheStorage.GetCacheList().Remove(j)
			delete(lru.CacheStorage.GetCacheMap(), kv.GetKey())
			delete(lru.CacheStorage.GetExpiredTimeMap(), kv.GetKey())
			//修改内存容量
			var tmpBytes int64 = int64(len(kv.GetKey()) + len(kv.GetValue()))
			lru.CacheStorage.AddNBytes(tmpBytes)
		}
	}
	logrus.Info("剩余Element expiredTime ")
	for i := lru.CacheStorage.GetCacheList().Front(); i != nil; i = i.Next() {
		kv := i.Value.(*models.Entry)
		logrus.Info(kv.GetKey(), "--", lru.CacheStorage.GetExpiredTimeMap()[kv.GetKey()])
	}
	return true
}
func NewLRUStrategy() *LRUStrategy {
	return &LRUStrategy{
		CacheStorage: models.NewCacheStorage(conf.DEFAULT_MAX_BYTES),
	}
}
