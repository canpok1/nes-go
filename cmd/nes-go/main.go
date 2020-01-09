package main

import (
	"fmt"
	"nes-go/pkg/infra"
	"nes-go/pkg/log"
	"nes-go/pkg/model"
	"nes-go/pkg/model/cpu"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLogLevel(log.LevelInfo)

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

	var rom *model.ROM
	rom, err = model.FetchROM(romPath)
	if err != nil {
		return
	}

	bus := model.NewBus()
	cpu := cpu.NewCPU()
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
			cycle, err := cpu.Run()
			if err != nil {
				log.Fatal("error: %v", err)
				break
			}

			imgs, err := ppu.Run(cycle * 3)
			if err != nil {
				log.Fatal("error: %v", err)
				break
			}

			if imgs != nil {
				err = m.Render(imgs)
				if err != nil {
					log.Fatal("error: %v", err)
					break
				}
			}
		}

		return
	}()

	err = m.Run()
}
