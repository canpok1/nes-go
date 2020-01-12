package main

import (
	"fmt"
	"nes-go/pkg/domain"
	"nes-go/pkg/impl"
	"nes-go/pkg/log"
	"os"

	"github.com/hajimehoshi/ebiten"
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

	var err error
	defer func() {
		if err != nil {
			log.Fatal("error:%#v", err)
		}
		log.Debug("========================================")
		log.Debug("program end")
		log.Debug("========================================")
	}()

	if len(os.Args) < 2 {
		err = fmt.Errorf("failed to start, rom is nil")
		return
	}

	romPath := os.Args[1]
	log.Info("rom: %v", romPath)

	var rom *domain.ROM
	rom, err = domain.FetchROM(romPath)
	if err != nil {
		return
	}

	bus := impl.NewBus()
	cpu := impl.NewCPU()
	ppu, err := impl.NewPPU()
	if err != nil {
		return
	}
	vram := domain.NewVRAM()

	bus.Setup(rom, ppu, cpu, vram, makePad1(), makePad2())

	cpu.SetBus(bus)
	ppu.SetBus(bus)

	r, err := impl.NewRenderer(
		SCALE,
		"nes-go",
		ENABLE_DEBUG_PRINT,
	)
	if err != nil {
		return
	}

	go func() {
		for {
			cycle, err := cpu.Run()
			if err != nil {
				log.Fatal("error: %v", err)
				break
			}

			tis, sis, err := ppu.Run(cycle * 3)
			if err != nil {
				log.Fatal("error: %v", err)
				break
			}

			if tis != nil && sis != nil {
				err = r.Render(tis, sis)
				if err != nil {
					log.Fatal("error: %v", err)
					break
				}
			}
		}

		return
	}()

	err = r.Run()
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
