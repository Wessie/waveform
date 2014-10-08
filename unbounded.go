package waveform

import "image"

type WaveFormer interface {
	AddPair(Pair)
}

type UnboundedWaveForm struct {
	img                         []Pair
	samplesPerPair, samplesDone int
}

func NewUnboundedWaveForm(samplesPerPair int) *UnboundedWaveForm {
	return &UnboundedWaveForm{
		img:            make([]Pair, 1, 256),
		samplesPerPair: samplesPerPair,
		samplesDone:    0,
	}
}

func (wf *UnboundedWaveForm) AddPair(p Pair) {
	if wf.samplesDone >= wf.samplesPerPair {
		wf.img = append(wf.img, p)
		wf.samplesDone = 1
		return
	}

	p1 := wf.img[len(wf.img)-1]

	p1.Max = max(p.Max, p1.Max)
	p1.Min = min(p.Min, p1.Min)

	wf.img[len(wf.img)-1] = p1
	wf.samplesDone++
}

func max(x, y float32) float32 {
	if x > y {
		return x
	}
	return y
}

func min(x, y float32) float32 {
	if x > y {
		return y
	}
	return x
}

func (wf *UnboundedWaveForm) Image(width, height int) image.Image {
	/*samples := 550 // 22050
	img := make([]Pair, len(wf.img)/samples+1)
	for j, i := 0, 0; i < len(wf.img); j, i = j+1, i+samples {
		p := Pair{}
		for k := 0; k < samples && i+k < len(wf.img); k++ {
			p1 := wf.img[i+k]

			p.Max = max(p1.Max, p.Max)
			p.Min = min(p1.Min, p.Min)
		}
		img[j] = p
	}*/

	return &WaveForm{
		FillColor:       DefaultFillColor,
		BackgroundColor: DefaultBGColor,
		img:             wf.img,
		Width:           len(wf.img),
		Height:          height,
	}
}
