package models

import (
	"container/list"
	"errors"
	"fmt"
	"time"
)

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
	maxBytes       int64                     //最大内存占用
	nbytes         int64                     //当前内存占用
	cacheList      *list.List                //存储缓存的链表
	cacheMap       map[KeyType]*list.Element //key到list.Element的映射map
	expiredTimeMap map[KeyType]int64         //key到过期时间的映射map
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

func (cacheStorage *CacheStorage) GetExpiredTimeMap() map[KeyType]int64 {
	return cacheStorage.expiredTimeMap
}

func (cacheStorage *CacheStorage) GetExpiredTime(key KeyType) (int64, bool) {
	if element, ok := cacheStorage.GetExpiredTimeMap()[key]; ok {
		return element, true
	}
	return 0, false
}

//判断key是否过期，如果过期则返回(true，nil)，没有过期返回(false,nil),否则返回(false，error)
func (cacheStorage *CacheStorage) JudgeKeyExpired(key KeyType) (bool, error) {
	if expiredTime, ok := cacheStorage.GetExpiredTime(key); ok {
		if expiredTime < int64(time.Now().UnixNano()) {
			return true, nil
		} else {
			return false, nil
		}
	}
	return false, errors.New("key is not in expiredTimeMap")
}

func (cacheStorage *CacheStorage) RemoveExpiredTimeKey() {

	fmt.Println(time.Now().UnixNano())
	//任务B
	for i := cacheStorage.cacheList.Front(); i != nil; i = i.Next() {
		kv := i.Value.(*Entry)
		isExpired, _ := cacheStorage.JudgeKeyExpired(kv.GetKey())
		//fmt.Println("Element = ", kv.GetKey(), kv.GetValue(), cacheStorage.expiredTimeMap[kv.GetKey()])
		if isExpired {
			//从list,cacheMap,expiredTimeMap中移除
			fmt.Println("移除Element = ", kv.GetKey(), kv.GetValue(), cacheStorage.expiredTimeMap[kv.GetKey()])
			cacheStorage.cacheList.Remove(i)
			delete(cacheStorage.cacheMap, kv.GetKey())
			delete(cacheStorage.expiredTimeMap, kv.key)
			//修改内存容量
			var tmpBytes int64 = int64(len(kv.GetKey()) + len(kv.GetValue()))
			cacheStorage.AddNBytes(tmpBytes)
		}
	}

	for i := cacheStorage.cacheList.Front(); i != nil; i = i.Next() {
		kv := i.Value.(*Entry)
		fmt.Println("剩余Element = ", kv.GetKey(), kv.GetValue(), cacheStorage.expiredTimeMap[kv.GetKey()])
	}

}

func NewCacheStorage(maxBytes int64) *CacheStorage {
	return &CacheStorage{
		maxBytes:       maxBytes,
		nbytes:         0,
		cacheList:      list.New(),
		cacheMap:       make(map[KeyType]*list.Element),
		expiredTimeMap: make(map[KeyType]int64),
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
