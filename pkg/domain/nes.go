package domain

import (
	"nes-go/pkg/log"
)

// NES ...
type NES struct {
	Bus      Bus
	CPU      CPU
	PPU      PPU
	Pad1     Pad
	Pad2     Pad
	Renderer Renderer
	Recorder *Recorder
}

// Setup ...
func (n *NES) Setup(p string) error {
	rom, err := FetchROM(p)
	if err != nil {
		return err
	}

	vram := NewVRAM()

	n.Bus.Setup(rom, n.PPU, n.CPU, vram, n.Pad1, n.Pad2)
	n.CPU.SetBus(n.Bus)
	n.PPU.SetBus(n.Bus)

	n.CPU.SetRecorder(n.Recorder)
	n.PPU.SetRecorder(n.Recorder)

	return nil
}

// Run ...
func (n *NES) Run() error {
	go func() {
		defer func() {
			log.Info("process end")
		}()
		for {
			cycle, err := n.CPU.Run()
			if err != nil {
				panic(err)
			}

			screen, err := n.PPU.Run(cycle * 3)
			if err != nil {
				panic(err)
			}

			if screen != nil {
				err = n.Renderer.Render(screen)
				if err != nil {
					panic(err)
				}
			}
			log.Debug(n.Recorder.String())
		}
	}()

	if err := n.Renderer.Run(); err != nil {
		return err
	}

	return nil
}
