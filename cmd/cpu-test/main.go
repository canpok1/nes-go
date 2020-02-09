package main

import (
	"nes-go/pkg/domain"
	"nes-go/pkg/impl"
	"nes-go/pkg/log"
	"os"

	"github.com/hajimehoshi/ebiten"
)

const (
	LOGLEVEL           = log.LevelDebug
	SCALE              = 2.5
	ENABLE_DEBUG_PRINT = true
	FIRST_PC           = 0xC000
	ROM_PATH           = "test/roms/cpu-test/nestest.nes"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLogLevel(LOGLEVEL)
	log.SetEnableLevelLabel(false)
	log.SetEnableTimestamp(false)

	bus := impl.NewBus()

	firstPC := uint16(FIRST_PC)
	cpu := impl.NewCPU(&firstPC)
	ppu := impl.NewPPU2()

	renderer, err := impl.NewRenderer(
		SCALE,
		"nes-go",
		ENABLE_DEBUG_PRINT,
	)
	if err != nil {
		return
	}

	nes := domain.NES{
		Bus:      bus,
		CPU:      cpu,
		PPU:      ppu,
		Pad1:     makePad1(),
		Pad2:     makePad2(),
		Renderer: renderer,
		Recorder: &domain.Recorder{},
	}
	if err := nes.Setup(ROM_PATH); err != nil {
		panic(err)
	}

	if err := nes.Run(); err != nil {
		panic(err)
	}
}

func makePad1() domain.Pad {
	return impl.NewPad(map[domain.ButtonType]ebiten.Key{
		domain.ButtonTypeA:      ebiten.KeyA,
		domain.ButtonTypeB:      ebiten.KeyB,
		domain.ButtonTypeSelect: ebiten.KeyShift,
		domain.ButtonTypeStart:  ebiten.KeySpace,
		domain.ButtonTypeUp:     ebiten.KeyUp,
		domain.ButtonTypeDown:   ebiten.KeyDown,
		domain.ButtonTypeLeft:   ebiten.KeyLeft,
		domain.ButtonTypeRight:  ebiten.KeyRight,
	})
}

func makePad2() domain.Pad {
	return impl.NewPad(map[domain.ButtonType]ebiten.Key{})
}
