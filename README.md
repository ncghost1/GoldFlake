# GoldFlake
在终端输入以下命令导入GoldFlake：
```
go get github.com/ncghost1/GoldFlake
```
这是一个突发奇想出来的，非连续毫秒时间戳增量版本的雪花算法~<br>
这是专门针对一个优先级不是很高的分布式 id 需求进行的改造：**增长的 ID 不能让竞争对手发现你每天的业务量**<br>
这个需求对业务很重要的话，那么这个业务堪称“金子”，这是一个“金子”才适用的分布式 ID 生成算法，所以叫做 **GoldFlake**.<br>
GoldFlake 适用于 id 可被用户搜索的，你想要增加一些非连续性来使增长的 ID 不能让竞争对手发现你每天的业务量的场景。<br>
其实我感觉没人会用这个玩意，看个乐子就好了🤣🤣<br>
对于 GoldFlake，你可以有四种使用方法：<br>
第一种是可以像雪花算法一样使用：
```
func main() {
	var workerid uint32 = 1
	Gf, err := GoldFlake.InitGfNode(workerid)
	if err != nil {
		fmt.Println(err)
		return
	}
	uid, err = Gf.Generate()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(uid)
}
```
第二种是生成非连续毫秒时间戳的id，该方式使用Sleep间歇执行生成随机毫秒时间戳增量的代码，具体实现方式请看源码：
```
func main() {
	var workerid uint32 = 0
	var stackSize uint32 = 5
	var Signal int8 = RandProcessSignalDisable
	var chanceNumerator uint64 = 1
	var chanceDenominator uint64 = 2
	var maxTimeOffset uint64 = 5
	Gf, err := GoldFlake.InitGfNode(workerid)
	if err != nil {
		fmt.Println(err)
		return
	}
	GoldFlake.InitRandProcess(stackSize, Signal)
	runtime.GOMAXPROCS(2) // Optional,but need at least "2" to get good performance
  
  // Make sure 'IntervalRandProcess' is always running and does not exit when actually using GoldFlake
  // We use it to continuously generate random millisecond timestamp increments.
	go func() {
		for {
			status, err := GoldFlake.IntervalRandProcess(1, 2, maxTimeOffset, time.Millisecond)
			if err != nil {
				t.Errorf("RandProcess error:%s", err)
			}
			if status == GoldFlake.RandProcessNotReady {
				runtime.Gosched()
			}
		}
	}()
	uid, err := Gf.Generate()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(uid)
}
```
第三种是生成非连续毫秒时间戳的id，该方式使用信号量的方式在每到新毫秒时执行生成随机毫秒时间戳增量的代码，具体实现方式请看源码：
```
func main() {
	var workerid uint32 = 0
	var stackSize uint32 = 5
	var Signal int8 = RandProcessSignalEnable
	var chanceNumerator uint64 = 1
	var chanceDenominator uint64 = 2
	var maxTimeOffset uint64 = 5
	Gf, err := GoldFlake.InitGfNode(workerid)
	if err != nil {
		fmt.Println(err)
		return
	}
	GoldFlake.InitRandProcess(stackSize, Signal)
	runtime.GOMAXPROCS(2) // Need at least 2 to get good performance
  
  // Make sure 'RandProcess' is always running and does not exit when actually using GoldFlake
  // We use it to continuously generate random millisecond timestamp increments.
	go func() {
		for {
			status, err := GoldFlake.RandProcess(1, 2, maxTimeOffset)
			if err != nil {
				t.Errorf("RandProcess error:%s", err)
			}
			if status == GoldFlake.RandProcessNotReady {
				runtime.Gosched()
			}
		}
	}()
	uid, err := Gf.Generate()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(uid)
}
```
首先说明一下第二，第三种的 InitRandProcess,RandProcess 在示例方法中的使用是错误的，实际上我们需要让它能在机器上一直执行，
它才能持续地为我们的 id 生成随机毫秒时间戳偏移量（将 RandValStack 填充满为止），但是这样在 main 中的写法是做不到的。
想要做到让它一直运行，我们可以无限循环，并使用事件驱动编程的方法调用生成 id 函数等等...示例如下：<br>
```
// pseudo code
for {
	go func() {
		for {
			status, err := GoldFlake.RandProcess(1, 2, maxtimeoffset)
			if err != nil {
				b.Errorf("RandProcess error:%s", err)
			}
			if status == GoldFlake.RandProcessNotReady {
				runtime.Gosched()
			}
		}
	}()
	if getGenerateIdRequest() != nil {
		uid, err := Gf.Generate()
		if err != nil {
			fmt.Println(err)
			return
		}
		sendGenerateIdResponse(uid)
	}
}
```
第四种是生成非连续毫秒时间戳的id，该方法是在生成id函数发现来到新毫秒时间戳时调用随机获取时间偏移量函数，和第二，第三种方法区别在于该方法是相当于生成id
和填充保存随机偏移量的栈是同步在同一个函数里的，而第二，第三种方法则是异步填充栈。我个人更推荐用简单的第四种，前面都是花里胡哨的方法...<br>
```
func main() {
	var workerid uint32 = 1
	var stackSize uint32 = 5
	var chanceNumerator uint64 = 1
	var chanceDenominator uint64 = 2
	var maxTimeOffset uint64 = 5
	var Signal int8 = RandProcessSync
	Gf, err := GoldFlake.InitGfNode(workerid)
	if err != nil {
		fmt.Println(err)
		return
	}
	InitRandProcess(stackSize, Signal)
    uid, err := Gf.SyncGenerateAndRand(chanceNumerator, chanceDenominator, maxTimeOffset)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(uid)
}
```
GoldFlake Benchmark 测试：<br>
数据相近的结果均在误差范围之内，不能从相近的数据确定哪个方法更快，而只能考虑是性能接近。实际上除了单核使用 RandProcess 的 4047 ns/op 以外，其它的方法都很接近。<br>
linux(Ubuntu20.04):
```
goos: linux
goarch: amd64
BenchmarkNormalGenerateId-2                      	 4902864	       244 ns/op
BenchmarkGenerateIdWithIntervalRandProcess-2     	 8060004	       161 ns/op
testing: BenchmarkGenerateIdWithIntervalRandProcess-2 left GOMAXPROCS set to 1
BenchmarkGenerateIdWithIntervalRandProcess_2-2   	 7229703	       168 ns/op
BenchmarkGenerateIdWithRandProcess-2             	  338768	      4047 ns/op
testing: BenchmarkGenerateIdWithRandProcess-2 left GOMAXPROCS set to 1
BenchmarkGenerateIdWithRandProcess_2-2           	 6711577	       156 ns/op
BenchmarkSyncGenerateAndRand-2                   	 5840563	       192 ns/op
testing: BenchmarkSyncGenerateAndRand-2 left GOMAXPROCS set to 1
BenchmarkSyncGenerateAndRand_2-2                 	 7214250	       178 ns/op
PASS
ok  	_/root/Gold	16.724s
```
2022/5/15:<br>
~~目前已经移除了 RandValStack 中的 mutex，使用的是无锁方式解决多线程冲突~~。已经重新将 mutex 加回 RandValStack（果然无锁还是不好保证线程安全），
原本使用 mutex 在单线程时若 RandProcess/IntervalRandProcess
未释放锁时切换了 goroutine，会导致生成 id 线程因获取不到锁而阻塞。现在的做法是在 RandValStack 中的 flag 增加了两个标志位，一个用来标志
RandValStack 被 GenerateId 所读写，另一个用来标志 RandValStack 被 RandProcess/IntervalRandProcess 所读写。<br>
除此之外还增加了一个新的函数：SyncGenerateAndRand，同步生成id和生成随机时间偏移量。该函数实现方法与 RandProcess 方案很像，
均是在生成id时来到了新的毫秒时间则调用一次随机获取时间偏移量函数，但是 RandProcess 给的方案是异步的，而这个 SyncGenerateAndRand 是同步的。
理论上它会比异步方案随机性更强。~~但注意使用该函数生成id时，请勿同时多线程使用 Generate 函数生成id，否则可能会导致线程冲突。~~（加回mutex已保证线程安全）<br>
在 RandValStack 被 RandProcess/IntervalRandProcess 所读写时，我们让 GenerateId 继续生成 id，但不进行偏移，
从而不会因为无法读写 RandValStack 而造成阻塞。<br>
在 RandValStack 被 GenerateId 所读写时，我们会返回状态码 RandProcessNotReady(宏，实际值为1) 表示 RandProcess/IntervalRandProcess 目前无法执行，
则我们使用 Gosched() 将 CPU 时间片分配给其他线程。<br>
IntervalRandProcess（非连续性）随机性较弱，因为我们是让 OS ”随缘“执行 IntervalRandProcess，不推荐使用。如果你问我为什么不把它删掉？因为也许可能
会有对随机性要求较弱，而性能要求较高的需求。<br>
而使用 RandProcess 方法对于生成 id 的性能相比 IntervalRandProcess 较低，但是随机性强。当然随机性和我们自定义设置的参数有关，这里所说的随机性高是因为
和 IntervalRandProcess 相比保证了更多的随机时间偏移量生成次数。<br>
新方法 SyncGenerateAndRand 与使用 RandProcess 相比具有更好的随机性，同时性能也更接近传统雪花算法，比起使用 RandProcess 更推荐使用此方法。<br>
要注意这三个方法都会有一种相同的损失，那就是可用id的数量，另外要注意一点本实现和网络上的雪花算法不一样，网络上只利用了41位毫秒时间戳，我们是使用uint64做id，可以利用42位，所以我们原本可用id的基础是可以用大约139年的，所以能够容忍一定损失。什么你跟我说unix时间戳用不了139年？不说139年，如果你的业务id真需要保持60年以上，你为什么不自己写一个新的时间戳啊？(╬▔皿▔)╯<br>
具体可以查看我的个人网站文章：[创造过程](https://www.eririspace.cn/2022/05/12/GoldFlake/)，[劣质のAPI使用文档](https://www.eririspace.cn/2022/05/15/GoldFlake_2/)<br>
虽然和文章的实现有些出入，但是原理是一样的。🍭🍭
