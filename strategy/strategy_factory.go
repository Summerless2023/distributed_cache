package strategy

import (
	"log"
	"main/models"
)

type StrategyFactory struct {
}

func (strategyFactory StrategyFactory) CreateStrategy(ext string) EliminationStrategy {
	switch ext {
	case "sample":
		{
			log.Println("启动Sample 淘汰策略")
			return &SampleStrategy{
				make(map[models.KeyType]models.ValueType),
			}
		}

	case "lru":
		return nil
	}
	return nil
}
