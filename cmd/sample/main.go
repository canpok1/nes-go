package main

import (
	"log"

	"github.com/canpok1/nes-go/pkg/model"
)

func main() {
	log.Printf("start")

	var err error
	defer func() {
		if err != nil {
			log.Printf("error: %v", err)
		}
		log.Printf("end")
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
