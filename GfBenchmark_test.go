package GoldFlake

import (
	"runtime"
	"testing"
	"time"
)

func BenchmarkNormalGenerateId(b *testing.B) {
	var workerid uint32
	Gf, err := InitGfNode(workerid)
	if err != nil {
		b.Errorf("Init GoldFlake node error:%s", err)
		return
	}
	for i := 0; i < b.N; i++ {
		_, err := Gf.Generate()
		if err != nil {
			b.Errorf("Generate id error:%s", err)
			return
		}
	}
}
func BenchmarkGenerateIdWithIntervalRandProcess(b *testing.B) {
	var workerid uint32 = 1
	var maxtimeoffset uint64 = 5
	var stacksize uint32 = 5
	var useSignal int8 = 1
	coreNum := 1
	runtime.GOMAXPROCS(coreNum)
	Gf, err := New(workerid)
	if err != nil {
		b.Errorf("Create Goldflake node error:%s", err)
		return
	}
	InitRandProcess(stacksize, useSignal)
	b.ResetTimer()
	go func() {
		for {
			err = IntervalRandProcess(1, 2, maxtimeoffset, time.Millisecond)
			if err != nil {
				b.Errorf("RandProcess error:%s", err)
			}
		}
	}()
	for i := 0; i < b.N; i++ {
		_, err := GenerateId(Gf)
		if err != nil {
			b.Errorf("GenerateId error:%s", err)
			return
		}
	}
}

func BenchmarkGenerateIdWithIntervalRandProcess_2(b *testing.B) {
	var workerid uint32 = 1
	var maxtimeoffset uint64 = 5
	var stacksize uint32 = 5
	var useSignal int8 = 1
	coreNum := 2
	runtime.GOMAXPROCS(coreNum)
	Gf, err := New(workerid)
	if err != nil {
		b.Errorf("Create Goldflake node error:%s", err)
		return
	}
	InitRandProcess(stacksize, useSignal)
	b.ResetTimer()
	go func() {
		for {
			err = IntervalRandProcess(1, 2, maxtimeoffset, time.Millisecond)
			if err != nil {
				b.Errorf("RandProcess error:%s", err)
			}
		}
	}()
	for i := 0; i < b.N; i++ {
		_, err := GenerateId(Gf)
		if err != nil {
			b.Errorf("GenerateId error:%s", err)
			return
		}
	}
}

func BenchmarkGenerateIdWithRandProcess(b *testing.B) {
	var workerid uint32 = 1
	var maxtimeoffset uint64 = 5
	var stacksize uint32 = 5
	var useSignal int8 = 1
	coreNum := 1
	runtime.GOMAXPROCS(coreNum)
	Gf, err := New(workerid)
	if err != nil {
		b.Errorf("Create Goldflake node error:%s", err)
		return
	}
	InitRandProcess(stacksize, useSignal)
	b.ResetTimer()
	go func() {
		for {
			err = RandProcess(1, 2, maxtimeoffset)
			if err != nil {
				b.Errorf("RandProcess error:%s", err)
			}
		}
	}()
	for i := 0; i < b.N; i++ {
		_, err := GenerateId(Gf)
		if err != nil {
			b.Errorf("GenerateId error:%s", err)
			return
		}
	}
}

func BenchmarkGenerateIdWithRandProcess_2(b *testing.B) {
	var workerid uint32 = 1
	var maxtimeoffset uint64 = 5
	var stacksize uint32 = 5
	var useSignal int8 = 1
	coreNum := 2
	runtime.GOMAXPROCS(coreNum)
	Gf, err := New(workerid)
	if err != nil {
		b.Errorf("Create Goldflake node error:%s", err)
		return
	}
	InitRandProcess(stacksize, useSignal)
	b.ResetTimer()
	go func() {
		for {
			err = RandProcess(1, 2, maxtimeoffset)
			if err != nil {
				b.Errorf("RandProcess error:%s", err)
			}
		}
	}()
	for i := 0; i < b.N; i++ {
		_, err := GenerateId(Gf)
		if err != nil {
			b.Errorf("GenerateId error:%s", err)
			return
		}
	}
}
