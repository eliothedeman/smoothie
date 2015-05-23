package smoothie

import (
	"fmt"
	"log"
	"math"

	"github.com/gonum/plot/plotter"
)

type DataFrame struct {
	pivot int
	data  []float64
}

func NewDataFrameFromSlice(f []float64) *DataFrame {
	return &DataFrame{
		data: f,
	}
}

func NewDataFrame(length int) *DataFrame {
	return &DataFrame{
		data: make([]float64, length),
	}
}

func EmptyDataFrame(size int) *DataFrame {
	df := NewDataFrame(size)
	for i := 0; i < df.Len(); i++ {
		df.Push(math.NaN())
	}

	return df
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
		d.Insert(i, 2.5*d.Index(i)*wf(i, d.Len()))
	}
	return d
}

func (d *DataFrame) WeightedMovingAverage(windowSize int, wf WeightingFunc) *DataFrame {
	ma := NewDataFrame(d.Len())

	for i := 0; i < d.Len(); i++ {
		if i+windowSize > d.Len() {
			ma.Insert(i, d.Slice(i, d.Len()).Weight(wf).Avg())
		} else {
			ma.Insert(i, d.Slice(i, i+windowSize).Weight(wf).Avg())
		}

	}

	return ma
}

// calculate the moving average of the dataframe
func (d *DataFrame) MovingAverage(windowSize int) *DataFrame {
	ma := NewDataFrame(d.Len())
	for i := 0; i < d.Len(); i++ {
		if i+windowSize > d.Len() {
			ma.Insert(i, d.Slice(i, d.Len()).Avg())
		} else {
			ma.Insert(i, d.Slice(i, i+windowSize).Avg())
		}
	}

	return ma
}

func (d *DataFrame) AddConst(f float64) *DataFrame {
	df := NewDataFrame(d.Len())
	for i := 0; i < d.Len(); i++ {
		df.Insert(i, f+d.Index(i))

	}

	return df
}

func (d *DataFrame) SubConst(f float64) *DataFrame {
	df := NewDataFrame(d.Len())
	for i := 0; i < d.Len(); i++ {
		df.Insert(i, d.Index(i)-f)

	}

	return df
}
func (d *DataFrame) MultiConst(f float64) *DataFrame {
	df := NewDataFrame(d.Len())
	for i := 0; i < d.Len(); i++ {
		df.Insert(i, f*d.Index(i))

	}

	return df
}
func (d *DataFrame) DivConst(f float64) *DataFrame {
	df := NewDataFrame(d.Len())
	for i := 0; i < d.Len(); i++ {
		df.Insert(i, d.Index(i)/f)

	}

	return df
}
func (d *DataFrame) Add(df *DataFrame) *DataFrame {
	if d.Len() != df.Len() {
		log.Panicf("Add: len %d and %d don't match", d.Len(), df.Len())
	}

	newDf := NewDataFrame(d.Len())

	for i := 0; i < d.Len(); i++ {
		newDf.Insert(i, d.Index(i)+df.Index(i))
	}

	return newDf
}

func (d *DataFrame) Sub(df *DataFrame) *DataFrame {
	if d.Len() != df.Len() {
		log.Panicf("Add: len %d and %d don't match", d.Len(), df.Len())
	}

	newDf := NewDataFrame(d.Len())

	for i := 0; i < d.Len(); i++ {
		newDf.Insert(i, d.Index(i)-df.Index(i))
	}

	return newDf
}

func (d *DataFrame) Mutli(df *DataFrame) *DataFrame {
	if d.Len() != df.Len() {
		log.Panicf("Add: len %d and %d don't match", d.Len(), df.Len())
	}

	newDf := NewDataFrame(d.Len())

	for i := 0; i < d.Len(); i++ {
		newDf.Insert(i, d.Index(i)*df.Index(i))
	}

	return newDf
}

func (d *DataFrame) Dev(df *DataFrame) *DataFrame {
	if d.Len() != df.Len() {
		log.Panicf("Add: len %d and %d don't match", d.Len(), df.Len())
	}

	newDf := NewDataFrame(d.Len())

	for i := 0; i < d.Len(); i++ {
		newDf.Insert(i, d.Index(i)/df.Index(i))
	}

	return newDf
}

func (d *DataFrame) Copy() *DataFrame {
	dst := NewDataFrame(d.Len())
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
		slice[i] = d.Index(b + i)
	}

	return NewDataFrameFromSlice(slice)
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
	var l int
	for _, e := range d.data {
		if !math.IsNaN(e) {
			t += e
			l++
		}
	}

	if l == 0 {
		return 0
	}

	return t / float64(l)
}

// standard deviation of the data frame
func (d *DataFrame) StdDev() float64 {
	var diff float64
	var l int
	avg := d.Avg()

	for _, e := range d.data {
		if !math.IsNaN(e) {
			diff += math.Abs(avg - e)
			l++

		}
	}

	if l == 0 {
		return 0
	}

	return diff / float64(l)
}

func (d *DataFrame) Push(e float64) float64 {
	d.data[d.pivot] = e
	d.incrPivot()
	return e
}

func (d *DataFrame) Insert(i int, val float64) float64 {
	if !d.hasIndex(i) {
		panic(fmt.Sprintf("DataFrame: index out of range. index: %d length: %d", i, d.Len()))
	}

	d.data[d.realIndex(i)] = val
	return val
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

func (d *DataFrame) PlotPoints() plotter.XYs {
	pts := make(plotter.XYs, d.Len())

	for i := range pts {
		pts[i].X = float64(i)
		pts[i].Y = d.Index(i)
	}

	return pts
}
