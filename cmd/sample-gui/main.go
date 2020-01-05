package main

import (
	"nes-go/pkg/infra"
	"nes-go/pkg/model"
	"os"
	"time"

	"github.com/canpok1/nes-go/pkg/log"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLogLevel(log.LevelDebug)

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

	romPath := "./test/roms/hello-world/hello-world.nes"

	rom, err := model.FetchROM(romPath)
	if err != nil {
		return
	}

	bus := model.NewBus()
	cpu := model.NewCPU()
	ppu, err := model.NewPPU()
	if err != nil {
		return
	}

	bus.Setup(rom, ppu)

	cpu.SetBus(bus)
	ppu.SetBus(bus)

	m, err := infra.NewMonitor(
		model.ResolutionWidth,
		model.ResolutionHeight,
		2,
		"nes-go",
	)
	if err != nil {
		return
	}

	go func() {
		for {
			err := cpu.Run()
			if err != nil {
				log.Fatal("error: %v", err)
				break
			}

			imgs, err := ppu.Run()
			if err != nil {
				log.Fatal("error: %v", err)
				break
			}

			err = m.Render(imgs)
			if err != nil {
				log.Fatal("error: %v", err)
				break
			}

			time.Sleep(time.Millisecond * 100)
		}

		return
	}()

	err = m.Run()
}
