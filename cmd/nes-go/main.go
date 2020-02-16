package main

import (
	"nes-go/pkg/domain"
	"nes-go/pkg/impl"
	"nes-go/pkg/log"
	"os"

	"github.com/hajimehoshi/ebiten"
	"golang.org/x/xerrors"
)

const (
	LOGLEVEL           = log.LevelInfo
	SCALE              = 2.5
	ENABLE_DEBUG_PRINT = true
	ENABLE_FUNC_NAME   = false
	FIRST_PC           = 0x0000
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLogLevel(LOGLEVEL)
	log.SetEnableFuncName(ENABLE_FUNC_NAME)

	log.Debug("========================================")
	log.Debug("program start")
	log.Debug("========================================")

	defer func() {
		log.Debug("========================================")
		log.Debug("program end")
		log.Debug("========================================")
	}()

	if len(os.Args) < 2 {
		panic(xerrors.Errorf("failed to start, rom is nil"))
	}

	romPath := os.Args[1]

	log.Info("rom: %v", romPath)

	bus := impl.NewBus()

	var cpu domain.CPU
	if FIRST_PC == 0 {
		cpu = impl.NewCPU(nil)
	} else {
		firstPC := uint16(FIRST_PC)
		cpu = impl.NewCPU(&firstPC)
	}

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

	if err := nes.Setup(romPath); err != nil {
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
