package GoldFlake

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const (
	testPath = "test.grf"
)

// Simple ASSERT function to check if a and b are equal
func ASSERT_EQUAL(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Assert error:%v is not equal to %v", a, b)
	}
}

func TestGRFSetFormatALL(t *testing.T) {
	SetGrfEnable()
	GRFSetPath(testPath)
	err := GRFSetFormat(FormatALL)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	ASSERT_EQUAL(t, GrfPath, testPath)
	ASSERT_EQUAL(t, GrfFormat, FormatALL)
	file, err := os.Open(testPath)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	err = readStrategyFromGRF(scanner)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	err = readFormatFromGRF(scanner)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	ASSERT_EQUAL(t, GrfFormat, FormatALL)
}

func TestGRFSetFormatMAX(t *testing.T) {
	err := GRFSetFormat(FormatMAX)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	ASSERT_EQUAL(t, GrfPath, testPath)
	ASSERT_EQUAL(t, GrfFormat, FormatMAX)
	file, err := os.Open(testPath)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	err = readStrategyFromGRF(scanner)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	err = readFormatFromGRF(scanner)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	ASSERT_EQUAL(t, GrfFormat, FormatMAX)
}

func TestGRFSetStrategyfSync(t *testing.T) {
	lineBytes, err := ioutil.ReadFile(GrfPath)
	var lines []string
	if err != nil {
		fmt.Println(err)
	} else {
		contents := string(lineBytes)
		lines = strings.Split(contents, "\n")
	}
	err = GRFSetStrategy(FSync)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	ASSERT_EQUAL(t, GrfPath, testPath)
	ASSERT_EQUAL(t, GrfStrategy, FSync)
	newlineBytes, err := ioutil.ReadFile(GrfPath)
	var newlines []string
	if err != nil {
		fmt.Println(err)
	} else {
		contents := string(newlineBytes)
		newlines = strings.Split(contents, "\n")
	}
	if len(newlines) != len(lines) {
		t.Error("error: Incorrect file writing.")
	}
	for i := 0; i < len(newlines); i++ {
		if i == 1 {
			if GrfStrategy == FSync {
				ASSERT_EQUAL(t, newlines[i], fSyncDecl)
				ASSERT_EQUAL(t, lines[i], tSyncDecl)
			} else {
				t.Error("error: Incorrect file writing.")
			}
		} else {
			if newlines[i] != lines[i] {
				t.Error("error: Incorrect file writing.")
			}
		}
	}
}

func TestGRFSetStrategytSync(t *testing.T) {
	lineBytes, err := ioutil.ReadFile(GrfPath)
	var lines []string
	if err != nil {
		fmt.Println(err)
	} else {
		contents := string(lineBytes)
		lines = strings.Split(contents, "\n")
	}
	err = GRFSetStrategy(TSync)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	ASSERT_EQUAL(t, GrfPath, testPath)
	ASSERT_EQUAL(t, GrfStrategy, TSync)
	newlineBytes, err := ioutil.ReadFile(GrfPath)
	var newlines []string
	if err != nil {
		fmt.Println(err)
	} else {
		contents := string(newlineBytes)
		newlines = strings.Split(contents, "\n")
	}
	if len(newlines) != len(lines) {
		t.Error("error: Incorrect file writing.")
	}
	for i := 0; i < len(newlines); i++ {
		if i == 1 {
			if GrfStrategy == TSync {
				ASSERT_EQUAL(t, newlines[i], tSyncDecl)
				ASSERT_EQUAL(t, lines[i], fSyncDecl)
			} else {
				t.Error("error: Incorrect file writing.")
			}
		} else {
			if newlines[i] != lines[i] {
				t.Error("error: Incorrect file writing.")
			}
		}
	}
}

func TestGRFUpdateTimeOffsetWithFormatALL(t *testing.T) {
	err := GRFSetFormat(FormatALL)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	GRFUpdateTimeOffset(1, 1)
	GRFUpdateTimeOffset(2, 2)
	GRFUpdateTimeOffset(3, 3)
	GRFUpdateTimeOffset(1, 5)
	w1_Offset := GRFGetTimeOffset(1) - tSyncThreshold
	w2_Offset := GRFGetTimeOffset(2) - tSyncThreshold
	w3_Offset := GRFGetTimeOffset(3) - tSyncThreshold
	ASSERT_EQUAL(t, w1_Offset, uint64(5))
	ASSERT_EQUAL(t, w2_Offset, uint64(2))
	ASSERT_EQUAL(t, w3_Offset, uint64(3))
}

func TestGRFUpdateTimeOffsetWithFormatMAX(t *testing.T) {
	err := GRFSetFormat(FormatMAX)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	GRFUpdateTimeOffset(1, 1)
	GRFUpdateTimeOffset(2, 2)
	GRFUpdateTimeOffset(3, 10)
	GRFUpdateTimeOffset(1, 5)
	w1_Offset := GRFGetTimeOffset(1) - tSyncThreshold
	ASSERT_EQUAL(t, w1_Offset, uint64(10))
}

func TestGRFRecoveryWithFormatALL(t *testing.T) {
	err := GRFSetFormat(FormatALL)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	GRFUpdateTimeOffset(1, 1)
	GRFUpdateTimeOffset(2, 2)
	GRFUpdateTimeOffset(3, 3)
	GRFUpdateTimeOffset(1, 5)

	gf1, err := InitGfNode(1)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	gf2, err := InitGfNode(2)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	gf3, err := InitGfNode(3)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	ASSERT_EQUAL(t, gf1.timeOffset, uint64(5)+tSyncThreshold)
	ASSERT_EQUAL(t, gf2.timeOffset, uint64(2)+tSyncThreshold)
	ASSERT_EQUAL(t, gf3.timeOffset, uint64(3)+tSyncThreshold)
}

func TestGRFRecoveryWithFormatMAX(t *testing.T) {
	err := GRFSetFormat(FormatMAX)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	GRFUpdateTimeOffset(1, 1)
	GRFUpdateTimeOffset(2, 2)
	GRFUpdateTimeOffset(3, 3)
	GRFUpdateTimeOffset(1, 10)

	gf1, err := InitGfNode(1)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	gf2, err := InitGfNode(2)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	gf3, err := InitGfNode(3)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	ASSERT_EQUAL(t, gf1.timeOffset, uint64(10)+tSyncThreshold)
	ASSERT_EQUAL(t, gf2.timeOffset, uint64(10)+tSyncThreshold)
	ASSERT_EQUAL(t, gf3.timeOffset, uint64(10)+tSyncThreshold)
}

func TestGRFUpdateTimeOffsetByGenerate_FormatMAX(t *testing.T) {
	err := GRFSetStrategy(FSync)
	err = GRFSetFormat(FormatMAX)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	gf1, err := InitGfNode(1)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	gf2, err := InitGfNode(2)
	if err != nil {
		t.Errorf("error:%v", err)
	}

	err = initRandValStack(5, RandProcessSignalDisable)
	if err != nil {
		t.Errorf("error:%v", err)
	}

	_, err = gf1.Generate()
	if err != nil {
		t.Errorf("error:%v", err)
	}

	_, err = fillWithRandValStack(1, 1, 10)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	_, err = gf1.Generate()
	if err != nil {
		t.Errorf("error:%v", err)
	}

	_, err = fillWithRandValStack(1, 1, 10)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	_, err = gf2.Generate()
	if err != nil {
		t.Errorf("error:%v", err)
	}
	w1_Offset := GRFGetTimeOffset(1)
	w2_Offset := GRFGetTimeOffset(2)

	var maxtimeoffset uint64
	if gf1.timeOffset >= gf2.timeOffset {
		maxtimeoffset = gf1.timeOffset
	} else {
		maxtimeoffset = gf2.timeOffset
	}

	ASSERT_EQUAL(t, maxtimeoffset, w1_Offset)
	ASSERT_EQUAL(t, maxtimeoffset, w2_Offset)
}

func TestGRFUpdateTimeOffsetByGenerate_FormatALL(t *testing.T) {
	err := GRFSetFormat(FormatALL)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	gf1, err := InitGfNode(1)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	gf2, err := InitGfNode(2)
	if err != nil {
		t.Errorf("error:%v", err)
	}

	err = initRandValStack(5, RandProcessSignalDisable)
	if err != nil {
		t.Errorf("error:%v", err)
	}

	_, err = gf1.Generate()
	if err != nil {
		t.Errorf("error:%v", err)
	}

	_, err = fillWithRandValStack(1, 1, 10)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	_, err = gf1.Generate()
	if err != nil {
		t.Errorf("error:%v", err)
	}

	_, err = fillWithRandValStack(1, 1, 10)
	if err != nil {
		t.Errorf("error:%v", err)
	}
	_, err = gf2.Generate()
	if err != nil {
		t.Errorf("error:%v", err)
	}

	w1_Offset := GRFGetTimeOffset(1)
	w2_Offset := GRFGetTimeOffset(2)
	ASSERT_EQUAL(t, gf1.timeOffset, w1_Offset)
	ASSERT_EQUAL(t, gf2.timeOffset, w2_Offset)
}
