package models

import (
	"container/list"
)

type KeyType interface{}
type ValueType interface{}

type Entry struct {
	key   KeyType
	value ValueType
}

type Cache struct {
	cacheList *list.List
	cacheMap  map[KeyType]*list.Element
}

func New() *Cache {
	return &Cache{
		cacheList: list.New(),
		cacheMap:  make(map[KeyType]*list.Element),
	}
}
