package strategy

import (
	"main/models"
)

type LRU struct {
	cache models.ValueType
}

func (lru LRU) Get(key models.KeyType) models.ValueType {
	return lru.cache
}

func (lru LRU) Add(key models.KeyType, value models.ValueType) {

}

func (lru LRU) Len() int {
	return 1
}
