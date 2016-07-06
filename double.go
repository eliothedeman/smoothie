package smoothie

// DoubleSmoothPredictN return a new DataFrame filled with n predictions
func (d *DataFrame) DoubleSmoothPredictN(n int, sf, tf float64) *DataFrame {
	out := make([]float64, d.Len()+n)
	out[0] = d.Index(0)
	var b, x, y float64
	b = d.Index(1) - d.Index(0)
	for i := 1; i < d.Len(); i++ {
		x = sf * d.Index(i)
		b = doubleSmoothBVal(i-1, sf, tf, out, b)
		y = (1 - sf) * (out[i-1] + b)
		out[i] = x + y
	}

	last := out[d.Len()-1]
	for i := d.Len(); i < len(out); i++ {
		out[i] = last + (b * float64(i-d.Len()))

	}

	return NewDataFrameFromSlice(out)
}

// DoubleExponentialSmooth given a smoothing factor, apply the hold-winters double exponential smoothing algorhythm
func (d *DataFrame) DoubleExponentialSmooth(sf, tf float64) *DataFrame {
	out := make([]float64, d.Len())
	out[0] = d.Index(0)
	var b, x, y float64
	b = d.Index(1) - d.Index(0)
	for i := 1; i < len(out); i++ {
		x = sf * d.Index(i)
		b = doubleSmoothBVal(i-1, sf, tf, out, b)
		y = (1 - sf) * (out[i-1] + b)
		out[i] = x + y
	}

	return NewDataFrameFromSlice(out)
}

func doubleSmoothBVal(i int, sf, tf float64, s []float64, b float64) float64 {
	if i == 0 {
		return b
	}

	x := tf * (s[i] - s[i-1])
	y := (1 - tf) * b

	return x + y
}
