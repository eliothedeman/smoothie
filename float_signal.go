package smoothie

import (
	"math"
	"math/cmplx"
	"math/rand"
)

// NewSignal creates a sinusoidal floating-point dataframe.
func NewSignal(length int, frequency float64) *FloatDataFrame {
	out := NewFloatDataFrame(length)
	for i := 0; i < length; i++ {
		out.Insert(i, math.Sin(float64(i)/float64(length)*math.Pi*2*frequency))
	}
	return out
}

// Noise creates a floating-point dataframe containing random values in [0, 1).
func Noise(length int) *FloatDataFrame {
	out := NewFloatDataFrame(length)
	for i := 0; i < length; i++ {
		out.Insert(i, rand.Float64())
	}
	return out
}

// FFT returns the discrete Fourier transform of the dataframe.
func (d *FloatDataFrame) FFT() []complex128 {
	result := make([]complex128, d.Len())
	if d.Len() == 0 {
		return result
	}

	for frequency := range result {
		for sample := 0; sample < d.Len(); sample++ {
			angle := -2 * math.Pi * float64(frequency*sample) / float64(d.Len())
			result[frequency] += complex(d.Index(sample), 0) * cmplx.Exp(complex(0, angle))
		}
	}
	return result
}

// FFTTopFreqs returns the indexes of frequencies whose magnitude is well
// above the dataframe's mean absolute deviation.
func (d *FloatDataFrame) FFTTopFreqs() *FloatDataFrame {
	frequencies := d.FFT()
	frequencies = frequencies[:len(frequencies)/2]
	magnitudes := make([]float64, len(frequencies))
	for i, frequency := range frequencies {
		magnitudes[i] = cmplx.Abs(frequency)
	}

	result := NewFloatDataFrame(0)
	if len(magnitudes) == 0 {
		return result
	}
	average := NewFloatDataFrameFromSlice(magnitudes).StdDev()
	for i := 1; i < len(magnitudes); i++ {
		if magnitudes[i] > average*3 {
			result.Grow(1)
			result.Push(float64(i))
		}
	}
	return result
}
