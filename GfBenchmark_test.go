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
	b.ResetTimer()
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
	var mode int8 = RandProcessSignalDisable
	coreNum := 1
	runtime.GOMAXPROCS(coreNum)
	Gf, err := InitGfNode(workerid)
	done := make(chan bool)
	if err != nil {
		b.Errorf("Create Goldflake node error:%s", err)
		return
	}
	err = InitRandProcess(stacksize, mode)
	if err != nil {
		b.Errorf("initialize RandValStack error:%s", err)
	}
	b.ResetTimer()
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				status, err := IntervalRandProcess(1, 2, maxtimeoffset, time.Millisecond)
				if err != nil {
					b.Errorf("RandProcess error:%s", err)
				}
				if status == RandProcessNotReady {
					runtime.Gosched()
				}
			}
		}
	}()
	for i := 0; i < b.N; i++ {
		select {
		case <-done:
			return
		default:
			_, err := Gf.Generate()
			if err != nil {
				b.Errorf("GenerateId error:%s", err)
				return
			}
		}
	}
	close(done)
}

func BenchmarkGenerateIdWithIntervalRandProcess_2(b *testing.B) {
	var workerid uint32 = 1
	var maxtimeoffset uint64 = 5
	var stacksize uint32 = 5
	var mode int8 = RandProcessSignalDisable
	coreNum := 2
	runtime.GOMAXPROCS(coreNum)
	Gf, err := InitGfNode(workerid)
	done := make(chan bool)
	if err != nil {
		b.Errorf("Create Goldflake node error:%s", err)
		return
	}
	err = InitRandProcess(stacksize, mode)
	if err != nil {
		b.Errorf("initialize RandValStack error:%s", err)
	}
	b.ResetTimer()
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				status, err := IntervalRandProcess(1, 2, maxtimeoffset, time.Millisecond)
				if err != nil {
					b.Errorf("RandProcess error:%s", err)
				}
				if status == RandProcessNotReady {
					runtime.Gosched()
				}
			}
		}
	}()
	for i := 0; i < b.N; i++ {
		select {
		case <-done:
			return
		default:
			_, err := Gf.Generate()
			if err != nil {
				b.Errorf("GenerateId error:%s", err)
				return
			}
		}
	}
	close(done)
}

func BenchmarkGenerateIdWithRandProcess(b *testing.B) {
	var workerid uint32 = 1
	var maxtimeoffset uint64 = 5
	var stacksize uint32 = 5
	var mode int8 = RandProcessSignalEnable
	coreNum := 1
	runtime.GOMAXPROCS(coreNum)
	Gf, err := InitGfNode(workerid)
	done := make(chan bool)
	if err != nil {
		b.Errorf("Create Goldflake node error:%s", err)
		return
	}
	err = InitRandProcess(stacksize, mode)
	if err != nil {
		b.Errorf("initialize RandValStack error:%s", err)
	}
	b.ResetTimer()
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				status, err := RandProcess(1, 2, maxtimeoffset)
				if err != nil {
					b.Errorf("RandProcess error:%s", err)
				}
				if status == RandProcessNotReady {
					runtime.Gosched()
				}
			}
		}
	}()
	for i := 0; i < b.N; i++ {
		select {
		case <-done:
			return
		default:
			_, err := Gf.Generate()
			if err != nil {
				b.Errorf("GenerateId error:%s", err)
				return
			}
		}
	}
	close(done)
}

func BenchmarkGenerateIdWithRandProcess_2(b *testing.B) {
	var workerid uint32 = 1
	var maxtimeoffset uint64 = 5
	var stacksize uint32 = 5
	var mode int8 = RandProcessSignalEnable
	coreNum := 2
	runtime.GOMAXPROCS(coreNum)
	Gf, err := InitGfNode(workerid)
	done := make(chan bool)
	if err != nil {
		b.Errorf("Create Goldflake node error:%s", err)
		return
	}
	err = InitRandProcess(stacksize, mode)
	if err != nil {
		b.Errorf("initialize RandValStack error:%s", err)
	}
	b.ResetTimer()
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				status, err := RandProcess(1, 2, maxtimeoffset)
				if err != nil {
					b.Errorf("RandProcess error:%s", err)
				}
				if status == RandProcessNotReady {
					runtime.Gosched()
				}
			}
		}
	}()
	for i := 0; i < b.N; i++ {
		select {
		case <-done:
			return
		default:
			_, err := Gf.Generate()
			if err != nil {
				b.Errorf("GenerateId error:%s", err)
				return
			}
		}
	}
	close(done)
}

func BenchmarkSyncGenerateAndRand(b *testing.B) {
	var workerid uint32 = 1
	var maxtimeoffset uint64 = 5
	var stacksize uint32 = 5
	var mode int8 = RandProcessSync
	coreNum := 1
	runtime.GOMAXPROCS(coreNum)
	Gf, err := InitGfNode(workerid)
	done := make(chan bool)
	if err != nil {
		b.Errorf("Create Goldflake node error:%s", err)
		return
	}
	err = InitRandProcess(stacksize, mode)
	if err != nil {
		b.Errorf("initialize RandValStack error:%s", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		select {
		case <-done:
			return
		default:
			_, err := Gf.SyncGenerateAndRand(1, 2, maxtimeoffset)
			if err != nil {
				b.Errorf("GenerateId error:%s", err)
				return
			}
		}
	}
	close(done)
}

func BenchmarkSyncGenerateAndRand_2(b *testing.B) {
	var workerid uint32 = 1
	var maxtimeoffset uint64 = 5
	var stacksize uint32 = 5
	var mode int8 = RandProcessSync
	coreNum := 2
	runtime.GOMAXPROCS(coreNum)
	Gf, err := InitGfNode(workerid)
	done := make(chan bool)
	if err != nil {
		b.Errorf("Create Goldflake node error:%s", err)
		return
	}
	err = InitRandProcess(stacksize, mode)
	if err != nil {
		b.Errorf("initialize RandValStack error:%s", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		select {
		case <-done:
			return
		default:
			_, err := Gf.SyncGenerateAndRand(1, 2, maxtimeoffset)
			if err != nil {
				b.Errorf("GenerateId error:%s", err)
				return
			}
		}
	}
	close(done)
}
