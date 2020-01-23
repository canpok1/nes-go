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
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLogLevel(LOGLEVEL)

	log.Debug("========================================")
	log.Debug("program start")
	log.Debug("========================================")

	defer func() {
		err := recover()
		if err != nil {
			log.Fatal("%+v", err)
		}
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
	cpu := impl.NewCPU()
	ppu, err := impl.NewPPU2()
	if err != nil {
		return
	}
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
	}

	if err := nes.Run(romPath); err != nil {
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
