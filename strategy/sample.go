package strategy

import "main/models"

type SampleStrategy struct {
	cache map[models.KeyType]models.ValueType
}

func (sampleStrategy SampleStrategy) Get(key models.KeyType) models.ValueType {
	return sampleStrategy.cache[key]
}

func (sampleStrategy SampleStrategy) Add(key models.KeyType, value models.ValueType) {
	sampleStrategy.cache[key] = value
}

func (sampleStrategy SampleStrategy) Len() int {
	return len(sampleStrategy.cache)
}

func New() *SampleStrategy {
	return &SampleStrategy{
		cache: make(map[models.KeyType]models.ValueType),
	}
}
