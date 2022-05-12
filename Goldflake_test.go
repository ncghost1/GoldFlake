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
	if err != nil {
		t.Errorf("Init GoldFlake node error:%s", err)
		return
	}
	go func() {
		for {
			_, err := Gf.Generate()
			if err != nil {
				t.Errorf("Generate id error:%s", err)
				return
			}
			count++
		}
	}()
	time.Sleep(time.Second)
	fmt.Println("TestNormalGenerateId: Number of generated ID:", count)
}
func TestGenerateIdWithIntervalRandProcess(t *testing.T) {
	var workerId uint32 = 1
	var maxTimeOffset uint64 = 5
	var stackSize uint32 = 5
	var useSignal int8 = 0
	coreNum := 1
	Randcnt := 0
	count := 0
	runtime.GOMAXPROCS(coreNum)
	Gf, err := New(workerId)
	if err != nil {
		t.Errorf("Create Goldflake node error:%s", err)
		return
	}
	InitRandProcess(stackSize, useSignal)
	go func() {
		for {
			err = IntervalRandProcess(1, 2, maxTimeOffset, time.Millisecond)
			if err != nil {
				t.Errorf("RandProcess error:%s", err)
			}
			Randcnt++
		}
	}()
	go func() {
		for {
			_, err := GenerateId(Gf)
			if err != nil {
				t.Errorf("GenerateId error:%s", err)
				return
			}
			count++
		}
	}()
	time.Sleep(time.Second)
	fmt.Println("TestGenerateIdWithIntervalRandProcess: Number of generated ID:", count)
	fmt.Println("TestGenerateIdWithIntervalRandProcess: IntervalRandProcess Execution Count:", Randcnt)
}

func TestGenerateIdWithIntervalRandProcess_2(t *testing.T) {
	var workerId uint32 = 1
	var maxTimeOffset uint64 = 5
	var stackSize uint32 = 5
	var useSignal int8 = 0
	coreNum := 2
	Randcnt := 0
	count := 0
	runtime.GOMAXPROCS(coreNum)
	Gf, err := New(workerId)
	if err != nil {
		t.Errorf("Create Goldflake node error:%s", err)
		return
	}
	InitRandProcess(stackSize, useSignal)
	go func() {
		for {
			err = IntervalRandProcess(1, 2, maxTimeOffset, time.Millisecond)
			if err != nil {
				t.Errorf("RandProcess error:%s", err)
			}
			Randcnt++
		}
	}()
	go func() {
		for {
			_, err := GenerateId(Gf)
			if err != nil {
				t.Errorf("GenerateId error:%s", err)
				return
			}
			count++
		}
	}()
	time.Sleep(time.Second)
	fmt.Println("TestGenerateIdWithIntervalRandProcess_2: Number of generated ID:", count)
	fmt.Println("TestGenerateIdWithIntervalRandProcess_2: IntervalRandProcess Execution Count:", Randcnt)
}

func TestGenerateIdWithRandProcess(t *testing.T) {
	var workerId uint32 = 1
	var maxTimeOffset uint64 = 5
	var stackSize uint32 = 5
	var useSignal int8 = 1
	coreNum := 1
	count := 0
	runtime.GOMAXPROCS(coreNum)
	Gf, err := New(workerId)
	if err != nil {
		t.Errorf("Create Goldflake node error:%s", err)
		return
	}
	InitRandProcess(stackSize, useSignal)
	go func() {
		for {
			err = RandProcess(1, 2, maxTimeOffset)
			if err != nil {
				t.Errorf("RandProcess error:%s", err)
			}
		}
	}()
	go func() {
		for {
			_, err := GenerateId(Gf)
			if err != nil {
				t.Errorf("GenerateId error:%s", err)
				return
			}
			count++
		}
	}()
	time.Sleep(time.Second)
	fmt.Println("TestGenerateIdWithRandProcess: Number of generated ID:", count)
}

// RandProcess can get better performance when at least two cores are used for parallel computing
// We use GOMAXPROCS(2) in this test.
func TestGenerateIdWithRandProcess_2(t *testing.T) {
	var workerId uint32 = 1
	var maxTimeOffset uint64 = 5
	var stackSize uint32 = 5
	var useSignal int8 = 1
	coreNum := 2
	count := 0
	runtime.GOMAXPROCS(coreNum)
	Gf, err := New(workerId)
	if err != nil {
		t.Errorf("Create Goldflake node error:%s", err)
		return
	}
	InitRandProcess(stackSize, useSignal)
	go func() {
		for {
			err = RandProcess(1, 2, maxTimeOffset)
			if err != nil {
				t.Errorf("RandProcess error:%s", err)
			}
		}
	}()
	go func() {
		for {
			_, err := GenerateId(Gf)
			if err != nil {
				t.Errorf("GenerateId error:%s", err)
				return
			}
			count++
		}
	}()
	time.Sleep(time.Second)
	fmt.Println("TestGenerateIdWithRandProcess_2: Number of generated ID:", count)
}
