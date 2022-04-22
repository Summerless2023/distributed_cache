package models

import "container/list"

type KeyType string   //cache的key值类型
type ValueType string //cache的value值类型

//存放元素实体的结构体
type Entry struct {
	Key   KeyType
	Value ValueType
}

type StorageCache struct {
	MaxBytes  int64      //最大内存占用
	Nbytes    int64      // 当前内存占用
	CacheList *list.List //存储缓存的链表
}

func (storageCache *StorageCache) GetMaxBytes() int64 {
	return storageCache.MaxBytes
}

func (storageCache *StorageCache) GetNbytes() int64 {
	return storageCache.Nbytes
}

func (storageCache *StorageCache) AddBytes(bytesNum int64) {
	storageCache.Nbytes += bytesNum
}

func NewStorageCache(maxBytes int64) *StorageCache {
	return &StorageCache{
		MaxBytes:  maxBytes,
		Nbytes:    0,
		CacheList: list.New(),
	}
}
