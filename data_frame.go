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

func (d *DataFrame) Len() int {
	return len(d.data)
}

func (d *DataFrame) Cap() int {
	return cap(d.data)
}

func (d *DataFrame) Avg() float64 {
	var t float64
	for _, e := range d.data {
		t += e
	}

	return t / float64(d.Len())
}

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
