package main

import (
	"nes-go/pkg/infra"
	"nes-go/pkg/model"
	"os"

	"github.com/canpok1/nes-go/pkg/log"
)

func main() {
	log.SetOutput(os.Stdout)

	log.Debug("========================================")
	log.Debug("program start")
	log.Debug("========================================")

	m := infra.NewMonitor(
		model.ResolutionWidth,
		model.ResolutionHeight,
		2,
		"Hello, World!",
	)
	if err := m.Run(); err != nil {
		log.Fatal("error:%#v", err)
	}

	log.Debug("========================================")
	log.Debug("program end")
	log.Debug("========================================")
}
