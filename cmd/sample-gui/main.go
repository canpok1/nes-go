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

	m := infra.NewMonitor(
		model.ResolutionWidth,
		model.ResolutionHeight,
		2,
		"nes-go",
	)

	go func() {
		size := 4 * model.ResolutionWidth * model.ResolutionHeight
		pixels := make([]byte, size)
		var color uint8

		for {
			color = (color + 1) & 0xFF
			for i := range pixels {
				pixels[i] = color
			}
			if err := m.Render(pixels); err != nil {
				log.Fatal("error: %v", err)
				break
			}

			log.Debug("color: %v", color)

			time.Sleep(time.Millisecond * 100)
		}

		return
	}()

	err = m.Run()
}
