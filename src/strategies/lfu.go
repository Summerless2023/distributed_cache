package strategies

import (
	"container/list"
	"main/conf"

	"main/src/models"

	"github.com/sirupsen/logrus"
)

type LFUStrategy struct {
	*models.CacheStorage
	//
	minFreq    int64                    //zk:元素最低使用频率
	freqMap    map[int64]*list.List     //zk:key-频率，value：存储相同频次元素的链表-LFU
	freqKeyMap map[models.KeyType]int64 //存储每一个元素的频率

}

func (lfu *LFUStrategy) GetFreqMap() map[int64]*list.List {
	return lfu.freqMap
}

func (lfu *LFUStrategy) GetFreqKeyMap() map[models.KeyType]int64 {
	return lfu.freqKeyMap
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
		kv := element.Value.(*models.Entry)
		lfu.FreqInc(kv) //增加key的使用频率
		var err error
		ok, err = lfu.JudgeKeyExpired(key)
		if err != nil { //该键不存在
			return "", false
		}
		if ok { //过期
			return "", false
		} else { //没过期
			return kv.GetValue(), true
		}
	}
	//查询到键不存在
	return "", false
}

//删除一个特定的键值对
func (lfu *LFUStrategy) RemoveKey(key models.KeyType) bool {
	logrus.Debug("调用LFU的RemoveKey方法")
	element, _ := lfu.GetCacheMap()[key]
	if element != nil {
		lfu.GetCacheList().Remove(element)
		kv := element.Value.(*models.Entry)
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
		logrus.Debug("移除元素为：", element.Value.(*models.Entry).GetKey(), "freq :", minFreq)
		//更新当前内存
		kv := element.Value.(*models.Entry)
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
		kv := element.Value.(*models.Entry)
		lfu.GetExpiredTimeMap()[key] = expiredTime
		lfu.FreqInc(kv)
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

		listnow := lfu.GetFreqMap()[1] //默认为空
		if listnow == nil {
			//logrus.Debug("创建freq对应的链表")
			listnow = list.New()
			lfu.GetFreqMap()[1] = listnow
		}
		//存入freqmap
		//logrus.Debug("将元素存入freq对应的链表")
		var ele = listnow.PushBack(models.NewEntry(key, value)) //问题
		//logrus.Debug("将元素存入cacheMap")
		lfu.GetFreqKeyMap()[key] = 1
		lfu.GetCacheMap()[key] = ele
		lfu.GetExpiredTimeMap()[key] = expiredTime
		lfu.SetMinFreq(1)
		return true
	}
}

func (lfu *LFUStrategy) FreqInc(lfuentry *models.Entry) bool {
	//将元素从原有freq对应的链表中移除，并更新minFreq
	freq := lfu.GetFreqKeyMap()[lfuentry.GetKey()]
	listold := lfu.GetFreqMap()[freq]
	//确定链表中需要删除的元素
	if listold != nil {
		removeele := listold.Front()
		listold.Remove(removeele)
	} else {
		logrus.Debug("error:Cache has no data,can't remove data")
		return false
	}
	//如果移除的是最低频率元素链表的最后一个元素
	if freq == lfu.GetMinFreq() && listold.Len() == 0 {
		//获取最新的频率
		lfu.SetMinFreq(freq + 1)
	}

	//将元素加入freq+1对应的链表
	freq += 1
	listnow := lfu.GetFreqMap()[freq]
	//如果当前链表为空
	if listnow == nil {
		listnow = list.New()
	}
	//
	element := listnow.PushFront(lfuentry)
	lfu.GetFreqKeyMap()[lfuentry.GetKey()] = freq
	lfu.GetCacheMap()[lfuentry.GetKey()] = element
	lfu.GetFreqMap()[freq] = listnow

	logrus.Debug("key:", lfuentry.GetKey(), " freq:", freq)
	return true
}

func (lfu *LFUStrategy) DeleteRegulary() bool {

	return true
}
func NewLFUStrategy() *LFUStrategy {
	logrus.Debug("Create LFUCache")
	return &LFUStrategy{
		CacheStorage: models.NewCacheStorage(conf.DEFAULT_MAX_BYTES),
		minFreq:      1,
		freqMap:      make(map[int64]*list.List),
		freqKeyMap:   make(map[models.KeyType]int64),
	}
}
