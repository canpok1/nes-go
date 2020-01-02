package main

import (
	"os"

	"github.com/canpok1/nes-go/pkg/log"
	"github.com/canpok1/nes-go/pkg/model"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLogLevel(log.LevelInfo)

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

	cpu := model.NewCPU(rom)
	err = cpu.Run()
	if err != nil {
		return
	}
}
