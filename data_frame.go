package smoothie

import (
	"fmt"
	"math"
)

type DataFrame struct {
	pivot int
	data  []float64
}

func NewDataFrame(f []float64) *DataFrame {
	return &DataFrame{
		data: f,
	}
}

func EmptyDataFrame(size int) *DataFrame {
	return NewDataFrame(make([]float64, size))
}

type WeightingFunc func(index, length int) float64

func LinearWeighting(index, length int) float64 {
	return float64(index) / float64(length)
}

func ReverseLinearWeighting(index, length int) float64 {
	return 1 - LinearWeighting(index, length)
}

func (d *DataFrame) Weight(wf WeightingFunc) *DataFrame {
	for i := 0; i < d.Len(); i++ {
		d.Insert(i, d.Index(i)*wf(i, d.Len()))
	}
	return d
}

func (d *DataFrame) WeightedMovingAverage(windowSize int, wf WeightingFunc) *DataFrame {
	ma := NewDataFrame(make([]float64, d.Len()))

	for i := 0; i < d.Len()-windowSize; i++ {
		ma.Push(d.Slice(i, i+windowSize).Weight(wf).Avg())
	}

	return ma
}

// calculate the moving average of the dataframe
func (d *DataFrame) MovingAverage(windowSize int) *DataFrame {
	ma := NewDataFrame(make([]float64, d.Len()))
	for i := windowSize; i < d.Len()-windowSize; i++ {
		ma.Push(d.Slice(i, i+windowSize).Avg())
	}

	return ma
}

func (d *DataFrame) SingleExponentialSmooth(sf float64) *DataFrame {
	smoothed := EmptyDataFrame(d.Len())

	for i := 0; i < d.Len(); i++ {
		smoothed.Push(d.SingleSmoothPoint(i, sf))
	}

	return smoothed
}

func (d *DataFrame) SingleSmoothPoint(i int, sf float64) float64 {
	if i <= 1 {
		return (sf * d.Index(i)) + (1 - sf)
	}

	return (sf * d.Index(i)) + ((1 - sf) * d.SingleSmoothPoint(i-1, sf))
}

// given a smoothing factor, apply the hold-winters double exponential smoothing algorhythm
func (d *DataFrame) DoubleExponentialSmooth(sf, tf float64) *DataFrame {
	smoothed := NewDataFrame(make([]float64, d.Len()))

	for i := 0; i < d.Len(); i++ {
		smoothed.Push(d.DoubleSmoothPoint(i, sf, tf))
	}

	return smoothed
}

func (d *DataFrame) DoubleSmoothPoint(i int, sf, tf float64) float64 {
	if i <= 1 {
		return d.Index(i)
	}

	return (sf * d.Index(i)) + ((1 - sf) * (d.DoubleSmoothPoint(i-1, sf, tf) + d.bVal(i-1, sf, tf)))
}

func (d *DataFrame) bVal(i int, sf, tf float64) float64 {
	if i <= 1 {
		return d.Index(1) - d.Index(0)
	}

	return (tf * (d.DoubleSmoothPoint(i, sf, tf) - d.DoubleSmoothPoint(i-1, sf, tf))) + ((1 - tf) * d.bVal(i-1, sf, tf))
}

func (d *DataFrame) Copy() *DataFrame {
	dst := NewDataFrame(make([]float64, d.Len()))
	copy(dst.data, d.data)
	dst.pivot = d.pivot
	return dst
}

// return the sub slice of the data as a data frame
func (d *DataFrame) Slice(b, e int) *DataFrame {
	if b >= e {
		panic(fmt.Sprintf("Dataframe: beginning cannot be larger than end in slice operaton. Begining: %d End: %d", b, e))
	}

	if e > d.Len() {
		panic(fmt.Sprintf("DataFrame: index out of range. index: %d length: %d", e, d.Len()))
	}

	slice := make([]float64, e-b)
	for i := range slice {
		slice[i] = d.Index(e + i)
	}

	return NewDataFrame(slice)
}

func (d *DataFrame) Len() int {
	return len(d.data)
}

func (d *DataFrame) Grow(amount int) *DataFrame {

	// separate the first and second halves
	first := d.data[d.pivot:]
	last := d.data[:d.pivot]

	// make a new slice that can accomodate the new space
	d.data = make([]float64, d.Len()+amount)
	copy(d.data[:d.pivot], first)
	copy(d.data[d.pivot:amount+d.pivot], last)
	d.pivot += amount
	return d
}

func (d *DataFrame) Shrink(amount int) *DataFrame {
	if amount > d.Len() {
		panic(fmt.Sprintf("DataFrame: unable to shrink frame. amount: %d length: %d", d.Len(), amount))
	}

	newData := make([]float64, d.Len()-amount)

	for i := range newData {
		newData[i] = d.Index(i)
	}

	d.data = newData
	return d
}

func (d *DataFrame) Avg() float64 {
	var t float64
	for _, e := range d.data {
		t += e
	}

	return t / float64(d.Len())
}

// standard deviation of the data frame
func (d *DataFrame) StdDev() float64 {
	var diff float64
	avg := d.Avg()

	for _, e := range d.data {
		diff += math.Abs(avg - e)
	}

	return diff / float64(d.Len())
}

func (d *DataFrame) Push(e float64) {
	d.data[d.pivot] = e
	d.incrPivot()
}

func (d *DataFrame) Insert(i int, val float64) {
	if !d.hasIndex(i) {
		panic(fmt.Sprintf("DataFrame: index out of range. index: %d length: %d", i, d.Len()))
	}

	d.data[d.realIndex(i)] = val
}

func (d *DataFrame) Index(i int) float64 {
	if !d.hasIndex(i) {
		panic(fmt.Sprintf("DataFrame: index out of range. index: %d length: %d", i, d.Len()))
	}

	return d.data[d.realIndex(i)]
}

func (d *DataFrame) hasIndex(i int) bool {
	return (i >= 0 && i < d.Len())
}

func (d *DataFrame) Data() []float64 {
	ord := make([]float64, d.Len())

	for i := range d.data {
		ord[i] = d.Index(i)
	}

	return ord
}

// returns the value the given index is actually pointing to
func (d *DataFrame) realIndex(i int) int {
	return (d.pivot + i) % d.Len()
}

func (d *DataFrame) incrPivot() {
	d.pivot += 1
	d.pivot = d.pivot % d.Len()
}
