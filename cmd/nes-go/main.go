package main

import (
	"errors"
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

	defer func() {
		err := recover()
		if err != nil {
			log.Fatal("error:%#v", err)
		}
		log.Debug("========================================")
		log.Debug("program end")
		log.Debug("========================================")
	}()

	if len(os.Args) < 2 {
		panic(fmt.Errorf("failed to start, rom is nil"))
	}

	romPath := os.Args[1]
	log.Info("rom: %v", romPath)

	rom, err := domain.FetchROM(romPath)
	if err != nil {
		panic(err)
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
		defer func() {
			if err := recover(); err != nil {
				switch err.(type) {
				case error:
					log.Fatal("process error: %v", err)
					inner := errors.Unwrap(err.(error))
					for {
						if inner == nil {
							break
						}
						log.Fatal("inner error: %v", inner)
						inner = errors.Unwrap(inner)
					}
				default:
					log.Fatal("process error: %v", err)
				}
			} else {
				log.Info("process end")
			}
		}()
		for {
			cycle, err := cpu.Run()
			if err != nil {
				panic(err)
			}

			screen, err := ppu.Run(cycle * 3)
			if err != nil {
				panic(err)
			}

			if screen != nil {
				err = r.Render(screen)
				if err != nil {
					panic(err)
				}
			}
		}
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
