package main

import (
	"nes-go/pkg/infra"
	"nes-go/pkg/log"
	"nes-go/pkg/model"
	"nes-go/pkg/model/cpu"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLogLevel(log.LevelWarn)

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

			// imgs := makeTestImage()

			if imgs != nil {
				err = m.Render(imgs)
				if err != nil {
					log.Fatal("error: %v", err)
					break
				}
			}

			// time.Sleep(time.Millisecond * 100)
		}

		return
	}()

	err = m.Run()
}

func makeTestImage() [][]model.SpriteImage {
	imgs := make([][]model.SpriteImage, (model.ResolutionHeight / model.SpriteHeight))

	ci := 0
	for y := 0; y < (model.ResolutionHeight / model.SpriteHeight); y++ {
		imgs[y] = make([]model.SpriteImage, (model.ResolutionWidth / model.SpriteWidth))
		for x := 0; x < (model.ResolutionWidth / model.SpriteWidth); x++ {
			sprite := model.Sprite([]byte{
				0x66, // 0b 0110 0110
				0x7F, // 0b 0111 1111
				0xFF, // 0b 1111 1111
				0xFF, // 0b 1111 1111
				0xFF, // 0b 1111 1111
				0x7E, // 0b 0111 1110
				0x3C, // 0b 0011 1100
				0x18, // 0b 0001 1000

				0x66, // 0b 0110 0110
				0x5F, // 0b 0101 1111
				0xBF, // 0b 1011 1111
				0xBF, // 0b 1011 1111
				0xFF, // 0b 1111 1111
				0x7E, // 0b 0111 1110
				0x3C, // 0b 0011 1100
				0x18, // 0b 0001 1000
			})
			palette := model.Palette([]byte{0x00, 0x01, 0x02, 0x03})
			img := *sprite.ToSpriteImage(&palette)

			// r := make([][]byte, model.SpriteHeight)
			// g := make([][]byte, model.SpriteHeight)
			// b := make([][]byte, model.SpriteHeight)

			// for sy := 0; sy < model.SpriteHeight; sy++ {
			// 	r[sy] = make([]byte, model.SpriteWidth)
			// 	g[sy] = make([]byte, model.SpriteWidth)
			// 	b[sy] = make([]byte, model.SpriteWidth)
			// 	for sx := 0; sx < model.SpriteWidth; sx++ {
			// 		switch ci {
			// 		case 0:
			// 			r[sy][sx] = 0xFF
			// 		case 1:
			// 			g[sy][sx] = 0xFF
			// 		case 2:
			// 			b[sy][sx] = 0xFF
			// 		}
			// 	}
			// }

			// img := model.SpriteImage{
			// 	R: r,
			// 	G: g,
			// 	B: b,
			// }

			imgs[y][x] = img

			ci = (ci + 1) % 3
		}
	}

	return imgs
}
