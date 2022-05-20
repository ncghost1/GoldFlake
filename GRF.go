// GoldFlake Recovery File
package GoldFlake

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
)

const (
	// In order to support more systems, we just store in the current path
	grfDefaultPath = "GoldRecovery.grf"

	defaultStoreThreshold uint64 = 200

	GRFEnable uint8 = 1

	GRFDisable uint8 = 0

	defaultGRFEnableConfig uint8 = GRFEnable

	// GRF Update Strategy:
	// FSync(FullSync): Update GRF when GoldFlake TimeOffset is updated.
	// TSync(ThresholdSync): Update GRF when GoldFlake TimeOffset reaches a multiple of GRF_StoreThreshold.
	FSync uint8 = 1

	TSync uint8 = 2

	defaultStrategy uint8 = TSync

	// GRF Format
	FormatALL uint8 = 1

	FormatMAX uint8 = 2

	defaultFormat = FormatMAX

	// GRF Content
	strategyDecl = "S:"

	fSyncDecl = "FSYNC"

	tSyncDecl = "TSYNC"

	formatDecl = "F:"

	formatALLDecl = "ALL"

	formatMAXDecl = "MAX"

	workeridDecl = "W:"

	timeOffsetDecl = "T:"
)

// GoldFlake node recovery file path
var GrfPath string = grfDefaultPath

// GoldFlake recovery Enable (if true),
// The default config is true
var GrfEnable uint8 = defaultGRFEnableConfig

// Update GRF when GoldFlake TimeOffset reaches a multiple of GRF_StoreThreshold
var tSyncThreshold uint64 = defaultStoreThreshold

// Default GRF Strategy is 'tSync'
var GrfStrategy uint8 = defaultStrategy

// Default GRF content format is 'MAX'
var GrfFormat uint8 = defaultFormat

var GrfLastupdatedtimeoffset []uint64 = make([]uint64, MaxWorkId)

var GrfCAS uint32 = 0

func GrfCasLock() bool {
	return atomic.CompareAndSwapUint32(&GrfCAS, 0, 1)
}

func GrfCasUnLock() bool {
	return atomic.CompareAndSwapUint32(&GrfCAS, 1, 0)
}

func Uint32ToString(p uint32) string {
	str := make([]byte, 0)
	for p > 0 {
		str = append(str, byte('0'+(p%10)))
		p /= 10
	}
	for i, j := 0, len(str)-1; i < j; i, j = i+1, j-1 {
		t := str[i]
		str[i] = str[j]
		str[j] = t
	}
	return string(str)
}

func Uint64ToString(p uint64) string {
	str := make([]byte, 0)
	for p > 0 {
		str = append(str, byte('0'+(p%10)))
		p /= 10
	}
	for i, j := 0, len(str)-1; i < j; i, j = i+1, j-1 {
		t := str[i]
		str[i] = str[j]
		str[j] = t
	}
	return string(str)
}

//
func readStrategyFromGRF(scanner *bufio.Scanner) error {
	if scanner.Scan() {
		if scanner.Text() != "S:" {
			return errors.New("GRF error: file content not detected 'strategy'.")
		}
	}
	if scanner.Scan() {
		strategy := scanner.Text()
		if strategy == fSyncDecl {
			GrfStrategy = FSync
		} else if strategy == tSyncDecl {
			GrfStrategy = TSync
		} else {
			return errors.New("GRF error: unknown 'strategy'.")
		}
	} else {
		return errors.New("GRF error: file content error.")
	}
	return nil
}

func readFormatFromGRF(scanner *bufio.Scanner) error {
	if scanner.Scan() {
		if scanner.Text() != "F:" {
			return errors.New("GRF error: file content not detected 'format'.")
		}
	}
	if scanner.Scan() {
		format := scanner.Text()
		if format == formatMAXDecl {
			GrfFormat = FormatMAX
		} else if format == formatALLDecl {
			GrfFormat = FormatALL
		} else {
			return errors.New("GRF error: unknown 'format'.")
		}
	} else {
		return errors.New("GRF error: file content error.")
	}
	return nil
}

func CreateGRF() {
	file, err := os.Create(GrfPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var strategy string
	var format string
	if GrfStrategy == FSync {
		strategy = fSyncDecl
	} else if GrfStrategy == TSync {
		strategy = tSyncDecl
	} else {
		panic("GRF error: unknown 'strategy'.")
	}

	if GrfFormat == FormatALL {
		format = formatALLDecl
	} else if GrfFormat == FormatMAX {
		format = formatMAXDecl
	} else {
		panic("GRF error: unknown 'format'")
	}

	var content string = strategyDecl + "\n" + strategy + "\n" + formatDecl + "\n" + format
	if GrfFormat == FormatMAX {
		content = content + "\n" + timeOffsetDecl + "0"
	}
	_, err = file.WriteString(content)
	if err != nil {
		panic(err)
	}
}

func readTimeOffsetFromGRF(workerId uint32) uint64 {
	for GrfCasLock() {
	}
	defer GrfCasUnLock()
	file, err := os.Open(GrfPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	wid := Uint32ToString(workerId)
	scanner := bufio.NewScanner(file)
	err = readStrategyFromGRF(scanner)
	if err != nil {
		panic(err)
	}
	err = readFormatFromGRF(scanner)
	if err != nil {
		panic(err)
	}

	if GrfFormat == FormatMAX {
		if scanner.Scan() {
			t := scanner.Text()
			if t[0:2] == timeOffsetDecl {
				timeoffset, err := strconv.ParseUint(t[2:], 10, 64)
				if err != nil {
					panic(err)
				}
				return timeoffset
			}
		} else {
			panic("GRF error: file content error.")
		}
	} else if GrfFormat == FormatALL {
		for scanner.Scan() {
			t := scanner.Text()
			if t[0:2] == workeridDecl {
				idx := 2
				for t[idx] != ' ' {
					idx++
				}
				if wid == t[2:idx] {
					idx++
					if t[idx:idx+2] == timeOffsetDecl {
						timeoffset, err := strconv.ParseUint(t[idx+2:], 10, 64)
						if err != nil {
							panic(err)
						}
						return timeoffset
					} else {
						panic("GRF error: file content error.")
					}
				}
			} else {
				panic("GRF error: file content error.")
			}
		}
	} else {
		panic("GRF error: unknown GRF format.")
	}
	return 0
}

func writeTimeOffsetInGRF(workerId uint32, timeOffset uint64) {
	for GrfCasLock() {
	}
	defer GrfCasUnLock()
	wid := Uint32ToString(workerId)
	timeoffset := Uint64ToString(timeOffset)
	lineBytes, err := ioutil.ReadFile(GrfPath)
	var lines []string
	if err != nil {
		fmt.Println(err)
	} else {
		contents := string(lineBytes)
		lines = strings.Split(contents, "\n")
	}
	var newLines []string
	var Updated bool = false
	for _, line := range lines {
		if GrfFormat == FormatMAX {
			isIn, err := regexp.MatchString("^"+timeOffsetDecl+".*", line)
			if err != nil {
				panic(err)
			}
			if isIn {
				oldtimeOffset, err := strconv.ParseUint(line[2:], 10, 64)
				if err != nil {
					panic(err)
				}
				if oldtimeOffset < timeOffset {
					newLines = append(newLines, timeOffsetDecl+timeoffset)
					Updated = true
					continue
				}
			}
		} else if GrfFormat == FormatALL {
			isIn, err := regexp.MatchString("^"+workeridDecl+wid+".*", line)
			if err != nil {
				panic(err)
			}
			if isIn {
				newContent := workeridDecl + wid + " " + timeOffsetDecl + timeoffset
				newLines = append(newLines, newContent)
				Updated = true
				continue
			}
		} else {
			panic("GRF error: unknown GRF format.")
		}
		newLines = append(newLines, line)
	}

	if !Updated && GrfFormat == FormatALL {
		newContent := workeridDecl + wid + " " + timeOffsetDecl + timeoffset
		newLines = append(newLines, newContent)
	}

	file, err := os.OpenFile(GrfPath, os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.WriteString(strings.Join(newLines, "\n"))
	if err != nil {
		panic(err)
	}
}

func writeStrategyInGRF() {
	for GrfCasLock() {
	}
	defer GrfCasUnLock()
	lineBytes, err := ioutil.ReadFile(GrfPath)
	var lines []string
	if err != nil {
		fmt.Println(err)
	} else {
		contents := string(lineBytes)
		lines = strings.Split(contents, "\n")
	}
	var newLines []string
	for idx, line := range lines {
		if idx == 0 {
			isIn, err := regexp.MatchString(strategyDecl, line)
			if err != nil {
				panic(err)
			}
			if !isIn {
				panic("GRF error: file content not detected 'strategy'.")
			}
		}
		if idx == 1 {
			if GrfStrategy == FSync {
				newLines = append(newLines, fSyncDecl)
			} else if GrfStrategy == TSync {
				newLines = append(newLines, tSyncDecl)
			} else {
				panic("GRF error: unknown 'strategy'.")
			}
		} else {
			newLines = append(newLines, line)
		}
	}

	file, err := os.OpenFile(GrfPath, os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.WriteString(strings.Join(newLines, "\n"))
	if err != nil {
		panic(err)
	}
}

func GRFGetTimeOffset(workerId uint32) uint64 {
	timeoffset := readTimeOffsetFromGRF(workerId)
	if GrfStrategy == TSync {
		timeoffset += tSyncThreshold
	}
	return timeoffset
}

func GRFUpdateTimeOffset(workerId uint32, timeOffset uint64) {
	if GrfStrategy == TSync && GrfLastupdatedtimeoffset[workerId]+tSyncThreshold <= timeOffset {
		GrfLastupdatedtimeoffset[workerId] = timeOffset
		writeTimeOffsetInGRF(workerId, timeOffset)
	} else {
		writeTimeOffsetInGRF(workerId, timeOffset)
	}
}

func GRFSetPath(path string) {
	GrfPath = path
}

func GRFSetStrategy(strategy uint8) error {
	if strategy != FSync && strategy != TSync {
		return errors.New("GRFSetStrategy error: unknown 'strategy'.")
	}
	GrfStrategy = strategy
	writeStrategyInGRF()
	return nil
}

// Attention: this function will recreate the GRF!
func GRFSetFormat(format uint8) error {
	if format != FormatALL && format != FormatMAX {
		return errors.New("GRFSetStrategy error: unknown 'format'.")
	}
	GrfFormat = format
	CreateGRF()
	return nil
}

func SetGrfEnable() {
	GrfEnable = GRFEnable
}

func SetGrfDisable() {
	GrfEnable = GRFDisable
}
