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

func (entry *Entry) GetEntryLen() int64 {
	return int64(len(entry.key)) + int64(len(entry.value))
}

func NewEntry(key KeyType, value ValueType) *Entry {
	return &Entry{
		key:   key,
		value: value,
	}
}

type CacheStorage struct {
	maxBytes  int64                     //最大内存占用
	nbytes    int64                     //当前内存占用
	cacheList *list.List                //存储缓存的链表
	cacheMap  map[KeyType]*list.Element //key到list.Element的映射map
}

func (cacheStorage *CacheStorage) GetMaxBytes() int64 {
	return cacheStorage.maxBytes
}

func (cacheStorage *CacheStorage) GetNbytes() int64 {
	return cacheStorage.nbytes
}

func (cacheStorage *CacheStorage) AddNBytes(bytesNum int64) {
	cacheStorage.nbytes += bytesNum
}

func (cacheStorage *CacheStorage) SubNBytes(bytesNum int64) {
	cacheStorage.nbytes -= bytesNum
}

func (cacheStorage *CacheStorage) GetCacheList() *list.List {
	return cacheStorage.cacheList
}

func (cacheStorage *CacheStorage) GetCacheMap() map[KeyType]*list.Element {
	return cacheStorage.cacheMap
}

func NewCacheStorage(maxBytes int64) *CacheStorage {
	return &CacheStorage{
		maxBytes:  maxBytes,
		nbytes:    0,
		cacheList: list.New(),
		cacheMap:  make(map[KeyType]*list.Element),
	}
}

type Getter interface {
	Get(key KeyType) (ValueType, error)
}

// A GetterFunc implements Getter with a function.
type GetterFunc func(key KeyType) (ValueType, error)

// Get implements Getter interface function
func (f GetterFunc) Get(key KeyType) (ValueType, error) {
	return f(key)
}
