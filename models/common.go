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
	maxBytes  int64                     //最大内存占用
	nbytes    int64                     //当前内存占用
	cacheList *list.List                //存储缓存的链表
	cacheMap  map[KeyType]*list.Element //key到list.Element的映射map
}

func (storageCache *StorageCache) GetMaxBytes() int64 {
	return storageCache.maxBytes
}

func (storageCache *StorageCache) GetNbytes() int64 {
	return storageCache.nbytes
}

func (storageCache *StorageCache) AddBytes(bytesNum int64) {
	storageCache.nbytes += bytesNum
}

func (storageCache *StorageCache) SubBytes(bytesNum int64) {
	storageCache.nbytes -= bytesNum
}

func (storageCache *StorageCache) GetCacheList() *list.List {
	return storageCache.cacheList
}

func (storageCache *StorageCache) GetCacheMap() map[KeyType]*list.Element {
	return storageCache.cacheMap
}

func NewStorageCache(maxBytes int64) *StorageCache {
	return &StorageCache{
		maxBytes:  maxBytes,
		nbytes:    0,
		cacheList: list.New(),
		cacheMap:  make(map[KeyType]*list.Element),
	}
}
