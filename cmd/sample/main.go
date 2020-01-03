package main

import (
	"os"
	"time"

	"github.com/canpok1/nes-go/pkg/log"
	"github.com/canpok1/nes-go/pkg/model"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLogLevel(log.LevelDebug)

	log.Debug("start")

	var err error
	defer func() {
		if err != nil {
			log.Fatal("error: %v", err)
		}
		log.Debug("end")
	}()

	romPath := "./test/roms/hello-world/hello-world.nes"

	rom, err := model.FetchROM(romPath)
	if err != nil {
		return
	}

	bus := model.NewBus((*rom).Prgrom)

	cpu := model.NewCPU()
	cpu.SetBus(bus)

	for {
		err = cpu.Run()
		if err != nil {
			return
		}

		time.Sleep(time.Second * 1)
	}
}
