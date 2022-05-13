package main

import (
	"main/conf"
	"main/src/concurrency"
	"main/src/models"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	//"net/http"
	//"github.com/sirupsen/logrus"
)

// type server int

// func (h *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	logrus.Println(r.URL.Path)
// 	w.Write([]byte("Hello World!"))
// }

func main() {
	conf.LoadConfig()                   //进行配置信息的加载
	logrus.Info(conf.Config.Max_Bytes)  //在conf中定义全局变量，其他包进行引用
	logrus.Info(conf.Config.Strategy)   //在conf中定义全局变量，其他包进行引用
	logrus.Info(conf.Config.UpdateTime) //获取key的定期删除时间
	//var s server
	//http.ListenAndServe("localhost:9999", &s)
	//test1()

	ConCache := concurrency.NewConcurrencyCache()
	for i := 0; i < 10; i++ {
		ConCache.Add(models.KeyType(strconv.Itoa(i)), models.ValueType(strconv.Itoa(i)), time.Now().UnixNano()+1000000000*int64(i))
	}
	timeout := time.After(time.Second * 100)
	timeRegular := conf.Config.UpdateTime
	//res := make(chan interface{}, 10000)
	go DoTickerWorkt(timeout, ConCache, timeRegular)

	test()

}
func DoTickerWorkt(timeout <-chan time.Time, cache *concurrency.ConcurrencyCache, ut int64) {
	ticker := time.NewTicker(time.Second * time.Duration(ut))
	//ticker := time.NewTicker(time.Second * 3)
	done := make(chan int, 1)
	go func() {
		//defer close(res)
		i := 1
		for {
			select {
			case <-ticker.C:
				logrus.Info("start ", i, "th worker")
				//功能
				cache.DeleteRegulary()
				//res <- i
				i++
			case <-timeout:
				logrus.Info("timeout", timeout)
				close(done)
				return
			}
		}
		// for range ticker.C {
		// 	i++
		// 	logrus.Info("start ", i, "th worker")
		// 	lru.RemoveExpiredTimeKey()
		// }
		// done <- 1
	}()
	<-done
}

// func DoTickerWorkt(res chan interface{}, timeout <-chan time.Time, lru *strategies.LRUStrategy) {
// 	//ticker := time.NewTicker(time.Second * time.Duration(conf.Config.UpdateTime))
// 	ticker := time.NewTicker(time.Second * 3)
// 	done := make(chan int, 1)
// 	go func() {
// 		defer close(res)
// 		i := 1
// 		for {
// 			select {
// 			case <-ticker.C:
// 				logrus.Info("start ", i, "th worker")
// 				lru.RemoveExpiredTimeKey()
// 				res <- i
// 				i++
// 			case <-timeout:
// 				logrus.Info("timeout", timeout)
// 				close(done)
// 				return
// 			}
// 		}
// 		// for range ticker.C {
// 		// 	i++
// 		// 	logrus.Info("start ", i, "th worker")
// 		// 	lru.RemoveExpiredTimeKey()
// 		// }
// 		// done <- 1
// 	}()
// 	<-done
// }

func test() {
	for i := 0; i < 20; i++ {
		logrus.Info("定时器后 hello " + strconv.Itoa(i))
		time.Sleep(time.Second)
	}
}

// func add(concurrencyCache *concurrency.ConcurrencyCache, key string, value string, wg *sync.WaitGroup) {
// 	concurrencyCache.Add(models.KeyType(key), models.ValueType(value))
// 	wg.Done()
// }

// var ToUseStrategy = flag.String("s", "lru", "调度策略")

// func main() {
// 	logrus.SetLevel(logrus.DebugLevel)
// 	logrus.Info("分布式缓存服务已启动...")
// 	// 解析参数
// 	flag.Parse()
// 	// 输出无用参数：若命令行指定未定义参数为 -b
// 	logrus.Info("无用参数：", flag.Args())
// 	if !conf.Strategy_Map[*ToUseStrategy] {
// 		logrus.Fatal("不存在该策略：", *ToUseStrategy)
// 	}
// 	factory := new(strategies.StrategyFactory)
// 	var mytest strategies.EliminationStrategy = factory.CreateStrategy(*ToUseStrategy)
// 	// mytest := strategy.NewLRUStrategy()
// 	for i := 0; i < 20; i++ {
// 		mytest.Add(models.KeyType(strconv.Itoa(i)), models.ValueType(strconv.Itoa(i)))
// 	}
// 	// fmt.Println()
// 	// value := reflect.ValueOf(mytest).MethodByName("GetCacheList").Call([]reflect.Value{})
// 	// utils.PrintList(value[0].Interface().(*list.List))
// 	// // factory := new(strategy.StrategyFactory)
// 	// // var mytest strategy.EliminationStrategy = factory.CreateStrategy("lru")
// 	// mytest := strategy.NewFIFOStrategy()
// 	// for i := 0; i < 100; i++ {
// 	// 	mytest.Add(models.KeyType(strconv.Itoa(i)), models.ValueType(strconv.Itoa(i)))
// 	// }
// 	// utils.PrintList(*mytest.GetCacheList())
// 	concurrencyCache := &concurrency.ConcurrencyCache{}
// 	var waitGroup sync.WaitGroup
// 	for i := 0; i < 1000; i++ {
// 		waitGroup.Add(1)
// 		go add(concurrencyCache, strconv.Itoa(i), strconv.Itoa(i), &waitGroup)
// 	}
// 	waitGroup.Wait()
// 	time.Sleep(3 * time.Second)
// 	// logrus.Info(concurrencyCache.Get(models.KeyType(strconv.Itoa(99))))
// 	cnt := 0
// 	for i := 0; i < 1000; i++ {
// 		_, ok := concurrencyCache.Get(models.KeyType(strconv.Itoa(i)))
// 		if ok {
// 			cnt++
// 		}
// 	}
// 	logrus.Debug(cnt)
// 	//测试保护分支
// }
