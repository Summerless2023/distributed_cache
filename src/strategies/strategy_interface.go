package strategies

import "main/src/models"

type EliminationStrategy interface {

	// Add adds a value to the cache.
	Add(key models.KeyType, value models.ValueType, expiredTime int64) bool

	// Get look ups a key's value
	Get(key models.KeyType) (models.ValueType, bool)

	//Remove a key's value
	RemoveKey(key models.KeyType) bool
}
