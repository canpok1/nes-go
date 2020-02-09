package it

import (
	"bufio"
	"nes-go/pkg/domain"
	"nes-go/pkg/impl"
	"nes-go/pkg/log"
	"nes-go/pkg/mock_domain"
	"os"
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

	recorder := &domain.Recorder{}

	nes := domain.NES{
		Bus:      bus,
		CPU:      cpu,
		PPU:      ppu,
		Pad1:     mock_domain.NewMockPad(ctrl),
		Pad2:     mock_domain.NewMockPad(ctrl),
		Renderer: mRenderer,
		Recorder: recorder,
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

	if err := nes.Run1Cycle(); err != nil {
		t.Errorf("failed to run, %v", err)
		return
	}

	var line int = 0
	for scanner.Scan() {
		want := scanner.Text()
		line++

		if err := nes.Run1Cycle(); err != nil {
			t.Errorf("failed to run, %v", err)
			return
		}

		got := recorder.String()
		if got != want {
			t.Errorf("wrong value\nline:%v\ngot :%s\nwant:%s\n", line, got, want)
			return
		}
	}
}
