package core_test

import (
	"testing"

	"github.com/dalloriam/synthia/core"
)

func TestNewMixerChannel(t *testing.T) {
	t.Run("initializes volume knob", func(t *testing.T) {
		c := core.NewMixerChannel()

		if c.Volume == nil {
			t.Error("newmixerchannel did not initialize the volume knob")
		}
	})

	t.Run("initializes input to nil", func(t *testing.T) {
		c := core.NewMixerChannel()

		if c.Input != nil {
			t.Error("newmixerchannel initialized the input line")
		}
	})
}
