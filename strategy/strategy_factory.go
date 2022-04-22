package strategy

import (
	"log"
)

type StrategyFactory struct {
}

func (strategyFactory StrategyFactory) CreateStrategy(ext string) EliminationStrategy {
	switch ext {
	case "sample":
		{
			log.Println("启动Sample 淘汰策略")
			return NewSampleStrategy()
		}

	case "lru":
		{
			log.Println("启动LRU淘汰策略")
			return NewLRUCache()
		}
	}
	return nil
}
