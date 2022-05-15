package GoldFlake

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestNormalGenerateId(t *testing.T) {
	var workerid uint32
	Gf, err := InitGfNode(workerid)
	count := 0
	done := make(chan bool)
	if err != nil {
		t.Errorf("Init GoldFlake node error:%s", err)
		return
	}
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				_, err := Gf.Generate()
				if err != nil {
					t.Errorf("Generate id error:%s", err)
					return
				}
				count++
			}
		}
	}()
	time.Sleep(time.Second)
	close(done)
	fmt.Println("TestNormalGenerateId: Number of generated ID:", count)
}

func TestGenerateIdWithIntervalRandProcess(t *testing.T) {
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
		for {
			select {
			case <-done:
				return
			default:
				_, err := Gf.Generate()
				if err != nil {
					t.Errorf("GenerateId error:%s", err)
					return
				}
				count++
			}
		}
	}()
	time.Sleep(time.Second)
	close(done)
	fmt.Println("TestGenerateIdWithIntervalRandProcess: Number of generated ID:", count)
	fmt.Println("TestGenerateIdWithIntervalRandProcess: IntervalRandProcess Execution Count:", Randcnt)
}

func TestGenerateIdWithIntervalRandProcess_2(t *testing.T) {
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
		for {
			select {
			case <-done:
				return
			default:
				_, err := Gf.Generate()
				if err != nil {
					t.Errorf("GenerateId error:%s", err)
					return
				}
				count++
			}
		}
	}()
	time.Sleep(time.Second)
	close(done)
	fmt.Println("TestGenerateIdWithIntervalRandProcess_2: Number of generated ID:", count)
	fmt.Println("TestGenerateIdWithIntervalRandProcess_2: IntervalRandProcess Execution Count:", Randcnt)
}

func TestGenerateIdWithRandProcess(t *testing.T) {
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
		for {
			select {
			case <-done:
				return
			default:
				_, err := Gf.Generate()
				if err != nil {
					t.Errorf("GenerateId error:%s", err)
					return
				}
				count++
			}
		}
	}()
	time.Sleep(time.Second)
	close(done)
	fmt.Println("TestGenerateIdWithRandProcess: Number of generated ID:", count)
}

// RandProcess can get better performance when at least two cores are used for parallel computing
// We use GOMAXPROCS(2) in this test.
func TestGenerateIdWithRandProcess_2(t *testing.T) {
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
		for {
			select {
			case <-done:
				return
			default:
				_, err := Gf.Generate()
				if err != nil {
					t.Errorf("GenerateId error:%s", err)
					return
				}
				count++
			}
		}
	}()
	time.Sleep(time.Second)
	close(done)
	fmt.Println("TestGenerateIdWithRandProcess_2: Number of generated ID:", count)
}

func TestSyncGenerateAndRand(t *testing.T) {
	var workerid uint32
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
		for {
			select {
			case <-done:
				return
			default:
				_, err := Gf.SyncGenerateAndRand(1, 2, maxTimeOffset)
				if err != nil {
					t.Errorf("Generate id error:%s", err)
					return
				}
				count++
			}
		}
	}()
	time.Sleep(time.Second)
	close(done)
	fmt.Println("TestSyncGenerateAndRand: Number of generated ID:", count)
}

// In fact, we won't get better performance in multi-core,
// but we still do a comparison test with other.
func TestSyncGenerateAndRand_2(t *testing.T) {
	var workerid uint32
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
		for {
			select {
			case <-done:
				return
			default:
				_, err := Gf.SyncGenerateAndRand(1, 2, maxTimeOffset)
				if err != nil {
					t.Errorf("Generate id error:%s", err)
					return
				}
				count++
			}
		}
	}()
	time.Sleep(time.Second)
	close(done)
	fmt.Println("TestSyncGenerateAndRand_2: Number of generated ID:", count)
}
