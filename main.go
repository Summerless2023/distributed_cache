package main

import (
	"fmt"
	"main/conf"
	"net/http"

	"github.com/sirupsen/logrus"
)

type server int

func (h *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logrus.Println(r.URL.Path)
	w.Write([]byte("Hello World!"))
}

func main() {
	conf.LoadConfig()                  //进行配置信息的加载
	fmt.Println(conf.Config.Max_Bytes) //在conf中定义全局变量，其他包进行引用
	fmt.Println(conf.Config.Strategy)  //在conf中定义全局变量，其他包进行引用
	//var s server
	//http.ListenAndServe("localhost:9999", &s)
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
