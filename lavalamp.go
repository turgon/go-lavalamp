package main

import (
	ui "github.com/nsf/termbox-go"
	"github.com/turgon/go-perlin/perlin"
	"math/rand"
	"time"
)

type Position struct {
	X, Y, Z float64
}

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	ui.SetOutputMode(ui.OutputNormal)

	p := perlin.New(r.Int63())

	var pos Position

	dimY, dimX := ui.Size()
	dimZ := dimX

	pos.X = 0
	pos.Y = 0
	pos.Z = 0

	events := make(chan ui.Event, 100)
	go func() {
		for {
			events <- ui.PollEvent()
		}
	}()

loop:
	for {
		ui.Clear(ui.ColorDefault, ui.ColorDefault)
		for x := 0; x < dimX; x++ {
			for y := 0; y < dimY; y++ {
				n := p.Noise(pos.X+float64(x)/float64(dimX), pos.Y+float64(y)/float64(dimY), pos.Z)
				if n >= 0.8 {
					ui.SetCell(y, x, '█', ui.ColorWhite, ui.ColorDefault)
				} else if n >= 0.5 {
					ui.SetCell(y, x, '▓', ui.ColorWhite, ui.ColorDefault)
				} else if n >= 0.2 {
					ui.SetCell(y, x, '▒', ui.ColorWhite, ui.ColorDefault)
				} else if n >= 0.0 {
					ui.SetCell(y, x, '░', ui.ColorWhite, ui.ColorDefault)
				}
			}
		}
		ui.Flush()

		select {
		case ev := <-events:
			if ev.Type == ui.EventKey && ev.Ch == 'q' {
				break loop
			}
		case <-time.After(50 * time.Millisecond):
			pos.Z = pos.Z + 1.0/float64(dimZ)
		}
	}
}
