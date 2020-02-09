package it

import (
	"bufio"
	"nes-go/pkg/domain"
	"nes-go/pkg/impl"
	"nes-go/pkg/log"
	"nes-go/pkg/mock_domain"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

const (
	LOGLEVEL           = log.LevelDebug
	SCALE              = 2.5
	ENABLE_DEBUG_PRINT = true
	FIRST_PC           = 0xC000
	ROM_PATH           = "../roms/cpu-test/nestest.nes"
	SUCCESS_LOG_PATH   = "../roms/cpu-test/nestest.log"
)

func TestCPU(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log.SetLogLevel(LOGLEVEL)
	log.SetEnableLevelLabel(false)
	log.SetEnableTimestamp(false)

	bus := impl.NewBus()

	firstPC := uint16(FIRST_PC)
	cpu := impl.NewCPU(&firstPC)
	ppu := impl.NewPPU2()

	mRenderer := mock_domain.NewMockRenderer(ctrl)
	mRenderer.EXPECT().Run().AnyTimes()

	nes := domain.NES{
		Bus:      bus,
		CPU:      cpu,
		PPU:      ppu,
		Pad1:     mock_domain.NewMockPad(ctrl),
		Pad2:     mock_domain.NewMockPad(ctrl),
		Renderer: mRenderer,
	}
	if err := nes.Setup(ROM_PATH); err != nil {
		t.Errorf("failed to setup; %v", err)
		return
	}

	file, err := os.Open(SUCCESS_LOG_PATH)
	if err != nil {
		t.Errorf("failed to open success log file; %v", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	tester := IOWriterTester{}
	log.SetOutput(&tester)

	if _, err := nes.CPU.Run(); err != nil {
		t.Errorf("failed to run, %v", err)
		return
	}

	var line int = 0
	for scanner.Scan() {
		tester.Buffer = nil
		want := scanner.Text()
		line++

		if _, err := nes.CPU.Run(); err != nil {
			t.Errorf("failed to run, %v", err)
			return
		}

		writeCount := len(tester.Buffer)
		if writeCount != 1 {
			t.Errorf("write count is too many, line:%v, count:%v", line, writeCount)
			return
		}

		got := strings.Trim(tester.Buffer[0], "\n")
		if got != want {
			t.Errorf("wrong value\nline:%v\ngot :%s\nwant:%s\n", line, got, want)
			return
		}
	}
}

// IOWriterTester ...
type IOWriterTester struct {
	Buffer []string
}

// Write ...
func (t *IOWriterTester) Write(b []byte) (int, error) {
	t.Buffer = append(t.Buffer, string(b))
	return len(b), nil
}
