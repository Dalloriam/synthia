package modular

import "github.com/dalloriam/synthia/synthia"

type Knob struct {
	Line  synthia.AudioStream
	value float64
}

func NewKnob(defaultVal float64) *Knob {
	k := &Knob{
		Line:  nil,
		value: defaultVal,
	}
	return k
}

func (k *Knob) SetValue(val float64) {
	k.value = val
}

func (k *Knob) Stream(p []float64) (int, error) {
	if k.Line == nil {
		for i := 0; i < len(p); i++ {
			p[i] = k.value
		}
		return len(p), nil
	}
	return k.Line.Stream(p)
}
