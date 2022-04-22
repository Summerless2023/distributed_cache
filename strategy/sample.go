package strategy

import "main/models"

type SampleStrategy struct {
	cache map[models.KeyType]models.ValueType
}

func (sampleStrategy SampleStrategy) Get(key models.KeyType) (models.ValueType, bool) {
	return sampleStrategy.cache[key], true
}

func (sampleStrategy SampleStrategy) Add(key models.KeyType, value models.ValueType) bool {
	sampleStrategy.cache[key] = value
	return true
}

func (sampleStrategy SampleStrategy) Len() int {
	return len(sampleStrategy.cache)
}

func NewSampleStrategy() *SampleStrategy {
	return &SampleStrategy{
		cache: make(map[models.KeyType]models.ValueType),
	}
}
