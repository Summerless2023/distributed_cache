package strategies

import (
	"container/list"
	"main/conf"
	"main/src/models"

	"github.com/sirupsen/logrus"
)

type LFUEntry struct {
	*models.Entry
	freq int64
}

func NewLFUEntry(key models.KeyType, value models.ValueType, freq int64) *LFUEntry {
	return &LFUEntry{
		Entry: models.NewEntry(key, value),
		freq:  freq,
	}
}
func (entry *LFUEntry) GetFreq() int64 {
	return entry.freq
}
func (entry *LFUEntry) SetFreq(freq int64) {
	entry.freq = freq
}

type LFUStrategy struct {
	*models.CacheStorage
	//
	minFreq int64                //zk:元素最低使用频率
	freqMap map[int64]*list.List //zk:key-频率，value：存储相同频次元素的链表-LFU
}

func (lfu *LFUStrategy) GetFreqMap() map[int64]*list.List {
	return lfu.freqMap
}

func (lfu *LFUStrategy) GetMinFreq() int64 {
	return lfu.minFreq
}

func (lfu *LFUStrategy) SetMinFreq(freq int64) {
	lfu.minFreq = freq
}

func (lfu *LFUStrategy) Get(key models.KeyType) (models.ValueType, bool) {
	logrus.Debug("调用LFU的Get操作，key值为", key)
	if element, ok := lfu.GetCacheMap()[key]; ok {
		kv := element.Value.(*LFUEntry)
		lfu.FreqInc(*kv) //增加key的使用频率
		var err error
		ok, err = lfu.JudgeKeyExpired(key)
		if err != nil { //该键不存在
			return "", false
		}
		if ok { //过期
			return kv.GetValue(), false
		} else { //没过期
			return kv.GetValue(), true
		}
	}
	return "", false
}

//删除一个特定的键值对
func (lfu *LFUStrategy) RemoveKey(key models.KeyType) bool {
	logrus.Debug("调用LFU的RemoveKey方法")
	element, _ := lfu.GetCacheMap()[key]
	if element != nil {
		lfu.GetCacheList().Remove(element)
		kv := element.Value.(*LFUEntry)
		var tmpBytes int64 = int64(len(kv.GetKey()) + len(kv.GetValue()))
		lfu.SubNBytes(tmpBytes)
		//map中移除元素
		delete(lfu.GetCacheMap(), kv.GetKey())
		delete(lfu.GetExpiredTimeMap(), kv.GetKey())
	}
	return true
}

func (lfu *LFUStrategy) Remove() bool {
	logrus.Debug("调用LFU的Remove方法")
	minFreq := lfu.GetMinFreq()
	minlist := lfu.GetFreqMap()[minFreq]
	if minlist != nil {
		//Cache中移除元素
		element := minlist.Back()
		lfu.GetFreqMap()[minFreq].Remove(element)
		//
		// logrus.Debug("打印当前缓存中最低频率对应元素")

		// for _, value := range lfu.GetCacheMap() {
		// 	lfuentry := value.Value.(*LFUEntry)
		// 	if lfuentry.GetFreq() == minFreq {
		// 		fmt.Println(lfuentry.GetKey(), ":", lfuentry.GetValue(), ":", lfuentry.GetFreq())
		// 	}
		// }
		logrus.Debug("移除元素为：", element.Value.(*LFUEntry).GetKey(), "freq :", minFreq)
		//更新当前内存
		kv := element.Value.(*LFUEntry)
		var tmpBytes int64 = int64(len(kv.GetKey()) + len(kv.GetValue()))
		lfu.SubNBytes(tmpBytes)
		//map中移除元素
		delete(lfu.GetCacheMap(), kv.GetKey())
	}
	return true
}

func (lfu *LFUStrategy) Add(key models.KeyType, value models.ValueType, expiredTime int64) bool {
	//如果值已存在
	if element, ok := lfu.GetCacheMap()[key]; ok {
		//更新频率
		logrus.Debug("key存在 LFU更新频率 ")
		kv := element.Value.(*LFUEntry)
		lfu.GetExpiredTimeMap()[key] = expiredTime
		lfu.FreqInc(*kv)
		//检查内存容量是否可用
		var tmpBytes int64 = int64(len(kv.GetValue())) - int64(len(value)) //如果值为正，则表示占用byte增加，值为负则表示占用减少
		kv.SetValue(value)
		for lfu.GetNbytes()+tmpBytes > lfu.GetMaxBytes() {
			lfu.Remove()
		}
		lfu.AddNBytes(tmpBytes)
		return true
	} else {
		logrus.Debug("key 不存在 LFU Add (", key, ",", value, ")")
		//值不存在
		var tmpBytes int64 = int64(len(key) + len(value))
		for lfu.GetNbytes()+tmpBytes > lfu.GetMaxBytes() {
			lfu.Remove()
		}
		lfu.AddNBytes(tmpBytes)

		//添加元素Cache
		//logrus.Debug("begin:Add to list")
		//lfu.AddtoList(*lfuentry) //问题

		listnow := lfu.GetFreqMap()[1] //默认为空
		if listnow == nil {
			//logrus.Debug("创建freq对应的链表")
			listnow = list.New()
			lfu.GetFreqMap()[1] = listnow
		}
		//存入freqmap
		//logrus.Debug("将元素存入freq对应的链表")
		var ele = listnow.PushBack(NewLFUEntry(key, value, 1)) //问题
		//logrus.Debug("将元素存入cacheMap")
		lfu.GetCacheMap()[key] = ele
		lfu.GetExpiredTimeMap()[key] = expiredTime
		//更新最小频率
		//logrus.Debug("更新minFreq")
		lfu.SetMinFreq(1)
		//logrus.Debug("end  :Add to list")

		return true
	}
}

// func (lfu *LFUStrategy) AddtoList(lfuentry LFUEntry) bool {
// 	listnow := lfu.GetFreqMap()[1]
// 	if listnow == nil {
// 		logrus.Debug("创建freq对应的链表")
// 		listnow := list.New()
// 		lfu.GetFreqMap()[1] = listnow
// 	}
// 	//存入freqmap
// 	logrus.Debug("将元素存入freq对应的链表")
// 	var element1 = listnow.PushFront(lfuentry) //问题
// 	logrus.Debug("将元素存入cacheMap")
// 	lfu.GetCacheMap()[lfuentry.GetKey()] = element1
// 	//更新最小频率
// 	logrus.Debug("更新minFreq")
// 	lfu.SetMinFreq(1)
// 	return true
// }

func (lfu *LFUStrategy) FreqInc(lfuentry LFUEntry) bool {
	//将元素从原有freq对应的链表中移除，并更新minFreq
	freq := lfuentry.GetFreq()
	listold := lfu.GetFreqMap()[freq]
	//确定链表中需要删除的元素
	if listold != nil {
		removeele := listold.Front()
		// if removeele == nil {
		// 	logrus.Debug("removeele is nil")
		// } else {
		// 	listold.Remove(removeele) //问题
		// }
		listold.Remove(removeele)

	} else {
		logrus.Debug("error:Cache has no data,can't remove data")
		return false
	}

	//如果移除的是最低频率元素链表的最后一个元素
	if freq == lfu.GetMinFreq() && listold.Len() == 0 {
		lfu.SetMinFreq(freq + 1)
	}

	//将元素加入freq+1对应的链表
	freq += 1
	listnow := lfu.GetFreqMap()[freq]
	//如果当前链表为空
	if listnow == nil {
		listnow = list.New()
	}

	element := listnow.PushFront(NewLFUEntry(lfuentry.GetKey(), lfuentry.GetValue(), freq))
	lfu.GetCacheMap()[lfuentry.GetKey()] = element
	lfu.GetFreqMap()[freq] = listnow

	logrus.Debug("key:", lfuentry.GetKey(), " freq:", freq)
	return true
}

func NewLFUStrategy() *LFUStrategy {
	logrus.Debug("Create LFUCache")
	return &LFUStrategy{
		CacheStorage: models.NewCacheStorage(conf.DEFAULT_MAX_BYTES),
		minFreq:      1,
		freqMap:      make(map[int64]*list.List),
	}
}
