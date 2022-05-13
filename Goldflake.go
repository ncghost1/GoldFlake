// +-------------------------------------------------------+
// | 42 Bit Timestamp | 10 Bit WorkID | 12 Bit Sequence ID |
// +-------------------------------------------------------+
package GoldFlake

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

const (
	// GMT+8 2022-5-12 00:00:00
	epoch int64 = 1652284800000

	numWorkerBits = 10

	numSequenceBits = 12

	MaxWorkId = -1 ^ (-1 << numWorkerBits)

	MaxSequence = -1 ^ (-1 << numSequenceBits)

	RandProcessSignalEnable = 1 << 0

	RandProcessSignalDisable = 1 << 1
)

type GoldFlake struct {
	lastTimestamp uint64
	sequence      uint32
	workerId      uint32
	timeOffset    uint64
	lock          sync.Mutex
}

// Random Value(TimeOffset) Stack
type RandValStack struct {
	RandVal []uint64 // We use slice to simulate stack space
	top     uint32
	Size    uint32
	flag    int8
	lock    sync.Mutex
}

var RVStack RandValStack

// Pack and return UUID
func (gf *GoldFlake) pack() uint64 {
	uuid := (gf.lastTimestamp << (numWorkerBits + numSequenceBits)) | (uint64(gf.workerId) << numSequenceBits) | (uint64(gf.sequence))
	return uuid
}

// Create and init a GoldFlake node
func New(workerId uint32) (*GoldFlake, error) {
	if workerId < 0 || workerId > MaxWorkId {
		return nil, errors.New("invalid worker Id")
	}
	return &GoldFlake{workerId: workerId}, nil
}

// Initialize Random Value Stack
// if 'UseSignal' is 1, we set the 'RandProcessSignalEnable' bit of flag,
// It means that we will use flag to notify whether a new millisecond has arrived.
// Else, we set the 'RandProcessSignalDisable' bit of flag,
// It means that we will not use flag to notify whether a new millisecond has arrived,
// but we will use time.Sleep(), which is not an accurate method,
// because it is affected by many factors such as OS and hardware.
func initRandValStack(Size uint32, UseSignal int8) {
	RVStack.lock.Lock()
	defer RVStack.lock.Unlock()
	RVStack.RandVal = make([]uint64, Size)
	RVStack.top = 0
	RVStack.Size = Size
	RVStack.flag = 0
	if UseSignal == 1 {
		RVStack.flag |= RandProcessSignalEnable
	} else {
		RVStack.flag |= RandProcessSignalDisable
	}
}

// Push random value into the stack with probability.
// We use fractional form to express probability,
// The probability is chanceNumerator / chanceDenominator.
// maxTimeOffset is the max millisecond time offset,
// We randomly pick an uint64 value from 1 to max and push into the stack.
func fillWithRandValStack(chanceNumerator, chanceDenominator, maxTimeOffset uint64) error {
	rand.Seed(time.Now().UnixNano())
	RVStack.lock.Lock()
	defer RVStack.lock.Unlock()
	if RVStack.flag&RandProcessSignalEnable != 0 && RVStack.flag&RandProcessSignalDisable != 0 {
		return errors.New("SignalEnable and SignalDisable are present at the same time")
	}
	if RVStack.flag&RandProcessSignalEnable != 0 || RVStack.flag&RandProcessSignalDisable != 0 {
		if RVStack.top < RVStack.Size && rand.Uint64()%chanceDenominator < chanceNumerator {
			offset := rand.Uint64()%maxTimeOffset + 1
			RVStack.RandVal[RVStack.top] = offset
			RVStack.top++
			RVStack.flag &= ^RandProcessSignalEnable
		} else {
			RVStack.flag &= ^RandProcessSignalEnable
		}
	}
	return nil
}

// Generate and return an UUID
func (gf *GoldFlake) Generate() (uint64, error) {
	gf.lock.Lock()
	defer gf.lock.Unlock()
	RVStack.lock.Lock()
	if RVStack.top > 0 {
		gf.timeOffset += RVStack.RandVal[RVStack.top-1]
		RVStack.top--
	}
	RVStack.lock.Unlock()
	ts := timestamp() + gf.timeOffset
	if ts == gf.lastTimestamp {
		gf.sequence = (gf.sequence + 1) & MaxSequence
		if gf.sequence == 0 {
			ts = gf.waitNextMilli(ts) + gf.timeOffset
		}
	} else {
		// It's a new millisecond.
		// If we use the signal method to remind the execution of RandProcess,
		// then we need to set RandProcessSignalEnable to 1.
		if RVStack.flag&RandProcessSignalDisable == 0 {
			RVStack.flag |= RandProcessSignalEnable
		}
		gf.sequence = 0
	}

	if ts < gf.lastTimestamp {
		return 0, errors.New("invalid system clock")
	}
	gf.lastTimestamp = ts
	return gf.pack(), nil
}

func (gf *GoldFlake) waitNextMilli(ts uint64) uint64 {
	for ts == gf.lastTimestamp {
		time.Sleep(100 * time.Microsecond)
		ts = timestamp()
	}
	return ts
}

// Get the timestamp,
// this is not the actual millisecond timestamp,
// because we subtract the epoch(The timestamp we want to use as the base)
func timestamp() uint64 {
	return uint64(time.Now().UnixNano()/int64(1000000) - epoch)
}

// *****************************
//        GoldFlake API
// *****************************

// Initialize a GoldFlake node
func InitGfNode(workerid uint32) (*GoldFlake, error) {
	GfNode, err := New(workerid)
	if err != nil {
		return nil, err
	}
	return GfNode, nil
}

// InitRandProcess
// "Size" is the RandValStack's Size,
// when "UseSignal" is 1, We will use the method of setting the signal bit
// of the flag to notify RandProcess to execute.
// when "UseSignal" is 0, We don't use flags to notify if the RandProcess need to execute,
// in this case we use IntervalRandProcess, which will use the Sleep function to run the RandProcess in intervals.
func InitRandProcess(Size uint32, UseSignal int8) {
	initRandValStack(Size, UseSignal)
}

// RandProcess
// When we choose to use the signal method to initialize the stack, we use this.
// Note: When using Goldflake, this function needs to keep running independently.
// Please Use a goroutine and loop this after InitRandProcess:
// For example:
//	go func() {
//		for {
//			err := RandProcess(1,2,5)
//			if err != nil {
//				return err
//			}
//		}
//	}
func RandProcess(chanceNumerator, chanceDenominator, maxTimeOffset uint64) error {
	for {
		err := fillWithRandValStack(chanceNumerator, chanceDenominator, maxTimeOffset)
		if err != nil {
			return err
		}
		return nil
	}
}

// IntervalRandProcess
// The parameter "Interval" is the time of Sleep.
// Attention: Sleep is not exact, if you want to use this function please test your own machine.
// We have provided the relevant output "IntervalRandProcess Execution Count" in the "Goldflake_test" file.
// Note: When using Goldflake, this function needs to keep running independently.
// Please Use a goroutine and loop this after InitRandProcess:
// For example:
//	go func() {
//		for {
//			err := IntervalRandProcess(1,2,5,time.Millisecond)
//			if err != nil {
//				return err
//			}
//		}
//	}
func IntervalRandProcess(chanceNumerator, chanceDenominator, maxTimeOffset uint64, Interval time.Duration) error {
	err := fillWithRandValStack(chanceNumerator, chanceDenominator, maxTimeOffset)
	if err != nil {
		return err
	}
	time.Sleep(Interval)
	return nil
}

// Generate an UUID.
func GenerateId(sf *GoldFlake) (uint64, error) {
	uuid, err := sf.Generate()
	if err != nil {
		return 0, err
	}
	return uuid, nil
}
