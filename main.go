package main

import (
	"image/color"
	"log"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// ADSRState describes the state of the ADSR envelope
type ADSRState int

const (
	// ADSRStateIdle idle state
	ADSRStateIdle ADSRState = iota
	// ADSRStateAttack attack state
	ADSRStateAttack
	// ADSRStateDecay decay state
	ADSRStateDecay
	// ADSRStateSustain sustain state
	ADSRStateSustain
	// ADSRStateRelease release state
	ADSRStateRelease
)

// ADSR envelope
type ADSR struct {
	segmentDuration  int64
	segmentRamp      float64
	segmentAmplitude float64
	segmentIncrement float64
	segmentShape     float64
	state            ADSRState
	Gated            bool
	Amplitude        float64
	AttackDuration   int64
	AttackShape      float64
	DecayDuration    int64
	DecayShape       float64
	DecayLevel       float64
	SustainDuration  int64
	ReleaseDuration  int64
	ReleaseShape     float64
}

// NewADSR creates a new adsr
func NewADSR(gated bool, attackDuration int64, attackShape float64, decayDuration int64, decayShape float64, decayLevel float64, sustainDuration int64, releaseDuration int64, releaseShape float64) *ADSR {
	return &ADSR{
		Gated:           gated,
		AttackDuration:  attackDuration,
		AttackShape:     attackShape,
		DecayDuration:   decayDuration,
		DecayShape:      decayShape,
		DecayLevel:      decayLevel,
		SustainDuration: sustainDuration,
		ReleaseDuration: releaseDuration,
		ReleaseShape:    releaseShape,
	}
}

// Gate on or off
func (adsr *ADSR) Gate(on bool) {
	if on {
		adsr.state = ADSRStateAttack
		adsr.segmentShape = adsr.AttackShape
		adsr.segmentDuration = adsr.AttackDuration
		adsr.segmentAmplitude = 1.0 - adsr.Amplitude
		adsr.segmentIncrement = 1.0 / float64(adsr.segmentDuration)
		adsr.segmentRamp = 0.0
	} else if adsr.Gated && adsr.state > ADSRStateIdle && adsr.state < ADSRStateRelease {
		adsr.state = ADSRStateRelease
		adsr.segmentShape = adsr.ReleaseShape
		adsr.segmentDuration = adsr.ReleaseDuration
		adsr.segmentAmplitude = -adsr.Amplitude
		adsr.segmentIncrement = 1.0 / float64(adsr.segmentDuration)
		adsr.segmentRamp = 0.0
	}
}

// Step to next sample
func (adsr *ADSR) Step() float64 {
	if adsr.state == ADSRStateIdle || (adsr.state == ADSRStateSustain && adsr.Gated) {
		return adsr.Amplitude
	}

	adsr.segmentRamp += adsr.segmentIncrement
	adsr.segmentDuration--

	// Output is base amplitude + segment amplitude * shape ramp
	output := adsr.Amplitude + adsr.segmentAmplitude*math.Pow(adsr.segmentRamp, adsr.segmentShape)

	// Clip output between 0 and 1
	if output < 0.0 {
		output = 0.0
	}

	if output > 1.0 {
		output = 1.0
	}

	// Check for end of segment
	if adsr.segmentDuration <= 0 {
		// Reset ramp
		adsr.segmentRamp = 0.0

		// Store base amplitude for next segment
		adsr.Amplitude = output

		// Update state
		adsr.state++

		if adsr.state == ADSRStateDecay {
			adsr.segmentDuration = adsr.DecayDuration
			adsr.segmentShape = adsr.DecayShape
			adsr.segmentAmplitude = -(adsr.Amplitude - adsr.DecayLevel)
			adsr.segmentIncrement = 1.0 / float64(adsr.segmentDuration)
		} else if adsr.state == ADSRStateSustain {
			adsr.segmentDuration = adsr.SustainDuration
			adsr.segmentShape = 1.0
			adsr.segmentAmplitude = 0.0
			adsr.segmentIncrement = 0.0
		} else if adsr.state == ADSRStateRelease {
			adsr.segmentDuration = adsr.ReleaseDuration
			adsr.segmentShape = adsr.ReleaseShape
			adsr.segmentAmplitude = -adsr.Amplitude
			adsr.segmentIncrement = 1.0 / float64(adsr.segmentDuration)
		} else {
			adsr.state = ADSRStateIdle
			adsr.segmentDuration = 0
			adsr.segmentShape = 1.0
			adsr.segmentAmplitude = 0.0
			adsr.segmentIncrement = 0.0
			adsr.Amplitude = 0.0
		}
	}

	return output
}

func main() {
	adsr := NewADSR(false, 10, 0.5, 6, 0.8, 0.4, 30, 30, 0.8)

	xys := make(plotter.XYs, 120)

	var i int

	adsr.Gate(true)

	for i < 120 {
		currentValue := adsr.Step()

		log.Printf("i = %v - currentValue %v\n", i, currentValue)

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

	err = p.Save(500, 200, "/Users/almerlucke/Desktop/test.png")
	if err != nil {
		log.Panic(err)
	}
}
