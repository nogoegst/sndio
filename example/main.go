package main

import (
	"encoding/binary"
	"io"
	"log"
	"math"

	"github.com/nogoegst/sndio"
	"zikichombo.org/sound"
	"zikichombo.org/sound/freq"
	"zikichombo.org/sound/gen"
)

type SoundReader struct {
	snd sound.Source
}

func (sr *SoundReader) Read(p []byte) (int, error) {
	floats := make([]float64, 1)
	_, err := sr.snd.Receive(floats)
	if err != nil {
		return 0, err
	}
	sampleFloat := floats[0]
	sampleFloat *= 0.09
	sample := int16(float64(math.MaxInt16) * sampleFloat)
	binary.BigEndian.PutUint16(p[0:2], uint16(sample))
	return 2, nil
}

func main() {
	d, err := sndio.Open(sndio.AnyDevice, sndio.Play, false)
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()
	p := &sndio.Parameters{
		Bits:         16,
		LittleEndian: false,
		Signed:       true,
		PlayChans:    1,
	}
	_, err = d.SetParameters(p)
	if err != nil {
		log.Fatal(err)
	}
	if err := d.Start(); err != nil {
		log.Fatal(err)
	}

	g := gen.New(48000 * freq.Hertz)
	sine := g.Sin(300 * freq.Hertz)

	_, err = io.Copy(d, &SoundReader{snd: sine})
	if err != nil {
		log.Fatal(err)
	}

}
