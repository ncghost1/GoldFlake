package GoldFlake

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestNormalGenerateId(t *testing.T) {
	baseNumGoroutine := runtime.NumGoroutine()
	var workerid uint32 = 1
	Gf, err := InitGfNode(workerid)
	count := 0
	done := make(chan bool)
	if err != nil {
		t.Errorf("Init GoldFlake node error:%s", err)
		return
	}

	go func() {
		var prev uint64 = 0
		for {
			select {
			case <-done:
				return
			default:
				cur, err := Gf.Generate()
				if err != nil {
					t.Errorf("Generate id error:%s", err)
					return
				}
				if cur <= prev {
					t.Errorf("Generate id error:The current id is less than or equal to the previous id.prev:%d,cur:%d", prev, cur)
				}
				prev = cur
				count++
			}
		}
	}()

	time.Sleep(time.Second)
	close(done)
	fmt.Println("TestNormalGenerateId: Number of generated ID:", count)

	for {
		if baseNumGoroutine == runtime.NumGoroutine() {
			break
		}
	}
}

func TestGenerateIdWithIntervalRandProcess(t *testing.T) {
	baseNumGoroutine := runtime.NumGoroutine()
	var workerId uint32 = 1
	var maxTimeOffset uint64 = 5
	var stackSize uint32 = 5
	var mode int8 = RandProcessSignalDisable
	coreNum := 1
	Randcnt := 0
	count := 0
	runtime.GOMAXPROCS(coreNum)
	Gf, err := InitGfNode(workerId)
	done := make(chan bool)
	if err != nil {
		t.Errorf("Create Goldflake node error:%s", err)
		return
	}
	err = InitRandProcess(stackSize, mode)
	if err != nil {
		t.Errorf("initialize RandValStack error:%s", err)
	}

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				status, err := IntervalRandProcess(1, 2, maxTimeOffset, time.Millisecond)
				if err != nil {
					t.Errorf("RandProcess error:%s", err)
				}
				if status == RandProcessNotReady {
					runtime.Gosched()
				}
				Randcnt++
			}

		}
	}()

	go func() {
		var prev uint64 = 0
		for {
			select {
			case <-done:
				return
			default:
				cur, err := Gf.Generate()
				if err != nil {
					t.Errorf("GenerateId error:%s", err)
					return
				}
				if cur <= prev {
					t.Errorf("Generate id error:The current id is less than or equal to the previous id.prev:%d,cur:%d", prev, cur)
				}
				prev = cur
				count++
			}
		}
	}()
	time.Sleep(time.Second)
	close(done)
	fmt.Println("TestGenerateIdWithIntervalRandProcess: Number of generated ID:", count)
	fmt.Println("TestGenerateIdWithIntervalRandProcess: IntervalRandProcess Execution Count:", Randcnt)

	for {
		// wait for goroutine exit
		if baseNumGoroutine == runtime.NumGoroutine() {
			break
		}
	}
}

func TestGenerateIdWithIntervalRandProcess_2(t *testing.T) {
	baseNumGoroutine := runtime.NumGoroutine()
	var workerId uint32 = 1
	var maxTimeOffset uint64 = 5
	var stackSize uint32 = 5
	var mode int8 = RandProcessSignalDisable
	coreNum := 2
	Randcnt := 0
	count := 0
	runtime.GOMAXPROCS(coreNum)
	Gf, err := InitGfNode(workerId)
	done := make(chan bool)
	if err != nil {
		t.Errorf("Create Goldflake node error:%s", err)
		return
	}
	err = InitRandProcess(stackSize, mode)
	if err != nil {
		t.Errorf("initialize RandValStack error:%s", err)
	}

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				status, err := IntervalRandProcess(1, 2, maxTimeOffset, time.Millisecond)
				if err != nil {
					t.Errorf("RandProcess error:%s", err)
				}
				if status == RandProcessNotReady {
					runtime.Gosched()
				}
				Randcnt++
			}
		}
	}()

	go func() {
		var prev uint64 = 0
		for {
			select {
			case <-done:
				return
			default:
				cur, err := Gf.Generate()
				if err != nil {
					t.Errorf("GenerateId error:%s", err)
					return
				}
				if cur <= prev {
					t.Errorf("Generate id error:The current id is less than or equal to the previous id.prev:%d,cur:%d", prev, cur)
				}
				prev = cur
				count++
			}
		}
	}()

	time.Sleep(time.Second)
	close(done)
	fmt.Println("TestGenerateIdWithIntervalRandProcess_2: Number of generated ID:", count)
	fmt.Println("TestGenerateIdWithIntervalRandProcess_2: IntervalRandProcess Execution Count:", Randcnt)

	for {
		if baseNumGoroutine == runtime.NumGoroutine() {
			break
		}
	}
}

func TestGenerateIdWithRandProcess(t *testing.T) {
	baseNumGoroutine := runtime.NumGoroutine()
	var workerId uint32 = 1
	var maxTimeOffset uint64 = 5
	var stackSize uint32 = 5
	var mode int8 = RandProcessSignalEnable
	coreNum := 1
	count := 0
	runtime.GOMAXPROCS(coreNum)
	Gf, err := InitGfNode(workerId)
	done := make(chan bool)
	if err != nil {
		t.Errorf("Create Goldflake node error:%s", err)
		return
	}
	err = InitRandProcess(stackSize, mode)
	if err != nil {
		t.Errorf("initialize RandValStack error:%s", err)
	}

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				status, err := RandProcess(1, 2, maxTimeOffset)
				if err != nil {
					t.Errorf("RandProcess error:%s", err)
				}
				if status == RandProcessNotReady {
					runtime.Gosched()
				}
			}
		}
	}()

	go func() {
		var prev uint64 = 0
		for {
			select {
			case <-done:
				return
			default:
				cur, err := Gf.Generate()
				if err != nil {
					t.Errorf("GenerateId error:%s", err)
					return
				}
				if cur <= prev {
					t.Errorf("Generate id error:The current id is less than or equal to the previous id.prev:%d,cur:%d", prev, cur)
				}
				prev = cur
				count++
			}
		}
	}()

	time.Sleep(time.Second)
	close(done)
	fmt.Println("TestGenerateIdWithRandProcess: Number of generated ID:", count)

	for {
		if baseNumGoroutine == runtime.NumGoroutine() {
			break
		}
	}
}

// RandProcess can get better performance when at least two cores are used for parallel computing
// We use GOMAXPROCS(2) in this test.
func TestGenerateIdWithRandProcess_2(t *testing.T) {
	baseNumGoroutine := runtime.NumGoroutine()
	var workerId uint32 = 1
	var maxTimeOffset uint64 = 5
	var stackSize uint32 = 5
	var mode int8 = RandProcessSignalEnable
	coreNum := 2
	count := 0
	runtime.GOMAXPROCS(coreNum)
	Gf, err := InitGfNode(workerId)
	done := make(chan bool)
	if err != nil {
		t.Errorf("Create Goldflake node error:%s", err)
		return
	}
	err = InitRandProcess(stackSize, mode)
	if err != nil {
		t.Errorf("initialize RandValStack error:%s", err)
	}

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				status, err := RandProcess(1, 2, maxTimeOffset)
				if err != nil {
					t.Errorf("RandProcess error:%s", err)
				}
				if status == RandProcessNotReady {
					runtime.Gosched()
				}
			}
		}
	}()

	go func() {
		var prev uint64 = 0
		for {
			select {
			case <-done:
				return
			default:
				cur, err := Gf.Generate()
				if err != nil {
					t.Errorf("GenerateId error:%s", err)
					return
				}
				if cur <= prev {
					t.Errorf("Generate id error:The current id is less than or equal to the previous id.prev:%d,cur:%d", prev, cur)
				}
				prev = cur
				count++
			}
		}
	}()

	time.Sleep(time.Second)
	close(done)
	fmt.Println("TestGenerateIdWithRandProcess_2: Number of generated ID:", count)

	for {
		if baseNumGoroutine == runtime.NumGoroutine() {
			break
		}
	}
}

func TestSyncGenerateAndRand(t *testing.T) {
	baseNumGoroutine := runtime.NumGoroutine()
	var workerid uint32 = 1
	var maxTimeOffset uint64 = 5
	var stackSize uint32 = 5
	var mode int8 = RandProcessSync
	Gf, err := InitGfNode(workerid)
	coreNum := 1
	count := 0
	runtime.GOMAXPROCS(coreNum)
	done := make(chan bool)
	if err != nil {
		t.Errorf("Init GoldFlake node error:%s", err)
		return
	}
	err = InitRandProcess(stackSize, mode)
	if err != nil {
		t.Errorf("initialize RandValStack error:%s", err)
	}

	go func() {
		var prev uint64 = 0
		for {
			select {
			case <-done:
				return
			default:
				cur, err := Gf.SyncGenerateAndRand(1, 2, maxTimeOffset)
				if err != nil {
					t.Errorf("Generate id error:%s", err)
					return
				}
				if cur <= prev {
					t.Errorf("Generate id error:The current id is less than or equal to the previous id.prev:%d,cur:%d", prev, cur)
				}
				prev = cur
				count++
			}
		}
	}()

	time.Sleep(time.Second)
	close(done)
	fmt.Println("TestSyncGenerateAndRand: Number of generated ID:", count)

	for {
		if baseNumGoroutine == runtime.NumGoroutine() {
			break
		}
	}
}

// In fact, we won't get better performance in multi-core,
// but we still do a comparison test with other.
func TestSyncGenerateAndRand_2(t *testing.T) {
	baseNumGoroutine := runtime.NumGoroutine()
	var workerid uint32 = 1
	var maxTimeOffset uint64 = 5
	var stackSize uint32 = 5
	var mode int8 = RandProcessSync
	Gf, err := InitGfNode(workerid)
	coreNum := 2
	count := 0
	runtime.GOMAXPROCS(coreNum)
	done := make(chan bool)
	if err != nil {
		t.Errorf("Init GoldFlake node error:%s", err)
		return
	}
	err = InitRandProcess(stackSize, mode)
	if err != nil {
		t.Errorf("initialize RandValStack error:%s", err)
	}

	go func() {
		var prev uint64 = 0
		for {
			select {
			case <-done:
				return
			default:
				cur, err := Gf.SyncGenerateAndRand(1, 2, maxTimeOffset)
				if err != nil {
					t.Errorf("Generate id error:%s", err)
					return
				}
				if cur <= prev {
					t.Errorf("Generate id error:The current id is less than or equal to the previous id.prev:%d,cur:%d", prev, cur)
				}
				prev = cur
				count++
			}
		}
	}()

	time.Sleep(time.Second)
	close(done)
	fmt.Println("TestSyncGenerateAndRand_2: Number of generated ID:", count)

	for {
		if baseNumGoroutine == runtime.NumGoroutine() {
			break
		}
	}
}

// This test didn't print anything, just testing a mix of all methods to check if the thread is safe.
func TestMixGenerate(t *testing.T) {
	baseNumGoroutine := runtime.NumGoroutine()
	var maxTimeOffset uint64 = 5
	var stackSize uint32 = 5
	var mode int8 = RandProcessSync
	Gf, err := InitGfNode(1)
	Gf_2, err := InitGfNode(2)
	coreNum := 4
	runtime.GOMAXPROCS(coreNum)
	done := make(chan bool)
	if err != nil {
		t.Errorf("Init GoldFlake node error:%s", err)
		return
	}
	err = InitRandProcess(stackSize, mode)
	if err != nil {
		t.Errorf("initialize RandValStack error:%s", err)
	}

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				status, err := RandProcess(1, 2, maxTimeOffset)
				if err != nil {
					t.Errorf("RandProcess error:%s", err)
				}
				if status == RandProcessNotReady {
					runtime.Gosched()
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				status, err := IntervalRandProcess(1, 2, maxTimeOffset, time.Millisecond)
				if err != nil {
					t.Errorf("RandProcess error:%s", err)
				}
				if status == RandProcessNotReady {
					runtime.Gosched()
				}
			}
		}
	}()

	go func() {
		var Gfprev uint64 = 0
		for {
			select {
			case <-done:
				return
			default:
				cur, err := Gf.SyncGenerateAndRand(1, 2, maxTimeOffset)
				if err != nil {
					t.Errorf("Generate id error:%s", err)
					return
				}
				if cur <= Gfprev {
					t.Errorf("Generate id error:The current id is less than or equal to the previous id.Gfprev:%d,cur:%d", Gfprev, cur)
				}
				Gfprev = cur
			}
		}
	}()

	go func() {
		var Gf_2prev uint64 = 0
		for {
			select {
			case <-done:
				return
			default:
				cur, err := Gf_2.Generate()
				if err != nil {
					t.Errorf("GenerateId error:%s", err)
					return
				}
				if cur <= Gf_2prev {
					t.Errorf("Generate id error:The current id is less than or equal to the previous id.Gf_2prev:%d,cur:%d", Gf_2prev, cur)
				}
				Gf_2prev = cur
			}
		}
	}()

	time.Sleep(time.Second)
	close(done)

	for {
		if baseNumGoroutine == runtime.NumGoroutine() {
			break
		}
	}
}
