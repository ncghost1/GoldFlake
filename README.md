# GoldFlake
在终端输入以下命令导入GoldFlake：
```
go get github.com/ncghost1/GoldFlake
```
&emsp;&emsp;这是一个突发奇想出来的，非连续毫秒时间戳增量版本的雪花算法~<br>

&emsp;&emsp;这是专门针对一个优先级不是很高的分布式 id 需求进行的改造：**增长的 ID 不能让竞争对手发现你每天的业务量**<br>

&emsp;&emsp;这个需求对业务很重要的话，那么这个业务堪称“金子”，这是一个“金子”才适用的分布式 ID 生成算法，所以叫做 **GoldFlake**.<br>

&emsp;&emsp;GoldFlake 适用于 id 可被用户搜索的，你想要增加一些非连续性来使增长的 ID 不能让竞争对手发现你每天的业务量的场景。<br>

&emsp;&emsp;其实我感觉没人会用这个玩意，看个乐子就好了🤣🤣<br>

&emsp;&emsp;对于 GoldFlake，你可以有四种使用方法：<br>

&emsp;&emsp;第一种是可以像雪花算法一样使用：
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
&emsp;&emsp;第二种是生成非连续毫秒时间戳的id，该方式使用Sleep间歇执行生成随机毫秒时间戳增量的代码，具体实现方式请看源码：
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
&emsp;&emsp;第三种是生成非连续毫秒时间戳的id，该方式使用信号量的方式在每到新毫秒时执行生成随机毫秒时间戳增量的代码，具体实现方式请看源码：
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
&emsp;&emsp;首先说明一下第二，第三种的 InitRandProcess,RandProcess 在示例方法中的使用是错误的，实际上我们需要让它能在机器上一直执行，
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
&emsp;&emsp;第四种是生成非连续毫秒时间戳的id，该方法是在生成id函数发现来到新毫秒时间戳时调用随机获取时间偏移量函数，和第二，第三种方法区别在于该方法是相当于生成id
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
**GoldFlake Benchmark 测试：**<br>

&emsp;&emsp;数据相近的结果均在误差范围之内，不能从相近的数据确定哪个方法更快，而只能考虑是性能接近。实际上除了单核使用 RandProcess 的 4351 ns/op 以外，其它的方法都很接近。<br>

&emsp;&emsp;其次 GoldFlake 可以超过理论上每秒能生成的最大数量（4096000），因为我们会有时间偏移发生，当发生时间偏移时即使实际时间还在同一毫秒，但是逻辑上id毫秒时间戳的部分已经刷新，序列号也已经归0，又可以从0开始生成id。<br>

linux(Ubuntu20.04)(双核2G)(GRF持久化开启 FSync+FormatALL ) 2022/5/20:

```
goos: linux
goarch: amd64
BenchmarkNormalGenerateId-2                      	 4914885	       244 ns/op
BenchmarkGenerateIdWithIntervalRandProcess-2     	 4922476	       255 ns/op
testing: BenchmarkGenerateIdWithIntervalRandProcess-2 left GOMAXPROCS set to 1
BenchmarkGenerateIdWithIntervalRandProcess_2-2   	 4987290	       241 ns/op
BenchmarkGenerateIdWithRandProcess-2             	  343227	      4351 ns/op
testing: BenchmarkGenerateIdWithRandProcess-2 left GOMAXPROCS set to 1
BenchmarkGenerateIdWithRandProcess_2-2           	 4794088	       253 ns/op
BenchmarkSyncGenerateAndRand-2                   	 4879827	       245 ns/op
testing: BenchmarkSyncGenerateAndRand-2 left GOMAXPROCS set to 1
BenchmarkSyncGenerateAndRand_2-2                 	 4923081	       244 ns/op
PASS
ok  	command-line-arguments	10.297s

```
***
2022/5/20:<br>
#### 新内容：
1. 增加了持久化功能GRF。
2. 目前已经使用无锁方式代替 RVStack 中的 mutex 来确保线程安全，大多数情况下无锁方式将开销更少。
#### 为什么需要持久化：
&emsp;&emsp;之前一直都没有在意这件事，因为我们的id毫秒时间戳会发生偏移，而且是随机的，这些数据一直保存在内存中。如果机器宕机重启后，
这些数据会丢失，机器这个时候的毫秒时间戳偏移量又从0开始，如果之前生成的id发生了偏移，是有可能会生成重复id的。为此我们需要一种方法，至少能够让id的毫秒时间戳偏移量大于之前的偏移量，这样才能有序且不重复。<br>
(呜呜本来就不会有人用了，生成个id还需要持久化这下更加没人用了╥﹏╥...)
#### 一些持久化方案：
1. 生成id的机器进行本地持久化，将毫秒时间戳偏移量本地保存。重启恢复的时候较为快速，但会降低生成id性能。

2. 重启恢复的时候通过数据库查表，找出最大（最新）的id，计算取出前42位毫秒时间戳部分，通过与当前毫秒时间戳相减，得出需要的偏移量。
但是这个方法还要注意其它因素，如果id是异步入库的，那么最新的id可能会在消息队列中，需要考虑这种情况做出调整。

3. 将最大（最新）的id放到缓存层中（如Redis），通过缓存查找id比数据库查表的速度通常要快得多。不过要注意更新顺序，如果id先在缓存上更新再到数据库更新，
这个方法才比较安全。
#### GRF:
1. GoldFlake 内部提供了本地持久化方案，默认是开启的。如果要关闭可以使用前调用 SetGrfDisable() 进行关闭，或者修改源码中的 defaultGRFEnableConfig 修改默认值。

2. GRF 提供了两种本地化策略：FSync 和 TSync.

3. FSync(FullSync) 完全同步：<br>
每一个 GoldFlake 节点的 timeoffset 更新时都进行持久化。

4. TSync(ThresholdSync) 阈值式同步（默认）：<br>
每一个 GoldFlake 节点的 timeoffset 超过阈值的倍数时进行持久化（默认阈值为200）。在 TSync 下重启恢复时，恢复的 timeoffset 为本地保存的 timeoffset 再加上 tSyncThreshold（阈值），这样做一定不会比之前的id小，以确保可用。

5. 注意需要持久化时，GRF 的顺序是先写入本地文件，再生成 id，以保证持久化数据一直是最新的。

6. GRF 还提供了两种持久化格式： ALL 和 MAX.

7. 首先两种持久化格式都需要在最开头存储策略（Strategy）和格式（Format）信息，示例如下：
```
S: // Strategy
TSYNC
F: // Format
MAX
```
8. Format ALL：<br>
&emsp;&emsp;ALL 格式会存储所有不同 GoldFlake 节点的 workerid 与 timeoffset 信息，格式如下：
```
S:
FSYNC
F:
ALL
W:1 T:6
W:2 T:2
```
8. Format MAX（默认）：<br>
&emsp;&emsp;MAX格式只存储所有 GoldFlake 节点中最大的 timeoffset 信息，所有 GoldFlake 节点恢复时 timeoffset 将都恢复成保存的最大值（TSync 下还要加上阈值），格式如下：
```
S:
TSYNC
F:
MAX
T:446
```
9. Format ALL 和 FSync 共同的优点是能够利用更多的可用id数量，缺陷都在于性能开销更大。<br>
而 Format MAX 和 TSync 共同的优点都是开销较少，缺陷都是能够利用的id数量更少。因为 TSync 的恢复方式将可能跳过一定的可用id，MAX 格式的恢复也可能会跳过其他 GoldFlake 节点的可用id。<br>

&emsp;&emsp;我们默认是选择了性能优先的组合方案：TSync + Format MAX.

10. 为了兼容不同系统，所以默认路径直接选用了当前路径，默认持久化文件名为"GoldRecovery.grf"。(本来是叫"dump.grf"，但是还是有特色点比较好嘻嘻！)
### 当 timeoffset 很大时 GRF 会发生什么？
&emsp;&emsp;想一想，如果 timeoffset 已经累积很大时，比如达到了1个小时的偏移量会发生什么？<br>

&emsp;&emsp;当我们的 timeoffset 有1个小时的时候，假如机器关机了2个小时重启了，timeoffset 恢复为之前的值，按理来说 timeoffset 此时为 0 也依然大于先前生成的id，而这么恢复我们会损失掉1个小时这么多的可用id数量！那么有办法解决这个问题吗？<br>

&emsp;&emsp;其实之前提到第二，第三种持久化方案就可以做到，取出最大id的前42位毫秒时间戳部分，为了更多利用可用id可以再取中间的机器id部分判断是否是当前要恢复的机器生成的，从而取出该机器生成过的最大id前42位毫秒时间戳部分。我们通过当前时间和上一次生成id的时间做减法即可获得准确需要的偏移量（注意结果为负值时要恒等于0）。<br>

&emsp;&emsp;那么本地持久化可以做到吗？答案是可以的，在每次生成id的时候都做持久化即可，但是这样做性能开销就太大了！所以我并不打算提供这样的持久化策略，另外 timeoffset 累积几个小时我认为还是能容忍的，而如果累积到了以天数来计的大小，能有这样的规模我认为也可以用得了第二或第三种持久化吧？<br>
***
&emsp;&emsp;原本使用 mutex 在单协程时若 RandProcess/IntervalRandProcess
未释放锁时切换了 goroutine，会导致生成 id 协程因获取不到锁而阻塞。现在的做法是在 RandValStack 中的 flag 增加了两个标志位，一个用来标志
RandValStack 被 GenerateId 所读写，另一个用来标志 RandValStack 被 RandProcess/IntervalRandProcess 所读写。<br>

&emsp;&emsp;除此之外还增加了一个新的函数：SyncGenerateAndRand，同步生成id和生成随机时间偏移量。该函数实现方法与 RandProcess 方案很像，
均是在生成id时来到了新的毫秒时间则调用一次随机获取时间偏移量函数，但是 RandProcess 给的方案是异步的，而这个 SyncGenerateAndRand 是同步的。
理论上它会比异步方案随机性（调用随机函数次数更多）更强。<br>

&emsp;&emsp;在 RandValStack 被 RandProcess/IntervalRandProcess 所读写时，我们让 GenerateId 继续生成 id，但不进行偏移，
从而不会因为无法读写 RandValStack 而造成阻塞。<br>

&emsp;&emsp;在 RandValStack 被 GenerateId 所读写时，我们会返回状态码 RandProcessNotReady(宏，实际值为1) 表示 RandProcess / IntervalRandProcess 目前无法执行，
则我们使用 Gosched() 将 CPU 时间片分配给其他协程。<br>

&emsp;&emsp;IntervalRandProcess（非连续性）随机性较弱，因为我们是让 OS ”随缘“执行 IntervalRandProcess，不推荐使用。如果你问我为什么不把它删掉？因为也许可能
会有对随机性要求较弱，而性能要求较高的需求。<br>

&emsp;&emsp;而使用 RandProcess 方法对于生成 id 的性能相比 IntervalRandProcess 较低，但是随机性强。当然随机性和我们自定义设置的参数有关，这里所说的随机性高是因为
和 IntervalRandProcess 相比保证了更多的随机时间偏移量生成次数。<br>

&emsp;&emsp;新方法 SyncGenerateAndRand 与使用 RandProcess 相比具有更好的随机性，同时性能也更接近传统雪花算法，比起使用 RandProcess 更推荐使用此方法。<br>

&emsp;&emsp;要注意这三个方法都会有一种相同的损失，那就是可用id的数量，另外要注意一点本实现和网络上的雪花算法不一样，网络上只利用了41位毫秒时间戳，我们是使用uint64做id，可以利用42位，所以我们原本可用id的基础是可以用大约139年的，所以能够容忍一定损失。什么你跟我说unix时间戳用不了139年？不说139年，如果你的业务id真需要保持60年以上，你为什么不自己写一个新的时间戳啊？(╬▔皿▔)╯<br>

&emsp;&emsp;具体可以查看我的个人网站文章：[创造过程](https://www.eririspace.cn/2022/05/12/GoldFlake/)，[劣质のAPI使用文档](https://www.eririspace.cn/2022/05/15/GoldFlake_2/)<br>
&emsp;&emsp;虽然和文章的实现有些出入，但是原理是一样的。🍭🍭
