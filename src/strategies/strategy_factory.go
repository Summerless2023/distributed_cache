package strategies

import "github.com/sirupsen/logrus"

type StrategyFactory struct {
}

func (strategyFactory StrategyFactory) CreateStrategy(ext string) EliminationStrategy {
	switch ext {
	case "lru":
		{
			logrus.Info("启动LRU淘汰策略")
			return NewLRUStrategy()
		}
	case "fifo":
		{
			logrus.Info("启动FIFO淘汰策略")
			return NewFIFOStrategy()
		}
	case "lfu":
		{

			logrus.Info("启动LFU淘汰策略")
			return NewLFUStrategy()

		}
	}

	return nil
}
