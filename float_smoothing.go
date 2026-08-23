package smoothie

// SingleExponentialSmooth applies simple exponential smoothing.
func (d *FloatDataFrame) SingleExponentialSmooth(smoothing float64) *FloatDataFrame {
	if d.Len() == 0 {
		return NewFloatDataFrame(0)
	}

	out := NewFloatDataFrame(d.Len())
	level := d.Index(0)
	out.Insert(0, level)
	for i := 1; i < d.Len(); i++ {
		level = smoothing*d.Index(i) + (1-smoothing)*level
		out.Insert(i, level)
	}
	return out
}

// DoubleExponentialSmooth applies Holt's double exponential smoothing.
func (d *FloatDataFrame) DoubleExponentialSmooth(smoothing, trendSmoothing float64) *FloatDataFrame {
	if d.Len() == 0 {
		return NewFloatDataFrame(0)
	}
	if d.Len() == 1 {
		return d.Copy()
	}

	out := NewFloatDataFrame(d.Len())
	level := d.Index(0)
	trend := d.Index(1) - d.Index(0)
	out.Insert(0, level)
	for i := 1; i < d.Len(); i++ {
		previousLevel := level
		level = smoothing*d.Index(i) + (1-smoothing)*(level+trend)
		trend = trendSmoothing*(level-previousLevel) + (1-trendSmoothing)*trend
		out.Insert(i, level)
	}
	return out
}

// DoubleSmoothPredictN smooths the dataframe and appends n trend predictions.
func (d *FloatDataFrame) DoubleSmoothPredictN(n int, smoothing, trendSmoothing float64) *FloatDataFrame {
	if n < 0 {
		panic("FloatDataFrame: prediction count cannot be negative")
	}
	if d.Len() == 0 {
		return NewFloatDataFrame(n)
	}
	if d.Len() == 1 {
		out := NewFloatDataFrame(d.Len() + n)
		out.Insert(0, d.Index(0))
		for i := 1; i < out.Len(); i++ {
			out.Insert(i, d.Index(0))
		}
		return out
	}

	smoothed := d.DoubleExponentialSmooth(smoothing, trendSmoothing)
	out := NewFloatDataFrame(d.Len() + n)
	for i := 0; i < d.Len(); i++ {
		out.Insert(i, smoothed.Index(i))
	}
	trend := smoothed.Index(d.Len()-1) - smoothed.Index(d.Len()-2)
	for i := d.Len(); i < out.Len(); i++ {
		out.Insert(i, smoothed.Index(d.Len()-1)+trend*float64(i-d.Len()+1))
	}
	return out
}

// TripleSmooth applies multiplicative Holt-Winters triple-exponential smoothing.
func (d *FloatDataFrame) TripleSmooth(smoothing, trendSmoothing, seasonalSmoothing float64, period int) *FloatDataFrame {
	if period <= 0 {
		panic("FloatDataFrame: seasonal period must be positive")
	}
	if d.Len() == 0 {
		return NewFloatDataFrame(0)
	}
	if d.Len() < period*2 {
		return d.Copy()
	}

	out := NewFloatDataFrame(d.Len())
	seasonal := make([]float64, period)
	var level float64
	for i := 0; i < period; i++ {
		level += d.Index(i)
	}
	level /= float64(period)
	if level == 0 {
		level = 1
	}
	for i := 0; i < period; i++ {
		seasonal[i] = d.Index(i) / level
		out.Insert(i, d.Index(i))
	}

	trend := 0.0
	for i := 0; i < period; i++ {
		trend += (d.Index(i+period) - d.Index(i)) / float64(period)
	}
	trend /= float64(period)

	for i := period; i < d.Len(); i++ {
		season := seasonal[i%period]
		if season == 0 {
			season = 1
		}
		previousLevel := level
		level = smoothing*(d.Index(i)/season) + (1-smoothing)*(level+trend)
		trend = trendSmoothing*(level-previousLevel) + (1-trendSmoothing)*trend
		if level != 0 {
			seasonal[i%period] = seasonalSmoothing*(d.Index(i)/level) + (1-seasonalSmoothing)*season
		}
		out.Insert(i, (level+trend)*seasonal[i%period])
	}
	return out
}
