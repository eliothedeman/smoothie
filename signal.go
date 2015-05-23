package smoothie

import (
	"math"
	"math/rand"
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
