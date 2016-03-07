package smoothie

import "math"

func sVal(i int, smooth, trend, season float64, period int, s, b, c, d *DataFrame) float64 {
	if i < 1 {
		if i < 0 {
			i = 0
		}
		return d.Index(i)
	}

	// check to see if this value has already been cached
	if f := s.Index(i); !math.IsNaN(f) {
		return f
	}

	x := d.Index(i)
	x = (smooth * (x / cVal(i-period, smooth, trend, season, period, s, b, c, d)))
	y := (1 - smooth) * (sVal(i-1, smooth, trend, season, period, s, b, c, d) + bVal(i-1, smooth, trend, season, period, s, b, c, d))

	// cache the value
	s.Insert(i, x+y)

	return x + y
}
func bVal(i int, smooth, trend, season float64, period int, s, b, c, d *DataFrame) float64 {
	if i < 1 {
		if i < 0 {
			i = 0
		}
		return d.Index(i)
	}

	// check to see if this value has already been cached
	if f := b.Index(i); !math.IsNaN(f) {
		return f
	}
	x := trend * (sVal(i, smooth, trend, season, period, s, b, c, d) - sVal(i-1, smooth, trend, season, period, s, b, c, d))
	y := (1 - trend) * (bVal(i-1, smooth, trend, season, period, s, b, c, d))
	// cache the value
	b.Insert(i, x+y)
	return x + y
}
func cVal(i int, smooth, trend, season float64, period int, s, b, c, d *DataFrame) float64 {
	// if we haven't been trained yet, just return the raw value
	if i < 0 {
		return d.Index(i + period)
	}

	// check to see if this value has already been cached
	if f := c.Index(i); !math.IsNaN(f) {
		return f
	}

	x := season * (d.Index(i) / sVal(i, smooth, trend, season, period, s, b, c, d))
	y := (1 - season) * cVal(i-period, smooth, trend, season, period, s, b, c, d)

	c.Insert(i, x+y)
	return x + y
}

// TripleSmooth  applies holt-winters tripple-exponential smoothing
func (d *DataFrame) TripleSmooth(smooth, trend, season float64, period int) *DataFrame {
	s := EmptyDataFrame(d.Len())
	b := EmptyDataFrame(d.Len())
	c := EmptyDataFrame(d.Len())

	for i := 0; i < d.Len(); i++ {

		// cache the produced values
		c.Insert(i, cVal(i, smooth, trend, season, period, s, b, c, d))
		b.Insert(i, bVal(i, smooth, trend, season, period, s, b, c, d))
		s.Insert(i, sVal(i, smooth, trend, season, period, s, b, c, d))
	}

	return s
}
