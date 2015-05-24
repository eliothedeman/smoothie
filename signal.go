package smoothie

import (
	"math"
	"math/cmplx"
	"math/rand"

	"github.com/mjibson/go-dsp/fft"
)

// create a new data frame filled with a signal given a frequencey and a phase
func NewSignal(length int, freq float64) *DataFrame {
	df := NewDataFrame(length)

	for i := 0; i < length; i++ {
		df.Insert(i, math.Sin(float64(i)/float64(df.Len())*math.Pi*2*freq))
	}

	return df
}

func Noise(length int) *DataFrame {
	df := NewDataFrame(length)

	for i := 0; i < length; i++ {
		df.Insert(i, rand.Float64())
	}

	return df
}

// Retusn a new dataframe populated with the real fft values from the current dataframe
func (d *DataFrame) FFT() *DataFrame {
	data := d.Data()

	f := fft.FFTReal(data)
	freqs := NewDataFrame(d.Len())
	for i := 0; i < d.Len(); i++ {
		freqs.Push(cmplx.Abs(f[i]))
	}

	// cut out imposible freqs
	freqs = freqs.Slice(0, freqs.Len()/2)

	stdDev := freqs.StdDev()
	top := NewDataFrame(0)
	for i := 0; i < freqs.Len(); i++ {
		if freqs.Index(i) > stdDev*5 {
			top.Grow(1)
			top.Push(float64(i))
		}
	}

	return top
}
