package domain

import (
	"golang.org/x/xerrors"
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

	ppuDelayCycle  int
	cpuBeforeCycle int
}

// Setup ...
func (n *NES) Setup(p string) error {
	rom, err := FetchROM(p)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	vram := NewVRAM()

	n.Bus.Setup(rom, n.PPU, n.CPU, vram, n.Pad1, n.Pad2)
	n.CPU.SetBus(n.Bus)
	n.PPU.SetBus(n.Bus)

	n.CPU.SetRecorder(n.Recorder)
	n.PPU.SetRecorder(n.Recorder)

	n.ppuDelayCycle = 7

	return nil
}

// Run1Cycle ...
func (n *NES) Run1Cycle() error {
	defer log.Debug(n.Recorder.String())

	if n.ppuDelayCycle <= 0 {
		screen, err := n.PPU.Run(n.cpuBeforeCycle * 3)
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}

		if screen != nil {
			err = n.Renderer.Render(screen)
			if err != nil {
				return xerrors.Errorf(": %w", err)
			}
		}
	} else {
		n.ppuDelayCycle = n.ppuDelayCycle - n.cpuBeforeCycle
	}

	n.Recorder.Cycle = n.Recorder.Cycle + n.cpuBeforeCycle

	cycle, err := n.CPU.Run()
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	n.cpuBeforeCycle = cycle

	return nil
}

// Run ...
func (n *NES) Run() error {
	go func() {
		defer func() {
			log.Info("process end")
		}()
		for {
			if err := n.Run1Cycle(); err != nil {
				log.Warn("error occured")
				log.Warn("%s", n.Recorder.String())
				panic(err)
			}
		}
	}()

	if err := n.Renderer.Run(); err != nil {
		return xerrors.Errorf(": %w", err)
	}

	return nil
}
