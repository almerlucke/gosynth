package synth

import "math"

// BLOscMode band limited oscillator mode
type BLOscMode int

const (
	// BLOscModeSine sine mode
	BLOscModeSine BLOscMode = iota
	// BLOscModeSaw saw mode
	BLOscModeSaw
	// BLOscModeSquare square mode
	BLOscModeSquare
	// BLOscModeTriangle triangle mode
	BLOscModeTriangle
)

// BLOsc band limited oscillator
type BLOsc struct {
	lastOutput float64
	phase      float64
	Mode       BLOscMode
}

func (blosc *BLOsc) polyBlep(t float64, dt float64) float64 {
	if t < dt {
		t /= dt
		return t + t - t*t - 1.0
	} else if t > 1.0-dt {
		t = (t - 1.0) / dt
		return t*t + t + t + 1.0
	}

	return 0.0
}

// Generate a sample
func (blosc *BLOsc) Generate(frequency float64, phaseOffset float64, sampleRate float64) float64 {
	dt := math.Abs(frequency) / sampleRate
	t := blosc.phase + phaseOffset
	v := 0.0

	// Fold phase back because of possible overshoot by adding phase offset
	for t >= 1.0 {
		t -= 1.0
	}

	for t < 0.0 {
		t += 1.0
	}

	if blosc.Mode == BLOscModeSine {
		v = math.Sin(t * math.Pi * 2.0)
	} else if blosc.Mode == BLOscModeSaw {
		v = 2.0*t - 1.0
		v -= blosc.polyBlep(t, dt)
	} else {
		// Generate square
		if t < 0.5 {
			v = 1.0
		} else {
			v = -1.0
		}

		v += blosc.polyBlep(t, dt)
		v -= blosc.polyBlep(math.Mod(t+0.5, 1.0), dt)
	}

	if blosc.Mode == BLOscModeTriangle {
		// Use square wave as input, leaky integration
		v = dt*v + (1.0-dt)*blosc.lastOutput
		blosc.lastOutput = v
		// Boost signal with triangle
		v *= 2.0
	}

	blosc.phase += dt

	// Keep phase within 0-1 bounds
	for blosc.phase >= 1.0 {
		blosc.phase -= 1.0
	}

	for blosc.phase < 0.0 {
		blosc.phase += 1.0
	}

	return v
}
