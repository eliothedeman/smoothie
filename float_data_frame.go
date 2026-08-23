package smoothie

import "math"

// FloatDataFrame wraps a float64 DataFrame and provides operations whose
// results cannot be represented by every numeric DataFrame type.
type FloatDataFrame struct {
	DataFrame[float64]
}

// NewFloatDataFrame creates a zero-valued floating-point dataframe.
func NewFloatDataFrame(length int) *FloatDataFrame {
	return &FloatDataFrame{DataFrame: *NewDataFrame[float64](length)}
}

// NewFloatDataFrameFromSlice creates a floating-point dataframe from a slice.
func NewFloatDataFrameFromSlice(data []float64) *FloatDataFrame {
	return &FloatDataFrame{DataFrame: *NewDataFrameFromSlice(data)}
}

// NewFloatDataFrameFromDataFrame wraps an existing float64 dataframe.
func NewFloatDataFrameFromDataFrame(data *DataFrame[float64]) *FloatDataFrame {
	return &FloatDataFrame{DataFrame: *data}
}

// EmptyFloatDataFrame creates a floating-point dataframe initialized with NaN.
func EmptyFloatDataFrame(size int) *FloatDataFrame {
	data := make([]float64, size)
	for i := range data {
		data[i] = math.NaN()
	}
	return NewFloatDataFrameFromSlice(data)
}

func floatFrame(data *DataFrame[float64]) *FloatDataFrame {
	return NewFloatDataFrameFromDataFrame(data)
}

func (d *FloatDataFrame) Add(other *FloatDataFrame) *FloatDataFrame {
	return floatFrame(d.DataFrame.Add(&other.DataFrame))
}

func (d *FloatDataFrame) AddConst(value float64) *FloatDataFrame {
	return floatFrame(d.DataFrame.AddConst(value))
}

func (d *FloatDataFrame) Sub(other *FloatDataFrame) *FloatDataFrame {
	return floatFrame(d.DataFrame.Sub(&other.DataFrame))
}

func (d *FloatDataFrame) SubConst(value float64) *FloatDataFrame {
	return floatFrame(d.DataFrame.SubConst(value))
}

func (d *FloatDataFrame) MultiConst(value float64) *FloatDataFrame {
	return floatFrame(d.DataFrame.MultiConst(value))
}

func (d *FloatDataFrame) DivConst(value float64) *FloatDataFrame {
	return floatFrame(d.DataFrame.DivConst(value))
}

func (d *FloatDataFrame) Mutli(other *FloatDataFrame) *FloatDataFrame {
	return floatFrame(d.DataFrame.Mutli(&other.DataFrame))
}

func (d *FloatDataFrame) Div(other *FloatDataFrame) *FloatDataFrame {
	return floatFrame(d.DataFrame.Div(&other.DataFrame))
}

func (d *FloatDataFrame) Shift(offset int) *FloatDataFrame {
	return floatFrame(d.DataFrame.Shift(offset))
}

func (d *FloatDataFrame) Derivative() *FloatDataFrame {
	return floatFrame(d.DataFrame.Derivative())
}

func (d *FloatDataFrame) Copy() *FloatDataFrame {
	return floatFrame(d.DataFrame.Copy())
}

func (d *FloatDataFrame) Slice(begin, end int) *FloatDataFrame {
	return floatFrame(d.DataFrame.Slice(begin, end))
}

func (d *FloatDataFrame) Grow(amount int) *FloatDataFrame {
	d.DataFrame.Grow(amount)
	return d
}

func (d *FloatDataFrame) Shrink(amount int) *FloatDataFrame {
	d.DataFrame.Shrink(amount)
	return d
}

func (d *FloatDataFrame) Cat(other *FloatDataFrame) *FloatDataFrame {
	return floatFrame(d.DataFrame.Cat(&other.DataFrame))
}

func (d *FloatDataFrame) Sort() *FloatDataFrame {
	return floatFrame(d.DataFrame.Sort())
}

func (d *FloatDataFrame) Reverse() *FloatDataFrame {
	return floatFrame(d.DataFrame.Reverse())
}

// Avg returns the arithmetic mean, ignoring NaN values.
func (d *FloatDataFrame) Avg() float64 {
	var total float64
	var count int
	for i := 0; i < d.Len(); i++ {
		value := d.Index(i)
		if !math.IsNaN(value) {
			total += value
			count++
		}
	}
	if count == 0 {
		return math.NaN()
	}
	return total / float64(count)
}

// StdDev returns the mean absolute deviation from the mean.
func (d *FloatDataFrame) StdDev() float64 {
	average := d.Avg()
	if math.IsNaN(average) {
		return math.NaN()
	}

	var difference float64
	var count int
	for i := 0; i < d.Len(); i++ {
		value := d.Index(i)
		if !math.IsNaN(value) {
			difference += math.Abs(average - value)
			count++
		}
	}
	return difference / float64(count)
}

// Median returns the middle value, ignoring NaN values.
func (d *FloatDataFrame) Median() float64 {
	values := make([]float64, 0, d.Len())
	for i := 0; i < d.Len(); i++ {
		if value := d.Index(i); !math.IsNaN(value) {
			values = append(values, value)
		}
	}
	if len(values) == 0 {
		return math.NaN()
	}
	return NewFloatDataFrameFromSlice(values).Sort().Index(len(values) / 2)
}

// Reduce groups values into newSize buckets and averages each bucket.
func (d *FloatDataFrame) Reduce(newSize int) *FloatDataFrame {
	if newSize <= 0 {
		panic("FloatDataFrame: reduce size must be positive")
	}

	out := NewFloatDataFrame(newSize)
	for i := 0; i < newSize; i++ {
		begin := i * d.Len() / newSize
		end := (i + 1) * d.Len() / newSize
		if begin == end {
			out.Insert(i, math.NaN())
			continue
		}
		out.Insert(i, d.Slice(begin, end).Avg())
	}
	return out
}

// WeightingFunc returns a weight for an index in a dataframe.
type WeightingFunc func(index, length int) float64

// LinearWeighting gives later values a higher weight.
func LinearWeighting(index, length int) float64 {
	return float64(index) / float64(length)
}

// ReverseLinearWeighting gives earlier values a higher weight.
func ReverseLinearWeighting(index, length int) float64 {
	return 1 - LinearWeighting(index, length)
}

// Weight applies a weighting function in place.
func (d *FloatDataFrame) Weight(weighting WeightingFunc) *FloatDataFrame {
	for i := 0; i < d.Len(); i++ {
		d.Insert(i, 2.5*d.Index(i)*weighting(i, d.Len()))
	}
	return d
}

// MovingAverage calculates the trailing moving average.
func (d *FloatDataFrame) MovingAverage(windowSize int) *FloatDataFrame {
	if windowSize <= 0 {
		panic("FloatDataFrame: moving-average window must be positive")
	}

	out := NewFloatDataFrame(d.Len())
	for i := 0; i < d.Len(); i++ {
		end := i + windowSize
		if end > d.Len() {
			end = d.Len()
		}
		out.Insert(i, d.Slice(i, end).Avg())
	}
	return out
}

// WeightedMovingAverage calculates a weighted trailing moving average.
func (d *FloatDataFrame) WeightedMovingAverage(windowSize int, weighting WeightingFunc) *FloatDataFrame {
	if windowSize <= 0 {
		panic("FloatDataFrame: weighted moving-average window must be positive")
	}

	out := NewFloatDataFrame(d.Len())
	for i := 0; i < d.Len(); i++ {
		end := i + windowSize
		if end > d.Len() {
			end = d.Len()
		}
		var total, weights float64
		for index := i; index < end; index++ {
			value := d.Index(index)
			if math.IsNaN(value) {
				continue
			}
			weight := weighting(index-i, end-i)
			total += value * weight
			weights += weight
		}
		if weights == 0 {
			out.Insert(i, math.NaN())
		} else {
			out.Insert(i, total/weights)
		}
	}
	return out
}
