package models

import "container/list"

type KeyType string   //cache的key值类型
type ValueType string //cache的value值类型

//存放元素实体的结构体
type Entry struct {
	key   KeyType
	value ValueType
}

func (entry *Entry) GetKey() KeyType {
	return entry.key
}

func (entry *Entry) GetValue() ValueType {
	return entry.value
}

func (entry *Entry) SetKey(key KeyType) {
	entry.key = key
}

func (entry *Entry) SetValue(value ValueType) {
	entry.value = value
}

func NewEntry(key KeyType, value ValueType) *Entry {
	return &Entry{
		key:   key,
		value: value,
	}
}

type StorageCache struct {
	MaxBytes  int64                     //最大内存占用
	Nbytes    int64                     // 当前内存占用
	CacheList *list.List                //存储缓存的链表
	CacheMap  map[KeyType]*list.Element //key到list.Element的映射map
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

func (storageCache *StorageCache) SubBytes(bytesNum int64) {
	storageCache.Nbytes -= bytesNum
}

func NewStorageCache(maxBytes int64) *StorageCache {
	return &StorageCache{
		MaxBytes:  maxBytes,
		Nbytes:    0,
		CacheList: list.New(),
		CacheMap:  make(map[KeyType]*list.Element),
	}
}
