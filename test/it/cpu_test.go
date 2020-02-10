package it

import (
	"bufio"
	"nes-go/pkg/domain"
	"nes-go/pkg/impl"
	"nes-go/pkg/log"
	"nes-go/pkg/mock_domain"
	"os"
	"strconv"
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
		wantFull := scanner.Text()
		line++

		if err := nes.Run1Cycle(); err != nil {
			t.Errorf("failed to run, %v", err)
			return
		}

		gotFull := recorder.String()

		// Pフラグは完全一致しなくてもいいため、分割して比較する
		{
			want := wantFull[0:65]
			got := gotFull[0:65]

			if want != got {
				t.Errorf("wrong value\nline  :%v\ngot :%#v\nwant:%#v\ngot  full:%#v\nwant full:%#v\n", line, got, want, gotFull, wantFull)
				return
			}
		}

		{
			want, err := strconv.ParseInt(wantFull[65:67], 16, 32)
			if err != nil {
				t.Errorf("failed to parse P\nline  :%v\nwant P:%#v\nerror :%#v\n", line, want, err)
				return
			}

			got, err := strconv.ParseInt(gotFull[65:67], 16, 32)
			if err != nil {
				t.Errorf("failed to parse P\nline  :%v\ngot P:%#v\nerror:%#v\n", line, got, err)
				return
			}

			// Pフラグのビット4,5は異なっててよいため、ビット4,5だけ0にして比較
			filteredWant := want & 0xCF
			filteredGot := got & 0xCF
			if filteredWant != filteredGot {
				t.Errorf("bit0-3 and bit6-7 is wrong\nline:%v\ngot  P:%02X => %02X\nwant P:%02X => %02X\ngot  full:%#v\nwant full:%#v\n", line, got, filteredGot, want, filteredWant, gotFull, wantFull)
				return
			}
		}

		{
			want := wantFull[67:]
			got := gotFull[67:]

			if want != got {
				t.Errorf("wrong value\nline:%v\ngot :%#v\nwant:%#v\ngot  full:%#v\nwant full:%#v\n", line, got, want, gotFull, wantFull)
				return
			}
		}
	}
}
