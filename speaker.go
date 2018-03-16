package synthia

import "fmt"

// A Speaker is the output device for the synthesizer.
type Speaker struct {
	bufferSize int
	InputR     AudioStream
	InputL     AudioStream
	player     StreamOutput
	status     chan bool
}

// NewSpeaker returns an initialized speaker instance
func NewSpeaker(output StreamOutput, bufferSize int) *Speaker {

	return &Speaker{bufferSize: bufferSize, player: output}
}

func (s *Speaker) convert(rightIn, leftIn []float64, p []byte) {
	offset := 0

	inLength := len(rightIn)

	fmt.Println(rightIn)

	chans := [][]float64{leftIn, rightIn}

	for i := 0; i < inLength; i++ {

		for chanIdx := 0; chanIdx < len(chans); chanIdx++ {
			v := uint16(chans[chanIdx][i])

			buf := []byte{uint8(v & 0xff), uint8(v >> 8)}
			for j := 0; j < len(buf); j++ {
				p[offset+j] = buf[j]
			}

			offset += len(buf)

		}
	}
}

func (s *Speaker) play() {
	stpChan := s.status
	for {
		select {
		default:
			rightBuf := make([]float64, s.bufferSize/2)
			if s.InputR != nil {
				s.InputR.Stream(rightBuf)
			}

			leftBuf := make([]float64, s.bufferSize/2)
			if s.InputL != nil {
				s.InputL.Stream(leftBuf)
			}

			outBuf := make([]byte, s.bufferSize*2)
			s.convert(rightBuf, leftBuf, outBuf)

			_, err := s.player.Write(outBuf)

			// TODO: Handler properly
			if err != nil {
				panic(err)
			}
		case <-stpChan:
			return
		}
	}
}

// Start starts the speaker
func (s *Speaker) Start() {
	s.status = make(chan bool)
	go s.play()
}

// Stop stops the speaker
func (s *Speaker) Stop() {
	s.status <- true
	close(s.status)
}
