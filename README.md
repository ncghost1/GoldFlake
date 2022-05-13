# GoldFlake
这是一个突发奇想出来的，非连续毫秒时间戳增量版本的雪花算法~<br>
这是专门针对一个优先级不是很高的分布式 id 需求进行的改造：**增长的 ID 不能让竞争对手发现你每天的业务量**<br>
这个需求对业务很重要的话，那么这个业务堪称“金子”，这是一个“金子”才适用的分布式 ID 生成算法，所以叫做 **GoldFlake**.<br>
其实我感觉没人会用这个玩意，看个乐子就好了🤣🤣<br>
对于 GoldFlake，你可以有三种使用方法：<br>
第一种是可以像雪花算法一样使用：
```
func main() {
	var workerid uint32 = 1
	Gf, err := Goldflake.InitGfNode(workerid)
	if err != nil {
		fmt.Println(err)
		return
	}
	uid, err := Gf.Generate()
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
	var Signal int8 = 0
	var chanceNumerator uint64 = 1
	var chanceDenominator uint64 = 2
	var maxTimeOffset uint64 = 5
	Gf, err := Goldflake.InitGfNode(workerid)
	if err != nil {
		fmt.Println(err)
		return
	}
	Goldflake.InitRandProcess(stackSize, Signal)
	runtime.GOMAXPROCS(2) // Optional,but need at least "2" to get good performance
  
  // Make sure 'IntervalRandProcess' is always running and does not exit when actually using GoldFlake
  // We use it to continuously generate random millisecond timestamp increments.
	go func() {
		for {
			Goldflake.IntervalRandProcess(chanceNumerator, chanceDenominator, maxTimeOffset, time.Millisecond)
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
	var Signal int8 = 1
	var chanceNumerator uint64 = 1
	var chanceDenominator uint64 = 2
	var maxTimeOffset uint64 = 5
	Gf, err := Goldflake.InitGfNode(workerid)
	if err != nil {
		fmt.Println(err)
		return
	}
	Goldflake.InitRandProcess(stackSize, Signal)
	runtime.GOMAXPROCS(2) // Need at least 2 to get good performance
  
  // Make sure 'RandProcess' is always running and does not exit when actually using GoldFlake
  // We use it to continuously generate random millisecond timestamp increments.
	go func() {
		for {
			Goldflake.RandProcess(chanceNumerator, chanceDenominator, maxTimeOffset)
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
GoldFlake Benchmark 测试：
```
root@iZj6c7ajkft5134m4zbf3vZ:~/Gold# go test GfBenchmark_test.go Goldflake.go -bench=.
goos: linux
goarch: amd64
BenchmarkNormalGenerateId-2                      	 4766114	       246 ns/op
BenchmarkGenerateIdWithIntervalRandProcess-2     	 7750366	       191 ns/op
testing: BenchmarkGenerateIdWithIntervalRandProcess-2 left GOMAXPROCS set to 1
BenchmarkGenerateIdWithIntervalRandProcess_2-2   	 6110949	       175 ns/op
BenchmarkGenerateIdWithRandProcess-2             	  167834	     14319 ns/op
testing: BenchmarkGenerateIdWithRandProcess-2 left GOMAXPROCS set to 1
BenchmarkGenerateIdWithRandProcess_2-2           	 2690144	       402 ns/op
PASS
ok  	command-line-arguments	10.348s
```
首先说明一下第二，第三种的 InitRandProcess,RandProcess 在示例方法中的使用是错误的，实际上我们需要让它能在机器上一直执行，它才能持续地为我们的 id 生成随机毫秒时间戳增量（将 RandValStack 填充满为止），但是这样在main 中的写法是做不到的。想要做到让它一直运行，我们可以无限循环，并使用事件驱动编程的方法调用生成id函数等等...<br>
其次，为什么Benchmark的结果显示使用了 IntervalRandProcess（上面列出的第二种方法）性能比不加偏移量（传统Snowflake）更高？我认为第一个也许只是测试结果的误差，实际上是差不多的。第二个，我认为理论上它确实能生成更多的id，首先 sleep 1ms 是不精确的，它们不能精确做到每 ms 执行一次获取随机毫秒偏移量并填充栈的函数(fillWithRandValStack)，这样可能会在每毫秒时间戳生成一些id后发生时间偏移。我们用的序列号是12位，每毫秒能生成4096个id，假如我们机器的性能能做到每毫秒能生成8000个id，那么在生成4095个id之后发生了时间偏移，强制跳到另一个毫秒时间戳（是逻辑上跳而不是实际时间发生跳跃），之后新的时间戳的序列号又从0开始生成，我们就可以在一个毫秒内将8000个id全部生成。<br>
然后还需要解释一下，单核情况下的 BenchmarkGenerateIdWithIntervalRandProcess 测试得到 191 ns/op 是不稳定的，因为我们每次对 RandValStack 进行读写都有加互斥锁，单核运行两个 goroutine，那么在切换 goroutine 时可能会造成 fillWithRandValStack 函数只执行到了一半尚未释放锁，所以 GenerateId 因为获取不到锁，本次分配到的 goroutine 执行时间会被一直阻塞没有操作，如果运气差的话应该是会像单核情况下的 BenchmarkGenerateIdWithRandProcess 测试一样慢。所以只有在多核情况下，GoldFlake才能真正发挥作用。<br>
具体原理可以查看我的个人网站文章：https://www.eririspace.cn/2022/05/12/GoldFlake/<br>
虽然和文章的实现有些出入，但是原理是一样的。🍭🍭
