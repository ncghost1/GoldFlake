# GoldFlake
这是一个突发奇想出来的，非连续毫秒时间戳增量版本的雪花算法~<br>
这是专门针对一个优先级不是很高的分布式 id 需求进行的改造：**增长的 ID 不能让竞争对手发现你每天的业务量**<br>
这个需求对业务很重要的话，那么这个业务堪称“金子”，这是一个“金子”才适用的分布式 ID 生成算法，所以叫做 **GoldFlake**.<br>
2022/5/14:已将 mutex 从 RandValStack 中移除，改成无锁方式解决 RandValStack 的多线程冲突，在单核情况下
按照示例使用 Gosched() 在 RandProcess/IntervalRandProcess 无法读写 RandValStack 时将 CPU 时间片分给其他线程，
即不会造成生成 id 线程阻塞。🍭🍭<br>
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
数据相近的结果均在误差范围之内，不能从相近的数据确定哪个方法更快，而只能考虑是性能接近。实际上除了 3924 ns/op 以外，其它的方法都很接近。<br>
linux(Ubuntu20.04):
```
goos: linux
goarch: amd64
BenchmarkNormalGenerateId-2                      	 4923379	       244 ns/op
BenchmarkGenerateIdWithIntervalRandProcess-2     	 8013470	       170 ns/op
testing: BenchmarkGenerateIdWithIntervalRandProcess-2 left GOMAXPROCS set to 1
BenchmarkGenerateIdWithIntervalRandProcess_2-2   	 8144731	       154 ns/op
BenchmarkGenerateIdWithRandProcess-2             	  344252	      3924 ns/op
testing: BenchmarkGenerateIdWithRandProcess-2 left GOMAXPROCS set to 1
BenchmarkGenerateIdWithRandProcess_2-2           	 6750685	       154 ns/op
BenchmarkSyncGenerateAndRand-2                   	 8289258	       206 ns/op
testing: BenchmarkSyncGenerateAndRand-2 left GOMAXPROCS set to 1
BenchmarkSyncGenerateAndRand_2-2                 	 5737431	       249 ns/op
PASS
ok  	_/root/Gold	18.443s
```
~~为什么Benchmark的结果显示使用了 IntervalRandProcess（上面列出的第二种方法）性能比不加偏移量（传统Snowflake）更高？
我认为第一个也许只是测试结果的误差，实际上是差不多的。第二个，我认为理论上它确实能生成更多的id，首先 sleep 1ms 是不精确的，
它们不能精确做到每 ms 执行一次获取随机毫秒偏移量并填充栈的函数(fillWithRandValStack)，这样可能会在每毫秒时间戳生成一些id后发生时间偏移。
我们用的序列号是12位，每毫秒能生成 4096 个id，假如我们机器的性能能做到每毫秒能生成8000个id，那么在生成4095个id之后发生了时间偏移，
强制跳到另一个毫秒时间戳（是逻辑上跳而不是实际时间发生跳跃），之后新的时间戳的序列号又从0开始生成，我们就可以在一个毫秒内将8000个id全部生成。<br>
然后还需要解释一下，单核情况下的 BenchmarkGenerateIdWithIntervalRandProcess 测试得到 191 ns/op 是不稳定的，
因为我们每次对 RandValStack 进行读写都有加互斥锁，单核运行两个 goroutine
，那么在切换 goroutine 时可能会造成 fillWithRandValStack 函数只执行到了一半尚未释放锁，
所以 GenerateId 因为获取不到锁，本次分配到的 goroutine 执行时间会被一直阻塞没有操作，
如果运气差的话应该是会像单核情况下的 BenchmarkGenerateIdWithRandProcess 测试一样慢。
所以只有在多核情况下，GoldFlake才能真正发挥作用，另外在机器性能高到每秒生成id数量超过序列号最大数量限制的情况下，
理论上 GoldFlake 能够生成比 SnowFlake 更多的id。~~<br>
2022/5/14:<br>
目前已经移除了 RandValStack 中的 mutex，使用的是无锁方式解决多线程冲突，原本使用 mutex 在单线程时若 RandProcess/IntervalRandProcess
未释放锁时切换了 goroutine，会导致生成 id 线程因获取不到锁而阻塞。现在的做法是在 RandValStack 中的 flag 增加了两个标志位，一个用来标志
RandValStack 被 GenerateId 所读写，另一个用来标志 RandValStack 被 RandProcess/IntervalRandProcess 所读写。<br>
除此之外还增加了一个新的函数：SyncGenerateAndRand，同步生成id和生成随机时间偏移量。该函数实现方法与 RandProcess 方案很像，
均是在生成id时来到了新的毫秒时间则调用一次随机获取时间偏移量函数，但是 RandProcess 给的方案是异步的，而这个 SyncGenerateAndRand 是同步的。
理论上它会比异步方案随机性更强，但注意使用该函数生成id时，请勿同时多线程使用 Generate 函数生成id，否则可能会导致线程冲突。<br>
在 RandValStack 被 RandProcess/IntervalRandProcess 所读写时，我们让 GenerateId 继续生成 id，但不进行偏移，
从而不会因为无法读写 RandValStack 而造成阻塞。<br>
在 RandValStack 被 GenerateId 所读写时，我们会返回状态码 RandProcessNotReady(宏，实际值为1) 表示 RandProcess/IntervalRandProcess 目前无法执行，
则我们使用 Gosched() 将 CPU 时间片分配给其他线程。<br>
IntervalRandProcess（非连续性）随机性较弱，因为我们是让 OS ”随缘“执行 IntervalRandProcess，不推荐使用。如果你问我为什么不把它删掉？因为也许可能
会有对随机性要求较弱，而性能要求较高的需求。<br>
而使用 RandProcess 方法对于生成 id 的性能相比 IntervalRandProcess 较低，但是随机性强。当然随机性和我们自定义设置的参数有关，这里所说的随机性高是因为
和 IntervalRandProcess 相比保证了更多的随机时间偏移量生成次数。<br>
新方法 SyncGenerateAndRand 与使用 RandProcess 相比具有更好的随机性，同时性能也更接近传统雪花算法，比起使用 RandProcess 更推荐使用此方法。<br>
要注意这两个方法都会有一种相同的损失，那就是可用id的数量，另外要注意一点本实现和网络上的雪花算法不一样，网络上只利用了41位毫秒时间戳，我们是使用uint64做id，可以利用42位，所以我们原本可用id的基础是可以用大约139年的，所以能够容忍一定损失。什么你跟我说unix时间戳用不了139年？不说139年，如果你的业务id真需要保持60年以上，你为什么不自己写一个新的时间戳啊？(╬▔皿▔)╯<br>
具体原理可以查看我的个人网站文章：https://www.eririspace.cn/2022/05/12/GoldFlake/<br>
虽然和文章的实现有些出入，但是原理是一样的。🍭🍭
