package main

import (
	"image/color"
	"log"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"

	"github.com/almerlucke/gosynth/synth"
)

func main() {
	// adsr := synth.NewADSR(false, 10, 0.5, 6, 0.8, 0.4, 30, 30, 0.8)

	osc := &synth.BLOsc{
		Mode: synth.BLOscModeTriangle,
	}

	xys := make(plotter.XYs, 4200)

	var i int

	for i < 4200 {
		currentValue := osc.Generate(5000.0, 0.0, 44100.0)

		xys[i].X = float64(i)
		xys[i].Y = currentValue

		i++
	}

	p, err := plot.New()
	if err != nil {

		log.Panic(err)
	}

	stepMid, err := plotter.NewLine(xys)
	if err != nil {
		log.Panic(err)
	}
	stepMid.LineStyle = draw.LineStyle{Color: color.RGBA{R: 196, B: 128, A: 255}, Width: vg.Points(1)}

	p.Add(stepMid)

	err = p.Save(1500, 200, "/Users/almerlucke/Desktop/test.png")
	if err != nil {
		log.Panic(err)
	}
}
