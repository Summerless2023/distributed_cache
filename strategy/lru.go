package eliminationstrategy

import (
	"../models"
)

type LRU struct {
	storge models.KeyType
}

func (lru LRU) Get(key models.KeyType) models.ValueType {
	return lru.storge
}

func (lru LRU) Add(key models.KeyType, value models.ValueType) {
	lru.storge = value
}

func (lru LRU) Len() int {
	return 1
}
