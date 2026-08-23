package smoothie

import (
	"fmt"
	"log"
	"sort"
)

type Num interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// A ring buffer implementation optimized for moving window statistic operations.
type DataFrame[T Num] struct {
	data  []T
	pivot int
}

// Create a new dataframe from an existing slice.
func NewDataFrameFromSlice[T Num](data []T) *DataFrame[T] {
	return &DataFrame[T]{data: data}
}

// Create a new zero-valued dataframe.
func NewDataFrame[T Num](length int) *DataFrame[T] {
	return &DataFrame[T]{data: make([]T, length)}
}

// Add a value to every value of the dataframe.
func (d *DataFrame[T]) AddConst(value T) *DataFrame[T] {
	out := NewDataFrame[T](d.Len())
	for i := 0; i < d.Len(); i++ {
		out.Insert(i, d.Index(i)+value)
	}
	return out
}

// Subtract a value from every value of the dataframe.
func (d *DataFrame[T]) SubConst(value T) *DataFrame[T] {
	out := NewDataFrame[T](d.Len())
	for i := 0; i < d.Len(); i++ {
		out.Insert(i, d.Index(i)-value)
	}
	return out
}

// Multiply every value of the dataframe by a value.
func (d *DataFrame[T]) MultiConst(value T) *DataFrame[T] {
	out := NewDataFrame[T](d.Len())
	for i := 0; i < d.Len(); i++ {
		out.Insert(i, d.Index(i)*value)
	}
	return out
}

// Divide every value of the dataframe by a value.
func (d *DataFrame[T]) DivConst(value T) *DataFrame[T] {
	out := NewDataFrame[T](d.Len())
	for i := 0; i < d.Len(); i++ {
		out.Insert(i, d.Index(i)/value)
	}
	return out
}

// Add two dataframes together.
func (d *DataFrame[T]) Add(other *DataFrame[T]) *DataFrame[T] {
	d.requireSameLength("Add", other)
	out := NewDataFrame[T](d.Len())
	for i := 0; i < d.Len(); i++ {
		out.Insert(i, d.Index(i)+other.Index(i))
	}
	return out
}

// Shift the pivot of the dataframe.
func (d *DataFrame[T]) Shift(offset int) *DataFrame[T] {
	out := d.Copy()
	if out.Len() == 0 {
		return out
	}
	out.pivot = (out.pivot + offset) % out.Len()
	if out.pivot < 0 {
		out.pivot += out.Len()
	}
	return out
}

// Subtract a dataframe by another dataframe.
func (d *DataFrame[T]) Sub(other *DataFrame[T]) *DataFrame[T] {
	d.requireSameLength("Sub", other)
	out := NewDataFrame[T](d.Len())
	for i := 0; i < d.Len(); i++ {
		out.Insert(i, d.Index(i)-other.Index(i))
	}
	return out
}

// Mutli multiplies two dataframes together.
func (d *DataFrame[T]) Mutli(other *DataFrame[T]) *DataFrame[T] {
	d.requireSameLength("Mutli", other)
	out := NewDataFrame[T](d.Len())
	for i := 0; i < d.Len(); i++ {
		out.Insert(i, d.Index(i)*other.Index(i))
	}
	return out
}

// Div divides a dataframe by another dataframe.
func (d *DataFrame[T]) Div(other *DataFrame[T]) *DataFrame[T] {
	d.requireSameLength("Div", other)
	out := NewDataFrame[T](d.Len())
	for i := 0; i < d.Len(); i++ {
		out.Insert(i, d.Index(i)/other.Index(i))
	}
	return out
}

// Derivative returns the difference between adjacent values.
func (d *DataFrame[T]) Derivative() *DataFrame[T] {
	if d.Len() == 0 {
		return NewDataFrame[T](0)
	}

	out := NewDataFrame[T](d.Len() - 1)
	for i := 1; i < d.Len(); i++ {
		out.Insert(i-1, d.Index(i)-d.Index(i-1))
	}
	return out
}

// Copy returns a copy of the dataframe.
func (d *DataFrame[T]) Copy() *DataFrame[T] {
	out := NewDataFrame[T](d.Len())
	copy(out.data, d.data)
	out.pivot = d.pivot
	return out
}

// Slice returns the values in the half-open range [begin, end].
func (d *DataFrame[T]) Slice(begin, end int) *DataFrame[T] {
	if begin >= end {
		panic(fmt.Sprintf("DataFrame: beginning cannot be larger than end in slice operation. Beginning: %d End: %d", begin, end))
	}
	if begin < 0 || end > d.Len() {
		panic(fmt.Sprintf("DataFrame: index out of range. index: %d length: %d", end, d.Len()))
	}

	values := make([]T, end-begin)
	for i := range values {
		values[i] = d.Index(begin + i)
	}
	return NewDataFrameFromSlice(values)
}

// Len returns the length of the dataframe.
func (d *DataFrame[T]) Len() int {
	return len(d.data)
}

// Grow increases the dataframe by a given amount with zero values.
func (d *DataFrame[T]) Grow(amount int) *DataFrame[T] {
	if amount < 0 {
		panic("DataFrame: grow amount cannot be negative")
	}
	if amount == 0 {
		return d
	}

	data := make([]T, amount+d.Len())
	copy(data[amount:], d.Data())
	d.data = data
	d.pivot = 0
	return d
}

// Shrink removes values from the dataframe.
func (d *DataFrame[T]) Shrink(amount int) *DataFrame[T] {
	if amount < 0 {
		panic("DataFrame: shrink amount cannot be negative")
	}
	if amount >= d.Len() {
		d.data = nil
		d.pivot = 0
		return d
	}

	values := make([]T, d.Len()-amount)
	for i := range values {
		values[i] = d.Index(i)
	}
	d.data = values
	d.pivot = 0
	return d
}

// Min returns the minimum value of the dataframe. It returns the zero value
// of T for an empty dataframe.
func (d *DataFrame[T]) Min() T {
	if d.Len() == 0 {
		var zero T
		return zero
	}

	minimum := d.Index(0)
	for i := 1; i < d.Len(); i++ {
		if value := d.Index(i); value < minimum {
			minimum = value
		}
	}
	return minimum
}

// Max returns the maximum value of the dataframe. It returns the zero value
// of T for an empty dataframe.
func (d *DataFrame[T]) Max() T {
	if d.Len() == 0 {
		var zero T
		return zero
	}

	maximum := d.Index(0)
	for i := 1; i < d.Len(); i++ {
		if value := d.Index(i); value > maximum {
			maximum = value
		}
	}
	return maximum
}

// Sum returns the sum of the values in the dataframe.
func (d *DataFrame[T]) Sum() T {
	var total T
	for i := 0; i < d.Len(); i++ {
		total += d.Index(i)
	}
	return total
}

func (d *DataFrame[T]) ForEach(f func(data T, index int)) {
	for i := 0; i < d.Len(); i++ {
		f(d.Index(i), i)
	}
}

// Cat concatenates two dataframes.
func (d *DataFrame[T]) Cat(other *DataFrame[T]) *DataFrame[T] {
	out := NewDataFrame[T](d.Len() + other.Len())
	for i := 0; i < d.Len(); i++ {
		out.Push(d.Index(i))
	}
	for i := 0; i < other.Len(); i++ {
		out.Push(other.Index(i))
	}
	return out
}

// Sort returns a sorted version of the dataframe.
func (d *DataFrame[T]) Sort() *DataFrame[T] {
	values := d.Data()
	sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })
	return NewDataFrameFromSlice(values)
}

// Reverse returns a reversed version of the dataframe.
func (d *DataFrame[T]) Reverse() *DataFrame[T] {
	out := NewDataFrame[T](d.Len())
	for i := d.Len() - 1; i >= 0; i-- {
		out.Push(d.Index(i))
	}
	return out
}

// Median returns the middle value of a sorted dataframe.
func (d *DataFrame[T]) Median() T {
	if d.Len() == 0 {
		var zero T
		return zero
	}
	return d.Sort().Index(d.Len() / 2)
}

// Push adds an element to the dataframe and advances the ring-buffer pivot.
func (d *DataFrame[T]) Push(value T) T {
	if d.Len() == 0 {
		panic("DataFrame: cannot push to an empty dataframe")
	}
	d.data[d.pivot] = value
	d.incrPivot()
	return value
}

// Insert inserts an element at a specific logical index.
func (d *DataFrame[T]) Insert(index int, value T) T {
	if !d.hasIndex(index) {
		panic(fmt.Sprintf("DataFrame: index out of range. index: %d length: %d", index, d.Len()))
	}
	d.data[d.realIndex(index)] = value
	return value
}

// Index returns the value at a logical index.
func (d *DataFrame[T]) Index(index int) T {
	if !d.hasIndex(index) {
		panic(fmt.Sprintf("DataFrame: index out of range. index: %d length: %d", index, d.Len()))
	}
	return d.data[d.realIndex(index)]
}

func (d *DataFrame[T]) hasIndex(index int) bool {
	return index >= 0 && index < d.Len()
}

// Data returns the contents of a dataframe flattened into a slice.
func (d *DataFrame[T]) Data() []T {
	values := make([]T, d.Len())
	for i := range values {
		values[i] = d.Index(i)
	}
	return values
}

func (d *DataFrame[T]) requireSameLength(operation string, other *DataFrame[T]) {
	if d.Len() != other.Len() {
		log.Panicf("%s: len %d and %d don't match", operation, d.Len(), other.Len())
	}
}

func (d *DataFrame[T]) realIndex(index int) int {
	return (d.pivot + index) % d.Len()
}

func (d *DataFrame[T]) incrPivot() {
	d.pivot = (d.pivot + 1) % d.Len()
}
