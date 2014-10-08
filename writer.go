package waveform

import (
	"image"
	"unsafe"
)

// NewWriter creates a Writer filled with unbounded waveforms. If you want
// to work with a different WaveFormer, construct the Writer yourself.
//
// The amount is given by channels, and the amount of samples per resulting
// pixel is given by samplesPerPair.
func NewWriter(channels, samplesPerPair int) *Writer {
	ch := make([]WaveFormer, channels)
	for i := range ch {
		ch[i] = NewUnboundedWaveForm(samplesPerPair)
	}

	return &Writer{
		WaveForms: ch,
	}
}

// Writer is an io.Writer that wraps around one or more WaveFormer.
//
// See the Write method for how it interprets data that is written.
type Writer struct {
	WaveForms []WaveFormer
}

// Image returns an image.Image that represents all WaveForms
// stacked vertically.
func (wf *Writer) Image(width, height int) image.Image {
	var imgs []image.Image

	for i := range wf.WaveForms {
		former := wf.WaveForms[i]
		i, ok := former.(image.Image)
		if ok {
			imgs = append(imgs, i)
		}

		type imager interface {
			Image(width, height int) image.Image
		}

		ir, ok := former.(imager)
		if ok {
			imgs = append(imgs, ir.Image(width, height))
		}
	}

	return NewStack(imgs...)
}

// Write adds len(p)/len(wf.WaveForms) amount of pairs to each WaveFormer
// in the WaveForms slice.
//
// Write works by interpreting p as audio data that has len(wf.WaveForms)
// amount of channels interleaved together.
//
// NOTE: Make sure to properly cutoff the end of your p on the end of samples,
// p should basically adhere to len(p) % (sampleSize * len(wf.WaveForms)) == 0
func (wf *Writer) Write(p []byte) (n int, err error) {
	// assume float32 and interleaved channels for now
	floats := *(*[]float32)(unsafe.Pointer(&p))
	floats = floats[0 : len(p)/4 : len(p)/4]

	for f := 0; f < len(floats); {
		for i := range wf.WaveForms {
			wf.WaveForms[i].AddPair(Pair{floats[f], floats[f]})
			f++
		}
	}

	return len(p), nil
}
