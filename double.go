package smoothie

import (
	"math"
)

// DoubleSmoothPredictPoint the Nth data point in the future
func (d *DataFrame) DoubleSmoothPredictPoint(n int, sf, tf float64) float64 {
	p := EmptyDataFrame(d.Len())
	b := EmptyDataFrame(d.Len())

	p.Insert(0, d.Index(0))
	b.Insert(0, d.Index(0))

	// so normal holt-winters for the beginning of the data frame
	for i := 0; i < d.Len(); i++ {
		p.Insert(i, d.doubleSmoothPoint(i, sf, tf, p, b))
	}

	return p.Index(d.Len()-1) + (p.bVal(p.Len()-1, sf, tf, p, b) * float64(n))
}

// DoubleSmoothPredictN return a new DataFrame filled with n predictions
func (d *DataFrame) DoubleSmoothPredictN(n int, sf, tf float64) *DataFrame {
	p := EmptyDataFrame(n)
	for i := 0; i < n; i++ {
		p.Insert(i, d.DoubleSmoothPredictPoint(i, sf, tf))
	}

	return p
}

// DoubleExponentialSmooth given a smoothing factor, apply the hold-winters double exponential smoothing algorhythm
func (d *DataFrame) DoubleExponentialSmooth(sf, tf float64) *DataFrame {
	smoothingScratch := EmptyDataFrame(d.Len())
	bValScratch := EmptyDataFrame(d.Len())

	// init 0 values
	smoothingScratch.Insert(0, d.Index(0))
	bValScratch.Insert(0, d.Index(0))

	for i := 1; i < d.Len(); i++ {
		smoothingScratch.Insert(i, d.doubleSmoothPoint(i, sf, tf, smoothingScratch, bValScratch))
	}

	return smoothingScratch
}

func (d *DataFrame) doubleSmoothPoint(i int, sf, tf float64, s, b *DataFrame) float64 {
	if i == 1 {
		return d.Index(1)
	}

	// check if the values has already been calculated
	if f := s.Index(i); !math.IsNaN(f) {
		return f
	}

	// if the value has not been calculated before, calc and return it
	return (sf * d.Index(i)) + ((1 - sf) * (d.doubleSmoothPoint(i-1, sf, tf, s, b) + d.bVal(i-1, sf, tf, s, b)))

}

func (d *DataFrame) bVal(i int, sf, tf float64, s, b *DataFrame) float64 {
	if i == 1 {
		return b.Insert(1, d.Index(1)-d.Index(0))
	}

	// check if the values has already been calculated
	if f := b.Index(i); !math.IsNaN(f) {
		return f
	}

	x := d.doubleSmoothPoint(i, sf, tf, s, b) - d.doubleSmoothPoint(i-1, sf, tf, s, b)
	y := (1 - tf) * d.bVal(i-1, sf, tf, s, b)

	return b.Insert(i, (tf*x)+y)
}
